package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/hasher"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMember メンバー管理サービス。
type ManageMember struct {
	DB      store.Store
	Hash    hasher.Hash
	Clocker clock.Clock
	Storage storage.Storage
}

// UpdateMember メンバーを更新する。
func (m *ManageMember) UpdateMember(
	ctx context.Context,
	id uuid.UUID,
	email,
	name string,
	firstName,
	lastName entity.String,
	profileImageID entity.UUID,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	// profile image check
	if profileImageID.Valid {
		image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, profileImageID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetProfileImages)
			}
			return entity.Member{}, fmt.Errorf("failed to find attachable item by id: %w", err)
		}
		if image.AttachableItem.OwnerID.Valid && image.AttachableItem.OwnerID.Bytes != id {
			return entity.Member{}, errhandle.NewCommonError(response.NotFileOwner, nil)
		}
	}
	e, err = m.DB.FindMemberByIDWithSd(ctx, sd, id)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	if name != e.Name {
		org, err := m.DB.FindOrganizationByIDWithSd(ctx, sd, e.PersonalOrganizationID)
		if err != nil {
			return entity.Member{}, fmt.Errorf("failed to find organization by id: %w", err)
		}
		_, err = m.DB.UpdateOrganizationWithSd(
			ctx, sd, e.PersonalOrganizationID, parameter.UpdateOrganizationParams{
				Name:        fmt.Sprintf("%s(personal)", name),
				Description: entity.String{Valid: true, String: fmt.Sprintf("%s (personal organization)", name)},
				Color:       org.Color,
			})
		if err != nil {
			return entity.Member{}, fmt.Errorf("failed to update organization: %w", err)
		}
	}
	p := parameter.UpdateMemberParams{
		Email:          email,
		Name:           name,
		FirstName:      firstName,
		LastName:       lastName,
		ProfileImageID: profileImageID,
	}
	e, err = m.DB.UpdateMemberWithSd(ctx, sd, id, p)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member: %w", err)
	}
	return e, nil
}

// DeleteMember メンバーを削除する。
func (m *ManageMember) DeleteMember(ctx context.Context, id uuid.UUID) (c int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	ec, err := m.DB.FindMemberWithDetailWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find member with detail: %w", err)
	}
	if ec.ProfileImageID.Valid {
		_, err = pluralDeleteImages(
			ctx,
			sd,
			m.DB,
			m.Storage,
			[]uuid.UUID{ec.ProfileImageID.Bytes},
			entity.UUID{
				Valid: true,
				Bytes: id,
			})
		if err != nil {
			return 0, fmt.Errorf("failed to delete images: %w", err)
		}
	}
	if ec.Student.Valid {
		_, err = m.DB.DeleteStudentWithSd(ctx, sd, ec.Student.Entity.StudentID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete student: %w", err)
		}
	}
	if ec.Professor.Valid {
		_, err = m.DB.DeleteProfessorWithSd(ctx, sd, ec.Professor.Entity.ProfessorID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete professor: %w", err)
		}
	}
	e, err := m.DB.FindMemberWithPersonalOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find member with personal organization: %w", err)
	}
	_, err = m.DB.DisbelongChatRoomOnMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on member: %w", err)
	}
	_, err = m.DB.DisbelongOrganizationOnMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong organization on member: %w", err)
	}
	c, err = m.DB.DeleteMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete member: %w", err)
	}
	_, err = m.DB.DeleteOrganizationWithSd(ctx, sd, e.PersonalOrganization.OrganizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	return c, nil
}

// UpdateMemberPassword メンバーのパスワードを更新する。
func (m *ManageMember) UpdateMemberPassword(
	ctx context.Context,
	id uuid.UUID,
	rawPassword string,
) (entity.Member, error) {
	hashPass, err := m.Hash.Encrypt(rawPassword)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to encrypt password: %w", err)
	}
	e, err := m.DB.UpdateMemberPassword(ctx, id, hashPass)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member password: %w", err)
	}
	return e, nil
}

// UpdateMemberRole メンバーのロールを更新する。
func (m *ManageMember) UpdateMemberRole(
	ctx context.Context,
	id uuid.UUID,
	roleID entity.UUID,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	// role check
	if roleID.Valid {
		_, err = m.DB.FindRoleByIDWithSd(ctx, sd, roleID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetRoles)
			}
			return entity.Member{}, fmt.Errorf("failed to find role by id: %w", err)
		}
	}
	e, err = m.DB.UpdateMemberRoleWithSd(ctx, sd, id, roleID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member role: %w", err)
	}
	return e, nil
}

// UpdateMemberLoginID メンバーのログインIDを更新する。
func (m *ManageMember) UpdateMemberLoginID(
	ctx context.Context,
	id uuid.UUID,
	loginID string,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	// loginID check
	_, err = m.DB.FindMemberByLoginIDWithSd(ctx, sd, loginID)
	if err == nil {
		return entity.Member{}, errhandle.NewModelDuplicatedError(MemberTargetLoginID)
	}
	e, err = m.DB.UpdateMemberLoginIDWithSd(ctx, sd, id, loginID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member login id: %w", err)
	}
	return e, nil
}

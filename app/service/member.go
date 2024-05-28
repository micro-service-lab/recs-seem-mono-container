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
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMember メンバー管理サービス。
type ManageMember struct {
	DB      store.Store
	Hash    hasher.Hash
	Clocker clock.Clock
}

// CreateMember メンバーを作成する。
func (m *ManageMember) CreateMember(
	ctx context.Context,
	loginID,
	rawPassword,
	email,
	name string,
	firstName,
	lastName entity.String,
	gradeID,
	groupID uuid.UUID,
	profileImageID,
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
	// grade check
	_, err = m.DB.FindGradeByIDWithSd(ctx, sd, gradeID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetGrades)
		}
		return entity.Member{}, fmt.Errorf("failed to find grade by id: %w", err)
	}
	// group check
	_, err = m.DB.FindGroupByIDWithSd(ctx, sd, groupID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetGroups)
		}
		return entity.Member{}, fmt.Errorf("failed to find group by id: %w", err)
	}
	// profile image check
	if profileImageID.Valid {
		_, err = m.DB.FindImageByIDWithSd(ctx, sd, profileImageID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetProfileImages)
			}
			return entity.Member{}, fmt.Errorf("failed to find attachable item by id: %w", err)
		}
	}
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
	defaultAttendStatus, err := m.DB.FindAttendStatusByKeyWithSd(ctx, sd, string(DefaultAttendStatusKey))
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find default attend status: %w", err)
	}
	col := entity.String{
		Valid:  true,
		String: RandomColor(),
	}
	org, err := m.DB.CreateOrganizationWithSd(ctx, sd, parameter.CreateOrganizationParam{
		Name: fmt.Sprintf("%s(personal)", name),
		Description: entity.String{
			Valid:  true,
			String: fmt.Sprintf("%s (personal organization)", name),
		},
		Color:      col,
		IsPersonal: true,
		IsWhole:    false,
		ChatRoomID: entity.UUID{},
	})
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to create organization: %w", err)
	}
	hashPass, err := m.Hash.Encrypt(rawPassword)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to encrypt password: %w", err)
	}
	p := parameter.CreateMemberParam{
		LoginID:                loginID,
		Password:               hashPass,
		Email:                  email,
		Name:                   name,
		FirstName:              firstName,
		LastName:               lastName,
		AttendStatusID:         defaultAttendStatus.AttendStatusID,
		GradeID:                gradeID,
		GroupID:                groupID,
		ProfileImageID:         profileImageID,
		RoleID:                 roleID,
		PersonalOrganizationID: org.OrganizationID,
	}
	e, err = m.DB.CreateMemberWithSd(ctx, sd, p)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to create organization: %w", err)
	}
	wholeOrg, err := m.DB.FindWholeOrganization(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	groupOrg, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, groupID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find group with organization: %w", err)
	}
	gradeOrg, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, gradeID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find grade with organization: %w", err)
	}
	// add some organization, chat room
	now := m.Clocker.Now()
	bop := []parameter.BelongOrganizationParam{
		{
			MemberID:       e.MemberID,
			OrganizationID: wholeOrg.OrganizationID,
			WorkPositionID: entity.UUID{},
			AddedAt:        now,
		}, {
			MemberID:       e.MemberID,
			OrganizationID: groupOrg.Organization.OrganizationID,
			WorkPositionID: entity.UUID{},
			AddedAt:        now,
		}, {
			MemberID:       e.MemberID,
			OrganizationID: gradeOrg.Organization.OrganizationID,
			WorkPositionID: entity.UUID{},
			AddedAt:        now,
		}, {
			MemberID:       e.MemberID,
			OrganizationID: org.OrganizationID,
			WorkPositionID: entity.UUID{},
			AddedAt:        now,
		},
	}
	bcrp := make([]parameter.BelongChatRoomParam, 0, len(bop))
	if groupOrg.Organization.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: groupOrg.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
	}
	if gradeOrg.Organization.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: gradeOrg.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
	}
	if wholeOrg.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: wholeOrg.ChatRoomID.Bytes,
			AddedAt:    now,
		})
	}
	if org.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: org.ChatRoomID.Bytes,
			AddedAt:    now,
		})
	}

	_, err = m.DB.BelongOrganizationsWithSd(ctx, sd, bop)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	_, err = m.DB.BelongChatRoomsWithSd(ctx, sd, bcrp)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}
	return e, nil
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
		_, err = m.DB.FindImageByIDWithSd(ctx, sd, profileImageID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetProfileImages)
			}
			return entity.Member{}, fmt.Errorf("failed to find attachable item by id: %w", err)
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
func (m *ManageMember) DeleteMember(ctx context.Context, id uuid.UUID) (int64, error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rerr)
		}
	}()
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
	c, err := m.DB.DeleteMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete member: %w", err)
	}
	_, err = m.DB.DeleteOrganizationWithSd(ctx, sd, e.PersonalOrganization.ChatRoomID.Bytes)
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
		if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rerr)
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
	e, err = m.DB.UpdateMemberLoginID(ctx, id, loginID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member login id: %w", err)
	}
	return e, nil
}

// UpdateMemberGrade メンバーの学年を更新する。
func (m *ManageMember) UpdateMemberGrade(
	ctx context.Context,
	id uuid.UUID,
	gradeID uuid.UUID,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rerr)
		}
	}()
	// grade check
	_, err = m.DB.FindGradeByIDWithSd(ctx, sd, gradeID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetGrades)
		}
		return entity.Member{}, fmt.Errorf("failed to find grade by id: %w", err)
	}
	e, err = m.DB.UpdateMemberGradeWithSd(ctx, sd, id, gradeID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member grade: %w", err)
	}
	return e, nil
}

// UpdateMemberGroup メンバーのグループを更新する。
func (m *ManageMember) UpdateMemberGroup(
	ctx context.Context,
	id uuid.UUID,
	groupID uuid.UUID,
) (e entity.Member, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rerr)
		}
	}()
	// group check
	_, err = m.DB.FindGroupByIDWithSd(ctx, sd, groupID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Member{}, errhandle.NewModelNotFoundError(MemberTargetGroups)
		}
		return entity.Member{}, fmt.Errorf("failed to find group by id: %w", err)
	}
	e, err = m.DB.UpdateMemberGroupWithSd(ctx, sd, id, groupID)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to update member group: %w", err)
	}
	return e, nil
}

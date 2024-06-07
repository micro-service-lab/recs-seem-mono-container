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

// CreateMember 生徒を作成する。
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
	roleID entity.UUID,
) (e entity.MemberWithDetail, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	now := m.Clocker.Now()
	// loginID check
	_, err = m.DB.FindMemberByLoginIDWithSd(ctx, sd, loginID)
	if err == nil {
		return entity.MemberWithDetail{}, errhandle.NewModelDuplicatedError(MemberTargetLoginID)
	}
	// grade check
	grade, err := m.DB.FindGradeByIDWithSd(ctx, sd, gradeID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.MemberWithDetail{}, errhandle.NewModelNotFoundError(MemberTargetGrades)
		}
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find grade by id: %w", err)
	}
	// group check
	group, err := m.DB.FindGroupByIDWithSd(ctx, sd, groupID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.MemberWithDetail{}, errhandle.NewModelNotFoundError(MemberTargetGroups)
		}
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find group by id: %w", err)
	}
	if grade.Key != string(GradeKeyProfessor) && group.Key == string(GroupKeyProfessor) {
		e := errhandle.NewCommonError(response.MustBeGradeProfessorIfGroupProfessor, nil)
		return entity.MemberWithDetail{}, e
	}
	if grade.Key == string(GradeKeyProfessor) && group.Key != string(GroupKeyProfessor) {
		e := errhandle.NewCommonError(response.MustBeGroupProfessorIfGradeProfessor, nil)
		return entity.MemberWithDetail{}, e
	}
	// role check
	if roleID.Valid {
		_, err = m.DB.FindRoleByIDWithSd(ctx, sd, roleID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.MemberWithDetail{}, errhandle.NewModelNotFoundError(MemberTargetRoles)
			}
			return entity.MemberWithDetail{}, fmt.Errorf("failed to find role by id: %w", err)
		}
	}
	defaultAttendStatus, err := m.DB.FindAttendStatusByKeyWithSd(ctx, sd, string(DefaultAttendStatusKey))
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find default attend status: %w", err)
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
		return entity.MemberWithDetail{}, fmt.Errorf("failed to create organization: %w", err)
	}
	hashPass, err := m.Hash.Encrypt(rawPassword)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to encrypt password: %w", err)
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
		RoleID:                 roleID,
		PersonalOrganizationID: org.OrganizationID,
	}
	member, err := m.DB.CreateMemberWithSd(ctx, sd, p)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to create organization: %w", err)
	}

	var professor entity.NullableEntity[entity.Professor]
	var student entity.NullableEntity[entity.Student]

	if grade.Key == string(GradeKeyProfessor) {
		if _, err := m.DB.CreateProfessorWithSd(ctx, sd, parameter.CreateProfessorParam{
			MemberID: member.MemberID,
		}); err != nil {
			return entity.MemberWithDetail{}, fmt.Errorf("failed to create professor: %w", err)
		}
	} else {
		if _, err := m.DB.CreateStudentWithSd(ctx, sd, parameter.CreateStudentParam{
			MemberID: member.MemberID,
		}); err != nil {
			return entity.MemberWithDetail{}, fmt.Errorf("failed to create student: %w", err)
		}
	}

	e = entity.MemberWithDetail{
		MemberID:               member.MemberID,
		Email:                  member.Email,
		Name:                   member.Name,
		FirstName:              member.FirstName,
		LastName:               member.LastName,
		AttendStatusID:         member.AttendStatusID,
		ProfileImageID:         member.ProfileImageID,
		GradeID:                member.GradeID,
		GroupID:                member.GroupID,
		PersonalOrganizationID: member.PersonalOrganizationID,
		RoleID:                 member.RoleID,
		Professor:              professor,
		Student:                student,
	}
	wholeOrg, err := m.DB.FindWholeOrganization(ctx)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	groupOrg, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, groupID)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find group with organization: %w", err)
	}
	gradeOrg, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, gradeID)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find grade with organization: %w", err)
	}
	// add some organization, chat room
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
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	bcrp := make([]parameter.BelongChatRoomParam, 0, len(bop))
	crap := make([]parameter.CreateChatRoomActionParam, 0, len(bop))
	if groupOrg.Organization.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: groupOrg.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		crap = append(crap, parameter.CreateChatRoomActionParam{
			ChatRoomID:           groupOrg.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
	}
	if gradeOrg.Organization.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: gradeOrg.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		crap = append(crap, parameter.CreateChatRoomActionParam{
			ChatRoomID:           gradeOrg.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
	}
	if wholeOrg.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: wholeOrg.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		crap = append(crap, parameter.CreateChatRoomActionParam{
			ChatRoomID:           wholeOrg.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
	}
	if org.ChatRoomID.Valid {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			MemberID:   e.MemberID,
			ChatRoomID: org.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		crap = append(crap, parameter.CreateChatRoomActionParam{
			ChatRoomID:           org.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
	}

	_, err = m.DB.BelongOrganizationsWithSd(ctx, sd, bop)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	_, err = m.DB.BelongChatRoomsWithSd(ctx, sd, bcrp)
	if err != nil {
		return entity.MemberWithDetail{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	for _, v := range crap {
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, v)
		if err != nil {
			return entity.MemberWithDetail{}, fmt.Errorf("failed to create chat room actions: %w", err)
		}
		craAdd, err := m.DB.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			AddedBy:          entity.UUID{},
		})
		if err != nil {
			return entity.MemberWithDetail{}, fmt.Errorf("failed to create chat room add member actions: %w", err)
		}
		_, err = m.DB.AddMemberToChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: e.MemberID},
		})
		if err != nil {
			return entity.MemberWithDetail{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
		}
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
	if e.ProfileImageID.Valid && e.ProfileImageID.Bytes != profileImageID.Bytes {
		defer func(ownerID, imageID uuid.UUID) {
			if err == nil {
				_, err = pluralDeleteImages(
					ctx,
					sd,
					m.DB,
					m.Storage,
					[]uuid.UUID{imageID},
					entity.UUID{
						Valid: true,
						Bytes: ownerID,
					},
					true,
				)
			}
		}(id, e.ProfileImageID.Bytes)
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
	now := m.Clocker.Now()
	ec, err := m.DB.FindMemberWithDetailWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find member with detail: %w", err)
	}
	var imageIDs []uuid.UUID
	var fileIDs []uuid.UUID
	attachableItems, err := m.DB.GetAttachableItemsWithSd(
		ctx,
		sd,
		parameter.WhereAttachableItemParam{
			WhereInOwner: true,
			InOwners:     []uuid.UUID{id},
		},
		parameter.AttachableItemOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get attachable items: %w", err)
	}
	for _, v := range attachableItems.Data {
		if v.Image.Valid {
			imageIDs = append(imageIDs, v.Image.Entity.ImageID)
		} else if v.File.Valid {
			fileIDs = append(fileIDs, v.File.Entity.FileID)
		}
	}
	if ec.ProfileImageID.Valid {
		imageIDs = append(imageIDs, ec.ProfileImageID.Bytes)
	}
	if len(imageIDs) > 0 {
		defer func(imageIDs []uuid.UUID, ownerID uuid.UUID) {
			if err == nil {
				_, err = pluralDeleteImages(
					ctx,
					sd,
					m.DB,
					m.Storage,
					imageIDs,
					entity.UUID{
						Valid: true,
						Bytes: ownerID,
					},
					true,
				)
			}
		}(imageIDs, id)
	}
	if len(fileIDs) > 0 {
		_, err = pluralDeleteFiles(
			ctx,
			sd,
			m.DB,
			m.Storage,
			fileIDs,
			entity.UUID{
				Valid: true,
				Bytes: id,
			},
			true,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to delete files: %w", err)
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
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyWithdraw))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	crs, err := m.DB.GetChatRoomsOnMemberWithSd(
		ctx,
		sd,
		id,
		parameter.WhereChatRoomOnMemberParam{},
		parameter.ChatRoomOnMemberOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get chat rooms on member: %w", err)
	}
	for _, v := range crs.Data {
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           v.ChatRoom.ChatRoomID,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomWithdrawActionWithSd(ctx, sd, parameter.CreateChatRoomWithdrawActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			MemberID:         entity.UUID{Valid: true, Bytes: id},
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room withdraw action: %w", err)
		}
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

// FindMemberByID メンバーを ID で取得する。
func (m *ManageMember) FindMemberByID(ctx context.Context, id uuid.UUID) (e entity.Member, err error) {
	e, err = m.DB.FindMemberByID(ctx, id)
	if err != nil {
		return entity.Member{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	return e, nil
}

// FindAuthMemberByID メンバーをIDで取得する。
func (m *ManageMember) FindAuthMemberByID(
	ctx context.Context,
	id uuid.UUID,
) (e entity.AuthMember, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.AuthMember{}, fmt.Errorf("failed to begin transaction: %w", err)
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

	member, err := m.DB.FindMemberWithProfileImageWithSd(ctx, sd, id)
	if err != nil {
		return entity.AuthMember{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	var role entity.NullableEntity[entity.RoleWithPolicies]
	if member.RoleID.Valid {
		re, err := m.DB.GetPoliciesOnRole(
			ctx,
			member.RoleID.Bytes,
			parameter.WherePolicyOnRoleParam{},
			parameter.PolicyOnRoleOrderMethodDefault,
			store.NumberedPaginationParam{},
			store.CursorPaginationParam{},
			store.WithCountParam{},
		)
		if err != nil {
			return entity.AuthMember{}, fmt.Errorf("failed to get policies on role: %w", err)
		}
		r, err := m.DB.FindRoleByID(ctx, member.RoleID.Bytes)
		if err != nil {
			return entity.AuthMember{}, fmt.Errorf("failed to find role by id: %w", err)
		}
		role = entity.NullableEntity[entity.RoleWithPolicies]{
			Valid: true,
			Entity: entity.RoleWithPolicies{
				Role: entity.Role{
					RoleID:      r.RoleID,
					Name:        r.Name,
					Description: r.Description,
				},
				Policies: re.Data,
			},
		}
	}
	return entity.AuthMember{
		MemberID:               id,
		Email:                  member.Email,
		Name:                   member.Name,
		FirstName:              member.FirstName,
		LastName:               member.LastName,
		AttendStatusID:         member.AttendStatusID,
		ProfileImage:           member.ProfileImage,
		GradeID:                member.GradeID,
		GroupID:                member.GroupID,
		PersonalOrganizationID: member.PersonalOrganizationID,
		Role:                   role,
	}, nil
}

// GetMembers メンバーを取得する。
func (m *ManageMember) GetMembers(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Member], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembers(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Member]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithAttendStatus メンバーを取得する。
func (m *ManageMember) GetMembersWithAttendStatus(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithAttendStatus(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithAttendStatus]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithDetail メンバーを取得する。
func (m *ManageMember) GetMembersWithDetail(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithDetail], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithDetail(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithDetail]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithCrew メンバーを取得する。
func (m *ManageMember) GetMembersWithCrew(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithCrew], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithCrew(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithCrew]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithProfileImage メンバーを取得する。
func (m *ManageMember) GetMembersWithProfileImage(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithProfileImage(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithProfileImage]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithPersonalOrganization メンバーを取得する。
func (m *ManageMember) GetMembersWithPersonalOrganization(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithPersonalOrganization(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithPersonalOrganization]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithRole メンバーを取得する。
func (m *ManageMember) GetMembersWithRole(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithRole], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMembersWithRole(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithRole]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return r, nil
}

// GetMembersWithCrewAndProfileImageAndAttendStatus メンバーを取得する。
func (m *ManageMember) GetMembersWithCrewAndProfileImageAndAttendStatus(
	ctx context.Context,
	whereSearchName string,
	whereHasInPolicies []uuid.UUID,
	whereInAttendStatuses []uuid.UUID,
	whereInGrades []uuid.UUID,
	whereInGroups []uuid.UUID,
	order parameter.MemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberParam{
		WhereLikeName:      whereSearchName != "",
		SearchName:         whereSearchName,
		WhereHasPolicy:     len(whereHasInPolicies) > 0,
		HasPolicyIDs:       whereHasInPolicies,
		WhenInAttendStatus: len(whereInAttendStatuses) > 0,
		InAttendStatusIDs:  whereInAttendStatuses,
		WhenInGrade:        len(whereInGrades) > 0,
		InGradeIDs:         whereInGrades,
		WhenInGroup:        len(whereInGroups) > 0,
		InGroupIDs:         whereInGroups,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	var e store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus]
	r, err := m.DB.GetMembersWithCrew(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus]{},
			fmt.Errorf("failed to get members: %w", err)
	}
	e.CursorPagination = r.CursorPagination
	e.WithCount = r.WithCount
	uniqueAttendStatues := make(map[uuid.UUID]entity.AttendStatus)
	uniqueProfileImages := make(map[uuid.UUID]entity.ImageWithAttachableItem)
	for _, v := range r.Data {
		uniqueAttendStatues[v.AttendStatusID] = entity.AttendStatus{}
		if v.ProfileImageID.Valid {
			uniqueProfileImages[v.ProfileImageID.Bytes] = entity.ImageWithAttachableItem{}
		}
	}
	attendStatuses := make([]uuid.UUID, 0, len(uniqueAttendStatues))
	for k := range uniqueAttendStatues {
		attendStatuses = append(attendStatuses, k)
	}
	profileImages := make([]uuid.UUID, 0, len(uniqueProfileImages))
	for k := range uniqueProfileImages {
		profileImages = append(profileImages, k)
	}
	attendStatusesData, err := m.DB.GetPluralAttendStatuses(
		ctx,
		attendStatuses,
		parameter.AttendStatusOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus]{},
			fmt.Errorf("failed to get plural attend statuses: %w", err)
	}
	for _, v := range attendStatusesData.Data {
		uniqueAttendStatues[v.AttendStatusID] = v
	}
	profileImagesData, err := m.DB.GetPluralImagesWithAttachableItem(
		ctx,
		profileImages,
		parameter.ImageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus]{},
			fmt.Errorf("failed to get plural images: %w", err)
	}
	for _, v := range profileImagesData.Data {
		uniqueProfileImages[v.ImageID] = v
	}

	e.Data = make([]entity.MemberWithCrewAndProfileImageAndAttendStatus, 0, len(r.Data))
	for _, v := range r.Data {
		attendStatus, ok := uniqueAttendStatues[v.AttendStatusID]
		if !ok {
			attendStatus = entity.AttendStatus{}
		}
		var profileImage entity.NullableEntity[entity.ImageWithAttachableItem]
		if v.ProfileImageID.Valid {
			if pi, ok := uniqueProfileImages[v.ProfileImageID.Bytes]; ok {
				profileImage = entity.NullableEntity[entity.ImageWithAttachableItem]{
					Valid:  true,
					Entity: pi,
				}
			}
		}
		e.Data = append(e.Data, entity.MemberWithCrewAndProfileImageAndAttendStatus{
			MemberID:               v.MemberID,
			Email:                  v.Email,
			Name:                   v.Name,
			FirstName:              v.FirstName,
			LastName:               v.LastName,
			AttendStatus:           attendStatus,
			Grade:                  v.Grade,
			Group:                  v.Group,
			PersonalOrganizationID: v.PersonalOrganizationID,
			ProfileImage:           profileImage,
		})
	}
	return e, nil
}

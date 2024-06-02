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

// ManageStudent 生徒管理サービス。
type ManageStudent struct {
	DB      store.Store
	Hash    hasher.Hash
	Clocker clock.Clock
	Storage storage.Storage
}

// CreateStudent 生徒を作成する。
func (m *ManageStudent) CreateStudent(
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
) (e entity.Student, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		return entity.Student{}, errhandle.NewModelDuplicatedError(MemberTargetLoginID)
	}
	// grade check
	grade, err := m.DB.FindGradeByIDWithSd(ctx, sd, gradeID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Student{}, errhandle.NewModelNotFoundError(MemberTargetGrades)
		}
		return entity.Student{}, fmt.Errorf("failed to find grade by id: %w", err)
	}
	if grade.Key == string(GradeKeyProfessor) {
		e := errhandle.NewCommonError(response.OnlyProfessorAction, nil)
		e.SetTarget(MemberTargetGrades)
		return entity.Student{}, e
	}
	// group check
	group, err := m.DB.FindGroupByIDWithSd(ctx, sd, groupID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Student{}, errhandle.NewModelNotFoundError(MemberTargetGroups)
		}
		return entity.Student{}, fmt.Errorf("failed to find group by id: %w", err)
	}
	if group.Key == string(GroupKeyProfessor) {
		e := errhandle.NewCommonError(response.OnlyProfessorAction, nil)
		e.SetTarget(MemberTargetGroups)
		return entity.Student{}, e
	}
	// role check
	if roleID.Valid {
		_, err = m.DB.FindRoleByIDWithSd(ctx, sd, roleID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Student{}, errhandle.NewModelNotFoundError(MemberTargetRoles)
			}
			return entity.Student{}, fmt.Errorf("failed to find role by id: %w", err)
		}
	}
	defaultAttendStatus, err := m.DB.FindAttendStatusByKeyWithSd(ctx, sd, string(DefaultAttendStatusKey))
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find default attend status: %w", err)
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
		return entity.Student{}, fmt.Errorf("failed to create organization: %w", err)
	}
	hashPass, err := m.Hash.Encrypt(rawPassword)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to encrypt password: %w", err)
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
		return entity.Student{}, fmt.Errorf("failed to create organization: %w", err)
	}
	e, err = m.DB.CreateStudentWithSd(ctx, sd, parameter.CreateStudentParam{
		MemberID: member.MemberID,
	})
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to create student: %w", err)
	}
	wholeOrg, err := m.DB.FindWholeOrganization(ctx)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	groupOrg, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, groupID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find group with organization: %w", err)
	}
	gradeOrg, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, gradeID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find grade with organization: %w", err)
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
		return entity.Student{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
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
		return entity.Student{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	_, err = m.DB.BelongChatRoomsWithSd(ctx, sd, bcrp)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	for _, v := range crap {
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, v)
		if err != nil {
			return entity.Student{}, fmt.Errorf("failed to create chat room actions: %w", err)
		}
		craAdd, err := m.DB.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			AddedBy:          entity.UUID{},
		})
		if err != nil {
			return entity.Student{}, fmt.Errorf("failed to create chat room add member actions: %w", err)
		}
		_, err = m.DB.AddMemberToChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: e.MemberID},
		})
		if err != nil {
			return entity.Student{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
		}
	}

	return e, nil
}

// DeleteStudent 生徒を削除する。
func (m *ManageStudent) DeleteStudent(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	st, err := m.DB.FindStudentWithMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find student by id: %w", err)
	}
	var imageIDs []uuid.UUID
	var fileIDs []uuid.UUID
	attachableItems, err := m.DB.GetAttachableItemsWithSd(
		ctx,
		sd,
		parameter.WhereAttachableItemParam{
			WhereInOwner: true,
			InOwners:     []uuid.UUID{st.Member.MemberID},
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
	if st.Member.ProfileImage.Valid {
		imageIDs = append(imageIDs, st.Member.ProfileImage.Entity.ImageID)
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
		}(imageIDs, st.Member.MemberID)
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
				Bytes: st.Member.MemberID,
			},
			true,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to delete files: %w", err)
		}
	}
	c, err = m.DB.DeleteStudentWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete student: %w", err)
	}
	e, err := m.DB.FindMemberWithPersonalOrganizationWithSd(ctx, sd, st.Member.MemberID)
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
		st.Member.MemberID,
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
			MemberID:         entity.UUID{Valid: true, Bytes: st.Member.MemberID},
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room withdraw action: %w", err)
		}
	}
	_, err = m.DB.DisbelongChatRoomOnMemberWithSd(ctx, sd, st.Member.MemberID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on member: %w", err)
	}
	_, err = m.DB.DisbelongOrganizationOnMemberWithSd(ctx, sd, st.Member.MemberID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong organization on member: %w", err)
	}
	_, err = m.DB.DeleteMemberWithSd(ctx, sd, st.Member.MemberID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete member: %w", err)
	}
	_, err = m.DB.DeleteOrganizationWithSd(ctx, sd, e.PersonalOrganization.OrganizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	return c, nil
}

// UpdateStudentGrade 生徒の学年を更新する。
func (m *ManageStudent) UpdateStudentGrade(
	ctx context.Context,
	id uuid.UUID,
	gradeID uuid.UUID,
) (e entity.StudentWithMember, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	// grade check
	grade, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, gradeID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.StudentWithMember{}, errhandle.NewModelNotFoundError(MemberTargetGrades)
		}
		return entity.StudentWithMember{}, fmt.Errorf("failed to find grade by id: %w", err)
	}
	if grade.Key == string(GradeKeyProfessor) {
		e := errhandle.NewCommonError(response.OnlyProfessorAction, nil)
		e.SetTarget(MemberTargetGrades)
		return entity.StudentWithMember{}, e
	}
	e, err = m.DB.FindStudentWithMemberWithSd(ctx, sd, id)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find student by id: %w", err)
	}
	if e.Member.GradeID == gradeID {
		return e, nil
	}
	member, err := m.DB.UpdateMemberGradeWithSd(ctx, sd, e.Member.MemberID, gradeID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to update member grade: %w", err)
	}
	addCraType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	withdrawCraType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyWithdraw))
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	originGrade, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, e.Member.GradeID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find grade with organization: %w", err)
	}

	_, err = m.DB.BelongOrganizationWithSd(ctx, sd, parameter.BelongOrganizationParam{
		MemberID:       e.Member.MemberID,
		OrganizationID: grade.Organization.OrganizationID,
		WorkPositionID: entity.UUID{},
		AddedAt:        now,
	})
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to belong organization: %w", err)
	}
	if grade.Organization.ChatRoomID.Valid {
		_, err = m.DB.BelongChatRoomWithSd(ctx, sd, parameter.BelongChatRoomParam{
			MemberID:   e.Member.MemberID,
			ChatRoomID: grade.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to belong chat room: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           grade.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		craAdd, err := m.DB.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			AddedBy:          entity.UUID{},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room add member actions: %w", err)
		}
		_, err = m.DB.AddMemberToChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: e.Member.MemberID},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
		}
	}

	_, err = m.DB.DisbelongOrganizationWithSd(ctx, sd, e.Member.MemberID, originGrade.Organization.OrganizationID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to disbelong organization: %w", err)
	}
	if originGrade.Organization.ChatRoomID.Valid {
		_, err = m.DB.DisbelongChatRoomWithSd(ctx, sd, e.Member.MemberID, originGrade.Organization.ChatRoomID.Bytes)
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to disbelong chat room: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           originGrade.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: withdrawCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomWithdrawActionWithSd(ctx, sd, parameter.CreateChatRoomWithdrawActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			MemberID:         entity.UUID{Valid: true, Bytes: e.Member.MemberID},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room withdraw action: %w", err)
		}
	}

	e.Member.GradeID = member.GradeID

	return e, nil
}

// UpdateStudentGroup 生徒の学年を更新する。
func (m *ManageStudent) UpdateStudentGroup(
	ctx context.Context,
	id uuid.UUID,
	groupID uuid.UUID,
) (e entity.StudentWithMember, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	// group check
	group, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, groupID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.StudentWithMember{}, errhandle.NewModelNotFoundError(MemberTargetGroups)
		}
		return entity.StudentWithMember{}, fmt.Errorf("failed to find group by id: %w", err)
	}
	if group.Key == string(GroupKeyProfessor) {
		e := errhandle.NewCommonError(response.OnlyProfessorAction, nil)
		e.SetTarget(MemberTargetGroups)
		return entity.StudentWithMember{}, e
	}
	e, err = m.DB.FindStudentWithMemberWithSd(ctx, sd, id)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find student by id: %w", err)
	}
	if e.Member.GroupID == groupID {
		return e, nil
	}
	member, err := m.DB.UpdateMemberGroupWithSd(ctx, sd, e.Member.MemberID, groupID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to update member group: %w", err)
	}
	addCraType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	withdrawCraType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyWithdraw))
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	originGroup, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, e.Member.GroupID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to find group with organization: %w", err)
	}

	_, err = m.DB.BelongOrganizationWithSd(ctx, sd, parameter.BelongOrganizationParam{
		MemberID:       e.Member.MemberID,
		OrganizationID: group.Organization.OrganizationID,
		WorkPositionID: entity.UUID{},
		AddedAt:        now,
	})
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to belong organization: %w", err)
	}
	if group.Organization.ChatRoomID.Valid {
		_, err = m.DB.BelongChatRoomWithSd(ctx, sd, parameter.BelongChatRoomParam{
			MemberID:   e.Member.MemberID,
			ChatRoomID: group.Organization.ChatRoomID.Bytes,
			AddedAt:    now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to belong chat room: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           group.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		craAdd, err := m.DB.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			AddedBy:          entity.UUID{},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room add member actions: %w", err)
		}
		_, err = m.DB.AddMemberToChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: e.Member.MemberID},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
		}
	}

	_, err = m.DB.DisbelongOrganizationWithSd(ctx, sd, e.Member.MemberID, originGroup.Organization.OrganizationID)
	if err != nil {
		return entity.StudentWithMember{}, fmt.Errorf("failed to disbelong organization: %w", err)
	}
	if originGroup.Organization.ChatRoomID.Valid {
		_, err = m.DB.DisbelongChatRoomWithSd(ctx, sd, e.Member.MemberID, originGroup.Organization.ChatRoomID.Bytes)
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to disbelong chat room: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           originGroup.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: withdrawCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomWithdrawActionWithSd(ctx, sd, parameter.CreateChatRoomWithdrawActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			MemberID:         entity.UUID{Valid: true, Bytes: e.Member.MemberID},
		})
		if err != nil {
			return entity.StudentWithMember{}, fmt.Errorf("failed to create chat room withdraw action: %w", err)
		}
	}

	e.Member.GroupID = member.GradeID

	return e, nil
}

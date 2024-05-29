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
		return entity.Student{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	_, err = m.DB.BelongChatRoomsWithSd(ctx, sd, bcrp)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to belong chat rooms: %w", err)
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
	st, err := m.DB.FindStudentWithMemberWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find student by id: %w", err)
	}
	if st.Member.ProfileImage.Valid {
		_, err = pluralDeleteImages(
			ctx,
			sd,
			m.DB,
			m.Storage,
			[]uuid.UUID{st.Member.ProfileImage.Entity.ImageID},
			entity.UUID{
				Valid: true,
				Bytes: id,
			})
		if err != nil {
			return 0, fmt.Errorf("failed to delete images: %w", err)
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
	e, err = m.DB.FindStudentByIDWithSd(ctx, sd, id)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find student by id: %w", err)
	}
	_, err = m.DB.UpdateMemberGradeWithSd(ctx, sd, e.MemberID, gradeID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to update member grade: %w", err)
	}
	return e, nil
}

// UpdateStudentGroup 生徒のグループを更新する。
func (m *ManageStudent) UpdateStudentGroup(
	ctx context.Context,
	id uuid.UUID,
	groupID uuid.UUID,
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
	e, err = m.DB.FindStudentByIDWithSd(ctx, sd, id)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to find student by id: %w", err)
	}
	_, err = m.DB.UpdateMemberGroupWithSd(ctx, sd, e.MemberID, groupID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("failed to update member group: %w", err)
	}
	return e, nil
}

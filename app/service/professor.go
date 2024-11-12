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
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageProfessor 教授管理サービス。
type ManageProfessor struct {
	DB      store.Store
	Hash    hasher.Hash
	Clocker clock.Clock
	Storage storage.Storage
	WsHub   ws.HubInterface
}

// CreateProfessor 教授を作成する。
func (m *ManageProfessor) CreateProfessor(
	ctx context.Context,
	loginID,
	rawPassword,
	email,
	name string,
	firstName,
	lastName entity.String,
	roleID entity.UUID,
) (e entity.Professor, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		return entity.Professor{}, errhandle.NewModelDuplicatedError(MemberTargetLoginID)
	}
	// grade check
	grade, err := m.DB.FindGradeByKeyWithSd(ctx, sd, string(GradeKeyProfessor))
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find grade by key: %w", err)
	}
	// group check
	group, err := m.DB.FindGroupByKeyWithSd(ctx, sd, string(GroupKeyProfessor))
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find group by key: %w", err)
	}
	// role check
	if roleID.Valid {
		_, err = m.DB.FindRoleByIDWithSd(ctx, sd, roleID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Professor{}, errhandle.NewModelNotFoundError(MemberTargetRoles)
			}
			return entity.Professor{}, fmt.Errorf("failed to find role by id: %w", err)
		}
	}
	defaultAttendStatus, err := m.DB.FindAttendStatusByKeyWithSd(ctx, sd, string(DefaultAttendStatusKey))
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find default attend status: %w", err)
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
		return entity.Professor{}, fmt.Errorf("failed to create organization: %w", err)
	}
	hashPass, err := m.Hash.Encrypt(rawPassword)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to encrypt password: %w", err)
	}
	p := parameter.CreateMemberParam{
		LoginID:                loginID,
		Password:               hashPass,
		Email:                  email,
		Name:                   name,
		FirstName:              firstName,
		LastName:               lastName,
		AttendStatusID:         defaultAttendStatus.AttendStatusID,
		GradeID:                grade.GradeID,
		GroupID:                group.GroupID,
		RoleID:                 roleID,
		PersonalOrganizationID: org.OrganizationID,
	}
	member, err := m.DB.CreateMemberWithSd(ctx, sd, p)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to create organization: %w", err)
	}
	e, err = m.DB.CreateProfessorWithSd(ctx, sd, parameter.CreateProfessorParam{
		MemberID: member.MemberID,
	})
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to create student: %w", err)
	}
	wholeOrg, err := m.DB.FindWholeOrganization(ctx)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	groupOrg, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, group.GroupID)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find group with organization: %w", err)
	}
	gradeOrg, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, grade.GradeID)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find grade with organization: %w", err)
	}
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
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
	bcrp := make([]parameter.BelongChatRoomParam, 0, len(bop))
	crap := make([]parameter.CreateChatRoomActionParam, 0, len(bop))
	wsTargets := make([]ws.Targets, 0, len(bop))
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
		belonging, err := m.DB.GetMembersOnChatRoomWithSd(
			ctx,
			sd,
			groupOrg.Organization.ChatRoomID.Bytes,
			parameter.WhereMemberOnChatRoomParam{},
			parameter.MemberOnChatRoomOrderMethodDefault,
			store.NumberedPaginationParam{},
			store.CursorPaginationParam{},
			store.WithCountParam{},
		)
		if err != nil {
			return entity.Professor{}, fmt.Errorf("failed to get members on chat room: %w", err)
		}
		memberIDs := make([]uuid.UUID, 0, len(belonging.Data))
		for _, v := range belonging.Data {
			memberIDs = append(memberIDs, v.Member.MemberID)
		}
		wsTargets = append(wsTargets, ws.Targets{
			Members: memberIDs,
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
		belonging, err := m.DB.GetMembersOnChatRoomWithSd(
			ctx,
			sd,
			gradeOrg.Organization.ChatRoomID.Bytes,
			parameter.WhereMemberOnChatRoomParam{},
			parameter.MemberOnChatRoomOrderMethodDefault,
			store.NumberedPaginationParam{},
			store.CursorPaginationParam{},
			store.WithCountParam{},
		)
		if err != nil {
			return entity.Professor{}, fmt.Errorf("failed to get members on chat room: %w", err)
		}
		memberIDs := make([]uuid.UUID, 0, len(belonging.Data))
		for _, v := range belonging.Data {
			memberIDs = append(memberIDs, v.Member.MemberID)
		}
		wsTargets = append(wsTargets, ws.Targets{
			Members: memberIDs,
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
		wsTargets = append(wsTargets, ws.Targets{
			All: true,
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
		return entity.Professor{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	_, err = m.DB.BelongChatRoomsWithSd(ctx, sd, bcrp)
	if err != nil {
		return entity.Professor{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	for i, v := range crap {
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, v)
		if err != nil {
			return entity.Professor{}, fmt.Errorf("failed to create chat room actions: %w", err)
		}
		craAdd, err := m.DB.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			AddedBy:          entity.UUID{},
		})
		if err != nil {
			return entity.Professor{}, fmt.Errorf("failed to create chat room add member actions: %w", err)
		}
		_, err = m.DB.AddMemberToChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: e.MemberID},
		})
		if err != nil {
			return entity.Professor{}, fmt.Errorf("failed to add member to chat room add member actions: %w", err)
		}
		action := entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{
			ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
			ChatRoomActionID:          cra.ChatRoomActionID,
			AddedBy:                   entity.NullableEntity[entity.SimpleMember]{},
			AddMembers: []entity.MemberOnChatRoomAddMemberAction{
				{
					ChatRoomAddMemberActionID: craAdd.ChatRoomAddMemberActionID,
					Member: entity.NullableEntity[entity.SimpleMember]{
						Valid: true,
						Entity: entity.SimpleMember{
							MemberID:       member.MemberID,
							Name:           member.Name,
							FirstName:      member.FirstName,
							LastName:       member.LastName,
							Email:          member.Email,
							ProfileImageID: member.ProfileImageID,
							GradeID:        member.GradeID,
							GroupID:        member.GroupID,
						},
					},
				},
			},
		}
		defer func(
			roomID uuid.UUID, wsTarget ws.Targets,
			action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMember, wsTarget,
					ws.ChatRoomAddedMemberEventData{
						ChatRoomID:           roomID,
						Action:               action,
						ChatRoomActionID:     actAttr.ChatRoomActionID,
						ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
						ActedAt:              actAttr.ActedAt,
					})
			}
		}(v.ChatRoomID, wsTargets[i], action, cra)
	}

	return e, nil
}

// DeleteProfessor 教授を削除する。
func (m *ManageProfessor) DeleteProfessor(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	st, err := m.DB.FindProfessorWithMemberWithSd(ctx, sd, id)
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
	c, err = m.DB.DeleteProfessorWithSd(ctx, sd, id)
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
	chatRoomIDs := make([]uuid.UUID, 0, len(crs.Data))
	for _, v := range crs.Data {
		chatRoomIDs = append(chatRoomIDs, v.ChatRoom.ChatRoomID)
	}
	belongings, err := m.DB.GetPluralMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoomIDs,
		store.NumberedPaginationParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural members on chat room: %w", err)
	}
	belongingsMap := make(map[uuid.UUID][]entity.MemberOnChatRoomWithChatRoomID, len(belongings.Data))
	for _, v := range belongings.Data {
		belongingsMap[v.ChatRoomID] = append(belongingsMap[v.ChatRoomID], v)
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
		belonging, ok := belongingsMap[v.ChatRoom.ChatRoomID]
		if !ok {
			continue
		}
		memberIDs := make([]uuid.UUID, 0, len(belonging))
		for _, v := range belonging {
			memberIDs = append(memberIDs, v.Member.MemberID)
		}
		action := entity.ChatRoomWithdrawActionWithMember{
			ChatRoomWithdrawActionID: cra.ChatRoomActionID,
			ChatRoomActionID:         cra.ChatRoomActionID,
			Member: entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       id,
					Name:           e.Name,
					FirstName:      e.FirstName,
					LastName:       e.LastName,
					Email:          e.Email,
					ProfileImageID: e.ProfileImageID,
					GradeID:        e.GradeID,
					GroupID:        e.GroupID,
				},
			},
		}
		defer func(
			roomID uuid.UUID, memberIDs []uuid.UUID,
			action entity.ChatRoomWithdrawActionWithMember,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomWithdrawnMember, ws.Targets{
					Members: memberIDs,
				}, ws.ChatRoomWithdrawnMemberEventData{
					ChatRoomID:           roomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
			}
		}(v.ChatRoom.ChatRoomID, memberIDs, action, cra)
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

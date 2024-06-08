package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageChatRoomBelonging チャットルーム所属管理サービス。
type ManageChatRoomBelonging struct {
	DB      store.Store
	Clocker clock.Clock
	WsHub   ws.HubInterface
}

func belongMembersOnChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	chatRoom entity.ChatRoom,
	owner entity.Member,
	members []entity.Member,
	force bool,
) (e int64, action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers,
	actAttr entity.ChatRoomAction, err error,
) {
	if !force && chatRoom.FromOrganization {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotAddMemberToOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotAddMemberToPrivateChatRoom, nil)
	}
	bcrp := make([]parameter.BelongChatRoomParam, len(members))
	for i, member := range members {
		bcrp[i] = parameter.BelongChatRoomParam{
			ChatRoomID: chatRoom.ChatRoomID,
			MemberID:   member.MemberID,
			AddedAt:    now,
		}
	}
	if e, err = str.BelongChatRoomsWithSd(
		ctx,
		sd,
		bcrp,
	); err != nil {
		var ufe errhandle.ModelDuplicatedError
		if errors.As(err, &ufe) {
			return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
				entity.ChatRoomAction{},
				errhandle.NewModelDuplicatedError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	addCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to find chat room action type: %w", err)
	}
	addCra, err := str.CreateChatRoomActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomActionParam{
			ChatRoomID:           chatRoom.ChatRoomID,
			ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		},
	)
	if err != nil {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room action: %w", err)
	}
	crama, err := str.CreateChatRoomAddMemberActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: addCra.ChatRoomActionID,
			AddedBy:          entity.UUID{Valid: true, Bytes: owner.MemberID},
		},
	)
	if err != nil {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room add member action: %w", err)
	}
	cramap := make([]parameter.CreateChatRoomAddedMemberParam, len(members))
	for i, member := range members {
		cramap[i] = parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: crama.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: member.MemberID},
		}
	}

	if _, err = str.AddMembersToChatRoomAddMemberActionWithSd(
		ctx,
		sd,
		cramap,
	); err != nil {
		return 0, entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to add members to chat room add member action: %w", err)
	}

	addMembers := make([]entity.MemberOnChatRoomAddMemberAction, len(members))
	for i, member := range members {
		addMembers[i] = entity.MemberOnChatRoomAddMemberAction{
			ChatRoomAddMemberActionID: crama.ChatRoomAddMemberActionID,
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
		}
	}

	action = entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{
		ChatRoomAddMemberActionID: crama.ChatRoomAddMemberActionID,
		ChatRoomActionID:          addCra.ChatRoomActionID,
		AddedBy: entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       owner.MemberID,
				Name:           owner.Name,
				FirstName:      owner.FirstName,
				LastName:       owner.LastName,
				Email:          owner.Email,
				ProfileImageID: owner.ProfileImageID,
				GradeID:        owner.GradeID,
				GroupID:        owner.GroupID,
			},
		},
	}

	return e, action, addCra, nil
}

func removeMembersFromChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	chatRoom entity.ChatRoom,
	owner entity.Member,
	members []entity.Member,
	force bool,
) (e int64, action entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers,
	actAttr entity.ChatRoomAction, err error,
) {
	if !force && chatRoom.FromOrganization {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotRemoveMemberFromOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotRemoveMemberFromPrivateChatRoom, nil)
	}
	memberIDs := make([]uuid.UUID, len(members))
	for i, member := range members {
		memberIDs[i] = member.MemberID
	}
	if e, err = str.DisbelongPluralMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoom.ChatRoomID,
		memberIDs,
	); err != nil {
		var ufe errhandle.ModelNotFoundError
		if errors.As(err, &ufe) {
			return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
				entity.ChatRoomAction{},
				errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to disbelong plural members on chat room: %w", err)
	}

	removeCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyRemoveMember))
	if err != nil {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to find chat room action type: %w", err)
	}
	removeCra, err := str.CreateChatRoomActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomActionParam{
			ChatRoomID:           chatRoom.ChatRoomID,
			ChatRoomActionTypeID: removeCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		},
	)
	if err != nil {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room action: %w", err)
	}
	crama, err := str.CreateChatRoomRemoveMemberActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomRemoveMemberActionParam{
			ChatRoomActionID: removeCra.ChatRoomActionID,
			RemovedBy:        entity.UUID{Valid: true, Bytes: owner.MemberID},
		},
	)
	if err != nil {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room remove member action: %w", err)
	}
	cramap := make([]parameter.CreateChatRoomRemovedMemberParam, len(members))
	for i, member := range members {
		cramap[i] = parameter.CreateChatRoomRemovedMemberParam{
			ChatRoomRemoveMemberActionID: crama.ChatRoomRemoveMemberActionID,
			MemberID:                     entity.UUID{Valid: true, Bytes: member.MemberID},
		}
	}
	if _, err = str.RemoveMembersToChatRoomRemoveMemberActionWithSd(
		ctx,
		sd,
		cramap,
	); err != nil {
		return 0, entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to remove members to chat room remove member action: %w", err)
	}

	removedMembers := make([]entity.MemberOnChatRoomRemoveMemberAction, len(members))
	for i, member := range members {
		removedMembers[i] = entity.MemberOnChatRoomRemoveMemberAction{
			ChatRoomRemoveMemberActionID: crama.ChatRoomRemoveMemberActionID,
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
		}
	}

	action = entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{
		ChatRoomRemoveMemberActionID: crama.ChatRoomRemoveMemberActionID,
		ChatRoomActionID:             removeCra.ChatRoomActionID,
		RemovedBy: entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       owner.MemberID,
				Name:           owner.Name,
				FirstName:      owner.FirstName,
				LastName:       owner.LastName,
				Email:          owner.Email,
				ProfileImageID: owner.ProfileImageID,
				GradeID:        owner.GradeID,
				GroupID:        owner.GroupID,
			},
		},
		RemoveMembers: removedMembers,
	}

	return e, action, removeCra, nil
}

func withdrawMemberFromChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	chatRoom entity.ChatRoom,
	member entity.Member,
	force bool,
) (e int64, action entity.ChatRoomWithdrawActionWithMember,
	actAttr entity.ChatRoomAction, err error,
) {
	if !force && chatRoom.FromOrganization {
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotWithdrawMemberFromOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			errhandle.NewCommonError(response.CannotWithdrawMemberFromPrivateChatRoom, nil)
	}
	if e, err = str.DisbelongPluralMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoom.ChatRoomID,
		[]uuid.UUID{member.MemberID},
	); err != nil {
		var ufe errhandle.ModelNotFoundError
		if errors.As(err, &ufe) {
			return 0, entity.ChatRoomWithdrawActionWithMember{},
				entity.ChatRoomAction{},
				errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to disbelong plural members on chat room: %w", err)
	}

	withdrawCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyWithdraw))
	if err != nil {
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to find chat room action type: %w", err)
	}
	withdrawCra, err := str.CreateChatRoomActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomActionParam{
			ChatRoomID:           chatRoom.ChatRoomID,
			ChatRoomActionTypeID: withdrawCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		},
	)
	if err != nil {
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room action: %w", err)
	}
	withdrawAct, err := str.CreateChatRoomWithdrawActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomWithdrawActionParam{
			ChatRoomActionID: withdrawCra.ChatRoomActionID,
			MemberID:         entity.UUID{Valid: true, Bytes: member.MemberID},
		},
	)
	if err != nil {
		return 0, entity.ChatRoomWithdrawActionWithMember{},
			entity.ChatRoomAction{},
			fmt.Errorf("failed to create chat room withdraw action: %w", err)
	}

	action = entity.ChatRoomWithdrawActionWithMember{
		ChatRoomWithdrawActionID: withdrawAct.ChatRoomWithdrawActionID,
		ChatRoomActionID:         withdrawCra.ChatRoomActionID,
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
	}

	return e, action, withdrawCra, nil
}

// BelongMembersOnChatRoom メンバーをチャットルームに所属させる。
func (m *ManageChatRoomBelonging) BelongMembersOnChatRoom(
	ctx context.Context,
	chatRoomID,
	ownerID uuid.UUID,
	memberIDs []uuid.UUID,
) (e int64, err error) {
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
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	room, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, chatRoomID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
		}
		return 0, fmt.Errorf("failed to find chat room: %w", err)
	}
	mm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		memberIDs,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(mm.Data) != len(memberIDs) {
		return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetMembers)
	}
	belongingMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		room.ChatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	alreadyMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		alreadyMemberIDs = append(alreadyMemberIDs, v.Member.MemberID)
	}

	var action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers
	var actAttr entity.ChatRoomAction
	e, action, actAttr, err = belongMembersOnChatRoom(ctx, sd, now, m.DB, room, owner, mm.Data, false)

	defer func(
		room entity.ChatRoom, membersIDs, alreadyMemberIDs []uuid.UUID,
		action entity.ChatRoomAddMemberActionWithAddedByAndAddMembers,
		actAttr entity.ChatRoomAction,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
				Members: membersIDs,
			}, ws.ChatRoomAddedMeEventData{
				ChatRoom: room,
			})
			m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMember, ws.Targets{
				Members: alreadyMemberIDs,
			}, ws.ChatRoomAddedMemberEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               action,
				ChatRoomActionID:     actAttr.ChatRoomActionID,
				ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
				ActedAt:              actAttr.ActedAt,
			})
		}
	}(room, memberIDs, alreadyMemberIDs, action, actAttr)

	return e, err
}

// RemoveMembersFromChatRoom チャットルームからメンバーを削除する。
func (m *ManageChatRoomBelonging) RemoveMembersFromChatRoom(
	ctx context.Context,
	chatRoomID,
	ownerID uuid.UUID,
	memberIDs []uuid.UUID,
) (e int64, err error) {
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
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	room, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, chatRoomID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
		}
		return 0, fmt.Errorf("failed to find chat room: %w", err)
	}
	mm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		memberIDs,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(mm.Data) != len(memberIDs) {
		return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetMembers)
	}
	belongingMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		room.ChatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	removeMembers := make(map[uuid.UUID]struct{}, len(mm.Data))
	for _, v := range memberIDs {
		removeMembers[v] = struct{}{}
	}
	leftMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		if _, ok := removeMembers[v.Member.MemberID]; ok {
			continue
		}
		leftMemberIDs = append(leftMemberIDs, v.Member.MemberID)
	}

	var action entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers
	var actAttr entity.ChatRoomAction
	e, action, actAttr, err = removeMembersFromChatRoom(ctx, sd, now, m.DB, room, owner, mm.Data, false)

	defer func(
		room entity.ChatRoom, membersIDs, leftMemberIDs []uuid.UUID,
		action entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers,
		actAttr entity.ChatRoomAction,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomRemovedMe, ws.Targets{
				Members: membersIDs,
			}, ws.ChatRoomRemovedMeEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               action,
				ChatRoomActionID:     actAttr.ChatRoomActionID,
				ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
				ActedAt:              actAttr.ActedAt,
			})
			m.WsHub.Dispatch(ws.EventTypeChatRoomRemovedMember, ws.Targets{
				Members: leftMemberIDs,
			}, ws.ChatRoomRemovedMemberEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               action,
				ChatRoomActionID:     actAttr.ChatRoomActionID,
				ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
				ActedAt:              actAttr.ActedAt,
			})
		}
	}(room, memberIDs, leftMemberIDs, action, actAttr)

	return e, err
}

// WithdrawMemberFromChatRoom チャットルームからメンバーを退会させる。
func (m *ManageChatRoomBelonging) WithdrawMemberFromChatRoom(
	ctx context.Context,
	chatRoomID,
	memberID uuid.UUID,
) (e int64, err error) {
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
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, memberID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	room, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, chatRoomID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoom)
		}
		return 0, fmt.Errorf("failed to find chat room: %w", err)
	}
	belongingMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		room.ChatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	leftMemberIDs := make([]uuid.UUID, 0, len(belongingMembers.Data))
	for _, v := range belongingMembers.Data {
		if v.Member.MemberID == memberID {
			continue
		}
		leftMemberIDs = append(leftMemberIDs, v.Member.MemberID)
	}

	var action entity.ChatRoomWithdrawActionWithMember
	var actAttr entity.ChatRoomAction
	e, action, actAttr, err = withdrawMemberFromChatRoom(ctx, sd, now, m.DB, room, owner, false)

	defer func(
		room entity.ChatRoom, leftMemberIDs []uuid.UUID,
		action entity.ChatRoomWithdrawActionWithMember,
		actAttr entity.ChatRoomAction,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomWithdrawnMember, ws.Targets{
				Members: leftMemberIDs,
			}, ws.ChatRoomWithdrawnMemberEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               action,
				ChatRoomActionID:     actAttr.ChatRoomActionID,
				ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
				ActedAt:              actAttr.ActedAt,
			})
		}
	}(room, leftMemberIDs, action, actAttr)

	return e, err
}

// GetChatRoomsOnMember メンバーに関連付けられたチャットルームを取得する。
func (m *ManageChatRoomBelonging) GetChatRoomsOnMember(
	ctx context.Context,
	memberID uuid.UUID,
	whereSearchName string,
	order parameter.ChatRoomOnMemberOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (es store.ListResult[entity.PracticalChatRoomOnMember], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereChatRoomOnMemberParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
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
	e, err := m.DB.GetChatRoomsOnMember(ctx, memberID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PracticalChatRoomOnMember]{}, fmt.Errorf("failed to get chat rooms on member: %w", err)
	}
	es = store.ListResult[entity.PracticalChatRoomOnMember]{
		CursorPagination: e.CursorPagination,
		WithCount:        e.WithCount,
	}
	privateRoomIDs := make([]uuid.UUID, 0, len(e.Data))
	es.Data = make([]entity.PracticalChatRoomOnMember, len(e.Data))
	for _, v := range e.Data {
		if v.ChatRoom.IsPrivate {
			privateRoomIDs = append(privateRoomIDs, v.ChatRoom.ChatRoomID)
		}
	}
	companions, err := m.DB.GetPluralPrivateChatRoomCompanions(
		ctx,
		privateRoomIDs,
		memberID,
		store.NumberedPaginationParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
	)
	if err != nil {
		return store.ListResult[entity.PracticalChatRoomOnMember]{},
			fmt.Errorf("failed to get plural private chat room companions: %w", err)
	}
	companionMap := make(map[uuid.UUID]entity.MemberOnChatRoom, len(companions.Data))
	for _, v := range companions.Data {
		companionMap[v.ChatRoomID] = entity.MemberOnChatRoom{
			Member:  v.Member,
			AddedAt: v.AddedAt,
		}
	}
	for i, v := range e.Data {
		var cp entity.NullableEntity[entity.MemberOnChatRoom]
		if companion, ok := companionMap[v.ChatRoom.ChatRoomID]; ok {
			cp = entity.NullableEntity[entity.MemberOnChatRoom]{
				Valid:  true,
				Entity: companion,
			}
		}
		es.Data[i] = entity.PracticalChatRoomOnMember{
			ChatRoom: entity.PracticalChatRoom{
				ChatRoomID:       v.ChatRoom.ChatRoomID,
				Name:             v.ChatRoom.Name,
				IsPrivate:        v.ChatRoom.IsPrivate,
				FromOrganization: v.ChatRoom.FromOrganization,
				CoverImage:       v.ChatRoom.CoverImage,
				OwnerID:          v.ChatRoom.OwnerID,
				LatestMessage:    v.ChatRoom.LatestMessage,
				LatestAction:     v.ChatRoom.LatestAction,
				Companion:        cp,
			},
			AddedAt: v.AddedAt,
		}
	}
	return es, nil
}

// GetChatRoomsOnMemberCount メンバーに関連付けられたチャットルームの数を取得する。
func (m *ManageChatRoomBelonging) GetChatRoomsOnMemberCount(
	ctx context.Context, memberID uuid.UUID,
	whereSearchName string,
) (es int64, err error) {
	where := parameter.WhereChatRoomOnMemberParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err = m.DB.CountChatRoomsOnMember(ctx, memberID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get chat rooms on member count: %w", err)
	}
	return es, nil
}

// GetMembersOnChatRoom チャットルームに関連付けられたメンバーを取得する。
func (m *ManageChatRoomBelonging) GetMembersOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	whereSearchName string,
	order parameter.MemberOnChatRoomOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (es store.ListResult[entity.MemberOnChatRoom], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMemberOnChatRoomParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
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
	es, err = m.DB.GetMembersOnChatRoom(ctx, chatRoomID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MemberOnChatRoom]{}, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	return es, nil
}

// GetMembersOnChatRoomCount チャットルームに関連付けられたメンバーの数を取得する。
func (m *ManageChatRoomBelonging) GetMembersOnChatRoomCount(
	ctx context.Context, chatRoomID uuid.UUID,
	whereSearchName string,
) (es int64, err error) {
	where := parameter.WhereMemberOnChatRoomParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err = m.DB.CountMembersOnChatRoom(ctx, chatRoomID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on chat room count: %w", err)
	}
	return es, nil
}

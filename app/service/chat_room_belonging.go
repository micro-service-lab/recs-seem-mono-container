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
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageChatRoomBelonging チャットルーム所属管理サービス。
type ManageChatRoomBelonging struct {
	DB      store.Store
	Clocker clock.Clock
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
) (e int64, err error) {
	if !force && chatRoom.FromOrganization {
		return 0, errhandle.NewCommonError(response.CannotAddMemberToOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, errhandle.NewCommonError(response.CannotAddMemberToPrivateChatRoom, nil)
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
			return 0, errhandle.NewModelDuplicatedError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	addCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type: %w", err)
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
		return 0, fmt.Errorf("failed to create chat room action: %w", err)
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
		return 0, fmt.Errorf("failed to create chat room add member action: %w", err)
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
		return 0, fmt.Errorf("failed to add members to chat room add member action: %w", err)
	}

	return e, nil
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
) (e int64, err error) {
	if !force && chatRoom.FromOrganization {
		return 0, errhandle.NewCommonError(response.CannotRemoveMemberFromOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, errhandle.NewCommonError(response.CannotRemoveMemberFromPrivateChatRoom, nil)
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
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, fmt.Errorf("failed to disbelong plural members on chat room: %w", err)
	}

	removeCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyRemoveMember))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type: %w", err)
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
		return 0, fmt.Errorf("failed to create chat room action: %w", err)
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
		return 0, fmt.Errorf("failed to create chat room remove member action: %w", err)
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
		return 0, fmt.Errorf("failed to remove members to chat room remove member action: %w", err)
	}

	return e, nil
}

func withdrawMemberFromChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	chatRoom entity.ChatRoom,
	member entity.Member,
	force bool,
) (e int64, err error) {
	if !force && chatRoom.FromOrganization {
		return 0, errhandle.NewCommonError(response.CannotWithdrawMemberFromOrganizationChatRoom, nil)
	}
	if !force && chatRoom.IsPrivate {
		return 0, errhandle.NewCommonError(response.CannotWithdrawMemberFromPrivateChatRoom, nil)
	}
	if e, err = str.DisbelongPluralMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoom.ChatRoomID,
		[]uuid.UUID{member.MemberID},
	); err != nil {
		var ufe errhandle.ModelNotFoundError
		if errors.As(err, &ufe) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomBelongingTargetChatRoomBelongings)
		}
		return 0, fmt.Errorf("failed to disbelong plural members on chat room: %w", err)
	}

	withdrawCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyWithdraw))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type: %w", err)
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
		return 0, fmt.Errorf("failed to create chat room action: %w", err)
	}
	if _, err := str.CreateChatRoomWithdrawActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomWithdrawActionParam{
			ChatRoomActionID: withdrawCra.ChatRoomActionID,
			MemberID:         entity.UUID{Valid: true, Bytes: member.MemberID},
		},
	); err != nil {
		return 0, fmt.Errorf("failed to create chat room withdraw action: %w", err)
	}

	return e, nil
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

	return belongMembersOnChatRoom(ctx, sd, now, m.DB, room, owner, mm.Data, false)
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

	return removeMembersFromChatRoom(ctx, sd, now, m.DB, room, owner, mm.Data, false)
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

	return withdrawMemberFromChatRoom(ctx, sd, now, m.DB, room, owner, false)
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

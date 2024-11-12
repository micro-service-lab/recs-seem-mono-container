package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ManageChatRoomAction チャットルームアクション管理サービス。
type ManageChatRoomAction struct {
	DB store.Store
}

// GetChatRoomActionsOnChatRoom チャットルームに紐づくチャットルームアクションを取得する。
func (m *ManageChatRoomAction) GetChatRoomActionsOnChatRoom(
	ctx context.Context,
	chatRoomID,
	ownerID uuid.UUID,
	whereInTypes []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (e store.ListResult[entity.ChatRoomActionPractical], err error) {
	_, err = m.DB.FindChatRoomByID(ctx, chatRoomID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return store.ListResult[entity.ChatRoomActionPractical]{}, errhandle.NewModelNotFoundError("chat room")
		}
		return store.ListResult[entity.ChatRoomActionPractical]{}, fmt.Errorf("failed to find chat room by id: %w", err)
	}

	_, err = m.DB.FindMemberByID(ctx, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return store.ListResult[entity.ChatRoomActionPractical]{}, errhandle.NewModelNotFoundError("owner")
		}
		return store.ListResult[entity.ChatRoomActionPractical]{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	if b, err := m.DB.ExistsChatRoomBelonging(ctx, ownerID, chatRoomID); err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{}, fmt.Errorf("failed to check chat room belonging: %w", err)
	} else if !b {
		fmt.Println("NotChatRoomMember")
		return store.ListResult[entity.ChatRoomActionPractical]{}, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereChatRoomActionParam{
		WhereInChatRoomActionTypeIDs: len(whereInTypes) > 0,
		InChatRoomActionTypeIDs:      whereInTypes,
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
	r, err := m.DB.GetChatRoomActionsOnChatRoom(
		ctx,
		chatRoomID,
		where,
		order,
		np,
		cp,
		wc,
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{}, fmt.Errorf("failed to get chat room action: %w", err)
	}
	ids := make([]uuid.UUID, len(r.Data))
	for _, v := range r.Data {
		ids = append(ids, v.ChatRoomActionID)
	}
	createActionsRet, err := m.DB.GetPluralChatRoomCreateActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomCreateActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room create action: %w", err)
	}
	createActions := make(map[uuid.UUID]entity.ChatRoomCreateActionWithCreatedBy, len(createActionsRet.Data))
	for _, v := range createActionsRet.Data {
		createActions[v.ChatRoomActionID] = v
	}

	updateNameActionsRet, err := m.DB.GetPluralChatRoomUpdateNameActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomUpdateNameActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room update name action: %w", err)
	}
	updateNameActions := make(map[uuid.UUID]entity.ChatRoomUpdateNameActionWithUpdatedBy, len(updateNameActionsRet.Data))
	for _, v := range updateNameActionsRet.Data {
		updateNameActions[v.ChatRoomActionID] = v
	}

	addMemberActionsRet, err := m.DB.GetPluralChatRoomAddMemberActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomAddMemberActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room add member action: %w", err)
	}
	addMemberActions := make(map[uuid.UUID]entity.ChatRoomAddMemberActionWithAddedBy, len(addMemberActionsRet.Data))
	addMemberActionIDs := make([]uuid.UUID, len(addMemberActionsRet.Data))
	for i, v := range addMemberActionsRet.Data {
		addMemberActions[v.ChatRoomActionID] = v
		addMemberActionIDs[i] = v.ChatRoomAddMemberActionID
	}

	removeMemberActionsRet, err := m.DB.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomRemoveMemberActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room remove member action: %w", err)
	}
	removeMemberActions := make(
		map[uuid.UUID]entity.ChatRoomRemoveMemberActionWithRemovedBy, len(removeMemberActionsRet.Data))
	removeMemberActionIDs := make([]uuid.UUID, len(removeMemberActionsRet.Data))
	for i, v := range removeMemberActionsRet.Data {
		removeMemberActions[v.ChatRoomActionID] = v
		removeMemberActionIDs[i] = v.ChatRoomRemoveMemberActionID
	}

	withdrawActionsRet, err := m.DB.GetPluralChatRoomWithdrawActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomWithdrawActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room withdraw action: %w", err)
	}
	withdrawActions := make(map[uuid.UUID]entity.ChatRoomWithdrawActionWithMember, len(withdrawActionsRet.Data))
	for _, v := range withdrawActionsRet.Data {
		withdrawActions[v.ChatRoomActionID] = v
	}

	deleteMessageActionsRet, err := m.DB.GetPluralChatRoomDeleteMessageActionsByChatRoomActionIDs(
		ctx,
		ids,
		parameter.ChatRoomDeleteMessageActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural chat room delete message action: %w", err)
	}
	deleteMessageActions := make(
		map[uuid.UUID]entity.ChatRoomDeleteMessageActionWithDeletedBy, len(deleteMessageActionsRet.Data))
	for _, v := range deleteMessageActionsRet.Data {
		deleteMessageActions[v.ChatRoomActionID] = v
	}

	messagesRet, err := m.DB.GetPluralMessagesWithSenderByChatRoomActionIDs(
		ctx,
		ids,
		parameter.MessageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural messages with sender: %w", err)
	}
	messages := make(map[uuid.UUID]entity.MessageWithSender, len(messagesRet.Data))
	messageIDs := make([]uuid.UUID, len(messagesRet.Data))
	for i, v := range messagesRet.Data {
		messages[v.ChatRoomActionID] = v
		messageIDs[i] = v.MessageID
	}

	am := make(map[uuid.UUID][]entity.MemberOnChatRoomAddMemberAction, len(addMemberActionIDs))
	addMembers, err := m.DB.GetPluralMembersOnChatRoomAddMemberAction(
		ctx,
		addMemberActionIDs,
		parameter.MemberOnChatRoomAddMemberActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural members on chat room add member action: %w", err)
	}
	for _, v := range addMembers.Data {
		am[v.ChatRoomAddMemberActionID] = append(am[v.ChatRoomAddMemberActionID], v)
	}
	rm := make(map[uuid.UUID][]entity.MemberOnChatRoomRemoveMemberAction, len(removeMemberActionIDs))
	removeMembers, err := m.DB.GetPluralMembersOnChatRoomRemoveMemberAction(
		ctx,
		removeMemberActionIDs,
		parameter.MemberOnChatRoomRemoveMemberActionOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural members on chat room remove member action: %w", err)
	}
	for _, v := range removeMembers.Data {
		rm[v.ChatRoomRemoveMemberActionID] = append(rm[v.ChatRoomRemoveMemberActionID], v)
	}

	rs := make(map[uuid.UUID]int64, len(messageIDs))
	reads, err := m.DB.CountReadsOnMessages(
		ctx,
		messageIDs,
		parameter.WhereReadsOnMessageParam{
			WhereIsRead: true,
		},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to count reads on messages: %w", err)
	}
	for _, v := range reads {
		rs[v.MessageID] = v.Count
	}

	ai := make(map[uuid.UUID][]entity.AttachedItemOnMessage, len(messageIDs))
	attachments, err := m.DB.GetPluralAttachedItemsOnMessage(
		ctx,
		messageIDs,
		parameter.AttachedItemOnMessageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.ChatRoomActionPractical]{},
			fmt.Errorf("failed to get plural attached items on message: %w", err)
	}
	for _, v := range attachments.Data {
		ai[v.MessageID] = append(ai[v.MessageID], v)
	}

	de := make([]entity.ChatRoomActionPractical, len(r.Data))
	for i, v := range r.Data {
		var createAction entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]
		var updateNameAction entity.NullableEntity[entity.ChatRoomUpdateNameActionWithUpdatedBy]
		var addMemberAction entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedByAndAddMembers]
		var removeMemberAction entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers]
		var withdrawAction entity.NullableEntity[entity.ChatRoomWithdrawActionWithMember]
		var deleteMessageAction entity.NullableEntity[entity.ChatRoomDeleteMessageActionWithDeletedBy]
		var message entity.NullableEntity[entity.MessageWithSenderAndReadReceiptCountAndAttachments]

		if v, ok := createActions[v.ChatRoomActionID]; ok {
			createAction = entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]{
				Valid:  true,
				Entity: v,
			}
		}

		if v, ok := updateNameActions[v.ChatRoomActionID]; ok {
			updateNameAction = entity.NullableEntity[entity.ChatRoomUpdateNameActionWithUpdatedBy]{
				Valid:  true,
				Entity: v,
			}
		}

		if v, ok := addMemberActions[v.ChatRoomActionID]; ok {
			attr, ok := am[v.ChatRoomAddMemberActionID]
			if !ok {
				attr = []entity.MemberOnChatRoomAddMemberAction{}
			}
			addMemberAction = entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedByAndAddMembers]{
				Valid: true,
				Entity: entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{
					ChatRoomAddMemberActionID: v.ChatRoomAddMemberActionID,
					ChatRoomActionID:          v.ChatRoomActionID,
					AddedBy:                   v.AddedBy,
					AddMembers:                attr,
				},
			}
		}

		if v, ok := removeMemberActions[v.ChatRoomActionID]; ok {
			attr, ok := rm[v.ChatRoomRemoveMemberActionID]
			if !ok {
				attr = []entity.MemberOnChatRoomRemoveMemberAction{}
			}
			removeMemberAction = entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers]{
				Valid: true,
				Entity: entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{
					ChatRoomRemoveMemberActionID: v.ChatRoomRemoveMemberActionID,
					ChatRoomActionID:             v.ChatRoomActionID,
					RemovedBy:                    v.RemovedBy,
					RemoveMembers:                attr,
				},
			}
		}

		if v, ok := withdrawActions[v.ChatRoomActionID]; ok {
			withdrawAction = entity.NullableEntity[entity.ChatRoomWithdrawActionWithMember]{
				Valid:  true,
				Entity: v,
			}
		}

		if v, ok := deleteMessageActions[v.ChatRoomActionID]; ok {
			deleteMessageAction = entity.NullableEntity[entity.ChatRoomDeleteMessageActionWithDeletedBy]{
				Valid:  true,
				Entity: v,
			}
		}

		if v, ok := messages[v.ChatRoomActionID]; ok {
			rc, ok := rs[v.MessageID]
			if !ok {
				rc = 0
			}
			att, ok := ai[v.MessageID]
			if !ok {
				att = []entity.AttachedItemOnMessage{}
			}
			message = entity.NullableEntity[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{
				Valid: true,
				Entity: entity.MessageWithSenderAndReadReceiptCountAndAttachments{
					MessageID:        v.MessageID,
					ChatRoomActionID: v.ChatRoomActionID,
					Sender:           v.Sender,
					Body:             v.Body,
					PostedAt:         v.PostedAt,
					LastEditedAt:     v.LastEditedAt,
					ReadReceiptCount: rc,
					Attachments:      att,
				},
			}
		}

		de[i] = entity.ChatRoomActionPractical{
			ChatRoomActionID:            v.ChatRoomActionID,
			ChatRoomID:                  v.ChatRoomID,
			ChatRoomActionTypeID:        v.ChatRoomActionTypeID,
			ActedAt:                     v.ActedAt,
			ChatRoomCreateAction:        createAction,
			ChatRoomUpdateNameAction:    updateNameAction,
			ChatRoomAddMemberAction:     addMemberAction,
			ChatRoomRemoveMemberAction:  removeMemberAction,
			ChatRoomWithdrawAction:      withdrawAction,
			ChatRoomDeleteMessageAction: deleteMessageAction,
			Message:                     message,
		}
	}

	e = store.ListResult[entity.ChatRoomActionPractical]{
		Data:             de,
		CursorPagination: r.CursorPagination,
		WithCount:        r.WithCount,
	}
	return e, nil
}

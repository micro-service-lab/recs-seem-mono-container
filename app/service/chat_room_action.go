package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageChatRoomAction チャットルームアクション管理サービス。
type ManageChatRoomAction struct {
	DB store.Store
}

// GetChatRoomActionsOnChatRoom チャットルームに紐づくチャットルームアクションを取得する。
func (m *ManageChatRoomAction) GetChatRoomActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	whereInTypes []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (e store.ListResult[entity.ChatRoomActionPractical], err error) {
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
	var addMemberActionIDs []uuid.UUID
	var removeMemberActionIDs []uuid.UUID
	var messageIDs []uuid.UUID
	for _, v := range r.Data {
		if v.ChatRoomAddMemberAction.Valid {
			addMemberActionIDs = append(
				addMemberActionIDs, v.ChatRoomAddMemberAction.Entity.ChatRoomAddMemberActionID)
		} else if v.ChatRoomRemoveMemberAction.Valid {
			removeMemberActionIDs = append(
				removeMemberActionIDs, v.ChatRoomRemoveMemberAction.Entity.ChatRoomRemoveMemberActionID)
		} else if v.Message.Valid {
			messageIDs = append(messageIDs, v.Message.Entity.MessageID)
		}
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
		var addMemberAction entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedByAndAddMembers]
		var removeMemberAction entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers]
		var message entity.NullableEntity[entity.MessageWithSenderAndReadReceiptCountAndAttachments]

		if v.ChatRoomAddMemberAction.Valid {
			attr, ok := am[v.ChatRoomAddMemberAction.Entity.ChatRoomAddMemberActionID]
			if !ok {
				attr = []entity.MemberOnChatRoomAddMemberAction{}
			}
			addMemberAction = entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedByAndAddMembers]{
				Valid: true,
				Entity: entity.ChatRoomAddMemberActionWithAddedByAndAddMembers{
					ChatRoomAddMemberActionID: v.ChatRoomAddMemberAction.Entity.ChatRoomAddMemberActionID,
					ChatRoomActionID:          v.ChatRoomAddMemberAction.Entity.ChatRoomActionID,
					AddedBy:                   v.ChatRoomAddMemberAction.Entity.AddedBy,
					AddMembers:                attr,
				},
			}
		}

		if v.ChatRoomRemoveMemberAction.Valid {
			attr, ok := rm[v.ChatRoomRemoveMemberAction.Entity.ChatRoomRemoveMemberActionID]
			if !ok {
				attr = []entity.MemberOnChatRoomRemoveMemberAction{}
			}
			removeMemberAction = entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers]{
				Valid: true,
				Entity: entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers{
					ChatRoomRemoveMemberActionID: v.ChatRoomRemoveMemberAction.Entity.ChatRoomRemoveMemberActionID,
					ChatRoomActionID:             v.ChatRoomRemoveMemberAction.Entity.ChatRoomActionID,
					RemovedBy:                    v.ChatRoomRemoveMemberAction.Entity.RemovedBy,
					RemoveMembers:                attr,
				},
			}
		}

		if v.Message.Valid {
			rc, ok := rs[v.Message.Entity.MessageID]
			if !ok {
				rc = 0
			}
			att, ok := ai[v.Message.Entity.MessageID]
			if !ok {
				att = []entity.AttachedItemOnMessage{}
			}
			message = entity.NullableEntity[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{
				Valid: true,
				Entity: entity.MessageWithSenderAndReadReceiptCountAndAttachments{
					MessageID:        v.Message.Entity.MessageID,
					ChatRoomActionID: v.Message.Entity.ChatRoomActionID,
					Sender:           v.Message.Entity.Sender,
					Body:             v.Message.Entity.Body,
					PostedAt:         v.Message.Entity.PostedAt,
					LastEditedAt:     v.Message.Entity.LastEditedAt,
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
			ChatRoomCreateAction:        v.ChatRoomCreateAction,
			ChatRoomUpdateNameAction:    v.ChatRoomUpdateNameAction,
			ChatRoomAddMemberAction:     addMemberAction,
			ChatRoomRemoveMemberAction:  removeMemberAction,
			ChatRoomWithdrawAction:      v.ChatRoomWithdrawAction,
			ChatRoomDeleteMessageAction: v.ChatRoomDeleteMessageAction,
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

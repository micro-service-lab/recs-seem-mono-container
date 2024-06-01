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
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMessage メッセージ管理サービス。
type ManageMessage struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
}

// CreateMessage メッセージを作成する。
func (m *ManageMessage) CreateMessage(
	ctx context.Context,
	senderID, chatRoomID uuid.UUID,
	content string,
	attachments []uuid.UUID,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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

	_, err = m.DB.FindMemberByIDWithSd(ctx, sd, senderID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetSender)
		}
		return entity.Message{}, fmt.Errorf("failed to find member: %w", err)
	}
	_, err = m.DB.FindChatRoomByIDWithSd(ctx, sd, chatRoomID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetChatRoom)
		}
		return entity.Message{}, fmt.Errorf("failed to find chat room: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	readableMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for i, v := range belongMembers.Data {
		if v.Member.MemberID == senderID {
			belong = true
		} else {
			readableMemberIDs[i] = v.Member.MemberID
		}
	}
	if !belong {
		return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetChatRoomBelongings)
	}
	if len(attachments) > 0 {
		ai, err := m.DB.GetPluralAttachableItemsWithSd(
			ctx,
			sd,
			attachments,
			parameter.AttachableItemOrderMethodDefault,
			store.NumberedPaginationParam{},
		)
		if err != nil {
			return entity.Message{}, fmt.Errorf("failed to get plural attachable items: %w", err)
		}
		if len(ai.Data) != len(attachments) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetAttachments)
		}
	}

	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyMessage))
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           chatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	e, err = m.DB.CreateMessageWithSd(
		ctx,
		sd,
		parameter.CreateMessageParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			SenderID:         entity.UUID{Valid: true, Bytes: senderID},
			Body:             content,
			PostedAt:         now,
		},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create message: %w", err)
	}
	if len(attachments) > 0 {
		aiParams := make([]parameter.AttachItemMessageParam, len(attachments))
		for i, v := range attachments {
			aiParams[i] = parameter.AttachItemMessageParam{
				MessageID:        e.MessageID,
				AttachableItemID: entity.UUID{Valid: true, Bytes: v},
			}
		}
		if _, err = m.DB.AttacheItemsOnMessagesWithSd(
			ctx,
			sd,
			aiParams,
		); err != nil {
			return entity.Message{}, fmt.Errorf("failed to attache items on messages: %w", err)
		}
	}
	if len(readableMemberIDs) > 0 {
		rrParams := make([]parameter.CreateReadReceiptParam, len(readableMemberIDs))
		for i, v := range readableMemberIDs {
			rrParams[i] = parameter.CreateReadReceiptParam{
				MessageID: e.MessageID,
				MemberID:  v,
			}
		}
		if _, err = m.DB.CreateReadReceiptsWithSd(
			ctx,
			sd,
			rrParams,
		); err != nil {
			return entity.Message{}, fmt.Errorf("failed to create read receipts: %w", err)
		}
	}
	return e, nil
}

// CreateMessageOnPrivateRoom 個人チャットルームにメッセージを作成する。
func (m *ManageMessage) CreateMessageOnPrivateRoom(
	ctx context.Context,
	senderID, receiverID uuid.UUID,
	content string,
	attachments []uuid.UUID,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	return e, nil
}

// DeleteMessage メッセージを削除する。
func (m *ManageMessage) DeleteMessage(
	ctx context.Context,
	ownerID, messageID uuid.UUID,
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
	return 0, nil
}

// EditMessage メッセージを編集する。
func (m *ManageMessage) EditMessage(
	ctx context.Context,
	ownerID, messageID uuid.UUID,
	content string,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	return e, nil
}

// GetMessagesOnChatRoom チャットルームのメッセージを取得する。
func (m *ManageMessage) GetMessagesOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	whereInSenders []uuid.UUID,
	whereSearchBody string,
	whereEarlierPostedAt time.Time,
	whereLaterPostedAt time.Time,
	whereEarlierLastEditedAt time.Time,
	whereLaterLastEditedAt time.Time,
	order parameter.MessageOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (e store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMessageParam{
		WhereInChatRoom:          true,
		InChatRoom:               []uuid.UUID{chatRoomID},
		WhereInSender:            len(whereInSenders) > 0,
		InSender:                 whereInSenders,
		WhereLikeBody:            whereSearchBody != "",
		SearchBody:               whereSearchBody,
		WhereEarlierPostedAt:     !whereEarlierPostedAt.IsZero(),
		EarlierPostedAt:          whereEarlierPostedAt,
		WhereLaterPostedAt:       !whereLaterPostedAt.IsZero(),
		LaterPostedAt:            whereLaterPostedAt,
		WhereEarlierLastEditedAt: !whereEarlierLastEditedAt.IsZero(),
		EarlierLastEditedAt:      whereEarlierLastEditedAt,
		WhereLaterLastEditedAt:   !whereLaterLastEditedAt.IsZero(),
		LaterLastEditedAt:        whereLaterLastEditedAt,
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
	var message store.ListResult[entity.MessageWithSender]
	if message, err = m.DB.GetMessagesWithSender(
		ctx,
		where,
		order,
		np,
		cp,
		wc,
	); err != nil {
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to get messages: %w", err)
	}
	messageIDs := make([]uuid.UUID, len(message.Data))
	for i, v := range message.Data {
		messageIDs[i] = v.MessageID
	}

	e.CursorPagination = message.CursorPagination
	e.WithCount = message.WithCount

	rs := make(map[uuid.UUID]int64, len(messageIDs))
	reads, err := m.DB.CountReadsOnMessages(
		ctx,
		messageIDs,
		parameter.WhereReadsOnMessageParam{
			WhereIsRead: true,
		},
	)
	if err != nil {
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to count reads: %w", err)
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
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to get plural attached items on message: %w", err)
	}
	for _, v := range attachments.Data {
		ai[v.MessageID] = append(ai[v.MessageID], v)
	}

	e.Data = make([]entity.MessageWithSenderAndReadReceiptCountAndAttachments, len(message.Data))

	for i, v := range message.Data {
		rc, ok := rs[v.MessageID]
		if !ok {
			rc = 0
		}
		att, ok := ai[v.MessageID]
		if !ok {
			att = []entity.AttachedItemOnMessage{}
		}
		e.Data[i] = entity.MessageWithSenderAndReadReceiptCountAndAttachments{
			MessageID:        v.MessageID,
			ChatRoomActionID: v.ChatRoomActionID,
			Sender:           v.Sender,
			Body:             v.Body,
			PostedAt:         v.PostedAt,
			LastEditedAt:     v.LastEditedAt,
			ReadReceiptCount: rc,
			Attachments:      att,
		}
	}

	return e, nil
}

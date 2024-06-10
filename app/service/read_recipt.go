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
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageReadReceipt 既読管理サービス。
type ManageReadReceipt struct {
	DB      store.Store
	Clocker clock.Clock
	WsHub   ws.HubInterface
}

// CountUnreadReceiptsOnMember メンバー上の未読既読数を取得する。
func (m *ManageReadReceipt) CountUnreadReceiptsOnMember(
	ctx context.Context,
	memberID uuid.UUID,
) (int64, error) {
	e, err := m.DB.CountReadableMessagesOnMember(
		ctx,
		memberID,
		parameter.WhereReadableMessageOnChatRoomAndMemberParam{
			WhereIsNotRead: true,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread receipts on member: %w", err)
	}
	return e, nil
}

// ReadMessage 既読処理を行う。
func (m *ManageReadReceipt) ReadMessage(
	ctx context.Context,
	chatRoomID,
	memberID, messageID uuid.UUID,
) (read bool, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
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
	msg, err := m.DB.FindMessageWithChatRoomActionWithSd(ctx, sd, messageID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return false, errhandle.NewModelNotFoundError(MessageTargetMessages)
		}
		return false, fmt.Errorf("failed to find message: %w", err)
	}
	if msg.SenderID.Bytes == memberID {
		return false, errhandle.NewCommonError(response.CannotReadOwnMessage, nil)
	}
	if msg.ChatRoomAction.ChatRoomID != chatRoomID {
		return false, errhandle.NewCommonError(response.NotMatchChatRoomMessage, nil)
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
		return false, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	belongMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == memberID {
			belong = true
		}
		belongMemberIDs = append(belongMemberIDs, v.Member.MemberID)
	}
	if !belong {
		return false, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	rr, err := m.DB.FindReadReceiptWithSd(ctx, sd, memberID, messageID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return false, errhandle.NewModelNotFoundError(ReadReceiptTargetReadReceipts)
		}
		return false, fmt.Errorf("failed to find read receipt: %w", err)
	}
	if rr.ReadAt.Valid {
		return false, nil
	}
	if _, err = m.DB.ReadReceiptWithSd(ctx, sd, parameter.ReadReceiptParam{
		MemberID:  memberID,
		MessageID: messageID,
		ReadAt:    entity.Timestamptz{Time: now, Valid: true},
	}); err != nil {
		return false, fmt.Errorf("failed to read message: %w", err)
	}

	defer func(
		roomID, memberID uuid.UUID, belongMemberIDs, messageIDs []uuid.UUID,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomReadMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomReadMessageEventData{
				ChatRoomID: roomID,
				MemberID:   memberID,
				MessageIDs: messageIDs,
			})
		}
	}(chatRoomID, memberID, belongMemberIDs, []uuid.UUID{messageID})

	return true, nil
}

// ReadMessagesOnMember メンバー上のメッセージを既読にする。
func (m *ManageReadReceipt) ReadMessagesOnMember(
	ctx context.Context,
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
	if e, err = m.DB.ReadReceiptsOnMemberWithSd(
		ctx,
		sd,
		memberID,
		now,
	); err != nil {
		return 0, fmt.Errorf("failed to read messages on member: %w", err)
	}

	// TODO: wsの通知処理を追加する

	return e, nil
}

// ReadMessagesOnChatRoomAndMember チャットルーム、メンバー上のメッセージを既読にする。
func (m *ManageReadReceipt) ReadMessagesOnChatRoomAndMember(
	ctx context.Context,
	chatRoomID, memberID uuid.UUID,
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
		return 0, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	belongMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == memberID {
			belong = true
		}
		belongMemberIDs = append(belongMemberIDs, v.Member.MemberID)
	}
	if !belong {
		return 0, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	readableMessages, err := m.DB.GetReadableMessagesOnChatRoomAndMemberWithSd(
		ctx,
		sd,
		chatRoomID,
		memberID,
		parameter.WhereReadableMessageOnChatRoomAndMemberParam{
			WhereIsNotRead: true,
		},
		parameter.ReadableMessageOnChatRoomAndMemberOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get readable messages on chat room and member: %w", err)
	}
	if len(readableMessages.Data) == 0 {
		return 0, nil
	}
	messageIDs := make([]uuid.UUID, len(readableMessages.Data))
	for i, v := range readableMessages.Data {
		messageIDs[i] = v.Message.MessageID
	}
	if e, err = m.DB.ReadReceiptsOnChatRoomAndMemberWithSd(
		ctx,
		sd,
		chatRoomID,
		memberID,
		now,
	); err != nil {
		return 0, fmt.Errorf("failed to read messages on chat room and member: %w", err)
	}

	defer func(
		roomID, memberID uuid.UUID, belongMemberIDs, messageIDs []uuid.UUID,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomReadMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomReadMessageEventData{
				ChatRoomID: roomID,
				MemberID:   memberID,
				MessageIDs: messageIDs,
			})
		}
	}(chatRoomID, memberID, belongMemberIDs, messageIDs)

	return e, nil
}

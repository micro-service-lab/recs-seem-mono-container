package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// AttachedMessage メッセージ添付を表すインターフェース。
type AttachedMessage interface {
	// CountAttachedItemsOnChatRoom チャットルームに関連付けられた添付アイテム数を取得する。
	CountAttachedItemsOnChatRoom(
		ctx context.Context, chatRoomID uuid.UUID, where parameter.WhereAttachedItemOnChatRoomParam) (int64, error)
	// CountAttachedItemsOnChatRoomWithSd SD付きでチャットルームに関連付けられた添付アイテム数を取得する。
	CountAttachedItemsOnChatRoomWithSd(
		ctx context.Context, sd Sd, chatRoomID uuid.UUID, where parameter.WhereAttachedItemOnChatRoomParam) (int64, error)
	// CountAttachedItemsOnMessage メッセージに関連付けられた添付アイテム数を取得する。
	CountAttachedItemsOnMessage(
		ctx context.Context, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam) (int64, error)
	// CountAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテム数を取得する。
	CountAttachedItemsOnMessageWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam) (int64, error)
	// AttacheItemOnMessage 添付アイテムを関連付ける。
	AttacheItemOnMessage(ctx context.Context, param parameter.AttachItemMessageParam) (entity.AttachedMessage, error)
	// AttacheItemOnMessageWithSd SD付きで添付アイテムを関連付ける。
	AttacheItemOnMessageWithSd(
		ctx context.Context, sd Sd, param parameter.AttachItemMessageParam) (entity.AttachedMessage, error)
	// AttacheItemsOnMessages 添付アイテムを複数関連付ける。
	AttacheItemsOnMessages(ctx context.Context, params []parameter.AttachItemMessageParam) (int64, error)
	// AttacheItemsOnMessagesWithSd SD付きで添付アイテムを複数関連付ける。
	AttacheItemsOnMessagesWithSd(ctx context.Context, sd Sd, params []parameter.AttachItemMessageParam) (int64, error)
	// DetachAttachedMessage 添付アイテムを関連付けを解除する。
	DetachAttachedMessage(ctx context.Context, attachedMessageID uuid.UUID) (int64, error)
	// DetachAttachedMessageWithSd SD付きで添付アイテムを関連付けを解除する。
	DetachAttachedMessageWithSd(ctx context.Context, sd Sd, attachedMessageID uuid.UUID) (int64, error)
	// DetachItemsOnMessage メッセージに関連付けられた添付アイテムを解除する。
	DetachItemsOnMessage(ctx context.Context, messageID uuid.UUID) (int64, error)
	// DetachItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテムを解除する。
	DetachItemsOnMessageWithSd(ctx context.Context, sd Sd, messageID uuid.UUID) (int64, error)
	// DetachItemsOnMessages メッセージに関連付けられた添付アイテムを複数解除する。
	DetachItemsOnMessages(ctx context.Context, messageIDs []uuid.UUID) (int64, error)
	// DetachItemsOnMessagesWithSd SD付きでメッセージに関連付けられた添付アイテムを複数解除する。
	DetachItemsOnMessagesWithSd(ctx context.Context, sd Sd, messageIDs []uuid.UUID) (int64, error)
	// PluralDetachItemsOnMessage メッセージに関連付けられた複数の添付アイテムを解除する。
	PluralDetachItemsOnMessage(ctx context.Context, messageID uuid.UUID, attachedItemIDs []uuid.UUID) (int64, error)
	// PluralDetachItemsOnMessageWithSd SD付きでメッセージに関連付けられた複数の添付アイテムを解除する。
	PluralDetachItemsOnMessageWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID, attachedItemIDs []uuid.UUID) (int64, error)
	// GetAttachedItemsOnChatRoom チャットルームに関連付けられた添付アイテムを取得する。
	GetAttachedItemsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereAttachedItemOnChatRoomParam,
		order parameter.AttachedItemOnChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnChatRoom], error)
	// GetAttachedItemsOnChatRoomWithSd SD付きでチャットルームに関連付けられた添付アイテムを取得する。
	GetAttachedItemsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereAttachedItemOnChatRoomParam,
		order parameter.AttachedItemOnChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnChatRoom], error)
	// GetAttachedItemsOnMessage メッセージに関連付けられた添付アイテムを取得する。
	GetAttachedItemsOnMessage(
		ctx context.Context,
		messageID uuid.UUID,
		where parameter.WhereAttachedItemOnMessageParam,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnMessage], error)
	// GetAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテムを取得する。
	GetAttachedItemsOnMessageWithSd(
		ctx context.Context,
		sd Sd,
		messageID uuid.UUID,
		where parameter.WhereAttachedItemOnMessageParam,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnMessage], error)
	// GetPluralAttachedItemsOnMessage メッセージに関連付けられた複数の添付アイテムを取得する。
	GetPluralAttachedItemsOnMessage(
		ctx context.Context,
		messageIDs []uuid.UUID,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachedItemOnMessage], error)
	// GetPluralAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた複数の添付アイテムを取得する。
	GetPluralAttachedItemsOnMessageWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachedItemOnMessage], error)
	// GetAttachedItemsOnMessageWithMimeType メッセージに関連付けられた添付アイテムとそのマイムタイプを取得する。
	GetAttachedItemsOnMessageWithMimeType(
		ctx context.Context,
		messageID uuid.UUID,
		where parameter.WhereAttachedItemOnMessageParam,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnMessageWithMimeType], error)
	// GetAttachedItemsOnMessageWithMimeTypeWithSd SD付きでメッセージに関連付けられた添付アイテムとそのマイムタイプを取得する。
	GetAttachedItemsOnMessageWithMimeTypeWithSd(
		ctx context.Context,
		sd Sd,
		messageID uuid.UUID,
		where parameter.WhereAttachedItemOnMessageParam,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachedItemOnMessageWithMimeType], error)
	// GetPluralAttachedItemsOnMessageWithMimeType メッセージに関連付けられた複数の添付アイテムとそのマイムタイプを取得する。
	GetPluralAttachedItemsOnMessageWithMimeType(
		ctx context.Context,
		messageIDs []uuid.UUID,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachedItemOnMessageWithMimeType], error)
	// GetPluralAttachedItemsOnMessageWithMimeTypeWithSd SD付きでメッセージに関連付けられた複数の添付アイテムとそのマイムタイプを取得する。
	GetPluralAttachedItemsOnMessageWithMimeTypeWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		order parameter.AttachedItemOnMessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachedItemOnMessageWithMimeType], error)
}

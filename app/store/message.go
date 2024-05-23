package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Message メッセージを表すインターフェース。
type Message interface {
	// CountMessages メッセージ数を取得する。
	CountMessages(ctx context.Context, where parameter.WhereMessageParam) (int64, error)
	// CountMessagesWithSd SD付きでメッセージ数を取得する。
	CountMessagesWithSd(ctx context.Context, sd Sd, where parameter.WhereMessageParam) (int64, error)
	// CreateMessage メッセージを作成する。
	CreateMessage(ctx context.Context, param parameter.CreateMessageParam) (entity.Message, error)
	// CreateMessageWithSd SD付きでメッセージを作成する。
	CreateMessageWithSd(
		ctx context.Context, sd Sd, param parameter.CreateMessageParam) (entity.Message, error)
	// CreateMessages メッセージを作成する。
	CreateMessages(ctx context.Context, params []parameter.CreateMessageParam) (int64, error)
	// CreateMessagesWithSd SD付きでメッセージを作成する。
	CreateMessagesWithSd(ctx context.Context, sd Sd, params []parameter.CreateMessageParam) (int64, error)
	// DeleteMessage メッセージを削除する。
	DeleteMessage(ctx context.Context, messageID uuid.UUID) (int64, error)
	// DeleteMessageWithSd SD付きでメッセージを削除する。
	DeleteMessageWithSd(ctx context.Context, sd Sd, messageID uuid.UUID) (int64, error)
	// DeleteMessagesOnChatRoom チャットルーム内のメッセージを削除する。
	DeleteMessagesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error)
	// DeleteMessagesOnChatRoomWithSd SD付きでチャットルーム内のメッセージを削除する。
	DeleteMessagesOnChatRoomWithSd(ctx context.Context, sd Sd, chatRoomID uuid.UUID) (int64, error)
	// PluralDeleteMessages メッセージを複数削除する。
	PluralDeleteMessages(ctx context.Context, messageIDs []uuid.UUID) (int64, error)
	// PluralDeleteMessagesWithSd SD付きでメッセージを複数削除する。
	PluralDeleteMessagesWithSd(ctx context.Context, sd Sd, messageIDs []uuid.UUID) (int64, error)
	// FindMessageByID メッセージを取得する。
	FindMessageByID(ctx context.Context, messageID uuid.UUID) (entity.Message, error)
	// FindMessageByIDWithSd SD付きでメッセージを取得する。
	FindMessageByIDWithSd(ctx context.Context, sd Sd, messageID uuid.UUID) (entity.Message, error)
	// FindMessageWithChatRoom メッセージを取得する。
	FindMessageWithChatRoom(ctx context.Context, messageID uuid.UUID) (entity.MessageWithChatRoom, error)
	// FindMessageWithChatRoomWithSd SD付きでメッセージを取得する。
	FindMessageWithChatRoomWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID) (entity.MessageWithChatRoom, error)
	// FindMessageWithSender メッセージを取得する。
	FindMessageWithSender(ctx context.Context, messageID uuid.UUID) (entity.MessageWithSender, error)
	// FindMessageWithSenderWithSd SD付きでメッセージを取得する。
	FindMessageWithSenderWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID) (entity.MessageWithSender, error)
	// GetMessages メッセージを取得する。
	GetMessages(
		ctx context.Context,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Message], error)
	// GetMessagesWithSd SD付きでメッセージを取得する。
	GetMessagesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Message], error)
	// GetPluralMessages メッセージを取得する。
	GetPluralMessages(
		ctx context.Context,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Message], error)
	// GetPluralMessagesWithSd SD付きでメッセージを取得する。
	GetPluralMessagesWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Message], error)
	// GetMessagesWithChatRoom メッセージを取得する。
	GetMessagesWithChatRoom(
		ctx context.Context,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MessageWithChatRoom], error)
	// GetMessagesWithChatRoomWithSd SD付きでメッセージを取得する。
	GetMessagesWithChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MessageWithChatRoom], error)
	// GetPluralMessagesWithChatRoom メッセージを取得する。
	GetPluralMessagesWithChatRoom(
		ctx context.Context,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MessageWithChatRoom], error)
	// GetPluralMessagesWithChatRoomWithSd SD付きでメッセージを取得する。
	GetPluralMessagesWithChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MessageWithChatRoom], error)
	// GetMessagesWithSender メッセージを取得する。
	GetMessagesWithSender(
		ctx context.Context,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MessageWithSender], error)
	// GetMessagesWithSenderWithSd SD付きでメッセージを取得する。
	GetMessagesWithSenderWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMessageParam,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MessageWithSender], error)
	// GetPluralMessagesWithSender メッセージを取得する。
	GetPluralMessagesWithSender(
		ctx context.Context,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MessageWithSender], error)
	// GetPluralMessagesWithSenderWithSd SD付きでメッセージを取得する。
	GetPluralMessagesWithSenderWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		order parameter.MessageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MessageWithSender], error)
	// UpdateMessage メッセージを更新する。
	UpdateMessage(
		ctx context.Context,
		messageID uuid.UUID,
		param parameter.UpdateMessageParams,
	) (entity.Message, error)
	// UpdateMessageWithSd SD付きでメッセージを更新する。
	UpdateMessageWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID,
		param parameter.UpdateMessageParams) (entity.Message, error)
}

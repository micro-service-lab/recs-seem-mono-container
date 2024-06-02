package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomDeleteMessageAction チャットルームメッセージ削除アクションを表すインターフェース。
type ChatRoomDeleteMessageAction interface {
	// CountChatRoomDeleteMessageActions チャットルームメッセージ削除アクション数を取得する。
	CountChatRoomDeleteMessageActions(
		ctx context.Context, where parameter.WhereChatRoomDeleteMessageActionParam) (int64, error)
	// CountChatRoomDeleteMessageActionsWithSd SD付きでチャットルームメッセージ削除アクション数を取得する。
	CountChatRoomDeleteMessageActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomDeleteMessageActionParam) (int64, error)
	// CreateChatRoomDeleteMessageAction チャットルームメッセージ削除アクションをメッセージ削除する。
	CreateChatRoomDeleteMessageAction(
		ctx context.Context, param parameter.CreateChatRoomDeleteMessageActionParam,
	) (entity.ChatRoomDeleteMessageAction, error)
	// CreateChatRoomDeleteMessageActionWithSd SD付きでチャットルームメッセージ削除アクションをメッセージ削除する。
	CreateChatRoomDeleteMessageActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomDeleteMessageActionParam,
	) (entity.ChatRoomDeleteMessageAction, error)
	// CreateChatRoomDeleteMessageActions チャットルームメッセージ削除アクションをメッセージ削除する。
	CreateChatRoomDeleteMessageActions(
		ctx context.Context, params []parameter.CreateChatRoomDeleteMessageActionParam) (int64, error)
	// CreateChatRoomDeleteMessageActionsWithSd SD付きでチャットルームメッセージ削除アクションをメッセージ削除する。
	CreateChatRoomDeleteMessageActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomDeleteMessageActionParam) (int64, error)
	// DeleteChatRoomDeleteMessageAction チャットルームメッセージ削除アクションを削除する。
	DeleteChatRoomDeleteMessageAction(
		ctx context.Context, chatRoomDeleteMessageActionID uuid.UUID) (int64, error)
	// DeleteChatRoomDeleteMessageActionWithSd SD付きでチャットルームメッセージ削除アクションを削除する。
	DeleteChatRoomDeleteMessageActionWithSd(
		ctx context.Context, sd Sd, chatRoomDeleteMessageActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomDeleteMessageActions チャットルームメッセージ削除アクションを複数削除する。
	PluralDeleteChatRoomDeleteMessageActions(
		ctx context.Context, chatRoomDeleteMessageActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomDeleteMessageActionsWithSd SD付きでチャットルームメッセージ削除アクションを複数削除する。
	PluralDeleteChatRoomDeleteMessageActionsWithSd(
		ctx context.Context, sd Sd, chatRoomDeleteMessageActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomDeleteMessageActionsOnChatRoom チャットルームメッセージ削除アクションを取得する。
	GetChatRoomDeleteMessageActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomDeleteMessageActionParam,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
	// GetChatRoomDeleteMessageActionsOnChatRoomWithSd SD付きでチャットルームメッセージ削除アクションを取得する。
	GetChatRoomDeleteMessageActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomDeleteMessageActionParam,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
	// GetPluralChatRoomDeleteMessageActions チャットルームメッセージ削除アクションを取得する。
	GetPluralChatRoomDeleteMessageActions(
		ctx context.Context,
		chatRoomDeleteMessageActionIDs []uuid.UUID,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
	// GetPluralChatRoomDeleteMessageActionsWithSd SD付きでチャットルームメッセージ削除アクションを取得する。
	GetPluralChatRoomDeleteMessageActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomDeleteMessageActionIDs []uuid.UUID,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
	// GetChatRoomDeleteMessageActionsOnChatRoom チャットルームに紐づくチャットルームメッセージ削除アクションを取得する。
	GetPluralChatRoomDeleteMessageActionsByChatRoomActionIDs(
		ctx context.Context,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
	// GetPluralChatRoomDeleteMessageActionsByChatRoomIDsWithSd SD付きでチャットルームメッセージ削除アクションを取得する。
	GetPluralChatRoomDeleteMessageActionsByChatRoomActionIDsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomDeleteMessageActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error)
}

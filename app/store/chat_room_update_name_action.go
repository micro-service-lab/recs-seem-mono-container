package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomUpdateNameAction チャットルーム名前更新アクションを表すインターフェース。
type ChatRoomUpdateNameAction interface {
	// CountChatRoomUpdateNameActions チャットルーム名前更新アクション数を取得する。
	CountChatRoomUpdateNameActions(
		ctx context.Context, where parameter.WhereChatRoomUpdateNameActionParam) (int64, error)
	// CountChatRoomUpdateNameActionsWithSd SD付きでチャットルーム名前更新アクション数を取得する。
	CountChatRoomUpdateNameActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomUpdateNameActionParam) (int64, error)
	// CreateChatRoomUpdateNameAction チャットルーム名前更新アクションを名前更新する。
	CreateChatRoomUpdateNameAction(
		ctx context.Context, param parameter.CreateChatRoomUpdateNameActionParam) (entity.ChatRoomUpdateNameAction, error)
	// CreateChatRoomUpdateNameActionWithSd SD付きでチャットルーム名前更新アクションを名前更新する。
	CreateChatRoomUpdateNameActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomUpdateNameActionParam,
	) (entity.ChatRoomUpdateNameAction, error)
	// CreateChatRoomUpdateNameActions チャットルーム名前更新アクションを名前更新する。
	CreateChatRoomUpdateNameActions(
		ctx context.Context, params []parameter.CreateChatRoomUpdateNameActionParam) (int64, error)
	// CreateChatRoomUpdateNameActionsWithSd SD付きでチャットルーム名前更新アクションを名前更新する。
	CreateChatRoomUpdateNameActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomUpdateNameActionParam) (int64, error)
	// DeleteChatRoomUpdateNameAction チャットルーム名前更新アクションを削除する。
	DeleteChatRoomUpdateNameAction(
		ctx context.Context, chatRoomUpdateNameActionID uuid.UUID) (int64, error)
	// DeleteChatRoomUpdateNameActionWithSd SD付きでチャットルーム名前更新アクションを削除する。
	DeleteChatRoomUpdateNameActionWithSd(
		ctx context.Context, sd Sd, chatRoomUpdateNameActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomUpdateNameActions チャットルーム名前更新アクションを複数削除する。
	PluralDeleteChatRoomUpdateNameActions(
		ctx context.Context, chatRoomUpdateNameActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomUpdateNameActionsWithSd SD付きでチャットルーム名前更新アクションを複数削除する。
	PluralDeleteChatRoomUpdateNameActionsWithSd(
		ctx context.Context, sd Sd, chatRoomUpdateNameActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomUpdateNameActionsOnChatRoom チャットルーム名前更新アクションを取得する。
	GetChatRoomUpdateNameActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomUpdateNameActionParam,
		order parameter.ChatRoomUpdateNameActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error)
	// GetChatRoomUpdateNameActionsOnChatRoomWithSd SD付きでチャットルーム名前更新アクションを取得する。
	GetChatRoomUpdateNameActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomUpdateNameActionParam,
		order parameter.ChatRoomUpdateNameActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error)
	// GetPluralChatRoomUpdateNameActions チャットルーム名前更新アクションを取得する。
	GetPluralChatRoomUpdateNameActions(
		ctx context.Context,
		chatRoomUpdateNameActionIDs []uuid.UUID,
		order parameter.ChatRoomUpdateNameActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error)
	// GetPluralChatRoomUpdateNameActionsWithSd SD付きでチャットルーム名前更新アクションを取得する。
	GetPluralChatRoomUpdateNameActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomUpdateNameActionIDs []uuid.UUID,
		order parameter.ChatRoomUpdateNameActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error)
}

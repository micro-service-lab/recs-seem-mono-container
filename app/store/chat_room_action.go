package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomAction チャットルームアクションを表すインターフェース。
type ChatRoomAction interface {
	// CountChatRoomActions チャットルームアクション数を取得する。
	CountChatRoomActions(
		ctx context.Context, where parameter.WhereChatRoomActionParam) (int64, error)
	// CountChatRoomActionsWithSd SD付きでチャットルームアクション数を取得する。
	CountChatRoomActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomActionParam) (int64, error)
	// CreateChatRoomAction チャットルームアクションを作成する。
	CreateChatRoomAction(
		ctx context.Context, param parameter.CreateChatRoomActionParam) (entity.ChatRoomAction, error)
	// CreateChatRoomActionWithSd SD付きでチャットルームアクションを作成する。
	CreateChatRoomActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomActionParam) (entity.ChatRoomAction, error)
	// CreateChatRoomActions チャットルームアクションを作成する。
	CreateChatRoomActions(
		ctx context.Context, params []parameter.CreateChatRoomActionParam) (int64, error)
	// CreateChatRoomActionsWithSd SD付きでチャットルームアクションを作成する。
	CreateChatRoomActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomActionParam) (int64, error)
	// UpdateChatRoomAction チャットルームアクションを更新する。
	UpdateChatRoomAction(
		ctx context.Context, chatRoomActionID uuid.UUID, param parameter.UpdateChatRoomActionParam,
	) (entity.ChatRoomAction, error)
	// UpdateChatRoomActionWithSd SD付きでチャットルームアクションを更新する。
	UpdateChatRoomActionWithSd(
		ctx context.Context, sd Sd, chatRoomActionID uuid.UUID, param parameter.UpdateChatRoomActionParam,
	) (entity.ChatRoomAction, error)
	// DeleteChatRoomAction チャットルームアクションを削除する。
	DeleteChatRoomAction(
		ctx context.Context, chatRoomActionID uuid.UUID) (int64, error)
	// DeleteChatRoomActionWithSd SD付きでチャットルームアクションを削除する。
	DeleteChatRoomActionWithSd(
		ctx context.Context, sd Sd, chatRoomActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomActions チャットルームアクションを複数削除する。
	PluralDeleteChatRoomActions(
		ctx context.Context, chatRoomActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomActionsWithSd SD付きでチャットルームアクションを複数削除する。
	PluralDeleteChatRoomActionsWithSd(
		ctx context.Context, sd Sd, chatRoomActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomActionsOnChatRoom チャットルームアクションを取得する。
	GetChatRoomActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomActionParam,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomAction], error)
	// GetChatRoomActionsOnChatRoomWithSd SD付きでチャットルームアクションを取得する。
	GetChatRoomActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomActionParam,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomAction], error)
	// GetPluralChatRoomActions チャットルームアクションを取得する。
	GetPluralChatRoomActions(
		ctx context.Context,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomAction], error)
	// GetPluralChatRoomActionsWithSd SD付きでチャットルームアクションを取得する。
	GetPluralChatRoomActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomAction], error)
	// GetChatRoomActionsWithDetailOnChatRoom チャットルームアクションを取得する。
	GetChatRoomActionsWithDetailOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomActionParam,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomActionWithDetail], error)
	// GetChatRoomActionsWithDetailOnChatRoomWithSd SD付きでチャットルームアクションを取得する。
	GetChatRoomActionsWithDetailOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomActionParam,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomActionWithDetail], error)
	// GetPluralChatRoomActionsWithDetail チャットルームアクションを取得する。
	GetPluralChatRoomActionsWithDetail(
		ctx context.Context,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomActionWithDetail], error)
	// GetPluralChatRoomActionsWithDetailWithSd SD付きでチャットルームアクションを取得する。
	GetPluralChatRoomActionsWithDetailWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomActionWithDetail], error)
}

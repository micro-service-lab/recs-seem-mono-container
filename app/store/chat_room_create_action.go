package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomCreateAction チャットルーム作成アクションを表すインターフェース。
type ChatRoomCreateAction interface {
	// CountChatRoomCreateActions チャットルーム作成アクション数を取得する。
	CountChatRoomCreateActions(
		ctx context.Context, where parameter.WhereChatRoomCreateActionParam) (int64, error)
	// CountChatRoomCreateActionsWithSd SD付きでチャットルーム作成アクション数を取得する。
	CountChatRoomCreateActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomCreateActionParam) (int64, error)
	// CreateChatRoomCreateAction チャットルーム作成アクションを作成する。
	CreateChatRoomCreateAction(
		ctx context.Context, param parameter.CreateChatRoomCreateActionParam) (entity.ChatRoomCreateAction, error)
	// CreateChatRoomCreateActionWithSd SD付きでチャットルーム作成アクションを作成する。
	CreateChatRoomCreateActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomCreateActionParam) (entity.ChatRoomCreateAction, error)
	// CreateChatRoomCreateActions チャットルーム作成アクションを作成する。
	CreateChatRoomCreateActions(
		ctx context.Context, params []parameter.CreateChatRoomCreateActionParam) (int64, error)
	// CreateChatRoomCreateActionsWithSd SD付きでチャットルーム作成アクションを作成する。
	CreateChatRoomCreateActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomCreateActionParam) (int64, error)
	// DeleteChatRoomCreateAction チャットルーム作成アクションを削除する。
	DeleteChatRoomCreateAction(
		ctx context.Context, chatRoomCreateActionID uuid.UUID) (int64, error)
	// DeleteChatRoomCreateActionWithSd SD付きでチャットルーム作成アクションを削除する。
	DeleteChatRoomCreateActionWithSd(
		ctx context.Context, sd Sd, chatRoomCreateActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomCreateActions チャットルーム作成アクションを複数削除する。
	PluralDeleteChatRoomCreateActions(
		ctx context.Context, chatRoomCreateActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomCreateActionsWithSd SD付きでチャットルーム作成アクションを複数削除する。
	PluralDeleteChatRoomCreateActionsWithSd(
		ctx context.Context, sd Sd, chatRoomCreateActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomCreateActionsOnChatRoom チャットルーム作成アクションを取得する。
	GetChatRoomCreateActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomCreateActionParam,
		order parameter.ChatRoomCreateActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomCreateActionOnChatRoom], error)
	// GetChatRoomCreateActionsOnChatRoomWithSd SD付きでチャットルーム作成アクションを取得する。
	GetChatRoomCreateActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomCreateActionParam,
		order parameter.ChatRoomCreateActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomCreateActionOnChatRoom], error)
	// GetPluralChatRoomCreateActions チャットルーム作成アクションを取得する。
	GetPluralChatRoomCreateActions(
		ctx context.Context,
		chatRoomCreateActionIDs []uuid.UUID,
		order parameter.ChatRoomCreateActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomCreateActionOnChatRoom], error)
	// GetPluralChatRoomCreateActionsWithSd SD付きでチャットルーム作成アクションを取得する。
	GetPluralChatRoomCreateActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomCreateActionIDs []uuid.UUID,
		order parameter.ChatRoomCreateActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomCreateActionOnChatRoom], error)
}

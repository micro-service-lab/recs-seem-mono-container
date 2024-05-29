package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomRemoveMemberAction チャットルームメンバー追放アクションを表すインターフェース。
type ChatRoomRemoveMemberAction interface {
	// CountChatRoomRemoveMemberActions チャットルームメンバー追放アクション数を取得する。
	CountChatRoomRemoveMemberActions(
		ctx context.Context, where parameter.WhereChatRoomRemoveMemberActionParam) (int64, error)
	// CountChatRoomRemoveMemberActionsWithSd SD付きでチャットルームメンバー追放アクション数を取得する。
	CountChatRoomRemoveMemberActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomRemoveMemberActionParam) (int64, error)
	// CreateChatRoomRemoveMemberAction チャットルームメンバー追放アクションをメンバー追放する。
	CreateChatRoomRemoveMemberAction(
		ctx context.Context, param parameter.CreateChatRoomRemoveMemberActionParam) (entity.ChatRoomRemoveMemberAction, error)
	// CreateChatRoomRemoveMemberActionWithSd SD付きでチャットルームメンバー追放アクションをメンバー追放する。
	CreateChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomRemoveMemberActionParam,
	) (entity.ChatRoomRemoveMemberAction, error)
	// CreateChatRoomRemoveMemberActions チャットルームメンバー追放アクションをメンバー追放する。
	CreateChatRoomRemoveMemberActions(
		ctx context.Context, params []parameter.CreateChatRoomRemoveMemberActionParam) (int64, error)
	// CreateChatRoomRemoveMemberActionsWithSd SD付きでチャットルームメンバー追放アクションをメンバー追放する。
	CreateChatRoomRemoveMemberActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomRemoveMemberActionParam) (int64, error)
	// DeleteChatRoomRemoveMemberAction チャットルームメンバー追放アクションを削除する。
	DeleteChatRoomRemoveMemberAction(
		ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomRemoveMemberActionWithSd SD付きでチャットルームメンバー追放アクションを削除する。
	DeleteChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomRemoveMemberActions チャットルームメンバー追放アクションを複数削除する。
	PluralDeleteChatRoomRemoveMemberActions(
		ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomRemoveMemberActionsWithSd SD付きでチャットルームメンバー追放アクションを複数削除する。
	PluralDeleteChatRoomRemoveMemberActionsWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomRemoveMemberActionsOnChatRoom チャットルームメンバー追放アクションを取得する。
	GetChatRoomRemoveMemberActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomRemoveMemberActionParam,
		order parameter.ChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomRemoveMemberActionOnChatRoom], error)
	// GetChatRoomRemoveMemberActionsOnChatRoomWithSd SD付きでチャットルームメンバー追放アクションを取得する。
	GetChatRoomRemoveMemberActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomRemoveMemberActionParam,
		order parameter.ChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomRemoveMemberActionOnChatRoom], error)
	// GetPluralChatRoomRemoveMemberActions チャットルームメンバー追放アクションを取得する。
	GetPluralChatRoomRemoveMemberActions(
		ctx context.Context,
		chatRoomRemoveMemberActionIDs []uuid.UUID,
		order parameter.ChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomRemoveMemberActionOnChatRoom], error)
	// GetPluralChatRoomRemoveMemberActionsWithSd SD付きでチャットルームメンバー追放アクションを取得する。
	GetPluralChatRoomRemoveMemberActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomRemoveMemberActionIDs []uuid.UUID,
		order parameter.ChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomRemoveMemberActionOnChatRoom], error)
}

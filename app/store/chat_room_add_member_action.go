package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomAddMemberAction チャットルームメンバー追加アクションを表すインターフェース。
type ChatRoomAddMemberAction interface {
	// CountChatRoomAddMemberActions チャットルームメンバー追加アクション数を取得する。
	CountChatRoomAddMemberActions(
		ctx context.Context, where parameter.WhereChatRoomAddMemberActionParam) (int64, error)
	// CountChatRoomAddMemberActionsWithSd SD付きでチャットルームメンバー追加アクション数を取得する。
	CountChatRoomAddMemberActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomAddMemberActionParam) (int64, error)
	// CreateChatRoomAddMemberAction チャットルームメンバー追加アクションをメンバー追加する。
	CreateChatRoomAddMemberAction(
		ctx context.Context, param parameter.CreateChatRoomAddMemberActionParam) (entity.ChatRoomAddMemberAction, error)
	// CreateChatRoomAddMemberActionWithSd SD付きでチャットルームメンバー追加アクションをメンバー追加する。
	CreateChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomAddMemberActionParam,
	) (entity.ChatRoomAddMemberAction, error)
	// CreateChatRoomAddMemberActions チャットルームメンバー追加アクションをメンバー追加する。
	CreateChatRoomAddMemberActions(
		ctx context.Context, params []parameter.CreateChatRoomAddMemberActionParam) (int64, error)
	// CreateChatRoomAddMemberActionsWithSd SD付きでチャットルームメンバー追加アクションをメンバー追加する。
	CreateChatRoomAddMemberActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomAddMemberActionParam) (int64, error)
	// DeleteChatRoomAddMemberAction チャットルームメンバー追加アクションを削除する。
	DeleteChatRoomAddMemberAction(
		ctx context.Context, chatRoomAddMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomAddMemberActionWithSd SD付きでチャットルームメンバー追加アクションを削除する。
	DeleteChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomAddMemberActions チャットルームメンバー追加アクションを複数削除する。
	PluralDeleteChatRoomAddMemberActions(
		ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomAddMemberActionsWithSd SD付きでチャットルームメンバー追加アクションを複数削除する。
	PluralDeleteChatRoomAddMemberActionsWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomAddMemberActionsOnChatRoom チャットルームメンバー追加アクションを取得する。
	GetChatRoomAddMemberActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomAddMemberActionParam,
		order parameter.ChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomAddMemberActionOnChatRoom], error)
	// GetChatRoomAddMemberActionsOnChatRoomWithSd SD付きでチャットルームメンバー追加アクションを取得する。
	GetChatRoomAddMemberActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomAddMemberActionParam,
		order parameter.ChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomAddMemberActionOnChatRoom], error)
	// GetPluralChatRoomAddMemberActions チャットルームメンバー追加アクションを取得する。
	GetPluralChatRoomAddMemberActions(
		ctx context.Context,
		chatRoomAddMemberActionIDs []uuid.UUID,
		order parameter.ChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomAddMemberActionOnChatRoom], error)
	// GetPluralChatRoomAddMemberActionsWithSd SD付きでチャットルームメンバー追加アクションを取得する。
	GetPluralChatRoomAddMemberActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomAddMemberActionIDs []uuid.UUID,
		order parameter.ChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomAddMemberActionOnChatRoom], error)
}

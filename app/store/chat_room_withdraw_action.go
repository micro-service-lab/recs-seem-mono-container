package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomWithdrawAction チャットルームメンバー脱退アクションを表すインターフェース。
type ChatRoomWithdrawAction interface {
	// CountChatRoomWithdrawActions チャットルームメンバー脱退アクション数を取得する。
	CountChatRoomWithdrawActions(
		ctx context.Context, where parameter.WhereChatRoomWithdrawActionParam) (int64, error)
	// CountChatRoomWithdrawActionsWithSd SD付きでチャットルームメンバー脱退アクション数を取得する。
	CountChatRoomWithdrawActionsWithSd(
		ctx context.Context, sd Sd, where parameter.WhereChatRoomWithdrawActionParam) (int64, error)
	// CreateChatRoomWithdrawAction チャットルームメンバー脱退アクションをメンバー脱退する。
	CreateChatRoomWithdrawAction(
		ctx context.Context, param parameter.CreateChatRoomWithdrawActionParam) (entity.ChatRoomWithdrawAction, error)
	// CreateChatRoomWithdrawActionWithSd SD付きでチャットルームメンバー脱退アクションをメンバー脱退する。
	CreateChatRoomWithdrawActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomWithdrawActionParam,
	) (entity.ChatRoomWithdrawAction, error)
	// CreateChatRoomWithdrawActions チャットルームメンバー脱退アクションをメンバー脱退する。
	CreateChatRoomWithdrawActions(
		ctx context.Context, params []parameter.CreateChatRoomWithdrawActionParam) (int64, error)
	// CreateChatRoomWithdrawActionsWithSd SD付きでチャットルームメンバー脱退アクションをメンバー脱退する。
	CreateChatRoomWithdrawActionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomWithdrawActionParam) (int64, error)
	// DeleteChatRoomWithdrawAction チャットルームメンバー脱退アクションを削除する。
	DeleteChatRoomWithdrawAction(
		ctx context.Context, chatRoomWithdrawActionID uuid.UUID) (int64, error)
	// DeleteChatRoomWithdrawActionWithSd SD付きでチャットルームメンバー脱退アクションを削除する。
	DeleteChatRoomWithdrawActionWithSd(
		ctx context.Context, sd Sd, chatRoomWithdrawActionID uuid.UUID) (int64, error)
	// PluralDeleteChatRoomWithdrawActions チャットルームメンバー脱退アクションを複数削除する。
	PluralDeleteChatRoomWithdrawActions(
		ctx context.Context, chatRoomWithdrawActionIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomWithdrawActionsWithSd SD付きでチャットルームメンバー脱退アクションを複数削除する。
	PluralDeleteChatRoomWithdrawActionsWithSd(
		ctx context.Context, sd Sd, chatRoomWithdrawActionIDs []uuid.UUID) (int64, error)
	// GetChatRoomWithdrawActionsOnChatRoom チャットルームメンバー脱退アクションを取得する。
	GetChatRoomWithdrawActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomWithdrawActionParam,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
	// GetChatRoomWithdrawActionsOnChatRoomWithSd SD付きでチャットルームメンバー脱退アクションを取得する。
	GetChatRoomWithdrawActionsOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereChatRoomWithdrawActionParam,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
	// GetPluralChatRoomWithdrawActions チャットルームメンバー脱退アクションを取得する。
	GetPluralChatRoomWithdrawActions(
		ctx context.Context,
		chatRoomWithdrawActionIDs []uuid.UUID,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
	// GetPluralChatRoomWithdrawActionsWithSd SD付きでチャットルームメンバー脱退アクションを取得する。
	GetPluralChatRoomWithdrawActionsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomWithdrawActionIDs []uuid.UUID,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
	// GetPluralChatRoomWithdrawActionsByChatRoomIDs チャットルームメンバー脱退アクションを取得する。
	GetPluralChatRoomWithdrawActionsByChatRoomActionIDs(
		ctx context.Context,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
	// GetPluralChatRoomWithdrawActionsByChatRoomIDsWithSd SD付きでチャットルームメンバー脱退アクションを取得する。
	GetPluralChatRoomWithdrawActionsByChatRoomActionIDsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomActionIDs []uuid.UUID,
		order parameter.ChatRoomWithdrawActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomWithdrawActionWithMember], error)
}

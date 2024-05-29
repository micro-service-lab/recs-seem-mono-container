package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomActionType チャットルームアクションタイプを表すインターフェース。
type ChatRoomActionType interface {
	// CountChatRoomActionTypes チャットルームアクションタイプ数を取得する。
	CountChatRoomActionTypes(ctx context.Context, where parameter.WhereChatRoomActionTypeParam) (int64, error)
	// CountChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプ数を取得する。
	CountChatRoomActionTypesWithSd(ctx context.Context, sd Sd, where parameter.WhereChatRoomActionTypeParam) (int64, error)
	// CreateChatRoomActionType チャットルームアクションタイプを作成する。
	CreateChatRoomActionType(
		ctx context.Context, param parameter.CreateChatRoomActionTypeParam) (entity.ChatRoomActionType, error)
	// CreateChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを作成する。
	CreateChatRoomActionTypeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomActionTypeParam) (entity.ChatRoomActionType, error)
	// CreateChatRoomActionTypes チャットルームアクションタイプを作成する。
	CreateChatRoomActionTypes(ctx context.Context, params []parameter.CreateChatRoomActionTypeParam) (int64, error)
	// CreateChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを作成する。
	CreateChatRoomActionTypesWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomActionTypeParam) (int64, error)
	// DeleteChatRoomActionType チャットルームアクションタイプを削除する。
	DeleteChatRoomActionType(ctx context.Context, recordTypeID uuid.UUID) (int64, error)
	// DeleteChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを削除する。
	DeleteChatRoomActionTypeWithSd(ctx context.Context, sd Sd, recordTypeID uuid.UUID) (int64, error)
	// DeleteChatRoomActionTypeByKey チャットルームアクションタイプを削除する。
	DeleteChatRoomActionTypeByKey(ctx context.Context, key string) (int64, error)
	// DeleteChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを削除する。
	DeleteChatRoomActionTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeleteChatRoomActionTypes チャットルームアクションタイプを複数削除する。
	PluralDeleteChatRoomActionTypes(ctx context.Context, recordTypeIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを複数削除する。
	PluralDeleteChatRoomActionTypesWithSd(ctx context.Context, sd Sd, recordTypeIDs []uuid.UUID) (int64, error)
	// FindChatRoomActionTypeByID チャットルームアクションタイプを取得する。
	FindChatRoomActionTypeByID(ctx context.Context, recordTypeID uuid.UUID) (entity.ChatRoomActionType, error)
	// FindChatRoomActionTypeByIDWithSd SD付きでチャットルームアクションタイプを取得する。
	FindChatRoomActionTypeByIDWithSd(ctx context.Context, sd Sd, recordTypeID uuid.UUID) (entity.ChatRoomActionType, error)
	// FindChatRoomActionTypeByKey チャットルームアクションタイプを取得する。
	FindChatRoomActionTypeByKey(ctx context.Context, key string) (entity.ChatRoomActionType, error)
	// FindChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを取得する。
	FindChatRoomActionTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.ChatRoomActionType, error)
	// GetChatRoomActionTypes チャットルームアクションタイプを取得する。
	GetChatRoomActionTypes(
		ctx context.Context,
		where parameter.WhereChatRoomActionTypeParam,
		order parameter.ChatRoomActionTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomActionType], error)
	// GetChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを取得する。
	GetChatRoomActionTypesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereChatRoomActionTypeParam,
		order parameter.ChatRoomActionTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomActionType], error)
	// GetPluralChatRoomActionTypes チャットルームアクションタイプを取得する。
	GetPluralChatRoomActionTypes(
		ctx context.Context,
		recordTypeIDs []uuid.UUID,
		order parameter.ChatRoomActionTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomActionType], error)
	// GetPluralChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを取得する。
	GetPluralChatRoomActionTypesWithSd(
		ctx context.Context,
		sd Sd,
		recordTypeIDs []uuid.UUID,
		order parameter.ChatRoomActionTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoomActionType], error)
	// UpdateChatRoomActionType チャットルームアクションタイプを更新する。
	UpdateChatRoomActionType(
		ctx context.Context,
		recordTypeID uuid.UUID,
		param parameter.UpdateChatRoomActionTypeParams,
	) (entity.ChatRoomActionType, error)
	// UpdateChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを更新する。
	UpdateChatRoomActionTypeWithSd(
		ctx context.Context, sd Sd, recordTypeID uuid.UUID,
		param parameter.UpdateChatRoomActionTypeParams) (entity.ChatRoomActionType, error)
	// UpdateChatRoomActionTypeByKey チャットルームアクションタイプを更新する。
	UpdateChatRoomActionTypeByKey(
		ctx context.Context, key string,
		param parameter.UpdateChatRoomActionTypeByKeyParams) (entity.ChatRoomActionType, error)
	// UpdateChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを更新する。
	UpdateChatRoomActionTypeByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdateChatRoomActionTypeByKeyParams,
	) (entity.ChatRoomActionType, error)
}

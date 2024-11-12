package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// EventType イベントタイプを表すインターフェース。
type EventType interface {
	// CountEventTypes イベントタイプ数を取得する。
	CountEventTypes(ctx context.Context, where parameter.WhereEventTypeParam) (int64, error)
	// CountEventTypesWithSd SD付きでイベントタイプ数を取得する。
	CountEventTypesWithSd(ctx context.Context, sd Sd, where parameter.WhereEventTypeParam) (int64, error)
	// CreateEventType イベントタイプを作成する。
	CreateEventType(ctx context.Context, param parameter.CreateEventTypeParam) (entity.EventType, error)
	// CreateEventTypeWithSd SD付きでイベントタイプを作成する。
	CreateEventTypeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateEventTypeParam) (entity.EventType, error)
	// CreateEventTypes イベントタイプを作成する。
	CreateEventTypes(ctx context.Context, params []parameter.CreateEventTypeParam) (int64, error)
	// CreateEventTypesWithSd SD付きでイベントタイプを作成する。
	CreateEventTypesWithSd(ctx context.Context, sd Sd, params []parameter.CreateEventTypeParam) (int64, error)
	// DeleteEventType イベントタイプを削除する。
	DeleteEventType(ctx context.Context, eventTypeID uuid.UUID) (int64, error)
	// DeleteEventTypeWithSd SD付きでイベントタイプを削除する。
	DeleteEventTypeWithSd(ctx context.Context, sd Sd, eventTypeID uuid.UUID) (int64, error)
	// DeleteEventTypeByKey イベントタイプを削除する。
	DeleteEventTypeByKey(ctx context.Context, key string) (int64, error)
	// DeleteEventTypeByKeyWithSd SD付きでイベントタイプを削除する。
	DeleteEventTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeleteEventTypes イベントタイプを複数削除する。
	PluralDeleteEventTypes(ctx context.Context, eventTypeIDs []uuid.UUID) (int64, error)
	// PluralDeleteEventTypesWithSd SD付きでイベントタイプを複数削除する。
	PluralDeleteEventTypesWithSd(ctx context.Context, sd Sd, eventTypeIDs []uuid.UUID) (int64, error)
	// FindEventTypeByID イベントタイプを取得する。
	FindEventTypeByID(ctx context.Context, eventTypeID uuid.UUID) (entity.EventType, error)
	// FindEventTypeByIDWithSd SD付きでイベントタイプを取得する。
	FindEventTypeByIDWithSd(ctx context.Context, sd Sd, eventTypeID uuid.UUID) (entity.EventType, error)
	// FindEventTypeByKey イベントタイプを取得する。
	FindEventTypeByKey(ctx context.Context, key string) (entity.EventType, error)
	// FindEventTypeByKeyWithSd SD付きでイベントタイプを取得する。
	FindEventTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.EventType, error)
	// GetEventTypes イベントタイプを取得する。
	GetEventTypes(
		ctx context.Context,
		where parameter.WhereEventTypeParam,
		order parameter.EventTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.EventType], error)
	// GetEventTypesWithSd SD付きでイベントタイプを取得する。
	GetEventTypesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereEventTypeParam,
		order parameter.EventTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.EventType], error)
	// GetPluralEventTypes イベントタイプを取得する。
	GetPluralEventTypes(
		ctx context.Context,
		eventTypeIDs []uuid.UUID,
		order parameter.EventTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.EventType], error)
	// GetPluralEventTypesWithSd SD付きでイベントタイプを取得する。
	GetPluralEventTypesWithSd(
		ctx context.Context,
		sd Sd,
		eventTypeIDs []uuid.UUID,
		order parameter.EventTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.EventType], error)
	// UpdateEventType イベントタイプを更新する。
	UpdateEventType(
		ctx context.Context,
		eventTypeID uuid.UUID,
		param parameter.UpdateEventTypeParams,
	) (entity.EventType, error)
	// UpdateEventTypeWithSd SD付きでイベントタイプを更新する。
	UpdateEventTypeWithSd(
		ctx context.Context, sd Sd, eventTypeID uuid.UUID,
		param parameter.UpdateEventTypeParams) (entity.EventType, error)
	// UpdateEventTypeByKey イベントタイプを更新する。
	UpdateEventTypeByKey(
		ctx context.Context, key string, param parameter.UpdateEventTypeByKeyParams) (entity.EventType, error)
	// UpdateEventTypeByKeyWithSd SD付きでイベントタイプを更新する。
	UpdateEventTypeByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdateEventTypeByKeyParams,
	) (entity.EventType, error)
}

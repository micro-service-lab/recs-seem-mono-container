package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// EventTypeKey イベントタイプキー。
type EventTypeKey string

const (
	// EventTypeKeyMeeting 会議。
	EventTypeKeyMeeting EventTypeKey = "meeting"
	// EventTypeKeyJournalClub 輪講。
	EventTypeKeyJournalClub EventTypeKey = "journal_club"
	// EventTypeKeyHoliday 休日。
	EventTypeKeyHoliday EventTypeKey = "holiday"
	// EventTypeKeyOther その他。
	EventTypeKeyOther EventTypeKey = "other"
)

// EventType イベントタイプ。
type EventType struct {
	Key   string
	Name  string
	Color string
}

// EventTypes イベントタイプ一覧。
var EventTypes = []EventType{
	{Key: string(EventTypeKeyMeeting), Name: "会議", Color: "#FFB866"},
	{Key: string(EventTypeKeyJournalClub), Name: "輪講", Color: "#ADFF66"},
	{Key: string(EventTypeKeyHoliday), Name: "休日", Color: "#FF4D4D"},
	{Key: string(EventTypeKeyOther), Name: "その他", Color: "#66B3FF"},
}

// ManageEventType イベントタイプ管理サービス。
type ManageEventType struct {
	DB store.Store
}

// CreateEventType イベントタイプを作成する。
func (m *ManageEventType) CreateEventType(
	ctx context.Context,
	name, key, color string,
) (entity.EventType, error) {
	p := parameter.CreateEventTypeParam{
		Name:  name,
		Key:   key,
		Color: color,
	}
	e, err := m.DB.CreateEventType(ctx, p)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to create event type: %w", err)
	}
	return e, nil
}

// CreateEventTypes イベントタイプを複数作成する。
func (m *ManageEventType) CreateEventTypes(
	ctx context.Context, ps []parameter.CreateEventTypeParam,
) (int64, error) {
	es, err := m.DB.CreateEventTypes(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create event types: %w", err)
	}
	return es, nil
}

// UpdateEventType イベントタイプを更新する。
func (m *ManageEventType) UpdateEventType(
	ctx context.Context, id uuid.UUID, name, key, color string,
) (entity.EventType, error) {
	p := parameter.UpdateEventTypeParams{
		Name:  name,
		Key:   key,
		Color: color,
	}
	e, err := m.DB.UpdateEventType(ctx, id, p)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	return e, nil
}

// DeleteEventType イベントタイプを削除する。
func (m *ManageEventType) DeleteEventType(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteEventType(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete event type: %w", err)
	}
	return c, nil
}

// PluralDeleteEventTypes イベントタイプを複数削除する。
func (m *ManageEventType) PluralDeleteEventTypes(ctx context.Context, ids []uuid.UUID) (int64, error) {
	c, err := m.DB.PluralDeleteEventTypes(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete event types: %w", err)
	}
	return c, nil
}

// FindEventTypeByID イベントタイプをIDで取得する。
func (m *ManageEventType) FindEventTypeByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.EventType, error) {
	e, err := m.DB.FindEventTypeByID(ctx, id)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type by id: %w", err)
	}
	return e, nil
}

// FindEventTypeByKey イベントタイプをキーで取得する。
func (m *ManageEventType) FindEventTypeByKey(ctx context.Context, key string) (entity.EventType, error) {
	e, err := m.DB.FindEventTypeByKey(ctx, key)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type by key: %w", err)
	}
	return e, nil
}

// GetEventTypes イベントタイプを取得する。
func (m *ManageEventType) GetEventTypes(
	ctx context.Context,
	whereSearchName string,
	order parameter.EventTypeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.EventType], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereEventTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetEventTypes(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get event types: %w", err)
	}
	return r, nil
}

// GetEventTypesCount イベントタイプの数を取得する。
func (m *ManageEventType) GetEventTypesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereEventTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountEventTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get event types count: %w", err)
	}
	return c, nil
}

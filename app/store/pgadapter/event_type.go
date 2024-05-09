package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func countEventTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereEventTypeParam,
) (int64, error) {
	p := query.CountEventTypesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountEventTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count event type: %w", err)
	}
	return c, nil
}

// CountEventTypes イベントタイプ数を取得する。
func (a *PgAdapter) CountEventTypes(ctx context.Context, where parameter.WhereEventTypeParam) (int64, error) {
	c, err := countEventTypes(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count event type: %w", err)
	}
	return c, nil
}

// CountEventTypesWithSd SD付きでイベントタイプ数を取得する。
func (a *PgAdapter) CountEventTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereEventTypeParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countEventTypes(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count event type: %w", err)
	}
	return c, nil
}

func createEventType(
	ctx context.Context, qtx *query.Queries, param parameter.CreateEventTypeParam,
) (entity.EventType, error) {
	p := query.CreateEventTypeParams{
		Name:  param.Name,
		Key:   param.Key,
		Color: param.Color,
	}
	e, err := qtx.CreateEventType(ctx, p)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to create event type: %w", err)
	}
	entity := entity.EventType{
		EventTypeID: e.EventTypeID,
		Name:        e.Name,
		Key:         e.Key,
		Color:       e.Color,
	}
	return entity, nil
}

// CreateEventType イベントタイプを作成する。
func (a *PgAdapter) CreateEventType(
	ctx context.Context, param parameter.CreateEventTypeParam,
) (entity.EventType, error) {
	e, err := createEventType(ctx, a.query, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to create event type: %w", err)
	}
	return e, nil
}

// CreateEventTypeWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreateEventTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateEventTypeParam,
) (entity.EventType, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	e, err := createEventType(ctx, qtx, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to create event type: %w", err)
	}
	return e, nil
}

func createEventTypes(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateEventTypeParam,
) (int64, error) {
	p := make([]query.CreateEventTypesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateEventTypesParams{
			Name:  param.Name,
			Key:   param.Key,
			Color: param.Color,
		}
	}
	c, err := qtx.CreateEventTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create event types: %w", err)
	}
	return c, nil
}

// CreateEventTypes イベントタイプを作成する。
func (a *PgAdapter) CreateEventTypes(
	ctx context.Context, params []parameter.CreateEventTypeParam,
) (int64, error) {
	c, err := createEventTypes(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create event types: %w", err)
	}
	return c, nil
}

// CreateEventTypesWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreateEventTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateEventTypeParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createEventTypes(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create event types: %w", err)
	}
	return c, nil
}

func deleteEventType(ctx context.Context, qtx *query.Queries, eventTypeID uuid.UUID) error {
	err := qtx.DeleteEventType(ctx, eventTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

// DeleteEventType イベントタイプを削除する。
func (a *PgAdapter) DeleteEventType(ctx context.Context, eventTypeID uuid.UUID) error {
	err := deleteEventType(ctx, a.query, eventTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

// DeleteEventTypeWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteEventType(ctx, qtx, eventTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

func deleteEventTypeByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeleteEventTypeByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

// DeleteEventTypeByKey イベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeByKey(ctx context.Context, key string) error {
	err := deleteEventTypeByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

// DeleteEventTypeByKeyWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteEventTypeByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}
	return nil
}

func pluralDeleteEventTypes(ctx context.Context, qtx *query.Queries, eventTypeIDs []uuid.UUID) error {
	err := qtx.PluralDeleteEventTypes(ctx, eventTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete event types: %w", err)
	}
	return nil
}

// PluralDeleteEventTypes イベントタイプを複数削除する。
func (a *PgAdapter) PluralDeleteEventTypes(ctx context.Context, eventTypeIDs []uuid.UUID) error {
	err := pluralDeleteEventTypes(ctx, a.query, eventTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete event types: %w", err)
	}
	return nil
}

// PluralDeleteEventTypesWithSd SD付きでイベントタイプを複数削除する。
func (a *PgAdapter) PluralDeleteEventTypesWithSd(
	ctx context.Context, sd store.Sd, eventTypeIDs []uuid.UUID,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeleteEventTypes(ctx, qtx, eventTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete event types: %w", err)
	}
	return nil
}

func findEventTypeByID(
	ctx context.Context, qtx *query.Queries, eventTypeID uuid.UUID,
) (entity.EventType, error) {
	e, err := qtx.FindEventTypeByID(ctx, eventTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, store.ErrDataNoRecord
		}
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	entity := entity.EventType{
		EventTypeID: e.EventTypeID,
		Name:        e.Name,
		Key:         e.Key,
		Color:       e.Color,
	}
	return entity, nil
}

// FindEventTypeByID イベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByID(
	ctx context.Context, eventTypeID uuid.UUID,
) (entity.EventType, error) {
	e, err := findEventTypeByID(ctx, a.query, eventTypeID)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	return e, nil
}

// FindEventTypeByIDWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByIDWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID,
) (entity.EventType, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	e, err := findEventTypeByID(ctx, qtx, eventTypeID)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	return e, nil
}

func findEventTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.EventType, error) {
	e, err := qtx.FindEventTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, store.ErrDataNoRecord
		}
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	entity := entity.EventType{
		EventTypeID: e.EventTypeID,
		Name:        e.Name,
		Key:         e.Key,
		Color:       e.Color,
	}
	return entity, nil
}

// FindEventTypeByKey イベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByKey(ctx context.Context, key string) (entity.EventType, error) {
	e, err := findEventTypeByKey(ctx, a.query, key)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	return e, nil
}

// FindEventTypeByKeyWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.EventType, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	e, err := findEventTypeByKey(ctx, qtx, key)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to find event type: %w", err)
	}
	return e, nil
}

// EventTypeCursor is a cursor for EventType.
type EventTypeCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getEventTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereEventTypeParam,
	order parameter.EventTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.EventType], error) {
	eConvFunc := func(e query.EventType) (entity.EventType, error) {
		return entity.EventType{
			EventTypeID: e.EventTypeID,
			Name:        e.Name,
			Key:         e.Key,
			Color:       e.Color,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountEventTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountEventTypes(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count event types: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.EventType, error) {
		p := query.GetEventTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetEventTypes(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.EventType{}, nil
			}
			return nil, fmt.Errorf("failed to get event types: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.EventType, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.EventTypeNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetEventTypesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetEventTypesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get event types: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.EventType, error) {
		p := query.GetEventTypesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetEventTypesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get event types: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.EventType) (entity.Int, any) {
		switch subCursor {
		case parameter.EventTypeDefaultCursorKey:
			return entity.Int(e.MEventTypesPkey), nil
		case parameter.EventTypeNameCursorKey:
			return entity.Int(e.MEventTypesPkey), e.Name
		}
		return entity.Int(e.MEventTypesPkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get event types: %w", err)
	}
	return res, nil
}

// GetEventTypes イベントタイプを取得する。
func (a *PgAdapter) GetEventTypes(
	ctx context.Context,
	where parameter.WhereEventTypeParam,
	order parameter.EventTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.EventType], error) {
	r, err := getEventTypes(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get event types: %w", err)
	}
	return r, nil
}

// GetEventTypesWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) GetEventTypesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereEventTypeParam,
	order parameter.EventTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.EventType], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.EventType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getEventTypes(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get event types: %w", err)
	}
	return r, nil
}

func getPluralEventTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.EventType], error) {
	p := query.GetPluralEventTypesParams{
		EventTypeIds: ids,
		Offset:       int32(np.Offset.Int64),
		Limit:        int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralEventTypes(ctx, p)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get plural event types: %w", err)
	}
	entities := make([]entity.EventType, len(e))
	for i, v := range e {
		entities[i] = entity.EventType{
			EventTypeID: v.EventTypeID,
			Name:        v.Name,
			Key:         v.Key,
			Color:       v.Color,
		}
	}
	return store.ListResult[entity.EventType]{Data: entities}, nil
}

// GetPluralEventTypes イベントタイプを取得する。
func (a *PgAdapter) GetPluralEventTypes(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.EventType], error) {
	r, err := getPluralEventTypes(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get plural event types: %w", err)
	}
	return r, nil
}

// GetPluralEventTypesWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) GetPluralEventTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.EventType], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.EventType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralEventTypes(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.EventType]{}, fmt.Errorf("failed to get plural event types: %w", err)
	}
	return r, nil
}

func updateEventType(
	ctx context.Context, qtx *query.Queries, eventTypeID uuid.UUID, param parameter.UpdateEventTypeParams,
) (entity.EventType, error) {
	p := query.UpdateEventTypeParams{
		EventTypeID: eventTypeID,
		Name:        param.Name,
		Key:         param.Key,
		Color:       param.Color,
	}
	e, err := qtx.UpdateEventType(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, store.ErrDataNoRecord
		}
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	entity := entity.EventType{
		EventTypeID: e.EventTypeID,
		Name:        e.Name,
		Key:         e.Key,
		Color:       e.Color,
	}
	return entity, nil
}

// UpdateEventType イベントタイプを更新する。
func (a *PgAdapter) UpdateEventType(
	ctx context.Context, eventTypeID uuid.UUID, param parameter.UpdateEventTypeParams,
) (entity.EventType, error) {
	e, err := updateEventType(ctx, a.query, eventTypeID, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	return e, nil
}

// UpdateEventTypeWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdateEventTypeWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID, param parameter.UpdateEventTypeParams,
) (entity.EventType, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateEventType(ctx, qtx, eventTypeID, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	return e, nil
}

func updateEventTypeByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateEventTypeByKeyParams,
) (entity.EventType, error) {
	p := query.UpdateEventTypeByKeyParams{
		Key:   key,
		Name:  param.Name,
		Color: param.Color,
	}
	e, err := qtx.UpdateEventTypeByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, store.ErrDataNoRecord
		}
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	entity := entity.EventType{
		EventTypeID: e.EventTypeID,
		Name:        e.Name,
		Key:         e.Key,
		Color:       e.Color,
	}
	return entity, nil
}

// UpdateEventTypeByKey イベントタイプを更新する。
func (a *PgAdapter) UpdateEventTypeByKey(
	ctx context.Context, key string, param parameter.UpdateEventTypeByKeyParams,
) (entity.EventType, error) {
	e, err := updateEventTypeByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	return e, nil
}

// UpdateEventTypeByKeyWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdateEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateEventTypeByKeyParams,
) (entity.EventType, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateEventTypeByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.EventType{}, fmt.Errorf("failed to update event type: %w", err)
	}
	return e, nil
}

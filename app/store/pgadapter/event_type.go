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
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
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
	return countEventTypes(ctx, a.query, where)
}

// CountEventTypesWithSd SD付きでイベントタイプ数を取得する。
func (a *PgAdapter) CountEventTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereEventTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countEventTypes(ctx, qtx, where)
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
	return createEventType(ctx, a.query, param)
}

// CreateEventTypeWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreateEventTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateEventTypeParam,
) (entity.EventType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	return createEventType(ctx, qtx, param)
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
	return createEventTypes(ctx, a.query, params)
}

// CreateEventTypesWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreateEventTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateEventTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createEventTypes(ctx, qtx, params)
}

func deleteEventType(ctx context.Context, qtx *query.Queries, eventTypeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteEventType(ctx, eventTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete event type: %w", err)
	}
	return c, nil
}

// DeleteEventType イベントタイプを削除する。
func (a *PgAdapter) DeleteEventType(ctx context.Context, eventTypeID uuid.UUID) (int64, error) {
	return deleteEventType(ctx, a.query, eventTypeID)
}

// DeleteEventTypeWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteEventType(ctx, qtx, eventTypeID)
}

func deleteEventTypeByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteEventTypeByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete event type: %w", err)
	}
	return c, nil
}

// DeleteEventTypeByKey イベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeByKey(ctx context.Context, key string) (int64, error) {
	return deleteEventTypeByKey(ctx, a.query, key)
}

// DeleteEventTypeByKeyWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeleteEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteEventTypeByKey(ctx, qtx, key)
}

func pluralDeleteEventTypes(ctx context.Context, qtx *query.Queries, eventTypeIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteEventTypes(ctx, eventTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete event types: %w", err)
	}
	return c, nil
}

// PluralDeleteEventTypes イベントタイプを複数削除する。
func (a *PgAdapter) PluralDeleteEventTypes(ctx context.Context, eventTypeIDs []uuid.UUID) (int64, error) {
	return pluralDeleteEventTypes(ctx, a.query, eventTypeIDs)
}

// PluralDeleteEventTypesWithSd SD付きでイベントタイプを複数削除する。
func (a *PgAdapter) PluralDeleteEventTypesWithSd(
	ctx context.Context, sd store.Sd, eventTypeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteEventTypes(ctx, qtx, eventTypeIDs)
}

func findEventTypeByID(
	ctx context.Context, qtx *query.Queries, eventTypeID uuid.UUID,
) (entity.EventType, error) {
	e, err := qtx.FindEventTypeByID(ctx, eventTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, errhandle.NewModelNotFoundError("event type")
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
	return findEventTypeByID(ctx, a.query, eventTypeID)
}

// FindEventTypeByIDWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByIDWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID,
) (entity.EventType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	return findEventTypeByID(ctx, qtx, eventTypeID)
}

func findEventTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.EventType, error) {
	e, err := qtx.FindEventTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, errhandle.NewModelNotFoundError("event type")
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
	return findEventTypeByKey(ctx, a.query, key)
}

// FindEventTypeByKeyWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.EventType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	return findEventTypeByKey(ctx, qtx, key)
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
	return getEventTypes(ctx, a.query, where, order, np, cp, wc)
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
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.EventType]{}, store.ErrNotFoundDescriptor
	}
	return getEventTypes(ctx, qtx, where, order, np, cp, wc)
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
	return getPluralEventTypes(ctx, a.query, ids, np)
}

// GetPluralEventTypesWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) GetPluralEventTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.EventType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.EventType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralEventTypes(ctx, qtx, ids, np)
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
			return entity.EventType{}, errhandle.NewModelNotFoundError("event type")
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
	return updateEventType(ctx, a.query, eventTypeID, param)
}

// UpdateEventTypeWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdateEventTypeWithSd(
	ctx context.Context, sd store.Sd, eventTypeID uuid.UUID, param parameter.UpdateEventTypeParams,
) (entity.EventType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	return updateEventType(ctx, qtx, eventTypeID, param)
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
			return entity.EventType{}, errhandle.NewModelNotFoundError("event type")
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
	return updateEventTypeByKey(ctx, a.query, key, param)
}

// UpdateEventTypeByKeyWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdateEventTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateEventTypeByKeyParams,
) (entity.EventType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.EventType{}, store.ErrNotFoundDescriptor
	}
	return updateEventTypeByKey(ctx, qtx, key, param)
}

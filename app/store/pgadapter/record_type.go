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

func countRecordTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereRecordTypeParam,
) (int64, error) {
	p := query.CountRecordTypesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountRecordTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count record type: %w", err)
	}
	return c, nil
}

// CountRecordTypes 議事録タイプ数を取得する。
func (a *PgAdapter) CountRecordTypes(ctx context.Context, where parameter.WhereRecordTypeParam) (int64, error) {
	c, err := countRecordTypes(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count record type: %w", err)
	}
	return c, nil
}

// CountRecordTypesWithSd SD付きで議事録タイプ数を取得する。
func (a *PgAdapter) CountRecordTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereRecordTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countRecordTypes(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count record type: %w", err)
	}
	return c, nil
}

func createRecordType(
	ctx context.Context, qtx *query.Queries, param parameter.CreateRecordTypeParam,
) (entity.RecordType, error) {
	p := query.CreateRecordTypeParams{
		Name: param.Name,
		Key:  param.Key,
	}
	e, err := qtx.CreateRecordType(ctx, p)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to create record type: %w", err)
	}
	entity := entity.RecordType{
		RecordTypeID: e.RecordTypeID,
		Name:         e.Name,
		Key:          e.Key,
	}
	return entity, nil
}

// CreateRecordType 議事録タイプを作成する。
func (a *PgAdapter) CreateRecordType(
	ctx context.Context, param parameter.CreateRecordTypeParam,
) (entity.RecordType, error) {
	e, err := createRecordType(ctx, a.query, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to create record type: %w", err)
	}
	return e, nil
}

// CreateRecordTypeWithSd SD付きで議事録タイプを作成する。
func (a *PgAdapter) CreateRecordTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateRecordTypeParam,
) (entity.RecordType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RecordType{}, store.ErrNotFoundDescriptor
	}
	e, err := createRecordType(ctx, qtx, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to create record type: %w", err)
	}
	return e, nil
}

func createRecordTypes(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateRecordTypeParam,
) (int64, error) {
	p := make([]query.CreateRecordTypesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateRecordTypesParams{
			Name: param.Name,
			Key:  param.Key,
		}
	}
	c, err := qtx.CreateRecordTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create record types: %w", err)
	}
	return c, nil
}

// CreateRecordTypes 議事録タイプを作成する。
func (a *PgAdapter) CreateRecordTypes(
	ctx context.Context, params []parameter.CreateRecordTypeParam,
) (int64, error) {
	c, err := createRecordTypes(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create record types: %w", err)
	}
	return c, nil
}

// CreateRecordTypesWithSd SD付きで議事録タイプを作成する。
func (a *PgAdapter) CreateRecordTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateRecordTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createRecordTypes(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create record types: %w", err)
	}
	return c, nil
}

func deleteRecordType(ctx context.Context, qtx *query.Queries, recordTypeID uuid.UUID) error {
	err := qtx.DeleteRecordType(ctx, recordTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

// DeleteRecordType 議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordType(ctx context.Context, recordTypeID uuid.UUID) error {
	err := deleteRecordType(ctx, a.query, recordTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

// DeleteRecordTypeWithSd SD付きで議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeWithSd(
	ctx context.Context, sd store.Sd, recordTypeID uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteRecordType(ctx, qtx, recordTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

func deleteRecordTypeByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeleteRecordTypeByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

// DeleteRecordTypeByKey 議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeByKey(ctx context.Context, key string) error {
	err := deleteRecordTypeByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

// DeleteRecordTypeByKeyWithSd SD付きで議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteRecordTypeByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete record type: %w", err)
	}
	return nil
}

func pluralDeleteRecordTypes(ctx context.Context, qtx *query.Queries, recordTypeIDs []uuid.UUID) error {
	err := qtx.PluralDeleteRecordTypes(ctx, recordTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete record types: %w", err)
	}
	return nil
}

// PluralDeleteRecordTypes 議事録タイプを複数削除する。
func (a *PgAdapter) PluralDeleteRecordTypes(ctx context.Context, recordTypeIDs []uuid.UUID) error {
	err := pluralDeleteRecordTypes(ctx, a.query, recordTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete record types: %w", err)
	}
	return nil
}

// PluralDeleteRecordTypesWithSd SD付きで議事録タイプを複数削除する。
func (a *PgAdapter) PluralDeleteRecordTypesWithSd(
	ctx context.Context, sd store.Sd, recordTypeIDs []uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeleteRecordTypes(ctx, qtx, recordTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete record types: %w", err)
	}
	return nil
}

func findRecordTypeByID(
	ctx context.Context, qtx *query.Queries, recordTypeID uuid.UUID,
) (entity.RecordType, error) {
	e, err := qtx.FindRecordTypeByID(ctx, recordTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, store.ErrDataNoRecord
		}
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	entity := entity.RecordType{
		RecordTypeID: e.RecordTypeID,
		Name:         e.Name,
		Key:          e.Key,
	}
	return entity, nil
}

// FindRecordTypeByID 議事録タイプを取得する。
func (a *PgAdapter) FindRecordTypeByID(
	ctx context.Context, recordTypeID uuid.UUID,
) (entity.RecordType, error) {
	e, err := findRecordTypeByID(ctx, a.query, recordTypeID)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	return e, nil
}

// FindRecordTypeByIDWithSd SD付きで議事録タイプを取得する。
func (a *PgAdapter) FindRecordTypeByIDWithSd(
	ctx context.Context, sd store.Sd, recordTypeID uuid.UUID,
) (entity.RecordType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RecordType{}, store.ErrNotFoundDescriptor
	}
	e, err := findRecordTypeByID(ctx, qtx, recordTypeID)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	return e, nil
}

func findRecordTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.RecordType, error) {
	e, err := qtx.FindRecordTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, store.ErrDataNoRecord
		}
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	entity := entity.RecordType{
		RecordTypeID: e.RecordTypeID,
		Name:         e.Name,
		Key:          e.Key,
	}
	return entity, nil
}

// FindRecordTypeByKey 議事録タイプを取得する。
func (a *PgAdapter) FindRecordTypeByKey(ctx context.Context, key string) (entity.RecordType, error) {
	e, err := findRecordTypeByKey(ctx, a.query, key)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	return e, nil
}

// FindRecordTypeByKeyWithSd SD付きで議事録タイプを取得する。
func (a *PgAdapter) FindRecordTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.RecordType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RecordType{}, store.ErrNotFoundDescriptor
	}
	e, err := findRecordTypeByKey(ctx, qtx, key)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to find record type: %w", err)
	}
	return e, nil
}

// RecordTypeCursor is a cursor for RecordType.
type RecordTypeCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getRecordTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereRecordTypeParam,
	order parameter.RecordTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.RecordType], error) {
	eConvFunc := func(e query.RecordType) (entity.RecordType, error) {
		return entity.RecordType{
			RecordTypeID: e.RecordTypeID,
			Name:         e.Name,
			Key:          e.Key,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountRecordTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountRecordTypes(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count record types: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.RecordType, error) {
		p := query.GetRecordTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetRecordTypes(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.RecordType{}, nil
			}
			return nil, fmt.Errorf("failed to get record types: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.RecordType, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.RecordTypeNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetRecordTypesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetRecordTypesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get record types: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.RecordType, error) {
		p := query.GetRecordTypesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetRecordTypesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get record types: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.RecordType) (entity.Int, any) {
		switch subCursor {
		case parameter.RecordTypeDefaultCursorKey:
			return entity.Int(e.MRecordTypesPkey), nil
		case parameter.RecordTypeNameCursorKey:
			return entity.Int(e.MRecordTypesPkey), e.Name
		}
		return entity.Int(e.MRecordTypesPkey), nil
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
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get record types: %w", err)
	}
	return res, nil
}

// GetRecordTypes 議事録タイプを取得する。
func (a *PgAdapter) GetRecordTypes(
	ctx context.Context,
	where parameter.WhereRecordTypeParam,
	order parameter.RecordTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.RecordType], error) {
	r, err := getRecordTypes(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get record types: %w", err)
	}
	return r, nil
}

// GetRecordTypesWithSd SD付きで議事録タイプを取得する。
func (a *PgAdapter) GetRecordTypesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereRecordTypeParam,
	order parameter.RecordTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.RecordType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.RecordType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getRecordTypes(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get record types: %w", err)
	}
	return r, nil
}

func getPluralRecordTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	p := query.GetPluralRecordTypesParams{
		RecordTypeIds: ids,
		Offset:        int32(np.Offset.Int64),
		Limit:         int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralRecordTypes(ctx, p)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get plural record types: %w", err)
	}
	entities := make([]entity.RecordType, len(e))
	for i, v := range e {
		entities[i] = entity.RecordType{
			RecordTypeID: v.RecordTypeID,
			Name:         v.Name,
			Key:          v.Key,
		}
	}
	return store.ListResult[entity.RecordType]{Data: entities}, nil
}

// GetPluralRecordTypes 議事録タイプを取得する。
func (a *PgAdapter) GetPluralRecordTypes(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	r, err := getPluralRecordTypes(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get plural record types: %w", err)
	}
	return r, nil
}

// GetPluralRecordTypesWithSd SD付きで議事録タイプを取得する。
func (a *PgAdapter) GetPluralRecordTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.RecordType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralRecordTypes(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.RecordType]{}, fmt.Errorf("failed to get plural record types: %w", err)
	}
	return r, nil
}

func updateRecordType(
	ctx context.Context, qtx *query.Queries, recordTypeID uuid.UUID, param parameter.UpdateRecordTypeParams,
) (entity.RecordType, error) {
	p := query.UpdateRecordTypeParams{
		RecordTypeID: recordTypeID,
		Name:         param.Name,
		Key:          param.Key,
	}
	e, err := qtx.UpdateRecordType(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, store.ErrDataNoRecord
		}
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	entity := entity.RecordType{
		RecordTypeID: e.RecordTypeID,
		Name:         e.Name,
		Key:          e.Key,
	}
	return entity, nil
}

// UpdateRecordType 議事録タイプを更新する。
func (a *PgAdapter) UpdateRecordType(
	ctx context.Context, recordTypeID uuid.UUID, param parameter.UpdateRecordTypeParams,
) (entity.RecordType, error) {
	e, err := updateRecordType(ctx, a.query, recordTypeID, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	return e, nil
}

// UpdateRecordTypeWithSd SD付きで議事録タイプを更新する。
func (a *PgAdapter) UpdateRecordTypeWithSd(
	ctx context.Context, sd store.Sd, recordTypeID uuid.UUID, param parameter.UpdateRecordTypeParams,
) (entity.RecordType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RecordType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateRecordType(ctx, qtx, recordTypeID, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	return e, nil
}

func updateRecordTypeByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateRecordTypeByKeyParams,
) (entity.RecordType, error) {
	p := query.UpdateRecordTypeByKeyParams{
		Key:  key,
		Name: param.Name,
	}
	e, err := qtx.UpdateRecordTypeByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, store.ErrDataNoRecord
		}
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	entity := entity.RecordType{
		RecordTypeID: e.RecordTypeID,
		Name:         e.Name,
		Key:          e.Key,
	}
	return entity, nil
}

// UpdateRecordTypeByKey 議事録タイプを更新する。
func (a *PgAdapter) UpdateRecordTypeByKey(
	ctx context.Context, key string, param parameter.UpdateRecordTypeByKeyParams,
) (entity.RecordType, error) {
	e, err := updateRecordTypeByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	return e, nil
}

// UpdateRecordTypeByKeyWithSd SD付きで議事録タイプを更新する。
func (a *PgAdapter) UpdateRecordTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateRecordTypeByKeyParams,
) (entity.RecordType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RecordType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateRecordTypeByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.RecordType{}, fmt.Errorf("failed to update record type: %w", err)
	}
	return e, nil
}

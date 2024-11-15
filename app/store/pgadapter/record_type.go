package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
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
	return countRecordTypes(ctx, a.query, where)
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
	return countRecordTypes(ctx, qtx, where)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.RecordType{}, errhandle.NewModelDuplicatedError("record type")
		}
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
	return createRecordType(ctx, a.query, param)
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
	return createRecordType(ctx, qtx, param)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("record type")
		}
		return 0, fmt.Errorf("failed to create record types: %w", err)
	}
	return c, nil
}

// CreateRecordTypes 議事録タイプを作成する。
func (a *PgAdapter) CreateRecordTypes(
	ctx context.Context, params []parameter.CreateRecordTypeParam,
) (int64, error) {
	return createRecordTypes(ctx, a.query, params)
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
	return createRecordTypes(ctx, qtx, params)
}

func deleteRecordType(ctx context.Context, qtx *query.Queries, recordTypeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteRecordType(ctx, recordTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete record type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("record type")
	}
	return c, nil
}

// DeleteRecordType 議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordType(ctx context.Context, recordTypeID uuid.UUID) (int64, error) {
	return deleteRecordType(ctx, a.query, recordTypeID)
}

// DeleteRecordTypeWithSd SD付きで議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeWithSd(
	ctx context.Context, sd store.Sd, recordTypeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteRecordType(ctx, qtx, recordTypeID)
}

func deleteRecordTypeByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteRecordTypeByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete record type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("record type")
	}
	return c, nil
}

// DeleteRecordTypeByKey 議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeByKey(ctx context.Context, key string) (int64, error) {
	return deleteRecordTypeByKey(ctx, a.query, key)
}

// DeleteRecordTypeByKeyWithSd SD付きで議事録タイプを削除する。
func (a *PgAdapter) DeleteRecordTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteRecordTypeByKey(ctx, qtx, key)
}

func pluralDeleteRecordTypes(ctx context.Context, qtx *query.Queries, recordTypeIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteRecordTypes(ctx, recordTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete record types: %w", err)
	}
	if c != int64(len(recordTypeIDs)) {
		return 0, errhandle.NewModelNotFoundError("record type")
	}
	return c, nil
}

// PluralDeleteRecordTypes 議事録タイプを複数削除する。
func (a *PgAdapter) PluralDeleteRecordTypes(ctx context.Context, recordTypeIDs []uuid.UUID) (int64, error) {
	return pluralDeleteRecordTypes(ctx, a.query, recordTypeIDs)
}

// PluralDeleteRecordTypesWithSd SD付きで議事録タイプを複数削除する。
func (a *PgAdapter) PluralDeleteRecordTypesWithSd(
	ctx context.Context, sd store.Sd, recordTypeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteRecordTypes(ctx, qtx, recordTypeIDs)
}

func findRecordTypeByID(
	ctx context.Context, qtx *query.Queries, recordTypeID uuid.UUID,
) (entity.RecordType, error) {
	e, err := qtx.FindRecordTypeByID(ctx, recordTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, errhandle.NewModelNotFoundError("record type")
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
	return findRecordTypeByID(ctx, a.query, recordTypeID)
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
	return findRecordTypeByID(ctx, qtx, recordTypeID)
}

func findRecordTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.RecordType, error) {
	e, err := qtx.FindRecordTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RecordType{}, errhandle.NewModelNotFoundError("record type")
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
	return findRecordTypeByKey(ctx, a.query, key)
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
	return findRecordTypeByKey(ctx, qtx, key)
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
	return getRecordTypes(ctx, a.query, where, order, np, cp, wc)
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
	return getRecordTypes(ctx, qtx, where, order, np, cp, wc)
}

func getPluralRecordTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID,
	order parameter.RecordTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	var e []query.RecordType
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralRecordTypes(ctx, query.GetPluralRecordTypesParams{
			RecordTypeIds: ids,
			OrderMethod:   order.GetStringValue(),
		})
	} else {
		p := query.GetPluralRecordTypesUseNumberedPaginateParams{
			RecordTypeIds: ids,
			Offset:        int32(np.Offset.Int64),
			Limit:         int32(np.Limit.Int64),
		}
		e, err = qtx.GetPluralRecordTypesUseNumberedPaginate(ctx, p)
	}
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
	ctx context.Context, ids []uuid.UUID, order parameter.RecordTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	return getPluralRecordTypes(ctx, a.query, ids, order, np)
}

// GetPluralRecordTypesWithSd SD付きで議事録タイプを取得する。
func (a *PgAdapter) GetPluralRecordTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID,
	order parameter.RecordTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.RecordType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.RecordType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralRecordTypes(ctx, qtx, ids, order, np)
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
			return entity.RecordType{}, errhandle.NewModelNotFoundError("record type")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.RecordType{}, errhandle.NewModelDuplicatedError("record type")
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
	return updateRecordType(ctx, a.query, recordTypeID, param)
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
	return updateRecordType(ctx, qtx, recordTypeID, param)
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
			return entity.RecordType{}, errhandle.NewModelNotFoundError("record type")
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
	return updateRecordTypeByKey(ctx, a.query, key, param)
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
	return updateRecordTypeByKey(ctx, qtx, key, param)
}

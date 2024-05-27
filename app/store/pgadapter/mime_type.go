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

func countMimeTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMimeTypeParam,
) (int64, error) {
	p := query.CountMimeTypesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountMimeTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count mime type: %w", err)
	}
	return c, nil
}

// CountMimeTypes マイムタイプ数を取得する。
func (a *PgAdapter) CountMimeTypes(ctx context.Context, where parameter.WhereMimeTypeParam) (int64, error) {
	return countMimeTypes(ctx, a.query, where)
}

// CountMimeTypesWithSd SD付きでマイムタイプ数を取得する。
func (a *PgAdapter) CountMimeTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMimeTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMimeTypes(ctx, qtx, where)
}

func createMimeType(
	ctx context.Context, qtx *query.Queries, param parameter.CreateMimeTypeParam,
) (entity.MimeType, error) {
	p := query.CreateMimeTypeParams{
		Name: param.Name,
		Key:  param.Key,
		Kind: param.Kind,
	}
	e, err := qtx.CreateMimeType(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.MimeType{}, errhandle.NewModelDuplicatedError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to create mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// CreateMimeType マイムタイプを作成する。
func (a *PgAdapter) CreateMimeType(
	ctx context.Context, param parameter.CreateMimeTypeParam,
) (entity.MimeType, error) {
	return createMimeType(ctx, a.query, param)
}

// CreateMimeTypeWithSd SD付きでマイムタイプを作成する。
func (a *PgAdapter) CreateMimeTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateMimeTypeParam,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return createMimeType(ctx, qtx, param)
}

func createMimeTypes(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateMimeTypeParam,
) (int64, error) {
	p := make([]query.CreateMimeTypesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateMimeTypesParams{
			Name: param.Name,
			Key:  param.Key,
			Kind: param.Kind,
		}
	}
	c, err := qtx.CreateMimeTypes(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("mime type")
		}
		return 0, fmt.Errorf("failed to create mime types: %w", err)
	}
	return c, nil
}

// CreateMimeTypes マイムタイプを作成する。
func (a *PgAdapter) CreateMimeTypes(
	ctx context.Context, params []parameter.CreateMimeTypeParam,
) (int64, error) {
	return createMimeTypes(ctx, a.query, params)
}

// CreateMimeTypesWithSd SD付きでマイムタイプを作成する。
func (a *PgAdapter) CreateMimeTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateMimeTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createMimeTypes(ctx, qtx, params)
}

func deleteMimeType(ctx context.Context, qtx *query.Queries, mimeTypeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteMimeType(ctx, mimeTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete mime type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("mime type")
	}
	return c, nil
}

// DeleteMimeType マイムタイプを削除する。
func (a *PgAdapter) DeleteMimeType(ctx context.Context, mimeTypeID uuid.UUID) (int64, error) {
	return deleteMimeType(ctx, a.query, mimeTypeID)
}

// DeleteMimeTypeWithSd SD付きでマイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeWithSd(
	ctx context.Context, sd store.Sd, mimeTypeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteMimeType(ctx, qtx, mimeTypeID)
}

func deleteMimeTypeByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteMimeTypeByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete mime type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("mime type")
	}
	return c, nil
}

// DeleteMimeTypeByKey マイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeByKey(ctx context.Context, key string) (int64, error) {
	return deleteMimeTypeByKey(ctx, a.query, key)
}

// DeleteMimeTypeByKeyWithSd SD付きでマイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteMimeTypeByKey(ctx, qtx, key)
}

func pluralDeleteMimeTypes(ctx context.Context, qtx *query.Queries, mimeTypeIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteMimeTypes(ctx, mimeTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete mime types: %w", err)
	}
	if c != int64(len(mimeTypeIDs)) {
		return 0, errhandle.NewModelNotFoundError("mime type")
	}
	return c, nil
}

// PluralDeleteMimeTypes マイムタイプを複数削除する。
func (a *PgAdapter) PluralDeleteMimeTypes(ctx context.Context, mimeTypeIDs []uuid.UUID) (int64, error) {
	return pluralDeleteMimeTypes(ctx, a.query, mimeTypeIDs)
}

// PluralDeleteMimeTypesWithSd SD付きでマイムタイプを複数削除する。
func (a *PgAdapter) PluralDeleteMimeTypesWithSd(
	ctx context.Context, sd store.Sd, mimeTypeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteMimeTypes(ctx, qtx, mimeTypeIDs)
}

func findMimeTypeByID(
	ctx context.Context, qtx *query.Queries, mimeTypeID uuid.UUID,
) (entity.MimeType, error) {
	e, err := qtx.FindMimeTypeByID(ctx, mimeTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, errhandle.NewModelNotFoundError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// FindMimeTypeByID マイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByID(
	ctx context.Context, mimeTypeID uuid.UUID,
) (entity.MimeType, error) {
	return findMimeTypeByID(ctx, a.query, mimeTypeID)
}

// FindMimeTypeByIDWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByIDWithSd(
	ctx context.Context, sd store.Sd, mimeTypeID uuid.UUID,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return findMimeTypeByID(ctx, qtx, mimeTypeID)
}

func findMimeTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.MimeType, error) {
	e, err := qtx.FindMimeTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, errhandle.NewModelNotFoundError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// FindMimeTypeByKey マイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByKey(ctx context.Context, key string) (entity.MimeType, error) {
	return findMimeTypeByKey(ctx, a.query, key)
}

// FindMimeTypeByKeyWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return findMimeTypeByKey(ctx, qtx, key)
}

func findMimeTypeByKind(ctx context.Context, qtx *query.Queries, kind string) (entity.MimeType, error) {
	e, err := qtx.FindMimeTypeByKind(ctx, kind)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, errhandle.NewModelNotFoundError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// FindMimeTypeByKind マイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByKind(ctx context.Context, kind string) (entity.MimeType, error) {
	return findMimeTypeByKind(ctx, a.query, kind)
}

// FindMimeTypeByKindWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) FindMimeTypeByKindWithSd(
	ctx context.Context, sd store.Sd, kind string,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return findMimeTypeByKind(ctx, qtx, kind)
}

// MimeTypeCursor is a cursor for MimeType.
type MimeTypeCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getMimeTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMimeTypeParam,
	order parameter.MimeTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MimeType], error) {
	eConvFunc := func(e query.MimeType) (entity.MimeType, error) {
		return entity.MimeType{
			MimeTypeID: e.MimeTypeID,
			Name:       e.Name,
			Key:        e.Key,
			Kind:       e.Kind,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountMimeTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountMimeTypes(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count mime types: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.MimeType, error) {
		p := query.GetMimeTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetMimeTypes(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.MimeType{}, nil
			}
			return nil, fmt.Errorf("failed to get mime types: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.MimeType, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.MimeTypeNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetMimeTypesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetMimeTypesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get mime types: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.MimeType, error) {
		p := query.GetMimeTypesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetMimeTypesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get mime types: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.MimeType) (entity.Int, any) {
		switch subCursor {
		case parameter.MimeTypeDefaultCursorKey:
			return entity.Int(e.MMimeTypesPkey), nil
		case parameter.MimeTypeNameCursorKey:
			return entity.Int(e.MMimeTypesPkey), e.Name
		}
		return entity.Int(e.MMimeTypesPkey), nil
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
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get mime types: %w", err)
	}
	return res, nil
}

// GetMimeTypes マイムタイプを取得する。
func (a *PgAdapter) GetMimeTypes(
	ctx context.Context,
	where parameter.WhereMimeTypeParam,
	order parameter.MimeTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MimeType], error) {
	return getMimeTypes(ctx, a.query, where, order, np, cp, wc)
}

// GetMimeTypesWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) GetMimeTypesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereMimeTypeParam,
	order parameter.MimeTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MimeType]{}, store.ErrNotFoundDescriptor
	}
	return getMimeTypes(ctx, qtx, where, order, np, cp, wc)
}

func getPluralMimeTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID,
	order parameter.MimeTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	var e []query.MimeType
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMimeTypes(ctx, query.GetPluralMimeTypesParams{
			MimeTypeIds: ids,
			OrderMethod: order.GetStringValue(),
		})
	} else {
		p := query.GetPluralMimeTypesUseNumberedPaginateParams{
			MimeTypeIds: ids,
			Offset:      int32(np.Offset.Int64),
			Limit:       int32(np.Limit.Int64),
		}
		e, err = qtx.GetPluralMimeTypesUseNumberedPaginate(ctx, p)
	}
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get plural mime types: %w", err)
	}
	entities := make([]entity.MimeType, len(e))
	for i, v := range e {
		entities[i] = entity.MimeType{
			MimeTypeID: v.MimeTypeID,
			Name:       v.Name,
			Key:        v.Key,
			Kind:       v.Kind,
		}
	}
	return store.ListResult[entity.MimeType]{Data: entities}, nil
}

// GetPluralMimeTypes マイムタイプを取得する。
func (a *PgAdapter) GetPluralMimeTypes(
	ctx context.Context, ids []uuid.UUID, order parameter.MimeTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	return getPluralMimeTypes(ctx, a.query, ids, order, np)
}

// GetPluralMimeTypesWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) GetPluralMimeTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID,
	order parameter.MimeTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MimeType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMimeTypes(ctx, qtx, ids, order, np)
}

func updateMimeType(
	ctx context.Context, qtx *query.Queries, mimeTypeID uuid.UUID, param parameter.UpdateMimeTypeParams,
) (entity.MimeType, error) {
	p := query.UpdateMimeTypeParams{
		MimeTypeID: mimeTypeID,
		Name:       param.Name,
		Key:        param.Key,
		Kind:       param.Kind,
	}
	e, err := qtx.UpdateMimeType(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, errhandle.NewModelNotFoundError("mime type")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.MimeType{}, errhandle.NewModelDuplicatedError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// UpdateMimeType マイムタイプを更新する。
func (a *PgAdapter) UpdateMimeType(
	ctx context.Context, mimeTypeID uuid.UUID, param parameter.UpdateMimeTypeParams,
) (entity.MimeType, error) {
	return updateMimeType(ctx, a.query, mimeTypeID, param)
}

// UpdateMimeTypeWithSd SD付きでマイムタイプを更新する。
func (a *PgAdapter) UpdateMimeTypeWithSd(
	ctx context.Context, sd store.Sd, mimeTypeID uuid.UUID, param parameter.UpdateMimeTypeParams,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return updateMimeType(ctx, qtx, mimeTypeID, param)
}

func updateMimeTypeByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateMimeTypeByKeyParams,
) (entity.MimeType, error) {
	p := query.UpdateMimeTypeByKeyParams{
		Key:  key,
		Name: param.Name,
		Kind: param.Kind,
	}
	e, err := qtx.UpdateMimeTypeByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, errhandle.NewModelNotFoundError("mime type")
		}
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	entity := entity.MimeType{
		MimeTypeID: e.MimeTypeID,
		Name:       e.Name,
		Key:        e.Key,
		Kind:       e.Kind,
	}
	return entity, nil
}

// UpdateMimeTypeByKey マイムタイプを更新する。
func (a *PgAdapter) UpdateMimeTypeByKey(
	ctx context.Context, key string, param parameter.UpdateMimeTypeByKeyParams,
) (entity.MimeType, error) {
	return updateMimeTypeByKey(ctx, a.query, key, param)
}

// UpdateMimeTypeByKeyWithSd SD付きでマイムタイプを更新する。
func (a *PgAdapter) UpdateMimeTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateMimeTypeByKeyParams,
) (entity.MimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MimeType{}, store.ErrNotFoundDescriptor
	}
	return updateMimeTypeByKey(ctx, qtx, key, param)
}

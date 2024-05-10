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
	c, err := countMimeTypes(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count mime type: %w", err)
	}
	return c, nil
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
	c, err := countMimeTypes(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count mime type: %w", err)
	}
	return c, nil
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
	e, err := createMimeType(ctx, a.query, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to create mime type: %w", err)
	}
	return e, nil
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
	e, err := createMimeType(ctx, qtx, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to create mime type: %w", err)
	}
	return e, nil
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
		return 0, fmt.Errorf("failed to create mime types: %w", err)
	}
	return c, nil
}

// CreateMimeTypes マイムタイプを作成する。
func (a *PgAdapter) CreateMimeTypes(
	ctx context.Context, params []parameter.CreateMimeTypeParam,
) (int64, error) {
	c, err := createMimeTypes(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create mime types: %w", err)
	}
	return c, nil
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
	c, err := createMimeTypes(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create mime types: %w", err)
	}
	return c, nil
}

func deleteMimeType(ctx context.Context, qtx *query.Queries, mimeTypeID uuid.UUID) error {
	err := qtx.DeleteMimeType(ctx, mimeTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

// DeleteMimeType マイムタイプを削除する。
func (a *PgAdapter) DeleteMimeType(ctx context.Context, mimeTypeID uuid.UUID) error {
	err := deleteMimeType(ctx, a.query, mimeTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

// DeleteMimeTypeWithSd SD付きでマイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeWithSd(
	ctx context.Context, sd store.Sd, mimeTypeID uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteMimeType(ctx, qtx, mimeTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

func deleteMimeTypeByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeleteMimeTypeByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

// DeleteMimeTypeByKey マイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeByKey(ctx context.Context, key string) error {
	err := deleteMimeTypeByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

// DeleteMimeTypeByKeyWithSd SD付きでマイムタイプを削除する。
func (a *PgAdapter) DeleteMimeTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteMimeTypeByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete mime type: %w", err)
	}
	return nil
}

func pluralDeleteMimeTypes(ctx context.Context, qtx *query.Queries, mimeTypeIDs []uuid.UUID) error {
	err := qtx.PluralDeleteMimeTypes(ctx, mimeTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete mime types: %w", err)
	}
	return nil
}

// PluralDeleteMimeTypes マイムタイプを複数削除する。
func (a *PgAdapter) PluralDeleteMimeTypes(ctx context.Context, mimeTypeIDs []uuid.UUID) error {
	err := pluralDeleteMimeTypes(ctx, a.query, mimeTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete mime types: %w", err)
	}
	return nil
}

// PluralDeleteMimeTypesWithSd SD付きでマイムタイプを複数削除する。
func (a *PgAdapter) PluralDeleteMimeTypesWithSd(
	ctx context.Context, sd store.Sd, mimeTypeIDs []uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeleteMimeTypes(ctx, qtx, mimeTypeIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete mime types: %w", err)
	}
	return nil
}

func findMimeTypeByID(
	ctx context.Context, qtx *query.Queries, mimeTypeID uuid.UUID,
) (entity.MimeType, error) {
	e, err := qtx.FindMimeTypeByID(ctx, mimeTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, store.ErrDataNoRecord
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
	e, err := findMimeTypeByID(ctx, a.query, mimeTypeID)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	return e, nil
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
	e, err := findMimeTypeByID(ctx, qtx, mimeTypeID)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	return e, nil
}

func findMimeTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.MimeType, error) {
	e, err := qtx.FindMimeTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MimeType{}, store.ErrDataNoRecord
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
	e, err := findMimeTypeByKey(ctx, a.query, key)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	return e, nil
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
	e, err := findMimeTypeByKey(ctx, qtx, key)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type: %w", err)
	}
	return e, nil
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
	r, err := getMimeTypes(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get mime types: %w", err)
	}
	return r, nil
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
	r, err := getMimeTypes(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get mime types: %w", err)
	}
	return r, nil
}

func getPluralMimeTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	p := query.GetPluralMimeTypesParams{
		MimeTypeIds: ids,
		Offset:      int32(np.Offset.Int64),
		Limit:       int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralMimeTypes(ctx, p)
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
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	r, err := getPluralMimeTypes(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get plural mime types: %w", err)
	}
	return r, nil
}

// GetPluralMimeTypesWithSd SD付きでマイムタイプを取得する。
func (a *PgAdapter) GetPluralMimeTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.MimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MimeType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralMimeTypes(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get plural mime types: %w", err)
	}
	return r, nil
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
			return entity.MimeType{}, store.ErrDataNoRecord
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
	e, err := updateMimeType(ctx, a.query, mimeTypeID, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	return e, nil
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
	e, err := updateMimeType(ctx, qtx, mimeTypeID, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	return e, nil
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
			return entity.MimeType{}, store.ErrDataNoRecord
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
	e, err := updateMimeTypeByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	return e, nil
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
	e, err := updateMimeTypeByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	return e, nil
}

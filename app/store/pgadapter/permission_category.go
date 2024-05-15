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

func countPermissionCategories(
	ctx context.Context, qtx *query.Queries, where parameter.WherePermissionCategoryParam,
) (int64, error) {
	p := query.CountPermissionCategoriesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountPermissionCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission category: %w", err)
	}
	return c, nil
}

// CountPermissionCategories 権限カテゴリー数を取得する。
func (a *PgAdapter) CountPermissionCategories(
	ctx context.Context, where parameter.WherePermissionCategoryParam,
) (int64, error) {
	return countPermissionCategories(ctx, a.query, where)
}

// CountPermissionCategoriesWithSd SD付きで権限カテゴリー数を取得する。
func (a *PgAdapter) CountPermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WherePermissionCategoryParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countPermissionCategories(ctx, qtx, where)
}

func createPermissionCategory(
	ctx context.Context, qtx *query.Queries, param parameter.CreatePermissionCategoryParam,
) (entity.PermissionCategory, error) {
	p := query.CreatePermissionCategoryParams{
		Name:        param.Name,
		Key:         param.Key,
		Description: param.Description,
	}
	e, err := qtx.CreatePermissionCategory(ctx, p)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to create permission category: %w", err)
	}
	entity := entity.PermissionCategory{
		PermissionCategoryID: e.PermissionCategoryID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
	}
	return entity, nil
}

// CreatePermissionCategory 権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionCategory(
	ctx context.Context, param parameter.CreatePermissionCategoryParam,
) (entity.PermissionCategory, error) {
	return createPermissionCategory(ctx, a.query, param)
}

// CreatePermissionCategoryWithSd SD付きで権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreatePermissionCategoryParam,
) (entity.PermissionCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	return createPermissionCategory(ctx, qtx, param)
}

func createPermissionCategories(
	ctx context.Context, qtx *query.Queries, params []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	p := make([]query.CreatePermissionCategoriesParams, len(params))
	for i, param := range params {
		p[i] = query.CreatePermissionCategoriesParams{
			Name:        param.Name,
			Key:         param.Key,
			Description: param.Description,
		}
	}
	c, err := qtx.CreatePermissionCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission categories: %w", err)
	}
	return c, nil
}

// CreatePermissionCategories 権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionCategories(
	ctx context.Context, params []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	return createPermissionCategories(ctx, a.query, params)
}

// CreatePermissionCategoriesWithSd SD付きで権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createPermissionCategories(ctx, qtx, params)
}

func deletePermissionCategory(ctx context.Context, qtx *query.Queries, permissionCategoryID uuid.UUID) (int64, error) {
	c, err := qtx.DeletePermissionCategory(ctx, permissionCategoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete permission category: %w", err)
	}
	return c, nil
}

// DeletePermissionCategory 権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionCategory(ctx context.Context, permissionCategoryID uuid.UUID) (int64, error) {
	return deletePermissionCategory(ctx, a.query, permissionCategoryID)
}

// DeletePermissionCategoryWithSd SD付きで権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePermissionCategory(ctx, qtx, permissionCategoryID)
}

func deletePermissionCategoryByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeletePermissionCategoryByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete permission category: %w", err)
	}
	return c, nil
}

// DeletePermissionCategoryByKey 権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionCategoryByKey(ctx context.Context, key string) (int64, error) {
	return deletePermissionCategoryByKey(ctx, a.query, key)
}

// DeletePermissionCategoryByKeyWithSd SD付きで権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePermissionCategoryByKey(ctx, qtx, key)
}

func pluralDeletePermissionCategories(
	ctx context.Context, qtx *query.Queries, permissionCategoryIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeletePermissionCategories(ctx, permissionCategoryIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete permission categories: %w", err)
	}
	return c, nil
}

// PluralDeletePermissionCategories 権限カテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePermissionCategories(
	ctx context.Context, permissionCategoryIDs []uuid.UUID,
) (int64, error) {
	return pluralDeletePermissionCategories(ctx, a.query, permissionCategoryIDs)
}

// PluralDeletePermissionCategoriesWithSd SD付きで権限カテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeletePermissionCategories(ctx, qtx, permissionCategoryIDs)
}

func findPermissionCategoryByID(
	ctx context.Context, qtx *query.Queries, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	e, err := qtx.FindPermissionCategoryByID(ctx, permissionCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, errhandle.NewModelNotFoundError("permission category")
		}
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	entity := entity.PermissionCategory{
		PermissionCategoryID: e.PermissionCategoryID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
	}
	return entity, nil
}

// FindPermissionCategoryByID 権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionCategoryByID(
	ctx context.Context, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	return findPermissionCategoryByID(ctx, a.query, permissionCategoryID)
}

// FindPermissionCategoryByIDWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionCategoryByIDWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	return findPermissionCategoryByID(ctx, qtx, permissionCategoryID)
}

func findPermissionCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PermissionCategory, error) {
	e, err := qtx.FindPermissionCategoryByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, errhandle.NewModelNotFoundError("permission category")
		}
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	entity := entity.PermissionCategory{
		PermissionCategoryID: e.PermissionCategoryID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
	}
	return entity, nil
}

// FindPermissionCategoryByKey 権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionCategoryByKey(ctx context.Context, key string) (entity.PermissionCategory, error) {
	return findPermissionCategoryByKey(ctx, a.query, key)
}

// FindPermissionCategoryByKeyWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.PermissionCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	return findPermissionCategoryByKey(ctx, qtx, key)
}

// PermissionCategoryCursor is a cursor for PermissionCategory.
type PermissionCategoryCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getPermissionCategories(
	ctx context.Context, qtx *query.Queries, where parameter.WherePermissionCategoryParam,
	order parameter.PermissionCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PermissionCategory], error) {
	eConvFunc := func(e query.PermissionCategory) (entity.PermissionCategory, error) {
		return entity.PermissionCategory{
			PermissionCategoryID: e.PermissionCategoryID,
			Name:                 e.Name,
			Key:                  e.Key,
			Description:          e.Description,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPermissionCategoriesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountPermissionCategories(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count permission categories: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.PermissionCategory, error) {
		p := query.GetPermissionCategoriesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetPermissionCategories(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.PermissionCategory{}, nil
			}
			return nil, fmt.Errorf("failed to get permission categories: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.PermissionCategory, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PermissionCategoryNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPermissionCategoriesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPermissionCategoriesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get permission categories: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.PermissionCategory, error) {
		p := query.GetPermissionCategoriesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetPermissionCategoriesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get permission categories: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.PermissionCategory) (entity.Int, any) {
		switch subCursor {
		case parameter.PermissionCategoryDefaultCursorKey:
			return entity.Int(e.MPermissionCategoriesPkey), nil
		case parameter.PermissionCategoryNameCursorKey:
			return entity.Int(e.MPermissionCategoriesPkey), e.Name
		}
		return entity.Int(e.MPermissionCategoriesPkey), nil
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
		return store.ListResult[entity.PermissionCategory]{}, fmt.Errorf("failed to get permission categories: %w", err)
	}
	return res, nil
}

// GetPermissionCategories 権限カテゴリーを取得する。
func (a *PgAdapter) GetPermissionCategories(
	ctx context.Context,
	where parameter.WherePermissionCategoryParam,
	order parameter.PermissionCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PermissionCategory], error) {
	return getPermissionCategories(ctx, a.query, where, order, np, cp, wc)
}

// GetPermissionCategoriesWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) GetPermissionCategoriesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePermissionCategoryParam,
	order parameter.PermissionCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PermissionCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PermissionCategory]{}, store.ErrNotFoundDescriptor
	}
	return getPermissionCategories(ctx, qtx, where, order, np, cp, wc)
}

func getPluralPermissionCategories(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PermissionCategory], error) {
	p := query.GetPluralPermissionCategoriesParams{
		PermissionCategoryIds: ids,
		Offset:                int32(np.Offset.Int64),
		Limit:                 int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralPermissionCategories(ctx, p)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{},
			fmt.Errorf("failed to get plural permission categories: %w", err)
	}
	entities := make([]entity.PermissionCategory, len(e))
	for i, v := range e {
		entities[i] = entity.PermissionCategory{
			PermissionCategoryID: v.PermissionCategoryID,
			Name:                 v.Name,
			Key:                  v.Key,
			Description:          v.Description,
		}
	}
	return store.ListResult[entity.PermissionCategory]{Data: entities}, nil
}

// GetPluralPermissionCategories 権限カテゴリーを取得する。
func (a *PgAdapter) GetPluralPermissionCategories(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PermissionCategory], error) {
	return getPluralPermissionCategories(ctx, a.query, ids, np)
}

// GetPluralPermissionCategoriesWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) GetPluralPermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PermissionCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PermissionCategory]{}, store.ErrNotFoundDescriptor
	}
	return getPluralPermissionCategories(ctx, qtx, ids, np)
}

func updatePermissionCategory(
	ctx context.Context, qtx *query.Queries,
	permissionCategoryID uuid.UUID, param parameter.UpdatePermissionCategoryParams,
) (entity.PermissionCategory, error) {
	p := query.UpdatePermissionCategoryParams{
		PermissionCategoryID: permissionCategoryID,
		Name:                 param.Name,
		Key:                  param.Key,
		Description:          param.Description,
	}
	e, err := qtx.UpdatePermissionCategory(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, errhandle.NewModelNotFoundError("permission category")
		}
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	entity := entity.PermissionCategory{
		PermissionCategoryID: e.PermissionCategoryID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
	}
	return entity, nil
}

// UpdatePermissionCategory 権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionCategory(
	ctx context.Context, permissionCategoryID uuid.UUID, param parameter.UpdatePermissionCategoryParams,
) (entity.PermissionCategory, error) {
	return updatePermissionCategory(ctx, a.query, permissionCategoryID, param)
}

// UpdatePermissionCategoryWithSd SD付きで権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID, param parameter.UpdatePermissionCategoryParams,
) (entity.PermissionCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	return updatePermissionCategory(ctx, qtx, permissionCategoryID, param)
}

func updatePermissionCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdatePermissionCategoryByKeyParams,
) (entity.PermissionCategory, error) {
	p := query.UpdatePermissionCategoryByKeyParams{
		Key:         key,
		Name:        param.Name,
		Description: param.Description,
	}
	e, err := qtx.UpdatePermissionCategoryByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, errhandle.NewModelNotFoundError("permission category")
		}
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	entity := entity.PermissionCategory{
		PermissionCategoryID: e.PermissionCategoryID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
	}
	return entity, nil
}

// UpdatePermissionCategoryByKey 権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionCategoryByKey(
	ctx context.Context, key string, param parameter.UpdatePermissionCategoryByKeyParams,
) (entity.PermissionCategory, error) {
	return updatePermissionCategoryByKey(ctx, a.query, key, param)
}

// UpdatePermissionCategoryByKeyWithSd SD付きで権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdatePermissionCategoryByKeyParams,
) (entity.PermissionCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	return updatePermissionCategoryByKey(ctx, qtx, key, param)
}

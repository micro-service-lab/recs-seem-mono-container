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

// CountPermissionCategories イベントタイプ数を取得する。
func (a *PgAdapter) CountPermissionCategories(
	ctx context.Context, where parameter.WherePermissionCategoryParam,
) (int64, error) {
	c, err := countPermissionCategories(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission category: %w", err)
	}
	return c, nil
}

// CountPermissionCategoriesWithSd SD付きでイベントタイプ数を取得する。
func (a *PgAdapter) CountPermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WherePermissionCategoryParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countPermissionCategories(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission category: %w", err)
	}
	return c, nil
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

// CreatePermissionCategory イベントタイプを作成する。
func (a *PgAdapter) CreatePermissionCategory(
	ctx context.Context, param parameter.CreatePermissionCategoryParam,
) (entity.PermissionCategory, error) {
	e, err := createPermissionCategory(ctx, a.query, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to create permission category: %w", err)
	}
	return e, nil
}

// CreatePermissionCategoryWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreatePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreatePermissionCategoryParam,
) (entity.PermissionCategory, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := createPermissionCategory(ctx, qtx, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to create permission category: %w", err)
	}
	return e, nil
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

// CreatePermissionCategories イベントタイプを作成する。
func (a *PgAdapter) CreatePermissionCategories(
	ctx context.Context, params []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	c, err := createPermissionCategories(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission categories: %w", err)
	}
	return c, nil
}

// CreatePermissionCategoriesWithSd SD付きでイベントタイプを作成する。
func (a *PgAdapter) CreatePermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createPermissionCategories(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission categories: %w", err)
	}
	return c, nil
}

func deletePermissionCategory(ctx context.Context, qtx *query.Queries, permissionCategoryID uuid.UUID) error {
	err := qtx.DeletePermissionCategory(ctx, permissionCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.ErrDataNoRecord
		}
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

// DeletePermissionCategory イベントタイプを削除する。
func (a *PgAdapter) DeletePermissionCategory(ctx context.Context, permissionCategoryID uuid.UUID) error {
	err := deletePermissionCategory(ctx, a.query, permissionCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

// DeletePermissionCategoryWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeletePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePermissionCategory(ctx, qtx, permissionCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

func deletePermissionCategoryByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeletePermissionCategoryByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.ErrDataNoRecord
		}
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

// DeletePermissionCategoryByKey イベントタイプを削除する。
func (a *PgAdapter) DeletePermissionCategoryByKey(ctx context.Context, key string) error {
	err := deletePermissionCategoryByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

// DeletePermissionCategoryByKeyWithSd SD付きでイベントタイプを削除する。
func (a *PgAdapter) DeletePermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePermissionCategoryByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

func findPermissionCategoryByID(
	ctx context.Context, qtx *query.Queries, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	e, err := qtx.FindPermissionCategoryByID(ctx, permissionCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, store.ErrDataNoRecord
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

// FindPermissionCategoryByID イベントタイプを取得する。
func (a *PgAdapter) FindPermissionCategoryByID(
	ctx context.Context, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	e, err := findPermissionCategoryByID(ctx, a.query, permissionCategoryID)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	return e, nil
}

// FindPermissionCategoryByIDWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindPermissionCategoryByIDWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID,
) (entity.PermissionCategory, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := findPermissionCategoryByID(ctx, qtx, permissionCategoryID)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	return e, nil
}

func findPermissionCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PermissionCategory, error) {
	e, err := qtx.FindPermissionCategoryByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PermissionCategory{}, store.ErrDataNoRecord
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

// FindPermissionCategoryByKey イベントタイプを取得する。
func (a *PgAdapter) FindPermissionCategoryByKey(ctx context.Context, key string) (entity.PermissionCategory, error) {
	e, err := findPermissionCategoryByKey(ctx, a.query, key)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	return e, nil
}

// FindPermissionCategoryByKeyWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) FindPermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.PermissionCategory, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := findPermissionCategoryByKey(ctx, qtx, key)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category: %w", err)
	}
	return e, nil
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

// GetPermissionCategories イベントタイプを取得する。
func (a *PgAdapter) GetPermissionCategories(
	ctx context.Context,
	where parameter.WherePermissionCategoryParam,
	order parameter.PermissionCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PermissionCategory], error) {
	r, err := getPermissionCategories(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{}, fmt.Errorf("failed to get permission categories: %w", err)
	}
	return r, nil
}

// GetPermissionCategoriesWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) GetPermissionCategoriesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePermissionCategoryParam,
	order parameter.PermissionCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PermissionCategory], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PermissionCategory]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPermissionCategories(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{}, fmt.Errorf("failed to get permission categories: %w", err)
	}
	return r, nil
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

// GetPluralPermissionCategories イベントタイプを取得する。
func (a *PgAdapter) GetPluralPermissionCategories(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PermissionCategory], error) {
	r, err := getPluralPermissionCategories(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{},
			fmt.Errorf("failed to get plural permission categories: %w", err)
	}
	return r, nil
}

// GetPluralPermissionCategoriesWithSd SD付きでイベントタイプを取得する。
func (a *PgAdapter) GetPluralPermissionCategoriesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PermissionCategory], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PermissionCategory]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralPermissionCategories(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{},
			fmt.Errorf("failed to get plural permission categories: %w", err)
	}
	return r, nil
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
			return entity.PermissionCategory{}, store.ErrDataNoRecord
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

// UpdatePermissionCategory イベントタイプを更新する。
func (a *PgAdapter) UpdatePermissionCategory(
	ctx context.Context, permissionCategoryID uuid.UUID, param parameter.UpdatePermissionCategoryParams,
) (entity.PermissionCategory, error) {
	e, err := updatePermissionCategory(ctx, a.query, permissionCategoryID, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	return e, nil
}

// UpdatePermissionCategoryWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdatePermissionCategoryWithSd(
	ctx context.Context, sd store.Sd, permissionCategoryID uuid.UUID, param parameter.UpdatePermissionCategoryParams,
) (entity.PermissionCategory, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePermissionCategory(ctx, qtx, permissionCategoryID, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	return e, nil
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
			return entity.PermissionCategory{}, store.ErrDataNoRecord
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

// UpdatePermissionCategoryByKey イベントタイプを更新する。
func (a *PgAdapter) UpdatePermissionCategoryByKey(
	ctx context.Context, key string, param parameter.UpdatePermissionCategoryByKeyParams,
) (entity.PermissionCategory, error) {
	e, err := updatePermissionCategoryByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	return e, nil
}

// UpdatePermissionCategoryByKeyWithSd SD付きでイベントタイプを更新する。
func (a *PgAdapter) UpdatePermissionCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdatePermissionCategoryByKeyParams,
) (entity.PermissionCategory, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PermissionCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePermissionCategoryByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	return e, nil
}

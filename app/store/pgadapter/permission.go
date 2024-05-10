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

func countPermissions(
	ctx context.Context, qtx *query.Queries, where parameter.WherePermissionParam,
) (int64, error) {
	p := query.CountPermissionsParams{
		WhereLikeName:   where.WhereLikeName,
		SearchName:      where.SearchName,
		WhereInCategory: where.WhereInCategory,
		InCategories:    where.InCategories,
	}
	c, err := qtx.CountPermissions(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission: %w", err)
	}
	return c, nil
}

// CountPermissions 権限カテゴリー数を取得する。
func (a *PgAdapter) CountPermissions(
	ctx context.Context, where parameter.WherePermissionParam,
) (int64, error) {
	c, err := countPermissions(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission: %w", err)
	}
	return c, nil
}

// CountPermissionsWithSd SD付きで権限カテゴリー数を取得する。
func (a *PgAdapter) CountPermissionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WherePermissionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countPermissions(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count permission: %w", err)
	}
	return c, nil
}

func createPermission(
	ctx context.Context, qtx *query.Queries, param parameter.CreatePermissionParam,
) (entity.Permission, error) {
	p := query.CreatePermissionParams{
		Name:                 param.Name,
		Key:                  param.Key,
		Description:          param.Description,
		PermissionCategoryID: param.PermissionCategoryID,
	}
	e, err := qtx.CreatePermission(ctx, p)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to create permission: %w", err)
	}
	entity := entity.Permission{
		PermissionID:         e.PermissionID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
		PermissionCategoryID: e.PermissionCategoryID,
	}
	return entity, nil
}

// CreatePermission 権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermission(
	ctx context.Context, param parameter.CreatePermissionParam,
) (entity.Permission, error) {
	e, err := createPermission(ctx, a.query, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to create permission: %w", err)
	}
	return e, nil
}

// CreatePermissionWithSd SD付きで権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreatePermissionParam,
) (entity.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Permission{}, store.ErrNotFoundDescriptor
	}
	e, err := createPermission(ctx, qtx, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to create permission: %w", err)
	}
	return e, nil
}

func createPermissions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreatePermissionParam,
) (int64, error) {
	p := make([]query.CreatePermissionsParams, len(params))
	for i, param := range params {
		p[i] = query.CreatePermissionsParams{
			Name:                 param.Name,
			Key:                  param.Key,
			Description:          param.Description,
			PermissionCategoryID: param.PermissionCategoryID,
		}
	}
	c, err := qtx.CreatePermissions(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}
	return c, nil
}

// CreatePermissions 権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissions(
	ctx context.Context, params []parameter.CreatePermissionParam,
) (int64, error) {
	c, err := createPermissions(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}
	return c, nil
}

// CreatePermissionsWithSd SD付きで権限カテゴリーを作成する。
func (a *PgAdapter) CreatePermissionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreatePermissionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createPermissions(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}
	return c, nil
}

func deletePermission(ctx context.Context, qtx *query.Queries, permissionID uuid.UUID) error {
	err := qtx.DeletePermission(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// DeletePermission 権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermission(ctx context.Context, permissionID uuid.UUID) error {
	err := deletePermission(ctx, a.query, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// DeletePermissionWithSd SD付きで権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionWithSd(
	ctx context.Context, sd store.Sd, permissionID uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePermission(ctx, qtx, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

func deletePermissionByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeletePermissionByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// DeletePermissionByKey 権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionByKey(ctx context.Context, key string) error {
	err := deletePermissionByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// DeletePermissionByKeyWithSd SD付きで権限カテゴリーを削除する。
func (a *PgAdapter) DeletePermissionByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePermissionByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

func pluralDeletePermissions(
	ctx context.Context, qtx *query.Queries, permissionIDs []uuid.UUID,
) error {
	err := qtx.PluralDeletePermissions(ctx, permissionIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete permission: %w", err)
	}
	return nil
}

// PluralDeletePermissions 権限カテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePermissions(
	ctx context.Context, permissionIDs []uuid.UUID,
) error {
	err := pluralDeletePermissions(ctx, a.query, permissionIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete permission: %w", err)
	}
	return nil
}

// PluralDeletePermissionsWithSd SD付きで権限カテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePermissionsWithSd(
	ctx context.Context, sd store.Sd, permissionIDs []uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeletePermissions(ctx, qtx, permissionIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete permission: %w", err)
	}
	return nil
}

func findPermissionByID(
	ctx context.Context, qtx *query.Queries, permissionID uuid.UUID,
) (entity.Permission, error) {
	e, err := qtx.FindPermissionByID(ctx, permissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Permission{}, store.ErrDataNoRecord
		}
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	entity := entity.Permission{
		PermissionID:         e.PermissionID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
		PermissionCategoryID: e.PermissionCategoryID,
	}
	return entity, nil
}

// FindPermissionByID 権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionByID(
	ctx context.Context, permissionID uuid.UUID,
) (entity.Permission, error) {
	e, err := findPermissionByID(ctx, a.query, permissionID)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	return e, nil
}

// FindPermissionByIDWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionByIDWithSd(
	ctx context.Context, sd store.Sd, permissionID uuid.UUID,
) (entity.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Permission{}, store.ErrNotFoundDescriptor
	}
	e, err := findPermissionByID(ctx, qtx, permissionID)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	return e, nil
}

func findPermissionByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.Permission, error) {
	e, err := qtx.FindPermissionByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Permission{}, store.ErrDataNoRecord
		}
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	entity := entity.Permission{
		PermissionID:         e.PermissionID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
		PermissionCategoryID: e.PermissionCategoryID,
	}
	return entity, nil
}

// FindPermissionByKey 権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionByKey(ctx context.Context, key string) (entity.Permission, error) {
	e, err := findPermissionByKey(ctx, a.query, key)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	return e, nil
}

// FindPermissionByKeyWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) FindPermissionByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Permission{}, store.ErrNotFoundDescriptor
	}
	e, err := findPermissionByKey(ctx, qtx, key)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission: %w", err)
	}
	return e, nil
}

// PermissionCursor is a cursor for Permission.
type PermissionCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getPermissions(
	ctx context.Context, qtx *query.Queries, where parameter.WherePermissionParam,
	order parameter.PermissionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Permission], error) {
	eConvFunc := func(e query.Permission) (entity.Permission, error) {
		return entity.Permission{
			PermissionID:         e.PermissionID,
			Name:                 e.Name,
			Key:                  e.Key,
			Description:          e.Description,
			PermissionCategoryID: e.PermissionCategoryID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPermissionsParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
		}
		r, err := qtx.CountPermissions(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count permission: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Permission, error) {
		p := query.GetPermissionsParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
		}
		r, err := qtx.GetPermissions(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Permission{}, nil
			}
			return nil, fmt.Errorf("failed to get permission: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Permission, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PermissionNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPermissionsUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPermissionsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get permission: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Permission, error) {
		p := query.GetPermissionsUseNumberedPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Offset:          offset,
			Limit:           limit,
		}
		r, err := qtx.GetPermissionsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get permission: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Permission) (entity.Int, any) {
		switch subCursor {
		case parameter.PermissionDefaultCursorKey:
			return entity.Int(e.MPermissionsPkey), nil
		case parameter.PermissionNameCursorKey:
			return entity.Int(e.MPermissionsPkey), e.Name
		}
		return entity.Int(e.MPermissionsPkey), nil
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
		return store.ListResult[entity.Permission]{}, fmt.Errorf("failed to get permission: %w", err)
	}
	return res, nil
}

// GetPermissions 権限カテゴリーを取得する。
func (a *PgAdapter) GetPermissions(
	ctx context.Context,
	where parameter.WherePermissionParam,
	order parameter.PermissionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Permission], error) {
	r, err := getPermissions(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Permission]{}, fmt.Errorf("failed to get permission: %w", err)
	}
	return r, nil
}

// GetPermissionsWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) GetPermissionsWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePermissionParam,
	order parameter.PermissionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Permission], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Permission]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPermissions(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Permission]{}, fmt.Errorf("failed to get permission: %w", err)
	}
	return r, nil
}

func getPluralPermissions(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Permission], error) {
	p := query.GetPluralPermissionsParams{
		PermissionIds: ids,
		Offset:        int32(np.Offset.Int64),
		Limit:         int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralPermissions(ctx, p)
	if err != nil {
		return store.ListResult[entity.Permission]{},
			fmt.Errorf("failed to get plural permission: %w", err)
	}
	entities := make([]entity.Permission, len(e))
	for i, v := range e {
		entities[i] = entity.Permission{
			PermissionID:         v.PermissionID,
			Name:                 v.Name,
			Key:                  v.Key,
			Description:          v.Description,
			PermissionCategoryID: v.PermissionCategoryID,
		}
	}
	return store.ListResult[entity.Permission]{Data: entities}, nil
}

// GetPluralPermissions 権限カテゴリーを取得する。
func (a *PgAdapter) GetPluralPermissions(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Permission], error) {
	r, err := getPluralPermissions(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.Permission]{},
			fmt.Errorf("failed to get plural permission: %w", err)
	}
	return r, nil
}

// GetPluralPermissionsWithSd SD付きで権限カテゴリーを取得する。
func (a *PgAdapter) GetPluralPermissionsWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Permission], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Permission]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralPermissions(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.Permission]{},
			fmt.Errorf("failed to get plural permission: %w", err)
	}
	return r, nil
}

func updatePermission(
	ctx context.Context, qtx *query.Queries,
	permissionID uuid.UUID, param parameter.UpdatePermissionParams,
) (entity.Permission, error) {
	p := query.UpdatePermissionParams{
		PermissionID:         permissionID,
		Name:                 param.Name,
		Key:                  param.Key,
		Description:          param.Description,
		PermissionCategoryID: param.PermissionCategoryID,
	}
	e, err := qtx.UpdatePermission(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Permission{}, store.ErrDataNoRecord
		}
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	entity := entity.Permission{
		PermissionID:         e.PermissionID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
		PermissionCategoryID: e.PermissionCategoryID,
	}
	return entity, nil
}

// UpdatePermission 権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermission(
	ctx context.Context, permissionID uuid.UUID, param parameter.UpdatePermissionParams,
) (entity.Permission, error) {
	e, err := updatePermission(ctx, a.query, permissionID, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	return e, nil
}

// UpdatePermissionWithSd SD付きで権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionWithSd(
	ctx context.Context, sd store.Sd, permissionID uuid.UUID, param parameter.UpdatePermissionParams,
) (entity.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Permission{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePermission(ctx, qtx, permissionID, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	return e, nil
}

func updatePermissionByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdatePermissionByKeyParams,
) (entity.Permission, error) {
	p := query.UpdatePermissionByKeyParams{
		Key:                  key,
		Name:                 param.Name,
		Description:          param.Description,
		PermissionCategoryID: param.PermissionCategoryID,
	}
	e, err := qtx.UpdatePermissionByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Permission{}, store.ErrDataNoRecord
		}
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	entity := entity.Permission{
		PermissionID:         e.PermissionID,
		Name:                 e.Name,
		Key:                  e.Key,
		Description:          e.Description,
		PermissionCategoryID: e.PermissionCategoryID,
	}
	return entity, nil
}

// UpdatePermissionByKey 権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionByKey(
	ctx context.Context, key string, param parameter.UpdatePermissionByKeyParams,
) (entity.Permission, error) {
	e, err := updatePermissionByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	return e, nil
}

// UpdatePermissionByKeyWithSd SD付きで権限カテゴリーを更新する。
func (a *PgAdapter) UpdatePermissionByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdatePermissionByKeyParams,
) (entity.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Permission{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePermissionByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	return e, nil
}

package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
)

func countRoles(
	ctx context.Context, qtx *query.Queries, where parameter.WhereRoleParam,
) (int64, error) {
	p := query.CountRolesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountRoles(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count role: %w", err)
	}
	return c, nil
}

// CountRoles ロール数を取得する。
func (a *PgAdapter) CountRoles(ctx context.Context, where parameter.WhereRoleParam) (int64, error) {
	return countRoles(ctx, a.query, where)
}

// CountRolesWithSd SD付きでロール数を取得する。
func (a *PgAdapter) CountRolesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereRoleParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countRoles(ctx, qtx, where)
}

func createRole(
	ctx context.Context, qtx *query.Queries, param parameter.CreateRoleParam, now time.Time,
) (entity.Role, error) {
	p := query.CreateRoleParams{
		Name:        param.Name,
		Description: param.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	e, err := qtx.CreateRole(ctx, p)
	if err != nil {
		return entity.Role{}, fmt.Errorf("failed to create role: %w", err)
	}
	entity := entity.Role{
		RoleID:      e.RoleID,
		Name:        e.Name,
		Description: e.Description,
	}
	return entity, nil
}

// CreateRole ロールを作成する。
func (a *PgAdapter) CreateRole(
	ctx context.Context, param parameter.CreateRoleParam,
) (entity.Role, error) {
	return createRole(ctx, a.query, param, a.clocker.Now())
}

// CreateRoleWithSd SD付きでロールを作成する。
func (a *PgAdapter) CreateRoleWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateRoleParam,
) (entity.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Role{}, store.ErrNotFoundDescriptor
	}
	return createRole(ctx, qtx, param, a.clocker.Now())
}

func createRoles(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateRoleParam, now time.Time,
) (int64, error) {
	p := make([]query.CreateRolesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateRolesParams{
			Name:        param.Name,
			Description: param.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
	c, err := qtx.CreateRoles(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create roles: %w", err)
	}
	return c, nil
}

// CreateRoles ロールを作成する。
func (a *PgAdapter) CreateRoles(
	ctx context.Context, params []parameter.CreateRoleParam,
) (int64, error) {
	return createRoles(ctx, a.query, params, a.clocker.Now())
}

// CreateRolesWithSd SD付きでロールを作成する。
func (a *PgAdapter) CreateRolesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateRoleParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createRoles(ctx, qtx, params, a.clocker.Now())
}

func deleteRole(ctx context.Context, qtx *query.Queries, roleID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteRole(ctx, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete role: %w", err)
	}
	return c, nil
}

// DeleteRole ロールを削除する。
func (a *PgAdapter) DeleteRole(ctx context.Context, roleID uuid.UUID) (int64, error) {
	return deleteRole(ctx, a.query, roleID)
}

// DeleteRoleWithSd SD付きでロールを削除する。
func (a *PgAdapter) DeleteRoleWithSd(
	ctx context.Context, sd store.Sd, roleID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteRole(ctx, qtx, roleID)
}

func pluralDeleteRoles(ctx context.Context, qtx *query.Queries, roleIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteRoles(ctx, roleIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete roles: %w", err)
	}
	return c, nil
}

// PluralDeleteRoles ロールを複数削除する。
func (a *PgAdapter) PluralDeleteRoles(ctx context.Context, roleIDs []uuid.UUID) (int64, error) {
	return pluralDeleteRoles(ctx, a.query, roleIDs)
}

// PluralDeleteRolesWithSd SD付きでロールを複数削除する。
func (a *PgAdapter) PluralDeleteRolesWithSd(
	ctx context.Context, sd store.Sd, roleIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteRoles(ctx, qtx, roleIDs)
}

func findRoleByID(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID,
) (entity.Role, error) {
	e, err := qtx.FindRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Role{}, errhandle.NewModelNotFoundError("role")
		}
		return entity.Role{}, fmt.Errorf("failed to find role: %w", err)
	}
	entity := entity.Role{
		RoleID:      e.RoleID,
		Name:        e.Name,
		Description: e.Description,
	}
	return entity, nil
}

// FindRoleByID ロールを取得する。
func (a *PgAdapter) FindRoleByID(
	ctx context.Context, roleID uuid.UUID,
) (entity.Role, error) {
	return findRoleByID(ctx, a.query, roleID)
}

// FindRoleByIDWithSd SD付きでロールを取得する。
func (a *PgAdapter) FindRoleByIDWithSd(
	ctx context.Context, sd store.Sd, roleID uuid.UUID,
) (entity.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Role{}, store.ErrNotFoundDescriptor
	}
	return findRoleByID(ctx, qtx, roleID)
}

// RoleCursor is a cursor for Role.
type RoleCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getRoles(
	ctx context.Context, qtx *query.Queries, where parameter.WhereRoleParam,
	order parameter.RoleOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Role], error) {
	eConvFunc := func(e query.Role) (entity.Role, error) {
		return entity.Role{
			RoleID:      e.RoleID,
			Name:        e.Name,
			Description: e.Description,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountRolesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountRoles(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count roles: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Role, error) {
		p := query.GetRolesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetRoles(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Role{}, nil
			}
			return nil, fmt.Errorf("failed to get roles: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Role, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.RoleNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetRolesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetRolesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Role, error) {
		p := query.GetRolesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetRolesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Role) (entity.Int, any) {
		switch subCursor {
		case parameter.RoleDefaultCursorKey:
			return entity.Int(e.MRolesPkey), nil
		case parameter.RoleNameCursorKey:
			return entity.Int(e.MRolesPkey), e.Name
		}
		return entity.Int(e.MRolesPkey), nil
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
		return store.ListResult[entity.Role]{}, fmt.Errorf("failed to get roles: %w", err)
	}
	return res, nil
}

// GetRoles ロールを取得する。
func (a *PgAdapter) GetRoles(
	ctx context.Context,
	where parameter.WhereRoleParam,
	order parameter.RoleOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Role], error) {
	return getRoles(ctx, a.query, where, order, np, cp, wc)
}

// GetRolesWithSd SD付きでロールを取得する。
func (a *PgAdapter) GetRolesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereRoleParam,
	order parameter.RoleOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Role], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Role]{}, store.ErrNotFoundDescriptor
	}
	return getRoles(ctx, qtx, where, order, np, cp, wc)
}

func getPluralRoles(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Role], error) {
	var e []query.Role
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralRoles(ctx, ids)
	} else {
		p := query.GetPluralRolesUseNumberedPaginateParams{
			RoleIds: ids,
			Offset:  int32(np.Offset.Int64),
			Limit:   int32(np.Limit.Int64),
		}
		e, err = qtx.GetPluralRolesUseNumberedPaginate(ctx, p)
	}
	if err != nil {
		return store.ListResult[entity.Role]{}, fmt.Errorf("failed to get plural roles: %w", err)
	}
	entities := make([]entity.Role, len(e))
	for i, v := range e {
		entities[i] = entity.Role{
			RoleID:      v.RoleID,
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return store.ListResult[entity.Role]{Data: entities}, nil
}

// GetPluralRoles ロールを取得する。
func (a *PgAdapter) GetPluralRoles(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Role], error) {
	return getPluralRoles(ctx, a.query, ids, np)
}

// GetPluralRolesWithSd SD付きでロールを取得する。
func (a *PgAdapter) GetPluralRolesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Role], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Role]{}, store.ErrNotFoundDescriptor
	}
	return getPluralRoles(ctx, qtx, ids, np)
}

func updateRole(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID, param parameter.UpdateRoleParams, now time.Time,
) (entity.Role, error) {
	p := query.UpdateRoleParams{
		RoleID:      roleID,
		Name:        param.Name,
		Description: param.Description,
		UpdatedAt:   now,
	}
	e, err := qtx.UpdateRole(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Role{}, errhandle.NewModelNotFoundError("role")
		}
		return entity.Role{}, fmt.Errorf("failed to update role: %w", err)
	}
	entity := entity.Role{
		RoleID:      e.RoleID,
		Name:        e.Name,
		Description: e.Description,
	}
	return entity, nil
}

// UpdateRole ロールを更新する。
func (a *PgAdapter) UpdateRole(
	ctx context.Context, roleID uuid.UUID, param parameter.UpdateRoleParams,
) (entity.Role, error) {
	return updateRole(ctx, a.query, roleID, param, a.clocker.Now())
}

// UpdateRoleWithSd SD付きでロールを更新する。
func (a *PgAdapter) UpdateRoleWithSd(
	ctx context.Context, sd store.Sd, roleID uuid.UUID, param parameter.UpdateRoleParams,
) (entity.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Role{}, store.ErrNotFoundDescriptor
	}
	return updateRole(ctx, qtx, roleID, param, a.clocker.Now())
}

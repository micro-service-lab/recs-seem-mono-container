package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageRole ロール管理サービス。
type ManageRole struct {
	DB store.Store
}

// CreateRole ロールを作成する。
func (m *ManageRole) CreateRole(
	ctx context.Context,
	name, description string,
) (entity.Role, error) {
	p := parameter.CreateRoleParam{
		Name:        name,
		Description: description,
	}
	e, err := m.DB.CreateRole(ctx, p)
	if err != nil {
		return entity.Role{}, fmt.Errorf("failed to create role: %w", err)
	}
	return e, nil
}

// CreateRoles ロールを複数作成する。
func (m *ManageRole) CreateRoles(
	ctx context.Context, ps []parameter.CreateRoleParam,
) (int64, error) {
	es, err := m.DB.CreateRoles(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create roles: %w", err)
	}
	return es, nil
}

// UpdateRole ロールを更新する。
func (m *ManageRole) UpdateRole(
	ctx context.Context, id uuid.UUID, name, description string,
) (entity.Role, error) {
	p := parameter.UpdateRoleParams{
		Name:        name,
		Description: description,
	}
	e, err := m.DB.UpdateRole(ctx, id, p)
	if err != nil {
		return entity.Role{}, fmt.Errorf("failed to update role: %w", err)
	}
	return e, nil
}

// DeleteRole ロールを削除する。
func (m *ManageRole) DeleteRole(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteRole(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete role: %w", err)
	}
	return c, nil
}

// PluralDeleteRoles ロールを複数削除する。
func (m *ManageRole) PluralDeleteRoles(ctx context.Context, ids []uuid.UUID) (int64, error) {
	c, err := m.DB.PluralDeleteRoles(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete roles: %w", err)
	}
	return c, nil
}

// FindRoleByID ロールをIDで取得する。
func (m *ManageRole) FindRoleByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Role, error) {
	e, err := m.DB.FindRoleByID(ctx, id)
	if err != nil {
		return entity.Role{}, fmt.Errorf("failed to find role by id: %w", err)
	}
	return e, nil
}

// GetRoles ロールを取得する。
func (m *ManageRole) GetRoles(
	ctx context.Context,
	whereSearchName string,
	order parameter.RoleOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Role], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereRoleParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetRoles(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Role]{}, fmt.Errorf("failed to get roles: %w", err)
	}
	return r, nil
}

// GetRolesCount ロールの数を取得する。
func (m *ManageRole) GetRolesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereRoleParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountRoles(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get roles count: %w", err)
	}
	return c, nil
}

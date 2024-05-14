package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Role ロールを表すインターフェース。
type Role interface {
	// CountRoles ロール数を取得する。
	CountRoles(ctx context.Context, where parameter.WhereRoleParam) (int64, error)
	// CountRolesWithSd SD付きでロール数を取得する。
	CountRolesWithSd(ctx context.Context, sd Sd, where parameter.WhereRoleParam) (int64, error)
	// CreateRole ロールを作成する。
	CreateRole(ctx context.Context, param parameter.CreateRoleParam) (entity.Role, error)
	// CreateRoleWithSd SD付きでロールを作成する。
	CreateRoleWithSd(
		ctx context.Context, sd Sd, param parameter.CreateRoleParam) (entity.Role, error)
	// CreateRoles ロールを作成する。
	CreateRoles(ctx context.Context, params []parameter.CreateRoleParam) (int64, error)
	// CreateRolesWithSd SD付きでロールを作成する。
	CreateRolesWithSd(ctx context.Context, sd Sd, params []parameter.CreateRoleParam) (int64, error)
	// DeleteRole ロールを削除する。
	DeleteRole(ctx context.Context, roleID uuid.UUID) (int64, error)
	// DeleteRoleWithSd SD付きでロールを削除する。
	DeleteRoleWithSd(ctx context.Context, sd Sd, roleID uuid.UUID) (int64, error)
	// PluralDeleteRoles ロールを複数削除する。
	PluralDeleteRoles(ctx context.Context, roleIDs []uuid.UUID) (int64, error)
	// PluralDeleteRolesWithSd SD付きでロールを複数削除する。
	PluralDeleteRolesWithSd(ctx context.Context, sd Sd, roleIDs []uuid.UUID) (int64, error)
	// FindRoleByID ロールを取得する。
	FindRoleByID(ctx context.Context, roleID uuid.UUID) (entity.Role, error)
	// FindRoleByIDWithSd SD付きでロールを取得する。
	FindRoleByIDWithSd(ctx context.Context, sd Sd, roleID uuid.UUID) (entity.Role, error)
	// GetRoles ロールを取得する。
	GetRoles(
		ctx context.Context,
		where parameter.WhereRoleParam,
		order parameter.RoleOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Role], error)
	// GetRolesWithSd SD付きでロールを取得する。
	GetRolesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereRoleParam,
		order parameter.RoleOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Role], error)
	// GetPluralRoles ロールを取得する。
	GetPluralRoles(
		ctx context.Context,
		roleIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Role], error)
	// GetPluralRolesWithSd SD付きでロールを取得する。
	GetPluralRolesWithSd(
		ctx context.Context,
		sd Sd,
		roleIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Role], error)
	// UpdateRole ロールを更新する。
	UpdateRole(
		ctx context.Context,
		roleID uuid.UUID,
		param parameter.UpdateRoleParams,
	) (entity.Role, error)
	// UpdateRoleWithSd SD付きでロールを更新する。
	UpdateRoleWithSd(
		ctx context.Context, sd Sd, roleID uuid.UUID,
		param parameter.UpdateRoleParams) (entity.Role, error)
}

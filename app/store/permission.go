package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Permission 権限を表すインターフェース。
type Permission interface {
	// CountPermissions 権限数を取得する。
	CountPermissions(ctx context.Context, where parameter.WherePermissionParam) (int64, error)
	// CountPermissionsWithSd SD付きで権限数を取得する。
	CountPermissionsWithSd(
		ctx context.Context, sd Sd, where parameter.WherePermissionParam) (int64, error)
	// CreatePermission 権限を作成する。
	CreatePermission(
		ctx context.Context, param parameter.CreatePermissionParam) (entity.Permission, error)
	// CreatePermissionWithSd SD付きで権限を作成する。
	CreatePermissionWithSd(
		ctx context.Context, sd Sd, param parameter.CreatePermissionParam) (entity.Permission, error)
	// CreatePermissions 権限を作成する。
	CreatePermissions(ctx context.Context, params []parameter.CreatePermissionParam) (int64, error)
	// CreatePermissionsWithSd SD付きで権限を作成する。
	CreatePermissionsWithSd(
		ctx context.Context, sd Sd, params []parameter.CreatePermissionParam) (int64, error)
	// DeletePermission 権限を削除する。
	DeletePermission(ctx context.Context, permissionID uuid.UUID) error
	// DeletePermissionWithSd SD付きで権限を削除する。
	DeletePermissionWithSd(ctx context.Context, sd Sd, permissionID uuid.UUID) error
	// DeletePermissionByKey 権限を削除する。
	DeletePermissionByKey(ctx context.Context, key string) error
	// DeletePermissionByKeyWithSd SD付きで権限を削除する。
	DeletePermissionByKeyWithSd(ctx context.Context, sd Sd, key string) error
	// PluralDeletePermissions 権限を複数削除する。
	PluralDeletePermissions(ctx context.Context, permissionIDs []uuid.UUID) error
	// PluralDeletePermissionsWithSd SD付きで権限を複数削除する。
	PluralDeletePermissionsWithSd(ctx context.Context, sd Sd, permissionIDs []uuid.UUID) error
	// FindPermissionByID 権限を取得する。
	FindPermissionByID(ctx context.Context, permissionID uuid.UUID) (entity.Permission, error)
	// FindPermissionByIDWithSd SD付きで権限を取得する。
	FindPermissionByIDWithSd(
		ctx context.Context, sd Sd, permissionID uuid.UUID) (entity.Permission, error)
	// FindPermissionByKey 権限を取得する。
	FindPermissionByKey(ctx context.Context, key string) (entity.Permission, error)
	// FindPermissionByKeyWithSd SD付きで権限を取得する。
	FindPermissionByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.Permission, error)
	// GetPermissions 権限を取得する。
	GetPermissions(
		ctx context.Context,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Permission], error)
	// GetPermissionsWithSd SD付きで権限を取得する。
	GetPermissionsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Permission], error)
	// GetPermissionsWithCategory 権限とそのカテゴリを取得する。
	GetPermissionsWithCategory(
		ctx context.Context,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PermissionWithCategory], error)
	// GetPermissionsWithCategoryWithSd SD付きで権限とそのカテゴリを取得する。
	GetPermissionsWithCategoryWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PermissionWithCategory], error)
	// GetPluralPermissions 権限を取得する。
	GetPluralPermissions(
		ctx context.Context,
		PermissionIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Permission], error)
	// GetPluralPermissionsWithSd SD付きで権限を取得する。
	GetPluralPermissionsWithSd(
		ctx context.Context,
		sd Sd,
		PermissionIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Permission], error)
	// UpdatePermission 権限を更新する。
	UpdatePermission(
		ctx context.Context,
		permissionID uuid.UUID,
		param parameter.UpdatePermissionParams,
	) (entity.Permission, error)
	// UpdatePermissionWithSd SD付きで権限を更新する。
	UpdatePermissionWithSd(
		ctx context.Context, sd Sd, permissionID uuid.UUID,
		param parameter.UpdatePermissionParams) (entity.Permission, error)
	// UpdatePermissionByKey 権限を更新する。
	UpdatePermissionByKey(
		ctx context.Context,
		key string, param parameter.UpdatePermissionByKeyParams) (entity.Permission, error)
	// UpdatePermissionByKeyWithSd SD付きで権限を更新する。
	UpdatePermissionByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdatePermissionByKeyParams,
	) (entity.Permission, error)
}

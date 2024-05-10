package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Permission 権限を表すインターフェース。
type Permission interface {
	// CountPermissionCategories 権限数を取得する。
	CountPermissionCategories(ctx context.Context, where parameter.WherePermissionParam) (int64, error)
	// CountPermissionCategoriesWithSd SD付きで権限数を取得する。
	CountPermissionCategoriesWithSd(
		ctx context.Context, sd Sd, where parameter.WherePermissionParam) (int64, error)
	// CreatePermission 権限を作成する。
	CreatePermission(
		ctx context.Context, param parameter.CreatePermissionParam) (entity.Permission, error)
	// CreatePermissionWithSd SD付きで権限を作成する。
	CreatePermissionWithSd(
		ctx context.Context, sd Sd, param parameter.CreatePermissionParam) (entity.Permission, error)
	// CreatePermissionCategories 権限を作成する。
	CreatePermissionCategories(ctx context.Context, params []parameter.CreatePermissionParam) (int64, error)
	// CreatePermissionCategoriesWithSd SD付きで権限を作成する。
	CreatePermissionCategoriesWithSd(
		ctx context.Context, sd Sd, params []parameter.CreatePermissionParam) (int64, error)
	// DeletePermission 権限を削除する。
	DeletePermission(ctx context.Context, permissionID uuid.UUID) error
	// DeletePermissionWithSd SD付きで権限を削除する。
	DeletePermissionWithSd(ctx context.Context, sd Sd, permissionID uuid.UUID) error
	// DeletePermissionByKey 権限を削除する。
	DeletePermissionByKey(ctx context.Context, key string) error
	// DeletePermissionByKeyWithSd SD付きで権限を削除する。
	DeletePermissionByKeyWithSd(ctx context.Context, sd Sd, key string) error
	// PluralDeletePermissionCategories 権限を複数削除する。
	PluralDeletePermissionCategories(ctx context.Context, permissionIDs []uuid.UUID) error
	// PluralDeletePermissionCategoriesWithSd SD付きで権限を複数削除する。
	PluralDeletePermissionCategoriesWithSd(ctx context.Context, sd Sd, permissionIDs []uuid.UUID) error
	// FindPermissionByID 権限を取得する。
	FindPermissionByID(ctx context.Context, permissionID uuid.UUID) (entity.Permission, error)
	// FindPermissionByIDWithSd SD付きで権限を取得する。
	FindPermissionByIDWithSd(
		ctx context.Context, sd Sd, permissionID uuid.UUID) (entity.Permission, error)
	// FindPermissionByKey 権限を取得する。
	FindPermissionByKey(ctx context.Context, key string) (entity.Permission, error)
	// FindPermissionByKeyWithSd SD付きで権限を取得する。
	FindPermissionByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.Permission, error)
	// GetPermissionCategories 権限を取得する。
	GetPermissionCategories(
		ctx context.Context,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Permission], error)
	// GetPermissionCategoriesWithSd SD付きで権限を取得する。
	GetPermissionCategoriesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePermissionParam,
		order parameter.PermissionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Permission], error)
	// GetPluralPermissionCategories 権限を取得する。
	GetPluralPermissionCategories(
		ctx context.Context,
		PermissionIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Permission], error)
	// GetPluralPermissionCategoriesWithSd SD付きで権限を取得する。
	GetPluralPermissionCategoriesWithSd(
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

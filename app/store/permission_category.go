package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// PermissionCategory 権限カテゴリーを表すインターフェース。
type PermissionCategory interface {
	// CountPermissionCategories 権限カテゴリー数を取得する。
	CountPermissionCategories(ctx context.Context, where parameter.WherePermissionCategoryParam) (int64, error)
	// CountPermissionCategoriesWithSd SD付きで権限カテゴリー数を取得する。
	CountPermissionCategoriesWithSd(
		ctx context.Context, sd Sd, where parameter.WherePermissionCategoryParam) (int64, error)
	// CreatePermissionCategory 権限カテゴリーを作成する。
	CreatePermissionCategory(
		ctx context.Context, param parameter.CreatePermissionCategoryParam) (entity.PermissionCategory, error)
	// CreatePermissionCategoryWithSd SD付きで権限カテゴリーを作成する。
	CreatePermissionCategoryWithSd(
		ctx context.Context, sd Sd, param parameter.CreatePermissionCategoryParam) (entity.PermissionCategory, error)
	// CreatePermissionCategories 権限カテゴリーを作成する。
	CreatePermissionCategories(ctx context.Context, params []parameter.CreatePermissionCategoryParam) (int64, error)
	// CreatePermissionCategoriesWithSd SD付きで権限カテゴリーを作成する。
	CreatePermissionCategoriesWithSd(
		ctx context.Context, sd Sd, params []parameter.CreatePermissionCategoryParam) (int64, error)
	// DeletePermissionCategory 権限カテゴリーを削除する。
	DeletePermissionCategory(ctx context.Context, permissionCategoryID uuid.UUID) (int64, error)
	// DeletePermissionCategoryWithSd SD付きで権限カテゴリーを削除する。
	DeletePermissionCategoryWithSd(ctx context.Context, sd Sd, permissionCategoryID uuid.UUID) (int64, error)
	// DeletePermissionCategoryByKey 権限カテゴリーを削除する。
	DeletePermissionCategoryByKey(ctx context.Context, key string) (int64, error)
	// DeletePermissionCategoryByKeyWithSd SD付きで権限カテゴリーを削除する。
	DeletePermissionCategoryByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeletePermissionCategories 権限カテゴリーを複数削除する。
	PluralDeletePermissionCategories(ctx context.Context, permissionCategoryIDs []uuid.UUID) (int64, error)
	// PluralDeletePermissionCategoriesWithSd SD付きで権限カテゴリーを複数削除する。
	PluralDeletePermissionCategoriesWithSd(ctx context.Context, sd Sd, permissionCategoryIDs []uuid.UUID) (int64, error)
	// FindPermissionCategoryByID 権限カテゴリーを取得する。
	FindPermissionCategoryByID(ctx context.Context, permissionCategoryID uuid.UUID) (entity.PermissionCategory, error)
	// FindPermissionCategoryByIDWithSd SD付きで権限カテゴリーを取得する。
	FindPermissionCategoryByIDWithSd(
		ctx context.Context, sd Sd, permissionCategoryID uuid.UUID) (entity.PermissionCategory, error)
	// FindPermissionCategoryByKey 権限カテゴリーを取得する。
	FindPermissionCategoryByKey(ctx context.Context, key string) (entity.PermissionCategory, error)
	// FindPermissionCategoryByKeyWithSd SD付きで権限カテゴリーを取得する。
	FindPermissionCategoryByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.PermissionCategory, error)
	// GetPermissionCategories 権限カテゴリーを取得する。
	GetPermissionCategories(
		ctx context.Context,
		where parameter.WherePermissionCategoryParam,
		order parameter.PermissionCategoryOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PermissionCategory], error)
	// GetPermissionCategoriesWithSd SD付きで権限カテゴリーを取得する。
	GetPermissionCategoriesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePermissionCategoryParam,
		order parameter.PermissionCategoryOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PermissionCategory], error)
	// GetPluralPermissionCategories 権限カテゴリーを取得する。
	GetPluralPermissionCategories(
		ctx context.Context,
		permissionCategoryIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.PermissionCategory], error)
	// GetPluralPermissionCategoriesWithSd SD付きで権限カテゴリーを取得する。
	GetPluralPermissionCategoriesWithSd(
		ctx context.Context,
		sd Sd,
		permissionCategoryIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.PermissionCategory], error)
	// UpdatePermissionCategory 権限カテゴリーを更新する。
	UpdatePermissionCategory(
		ctx context.Context,
		permissionCategoryID uuid.UUID,
		param parameter.UpdatePermissionCategoryParams,
	) (entity.PermissionCategory, error)
	// UpdatePermissionCategoryWithSd SD付きで権限カテゴリーを更新する。
	UpdatePermissionCategoryWithSd(
		ctx context.Context, sd Sd, permissionCategoryID uuid.UUID,
		param parameter.UpdatePermissionCategoryParams) (entity.PermissionCategory, error)
	// UpdatePermissionCategoryByKey 権限カテゴリーを更新する。
	UpdatePermissionCategoryByKey(
		ctx context.Context,
		key string, param parameter.UpdatePermissionCategoryByKeyParams) (entity.PermissionCategory, error)
	// UpdatePermissionCategoryByKeyWithSd SD付きで権限カテゴリーを更新する。
	UpdatePermissionCategoryByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdatePermissionCategoryByKeyParams,
	) (entity.PermissionCategory, error)
}

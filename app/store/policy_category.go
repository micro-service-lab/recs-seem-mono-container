package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// PolicyCategory ポリシーカテゴリーを表すインターフェース。
type PolicyCategory interface {
	// CountPolicyCategories ポリシーカテゴリー数を取得する。
	CountPolicyCategories(ctx context.Context, where parameter.WherePolicyCategoryParam) (int64, error)
	// CountPolicyCategoriesWithSd SD付きでポリシーカテゴリー数を取得する。
	CountPolicyCategoriesWithSd(
		ctx context.Context, sd Sd, where parameter.WherePolicyCategoryParam) (int64, error)
	// CreatePolicyCategory ポリシーカテゴリーを作成する。
	CreatePolicyCategory(
		ctx context.Context, param parameter.CreatePolicyCategoryParam) (entity.PolicyCategory, error)
	// CreatePolicyCategoryWithSd SD付きでポリシーカテゴリーを作成する。
	CreatePolicyCategoryWithSd(
		ctx context.Context, sd Sd, param parameter.CreatePolicyCategoryParam) (entity.PolicyCategory, error)
	// CreatePolicyCategories ポリシーカテゴリーを作成する。
	CreatePolicyCategories(ctx context.Context, params []parameter.CreatePolicyCategoryParam) (int64, error)
	// CreatePolicyCategoriesWithSd SD付きでポリシーカテゴリーを作成する。
	CreatePolicyCategoriesWithSd(
		ctx context.Context, sd Sd, params []parameter.CreatePolicyCategoryParam) (int64, error)
	// DeletePolicyCategory ポリシーカテゴリーを削除する。
	DeletePolicyCategory(ctx context.Context, policyCategoryID uuid.UUID) error
	// DeletePolicyCategoryWithSd SD付きでポリシーカテゴリーを削除する。
	DeletePolicyCategoryWithSd(ctx context.Context, sd Sd, policyCategoryID uuid.UUID) error
	// DeletePolicyCategoryByKey ポリシーカテゴリーを削除する。
	DeletePolicyCategoryByKey(ctx context.Context, key string) error
	// DeletePolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを削除する。
	DeletePolicyCategoryByKeyWithSd(ctx context.Context, sd Sd, key string) error
	// PluralDeletePolicyCategories ポリシーカテゴリーを複数削除する。
	PluralDeletePolicyCategories(ctx context.Context, policyCategoryIDs []uuid.UUID) error
	// PluralDeletePolicyCategoriesWithSd SD付きでポリシーカテゴリーを複数削除する。
	PluralDeletePolicyCategoriesWithSd(ctx context.Context, sd Sd, policyCategoryIDs []uuid.UUID) error
	// FindPolicyCategoryByID ポリシーカテゴリーを取得する。
	FindPolicyCategoryByID(ctx context.Context, policyCategoryID uuid.UUID) (entity.PolicyCategory, error)
	// FindPolicyCategoryByIDWithSd SD付きでポリシーカテゴリーを取得する。
	FindPolicyCategoryByIDWithSd(
		ctx context.Context, sd Sd, policyCategoryID uuid.UUID) (entity.PolicyCategory, error)
	// FindPolicyCategoryByKey ポリシーカテゴリーを取得する。
	FindPolicyCategoryByKey(ctx context.Context, key string) (entity.PolicyCategory, error)
	// FindPolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを取得する。
	FindPolicyCategoryByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.PolicyCategory, error)
	// GetPolicyCategories ポリシーカテゴリーを取得する。
	GetPolicyCategories(
		ctx context.Context,
		where parameter.WherePolicyCategoryParam,
		order parameter.PolicyCategoryOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyCategory], error)
	// GetPolicyCategoriesWithSd SD付きでポリシーカテゴリーを取得する。
	GetPolicyCategoriesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePolicyCategoryParam,
		order parameter.PolicyCategoryOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyCategory], error)
	// GetPluralPolicyCategories ポリシーカテゴリーを取得する。
	GetPluralPolicyCategories(
		ctx context.Context,
		policyCategoryIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.PolicyCategory], error)
	// GetPluralPolicyCategoriesWithSd SD付きでポリシーカテゴリーを取得する。
	GetPluralPolicyCategoriesWithSd(
		ctx context.Context,
		sd Sd,
		policyCategoryIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.PolicyCategory], error)
	// UpdatePolicyCategory ポリシーカテゴリーを更新する。
	UpdatePolicyCategory(
		ctx context.Context,
		policyCategoryID uuid.UUID,
		param parameter.UpdatePolicyCategoryParams,
	) (entity.PolicyCategory, error)
	// UpdatePolicyCategoryWithSd SD付きでポリシーカテゴリーを更新する。
	UpdatePolicyCategoryWithSd(
		ctx context.Context, sd Sd, policyCategoryID uuid.UUID,
		param parameter.UpdatePolicyCategoryParams) (entity.PolicyCategory, error)
	// UpdatePolicyCategoryByKey ポリシーカテゴリーを更新する。
	UpdatePolicyCategoryByKey(
		ctx context.Context,
		key string, param parameter.UpdatePolicyCategoryByKeyParams) (entity.PolicyCategory, error)
	// UpdatePolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを更新する。
	UpdatePolicyCategoryByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdatePolicyCategoryByKeyParams,
	) (entity.PolicyCategory, error)
}

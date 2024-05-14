package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Policy ポリシーを表すインターフェース。
type Policy interface {
	// CountPolicies ポリシー数を取得する。
	CountPolicies(ctx context.Context, where parameter.WherePolicyParam) (int64, error)
	// CountPoliciesWithSd SD付きでポリシー数を取得する。
	CountPoliciesWithSd(
		ctx context.Context, sd Sd, where parameter.WherePolicyParam) (int64, error)
	// CreatePolicy ポリシーを作成する。
	CreatePolicy(
		ctx context.Context, param parameter.CreatePolicyParam) (entity.Policy, error)
	// CreatePolicyWithSd SD付きでポリシーを作成する。
	CreatePolicyWithSd(
		ctx context.Context, sd Sd, param parameter.CreatePolicyParam) (entity.Policy, error)
	// CreatePolicies ポリシーを作成する。
	CreatePolicies(ctx context.Context, params []parameter.CreatePolicyParam) (int64, error)
	// CreatePoliciesWithSd SD付きでポリシーを作成する。
	CreatePoliciesWithSd(
		ctx context.Context, sd Sd, params []parameter.CreatePolicyParam) (int64, error)
	// DeletePolicy ポリシーを削除する。
	DeletePolicy(ctx context.Context, policyID uuid.UUID) (int64, error)
	// DeletePolicyWithSd SD付きでポリシーを削除する。
	DeletePolicyWithSd(ctx context.Context, sd Sd, policyID uuid.UUID) (int64, error)
	// DeletePolicyByKey ポリシーを削除する。
	DeletePolicyByKey(ctx context.Context, key string) (int64, error)
	// DeletePolicyByKeyWithSd SD付きでポリシーを削除する。
	DeletePolicyByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeletePolicies ポリシーを複数削除する。
	PluralDeletePolicies(ctx context.Context, policyIDs []uuid.UUID) (int64, error)
	// PluralDeletePoliciesWithSd SD付きでポリシーを複数削除する。
	PluralDeletePoliciesWithSd(ctx context.Context, sd Sd, policyIDs []uuid.UUID) (int64, error)
	// FindPolicyByID ポリシーを取得する。
	FindPolicyByID(ctx context.Context, policyID uuid.UUID) (entity.Policy, error)
	// FindPolicyByIDWithSd SD付きでポリシーを取得する。
	FindPolicyByIDWithSd(
		ctx context.Context, sd Sd, policyID uuid.UUID) (entity.Policy, error)
	// FindPolicyByIDWithCategory ポリシーとそのカテゴリを取得する。
	FindPolicyByIDWithCategory(ctx context.Context, policyID uuid.UUID) (entity.PolicyWithCategory, error)
	// FindPolicyByIDWithCategoryWithSd SD付きでポリシーとそのカテゴリを取得する。
	FindPolicyByIDWithCategoryWithSd(
		ctx context.Context, sd Sd, policyID uuid.UUID) (entity.PolicyWithCategory, error)
	// FindPolicyByKey ポリシーを取得する。
	FindPolicyByKey(ctx context.Context, key string) (entity.Policy, error)
	// FindPolicyByKeyWithSd SD付きでポリシーを取得する。
	FindPolicyByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.Policy, error)
	// FindPolicyByKeyWithCategory ポリシーを取得する。
	FindPolicyByKeyWithCategory(ctx context.Context, key string) (entity.PolicyWithCategory, error)
	// FindPolicyByKeyWithCategoryWithSd SD付きでポリシーを取得する。
	FindPolicyByKeyWithCategoryWithSd(ctx context.Context, sd Sd, key string) (entity.PolicyWithCategory, error)
	// GetPolicies ポリシーを取得する。
	GetPolicies(
		ctx context.Context,
		where parameter.WherePolicyParam,
		order parameter.PolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Policy], error)
	// GetPoliciesWithSd SD付きでポリシーを取得する。
	GetPoliciesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePolicyParam,
		order parameter.PolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Policy], error)
	// GetPoliciesWithCategory ポリシーとそのカテゴリを取得する。
	GetPoliciesWithCategory(
		ctx context.Context,
		where parameter.WherePolicyParam,
		order parameter.PolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyWithCategory], error)
	// GetPoliciesWithCategoryWithSd SD付きでポリシーとそのカテゴリを取得する。
	GetPoliciesWithCategoryWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WherePolicyParam,
		order parameter.PolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyWithCategory], error)
	// GetPluralPolicies ポリシーを取得する。
	GetPluralPolicies(
		ctx context.Context,
		policyIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Policy], error)
	// GetPluralPoliciesWithSd SD付きでポリシーを取得する。
	GetPluralPoliciesWithSd(
		ctx context.Context,
		sd Sd,
		policyIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.Policy], error)
	// UpdatePolicy ポリシーを更新する。
	UpdatePolicy(
		ctx context.Context,
		policyID uuid.UUID,
		param parameter.UpdatePolicyParams,
	) (entity.Policy, error)
	// UpdatePolicyWithSd SD付きでポリシーを更新する。
	UpdatePolicyWithSd(
		ctx context.Context, sd Sd, policyID uuid.UUID,
		param parameter.UpdatePolicyParams) (entity.Policy, error)
	// UpdatePolicyByKey ポリシーを更新する。
	UpdatePolicyByKey(
		ctx context.Context,
		key string, param parameter.UpdatePolicyByKeyParams) (entity.Policy, error)
	// UpdatePolicyByKeyWithSd SD付きでポリシーを更新する。
	UpdatePolicyByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdatePolicyByKeyParams,
	) (entity.Policy, error)
}

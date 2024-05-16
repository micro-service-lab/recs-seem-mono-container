package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// RoleAssociation ロールの関連付けを表すインターフェース。
type RoleAssociation interface {
	// CountPoliciesOnRole ロールに関連付けられたポリシー数を取得する。
	CountPoliciesOnRole(ctx context.Context, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam) (int64, error)
	// CountPoliciesOnRoleWithSd SD付きでロールに関連付けられたポリシー数を取得する。
	CountPoliciesOnRoleWithSd(
		ctx context.Context, sd Sd, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam) (int64, error)
	// CountRolesOnPolicy ポリシーに関連付けられたロール数を取得する。
	CountRolesOnPolicy(ctx context.Context, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam) (int64, error)
	// CountRolesOnPolicyWithSd SD付きでポリシーに関連付けられたロール数を取得する。
	CountRolesOnPolicyWithSd(
		ctx context.Context, sd Sd, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam) (int64, error)
	// AssociateRole ロールを関連付ける。
	AssociateRole(ctx context.Context, param parameter.AssociationRoleParam) (entity.RoleAssociation, error)
	// AssociateRoleWithSd SD付きでロールを関連付ける。
	AssociateRoleWithSd(ctx context.Context, sd Sd, param parameter.AssociationRoleParam) (entity.RoleAssociation, error)
	// AssociateRoles ロールを複数関連付ける。
	AssociateRoles(ctx context.Context, params []parameter.AssociationRoleParam) (int64, error)
	// AssociateRolesWithSd SD付きでロールを複数関連付ける。
	AssociateRolesWithSd(ctx context.Context, sd Sd, params []parameter.AssociationRoleParam) (int64, error)
	// DisassociateRole ロールの関連付けを解除する。
	DisassociateRole(ctx context.Context, roleID, policyID uuid.UUID) (int64, error)
	// DisassociateRoleWithSd SD付きでロールの関連付けを解除する。
	DisassociateRoleWithSd(ctx context.Context, sd Sd, roleID, policyID uuid.UUID) (int64, error)
	// DisassociateRoleOnPolicy ポリシーに関連付けられたロールを解除する。
	DisassociateRoleOnPolicy(ctx context.Context, policyID uuid.UUID) (int64, error)
	// DisassociateRoleOnPolicyWithSd SD付きでポリシーに関連付けられたロールを解除する。
	DisassociateRoleOnPolicyWithSd(ctx context.Context, sd Sd, policyID uuid.UUID) (int64, error)
	// DisassociateRoleOnPolicies ポリシーに関連付けられたロールを複数解除する。
	DisassociateRoleOnPolicies(ctx context.Context, policyIDs []uuid.UUID) (int64, error)
	// DisassociateRoleOnPoliciesWithSd SD付きでポリシーに関連付けられたロールを複数解除する。
	DisassociateRoleOnPoliciesWithSd(ctx context.Context, sd Sd, policyIDs []uuid.UUID) (int64, error)
	// PluralDisassociateRoleOnPolicy ポリシーに関連付けられたロールを複数解除する。
	PluralDisassociateRoleOnPolicy(ctx context.Context, policyID uuid.UUID, roleIDs []uuid.UUID) (int64, error)
	// PluralDisassociateRoleOnPolicyWithSd SD付きでポリシーに関連付けられたロールを複数解除する。
	PluralDisassociateRoleOnPolicyWithSd(
		ctx context.Context, sd Sd, policyID uuid.UUID, roleIDs []uuid.UUID) (int64, error)
	// DisassociatePolicyOnRole ロールに関連付けられたポリシーを解除する。
	DisassociatePolicyOnRole(ctx context.Context, roleID uuid.UUID) (int64, error)
	// DisassociatePolicyOnRoleWithSd SD付きでロールに関連付けられたポリシーを解除する。
	DisassociatePolicyOnRoleWithSd(ctx context.Context, sd Sd, roleID uuid.UUID) (int64, error)
	// DisassociatePolicyOnRoles ロールに関連付けられたポリシーを複数解除する。
	DisassociatePolicyOnRoles(ctx context.Context, roleIDs []uuid.UUID) (int64, error)
	// DisassociatePolicyOnRolesWithSd SD付きでロールに関連付けられたポリシーを複数解除する。
	DisassociatePolicyOnRolesWithSd(ctx context.Context, sd Sd, roleIDs []uuid.UUID) (int64, error)
	// PluralDisassociatePolicyOnRole ロールに関連付けられたポリシーを複数解除する。
	PluralDisassociatePolicyOnRole(ctx context.Context, roleID uuid.UUID, policyIDs []uuid.UUID) (int64, error)
	// PluralDisassociatePolicyOnRoleWithSd SD付きでロールに関連付けられたポリシーを複数解除する。
	PluralDisassociatePolicyOnRoleWithSd(
		ctx context.Context, sd Sd, roleID uuid.UUID, policyIDs []uuid.UUID) (int64, error)
	// GetRolesOnPolicy ポリシーに関連付けられたロールを取得する。
	GetRolesOnPolicy(
		ctx context.Context,
		policyID uuid.UUID,
		where parameter.WhereRoleOnPolicyParam,
		order parameter.RoleOnPolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.RoleOnPolicy], error)
	// GetRolesOnPolicyWithSd SD付きでポリシーに関連付けられたロールを取得する。
	GetRolesOnPolicyWithSd(
		ctx context.Context,
		sd Sd,
		policyID uuid.UUID,
		where parameter.WhereRoleOnPolicyParam,
		order parameter.RoleOnPolicyOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.RoleOnPolicy], error)
	// GetPoliciesOnRole ロールに関連付けられたポリシーを取得する。
	GetPoliciesOnRole(
		ctx context.Context,
		roleID uuid.UUID,
		where parameter.WherePolicyOnRoleParam,
		order parameter.PolicyOnRoleOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyOnRole], error)
	// GetPoliciesOnRoleWithSd SD付きでロールに関連付けられたポリシーを取得する。
	GetPoliciesOnRoleWithSd(
		ctx context.Context,
		sd Sd,
		roleID uuid.UUID,
		where parameter.WherePolicyOnRoleParam,
		order parameter.PolicyOnRoleOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.PolicyOnRole], error)
}

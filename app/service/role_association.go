package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
)

// ManageRoleAssociation ロール紐付け管理サービス。
type ManageRoleAssociation struct {
	DB store.Store
}

// AssociateRole ロールを関連付ける。
func (m *ManageRoleAssociation) AssociateRole(
	ctx context.Context, roleID, policyID uuid.UUID,
) (e entity.RoleAssociation, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.RoleAssociation{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	np := store.NumberedPaginationParam{}
	var ro store.ListResult[entity.Role]
	if ro, err = m.DB.GetPluralRolesWithSd(ctx, sd, []uuid.UUID{roleID}, np); err != nil {
		return entity.RoleAssociation{}, fmt.Errorf("failed to get plural roles: %w", err)
	}
	if len(ro.Data) != 1 {
		return entity.RoleAssociation{}, errhandle.NewModelNotFoundError(AssociateRoleTargetRoles)
	}
	var po store.ListResult[entity.Policy]
	if po, err = m.DB.GetPluralPoliciesWithSd(ctx, sd, []uuid.UUID{policyID}, np); err != nil {
		return entity.RoleAssociation{}, fmt.Errorf("failed to get plural policies: %w", err)
	}
	if len(po.Data) != 1 {
		return entity.RoleAssociation{}, errhandle.NewModelNotFoundError(AssociateRoleTargetPolicies)
	}
	p := parameter.AssociationRoleParam{
		RoleID:   roleID,
		PolicyID: policyID,
	}
	e, err = m.DB.AssociateRoleWithSd(ctx, sd, p)
	if err != nil {
		return entity.RoleAssociation{}, fmt.Errorf("failed to associate role: %w", err)
	}
	return e, nil
}

// AssociateRoles ロールを複数関連付ける。
func (m *ManageRoleAssociation) AssociateRoles(
	ctx context.Context, params []parameter.AssociationRoleParam,
) (int64, error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	roleIDs := make([]uuid.UUID, len(params))
	policyIDs := make([]uuid.UUID, len(params))
	for i, p := range params {
		roleIDs[i] = p.RoleID
		policyIDs[i] = p.PolicyID
	}
	np := store.NumberedPaginationParam{}
	var ro store.ListResult[entity.Role]
	if ro, err = m.DB.GetPluralRolesWithSd(ctx, sd, roleIDs, np); err != nil {
		return 0, fmt.Errorf("failed to get plural roles: %w", err)
	}
	fmt.Println(roleIDs, policyIDs)
	if len(ro.Data) != 1 {
		return 0, errhandle.NewModelNotFoundError(AssociateRoleTargetRoles)
	}
	var po store.ListResult[entity.Policy]
	if po, err = m.DB.GetPluralPoliciesWithSd(ctx, sd, policyIDs, np); err != nil {
		return 0, fmt.Errorf("failed to get plural policies: %w", err)
	}
	if len(po.Data) != 1 {
		return 0, errhandle.NewModelNotFoundError(AssociateRoleTargetPolicies)
	}
	es, err := m.DB.AssociateRoles(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to associate roles: %w", err)
	}
	return es, nil
}

// DisassociateRole ロールの関連付けを解除する。
func (m *ManageRoleAssociation) DisassociateRole(
	ctx context.Context, roleID, policyID uuid.UUID,
) (int64, error) {
	es, err := m.DB.DisassociateRole(ctx, roleID, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role: %w", err)
	}
	return es, nil
}

// DisassociateRoleOnPolicy ポリシーに関連付けられたロールを解除する。
func (m *ManageRoleAssociation) DisassociateRoleOnPolicy(
	ctx context.Context, policyID uuid.UUID,
) (int64, error) {
	es, err := m.DB.DisassociateRoleOnPolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on policy: %w", err)
	}
	return es, nil
}

// DisassociateRoleOnPolicies 複数のポリシーに関連付けられたロールを解除する。
func (m *ManageRoleAssociation) DisassociateRoleOnPolicies(
	ctx context.Context, policyIDs []uuid.UUID,
) (int64, error) {
	es, err := m.DB.DisassociateRoleOnPolicies(ctx, policyIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on policies: %w", err)
	}
	return es, nil
}

// PluralDisassociateRoleOnPolicy ポリシーに関連付けられたロールを複数解除する。
func (m *ManageRoleAssociation) PluralDisassociateRoleOnPolicy(
	ctx context.Context, policyID uuid.UUID, roleIDs []uuid.UUID,
) (int64, error) {
	es, err := m.DB.PluralDisassociateRoleOnPolicy(ctx, policyID, roleIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural disassociate role on policy: %w", err)
	}
	return es, nil
}

// DisassociatePolicyOnRole ロールに関連付けられたポリシーを解除する。
func (m *ManageRoleAssociation) DisassociatePolicyOnRole(
	ctx context.Context, roleID uuid.UUID,
) (int64, error) {
	es, err := m.DB.DisassociatePolicyOnRole(ctx, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate policy on role: %w", err)
	}
	return es, nil
}

// DisassociatePolicyOnRoles ロールに関連付けられたポリシーを複数解除する。
func (m *ManageRoleAssociation) DisassociatePolicyOnRoles(
	ctx context.Context, roleIDs []uuid.UUID,
) (int64, error) {
	es, err := m.DB.DisassociatePolicyOnRoles(ctx, roleIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate policy on roles: %w", err)
	}
	return es, nil
}

// PluralDisassociatePolicyOnRole ロールに関連付けられたポリシーを複数解除する。
func (m *ManageRoleAssociation) PluralDisassociatePolicyOnRole(
	ctx context.Context, roleID uuid.UUID, policyIDs []uuid.UUID,
) (int64, error) {
	es, err := m.DB.PluralDisassociatePolicyOnRole(ctx, roleID, policyIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural disassociate policy on role: %w", err)
	}
	return es, nil
}

// GetRolesOnPolicy ポリシーに関連付けられたロールを取得する。
func (m *ManageRoleAssociation) GetRolesOnPolicy(
	ctx context.Context, policyID uuid.UUID,
	whereSearchName string,
	order parameter.RoleOnPolicyOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.RoleOnPolicy], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereRoleOnPolicyParam{
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
	es, err := m.DB.GetRolesOnPolicy(ctx, policyID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.RoleOnPolicy]{}, fmt.Errorf("failed to get roles on policy: %w", err)
	}
	return es, nil
}

// GetPoliciesOnRole ロールに関連付けられたポリシーを取得する。
func (m *ManageRoleAssociation) GetPoliciesOnRole(
	ctx context.Context, roleID uuid.UUID,
	whereSearchName string,
	order parameter.PolicyOnRoleOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.PolicyOnRole], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePolicyOnRoleParam{
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
	es, err := m.DB.GetPoliciesOnRole(ctx, roleID, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyOnRole]{}, fmt.Errorf("failed to get policies on role: %w", err)
	}
	return es, nil
}

// GetRolesOnPolicyCount ポリシーに関連付けられたロールの数を取得する。
func (m *ManageRoleAssociation) GetRolesOnPolicyCount(
	ctx context.Context, policyID uuid.UUID,
	whereSearchName string,
) (int64, error) {
	where := parameter.WhereRoleOnPolicyParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err := m.DB.CountRolesOnPolicy(ctx, policyID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get roles on policy count: %w", err)
	}
	return es, nil
}

// GetPoliciesOnRoleCount ロールに関連付けられたポリシーの数を取得する。
func (m *ManageRoleAssociation) GetPoliciesOnRoleCount(
	ctx context.Context, roleID uuid.UUID,
	whereSearchName string,
) (int64, error) {
	where := parameter.WherePolicyOnRoleParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	es, err := m.DB.CountPoliciesOnRole(ctx, roleID, where)
	if err != nil {
		return 0, fmt.Errorf("failed to get policies on role count: %w", err)
	}
	return es, nil
}
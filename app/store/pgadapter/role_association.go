package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func countPoliciesOnRole(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam,
) (int64, error) {
	p := query.CountPoliciesOnRoleParams{
		RoleID:        roleID,
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountPoliciesOnRole(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count policies on role: %w", err)
	}
	return c, nil
}

// CountPoliciesOnRole ロールに関連付けられたポリシー数を取得する。
func (a *PgAdapter) CountPoliciesOnRole(
	ctx context.Context, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam,
) (int64, error) {
	return countPoliciesOnRole(ctx, a.query, roleID, where)
}

// CountPoliciesOnRoleWithSd SD付きでロールに関連付けられたポリシー数を取得する。
func (a *PgAdapter) CountPoliciesOnRoleWithSd(
	ctx context.Context, sd store.Sd, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countPoliciesOnRole(ctx, qtx, roleID, where)
}

func countRolesOnPolicy(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam,
) (int64, error) {
	p := query.CountRolesOnPolicyParams{
		PolicyID:      policyID,
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountRolesOnPolicy(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count role on policy: %w", err)
	}
	return c, nil
}

// CountRolesOnPolicy ポリシーに関連付けられたロール数を取得する。
func (a *PgAdapter) CountRolesOnPolicy(
	ctx context.Context, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam,
) (int64, error) {
	return countRolesOnPolicy(ctx, a.query, policyID, where)
}

// CountRolesOnPolicyWithSd SD付きでポリシーに関連付けられたロール数を取得する。
func (a *PgAdapter) CountRolesOnPolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countRolesOnPolicy(ctx, qtx, policyID, where)
}

func associateRole(
	ctx context.Context, qtx *query.Queries, param parameter.AssociationRoleParam,
) (entity.RoleAssociation, error) {
	p := query.CreateRoleAssociationParams{
		RoleID:   param.RoleID,
		PolicyID: param.PolicyID,
	}
	e, err := qtx.CreateRoleAssociation(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.RoleAssociation{}, errhandle.NewModelDuplicatedError("role association")
		}
		return entity.RoleAssociation{}, fmt.Errorf("failed to associate role: %w", err)
	}
	entity := entity.RoleAssociation{
		RoleID:   e.RoleID,
		PolicyID: e.PolicyID,
	}
	return entity, nil
}

// AssociateRole ロールを関連付ける。
func (a *PgAdapter) AssociateRole(
	ctx context.Context, param parameter.AssociationRoleParam,
) (entity.RoleAssociation, error) {
	return associateRole(ctx, a.query, param)
}

// AssociateRoleWithSd SD付きでロールを関連付ける。
func (a *PgAdapter) AssociateRoleWithSd(
	ctx context.Context, sd store.Sd, param parameter.AssociationRoleParam,
) (entity.RoleAssociation, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.RoleAssociation{}, store.ErrNotFoundDescriptor
	}
	return associateRole(ctx, qtx, param)
}

func associateRoles(
	ctx context.Context, qtx *query.Queries, params []parameter.AssociationRoleParam,
) (int64, error) {
	p := make([]query.CreateRoleAssociationsParams, len(params))
	for i, param := range params {
		p[i] = query.CreateRoleAssociationsParams{
			RoleID:   param.RoleID,
			PolicyID: param.PolicyID,
		}
	}
	c, err := qtx.CreateRoleAssociations(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("role association")
		}
		return 0, fmt.Errorf("failed to associate roles: %w", err)
	}
	return c, nil
}

// AssociateRoles ロールを複数関連付ける。
func (a *PgAdapter) AssociateRoles(ctx context.Context, params []parameter.AssociationRoleParam) (int64, error) {
	return associateRoles(ctx, a.query, params)
}

// AssociateRolesWithSd SD付きでロールを複数関連付ける。
func (a *PgAdapter) AssociateRolesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.AssociationRoleParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return associateRoles(ctx, qtx, params)
}

func disassociateRole(
	ctx context.Context, qtx *query.Queries, roleID, policyID uuid.UUID,
) (int64, error) {
	p := query.DeleteRoleAssociationParams{
		RoleID:   roleID,
		PolicyID: policyID,
	}
	c, err := qtx.DeleteRoleAssociation(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("role association")
	}
	return c, nil
}

// DisassociateRole ロールの関連付けを解除する。
func (a *PgAdapter) DisassociateRole(
	ctx context.Context, roleID, policyID uuid.UUID,
) (int64, error) {
	return disassociateRole(ctx, a.query, roleID, policyID)
}

// DisassociateRoleWithSd SD付きでロールの関連付けを解除する。
func (a *PgAdapter) DisassociateRoleWithSd(
	ctx context.Context, sd store.Sd, roleID, policyID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disassociateRole(ctx, qtx, roleID, policyID)
}

func disassociateRoleOnPolicy(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteRoleAssociationsOnPolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on policy: %w", err)
	}
	return c, nil
}

// DisassociateRoleOnPolicy ポリシーに関連付けられたロールを解除する。
func (a *PgAdapter) DisassociateRoleOnPolicy(
	ctx context.Context, policyID uuid.UUID,
) (int64, error) {
	return disassociateRoleOnPolicy(ctx, a.query, policyID)
}

// DisassociateRoleOnPolicyWithSd SD付きでポリシーに関連付けられたロールを解除する。
func (a *PgAdapter) DisassociateRoleOnPolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disassociateRoleOnPolicy(ctx, qtx, policyID)
}

func disassociateRoleOnPolicies(
	ctx context.Context, qtx *query.Queries, policyIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteRoleAssociationsOnPolicies(ctx, policyIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on policies: %w", err)
	}
	return c, nil
}

// DisassociateRoleOnPolicies ポリシーに関連付けられたロールを複数解除する。
func (a *PgAdapter) DisassociateRoleOnPolicies(
	ctx context.Context, policyIDs []uuid.UUID,
) (int64, error) {
	return disassociateRoleOnPolicies(ctx, a.query, policyIDs)
}

// DisassociateRoleOnPoliciesWithSd SD付きでポリシーに関連付けられたロールを複数解除する。
func (a *PgAdapter) DisassociateRoleOnPoliciesWithSd(
	ctx context.Context, sd store.Sd, policyIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disassociateRoleOnPolicies(ctx, qtx, policyIDs)
}

func pluralDisassociateRoleOnPolicy(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID, roleIDs []uuid.UUID,
) (int64, error) {
	p := query.PluralDeleteRoleAssociationsOnPolicyParams{
		PolicyID: policyID,
		RoleIds:  roleIDs,
	}
	c, err := qtx.PluralDeleteRoleAssociationsOnPolicy(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on policy: %w", err)
	}
	if c != int64(len(roleIDs)) {
		return 0, errhandle.NewModelNotFoundError("role association")
	}
	return c, nil
}

// PluralDisassociateRoleOnPolicy ポリシーに関連付けられたロールを複数解除する。
func (a *PgAdapter) PluralDisassociateRoleOnPolicy(
	ctx context.Context, policyID uuid.UUID, roleIDs []uuid.UUID,
) (int64, error) {
	return pluralDisassociateRoleOnPolicy(ctx, a.query, policyID, roleIDs)
}

// PluralDisassociateRoleOnPolicyWithSd SD付きでポリシーに関連付けられたロールを複数解除する。
func (a *PgAdapter) PluralDisassociateRoleOnPolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID, roleIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDisassociateRoleOnPolicy(ctx, qtx, policyID, roleIDs)
}

func disassociatePolicyOnRole(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteRoleAssociationsOnRole(ctx, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on role: %w", err)
	}
	return c, nil
}

// DisassociatePolicyOnRole ロールに関連付けられたポリシーを解除する。
func (a *PgAdapter) DisassociatePolicyOnRole(ctx context.Context, roleID uuid.UUID) (int64, error) {
	return disassociatePolicyOnRole(ctx, a.query, roleID)
}

// DisassociatePolicyOnRoleWithSd SD付きでロールに関連付けられたポリシーを解除する。
func (a *PgAdapter) DisassociatePolicyOnRoleWithSd(ctx context.Context, sd store.Sd, roleID uuid.UUID) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disassociatePolicyOnRole(ctx, qtx, roleID)
}

func disassociatePolicyOnRoles(
	ctx context.Context, qtx *query.Queries, roleIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteRoleAssociationsOnRoles(ctx, roleIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on roles: %w", err)
	}
	return c, nil
}

// DisassociatePolicyOnRoles ロールに関連付けられたポリシーを複数解除する。
func (a *PgAdapter) DisassociatePolicyOnRoles(ctx context.Context, roleIDs []uuid.UUID) (int64, error) {
	return disassociatePolicyOnRoles(ctx, a.query, roleIDs)
}

// DisassociatePolicyOnRolesWithSd SD付きでロールに関連付けられたポリシーを複数解除する。
func (a *PgAdapter) DisassociatePolicyOnRolesWithSd(
	ctx context.Context, sd store.Sd, roleIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disassociatePolicyOnRoles(ctx, qtx, roleIDs)
}

func pluralDisassociatePolicyOnRole(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID, policyIDs []uuid.UUID,
) (int64, error) {
	p := query.PluralDeleteRoleAssociationsOnRoleParams{
		RoleID:    roleID,
		PolicyIds: policyIDs,
	}
	c, err := qtx.PluralDeleteRoleAssociationsOnRole(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to disassociate role on role: %w", err)
	}
	if c != int64(len(policyIDs)) {
		return 0, errhandle.NewModelNotFoundError("role association")
	}
	return c, nil
}

// PluralDisassociatePolicyOnRole ロールに関連付けられたポリシーを複数解除する。
func (a *PgAdapter) PluralDisassociatePolicyOnRole(
	ctx context.Context, roleID uuid.UUID, policyIDs []uuid.UUID,
) (int64, error) {
	return pluralDisassociatePolicyOnRole(ctx, a.query, roleID, policyIDs)
}

// PluralDisassociatePolicyOnRoleWithSd SD付きでロールに関連付けられたポリシーを複数解除する。
func (a *PgAdapter) PluralDisassociatePolicyOnRoleWithSd(
	ctx context.Context, sd store.Sd, roleID uuid.UUID, policyIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDisassociatePolicyOnRole(ctx, qtx, roleID, policyIDs)
}

func getRolesOnPolicy(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID, where parameter.WhereRoleOnPolicyParam,
	order parameter.RoleOnPolicyOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.RoleOnPolicy], error) {
	eConvFunc := func(e entity.RoleOnPolicyForQuery) (entity.RoleOnPolicy, error) {
		return entity.RoleOnPolicy{
			Role: e.Role,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountRolesOnPolicyParams{
			PolicyID:      policyID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountRolesOnPolicy(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count roles on policy: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.RoleOnPolicyForQuery, error) {
		p := query.GetRolesOnPolicyParams{
			PolicyID:      policyID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetRolesOnPolicy(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.RoleOnPolicyForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get roles on policy: %w", err)
		}
		fq := make([]entity.RoleOnPolicyForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.RoleOnPolicyForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				RoleOnPolicy: entity.RoleOnPolicy{
					Role: entity.Role{
						RoleID:      e.RoleID,
						Name:        e.RoleName.String,
						Description: e.RoleDescription.String,
					},
				},
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.RoleOnPolicyForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.RoleOnPolicyNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetRolesOnPolicyUseKeysetPaginateParams{
			PolicyID:        policyID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetRolesOnPolicyUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles on policy: %w", err)
		}
		fq := make([]entity.RoleOnPolicyForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.RoleOnPolicyForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				RoleOnPolicy: entity.RoleOnPolicy{
					Role: entity.Role{
						RoleID:      e.RoleID,
						Name:        e.RoleName.String,
						Description: e.RoleDescription.String,
					},
				},
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.RoleOnPolicyForQuery, error) {
		p := query.GetRolesOnPolicyUseNumberedPaginateParams{
			PolicyID:      policyID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetRolesOnPolicyUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles on policy: %w", err)
		}
		fq := make([]entity.RoleOnPolicyForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.RoleOnPolicyForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				RoleOnPolicy: entity.RoleOnPolicy{
					Role: entity.Role{
						RoleID:      e.RoleID,
						Name:        e.RoleName.String,
						Description: e.RoleDescription.String,
					},
				},
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.RoleOnPolicyForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.RoleOnPolicyDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.RoleOnPolicyNameCursorKey:
			return entity.Int(e.Pkey), e.Role.Name
		}
		return entity.Int(e.Pkey), nil
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
		return store.ListResult[entity.RoleOnPolicy]{}, fmt.Errorf("failed to get roles on policy: %w", err)
	}
	return res, nil
}

// GetRolesOnPolicy ポリシーに関連付けられたロールを取得する。
func (a *PgAdapter) GetRolesOnPolicy(
	ctx context.Context,
	policyID uuid.UUID,
	where parameter.WhereRoleOnPolicyParam,
	order parameter.RoleOnPolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.RoleOnPolicy], error) {
	return getRolesOnPolicy(ctx, a.query, policyID, where, order, np, cp, wc)
}

// GetRolesOnPolicyWithSd SD付きでポリシーに関連付けられたロールを取得する。
func (a *PgAdapter) GetRolesOnPolicyWithSd(
	ctx context.Context,
	sd store.Sd,
	policyID uuid.UUID,
	where parameter.WhereRoleOnPolicyParam,
	order parameter.RoleOnPolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.RoleOnPolicy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.RoleOnPolicy]{}, store.ErrNotFoundDescriptor
	}
	return getRolesOnPolicy(ctx, qtx, policyID, where, order, np, cp, wc)
}

func getPoliciesOnRole(
	ctx context.Context, qtx *query.Queries, roleID uuid.UUID, where parameter.WherePolicyOnRoleParam,
	order parameter.PolicyOnRoleOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.PolicyOnRole], error) {
	eConvFunc := func(e entity.PolicyOnRoleForQuery) (entity.PolicyOnRole, error) {
		return entity.PolicyOnRole{
			Policy: e.Policy,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPoliciesOnRoleParams{
			RoleID:        roleID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountPoliciesOnRole(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count policies on role: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.PolicyOnRoleForQuery, error) {
		p := query.GetPoliciesOnRoleParams{
			RoleID:        roleID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetPoliciesOnRole(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.PolicyOnRoleForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get policies on role: %w", err)
		}
		fq := make([]entity.PolicyOnRoleForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.PolicyOnRoleForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				PolicyOnRole: entity.PolicyOnRole{
					Policy: entity.Policy{
						PolicyID:         e.PolicyID,
						Name:             e.PolicyName.String,
						Description:      e.PolicyDescription.String,
						Key:              e.PolicyKey.String,
						PolicyCategoryID: e.PolicyCategoryID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.PolicyOnRoleForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PolicyOnRoleNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPoliciesOnRoleUseKeysetPaginateParams{
			RoleID:          roleID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPoliciesOnRoleUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policies on role: %w", err)
		}
		fq := make([]entity.PolicyOnRoleForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.PolicyOnRoleForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				PolicyOnRole: entity.PolicyOnRole{
					Policy: entity.Policy{
						PolicyID:         e.PolicyID,
						Name:             e.PolicyName.String,
						Description:      e.PolicyDescription.String,
						Key:              e.PolicyKey.String,
						PolicyCategoryID: e.PolicyCategoryID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.PolicyOnRoleForQuery, error) {
		p := query.GetPoliciesOnRoleUseNumberedPaginateParams{
			RoleID:        roleID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetPoliciesOnRoleUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policies on role: %w", err)
		}
		fq := make([]entity.PolicyOnRoleForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.PolicyOnRoleForQuery{
				Pkey: entity.Int(e.MRoleAssociationsPkey),
				PolicyOnRole: entity.PolicyOnRole{
					Policy: entity.Policy{
						PolicyID:         e.PolicyID,
						Name:             e.PolicyName.String,
						Description:      e.PolicyDescription.String,
						Key:              e.PolicyKey.String,
						PolicyCategoryID: e.PolicyCategoryID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.PolicyOnRoleForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.PolicyOnRoleDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.PolicyOnRoleNameCursorKey:
			return entity.Int(e.Pkey), e.Policy.Name
		}
		return entity.Int(e.Pkey), nil
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
		return store.ListResult[entity.PolicyOnRole]{}, fmt.Errorf("failed to get policies on role: %w", err)
	}
	return res, nil
}

// GetPoliciesOnRole ロールに関連付けられたポリシーを取得する。
func (a *PgAdapter) GetPoliciesOnRole(
	ctx context.Context,
	roleID uuid.UUID,
	where parameter.WherePolicyOnRoleParam,
	order parameter.PolicyOnRoleOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyOnRole], error) {
	return getPoliciesOnRole(ctx, a.query, roleID, where, order, np, cp, wc)
}

// GetPoliciesOnRoleWithSd SD付きでロールに関連付けられたポリシーを取得する。
func (a *PgAdapter) GetPoliciesOnRoleWithSd(
	ctx context.Context,
	sd store.Sd,
	roleID uuid.UUID,
	where parameter.WherePolicyOnRoleParam,
	order parameter.PolicyOnRoleOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyOnRole], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PolicyOnRole]{}, store.ErrNotFoundDescriptor
	}
	return getPoliciesOnRole(ctx, qtx, roleID, where, order, np, cp, wc)
}

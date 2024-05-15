package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
)

func countPolicies(
	ctx context.Context, qtx *query.Queries, where parameter.WherePolicyParam,
) (int64, error) {
	p := query.CountPoliciesParams{
		WhereLikeName:   where.WhereLikeName,
		SearchName:      where.SearchName,
		WhereInCategory: where.WhereInCategory,
		InCategories:    where.InCategories,
	}
	c, err := qtx.CountPolicies(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy: %w", err)
	}
	return c, nil
}

// CountPolicies ポリシーカテゴリー数を取得する。
func (a *PgAdapter) CountPolicies(
	ctx context.Context, where parameter.WherePolicyParam,
) (int64, error) {
	return countPolicies(ctx, a.query, where)
}

// CountPoliciesWithSd SD付きでポリシーカテゴリー数を取得する。
func (a *PgAdapter) CountPoliciesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WherePolicyParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countPolicies(ctx, qtx, where)
}

func createPolicy(
	ctx context.Context, qtx *query.Queries, param parameter.CreatePolicyParam,
) (entity.Policy, error) {
	p := query.CreatePolicyParams{
		Name:             param.Name,
		Key:              param.Key,
		Description:      param.Description,
		PolicyCategoryID: param.PolicyCategoryID,
	}
	e, err := qtx.CreatePolicy(ctx, p)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to create policy: %w", err)
	}
	entity := entity.Policy{
		PolicyID:         e.PolicyID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
		PolicyCategoryID: e.PolicyCategoryID,
	}
	return entity, nil
}

// CreatePolicy ポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicy(
	ctx context.Context, param parameter.CreatePolicyParam,
) (entity.Policy, error) {
	return createPolicy(ctx, a.query, param)
}

// CreatePolicyWithSd SD付きでポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreatePolicyParam,
) (entity.Policy, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Policy{}, store.ErrNotFoundDescriptor
	}
	return createPolicy(ctx, qtx, param)
}

func createPolicies(
	ctx context.Context, qtx *query.Queries, params []parameter.CreatePolicyParam,
) (int64, error) {
	p := make([]query.CreatePoliciesParams, len(params))
	for i, param := range params {
		p[i] = query.CreatePoliciesParams{
			Name:             param.Name,
			Key:              param.Key,
			Description:      param.Description,
			PolicyCategoryID: param.PolicyCategoryID,
		}
	}
	c, err := qtx.CreatePolicies(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy: %w", err)
	}
	return c, nil
}

// CreatePolicies ポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicies(
	ctx context.Context, params []parameter.CreatePolicyParam,
) (int64, error) {
	return createPolicies(ctx, a.query, params)
}

// CreatePoliciesWithSd SD付きでポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePoliciesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreatePolicyParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createPolicies(ctx, qtx, params)
}

func deletePolicy(ctx context.Context, qtx *query.Queries, policyID uuid.UUID) (int64, error) {
	c, err := qtx.DeletePolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	return c, nil
}

// DeletePolicy ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicy(ctx context.Context, policyID uuid.UUID) (int64, error) {
	return deletePolicy(ctx, a.query, policyID)
}

// DeletePolicyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePolicy(ctx, qtx, policyID)
}

func deletePolicyByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeletePolicyByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	return c, nil
}

// DeletePolicyByKey ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyByKey(ctx context.Context, key string) (int64, error) {
	return deletePolicyByKey(ctx, a.query, key)
}

// DeletePolicyByKeyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePolicyByKey(ctx, qtx, key)
}

func pluralDeletePolicies(
	ctx context.Context, qtx *query.Queries, policyIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeletePolicies(ctx, policyIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return c, nil
}

// PluralDeletePolicies ポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicies(
	ctx context.Context, policyIDs []uuid.UUID,
) (int64, error) {
	return pluralDeletePolicies(ctx, a.query, policyIDs)
}

// PluralDeletePoliciesWithSd SD付きでポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePoliciesWithSd(
	ctx context.Context, sd store.Sd, policyIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeletePolicies(ctx, qtx, policyIDs)
}

func findPolicyByID(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID,
) (entity.Policy, error) {
	e, err := qtx.FindPolicyByID(ctx, policyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	entity := entity.Policy{
		PolicyID:         e.PolicyID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
		PolicyCategoryID: e.PolicyCategoryID,
	}
	return entity, nil
}

// FindPolicyByID ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByID(
	ctx context.Context, policyID uuid.UUID,
) (entity.Policy, error) {
	return findPolicyByID(ctx, a.query, policyID)
}

// FindPolicyByIDWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByIDWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID,
) (entity.Policy, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Policy{}, store.ErrNotFoundDescriptor
	}
	return findPolicyByID(ctx, qtx, policyID)
}

func findPolicyByIDWithCategory(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID,
) (entity.PolicyWithCategory, error) {
	e, err := qtx.FindPolicyByIDWithCategory(ctx, policyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyWithCategory{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	entity := entity.PolicyWithCategory{
		Policy: entity.Policy{
			PolicyID:         e.PolicyID,
			Name:             e.Name,
			Key:              e.Key,
			Description:      e.Description,
			PolicyCategoryID: e.PolicyCategoryID,
		},
		PolicyCategory: entity.PolicyCategory{
			PolicyCategoryID: e.PolicyCategoryID,
			Name:             e.PolicyCategoryName,
			Key:              e.PolicyCategoryKey,
			Description:      e.PolicyCategoryDescription,
		},
	}
	return entity, nil
}

// FindPolicyByIDWithCategory ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByIDWithCategory(
	ctx context.Context, policyID uuid.UUID,
) (entity.PolicyWithCategory, error) {
	return findPolicyByIDWithCategory(ctx, a.query, policyID)
}

// FindPolicyByIDWithCategoryWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByIDWithCategoryWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID,
) (entity.PolicyWithCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyWithCategory{}, store.ErrNotFoundDescriptor
	}
	return findPolicyByIDWithCategory(ctx, qtx, policyID)
}

func findPolicyByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.Policy, error) {
	e, err := qtx.FindPolicyByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	entity := entity.Policy{
		PolicyID:         e.PolicyID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
		PolicyCategoryID: e.PolicyCategoryID,
	}
	return entity, nil
}

// FindPolicyByKey ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByKey(ctx context.Context, key string) (entity.Policy, error) {
	return findPolicyByKey(ctx, a.query, key)
}

// FindPolicyByKeyWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.Policy, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Policy{}, store.ErrNotFoundDescriptor
	}
	return findPolicyByKey(ctx, qtx, key)
}

func findPolicyByKeyWithCategory(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PolicyWithCategory, error) {
	e, err := qtx.FindPolicyByKeyWithCategory(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyWithCategory{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	entity := entity.PolicyWithCategory{
		Policy: entity.Policy{
			PolicyID:         e.PolicyID,
			Name:             e.Name,
			Key:              e.Key,
			Description:      e.Description,
			PolicyCategoryID: e.PolicyCategoryID,
		},
		PolicyCategory: entity.PolicyCategory{
			PolicyCategoryID: e.PolicyCategoryID,
			Name:             e.PolicyCategoryName,
			Key:              e.PolicyCategoryKey,
			Description:      e.PolicyCategoryDescription,
		},
	}
	return entity, nil
}

// FindPolicyByKeyWithCategory ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByKeyWithCategory(
	ctx context.Context, key string,
) (entity.PolicyWithCategory, error) {
	return findPolicyByKeyWithCategory(ctx, a.query, key)
}

// FindPolicyByKeyWithCategoryWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyByKeyWithCategoryWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.PolicyWithCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyWithCategory{}, store.ErrNotFoundDescriptor
	}
	return findPolicyByKeyWithCategory(ctx, qtx, key)
}

// PolicyCursor is a cursor for Policy.
type PolicyCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getPolicies(
	ctx context.Context, qtx *query.Queries, where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Policy], error) {
	eConvFunc := func(e query.Policy) (entity.Policy, error) {
		return entity.Policy{
			PolicyID:         e.PolicyID,
			Name:             e.Name,
			Key:              e.Key,
			Description:      e.Description,
			PolicyCategoryID: e.PolicyCategoryID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPoliciesParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
		}
		r, err := qtx.CountPolicies(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count policy: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Policy, error) {
		p := query.GetPoliciesParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
		}
		r, err := qtx.GetPolicies(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Policy{}, nil
			}
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Policy, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PolicyNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPoliciesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPoliciesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Policy, error) {
		p := query.GetPoliciesUseNumberedPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Offset:          offset,
			Limit:           limit,
		}
		r, err := qtx.GetPoliciesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Policy) (entity.Int, any) {
		switch subCursor {
		case parameter.PolicyDefaultCursorKey:
			return entity.Int(e.MPoliciesPkey), nil
		case parameter.PolicyNameCursorKey:
			return entity.Int(e.MPoliciesPkey), e.Name
		}
		return entity.Int(e.MPoliciesPkey), nil
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
		return store.ListResult[entity.Policy]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return res, nil
}

// GetPolicies ポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPolicies(
	ctx context.Context,
	where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Policy], error) {
	return getPolicies(ctx, a.query, where, order, np, cp, wc)
}

// GetPoliciesWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPoliciesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Policy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Policy]{}, store.ErrNotFoundDescriptor
	}
	return getPolicies(ctx, qtx, where, order, np, cp, wc)
}

func getPoliciesWithCategory(
	ctx context.Context, qtx *query.Queries, where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyWithCategory], error) {
	eConvFunc := func(e entity.PolicyWithCategoryForQuery) (entity.PolicyWithCategory, error) {
		return e.PolicyWithCategory, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPoliciesParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
		}
		r, err := qtx.CountPolicies(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count policy: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.PolicyWithCategoryForQuery, error) {
		p := query.GetPoliciesWithCategoryParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
		}
		r, err := qtx.GetPoliciesWithCategory(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.PolicyWithCategoryForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		e := make([]entity.PolicyWithCategoryForQuery, len(r))
		for i, v := range r {
			e[i] = entity.PolicyWithCategoryForQuery{
				Pkey: entity.Int{
					Int64: v.MPoliciesPkey.Int64,
					Valid: v.MPoliciesPkey.Valid,
				},
				PolicyWithCategory: entity.PolicyWithCategory{
					Policy: entity.Policy{
						PolicyID:         v.PolicyID,
						Name:             v.Name,
						Key:              v.Key,
						Description:      v.Description,
						PolicyCategoryID: v.PolicyCategoryID,
					},
					PolicyCategory: entity.PolicyCategory{
						PolicyCategoryID: v.PolicyCategoryID,
						Name:             v.PolicyCategoryName,
						Key:              v.PolicyCategoryKey,
						Description:      v.PolicyCategoryDescription,
					},
				},
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.PolicyWithCategoryForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PolicyNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPoliciesWithCategoryUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPoliciesWithCategoryUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		e := make([]entity.PolicyWithCategoryForQuery, len(r))
		for i, v := range r {
			e[i] = entity.PolicyWithCategoryForQuery{
				Pkey: entity.Int{
					Int64: v.MPoliciesPkey.Int64,
					Valid: v.MPoliciesPkey.Valid,
				},
				PolicyWithCategory: entity.PolicyWithCategory{
					Policy: entity.Policy{
						PolicyID:         v.PolicyID,
						Name:             v.Name,
						Key:              v.Key,
						Description:      v.Description,
						PolicyCategoryID: v.PolicyCategoryID,
					},
					PolicyCategory: entity.PolicyCategory{
						PolicyCategoryID: v.PolicyCategoryID,
						Name:             v.PolicyCategoryName,
						Key:              v.PolicyCategoryKey,
						Description:      v.PolicyCategoryDescription,
					},
				},
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.PolicyWithCategoryForQuery, error) {
		p := query.GetPoliciesWithCategoryUseNumberedPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			WhereInCategory: where.WhereInCategory,
			InCategories:    where.InCategories,
			OrderMethod:     orderMethod,
			Offset:          offset,
			Limit:           limit,
		}
		r, err := qtx.GetPoliciesWithCategoryUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy: %w", err)
		}
		e := make([]entity.PolicyWithCategoryForQuery, len(r))
		for i, v := range r {
			e[i] = entity.PolicyWithCategoryForQuery{
				Pkey: entity.Int{
					Int64: v.MPoliciesPkey.Int64,
					Valid: v.MPoliciesPkey.Valid,
				},
				PolicyWithCategory: entity.PolicyWithCategory{
					Policy: entity.Policy{
						PolicyID:         v.PolicyID,
						Name:             v.Name,
						Key:              v.Key,
						Description:      v.Description,
						PolicyCategoryID: v.PolicyCategoryID,
					},
					PolicyCategory: entity.PolicyCategory{
						PolicyCategoryID: v.PolicyCategoryID,
						Name:             v.PolicyCategoryName,
						Key:              v.PolicyCategoryKey,
						Description:      v.PolicyCategoryDescription,
					},
				},
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.PolicyWithCategoryForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.PolicyDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.PolicyNameCursorKey:
			return entity.Int(e.Pkey), e.Name
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
		return store.ListResult[entity.PolicyWithCategory]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return res, nil
}

// GetPoliciesWithCategory ポリシーとそのカテゴリーを取得する。
func (a *PgAdapter) GetPoliciesWithCategory(
	ctx context.Context,
	where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyWithCategory], error) {
	return getPoliciesWithCategory(ctx, a.query, where, order, np, cp, wc)
}

// GetPoliciesWithCategoryWithSd SD付きでポリシーとそのカテゴリーを取得する。
func (a *PgAdapter) GetPoliciesWithCategoryWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePolicyParam,
	order parameter.PolicyOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyWithCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PolicyWithCategory]{}, store.ErrNotFoundDescriptor
	}
	return getPoliciesWithCategory(ctx, qtx, where, order, np, cp, wc)
}

func getPluralPolicies(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Policy], error) {
	var e []query.Policy
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralPolicies(ctx, ids)
	} else {
		e, err = qtx.GetPluralPoliciesUseNumberedPaginate(ctx, query.GetPluralPoliciesUseNumberedPaginateParams{
			PolicyIds: ids,
			Offset:    int32(np.Offset.Int64),
			Limit:     int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Policy]{},
			fmt.Errorf("failed to get plural policy: %w", err)
	}
	entities := make([]entity.Policy, len(e))
	for i, v := range e {
		entities[i] = entity.Policy{
			PolicyID:         v.PolicyID,
			Name:             v.Name,
			Key:              v.Key,
			Description:      v.Description,
			PolicyCategoryID: v.PolicyCategoryID,
		}
	}
	return store.ListResult[entity.Policy]{Data: entities}, nil
}

// GetPluralPolicies ポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPluralPolicies(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Policy], error) {
	return getPluralPolicies(ctx, a.query, ids, np)
}

// GetPluralPoliciesWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPluralPoliciesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Policy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Policy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralPolicies(ctx, qtx, ids, np)
}

func updatePolicy(
	ctx context.Context, qtx *query.Queries,
	policyID uuid.UUID, param parameter.UpdatePolicyParams,
) (entity.Policy, error) {
	p := query.UpdatePolicyParams{
		PolicyID:         policyID,
		Name:             param.Name,
		Key:              param.Key,
		Description:      param.Description,
		PolicyCategoryID: param.PolicyCategoryID,
	}
	e, err := qtx.UpdatePolicy(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	entity := entity.Policy{
		PolicyID:         e.PolicyID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
		PolicyCategoryID: e.PolicyCategoryID,
	}
	return entity, nil
}

// UpdatePolicy ポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicy(
	ctx context.Context, policyID uuid.UUID, param parameter.UpdatePolicyParams,
) (entity.Policy, error) {
	return updatePolicy(ctx, a.query, policyID, param)
}

// UpdatePolicyWithSd SD付きでポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID, param parameter.UpdatePolicyParams,
) (entity.Policy, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Policy{}, store.ErrNotFoundDescriptor
	}
	return updatePolicy(ctx, qtx, policyID, param)
}

func updatePolicyByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdatePolicyByKeyParams,
) (entity.Policy, error) {
	p := query.UpdatePolicyByKeyParams{
		Key:              key,
		Name:             param.Name,
		Description:      param.Description,
		PolicyCategoryID: param.PolicyCategoryID,
	}
	e, err := qtx.UpdatePolicyByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, errhandle.NewModelNotFoundError("policy")
		}
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	entity := entity.Policy{
		PolicyID:         e.PolicyID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
		PolicyCategoryID: e.PolicyCategoryID,
	}
	return entity, nil
}

// UpdatePolicyByKey ポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyByKey(
	ctx context.Context, key string, param parameter.UpdatePolicyByKeyParams,
) (entity.Policy, error) {
	return updatePolicyByKey(ctx, a.query, key, param)
}

// UpdatePolicyByKeyWithSd SD付きでポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdatePolicyByKeyParams,
) (entity.Policy, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Policy{}, store.ErrNotFoundDescriptor
	}
	return updatePolicyByKey(ctx, qtx, key, param)
}

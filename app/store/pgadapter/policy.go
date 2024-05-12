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
	c, err := countPolicies(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy: %w", err)
	}
	return c, nil
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
	c, err := countPolicies(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy: %w", err)
	}
	return c, nil
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
	e, err := createPolicy(ctx, a.query, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to create policy: %w", err)
	}
	return e, nil
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
	e, err := createPolicy(ctx, qtx, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to create policy: %w", err)
	}
	return e, nil
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
	c, err := createPolicies(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy: %w", err)
	}
	return c, nil
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
	c, err := createPolicies(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy: %w", err)
	}
	return c, nil
}

func deletePolicy(ctx context.Context, qtx *query.Queries, policyID uuid.UUID) error {
	err := qtx.DeletePolicy(ctx, policyID)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

// DeletePolicy ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicy(ctx context.Context, policyID uuid.UUID) error {
	err := deletePolicy(ctx, a.query, policyID)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

// DeletePolicyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyWithSd(
	ctx context.Context, sd store.Sd, policyID uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePolicy(ctx, qtx, policyID)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

func deletePolicyByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeletePolicyByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

// DeletePolicyByKey ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyByKey(ctx context.Context, key string) error {
	err := deletePolicyByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

// DeletePolicyByKeyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePolicyByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}

func pluralDeletePolicies(
	ctx context.Context, qtx *query.Queries, policyIDs []uuid.UUID,
) error {
	err := qtx.PluralDeletePolicies(ctx, policyIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return nil
}

// PluralDeletePolicies ポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicies(
	ctx context.Context, policyIDs []uuid.UUID,
) error {
	err := pluralDeletePolicies(ctx, a.query, policyIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return nil
}

// PluralDeletePoliciesWithSd SD付きでポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePoliciesWithSd(
	ctx context.Context, sd store.Sd, policyIDs []uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeletePolicies(ctx, qtx, policyIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return nil
}

func findPolicyByID(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID,
) (entity.Policy, error) {
	e, err := qtx.FindPolicyByID(ctx, policyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, store.ErrDataNoRecord
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
	e, err := findPolicyByID(ctx, a.query, policyID)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
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
	e, err := findPolicyByID(ctx, qtx, policyID)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
}

func findPolicyByIDWithCategory(
	ctx context.Context, qtx *query.Queries, policyID uuid.UUID,
) (entity.PolicyWithCategory, error) {
	e, err := qtx.FindPolicyByIDWithCategory(ctx, policyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyWithCategory{}, store.ErrDataNoRecord
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
	e, err := findPolicyByIDWithCategory(ctx, a.query, policyID)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
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
	e, err := findPolicyByIDWithCategory(ctx, qtx, policyID)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
}

func findPolicyByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.Policy, error) {
	e, err := qtx.FindPolicyByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Policy{}, store.ErrDataNoRecord
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
	e, err := findPolicyByKey(ctx, a.query, key)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
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
	e, err := findPolicyByKey(ctx, qtx, key)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
}

func findPolicyByKeyWithCategory(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PolicyWithCategory, error) {
	e, err := qtx.FindPolicyByKeyWithCategory(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyWithCategory{}, store.ErrDataNoRecord
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
	e, err := findPolicyByKeyWithCategory(ctx, a.query, key)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
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
	e, err := findPolicyByKeyWithCategory(ctx, qtx, key)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy: %w", err)
	}
	return e, nil
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
	r, err := getPolicies(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Policy]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return r, nil
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
	r, err := getPolicies(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Policy]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return r, nil
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
	r, err := getPoliciesWithCategory(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyWithCategory]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return r, nil
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
	r, err := getPoliciesWithCategory(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyWithCategory]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return r, nil
}

func getPluralPolicies(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Policy], error) {
	p := query.GetPluralPoliciesParams{
		PolicyIds: ids,
		Offset:    int32(np.Offset.Int64),
		Limit:     int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralPolicies(ctx, p)
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
	r, err := getPluralPolicies(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.Policy]{},
			fmt.Errorf("failed to get plural policy: %w", err)
	}
	return r, nil
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
	r, err := getPluralPolicies(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.Policy]{},
			fmt.Errorf("failed to get plural policy: %w", err)
	}
	return r, nil
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
			return entity.Policy{}, store.ErrDataNoRecord
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
	e, err := updatePolicy(ctx, a.query, policyID, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	return e, nil
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
	e, err := updatePolicy(ctx, qtx, policyID, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	return e, nil
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
			return entity.Policy{}, store.ErrDataNoRecord
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
	e, err := updatePolicyByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	return e, nil
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
	e, err := updatePolicyByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	return e, nil
}

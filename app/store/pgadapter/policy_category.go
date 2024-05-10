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

func countPolicyCategories(
	ctx context.Context, qtx *query.Queries, where parameter.WherePolicyCategoryParam,
) (int64, error) {
	p := query.CountPolicyCategoriesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountPolicyCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy category: %w", err)
	}
	return c, nil
}

// CountPolicyCategories ポリシーカテゴリー数を取得する。
func (a *PgAdapter) CountPolicyCategories(
	ctx context.Context, where parameter.WherePolicyCategoryParam,
) (int64, error) {
	c, err := countPolicyCategories(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy category: %w", err)
	}
	return c, nil
}

// CountPolicyCategoriesWithSd SD付きでポリシーカテゴリー数を取得する。
func (a *PgAdapter) CountPolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WherePolicyCategoryParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countPolicyCategories(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count policy category: %w", err)
	}
	return c, nil
}

func createPolicyCategory(
	ctx context.Context, qtx *query.Queries, param parameter.CreatePolicyCategoryParam,
) (entity.PolicyCategory, error) {
	p := query.CreatePolicyCategoryParams{
		Name:        param.Name,
		Key:         param.Key,
		Description: param.Description,
	}
	e, err := qtx.CreatePolicyCategory(ctx, p)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to create policy category: %w", err)
	}
	entity := entity.PolicyCategory{
		PolicyCategoryID: e.PolicyCategoryID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
	}
	return entity, nil
}

// CreatePolicyCategory ポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyCategory(
	ctx context.Context, param parameter.CreatePolicyCategoryParam,
) (entity.PolicyCategory, error) {
	e, err := createPolicyCategory(ctx, a.query, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to create policy category: %w", err)
	}
	return e, nil
}

// CreatePolicyCategoryWithSd SD付きでポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyCategoryWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreatePolicyCategoryParam,
) (entity.PolicyCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := createPolicyCategory(ctx, qtx, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to create policy category: %w", err)
	}
	return e, nil
}

func createPolicyCategories(
	ctx context.Context, qtx *query.Queries, params []parameter.CreatePolicyCategoryParam,
) (int64, error) {
	p := make([]query.CreatePolicyCategoriesParams, len(params))
	for i, param := range params {
		p[i] = query.CreatePolicyCategoriesParams{
			Name:        param.Name,
			Key:         param.Key,
			Description: param.Description,
		}
	}
	c, err := qtx.CreatePolicyCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy categories: %w", err)
	}
	return c, nil
}

// CreatePolicyCategories ポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyCategories(
	ctx context.Context, params []parameter.CreatePolicyCategoryParam,
) (int64, error) {
	c, err := createPolicyCategories(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy categories: %w", err)
	}
	return c, nil
}

// CreatePolicyCategoriesWithSd SD付きでポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreatePolicyCategoryParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createPolicyCategories(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy categories: %w", err)
	}
	return c, nil
}

func deletePolicyCategory(ctx context.Context, qtx *query.Queries, policyCategoryID uuid.UUID) error {
	err := qtx.DeletePolicyCategory(ctx, policyCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

// DeletePolicyCategory ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategory(ctx context.Context, policyCategoryID uuid.UUID) error {
	err := deletePolicyCategory(ctx, a.query, policyCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

// DeletePolicyCategoryWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryWithSd(
	ctx context.Context, sd store.Sd, policyCategoryID uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePolicyCategory(ctx, qtx, policyCategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

func deletePolicyCategoryByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeletePolicyCategoryByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

// DeletePolicyCategoryByKey ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryByKey(ctx context.Context, key string) error {
	err := deletePolicyCategoryByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

// DeletePolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deletePolicyCategoryByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

func pluralDeletePolicyCategories(
	ctx context.Context, qtx *query.Queries, policyCategoryIDs []uuid.UUID,
) error {
	err := qtx.PluralDeletePolicyCategories(ctx, policyCategoryIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy categories: %w", err)
	}
	return nil
}

// PluralDeletePolicyCategories ポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicyCategories(
	ctx context.Context, policyCategoryIDs []uuid.UUID,
) error {
	err := pluralDeletePolicyCategories(ctx, a.query, policyCategoryIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy categories: %w", err)
	}
	return nil
}

// PluralDeletePolicyCategoriesWithSd SD付きでポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, policyCategoryIDs []uuid.UUID,
) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := pluralDeletePolicyCategories(ctx, qtx, policyCategoryIDs)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy categories: %w", err)
	}
	return nil
}

func findPolicyCategoryByID(
	ctx context.Context, qtx *query.Queries, policyCategoryID uuid.UUID,
) (entity.PolicyCategory, error) {
	e, err := qtx.FindPolicyCategoryByID(ctx, policyCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, store.ErrDataNoRecord
		}
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	entity := entity.PolicyCategory{
		PolicyCategoryID: e.PolicyCategoryID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
	}
	return entity, nil
}

// FindPolicyCategoryByID ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyCategoryByID(
	ctx context.Context, policyCategoryID uuid.UUID,
) (entity.PolicyCategory, error) {
	e, err := findPolicyCategoryByID(ctx, a.query, policyCategoryID)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	return e, nil
}

// FindPolicyCategoryByIDWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyCategoryByIDWithSd(
	ctx context.Context, sd store.Sd, policyCategoryID uuid.UUID,
) (entity.PolicyCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := findPolicyCategoryByID(ctx, qtx, policyCategoryID)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	return e, nil
}

func findPolicyCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PolicyCategory, error) {
	e, err := qtx.FindPolicyCategoryByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, store.ErrDataNoRecord
		}
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	entity := entity.PolicyCategory{
		PolicyCategoryID: e.PolicyCategoryID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
	}
	return entity, nil
}

// FindPolicyCategoryByKey ポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyCategoryByKey(ctx context.Context, key string) (entity.PolicyCategory, error) {
	e, err := findPolicyCategoryByKey(ctx, a.query, key)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	return e, nil
}

// FindPolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) FindPolicyCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.PolicyCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := findPolicyCategoryByKey(ctx, qtx, key)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category: %w", err)
	}
	return e, nil
}

// PolicyCategoryCursor is a cursor for PolicyCategory.
type PolicyCategoryCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getPolicyCategories(
	ctx context.Context, qtx *query.Queries, where parameter.WherePolicyCategoryParam,
	order parameter.PolicyCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyCategory], error) {
	eConvFunc := func(e query.PolicyCategory) (entity.PolicyCategory, error) {
		return entity.PolicyCategory{
			PolicyCategoryID: e.PolicyCategoryID,
			Name:             e.Name,
			Key:              e.Key,
			Description:      e.Description,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountPolicyCategoriesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountPolicyCategories(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count policy categories: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.PolicyCategory, error) {
		p := query.GetPolicyCategoriesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetPolicyCategories(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.PolicyCategory{}, nil
			}
			return nil, fmt.Errorf("failed to get policy categories: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.PolicyCategory, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.PolicyCategoryNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetPolicyCategoriesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetPolicyCategoriesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy categories: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.PolicyCategory, error) {
		p := query.GetPolicyCategoriesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetPolicyCategoriesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy categories: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.PolicyCategory) (entity.Int, any) {
		switch subCursor {
		case parameter.PolicyCategoryDefaultCursorKey:
			return entity.Int(e.MPolicyCategoriesPkey), nil
		case parameter.PolicyCategoryNameCursorKey:
			return entity.Int(e.MPolicyCategoriesPkey), e.Name
		}
		return entity.Int(e.MPolicyCategoriesPkey), nil
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
		return store.ListResult[entity.PolicyCategory]{}, fmt.Errorf("failed to get policy categories: %w", err)
	}
	return res, nil
}

// GetPolicyCategories ポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPolicyCategories(
	ctx context.Context,
	where parameter.WherePolicyCategoryParam,
	order parameter.PolicyCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyCategory], error) {
	r, err := getPolicyCategories(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{}, fmt.Errorf("failed to get policy categories: %w", err)
	}
	return r, nil
}

// GetPolicyCategoriesWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPolicyCategoriesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WherePolicyCategoryParam,
	order parameter.PolicyCategoryOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.PolicyCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PolicyCategory]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPolicyCategories(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{}, fmt.Errorf("failed to get policy categories: %w", err)
	}
	return r, nil
}

func getPluralPolicyCategories(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	p := query.GetPluralPolicyCategoriesParams{
		PolicyCategoryIds: ids,
		Offset:            int32(np.Offset.Int64),
		Limit:             int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralPolicyCategories(ctx, p)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{},
			fmt.Errorf("failed to get plural policy categories: %w", err)
	}
	entities := make([]entity.PolicyCategory, len(e))
	for i, v := range e {
		entities[i] = entity.PolicyCategory{
			PolicyCategoryID: v.PolicyCategoryID,
			Name:             v.Name,
			Key:              v.Key,
			Description:      v.Description,
		}
	}
	return store.ListResult[entity.PolicyCategory]{Data: entities}, nil
}

// GetPluralPolicyCategories ポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPluralPolicyCategories(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	r, err := getPluralPolicyCategories(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{},
			fmt.Errorf("failed to get plural policy categories: %w", err)
	}
	return r, nil
}

// GetPluralPolicyCategoriesWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPluralPolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PolicyCategory]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralPolicyCategories(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{},
			fmt.Errorf("failed to get plural policy categories: %w", err)
	}
	return r, nil
}

func updatePolicyCategory(
	ctx context.Context, qtx *query.Queries,
	policyCategoryID uuid.UUID, param parameter.UpdatePolicyCategoryParams,
) (entity.PolicyCategory, error) {
	p := query.UpdatePolicyCategoryParams{
		PolicyCategoryID: policyCategoryID,
		Name:             param.Name,
		Key:              param.Key,
		Description:      param.Description,
	}
	e, err := qtx.UpdatePolicyCategory(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, store.ErrDataNoRecord
		}
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	entity := entity.PolicyCategory{
		PolicyCategoryID: e.PolicyCategoryID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
	}
	return entity, nil
}

// UpdatePolicyCategory ポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyCategory(
	ctx context.Context, policyCategoryID uuid.UUID, param parameter.UpdatePolicyCategoryParams,
) (entity.PolicyCategory, error) {
	e, err := updatePolicyCategory(ctx, a.query, policyCategoryID, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	return e, nil
}

// UpdatePolicyCategoryWithSd SD付きでポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyCategoryWithSd(
	ctx context.Context, sd store.Sd, policyCategoryID uuid.UUID, param parameter.UpdatePolicyCategoryParams,
) (entity.PolicyCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePolicyCategory(ctx, qtx, policyCategoryID, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	return e, nil
}

func updatePolicyCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdatePolicyCategoryByKeyParams,
) (entity.PolicyCategory, error) {
	p := query.UpdatePolicyCategoryByKeyParams{
		Key:         key,
		Name:        param.Name,
		Description: param.Description,
	}
	e, err := qtx.UpdatePolicyCategoryByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, store.ErrDataNoRecord
		}
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	entity := entity.PolicyCategory{
		PolicyCategoryID: e.PolicyCategoryID,
		Name:             e.Name,
		Key:              e.Key,
		Description:      e.Description,
	}
	return entity, nil
}

// UpdatePolicyCategoryByKey ポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyCategoryByKey(
	ctx context.Context, key string, param parameter.UpdatePolicyCategoryByKeyParams,
) (entity.PolicyCategory, error) {
	e, err := updatePolicyCategoryByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	return e, nil
}

// UpdatePolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを更新する。
func (a *PgAdapter) UpdatePolicyCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdatePolicyCategoryByKeyParams,
) (entity.PolicyCategory, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.PolicyCategory{}, store.ErrNotFoundDescriptor
	}
	e, err := updatePolicyCategoryByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	return e, nil
}

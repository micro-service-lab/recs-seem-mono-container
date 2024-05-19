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
	return countPolicyCategories(ctx, a.query, where)
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
	return countPolicyCategories(ctx, qtx, where)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.PolicyCategory{}, errhandle.NewModelDuplicatedError("policy category")
		}
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
	return createPolicyCategory(ctx, a.query, param)
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
	return createPolicyCategory(ctx, qtx, param)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("policy category")
		}
		return 0, fmt.Errorf("failed to create policy categories: %w", err)
	}
	return c, nil
}

// CreatePolicyCategories ポリシーカテゴリーを作成する。
func (a *PgAdapter) CreatePolicyCategories(
	ctx context.Context, params []parameter.CreatePolicyCategoryParam,
) (int64, error) {
	return createPolicyCategories(ctx, a.query, params)
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
	return createPolicyCategories(ctx, qtx, params)
}

func deletePolicyCategory(ctx context.Context, qtx *query.Queries, policyCategoryID uuid.UUID) (int64, error) {
	c, err := qtx.DeletePolicyCategory(ctx, policyCategoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy category: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("policy category")
	}
	return c, nil
}

// DeletePolicyCategory ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategory(ctx context.Context, policyCategoryID uuid.UUID) (int64, error) {
	return deletePolicyCategory(ctx, a.query, policyCategoryID)
}

// DeletePolicyCategoryWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryWithSd(
	ctx context.Context, sd store.Sd, policyCategoryID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePolicyCategory(ctx, qtx, policyCategoryID)
}

func deletePolicyCategoryByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeletePolicyCategoryByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy category: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("policy category")
	}
	return c, nil
}

// DeletePolicyCategoryByKey ポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryByKey(ctx context.Context, key string) (int64, error) {
	return deletePolicyCategoryByKey(ctx, a.query, key)
}

// DeletePolicyCategoryByKeyWithSd SD付きでポリシーカテゴリーを削除する。
func (a *PgAdapter) DeletePolicyCategoryByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deletePolicyCategoryByKey(ctx, qtx, key)
}

func pluralDeletePolicyCategories(
	ctx context.Context, qtx *query.Queries, policyCategoryIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeletePolicyCategories(ctx, policyCategoryIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete policy categories: %w", err)
	}
	if c != int64(len(policyCategoryIDs)) {
		return 0, errhandle.NewModelNotFoundError("policy category")
	}
	return c, nil
}

// PluralDeletePolicyCategories ポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicyCategories(
	ctx context.Context, policyCategoryIDs []uuid.UUID,
) (int64, error) {
	return pluralDeletePolicyCategories(ctx, a.query, policyCategoryIDs)
}

// PluralDeletePolicyCategoriesWithSd SD付きでポリシーカテゴリーを複数削除する。
func (a *PgAdapter) PluralDeletePolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, policyCategoryIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeletePolicyCategories(ctx, qtx, policyCategoryIDs)
}

func findPolicyCategoryByID(
	ctx context.Context, qtx *query.Queries, policyCategoryID uuid.UUID,
) (entity.PolicyCategory, error) {
	e, err := qtx.FindPolicyCategoryByID(ctx, policyCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, errhandle.NewModelNotFoundError("policy category")
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
	return findPolicyCategoryByID(ctx, a.query, policyCategoryID)
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
	return findPolicyCategoryByID(ctx, qtx, policyCategoryID)
}

func findPolicyCategoryByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.PolicyCategory, error) {
	e, err := qtx.FindPolicyCategoryByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PolicyCategory{}, errhandle.NewModelNotFoundError("policy category")
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
	return findPolicyCategoryByKey(ctx, a.query, key)
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
	return findPolicyCategoryByKey(ctx, qtx, key)
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
	return getPolicyCategories(ctx, a.query, where, order, np, cp, wc)
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
	return getPolicyCategories(ctx, qtx, where, order, np, cp, wc)
}

func getPluralPolicyCategories(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID,
	order parameter.PolicyCategoryOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	var e []query.PolicyCategory
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralPolicyCategories(ctx, query.GetPluralPolicyCategoriesParams{
			PolicyCategoryIds: ids,
			OrderMethod:       order.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralPolicyCategoriesUseNumberedPaginate(
			ctx, query.GetPluralPolicyCategoriesUseNumberedPaginateParams{
				PolicyCategoryIds: ids,
				Offset:            int32(np.Offset.Int64),
				Limit:             int32(np.Limit.Int64),
			})
	}
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
	ctx context.Context, ids []uuid.UUID, order parameter.PolicyCategoryOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	return getPluralPolicyCategories(ctx, a.query, ids, order, np)
}

// GetPluralPolicyCategoriesWithSd SD付きでポリシーカテゴリーを取得する。
func (a *PgAdapter) GetPluralPolicyCategoriesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID,
	order parameter.PolicyCategoryOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.PolicyCategory], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.PolicyCategory]{}, store.ErrNotFoundDescriptor
	}
	return getPluralPolicyCategories(ctx, qtx, ids, order, np)
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
			return entity.PolicyCategory{}, errhandle.NewModelNotFoundError("policy category")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.PolicyCategory{}, errhandle.NewModelDuplicatedError("policy category")
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
	return updatePolicyCategory(ctx, a.query, policyCategoryID, param)
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
	return updatePolicyCategory(ctx, qtx, policyCategoryID, param)
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
			return entity.PolicyCategory{}, errhandle.NewModelNotFoundError("policy category")
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
	return updatePolicyCategoryByKey(ctx, a.query, key, param)
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
	return updatePolicyCategoryByKey(ctx, qtx, key, param)
}

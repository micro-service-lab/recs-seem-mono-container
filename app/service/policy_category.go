package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// PolicyCategoryKey ポリシーカテゴリーキー。
type PolicyCategoryKey string

const (
	// PolicyCategoryKeyUser ユーザー(ポリシー含む)に関するポリシー。
	PolicyCategoryKeyUser PolicyCategoryKey = "user"
	// PolicyCategoryKeyOrganization オーガナイゼーションに関するポリシー。
	PolicyCategoryKeyOrganization PolicyCategoryKey = "organization"
	// PolicyCategoryKeyAttendance 出欠関連に関するポリシー。
	PolicyCategoryKeyAttendance PolicyCategoryKey = "attendance"
	// PolicyCategoryKeyPosition 位置情報に関するポリシー。
	PolicyCategoryKeyPosition PolicyCategoryKey = "position"
)

// PolicyCategory ポリシーカテゴリー。
type PolicyCategory struct {
	Key         string
	Name        string
	Description string
}

// PolicyCategories ポリシーカテゴリー一覧。
var PolicyCategories = []PolicyCategory{
	{
		Key:         string(PolicyCategoryKeyOrganization),
		Name:        "オーガナイゼーション",
		Description: "オーガナイゼーションの設定に関するポリシー",
	},
	{
		Key:         string(PolicyCategoryKeyUser),
		Name:        "ユーザー",
		Description: "オーガナイゼーション所属ユーザー(ポリシー含む)に関するポリシー",
	},
	{
		Key:         string(PolicyCategoryKeyAttendance),
		Name:        "出欠関連",
		Description: "出欠関連に関するポリシー",
	},
	{
		Key:         string(PolicyCategoryKeyPosition),
		Name:        "位置情報",
		Description: "位置情報に関するポリシー",
	},
}

// ManagePolicyCategory ポリシーカテゴリー管理サービス。
type ManagePolicyCategory struct {
	DB store.Store
}

// CreatePolicyCategory ポリシーカテゴリーを作成する。
func (m *ManagePolicyCategory) CreatePolicyCategory(
	ctx context.Context,
	name, key, description string,
) (entity.PolicyCategory, error) {
	p := parameter.CreatePolicyCategoryParam{
		Name:        name,
		Key:         key,
		Description: description,
	}
	e, err := m.DB.CreatePolicyCategory(ctx, p)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to create policy category: %w", err)
	}
	return e, nil
}

// CreatePolicyCategories ポリシーカテゴリーを複数作成する。
func (m *ManagePolicyCategory) CreatePolicyCategories(
	ctx context.Context, ps []parameter.CreatePolicyCategoryParam,
) (int64, error) {
	es, err := m.DB.CreatePolicyCategories(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy categories: %w", err)
	}
	return es, nil
}

// UpdatePolicyCategory ポリシーカテゴリーを更新する。
func (m *ManagePolicyCategory) UpdatePolicyCategory(
	ctx context.Context, id uuid.UUID, name, key, description string,
) (entity.PolicyCategory, error) {
	p := parameter.UpdatePolicyCategoryParams{
		Name:        name,
		Key:         key,
		Description: description,
	}
	e, err := m.DB.UpdatePolicyCategory(ctx, id, p)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to update policy category: %w", err)
	}
	return e, nil
}

// DeletePolicyCategory ポリシーカテゴリーを削除する。
func (m *ManagePolicyCategory) DeletePolicyCategory(ctx context.Context, id uuid.UUID) error {
	err := m.DB.DeletePolicyCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete policy category: %w", err)
	}
	return nil
}

// PluralDeletePolicyCategories ポリシーカテゴリーを複数削除する。
func (m *ManagePolicyCategory) PluralDeletePolicyCategories(
	ctx context.Context, ids []uuid.UUID,
) error {
	err := m.DB.PluralDeletePolicyCategories(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to plural delete policy categories: %w", err)
	}
	return nil
}

// FindPolicyCategoryByID ポリシーカテゴリーをIDで取得する。
func (m *ManagePolicyCategory) FindPolicyCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.PolicyCategory, error) {
	e, err := m.DB.FindPolicyCategoryByID(ctx, id)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category by id: %w", err)
	}
	return e, nil
}

// FindPolicyCategoryByKey ポリシーカテゴリーをキーで取得する。
func (m *ManagePolicyCategory) FindPolicyCategoryByKey(
	ctx context.Context, key string,
) (entity.PolicyCategory, error) {
	e, err := m.DB.FindPolicyCategoryByKey(ctx, key)
	if err != nil {
		return entity.PolicyCategory{}, fmt.Errorf("failed to find policy category by key: %w", err)
	}
	return e, nil
}

// GetPolicyCategories ポリシーカテゴリーを取得する。
func (m *ManagePolicyCategory) GetPolicyCategories(
	ctx context.Context,
	whereSearchName string,
	order parameter.PolicyCategoryOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.PolicyCategory], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePolicyCategoryParam{
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
	r, err := m.DB.GetPolicyCategories(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyCategory]{}, fmt.Errorf("failed to get policy categories: %w", err)
	}
	return r, nil
}

// GetPolicyCategoriesCount ポリシーカテゴリーの数を取得する。
func (m *ManagePolicyCategory) GetPolicyCategoriesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WherePolicyCategoryParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountPolicyCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get policy categories count: %w", err)
	}
	return c, nil
}

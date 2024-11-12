package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// PolicyKey ポリシーキー。
type PolicyKey string

const (
	// PolicyKeyOrganizationCreate オーガナイゼーション作成。
	PolicyKeyOrganizationCreate PolicyKey = "organization.create"
	// PolicyKeyMemberCreate メンバー作成。
	PolicyKeyMemberCreate PolicyKey = "member.create"
	// PolicyKeyMemberDelete メンバー削除。
	PolicyKeyMemberDelete PolicyKey = "member.delete"
	// PolicyKeyRoleCreateRole ロールを作成。
	PolicyKeyRoleCreateRole PolicyKey = "role.create"
	// PolicyKeyRoleDeleteRole ロールを削除。
	PolicyKeyRoleDeleteRole PolicyKey = "role.delete"
	// PolicyKeyRoleUpdateRole ロールを更新。
	PolicyKeyRoleUpdateRole PolicyKey = "role.update"
	// PolicyKeyRoleAttachRole ロールを付与。
	PolicyKeyRoleAttachRole PolicyKey = "role.attach"
	// PolicyKeyRoleDetachRole ロールを剥奪。
	PolicyKeyRoleDetachRole PolicyKey = "role.detach"
	// PolicyKeyAttendanceViewLabIOHistory LabIO履歴を閲覧。
	PolicyKeyAttendanceViewLabIOHistory PolicyKey = "attendance.view_lab_io_history"
	// PolicyKeyPositionViewPositionHistory 位置情報履歴を閲覧。
	PolicyKeyPositionViewPositionHistory PolicyKey = "position.view_position_history"
)

// Policy ポリシー。
type Policy struct {
	Key               string
	Name              string
	Description       string
	PolicyCategoryKey PolicyCategoryKey
}

// Policies ポリシー一覧。
var Policies = []Policy{
	{
		Key:               string(PolicyKeyOrganizationCreate),
		Name:              "オーガナイゼーション作成",
		Description:       "オーガナイゼーションを作成する権限",
		PolicyCategoryKey: PolicyCategoryKeyOrganization,
	},
	{
		Key:               string(PolicyKeyMemberCreate),
		Name:              "メンバー作成",
		Description:       "メンバーを作成する権限",
		PolicyCategoryKey: PolicyCategoryKeyMember,
	},
	{
		Key:               string(PolicyKeyMemberDelete),
		Name:              "メンバー削除",
		Description:       "メンバーを削除する権限",
		PolicyCategoryKey: PolicyCategoryKeyMember,
	},
	{
		Key:               string(PolicyKeyRoleCreateRole),
		Name:              "ロール作成",
		Description:       "ロールを作成する権限",
		PolicyCategoryKey: PolicyCategoryKeyRole,
	},
	{
		Key:               string(PolicyKeyRoleDeleteRole),
		Name:              "ロール削除",
		Description:       "ロールを削除する権限",
		PolicyCategoryKey: PolicyCategoryKeyRole,
	},
	{
		Key:               string(PolicyKeyRoleUpdateRole),
		Name:              "ロール更新",
		Description:       "ロールを更新する権限",
		PolicyCategoryKey: PolicyCategoryKeyRole,
	},
	{
		Key:               string(PolicyKeyRoleAttachRole),
		Name:              "ロール付与",
		Description:       "ロールを付与する権限",
		PolicyCategoryKey: PolicyCategoryKeyRole,
	},
	{
		Key:               string(PolicyKeyRoleDetachRole),
		Name:              "ロール剥奪",
		Description:       "ロールを剥奪する権限",
		PolicyCategoryKey: PolicyCategoryKeyRole,
	},
	{
		Key:               string(PolicyKeyAttendanceViewLabIOHistory),
		Name:              "研究室入退室履歴閲覧",
		Description:       "研究室入退室履歴を閲覧する権限",
		PolicyCategoryKey: PolicyCategoryKeyAttendance,
	},
	{
		Key:               string(PolicyKeyPositionViewPositionHistory),
		Name:              "位置情報履歴閲覧",
		Description:       "位置情報履歴を閲覧する権限",
		PolicyCategoryKey: PolicyCategoryKeyPosition,
	},
}

// ManagePolicy ポリシー管理サービス。
type ManagePolicy struct {
	DB store.Store
}

// CreatePolicy ポリシーを作成する。
func (m *ManagePolicy) CreatePolicy(
	ctx context.Context,
	name, key, description string,
	categoryID uuid.UUID,
) (entity.Policy, error) {
	p := parameter.CreatePolicyParam{
		Name:             name,
		Key:              key,
		Description:      description,
		PolicyCategoryID: categoryID,
	}
	e, err := m.DB.CreatePolicy(ctx, p)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to create policy: %w", err)
	}
	return e, nil
}

// CreatePolicies ポリシーを複数作成する。
func (m *ManagePolicy) CreatePolicies(
	ctx context.Context, ps []parameter.CreatePolicyParam,
) (int64, error) {
	es, err := m.DB.CreatePolicies(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy: %w", err)
	}
	return es, nil
}

// UpdatePolicy ポリシーを更新する。
func (m *ManagePolicy) UpdatePolicy(
	ctx context.Context, id uuid.UUID, name, key, description string,
	categoryID uuid.UUID,
) (entity.Policy, error) {
	p := parameter.UpdatePolicyParams{
		Name:             name,
		Key:              key,
		Description:      description,
		PolicyCategoryID: categoryID,
	}
	e, err := m.DB.UpdatePolicy(ctx, id, p)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to update policy: %w", err)
	}
	return e, nil
}

// DeletePolicy ポリシーを削除する。
func (m *ManagePolicy) DeletePolicy(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeletePolicy(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	return c, nil
}

// PluralDeletePolicies ポリシーを複数削除する。
func (m *ManagePolicy) PluralDeletePolicies(
	ctx context.Context, ids []uuid.UUID,
) (int64, error) {
	c, err := m.DB.PluralDeletePolicies(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return c, nil
}

// FindPolicyByID ポリシーをIDで取得する。
func (m *ManagePolicy) FindPolicyByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Policy, error) {
	e, err := m.DB.FindPolicyByID(ctx, id)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy by id: %w", err)
	}
	return e, nil
}

// FindPolicyByIDWithCategory ポリシーとそのカテゴリをIDで取得する。
func (m *ManagePolicy) FindPolicyByIDWithCategory(
	ctx context.Context,
	id uuid.UUID,
) (entity.PolicyWithCategory, error) {
	e, err := m.DB.FindPolicyByIDWithCategory(ctx, id)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy by id with category: %w", err)
	}
	return e, nil
}

// FindPolicyByKey ポリシーをキーで取得する。
func (m *ManagePolicy) FindPolicyByKey(
	ctx context.Context, key string,
) (entity.Policy, error) {
	e, err := m.DB.FindPolicyByKey(ctx, key)
	if err != nil {
		return entity.Policy{}, fmt.Errorf("failed to find policy by key: %w", err)
	}
	return e, nil
}

// FindPolicyByKeyWithCategory ポリシーとそのカテゴリをキーで取得する。
func (m *ManagePolicy) FindPolicyByKeyWithCategory(
	ctx context.Context, key string,
) (entity.PolicyWithCategory, error) {
	e, err := m.DB.FindPolicyByKeyWithCategory(ctx, key)
	if err != nil {
		return entity.PolicyWithCategory{}, fmt.Errorf("failed to find policy by key with category: %w", err)
	}
	return e, nil
}

// GetPolicies ポリシーを取得する。
func (m *ManagePolicy) GetPolicies(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
	order parameter.PolicyOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Policy], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePolicyParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetPolicies(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Policy]{}, fmt.Errorf("failed to get policy: %w", err)
	}
	return r, nil
}

// GetPoliciesWithCategory ポリシーを取得する。
func (m *ManagePolicy) GetPoliciesWithCategory(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
	order parameter.PolicyOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.PolicyWithCategory], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePolicyParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetPoliciesWithCategory(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PolicyWithCategory]{},
			fmt.Errorf("failed to get policy with category: %w", err)
	}
	return r, nil
}

// GetPoliciesCount ポリシーの数を取得する。
func (m *ManagePolicy) GetPoliciesCount(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
) (int64, error) {
	p := parameter.WherePolicyParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
	}
	c, err := m.DB.CountPolicies(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get policy count: %w", err)
	}
	return c, nil
}

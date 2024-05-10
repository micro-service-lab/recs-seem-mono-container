package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// PermissionCategoryKey 権限カテゴリーキー。
type PermissionCategoryKey string

const (
	// PermissionCategoryKeyOrganization オーガナイゼーションの設定に関する権限。
	PermissionCategoryKeyOrganization PermissionCategoryKey = "organization"
	// PermissionCategoryKeyMember オーガナイゼーション所属メンバー(権限含む)に関する権限。
	PermissionCategoryKeyMember PermissionCategoryKey = "member"
	// PermissionCategoryKeyEvent イベントに関する権限。
	PermissionCategoryKeyEvent PermissionCategoryKey = "event"
	// PermissionCategoryKeyChatRoom チャットルームに関する権限。
	PermissionCategoryKeyChatRoom PermissionCategoryKey = "chat_room"
	// PermissionCategoryKeyRecord 議事録に関する権限。
	PermissionCategoryKeyRecord PermissionCategoryKey = "record"
)

// PermissionCategory 権限カテゴリー。
type PermissionCategory struct {
	Key         string
	Name        string
	Description string
}

// PermissionCategories 権限カテゴリー一覧。
var PermissionCategories = []PermissionCategory{
	{
		Key:         string(PermissionCategoryKeyOrganization),
		Name:        "オーガナイゼーション",
		Description: "オーガナイゼーションの設定に関する権限",
	},
	{
		Key:         string(PermissionCategoryKeyMember),
		Name:        "メンバー",
		Description: "オーガナイゼーション所属メンバー(権限含む)に関する権限",
	},
	{
		Key:         string(PermissionCategoryKeyEvent),
		Name:        "イベント",
		Description: "イベントに関する権限",
	},
	{
		Key:         string(PermissionCategoryKeyChatRoom),
		Name:        "チャットルーム",
		Description: "チャットルームに関する権限",
	},
	{
		Key:         string(PermissionCategoryKeyRecord),
		Name:        "議事録",
		Description: "議事録に関する権限",
	},
}

// ManagePermissionCategory 権限カテゴリー管理サービス。
type ManagePermissionCategory struct {
	DB store.Store
}

// CreatePermissionCategory 権限カテゴリーを作成する。
func (m *ManagePermissionCategory) CreatePermissionCategory(
	ctx context.Context,
	name, key, description string,
) (entity.PermissionCategory, error) {
	p := parameter.CreatePermissionCategoryParam{
		Name:        name,
		Key:         key,
		Description: description,
	}
	e, err := m.DB.CreatePermissionCategory(ctx, p)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to create permission category: %w", err)
	}
	return e, nil
}

// CreatePermissionCategories 権限カテゴリーを複数作成する。
func (m *ManagePermissionCategory) CreatePermissionCategories(
	ctx context.Context, ps []parameter.CreatePermissionCategoryParam,
) (int64, error) {
	es, err := m.DB.CreatePermissionCategories(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission categories: %w", err)
	}
	return es, nil
}

// UpdatePermissionCategory 権限カテゴリーを更新する。
func (m *ManagePermissionCategory) UpdatePermissionCategory(
	ctx context.Context, id uuid.UUID, name, key, description string,
) (entity.PermissionCategory, error) {
	p := parameter.UpdatePermissionCategoryParams{
		Name:        name,
		Key:         key,
		Description: description,
	}
	e, err := m.DB.UpdatePermissionCategory(ctx, id, p)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to update permission category: %w", err)
	}
	return e, nil
}

// DeletePermissionCategory 権限カテゴリーを削除する。
func (m *ManagePermissionCategory) DeletePermissionCategory(ctx context.Context, id uuid.UUID) error {
	err := m.DB.DeletePermissionCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission category: %w", err)
	}
	return nil
}

// PluralDeletePermissionCategories 権限カテゴリーを複数削除する。
func (m *ManagePermissionCategory) PluralDeletePermissionCategories(
	ctx context.Context, ids []uuid.UUID,
) error {
	err := m.DB.PluralDeletePermissionCategories(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to plural delete permission categories: %w", err)
	}
	return nil
}

// FindPermissionCategoryByID 権限カテゴリーをIDで取得する。
func (m *ManagePermissionCategory) FindPermissionCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.PermissionCategory, error) {
	e, err := m.DB.FindPermissionCategoryByID(ctx, id)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category by id: %w", err)
	}
	return e, nil
}

// FindPermissionCategoryByKey 権限カテゴリーをキーで取得する。
func (m *ManagePermissionCategory) FindPermissionCategoryByKey(
	ctx context.Context, key string,
) (entity.PermissionCategory, error) {
	e, err := m.DB.FindPermissionCategoryByKey(ctx, key)
	if err != nil {
		return entity.PermissionCategory{}, fmt.Errorf("failed to find permission category by key: %w", err)
	}
	return e, nil
}

// GetPermissionCategories 権限カテゴリーを取得する。
func (m *ManagePermissionCategory) GetPermissionCategories(
	ctx context.Context,
	whereSearchName string,
	order parameter.PermissionCategoryOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.PermissionCategory], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePermissionCategoryParam{
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
	r, err := m.DB.GetPermissionCategories(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PermissionCategory]{}, fmt.Errorf("failed to get permission categories: %w", err)
	}
	return r, nil
}

// GetPermissionCategoriesCount 権限カテゴリーの数を取得する。
func (m *ManagePermissionCategory) GetPermissionCategoriesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WherePermissionCategoryParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountPermissionCategories(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get permission categories count: %w", err)
	}
	return c, nil
}

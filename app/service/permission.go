package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// PermissionKey 権限キー。
type PermissionKey string

const (
	// PermissionKeyOrganizationDelete オーガナイゼーション削除。
	PermissionKeyOrganizationDelete PermissionKey = "organization.delete"
	// PermissionKeyOrganizationSetting オーガナイゼーション設定。
	PermissionKeyOrganizationSetting PermissionKey = "organization.setting"
	// PermissionKeyMemberInvite メンバー招待。
	PermissionKeyMemberInvite PermissionKey = "member.invite"
	// PermissionKeyMemberDelete メンバー削除。
	PermissionKeyMemberDelete PermissionKey = "member.delete"
	// PermissionKeyWorkPositionCreateWorkPosition メンバーワークポジション作成。
	PermissionKeyWorkPositionCreateWorkPosition PermissionKey = "work_position.create"
	// PermissionKeyWorkPositionDeleteWorkPosition メンバーワークポジション削除。
	PermissionKeyWorkPositionDeleteWorkPosition PermissionKey = "work_position.delete"
	// PermissionKeyWorkPositionUpdateWorkPosition メンバーワークポジション更新。
	PermissionKeyWorkPositionUpdateWorkPosition PermissionKey = "work_position.update"
	// PermissionKeyWorkPositionAttachWorkPosition メンバーワークポジション付与。
	PermissionKeyWorkPositionAttachWorkPosition PermissionKey = "work_position.attach"
	// PermissionKeyWorkPositionDetachWorkPosition メンバーワークポジション剥奪。
	PermissionKeyWorkPositionDetachWorkPosition PermissionKey = "work_position.detach"
	// PermissionKeyEventCreate イベント作成。
	PermissionKeyEventCreate PermissionKey = "event.create"
	// PermissionKeyEventDelete イベント削除。
	PermissionKeyEventDelete PermissionKey = "event.delete"
	// PermissionKeyEventUpdate イベント更新。
	PermissionKeyEventUpdate PermissionKey = "event.update"
	// PermissionKeyChatRoomSetting チャットルーム設定。
	PermissionKeyChatRoomSetting PermissionKey = "chat_room.setting"
	// PermissionKeyRecordCreate 議事録作成。
	PermissionKeyRecordCreate PermissionKey = "record.create"
	// PermissionKeyRecordDelete 議事録削除。
	PermissionKeyRecordDelete PermissionKey = "record.delete"
	// PermissionKeyRecordUpdate 議事録更新。
	PermissionKeyRecordUpdate PermissionKey = "record.update"
)

// Permission 権限。
type Permission struct {
	Key                   string
	Name                  string
	Description           string
	PermissionCategoryKey PermissionCategoryKey
}

// Permissions 権限一覧。
var Permissions = []Permission{
	{
		Key:                   string(PermissionKeyOrganizationDelete),
		Name:                  "オーガナイゼーション削除",
		Description:           "オーガナイゼーションを削除する権限",
		PermissionCategoryKey: PermissionCategoryKeyOrganization,
	},
	{
		Key:                   string(PermissionKeyOrganizationSetting),
		Name:                  "オーガナイゼーション設定",
		Description:           "オーガナイゼーションの設定を変更する権限",
		PermissionCategoryKey: PermissionCategoryKeyOrganization,
	},
	{
		Key:                   string(PermissionKeyMemberInvite),
		Name:                  "メンバー招待",
		Description:           "オーガナイゼーションにメンバーを招待する権限",
		PermissionCategoryKey: PermissionCategoryKeyMember,
	},
	{
		Key:                   string(PermissionKeyMemberDelete),
		Name:                  "メンバー削除",
		Description:           "オーガナイゼーションからメンバーを削除する権限",
		PermissionCategoryKey: PermissionCategoryKeyMember,
	},
	{
		Key:                   string(PermissionKeyWorkPositionCreateWorkPosition),
		Name:                  "メンバーワークポジション作成",
		Description:           "オーガナイゼーションのワークポジションを作成する権限",
		PermissionCategoryKey: PermissionCategoryKeyWorkPosition,
	},
	{
		Key:                   string(PermissionKeyWorkPositionDeleteWorkPosition),
		Name:                  "メンバーワークポジション削除",
		Description:           "オーガナイゼーションのワークポジションを削除する権限",
		PermissionCategoryKey: PermissionCategoryKeyWorkPosition,
	},
	{
		Key:                   string(PermissionKeyWorkPositionUpdateWorkPosition),
		Name:                  "メンバーワークポジション更新",
		Description:           "オーガナイゼーションのワークポジションを更新する権限",
		PermissionCategoryKey: PermissionCategoryKeyWorkPosition,
	},
	{
		Key:                   string(PermissionKeyWorkPositionAttachWorkPosition),
		Name:                  "メンバーワークポジション付与",
		Description:           "オーガナイゼーションのワークポジションをメンバーに付与する権限",
		PermissionCategoryKey: PermissionCategoryKeyWorkPosition,
	},
	{
		Key:                   string(PermissionKeyWorkPositionDetachWorkPosition),
		Name:                  "メンバーワークポジション剥奪",
		Description:           "オーガナイゼーションのワークポジションをメンバーから剥奪する権限",
		PermissionCategoryKey: PermissionCategoryKeyWorkPosition,
	},
	{
		Key:                   string(PermissionKeyEventCreate),
		Name:                  "イベント作成",
		Description:           "イベントを作成する権限",
		PermissionCategoryKey: PermissionCategoryKeyEvent,
	},
	{
		Key:                   string(PermissionKeyEventDelete),
		Name:                  "イベント削除",
		Description:           "イベントを削除する権限",
		PermissionCategoryKey: PermissionCategoryKeyEvent,
	},
	{
		Key:                   string(PermissionKeyEventUpdate),
		Name:                  "イベント更新",
		Description:           "イベント作成者でなくてもイベントを更新する権限",
		PermissionCategoryKey: PermissionCategoryKeyEvent,
	},
	{
		Key:                   string(PermissionKeyChatRoomSetting),
		Name:                  "チャットルーム設定",
		Description:           "オーガナイゼーションに紐付いたチャットルームの設定を変更する権限",
		PermissionCategoryKey: PermissionCategoryKeyChatRoom,
	},
	{
		Key:                   string(PermissionKeyRecordCreate),
		Name:                  "議事録作成",
		Description:           "議事録を作成する権限",
		PermissionCategoryKey: PermissionCategoryKeyRecord,
	},
	{
		Key:                   string(PermissionKeyRecordDelete),
		Name:                  "議事録削除",
		Description:           "議事録を削除する権限",
		PermissionCategoryKey: PermissionCategoryKeyRecord,
	},
	{
		Key:                   string(PermissionKeyRecordUpdate),
		Name:                  "議事録更新",
		Description:           "投稿者でなくても議事録を更新する権限",
		PermissionCategoryKey: PermissionCategoryKeyRecord,
	},
}

// ManagePermission 権限管理サービス。
type ManagePermission struct {
	DB store.Store
}

// CreatePermission 権限を作成する。
func (m *ManagePermission) CreatePermission(
	ctx context.Context,
	name, key, description string,
	categoryID uuid.UUID,
) (entity.Permission, error) {
	p := parameter.CreatePermissionParam{
		Name:                 name,
		Key:                  key,
		Description:          description,
		PermissionCategoryID: categoryID,
	}
	e, err := m.DB.CreatePermission(ctx, p)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to create permission: %w", err)
	}
	return e, nil
}

// CreatePermissions 権限を複数作成する。
func (m *ManagePermission) CreatePermissions(
	ctx context.Context, ps []parameter.CreatePermissionParam,
) (int64, error) {
	es, err := m.DB.CreatePermissions(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}
	return es, nil
}

// UpdatePermission 権限を更新する。
func (m *ManagePermission) UpdatePermission(
	ctx context.Context, id uuid.UUID, name, key, description string,
	categoryID uuid.UUID,
) (entity.Permission, error) {
	p := parameter.UpdatePermissionParams{
		Name:                 name,
		Key:                  key,
		Description:          description,
		PermissionCategoryID: categoryID,
	}
	e, err := m.DB.UpdatePermission(ctx, id, p)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to update permission: %w", err)
	}
	return e, nil
}

// DeletePermission 権限を削除する。
func (m *ManagePermission) DeletePermission(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeletePermission(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete permission: %w", err)
	}
	return c, nil
}

// PluralDeletePermissions 権限を複数削除する。
func (m *ManagePermission) PluralDeletePermissions(
	ctx context.Context, ids []uuid.UUID,
) (int64, error) {
	c, err := m.DB.PluralDeletePermissions(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete permission: %w", err)
	}
	return c, nil
}

// FindPermissionByID 権限をIDで取得する。
func (m *ManagePermission) FindPermissionByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Permission, error) {
	e, err := m.DB.FindPermissionByID(ctx, id)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission by id: %w", err)
	}
	return e, nil
}

// FindPermissionByIDWithCategory 権限とそのカテゴリをIDで取得する。
func (m *ManagePermission) FindPermissionByIDWithCategory(
	ctx context.Context,
	id uuid.UUID,
) (entity.PermissionWithCategory, error) {
	e, err := m.DB.FindPermissionByIDWithCategory(ctx, id)
	if err != nil {
		return entity.PermissionWithCategory{}, fmt.Errorf("failed to find permission by id with category: %w", err)
	}
	return e, nil
}

// FindPermissionByKey 権限をキーで取得する。
func (m *ManagePermission) FindPermissionByKey(
	ctx context.Context, key string,
) (entity.Permission, error) {
	e, err := m.DB.FindPermissionByKey(ctx, key)
	if err != nil {
		return entity.Permission{}, fmt.Errorf("failed to find permission by key: %w", err)
	}
	return e, nil
}

// FindPermissionByKeyWithCategory 権限とそのカテゴリをキーで取得する。
func (m *ManagePermission) FindPermissionByKeyWithCategory(
	ctx context.Context, key string,
) (entity.PermissionWithCategory, error) {
	e, err := m.DB.FindPermissionByKeyWithCategory(ctx, key)
	if err != nil {
		return entity.PermissionWithCategory{}, fmt.Errorf("failed to find permission by key with category: %w", err)
	}
	return e, nil
}

// GetPermissions 権限を取得する。
func (m *ManagePermission) GetPermissions(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
	order parameter.PermissionOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Permission], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePermissionParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
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
	r, err := m.DB.GetPermissions(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Permission]{}, fmt.Errorf("failed to get permission: %w", err)
	}
	return r, nil
}

// GetPermissionsWithCategory 権限を取得する。
func (m *ManagePermission) GetPermissionsWithCategory(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
	order parameter.PermissionOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.PermissionWithCategory], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WherePermissionParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
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
	r, err := m.DB.GetPermissionsWithCategory(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.PermissionWithCategory]{},
			fmt.Errorf("failed to get permission with category: %w", err)
	}
	return r, nil
}

// GetPermissionsCount 権限の数を取得する。
func (m *ManagePermission) GetPermissionsCount(
	ctx context.Context,
	whereSearchName string,
	whereInCategories []uuid.UUID,
) (int64, error) {
	p := parameter.WherePermissionParam{
		WhereLikeName:   whereSearchName != "",
		SearchName:      whereSearchName,
		WhereInCategory: len(whereInCategories) > 0,
		InCategories:    whereInCategories,
	}
	c, err := m.DB.CountPermissions(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get permission count: %w", err)
	}
	return c, nil
}

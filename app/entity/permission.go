package entity

import "github.com/google/uuid"

// Permission 権限を表す構造体。
type Permission struct {
	PermissionID         uuid.UUID `json:"permission_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Key                  string    `json:"key"`
	PermissionCategoryID uuid.UUID `json:"permission_category_id"`
}

// PermissionWithCategory 権限とカテゴリを表す構造体。
type PermissionWithCategory struct {
	Permission
	PermissionCategory PermissionCategory `json:"permission_category"`
}

// PermissionWithCategoryForQuery 権限とカテゴリを表す構造体(クエリ用)。
type PermissionWithCategoryForQuery struct {
	Pkey Int `json:"-"`
	PermissionWithCategory
}

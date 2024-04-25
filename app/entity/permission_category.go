package entity

import "github.com/google/uuid"

// PermissionCategory 権限カテゴリを表す構造体。
type PermissionCategory struct {
	PermissionCategoryID uuid.UUID `json:"permission_category_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Key                  string    `json:"key"`
}

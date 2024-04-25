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

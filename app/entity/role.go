package entity

import (
	"github.com/google/uuid"
)

// Role ロールを表す構造体。
type Role struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

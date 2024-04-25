package entity

import (
	"github.com/google/uuid"
)

// RoleAssociation ロールとポリシーの関連を表す構造体。
type RoleAssociation struct {
	RoleID   uuid.UUID `json:"role_id"`
	PolicyID uuid.UUID `json:"policy_id"`
}

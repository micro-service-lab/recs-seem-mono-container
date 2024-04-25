package entity

import "github.com/google/uuid"

// Group グループを表す構造体。
type Group struct {
	GroupID        uuid.UUID `json:"group_id"`
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

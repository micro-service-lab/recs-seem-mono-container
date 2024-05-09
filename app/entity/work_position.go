package entity

import (
	"github.com/google/uuid"
)

// WorkPosition 役職を表す構造体。
type WorkPosition struct {
	WorkPositionID uuid.UUID `json:"work_position_id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Description    string    `json:"description"`
}

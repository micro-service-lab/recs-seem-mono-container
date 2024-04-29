package entity

import (
	"github.com/google/uuid"
)

// Organization 組織を表す構造体。
type Organization struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    String    `json:"description"`
	IsPersonal     bool      `json:"is_personal"`
	IsWhole        bool      `json:"is_whole"`
	ChatRoomID     UUID      `json:"chat_room_id"`
}

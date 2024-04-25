package entity

import (
	"time"

	"github.com/google/uuid"
)

// Record レコードを表す構造体。
type Record struct {
	RecordID       uuid.UUID `json:"record_id"`
	RecordTypeID   uuid.UUID `json:"record_type_id"`
	Title          string    `json:"title"`
	Body           String    `json:"body"`
	OrganizationID UUID      `json:"organization_id"`
	PostedBy       UUID      `json:"posted_by"`
	LastEditedBy   UUID      `json:"last_edited_by"`
	PostedAt       time.Time `json:"posted_at"`
	LastEditedAt   time.Time `json:"last_edited_at"`
}

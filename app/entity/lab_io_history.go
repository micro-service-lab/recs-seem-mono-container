package entity

import (
	"time"

	"github.com/google/uuid"
)

// LabIOHistory ラボ出入履歴を表す構造体。
type LabIOHistory struct {
	LabIoHistoryID uuid.UUID   `json:"lab_io_history_id"`
	MemberID       uuid.UUID   `json:"member_id"`
	EnteredAt      time.Time   `json:"entered_at"`
	ExitedAt       Timestamptz `json:"exited_at"`
}

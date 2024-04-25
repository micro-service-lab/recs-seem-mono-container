package entity

import (
	"time"

	"github.com/google/uuid"
)

// EarlyLeaving 早退を表す構造体。
type EarlyLeaving struct {
	EarlyLeavingID uuid.UUID `json:"early_leaving_id"`
	AttendanceID   uuid.UUID `json:"attendance_id"`
	LeaveTime      time.Time `json:"leave_time"`
}

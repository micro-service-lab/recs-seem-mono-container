package entity

import "github.com/google/uuid"

// AttendStatus 出欠ステータスを表す構造体。
type AttendStatus struct {
	AttendStatusID uuid.UUID `json:"attend_status_id"`
	Name           string    `json:"name"`
	Key            string    `json:"key"`
}

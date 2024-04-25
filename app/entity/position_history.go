package entity

import (
	"time"

	"github.com/google/uuid"
)

// PositionHistory 位置情報の履歴を表す構造体。
type PositionHistory struct {
	PositionHistoryID uuid.UUID `json:"position_history_id"`
	MemberID          uuid.UUID `json:"member_id"`
	XPos              float64   `json:"x_pos"`
	YPos              float64   `json:"y_pos"`
	SentAt            time.Time `json:"sent_at"`
}

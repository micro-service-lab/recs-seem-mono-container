package entity

import (
	"time"

	"github.com/google/uuid"
)

// LateArrival 遅刻を表す構造体。
type LateArrival struct {
	LateArrivalID uuid.UUID `json:"late_arrival_id"`
	AttendanceID  uuid.UUID `json:"attendance_id"`
	ArriveTime    time.Time `json:"arrive_time"`
}

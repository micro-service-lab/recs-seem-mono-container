package entity

import "github.com/google/uuid"

// AttendanceType 出席種別を表す構造体。
type AttendanceType struct {
	AttendanceTypeID uuid.UUID `json:"attendance_type_id"`
	Name             string    `json:"name"`
	Key              string    `json:"key"`
	Color            string    `json:"color"`
}

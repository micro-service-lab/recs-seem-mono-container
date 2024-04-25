package entity

import "github.com/google/uuid"

// Absence 欠席を表す構造体。
type Absence struct {
	AbsenceID    uuid.UUID `json:"absence_id"`
	AttendanceID uuid.UUID `json:"attendance_id"`
}

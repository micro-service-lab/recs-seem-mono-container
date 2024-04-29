package parameter

import "github.com/google/uuid"

// CreateAbsenceParam 欠席作成のパラメータ。
type CreateAbsenceParam struct {
	AttendanceID uuid.UUID
}

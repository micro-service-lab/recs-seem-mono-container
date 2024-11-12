package entity

import (
	"time"

	"github.com/google/uuid"
)

// Attendance 出席を表す構造体。
type Attendance struct {
	AttendanceID       uuid.UUID `json:"attendance_id"`
	AttendanceTypeID   uuid.UUID `json:"attendance_type_id"`
	MemberID           uuid.UUID `json:"member_id"`
	Description        string    `json:"description"`
	StartDate          Date      `json:"start_date"`
	EndDate            Date      `json:"end_date"`
	MailSendFlag       bool      `json:"mail_send_flag"`
	SendOrganizationID UUID      `json:"send_organization_id"`
	PostedAt           time.Time `json:"posted_at"`
	LastEditedAt       time.Time `json:"last_edited_at"`
}

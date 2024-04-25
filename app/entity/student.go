package entity

import "github.com/google/uuid"

// Student 生徒を表す構造体。
type Student struct {
	StudentID uuid.UUID `json:"student_id"`
	MemberID  uuid.UUID `json:"member_id"`
}

package entity

import "github.com/google/uuid"

// Professor 教授を表す構造体。
type Professor struct {
	ProfessorID uuid.UUID `json:"professor_id"`
	MemberID    uuid.UUID `json:"member_id"`
}

package entity

import "github.com/google/uuid"

// Grade 学年を表す構造体。
type Grade struct {
	GradeID        uuid.UUID `json:"grade_id"`
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

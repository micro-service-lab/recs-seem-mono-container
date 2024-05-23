package entity

import "github.com/google/uuid"

// Grade 学年を表す構造体。
type Grade struct {
	GradeID        uuid.UUID `json:"grade_id"`
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

// GradeWithOrganization 学年と組織を表す構造体。
type GradeWithOrganization struct {
	GradeID      uuid.UUID `json:"grade_id"`
	Key          string    `json:"key"`
	Organization Organization
}

// GradeWithOrganizationForQuery 学年と組織を表す構造体(クエリー用)。
type GradeWithOrganizationForQuery struct {
	Pkey Int `json:"-"`
	GradeWithOrganization
}

package entity

import "github.com/google/uuid"

// Professor 教授を表す構造体。
type Professor struct {
	ProfessorID uuid.UUID `json:"professor_id"`
	MemberID    uuid.UUID `json:"member_id"`
}

// ProfessorWithMember 教授とメンバーを表す構造体。
type ProfessorWithMember struct {
	ProfessorID uuid.UUID  `json:"professor_id"`
	Member      MemberCard `json:"member"`
}

// ProfessorWithMemberForQuery 教授とメンバーを表す構造体(クエリー用)。
type ProfessorWithMemberForQuery struct {
	Pkey Int `json:"-"`
	ProfessorWithMember
}

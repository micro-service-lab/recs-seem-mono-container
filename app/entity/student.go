package entity

import "github.com/google/uuid"

// Student 生徒を表す構造体。
type Student struct {
	StudentID uuid.UUID `json:"student_id"`
	MemberID  uuid.UUID `json:"member_id"`
}

// StudentWithMember 生徒とメンバーを表す構造体。
type StudentWithMember struct {
	StudentID uuid.UUID  `json:"student_id"`
	Member    MemberCard `json:"member"`
}

// StudentWithMemberForQuery 生徒とメンバーを表す構造体(クエリー用)。
type StudentWithMemberForQuery struct {
	Pkey Int `json:"-"`
	StudentWithMember
}

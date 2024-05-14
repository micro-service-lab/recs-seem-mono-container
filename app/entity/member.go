package entity

import "github.com/google/uuid"

// Member メンバーを表す構造体。
type Member struct {
	MemberID               uuid.UUID `json:"member_id"`
	Email                  string    `json:"email"`
	Name                   string    `json:"name"`
	AttendStatusID         uuid.UUID `json:"attend_status_id"`
	ProfileImageID         UUID      `json:"profile_image_id"`
	GradeID                uuid.UUID `json:"grade_id"`
	GroupID                uuid.UUID `json:"group_id"`
	PersonalOrganizationID uuid.UUID `json:"personal_organization_id"`
	RoleID                 UUID      `json:"role_id"`
}

// PracticeMember 練習メンバーを表す構造体。
type PracticeMember struct {
	MemberID               uuid.UUID               `json:"member_id"`
	Email                  string                  `json:"email"`
	Name                   string                  `json:"name"`
	AttendStatusID         uuid.UUID               `json:"attend_status_id"`
	ProfileImage           ImageWithAttachableItem `json:"profile_image"`
	GradeID                uuid.UUID               `json:"grade_id"`
	GroupID                uuid.UUID               `json:"group_id"`
	PersonalOrganizationID uuid.UUID               `json:"personal_organization_id"`
	RoleID                 UUID                    `json:"role_id"`
}

// MemberCard メンバーカードを表す構造体。
type MemberCard struct {
	MemberID     uuid.UUID               `json:"member_id"`
	Name         String                  `json:"name"`
	Email        String                  `json:"email"`
	ProfileImage ImageWithAttachableItem `json:"profile_image"`
}

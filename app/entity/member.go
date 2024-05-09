package entity

import "github.com/google/uuid"

// Member メンバーを表す構造体。
type Member struct {
	MemberID               uuid.UUID `json:"member_id"`
	Email                  string    `json:"email"`
	Name                   string    `json:"name"`
	AttendStatusID         uuid.UUID `json:"attend_status_id"`
	ProfileImageURL        String    `json:"profile_image_url"`
	GradeID                uuid.UUID `json:"grade_id"`
	GroupID                uuid.UUID `json:"group_id"`
	PersonalOrganizationID uuid.UUID `json:"personal_organization_id"`
	RoleID                 UUID      `json:"role_id"`
}

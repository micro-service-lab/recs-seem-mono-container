package entity

import (
	"github.com/google/uuid"
)

// Member メンバーを表す構造体。
type Member struct {
	MemberID               uuid.UUID `json:"member_id"`
	Email                  string    `json:"email"`
	Name                   string    `json:"name"`
	FirstName              String    `json:"first_name"`
	LastName               String    `json:"last_name"`
	AttendStatusID         uuid.UUID `json:"attend_status_id"`
	ProfileImageID         UUID      `json:"profile_image_id"`
	GradeID                uuid.UUID `json:"grade_id"`
	GroupID                uuid.UUID `json:"group_id"`
	PersonalOrganizationID uuid.UUID `json:"personal_organization_id"`
	RoleID                 UUID      `json:"role_id"`
}

// AuthMember 認証済みのメンバーを表す構造体。
type AuthMember struct {
	MemberID               uuid.UUID                        `json:"member_id"`
	Email                  string                           `json:"email"`
	Name                   string                           `json:"name"`
	FirstName              String                           `json:"first_name"`
	LastName               String                           `json:"last_name"`
	AttendStatusID         uuid.UUID                        `json:"attend_status_id"`
	ProfileImageID         UUID                             `json:"profile_image_id"`
	GradeID                uuid.UUID                        `json:"grade_id"`
	GroupID                uuid.UUID                        `json:"group_id"`
	PersonalOrganizationID uuid.UUID                        `json:"personal_organization_id"`
	Role                   NullableEntity[RoleWithPolicies] `json:"role"`
}

// AuthPayload 認証情報ペイロード
type AuthPayload struct {
	MemberID uuid.UUID `json:"member_id"`
}

// MemberCredentials メンバーの認証情報を表す構造体。
type MemberCredentials struct {
	MemberID uuid.UUID `json:"member_id"`
	LoginID  string    `json:"login_id"`
	Password string    `json:"password"`
}

// MemberWithAttendStatus メンバーと出席状況を表す構造体。
type MemberWithAttendStatus struct {
	MemberID               uuid.UUID    `json:"member_id"`
	Email                  string       `json:"email"`
	Name                   string       `json:"name"`
	FirstName              String       `json:"first_name"`
	LastName               String       `json:"last_name"`
	AttendStatus           AttendStatus `json:"attend_status"`
	ProfileImageID         UUID         `json:"profile_image_id"`
	GradeID                uuid.UUID    `json:"grade_id"`
	GroupID                uuid.UUID    `json:"group_id"`
	PersonalOrganizationID uuid.UUID    `json:"personal_organization_id"`
	RoleID                 UUID         `json:"role_id"`
}

// MemberWithAttendStatusForQuery メンバーと出席状況を表す構造体(クエリー用)。
type MemberWithAttendStatusForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithAttendStatus
}

// MemberWithProfileImage プロフィール画像付きのメンバーを表す構造体。
type MemberWithProfileImage struct {
	MemberID               uuid.UUID                               `json:"member_id"`
	Email                  string                                  `json:"email"`
	Name                   string                                  `json:"name"`
	FirstName              String                                  `json:"first_name"`
	LastName               String                                  `json:"last_name"`
	AttendStatusID         uuid.UUID                               `json:"attend_status_id"`
	ProfileImage           NullableEntity[ImageWithAttachableItem] `json:"profile_image"`
	GradeID                uuid.UUID                               `json:"grade_id"`
	GroupID                uuid.UUID                               `json:"group_id"`
	PersonalOrganizationID uuid.UUID                               `json:"personal_organization_id"`
	RoleID                 UUID                                    `json:"role_id"`
}

// MemberWithProfileImageForQuery プロフィール画像付きのメンバーを表す構造体(クエリー用)。
type MemberWithProfileImageForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithProfileImage
}

// MemberWithDetail 詳細情報付きのメンバーを表す構造体。
type MemberWithDetail struct {
	MemberID               uuid.UUID                 `json:"member_id"`
	Email                  string                    `json:"email"`
	Name                   string                    `json:"name"`
	FirstName              String                    `json:"first_name"`
	LastName               String                    `json:"last_name"`
	AttendStatusID         uuid.UUID                 `json:"attend_status_id"`
	ProfileImageID         UUID                      `json:"profile_image_id"`
	GradeID                uuid.UUID                 `json:"grade_id"`
	GroupID                uuid.UUID                 `json:"group_id"`
	PersonalOrganizationID uuid.UUID                 `json:"personal_organization_id"`
	RoleID                 UUID                      `json:"role_id"`
	Student                NullableEntity[Student]   `json:"student,omitempty"`
	Professor              NullableEntity[Professor] `json:"professor,omitempty"`
}

// MemberWithDetailForQuery 詳細情報付きのメンバーを表す構造体(クエリー用)。
type MemberWithDetailForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithDetail
}

// MemberWithCrew クルー付きのメンバーを表す構造体。
type MemberWithCrew struct {
	MemberID               uuid.UUID             `json:"member_id"`
	Email                  string                `json:"email"`
	Name                   string                `json:"name"`
	FirstName              String                `json:"first_name"`
	LastName               String                `json:"last_name"`
	AttendStatusID         uuid.UUID             `json:"attend_status_id"`
	ProfileImageID         UUID                  `json:"profile_image_id"`
	Grade                  GradeWithOrganization `json:"grade"`
	Group                  GroupWithOrganization `json:"group"`
	PersonalOrganizationID uuid.UUID             `json:"personal_organization_id"`
	RoleID                 UUID                  `json:"role_id"`
}

// MemberWithCrewForQuery クルー付きのメンバーを表す構造体(クエリー用)。
type MemberWithCrewForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithCrew
}

// MemberWithRole ロール付きのメンバーを表す構造体。
type MemberWithRole struct {
	MemberID               uuid.UUID            `json:"member_id"`
	Email                  string               `json:"email"`
	Name                   string               `json:"name"`
	FirstName              String               `json:"first_name"`
	LastName               String               `json:"last_name"`
	AttendStatusID         uuid.UUID            `json:"attend_status_id"`
	ProfileImageID         UUID                 `json:"profile_image_id"`
	GradeID                uuid.UUID            `json:"grade_id"`
	GroupID                uuid.UUID            `json:"group_id"`
	PersonalOrganizationID uuid.UUID            `json:"personal_organization_id"`
	Role                   NullableEntity[Role] `json:"role"`
}

// MemberWithRoleForQuery ロール付きのメンバーを表す構造体(クエリー用)。
type MemberWithRoleForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithRole
}

// MemberWithPersonalOrganization 個人組織付きのメンバーを表す構造体。
type MemberWithPersonalOrganization struct {
	MemberID             uuid.UUID    `json:"member_id"`
	Email                string       `json:"email"`
	Name                 string       `json:"name"`
	FirstName            String       `json:"first_name"`
	LastName             String       `json:"last_name"`
	AttendStatusID       uuid.UUID    `json:"attend_status_id"`
	ProfileImageID       UUID         `json:"profile_image_id"`
	GradeID              uuid.UUID    `json:"grade_id"`
	GroupID              uuid.UUID    `json:"group_id"`
	PersonalOrganization Organization `json:"personal_organization"`
	RoleID               UUID         `json:"role_id"`
}

// MemberWithPersonalOrganizationForQuery 個人組織付きのメンバーを表す構造体(クエリー用)。
type MemberWithPersonalOrganizationForQuery struct {
	Pkey Int `json:"pkey"`
	MemberWithPersonalOrganization
}

// MemberCard メンバーカードを表す構造体。
type MemberCard struct {
	MemberID     uuid.UUID                               `json:"member_id"`
	Name         string                                  `json:"name"`
	FirstName    String                                  `json:"first_name"`
	LastName     String                                  `json:"last_name"`
	Email        string                                  `json:"email"`
	ProfileImage NullableEntity[ImageWithAttachableItem] `json:"profile_image"`
	GradeID      uuid.UUID                               `json:"grade_id"`
	GroupID      uuid.UUID                               `json:"group_id"`
}

// SimpleMember シンプルなメンバーを表す構造体。
type SimpleMember struct {
	MemberID       uuid.UUID `json:"member_id"`
	Name           string    `json:"name"`
	FirstName      String    `json:"first_name"`
	LastName       String    `json:"last_name"`
	Email          string    `json:"email"`
	ProfileImageID UUID      `json:"profile_image_id"`
	GradeID        uuid.UUID `json:"grade_id"`
	GroupID        uuid.UUID `json:"group_id"`
}

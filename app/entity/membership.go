package entity

import (
	"time"

	"github.com/google/uuid"
)

// Membership メンバーシップを表す構造体。
type Membership struct {
	MemberID       uuid.UUID `json:"member_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	WorkPositionID UUID      `json:"work_position_id"`
	AddedAt        time.Time `json:"added_at"`
}

// MembershipMember メンバーシップのメンバーを表す構造体。
type MembershipMember struct {
	Member         MemberCard `json:"member"`
	WorkPositionID UUID       `json:"work_position_id"`
	AddedAt        time.Time  `json:"added_at"`
}

// MemberOnOrganization オーガナイゼーション上のメンバーを表す構造体。
type MemberOnOrganization struct {
	Member         MemberCard `json:"member"`
	WorkPositionID UUID       `json:"work_position_id"`
	AddedAt        time.Time  `json:"added_at"`
}

// MemberOnOrganizationForQuery オーガナイゼーション上のメンバーを表す構造体(クエリ用)。
type MemberOnOrganizationForQuery struct {
	Pkey Int `json:"-"`
	MemberOnOrganization
}

// OrganizationOnMember メンバー上のオーガナイゼーションを表す構造体。
type OrganizationOnMember struct {
	Organization   Organization `json:"organization"`
	WorkPositionID UUID         `json:"work_position_id"`
	AddedAt        time.Time    `json:"added_at"`
}

// OrganizationOnMemberForQuery メンバー上のオーガナイゼーションを表す構造体(クエリ用)。
type OrganizationOnMemberForQuery struct {
	Pkey Int `json:"-"`
	OrganizationOnMember
}

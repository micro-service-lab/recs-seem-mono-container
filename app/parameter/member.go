package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateMemberParam メンバー作成のパラメータ。
type CreateMemberParam struct {
	LoginID                string
	Password               string
	Email                  string
	Name                   string
	FirstName              entity.String
	LastName               entity.String
	AttendStatusID         uuid.UUID
	GradeID                uuid.UUID
	GroupID                uuid.UUID
	ProfileImageID         entity.UUID
	RoleID                 entity.UUID
	PersonalOrganizationID uuid.UUID
}

// UpdateMemberParams メンバー更新のパラメータ。
type UpdateMemberParams struct {
	Email          string
	Name           string
	FirstName      entity.String
	LastName       entity.String
	ProfileImageID entity.UUID
}

// MemberWith オーガナイゼーションの付加情報。
type MemberWith struct {
	withDetail               bool
	withProfileImage         bool
	withPersonalOrganization bool
	withCrew                 bool
	withAttendStatus         bool
	withRole                 bool
}

// MemberWithParams オーガナイゼーションの付加情報。
type MemberWithParams []MemberWith

// ParseMemberWithParam オーガナイゼーションの付加情報をパースする。
func ParseMemberWithParam(v string) (any, error) {
	if v == "" {
		return MemberWith{}, nil
	}
	switch v {
	case "detail":
		return MemberWith{withDetail: true}, nil
	case "profile_image":
		return MemberWith{withProfileImage: true}, nil
	case "personal_organization":
		return MemberWith{withPersonalOrganization: true}, nil
	case "crew":
		return MemberWith{withCrew: true}, nil
	case "attend_status":
		return MemberWith{withAttendStatus: true}, nil
	case "role":
		return MemberWith{withRole: true}, nil
	default:
		return MemberWith{}, nil
	}
}

// MemberWithCase オーガナイゼーションの付加情報のケース。
type MemberWithCase int8

const (
	// MemberWithCaseDefault はデフォルト。
	MemberWithCaseDefault MemberWithCase = 0b0
	// MemberWithCaseDetail はカテゴリを含む。
	MemberWithCaseDetail MemberWithCase = 0b1
	// MemberWithCaseProfileImage はプロフィール画像を含む。
	MemberWithCaseProfileImage MemberWithCase = 0b10
	// MemberWithCasePersonalOrganization は個人組織を含む。
	MemberWithCasePersonalOrganization MemberWithCase = 0b100
	// MemberWithCaseCrew はクルーを含む。
	MemberWithCaseCrew MemberWithCase = 0b1000
	// MemberWithCaseAttendStatus は出席状況を含む。
	MemberWithCaseAttendStatus MemberWithCase = 0b10000
	// MemberWithCaseRole は役割を含む。
	MemberWithCaseRole MemberWithCase = 0b100000
)

// Case はケースを取得する。
func (p MemberWithParams) Case() MemberWithCase {
	var c MemberWithCase
	for _, v := range p {
		if v.withDetail {
			c |= MemberWithCaseDetail
		}
		if v.withProfileImage {
			c |= MemberWithCaseProfileImage
		}
		if v.withPersonalOrganization {
			c |= MemberWithCasePersonalOrganization
		}
		if v.withCrew {
			c |= MemberWithCaseCrew
		}
		if v.withAttendStatus {
			c |= MemberWithCaseAttendStatus
		}
		if v.withRole {
			c |= MemberWithCaseRole
		}
	}
	return c
}

// WhereMemberParam メンバー検索のパラメータ。
type WhereMemberParam struct {
	WhereLikeName      bool
	SearchName         string
	WhereHasPolicy     bool
	HasPolicyIDs       []uuid.UUID
	WhenInAttendStatus bool
	InAttendStatusIDs  []uuid.UUID
	WhenInGrade        bool
	InGradeIDs         []uuid.UUID
	WhenInGroup        bool
	InGroupIDs         []uuid.UUID
}

// MemberOrderMethod メンバーの並び替え方法。
type MemberOrderMethod string

// ParseMemberOrderMethod はメンバーの並び替え方法をパースする。
func ParseMemberOrderMethod(v string) (any, error) {
	if v == "" {
		return MemberOrderMethodDefault, nil
	}
	switch v {
	case string(MemberOrderMethodDefault):
		return MemberOrderMethodDefault, nil
	case string(MemberOrderMethodName):
		return MemberOrderMethodName, nil
	case string(MemberOrderMethodReverseName):
		return MemberOrderMethodReverseName, nil
	default:
		return MemberOrderMethodDefault, nil
	}
}

const (
	// MemberDefaultCursorKey はデフォルトカーソルキー。
	MemberDefaultCursorKey = "default"
	// MemberNameCursorKey は名前カーソルキー。
	MemberNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MemberOrderMethod) GetCursorKeyName() string {
	switch m {
	case MemberOrderMethodDefault:
		return MemberDefaultCursorKey
	case MemberOrderMethodName, MemberOrderMethodReverseName:
		return MemberNameCursorKey
	default:
		return MemberDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MemberOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MemberOrderMethodDefault はデフォルト。
	MemberOrderMethodDefault MemberOrderMethod = "default"
	// MemberOrderMethodName は名前順。
	MemberOrderMethodName MemberOrderMethod = "name"
	// MemberOrderMethodReverseName は名前逆順。
	MemberOrderMethodReverseName MemberOrderMethod = "r_name"
)

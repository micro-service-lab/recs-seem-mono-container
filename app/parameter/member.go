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
	FirstName              string
	LastName               string
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
	FirstName      string
	LastName       string
	ProfileImageID entity.UUID
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

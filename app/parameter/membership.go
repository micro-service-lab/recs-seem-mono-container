package parameter

import (
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// BelongOrganizationParam オーガナイゼーション所属のパラメータ。
type BelongOrganizationParam struct {
	MemberID       uuid.UUID
	OrganizationID uuid.UUID
	WorkPositionID entity.UUID
	AddedAt        time.Time
}

// WhereOrganizationOnMemberParam メンバー上のオーガナイゼーション検索のパラメータ。
type WhereOrganizationOnMemberParam struct {
	WhereLikeName bool
	SearchName    string
}

// OrganizationOnMemberOrderMethod メンバー上のオーガナイゼーションの並び替え方法。
type OrganizationOnMemberOrderMethod string

// ParseOrganizationOnMemberOrderMethod はメンバー上のオーガナイゼーションの並び替え方法をパースする。
func ParseOrganizationOnMemberOrderMethod(v string) (any, error) {
	if v == "" {
		return OrganizationOnMemberOrderMethodDefault, nil
	}
	switch v {
	case string(OrganizationOnMemberOrderMethodDefault):
		return OrganizationOnMemberOrderMethodDefault, nil
	case string(OrganizationOnMemberOrderMethodName):
		return OrganizationOnMemberOrderMethodName, nil
	case string(OrganizationOnMemberOrderMethodReverseName):
		return OrganizationOnMemberOrderMethodReverseName, nil
	case string(OrganizationOnMemberOrderMethodOldAdd):
		return OrganizationOnMemberOrderMethodOldAdd, nil
	case string(OrganizationOnMemberOrderMethodLateAdd):
		return OrganizationOnMemberOrderMethodLateAdd, nil
	default:
		return OrganizationOnMemberOrderMethodDefault, nil
	}
}

const (
	// OrganizationOnMemberDefaultCursorKey はデフォルトカーソルキー。
	OrganizationOnMemberDefaultCursorKey = "default"
	// OrganizationOnMemberNameCursorKey は名前カーソルキー。
	OrganizationOnMemberNameCursorKey = "name"
	// OrganizationOnMemberAddedAtCursorKey は追加日時カーソルキー。
	OrganizationOnMemberAddedAtCursorKey = "added_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m OrganizationOnMemberOrderMethod) GetCursorKeyName() string {
	switch m {
	case OrganizationOnMemberOrderMethodDefault:
		return OrganizationOnMemberDefaultCursorKey
	case OrganizationOnMemberOrderMethodName, OrganizationOnMemberOrderMethodReverseName:
		return OrganizationOnMemberNameCursorKey
	case OrganizationOnMemberOrderMethodOldAdd, OrganizationOnMemberOrderMethodLateAdd:
		return OrganizationOnMemberAddedAtCursorKey
	default:
		return OrganizationOnMemberDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m OrganizationOnMemberOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// OrganizationOnMemberOrderMethodDefault はデフォルト。
	OrganizationOnMemberOrderMethodDefault OrganizationOnMemberOrderMethod = "default"
	// OrganizationOnMemberOrderMethodName は名前順。
	OrganizationOnMemberOrderMethodName OrganizationOnMemberOrderMethod = "name"
	// OrganizationOnMemberOrderMethodReverseName は名前逆順。
	OrganizationOnMemberOrderMethodReverseName OrganizationOnMemberOrderMethod = "r_name"
	// OrganizationOnMemberOrderMethodOldAdd は追加古い順。
	OrganizationOnMemberOrderMethodOldAdd OrganizationOnMemberOrderMethod = "old_add"
	// OrganizationOnMemberOrderMethodLateAdd は追加新しい順。
	OrganizationOnMemberOrderMethodLateAdd OrganizationOnMemberOrderMethod = "late_add"
)

// WhereMemberOnOrganizationParam オーガナイゼーション上のメンバー検索のパラメータ。
type WhereMemberOnOrganizationParam struct {
	WhereLikeName bool
	SearchName    string
}

// MemberOnOrganizationOrderMethod オーガナイゼーション上のメンバーの並び替え方法。
type MemberOnOrganizationOrderMethod string

// ParseMemberOnOrganizationOrderMethod はオーガナイゼーション上のメンバーの並び替え方法をパースする。
func ParseMemberOnOrganizationOrderMethod(v string) (any, error) {
	if v == "" {
		return MemberOnOrganizationOrderMethodDefault, nil
	}
	switch v {
	case string(MemberOnOrganizationOrderMethodDefault):
		return MemberOnOrganizationOrderMethodDefault, nil
	case string(MemberOnOrganizationOrderMethodName):
		return MemberOnOrganizationOrderMethodName, nil
	case string(MemberOnOrganizationOrderMethodReverseName):
		return MemberOnOrganizationOrderMethodReverseName, nil
	case string(MemberOnOrganizationOrderMethodOldAdd):
		return MemberOnOrganizationOrderMethodOldAdd, nil
	case string(MemberOnOrganizationOrderMethodLateAdd):
		return MemberOnOrganizationOrderMethodLateAdd, nil
	default:
		return MemberOnOrganizationOrderMethodDefault, nil
	}
}

const (
	// MemberOnOrganizationDefaultCursorKey はデフォルトカーソルキー。
	MemberOnOrganizationDefaultCursorKey = "default"
	// MemberOnOrganizationNameCursorKey は名前カーソルキー。
	MemberOnOrganizationNameCursorKey = "name"
	// MemberOnOrganizationAddedAtCursorKey は追加日時カーソルキー。
	MemberOnOrganizationAddedAtCursorKey = "added_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MemberOnOrganizationOrderMethod) GetCursorKeyName() string {
	switch m {
	case MemberOnOrganizationOrderMethodDefault:
		return MemberOnOrganizationDefaultCursorKey
	case MemberOnOrganizationOrderMethodName, MemberOnOrganizationOrderMethodReverseName:
		return MemberOnOrganizationNameCursorKey
	case MemberOnOrganizationOrderMethodOldAdd, MemberOnOrganizationOrderMethodLateAdd:
		return MemberOnOrganizationAddedAtCursorKey
	default:
		return MemberOnOrganizationDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MemberOnOrganizationOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MemberOnOrganizationOrderMethodDefault はデフォルト。
	MemberOnOrganizationOrderMethodDefault MemberOnOrganizationOrderMethod = "default"
	// MemberOnOrganizationOrderMethodName は名前順。
	MemberOnOrganizationOrderMethodName MemberOnOrganizationOrderMethod = "name"
	// MemberOnOrganizationOrderMethodReverseName は名前逆順。
	MemberOnOrganizationOrderMethodReverseName MemberOnOrganizationOrderMethod = "r_name"
	// MemberOnOrganizationOrderMethodOldAdd は追加古い順。
	MemberOnOrganizationOrderMethodOldAdd MemberOnOrganizationOrderMethod = "old_add"
	// MemberOnOrganizationOrderMethodLateAdd は追加新しい順。
	MemberOnOrganizationOrderMethodLateAdd MemberOnOrganizationOrderMethod = "late_add"
)

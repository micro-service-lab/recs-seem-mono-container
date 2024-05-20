package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateOrganizationParam オーガナイゼーション作成のパラメータ。
type CreateOrganizationParam struct {
	Name        string
	Description entity.String
	Color       entity.String
	IsPersonal  bool
	IsWhole     bool
	ChatRoomID  entity.UUID
}

// UpdateOrganizationParams オーガナイゼーション更新のパラメータ。
type UpdateOrganizationParams struct {
	Name        string
	Description entity.String
	Color       entity.String
}

// WhereOrganizationParam オーガナイゼーション検索のパラメータ。
type WhereOrganizationParam struct {
	WhereLikeName    bool
	SearchName       string
	WhereIsWhole     bool
	IsWhole          bool
	WhereIsPersonal  bool
	IsPersonal       bool
	PersonalMemberID uuid.UUID
	WhereIsGroup     bool
	WhereIsGrade     bool
}

// OrganizationOrderMethod オーガナイゼーションの並び替え方法。
type OrganizationOrderMethod string

// ParseOrganizationOrderMethod はオーガナイゼーションの並び替え方法をパースする。
func ParseOrganizationOrderMethod(v string) (any, error) {
	if v == "" {
		return OrganizationOrderMethodDefault, nil
	}
	switch v {
	case string(OrganizationOrderMethodDefault):
		return OrganizationOrderMethodDefault, nil
	case string(OrganizationOrderMethodName):
		return OrganizationOrderMethodName, nil
	case string(OrganizationOrderMethodReverseName):
		return OrganizationOrderMethodReverseName, nil
	default:
		return OrganizationOrderMethodDefault, nil
	}
}

const (
	// OrganizationDefaultCursorKey はデフォルトカーソルキー。
	OrganizationDefaultCursorKey = "default"
	// OrganizationNameCursorKey は名前カーソルキー。
	OrganizationNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m OrganizationOrderMethod) GetCursorKeyName() string {
	switch m {
	case OrganizationOrderMethodDefault:
		return OrganizationDefaultCursorKey
	case OrganizationOrderMethodName, OrganizationOrderMethodReverseName:
		return OrganizationNameCursorKey
	default:
		return OrganizationDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m OrganizationOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// OrganizationOrderMethodDefault はデフォルト。
	OrganizationOrderMethodDefault OrganizationOrderMethod = "default"
	// OrganizationOrderMethodName は名前順。
	OrganizationOrderMethodName OrganizationOrderMethod = "name"
	// OrganizationOrderMethodReverseName は名前逆順。
	OrganizationOrderMethodReverseName OrganizationOrderMethod = "r_name"
)

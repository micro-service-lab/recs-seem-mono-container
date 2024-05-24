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

// OrganizationWith オーガナイゼーションの付加情報。
type OrganizationWith struct {
	withDetail   bool
	withChatRoom bool
}

// OrganizationWithParams オーガナイゼーションの付加情報。
type OrganizationWithParams []OrganizationWith

// ParseOrganizationWithParam オーガナイゼーションの付加情報をパースする。
func ParseOrganizationWithParam(v string) (any, error) {
	if v == "" {
		return OrganizationWith{}, nil
	}
	switch v {
	case "detail":
		return OrganizationWith{withDetail: true}, nil
	case "chat_room":
		return OrganizationWith{withChatRoom: true}, nil
	default:
		return OrganizationWith{}, nil
	}
}

// OrganizationWithCase オーガナイゼーションの付加情報のケース。
type OrganizationWithCase int8

const (
	// OrganizationWithCaseChatRoomAndDetail はチャットルームとカテゴリを含む。
	OrganizationWithCaseChatRoomAndDetail OrganizationWithCase = 0b11
	// OrganizationWithCaseChatRoom はチャットルームを含む。
	OrganizationWithCaseChatRoom OrganizationWithCase = 0b10
	// OrganizationWithCaseDetail はカテゴリを含む。
	OrganizationWithCaseDetail OrganizationWithCase = 0b1
	// OrganizationWithCaseDefault はデフォルト。
	OrganizationWithCaseDefault OrganizationWithCase = 0b0
)

// Case はケースを取得する。
func (p OrganizationWithParams) Case() OrganizationWithCase {
	var c OrganizationWithCase
	for _, v := range p {
		if v.withDetail {
			c |= OrganizationWithCaseDetail
		}
		if v.withChatRoom {
			c |= OrganizationWithCaseChatRoom
		}
	}
	return c
}

// WhereOrganizationType オーガナイゼーション検索のタイプ。
type WhereOrganizationType string

const (
	// WhereOrganizationTypeDefault デフォルト。
	WhereOrganizationTypeDefault WhereOrganizationType = "default"
	// WhereOrganizationTypePersonal パーソナル。
	WhereOrganizationTypePersonal WhereOrganizationType = "personal"
	// WhereOrganizationTypeWhole 全体。
	WhereOrganizationTypeWhole WhereOrganizationType = "whole"
	// WhereOrganizationTypeGroup グループ。
	WhereOrganizationTypeGroup WhereOrganizationType = "group"
	// WhereOrganizationTypeGrade 学年。
	WhereOrganizationTypeGrade WhereOrganizationType = "grade"
)

// ParseWhereOrganizationType はオーガナイゼーション検索のタイプをパースする。
func ParseWhereOrganizationType(v string) (any, error) {
	if v == "" {
		return WhereOrganizationTypeDefault, nil
	}
	switch v {
	case string(WhereOrganizationTypeDefault):
		return WhereOrganizationTypeDefault, nil
	case string(WhereOrganizationTypePersonal):
		return WhereOrganizationTypePersonal, nil
	case string(WhereOrganizationTypeWhole):
		return WhereOrganizationTypeWhole, nil
	case string(WhereOrganizationTypeGroup):
		return WhereOrganizationTypeGroup, nil
	case string(WhereOrganizationTypeGrade):
		return WhereOrganizationTypeGrade, nil
	default:
		return WhereOrganizationTypeDefault, nil
	}
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

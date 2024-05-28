package parameter

import (
	"github.com/google/uuid"
	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateGroupServiceParam グループ作成のパラメータ。
type CreateGroupServiceParam struct {
	Name         string
	Key          string
	Description  entity.String
	Color        entity.String
	CoverImageID entity.UUID
}

// CreateGroupParam グループ作成のパラメータ。
type CreateGroupParam struct {
	Key            string
	OrganizationID uuid.UUID
}

// WhereGroupParam グループ検索のパラメータ。
type WhereGroupParam struct{}

// GroupOrderMethod グループの並び替え方法。
type GroupOrderMethod string

// ParseGroupOrderMethod はグループの並び替え方法をパースする。
func ParseGroupOrderMethod(v string) (any, error) {
	if v == "" {
		return GroupOrderMethodDefault, nil
	}
	switch v {
	case string(GroupOrderMethodDefault):
		return GroupOrderMethodDefault, nil
	case string(GroupOrderMethodName):
		return GroupOrderMethodName, nil
	case string(GroupOrderMethodReverseName):
		return GroupOrderMethodReverseName, nil
	default:
		return GroupOrderMethodDefault, nil
	}
}

const (
	// GroupDefaultCursorKey はデフォルトカーソルキー。
	GroupDefaultCursorKey = "default"
	// GroupNameCursorKey は名前カーソルキー。
	GroupNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m GroupOrderMethod) GetCursorKeyName() string {
	switch m {
	case GroupOrderMethodDefault:
		return GroupDefaultCursorKey
	case GroupOrderMethodName, GroupOrderMethodReverseName:
		return GroupNameCursorKey
	default:
		return GroupDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m GroupOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// GroupOrderMethodDefault はデフォルト。
	GroupOrderMethodDefault GroupOrderMethod = "default"
	// GroupOrderMethodName は名前順。
	GroupOrderMethodName GroupOrderMethod = "name"
	// GroupOrderMethodReverseName は名前逆順。
	GroupOrderMethodReverseName GroupOrderMethod = "r_name"
)

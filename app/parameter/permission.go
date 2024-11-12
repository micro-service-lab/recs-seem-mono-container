package parameter

import "github.com/google/uuid"

// CreatePermissionParam 権限作成のパラメータ。
type CreatePermissionParam struct {
	Name                 string
	Key                  string
	Description          string
	PermissionCategoryID uuid.UUID
}

// UpdatePermissionParams 権限更新のパラメータ。
type UpdatePermissionParams struct {
	Name                 string
	Key                  string
	Description          string
	PermissionCategoryID uuid.UUID
}

// UpdatePermissionByKeyParams 権限更新のパラメータ。
type UpdatePermissionByKeyParams struct {
	Name                 string
	Description          string
	PermissionCategoryID uuid.UUID
}

// WherePermissionParam 権限検索のパラメータ。
type WherePermissionParam struct {
	WhereLikeName   bool
	SearchName      string
	WhereInCategory bool
	InCategories    []uuid.UUID
}

// PermissionWith 権限の付加情報。
type PermissionWith struct {
	withCategory bool
}

// PermissionWithParams 権限の付加情報。
type PermissionWithParams []PermissionWith

// ParsePermissionWithParam 権限の付加情報をパースする。
func ParsePermissionWithParam(v string) (any, error) {
	if v == "" {
		return PermissionWith{}, nil
	}
	switch v {
	case "category":
		return PermissionWith{withCategory: true}, nil
	default:
		return PermissionWith{}, nil
	}
}

// PermissionWithCase 権限の付加情報のケース。
type PermissionWithCase int8

const (
	// PermissionWithCaseCategory はカテゴリを含む。
	PermissionWithCaseCategory PermissionWithCase = 0b1
	// PermissionWithCaseDefault はデフォルト。
	PermissionWithCaseDefault PermissionWithCase = 0b0
)

// Case はケースを取得する。
func (p PermissionWithParams) Case() PermissionWithCase {
	var c PermissionWithCase
	for _, v := range p {
		if v.withCategory {
			c |= PermissionWithCaseCategory
		}
	}
	return c
}

// PermissionOrderMethod 権限の並び替え方法。
type PermissionOrderMethod string

// ParsePermissionOrderMethod は権限の並び替え方法をパースする。
func ParsePermissionOrderMethod(v string) (any, error) {
	if v == "" {
		return PermissionOrderMethodDefault, nil
	}
	switch v {
	case string(PermissionOrderMethodDefault):
		return PermissionOrderMethodDefault, nil
	case string(PermissionOrderMethodName):
		return PermissionOrderMethodName, nil
	case string(PermissionOrderMethodReverseName):
		return PermissionOrderMethodReverseName, nil
	default:
		return PermissionOrderMethodDefault, nil
	}
}

const (
	// PermissionDefaultCursorKey はデフォルトカーソルキー。
	PermissionDefaultCursorKey = "default"
	// PermissionNameCursorKey は名前カーソルキー。
	PermissionNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m PermissionOrderMethod) GetCursorKeyName() string {
	switch m {
	case PermissionOrderMethodDefault:
		return PermissionDefaultCursorKey
	case PermissionOrderMethodName:
		return PermissionNameCursorKey
	case PermissionOrderMethodReverseName:
		return PermissionNameCursorKey
	default:
		return PermissionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m PermissionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// PermissionOrderMethodDefault はデフォルト。
	PermissionOrderMethodDefault PermissionOrderMethod = "default"
	// PermissionOrderMethodName は名前順。
	PermissionOrderMethodName PermissionOrderMethod = "name"
	// PermissionOrderMethodReverseName は名前逆順。
	PermissionOrderMethodReverseName PermissionOrderMethod = "r_name"
)

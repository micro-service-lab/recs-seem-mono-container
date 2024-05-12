package parameter

import "github.com/google/uuid"

// CreatePolicyParam ポリシー作成のパラメータ。
type CreatePolicyParam struct {
	Name             string
	Key              string
	Description      string
	PolicyCategoryID uuid.UUID
}

// UpdatePolicyParams ポリシー更新のパラメータ。
type UpdatePolicyParams struct {
	Name             string
	Key              string
	Description      string
	PolicyCategoryID uuid.UUID
}

// UpdatePolicyByKeyParams ポリシー更新のパラメータ。
type UpdatePolicyByKeyParams struct {
	Name             string
	Description      string
	PolicyCategoryID uuid.UUID
}

// WherePolicyParam ポリシー検索のパラメータ。
type WherePolicyParam struct {
	WhereLikeName   bool
	SearchName      string
	WhereInCategory bool
	InCategories    []uuid.UUID
}

// PolicyWith ポリシーの付加情報。
type PolicyWith struct {
	withCategory bool
}

// PolicyWithParams ポリシーの付加情報。
type PolicyWithParams []PolicyWith

// ParsePolicyWithParam ポリシーの付加情報をパースする。
func ParsePolicyWithParam(v string) (any, error) {
	if v == "" {
		return PolicyWith{}, nil
	}
	switch v {
	case "category":
		return PolicyWith{withCategory: true}, nil
	default:
		return PolicyWith{}, nil
	}
}

// PolicyWithCase ポリシーの付加情報のケース。
type PolicyWithCase int8

const (
	// PolicyWithCaseCategory はカテゴリを含む。
	PolicyWithCaseCategory PolicyWithCase = 0b1
	// PolicyWithCaseDefault はデフォルト。
	PolicyWithCaseDefault PolicyWithCase = 0b0
)

// Case はケースを取得する。
func (p PolicyWithParams) Case() PolicyWithCase {
	var c PolicyWithCase
	for _, v := range p {
		if v.withCategory {
			c |= PolicyWithCaseCategory
		}
	}
	return c
}

// PolicyOrderMethod ポリシーの並び替え方法。
type PolicyOrderMethod string

// ParsePolicyOrderMethod はポリシーの並び替え方法をパースする。
func ParsePolicyOrderMethod(v string) (any, error) {
	if v == "" {
		return PolicyOrderMethodDefault, nil
	}
	switch v {
	case string(PolicyOrderMethodDefault):
		return PolicyOrderMethodDefault, nil
	case string(PolicyOrderMethodName):
		return PolicyOrderMethodName, nil
	case string(PolicyOrderMethodReverseName):
		return PolicyOrderMethodReverseName, nil
	default:
		return PolicyOrderMethodDefault, nil
	}
}

const (
	// PolicyDefaultCursorKey はデフォルトカーソルキー。
	PolicyDefaultCursorKey = "default"
	// PolicyNameCursorKey は名前カーソルキー。
	PolicyNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m PolicyOrderMethod) GetCursorKeyName() string {
	switch m {
	case PolicyOrderMethodDefault:
		return PolicyDefaultCursorKey
	case PolicyOrderMethodName:
		return PolicyNameCursorKey
	case PolicyOrderMethodReverseName:
		return PolicyNameCursorKey
	default:
		return PolicyDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m PolicyOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// PolicyOrderMethodDefault はデフォルト。
	PolicyOrderMethodDefault PolicyOrderMethod = "default"
	// PolicyOrderMethodName は名前順。
	PolicyOrderMethodName PolicyOrderMethod = "name"
	// PolicyOrderMethodReverseName は名前逆順。
	PolicyOrderMethodReverseName PolicyOrderMethod = "r_name"
)

package parameter

// CreatePolicyCategoryParam ポリシーカテゴリー作成のパラメータ。
type CreatePolicyCategoryParam struct {
	Name        string
	Key         string
	Description string
}

// UpdatePolicyCategoryParams ポリシーカテゴリー更新のパラメータ。
type UpdatePolicyCategoryParams struct {
	Name        string
	Key         string
	Description string
}

// UpdatePolicyCategoryByKeyParams ポリシーカテゴリー更新のパラメータ。
type UpdatePolicyCategoryByKeyParams struct {
	Name        string
	Description string
}

// WherePolicyCategoryParam ポリシーカテゴリー検索のパラメータ。
type WherePolicyCategoryParam struct {
	WhereLikeName bool
	SearchName    string
}

// PolicyCategoryOrderMethod ポリシーカテゴリーの並び替え方法。
type PolicyCategoryOrderMethod string

// ParsePolicyCategoryOrderMethod はポリシーカテゴリーの並び替え方法をパースする。
func ParsePolicyCategoryOrderMethod(v string) (any, error) {
	if v == "" {
		return PolicyCategoryOrderMethodDefault, nil
	}
	switch v {
	case string(PolicyCategoryOrderMethodDefault):
		return PolicyCategoryOrderMethodDefault, nil
	case string(PolicyCategoryOrderMethodName):
		return PolicyCategoryOrderMethodName, nil
	case string(PolicyCategoryOrderMethodReverseName):
		return PolicyCategoryOrderMethodReverseName, nil
	default:
		return PolicyCategoryOrderMethodDefault, nil
	}
}

const (
	// PolicyCategoryDefaultCursorKey はデフォルトカーソルキー。
	PolicyCategoryDefaultCursorKey = "default"
	// PolicyCategoryNameCursorKey は名前カーソルキー。
	PolicyCategoryNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m PolicyCategoryOrderMethod) GetCursorKeyName() string {
	switch m {
	case PolicyCategoryOrderMethodDefault:
		return PolicyCategoryDefaultCursorKey
	case PolicyCategoryOrderMethodName:
		return PolicyCategoryNameCursorKey
	case PolicyCategoryOrderMethodReverseName:
		return PolicyCategoryNameCursorKey
	default:
		return PolicyCategoryDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m PolicyCategoryOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// PolicyCategoryOrderMethodDefault はデフォルト。
	PolicyCategoryOrderMethodDefault PolicyCategoryOrderMethod = "default"
	// PolicyCategoryOrderMethodName は名前順。
	PolicyCategoryOrderMethodName PolicyCategoryOrderMethod = "name"
	// PolicyCategoryOrderMethodReverseName は名前逆順。
	PolicyCategoryOrderMethodReverseName PolicyCategoryOrderMethod = "r_name"
)

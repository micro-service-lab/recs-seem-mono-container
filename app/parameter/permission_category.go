package parameter

// CreatePermissionCategoryParam 権限カテゴリー作成のパラメータ。
type CreatePermissionCategoryParam struct {
	Name        string
	Key         string
	Description string
}

// UpdatePermissionCategoryParams 権限カテゴリー更新のパラメータ。
type UpdatePermissionCategoryParams struct {
	Name        string
	Key         string
	Description string
}

// UpdatePermissionCategoryByKeyParams 権限カテゴリー更新のパラメータ。
type UpdatePermissionCategoryByKeyParams struct {
	Name        string
	Description string
}

// WherePermissionCategoryParam 権限カテゴリー検索のパラメータ。
type WherePermissionCategoryParam struct {
	WhereLikeName bool
	SearchName    string
}

// PermissionCategoryOrderMethod 権限カテゴリーの並び替え方法。
type PermissionCategoryOrderMethod string

// ParsePermissionCategoryOrderMethod は権限カテゴリーの並び替え方法をパースする。
func ParsePermissionCategoryOrderMethod(v string) (any, error) {
	if v == "" {
		return PermissionCategoryOrderMethodDefault, nil
	}
	switch v {
	case string(PermissionCategoryOrderMethodDefault):
		return PermissionCategoryOrderMethodDefault, nil
	case string(PermissionCategoryOrderMethodName):
		return PermissionCategoryOrderMethodName, nil
	case string(PermissionCategoryOrderMethodReverseName):
		return PermissionCategoryOrderMethodReverseName, nil
	default:
		return PermissionCategoryOrderMethodDefault, nil
	}
}

const (
	// PermissionCategoryDefaultCursorKey はデフォルトカーソルキー。
	PermissionCategoryDefaultCursorKey = "default"
	// PermissionCategoryNameCursorKey は名前カーソルキー。
	PermissionCategoryNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m PermissionCategoryOrderMethod) GetCursorKeyName() string {
	switch m {
	case PermissionCategoryOrderMethodDefault:
		return PermissionCategoryDefaultCursorKey
	case PermissionCategoryOrderMethodName:
		return PermissionCategoryNameCursorKey
	case PermissionCategoryOrderMethodReverseName:
		return PermissionCategoryNameCursorKey
	default:
		return PermissionCategoryDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m PermissionCategoryOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// PermissionCategoryOrderMethodDefault はデフォルト。
	PermissionCategoryOrderMethodDefault PermissionCategoryOrderMethod = "default"
	// PermissionCategoryOrderMethodName は名前順。
	PermissionCategoryOrderMethodName PermissionCategoryOrderMethod = "name"
	// PermissionCategoryOrderMethodReverseName は名前逆順。
	PermissionCategoryOrderMethodReverseName PermissionCategoryOrderMethod = "r_name"
)

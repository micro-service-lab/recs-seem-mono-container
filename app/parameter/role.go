package parameter

// CreateRoleParam ロール作成のパラメータ。
type CreateRoleParam struct {
	Name        string
	Description string
}

// UpdateRoleParams ロール更新のパラメータ。
type UpdateRoleParams struct {
	Name        string
	Description string
}

// WhereRoleParam ロール検索のパラメータ。
type WhereRoleParam struct {
	WhereLikeName bool
	SearchName    string
}

// RoleOrderMethod ロールの並び替え方法。
type RoleOrderMethod string

// ParseRoleOrderMethod はロールの並び替え方法をパースする。
func ParseRoleOrderMethod(v string) (any, error) {
	if v == "" {
		return RoleOrderMethodDefault, nil
	}
	switch v {
	case string(RoleOrderMethodDefault):
		return RoleOrderMethodDefault, nil
	case string(RoleOrderMethodName):
		return RoleOrderMethodName, nil
	case string(RoleOrderMethodReverseName):
		return RoleOrderMethodReverseName, nil
	default:
		return RoleOrderMethodDefault, nil
	}
}

const (
	// RoleDefaultCursorKey はデフォルトカーソルキー。
	RoleDefaultCursorKey = "default"
	// RoleNameCursorKey は名前カーソルキー。
	RoleNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m RoleOrderMethod) GetCursorKeyName() string {
	switch m {
	case RoleOrderMethodDefault:
		return RoleDefaultCursorKey
	case RoleOrderMethodName:
		return RoleNameCursorKey
	case RoleOrderMethodReverseName:
		return RoleNameCursorKey
	default:
		return RoleDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m RoleOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// RoleOrderMethodDefault はデフォルト。
	RoleOrderMethodDefault RoleOrderMethod = "default"
	// RoleOrderMethodName は名前順。
	RoleOrderMethodName RoleOrderMethod = "name"
	// RoleOrderMethodReverseName は名前逆順。
	RoleOrderMethodReverseName RoleOrderMethod = "r_name"
)

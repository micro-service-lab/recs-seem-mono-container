package parameter

import "github.com/google/uuid"

// AssociationRoleParam ロール関連付けのパラメータ。
type AssociationRoleParam struct {
	RoleID   uuid.UUID
	PolicyID uuid.UUID
}

// WhereRoleOnPolicyParam ポリシー上ロール検索のパラメータ。
type WhereRoleOnPolicyParam struct {
	WhereLikeName bool
	SearchName    string
}

// RoleOnPolicyOrderMethod ポリシー上ロールの並び替え方法。
type RoleOnPolicyOrderMethod string

// ParseRoleOnPolicyOrderMethod はポリシー上ロールの並び替え方法をパースする。
func ParseRoleOnPolicyOrderMethod(v string) (any, error) {
	if v == "" {
		return RoleOnPolicyOrderMethodDefault, nil
	}
	switch v {
	case string(RoleOnPolicyOrderMethodDefault):
		return RoleOnPolicyOrderMethodDefault, nil
	case string(RoleOnPolicyOrderMethodName):
		return RoleOnPolicyOrderMethodName, nil
	case string(RoleOnPolicyOrderMethodReverseName):
		return RoleOnPolicyOrderMethodReverseName, nil
	default:
		return RoleOnPolicyOrderMethodDefault, nil
	}
}

const (
	// RoleOnPolicyDefaultCursorKey はデフォルトカーソルキー。
	RoleOnPolicyDefaultCursorKey = "default"
	// RoleOnPolicyNameCursorKey は名前カーソルキー。
	RoleOnPolicyNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m RoleOnPolicyOrderMethod) GetCursorKeyName() string {
	switch m {
	case RoleOnPolicyOrderMethodDefault:
		return RoleOnPolicyDefaultCursorKey
	case RoleOnPolicyOrderMethodName:
		return RoleOnPolicyNameCursorKey
	case RoleOnPolicyOrderMethodReverseName:
		return RoleOnPolicyNameCursorKey
	default:
		return RoleOnPolicyDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m RoleOnPolicyOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// RoleOnPolicyOrderMethodDefault はデフォルト。
	RoleOnPolicyOrderMethodDefault RoleOnPolicyOrderMethod = "default"
	// RoleOnPolicyOrderMethodName は名前順。
	RoleOnPolicyOrderMethodName RoleOnPolicyOrderMethod = "name"
	// RoleOnPolicyOrderMethodReverseName は名前逆順。
	RoleOnPolicyOrderMethodReverseName RoleOnPolicyOrderMethod = "r_name"
)

// WherePolicyOnRoleParam ロール上のポリシー検索のパラメータ。
type WherePolicyOnRoleParam struct {
	WhereLikeName bool
	SearchName    string
}

// PolicyOnRoleOrderMethod ロール上のポリシーの並び替え方法。
type PolicyOnRoleOrderMethod string

// ParsePolicyOnRoleOrderMethod はロール上のポリシーの並び替え方法をパースする。
func ParsePolicyOnRoleOrderMethod(v string) (any, error) {
	if v == "" {
		return PolicyOnRoleOrderMethodDefault, nil
	}
	switch v {
	case string(PolicyOnRoleOrderMethodDefault):
		return PolicyOnRoleOrderMethodDefault, nil
	case string(PolicyOnRoleOrderMethodName):
		return PolicyOnRoleOrderMethodName, nil
	case string(PolicyOnRoleOrderMethodReverseName):
		return PolicyOnRoleOrderMethodReverseName, nil
	default:
		return PolicyOnRoleOrderMethodDefault, nil
	}
}

const (
	// PolicyOnRoleDefaultCursorKey はデフォルトカーソルキー。
	PolicyOnRoleDefaultCursorKey = "default"
	// PolicyOnRoleNameCursorKey は名前カーソルキー。
	PolicyOnRoleNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m PolicyOnRoleOrderMethod) GetCursorKeyName() string {
	switch m {
	case PolicyOnRoleOrderMethodDefault:
		return PolicyOnRoleDefaultCursorKey
	case PolicyOnRoleOrderMethodName:
		return PolicyOnRoleNameCursorKey
	case PolicyOnRoleOrderMethodReverseName:
		return PolicyOnRoleNameCursorKey
	default:
		return PolicyOnRoleDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m PolicyOnRoleOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// PolicyOnRoleOrderMethodDefault はデフォルト。
	PolicyOnRoleOrderMethodDefault PolicyOnRoleOrderMethod = "default"
	// PolicyOnRoleOrderMethodName は名前順。
	PolicyOnRoleOrderMethodName PolicyOnRoleOrderMethod = "name"
	// PolicyOnRoleOrderMethodReverseName は名前逆順。
	PolicyOnRoleOrderMethodReverseName PolicyOnRoleOrderMethod = "r_name"
)

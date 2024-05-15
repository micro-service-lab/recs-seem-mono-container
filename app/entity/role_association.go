package entity

import "github.com/google/uuid"

// RoleAssociation ロールの関連付けを表すインターフェース。
type RoleAssociation struct {
	RoleID   uuid.UUID `json:"role_id"`
	PolicyID uuid.UUID `json:"policy_id"`
}

// RoleOnPolicy ポリシー上のロールを表す構造体。
type RoleOnPolicy struct {
	Role Role `json:"role"`
}

// RoleOnPolicyForQuery ポリシー上のロールを表す構造体(クエリ用)。
type RoleOnPolicyForQuery struct {
	Pkey Int `json:"-"`
	RoleOnPolicy
}

// PolicyOnRole ロール上のポリシーを表す構造体。
type PolicyOnRole struct {
	Policy Policy `json:"policy"`
}

// PolicyOnRoleForQuery ロール上のポリシーを表す構造体(クエリ用)。
type PolicyOnRoleForQuery struct {
	Pkey Int `json:"-"`
	PolicyOnRole
}

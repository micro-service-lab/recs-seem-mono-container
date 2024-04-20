package session

import "strings"

// ClientType クライアントの種別を表す。
type ClientType int

const (
	// ClientTypeInvalid 不明なクライアント種別。
	ClientTypeInvalid ClientType = iota
	// ClientTypeAdmin 管理者アカウントを表すクライアント種別。
	ClientTypeAdmin
	// ClientTypeUser ゲームのユーザーを表すクライアント種別。
	ClientTypeUser
)

// ClientTypeFromString 文字列を ClientType に変換する。
func ClientTypeFromString(s string) ClientType {
	switch strings.ToLower(s) {
	case "admin":
		return ClientTypeAdmin
	case "user":
		return ClientTypeUser
	default:
		return ClientTypeInvalid
	}
}

// String クライアント種別を表す文字列を返す。
func (s ClientType) String() string {
	switch s {
	case ClientTypeInvalid:
		return "invalid"
	case ClientTypeAdmin:
		return "admin"
	case ClientTypeUser:
		return "user"
	default:
		return "invalid"
	}
}

// IsValid 正しいクライアント種別かどうかを判定する。
func (s ClientType) IsValid() bool {
	if s == ClientTypeAdmin || s == ClientTypeUser {
		return true
	}
	return false
}

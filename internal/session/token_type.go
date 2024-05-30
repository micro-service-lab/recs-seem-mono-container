package session

import "strings"

// TokenType トークンの種別を表す。
type TokenType int

const (
	// TokenTypeInvalid 不明なトークン種別。
	TokenTypeInvalid TokenType = iota
	// TokenTypeAccess アクセスを表すトークン種別。
	TokenTypeAccess
	// TokenTypeRefresh リフレッシュを表すトークン種別。
	TokenTypeRefresh
)

// TokenTypeFromString 文字列を TokenType に変換する。
func TokenTypeFromString(s string) TokenType {
	switch strings.ToLower(s) {
	case "student":
		return TokenTypeAccess
	case "professor":
		return TokenTypeRefresh
	default:
		return TokenTypeInvalid
	}
}

// String トークン種別を表す文字列を返す。
func (s TokenType) String() string {
	switch s {
	case TokenTypeInvalid:
		return "invalid"
	case TokenTypeAccess:
		return "student"
	case TokenTypeRefresh:
		return "professor"
	default:
		return "invalid"
	}
}

// IsValid 正しいトークン種別かどうかを判定する。
func (s TokenType) IsValid() bool {
	if s == TokenTypeAccess || s == TokenTypeRefresh {
		return true
	}
	return false
}

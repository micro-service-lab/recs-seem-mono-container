package session

import "strings"

// MemberType メンバーの種別を表す。
type MemberType int

const (
	// MemberTypeInvalid 不明なメンバー種別。
	MemberTypeInvalid MemberType = iota
	// MemberTypeStudent 生徒を表すメンバー種別。
	MemberTypeStudent
	// MemberTypeProfessor 教授を表すメンバー種別。
	MemberTypeProfessor
)

// MemberTypeFromString 文字列を MemberType に変換する。
func MemberTypeFromString(s string) MemberType {
	switch strings.ToLower(s) {
	case "student":
		return MemberTypeStudent
	case "professor":
		return MemberTypeProfessor
	default:
		return MemberTypeInvalid
	}
}

// String メンバー種別を表す文字列を返す。
func (s MemberType) String() string {
	switch s {
	case MemberTypeInvalid:
		return "invalid"
	case MemberTypeStudent:
		return "student"
	case MemberTypeProfessor:
		return "professor"
	default:
		return "invalid"
	}
}

// IsValid 正しいメンバー種別かどうかを判定する。
func (s MemberType) IsValid() bool {
	if s == MemberTypeStudent || s == MemberTypeProfessor {
		return true
	}
	return false
}

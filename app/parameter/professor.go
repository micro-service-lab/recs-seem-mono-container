package parameter

import (
	"github.com/google/uuid"
)

// CreateProfessorParam 教授作成のパラメータ。
type CreateProfessorParam struct {
	MemberID uuid.UUID
}

// WhereProfessorParam 教授検索のパラメータ。
type WhereProfessorParam struct{}

// ProfessorOrderMethod 教授の並び替え方法。
type ProfessorOrderMethod string

// ParseProfessorOrderMethod は教授の並び替え方法をパースする。
func ParseProfessorOrderMethod(v string) (any, error) {
	if v == "" {
		return ProfessorOrderMethodDefault, nil
	}
	switch v {
	case string(ProfessorOrderMethodDefault):
		return ProfessorOrderMethodDefault, nil
	default:
		return ProfessorOrderMethodDefault, nil
	}
}

const (
	// ProfessorDefaultCursorKey はデフォルトカーソルキー。
	ProfessorDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ProfessorOrderMethod) GetCursorKeyName() string {
	switch m {
	case ProfessorOrderMethodDefault:
		return ProfessorDefaultCursorKey
	default:
		return ProfessorDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ProfessorOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ProfessorOrderMethodDefault はデフォルト。
	ProfessorOrderMethodDefault ProfessorOrderMethod = "default"
)

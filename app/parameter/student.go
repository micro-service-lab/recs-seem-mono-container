package parameter

import (
	"github.com/google/uuid"
)

// CreateStudentParam 生徒作成のパラメータ。
type CreateStudentParam struct {
	MemberID uuid.UUID
}

// WhereStudentParam 生徒検索のパラメータ。
type WhereStudentParam struct{}

// StudentOrderMethod 生徒の並び替え方法。
type StudentOrderMethod string

// ParseStudentOrderMethod は生徒の並び替え方法をパースする。
func ParseStudentOrderMethod(v string) (any, error) {
	if v == "" {
		return StudentOrderMethodDefault, nil
	}
	switch v {
	case string(StudentOrderMethodDefault):
		return StudentOrderMethodDefault, nil
	default:
		return StudentOrderMethodDefault, nil
	}
}

const (
	// StudentDefaultCursorKey はデフォルトカーソルキー。
	StudentDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m StudentOrderMethod) GetCursorKeyName() string {
	switch m {
	case StudentOrderMethodDefault:
		return StudentDefaultCursorKey
	default:
		return StudentDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m StudentOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// StudentOrderMethodDefault はデフォルト。
	StudentOrderMethodDefault StudentOrderMethod = "default"
)

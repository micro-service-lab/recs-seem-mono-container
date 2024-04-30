package parameter

import "github.com/google/uuid"

// CreateAbsenceParam 欠席作成のパラメータ。
type CreateAbsenceParam struct {
	AttendanceID uuid.UUID
}

// AbsenceOrderMethod 欠席の並び替え方法。
type AbsenceOrderMethod string

// ParseAbsenceOrderMethod は欠席の並び替え方法をパースする。
func ParseAbsenceOrderMethod(v string) (any, error) {
	if v == "" {
		return AbsenceOrderMethodDefault, nil
	}
	switch v {
	case string(AbsenceOrderMethodDefault):
		return AbsenceOrderMethodDefault, nil
	default:
		return AbsenceOrderMethodDefault, nil
	}
}

const (
	// AbsenceDefaultCursorKey はデフォルトカーソルキー。
	AbsenceDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AbsenceOrderMethod) GetCursorKeyName() string {
	switch m {
	case AbsenceOrderMethodDefault:
		return AbsenceDefaultCursorKey
	default:
		return AbsenceDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AbsenceOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AbsenceOrderMethodDefault はデフォルト。
	AbsenceOrderMethodDefault AbsenceOrderMethod = "default"
)

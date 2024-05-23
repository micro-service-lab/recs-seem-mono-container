package parameter

import (
	"github.com/google/uuid"
)

// CreateGradeParam グループ作成のパラメータ。
type CreateGradeParam struct {
	Key            string
	OrganizationID uuid.UUID
}

// WhereGradeParam グループ検索のパラメータ。
type WhereGradeParam struct{}

// GradeOrderMethod グループの並び替え方法。
type GradeOrderMethod string

// ParseGradeOrderMethod はグループの並び替え方法をパースする。
func ParseGradeOrderMethod(v string) (any, error) {
	if v == "" {
		return GradeOrderMethodDefault, nil
	}
	switch v {
	case string(GradeOrderMethodDefault):
		return GradeOrderMethodDefault, nil
	case string(GradeOrderMethodName):
		return GradeOrderMethodName, nil
	case string(GradeOrderMethodReverseName):
		return GradeOrderMethodReverseName, nil
	default:
		return GradeOrderMethodDefault, nil
	}
}

const (
	// GradeDefaultCursorKey はデフォルトカーソルキー。
	GradeDefaultCursorKey = "default"
	// GradeNameCursorKey は名前カーソルキー。
	GradeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m GradeOrderMethod) GetCursorKeyName() string {
	switch m {
	case GradeOrderMethodDefault:
		return GradeDefaultCursorKey
	case GradeOrderMethodName, GradeOrderMethodReverseName:
		return GradeNameCursorKey
	default:
		return GradeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m GradeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// GradeOrderMethodDefault はデフォルト。
	GradeOrderMethodDefault GradeOrderMethod = "default"
	// GradeOrderMethodName は名前順。
	GradeOrderMethodName GradeOrderMethod = "name"
	// GradeOrderMethodReverseName は名前逆順。
	GradeOrderMethodReverseName GradeOrderMethod = "r_name"
)

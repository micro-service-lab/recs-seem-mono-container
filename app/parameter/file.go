package parameter

import (
	"github.com/google/uuid"
)

// CreateFileParam ファイル作成のパラメータ。
type CreateFileParam struct {
	AttachableItemID uuid.UUID
}

// WhereFileParam ファイル検索のパラメータ。
type WhereFileParam struct{}

// FileOrderMethod ファイルの並び替え方法。
type FileOrderMethod string

// ParseFileOrderMethod はファイルの並び替え方法をパースする。
func ParseFileOrderMethod(v string) (any, error) {
	if v == "" {
		return FileOrderMethodDefault, nil
	}
	switch v {
	case string(FileOrderMethodDefault):
		return FileOrderMethodDefault, nil
	default:
		return FileOrderMethodDefault, nil
	}
}

const (
	// FileDefaultCursorKey はデフォルトカーソルキー。
	FileDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m FileOrderMethod) GetCursorKeyName() string {
	switch m {
	case FileOrderMethodDefault:
		return FileDefaultCursorKey
	default:
		return FileDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m FileOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// FileOrderMethodDefault はデフォルト。
	FileOrderMethodDefault FileOrderMethod = "default"
)
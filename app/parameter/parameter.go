// Package parameter layer common parameters for the application.
package parameter

import (
	"strconv"

	"github.com/google/uuid"
)

// Pagination ページネーション。
type Pagination string

const (
	// NumberedPagination ページネーション。
	NumberedPagination Pagination = "numbered"
	// CursorPagination カーソルページネーション。
	CursorPagination Pagination = "cursor"
	// NonePagination ページネーションなし。
	NonePagination Pagination = "none"
)

// ParsePaginationParam はページネーションをパースする。
func ParsePaginationParam(v string) (any, error) {
	if v == "" {
		return NonePagination, nil
	}
	switch v {
	case string(NumberedPagination):
		return NumberedPagination, nil
	case string(CursorPagination):
		return CursorPagination, nil
	default:
		return NonePagination, nil
	}
}

// Limit リミット。
type Limit int

// DefaultLimit デフォルトリミット。
const DefaultLimit Limit = 10

// ParseLimitParam はリミットをパースする。
func ParseLimitParam(v string) (any, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return DefaultLimit, nil
	}
	if i < 0 {
		return DefaultLimit, nil
	}
	return Limit(i), nil
}

// Offset オフセット。
type Offset int

// DefaultOffset デフォルトオフセット。
const DefaultOffset Offset = 0

// ParseOffsetParam はオフセットをパースする。
func ParseOffsetParam(v string) (any, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return DefaultOffset, nil
	}
	if i < 0 {
		return DefaultOffset, nil
	}
	return Offset(i), nil
}

// Cursor カーソル。
type Cursor string

// ParseCursorParam はカーソルをパースする。
func ParseCursorParam(v string) (any, error) {
	return Cursor(v), nil
}

// WithCount カウント付きかどうか。
type WithCount bool

// ParseWithCountParam はカウント付きかどうかをパースする。
func ParseWithCountParam(v string) (any, error) {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return WithCount(false), nil
	}
	return WithCount(b), nil
}

// OrderMethod 並び替え方法。
type OrderMethod interface {
	// GetCursorKeyName はカーソルキー名を取得する。
	GetCursorKeyName() string
	// GetStringValue は文字列に変換する。
	GetStringValue() string
}

// EntityID エンティティID。
type EntityID uuid.UUID

// ParseEntityIDParam はエンティティIDをパースする。
func ParseEntityIDParam(v string) (any, error) {
	id, err := uuid.Parse(v)
	if err != nil {
		return EntityID(uuid.Nil), nil
	}
	return EntityID(id), nil
}

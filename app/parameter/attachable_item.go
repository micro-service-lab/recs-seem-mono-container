package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateAttachableItemParam 添付可能アイテム作成のパラメータ。
type CreateAttachableItemParam struct {
	URL        string
	Size       entity.Float
	OwnerID    uuid.UUID
	FromOuter  bool
	MimeTypeID uuid.UUID
}

// UpdateAttachableItemParams 添付可能アイテム更新のパラメータ。
type UpdateAttachableItemParams struct {
	URL        string
	Size       entity.Float
	MimeTypeID uuid.UUID
}

// WhereAttachableItemParam 添付可能アイテム検索のパラメータ。
type WhereAttachableItemParam struct {
	WhereInMimeType bool
	InMimeTypes     []uuid.UUID
	WhereInOwner    bool
	InOwners        []uuid.UUID
}

// AttachableItemOrderMethod 添付可能アイテムの並び替え方法。
type AttachableItemOrderMethod string

// ParseAttachableItemOrderMethod は添付可能アイテムの並び替え方法をパースする。
func ParseAttachableItemOrderMethod(v string) (any, error) {
	if v == "" {
		return AttachableItemOrderMethodDefault, nil
	}
	switch v {
	case string(AttachableItemOrderMethodDefault):
		return AttachableItemOrderMethodDefault, nil
	default:
		return AttachableItemOrderMethodDefault, nil
	}
}

const (
	// AttachableItemDefaultCursorKey はデフォルトカーソルキー。
	AttachableItemDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AttachableItemOrderMethod) GetCursorKeyName() string {
	switch m {
	case AttachableItemOrderMethodDefault:
		return AttachableItemDefaultCursorKey
	default:
		return AttachableItemDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AttachableItemOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AttachableItemOrderMethodDefault はデフォルト。
	AttachableItemOrderMethodDefault AttachableItemOrderMethod = "default"
)

package entity

import (
	"github.com/google/uuid"
)

// AttachableItem 添付アイテムを表す構造体。
type AttachableItem struct {
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
	URL              string    `json:"url"`
	Size             Float     `json:"size"`
	MimeTypeID       uuid.UUID `json:"mime_type_id"`
}

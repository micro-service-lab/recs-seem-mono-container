package entity

import (
	"github.com/google/uuid"
)

// AttachableItem 添付アイテムを表す構造体。
type AttachableItem struct {
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
	OwnerID          UUID      `json:"owner_id"`
	URL              string    `json:"url"`
	Size             Float     `json:"size"`
	MimeTypeID       uuid.UUID `json:"mime_type_id"`
	Image            Image     `json:"image,omitempty"`
	File             File      `json:"file,omitempty"`
}

// AttachableItemWithMimeType 添付アイテムと MIME タイプを表す構造体。
type AttachableItemWithMimeType struct {
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
	OwnerID          UUID      `json:"owner_id"`
	URL              string    `json:"url"`
	Size             Float     `json:"size"`
	MimeType         MimeType  `json:"mime_type"`
	Image            Image     `json:"image,omitempty"`
	File             File      `json:"file,omitempty"`
}

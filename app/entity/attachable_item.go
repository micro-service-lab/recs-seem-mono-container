package entity

import (
	"github.com/google/uuid"
)

// AttachableItem 添付アイテムを表す構造体。
type AttachableItem struct {
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
	OwnerID          UUID      `json:"owner_id"`
	FromOuter        bool      `json:"from_outer"`
	URL              string    `json:"url"`
	Size             Float     `json:"size"`
	MimeTypeID       uuid.UUID `json:"mime_type_id"`
	ImageID          UUID      `json:"image_id,omitempty"`
	FileID           UUID      `json:"file_id,omitempty"`
}

// AttachableItemWithContent 添付アイテムを表す構造体。
type AttachableItemWithContent struct {
	AttachableItemID uuid.UUID             `json:"attachable_item_id"`
	OwnerID          UUID                  `json:"owner_id"`
	FromOuter        bool                  `json:"from_outer"`
	URL              string                `json:"url"`
	Size             Float                 `json:"size"`
	MimeTypeID       uuid.UUID             `json:"mime_type_id"`
	Image            NullableEntity[Image] `json:"image,omitempty"`
	File             NullableEntity[File]  `json:"file,omitempty"`
}

// AttachableItemWithContentForQuery 添付アイテムを表す構造体(クエリー用)。
type AttachableItemWithContentForQuery struct {
	Pkey Int `json:"-"`
	AttachableItemWithContent
}

// AttachableItemWithMimeType 添付アイテムと MIME タイプを表す構造体。
type AttachableItemWithMimeType struct {
	AttachableItemID uuid.UUID                `json:"attachable_item_id"`
	OwnerID          UUID                     `json:"owner_id"`
	FromOuter        bool                     `json:"from_outer"`
	URL              string                   `json:"url"`
	Size             Float                    `json:"size"`
	MimeType         NullableEntity[MimeType] `json:"mime_type"`
	Image            NullableEntity[Image]    `json:"image,omitempty"`
	File             NullableEntity[File]     `json:"file,omitempty"`
}

// AttachableItemWithMimeTypeForQuery 添付アイテムと MIME タイプを表す構造体(クエリー用)。
type AttachableItemWithMimeTypeForQuery struct {
	Pkey Int `json:"-"`
	AttachableItemWithMimeType
}

package entity

import "github.com/google/uuid"

// File ファイルを表す構造体。
type File struct {
	FileID           uuid.UUID `json:"file_id"`
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
}

// FileWithAttachableItem ファイルと添付アイテムを表す構造体。
type FileWithAttachableItem struct {
	FileID         uuid.UUID      `json:"file_id"`
	AttachableItem AttachableItem `json:"attachable_item"`
}

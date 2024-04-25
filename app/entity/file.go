package entity

import "github.com/google/uuid"

// File ファイルを表す構造体。
type File struct {
	FileID           uuid.UUID `json:"file_id"`
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
}

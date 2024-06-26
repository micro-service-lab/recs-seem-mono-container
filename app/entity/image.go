package entity

import "github.com/google/uuid"

// Image 画像を表す構造体。
type Image struct {
	ImageID          uuid.UUID `json:"image_id"`
	Height           Float     `json:"height"`
	Width            Float     `json:"width"`
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
}

package entity

import "github.com/google/uuid"

// AttachedMessage 添付メッセージを表す構造体。
type AttachedMessage struct {
	MessageID        uuid.UUID `json:"message_id"`
	AttachableItemID uuid.UUID `json:"attachable_item_id"`
}

package entity

import "github.com/google/uuid"

// AttachedMessageOnMessage メッセージに添付された添付を表す構造体。
type AttachedMessageOnMessage struct {
	AttachedMessageID uuid.UUID `json:"attached_message_id"`
	AttachableItemID  UUID      `json:"attachable_item_id"`
}

// AttachedMessageOnMessageWithAttachableItem メッセージに添付された添付を表す構造体。
type AttachedMessageOnMessageWithAttachableItem struct {
	AttachedMessageID uuid.UUID      `json:"attached_message_id"`
	AttachableItem    AttachableItem `json:"attachable_item"`
}

// AttachableItemWithOnChatRoom チャットルームに添付された添付を表す構造体。
type AttachableItemWithOnChatRoom struct {
	AttachedMessageID uuid.UUID `json:"attached_message_id"`
	MessageID         uuid.UUID `json:"message_id"`
	AttachableItemID  UUID      `json:"attachable_item_id"`
}

// AttachableItemWithOnChatRoomWithAttachableItem チャットルームに添付された添付を表す構造体。
type AttachableItemWithOnChatRoomWithAttachableItem struct {
	AttachedMessageID uuid.UUID      `json:"attached_message_id"`
	MessageID         uuid.UUID      `json:"message_id"`
	AttachableItem    AttachableItem `json:"attachable_item"`
}

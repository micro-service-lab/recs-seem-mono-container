package entity

import "github.com/google/uuid"

// AttachedMessageOnMessage メッセージに添付された添付を表す構造体。
type AttachedMessageOnMessage struct {
	AttachedMessageID uuid.UUID `json:"attached_message_id"`
	FileURL           string    `json:"file_url"`
}

// AttachableItemWithOnChatRoom チャットルームに添付された添付を表す構造体。
type AttachableItemWithOnChatRoom struct {
	AttachedMessageID uuid.UUID `json:"attached_message_id"`
	MessageID         uuid.UUID `json:"message_id"`
	FileURL           string    `json:"file_url"`
}

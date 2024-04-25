package entity

import (
	"time"

	"github.com/google/uuid"
)

// Message メッセージを表す構造体。
type Message struct {
	MessageID    uuid.UUID `json:"message_id"`
	ChatRoomID   uuid.UUID `json:"chat_room_id"`
	SenderID     UUID      `json:"sender_id"`
	Body         string    `json:"body"`
	PostedAt     time.Time `json:"posted_at"`
	LastEditedAt time.Time `json:"last_edited_at"`
}

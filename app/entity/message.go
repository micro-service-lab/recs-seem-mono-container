package entity

import (
	"time"

	"github.com/google/uuid"
)

// Message メッセージを表す構造体。
type Message struct {
	MessageID        uuid.UUID `json:"message_id"`
	ChatRoomActionID uuid.UUID `json:"chat_room_action_id"`
	SenderID         UUID      `json:"sender_id"`
	Body             string    `json:"body"`
	PostedAt         time.Time `json:"posted_at"`
	LastEditedAt     time.Time `json:"last_edited_at"`
}

// MessageCard メッセージのカードを表す構造体。
type MessageCard struct {
	MessageID uuid.UUID `json:"message_id"`
	Body      string    `json:"body"`
	PostedAt  time.Time `json:"posted_at"`
}

// MessageWithSender メッセージと送信者を表す構造体。
type MessageWithSender struct {
	MessageID        uuid.UUID                  `json:"message_id"`
	ChatRoomActionID uuid.UUID                  `json:"chat_room_action_id"`
	Sender           NullableEntity[MemberCard] `json:"sender"`
	Body             string                     `json:"body"`
	PostedAt         time.Time                  `json:"posted_at"`
	LastEditedAt     time.Time                  `json:"last_edited_at"`
}

// MessageWithSenderForQuery メッセージと送信者を表す構造体(クエリ用)。
type MessageWithSenderForQuery struct {
	Pkey Int `json:"-"`
	MessageWithSender
}

// MessageWithSenderAndReadReceiptCountAndAttachments メッセージと送信者、既読情報、添付ファイルを表す構造体。
type MessageWithSenderAndReadReceiptCountAndAttachments struct {
	MessageID        uuid.UUID                  `json:"message_id"`
	ChatRoomActionID uuid.UUID                  `json:"chat_room_action_id"`
	Sender           NullableEntity[MemberCard] `json:"sender"`
	Body             string                     `json:"body"`
	PostedAt         time.Time                  `json:"posted_at"`
	LastEditedAt     time.Time                  `json:"last_edited_at"`
	ReadReceiptCount int64                      `json:"read_receipt_count"`
	Attachments      []AttachedItemOnMessage    `json:"attachments"`
}

// MessageWithChatRoom メッセージとチャットルームを表す構造体。
type MessageWithChatRoom struct {
	MessageID    uuid.UUID              `json:"message_id"`
	ChatRoom     ChatRoomWithCoverImage `json:"chat_room"`
	SenderID     UUID                   `json:"sender_id"`
	Body         string                 `json:"body"`
	PostedAt     time.Time              `json:"posted_at"`
	LastEditedAt time.Time              `json:"last_edited_at"`
}

// MessageWithChatRoomForQuery メッセージとチャットルームを表す構造体(クエリ用)。
type MessageWithChatRoomForQuery struct {
	Pkey Int `json:"-"`
	MessageWithChatRoom
}

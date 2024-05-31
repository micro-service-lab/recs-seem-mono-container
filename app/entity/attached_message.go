package entity

import (
	"github.com/google/uuid"
)

// AttachedMessage ロールの関連付けを表すインターフェース。
type AttachedMessage struct {
	AttachedMessageID uuid.UUID `json:"attached_message_id"`
	MessageID         uuid.UUID `json:"message_id"`
	AttachableItemID  UUID      `json:"attachable_item_id"`
}

// AttachedItemOnMessage メッセージに添付された添付を表す構造体。
type AttachedItemOnMessage struct {
	AttachedMessageID uuid.UUID      `json:"attached_message_id"`
	MessageID         uuid.UUID      `json:"message_id"`
	AttachableItem    AttachableItem `json:"attachable_item"`
}

// AttachedItemOnMessageForQuery メッセージに添付された添付を表す構造体(クエリー用)。
type AttachedItemOnMessageForQuery struct {
	Pkey Int `json:"-"`
	AttachedItemOnMessage
}

// AttachedItemOnMessageWithMimeType メッセージに添付された添付を表す構造体。
type AttachedItemOnMessageWithMimeType struct {
	AttachedMessageID uuid.UUID                  `json:"attached_message_id"`
	MessageID         uuid.UUID                  `json:"message_id"`
	AttachableItem    AttachableItemWithMimeType `json:"attachable_item"`
}

// AttachedItemOnMessageWithMimeTypeForQuery メッセージに添付された添付を表す構造体(クエリー用)。
type AttachedItemOnMessageWithMimeTypeForQuery struct {
	Pkey Int `json:"-"`
	AttachedItemOnMessageWithMimeType
}

// AttachedItemOnChatRoom チャットルームに添付された添付を表す構造体。
type AttachedItemOnChatRoom struct {
	AttachedMessageID uuid.UUID      `json:"attached_message_id"`
	AttachableItem    AttachableItem `json:"attachable_item"`
	Message           Message        `json:"message"`
}

// AttachedItemOnChatRoomForQuery チャットルームに添付された添付を表す構造体(クエリー用)。
type AttachedItemOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	AttachedItemOnChatRoom
}

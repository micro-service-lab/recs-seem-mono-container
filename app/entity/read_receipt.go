package entity

import (
	"github.com/google/uuid"
)

// ReadReceipt 既読情報を表す構造体。
type ReadReceipt struct {
	MemberID  uuid.UUID   `json:"member_id"`
	MessageID uuid.UUID   `json:"message_id"`
	ReadAt    Timestamptz `json:"read_at"`
}

// ReadReceiptGroupByChatRoom チャットルームごとの既読情報を表す構造体。
type ReadReceiptGroupByChatRoom struct {
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Count      int64     `json:"count"`
}

// ReadableMemberOnMessage メッセージ上のメンバーを表す構造体。
type ReadableMemberOnMessage struct {
	Member MemberCard  `json:"member"`
	ReadAt Timestamptz `json:"read_at"`
}

// ReadableMemberOnMessageForQuery メッセージ上のメンバーを表す構造体(クエリ用)。
type ReadableMemberOnMessageForQuery struct {
	Pkey Int `json:"-"`
	ReadableMemberOnMessage
}

// ReadableMessageOnMember メンバー上のメッセージを表す構造体。
type ReadableMessageOnMember struct {
	Message Message     `json:"message"`
	ReadAt  Timestamptz `json:"read_at"`
}

// ReadableMessageOnMemberForQuery メンバー上のメッセージを表す構造体(クエリ用)。
type ReadableMessageOnMemberForQuery struct {
	Pkey Int `json:"-"`
	ReadableMessageOnMember
}

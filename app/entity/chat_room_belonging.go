package entity

import (
	"time"

	"github.com/google/uuid"
)

// ChatRoomBelonging チャットルーム所属を表す構造体。
type ChatRoomBelonging struct {
	MemberID   uuid.UUID `json:"member_id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	AddedAt    time.Time `json:"added_at"`
}

// PrivateChatRoomCompanions プライベートチャットルームの相手を表す構造体。
type PrivateChatRoomCompanions struct {
	Member     MemberCard `json:"member"`
	ChatRoomID uuid.UUID  `json:"chat_room_id"`
	AddedAt    time.Time  `json:"added_at"`
}

// ChatRoomBelongingMember チャットルーム所属のメンバーを表す構造体。
type ChatRoomBelongingMember struct {
	Member  MemberCard `json:"member"`
	AddedAt time.Time  `json:"added_at"`
}

// MemberOnChatRoom チャットルーム上のメンバーを表す構造体。
type MemberOnChatRoom struct {
	Member  MemberCard `json:"member"`
	AddedAt time.Time  `json:"added_at"`
}

// MemberOnChatRoomForQuery チャットルーム上のメンバーを表す構造体(クエリ用)。
type MemberOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	MemberOnChatRoom
}

// PracticalChatRoomOnMember メンバー上の実用的なチャットルームを表す構造体。
type PracticalChatRoomOnMember struct {
	ChatRoom PracticalChatRoom `json:"chat_room"`
	AddedAt  time.Time         `json:"added_at"`
}

// ChatRoomOnMember メンバー上のチャットルームを表す構造体。
type ChatRoomOnMember struct {
	ChatRoom ChatRoomWithLatestAndCoverImage `json:"chat_room"`
	AddedAt  time.Time                       `json:"added_at"`
}

// ChatRoomOnMemberForQuery メンバー上のチャットルームを表す構造体(クエリ用)。
type ChatRoomOnMemberForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomOnMember
}

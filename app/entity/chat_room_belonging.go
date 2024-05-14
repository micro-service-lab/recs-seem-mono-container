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

// ChatRoomBelongingMember チャットルーム所属のメンバーを表す構造体。
type ChatRoomBelongingMember struct {
	Member  MemberCard `json:"member"`
	AddedAt time.Time  `json:"added_at"`
}

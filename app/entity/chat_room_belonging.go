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

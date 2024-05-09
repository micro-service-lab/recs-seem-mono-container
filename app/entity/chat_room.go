package entity

import "github.com/google/uuid"

// ChatRoom チャットルームを表す構造体。
type ChatRoom struct {
	ChatRoomID       uuid.UUID `json:"chat_room_id"`
	Name             String    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	FromOrganization bool      `json:"from_organization"`
	CoverImageURL    String    `json:"cover_image_url"`
	OwnerID          UUID      `json:"owner_id"`
}

package entity

import (
	"github.com/google/uuid"
)

// ChatRoom チャットルームを表す構造体。
type ChatRoom struct {
	ChatRoomID       uuid.UUID `json:"chat_room_id"`
	Name             string    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	FromOrganization bool      `json:"from_organization"`
	CoverImageID     UUID      `json:"cover_image_id"`
	OwnerID          UUID      `json:"owner_id"`
}

// ChatRoomWithCoverImage チャットルームとカバー画像を表す構造体。
type ChatRoomWithCoverImage struct {
	ChatRoomID       uuid.UUID                               `json:"chat_room_id"`
	Name             string                                  `json:"name"`
	IsPrivate        bool                                    `json:"is_private"`
	FromOrganization bool                                    `json:"from_organization"`
	CoverImage       NullableEntity[ImageWithAttachableItem] `json:"cover_image"`
	OwnerID          UUID                                    `json:"owner_id"`
}

// PracticalChatRoom 実用的なチャットルームを表す構造体。
type PracticalChatRoom struct {
	ChatRoomID       uuid.UUID                               `json:"chat_room_id"`
	Name             string                                  `json:"name"`
	IsPrivate        bool                                    `json:"is_private"`
	FromOrganization bool                                    `json:"from_organization"`
	CoverImage       NullableEntity[ImageWithAttachableItem] `json:"cover_image"`
	OwnerID          UUID                                    `json:"owner_id"`
	LatestMessage    NullableEntity[MessageCard]             `json:"latest_message"`
}

// ChatRoomWithOwner チャットルームとオーナーを表す構造体。
type ChatRoomWithOwner struct {
	ChatRoomID       uuid.UUID `json:"chat_room_id"`
	Name             string    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	FromOrganization bool      `json:"from_organization"`
	CoverImageID     UUID      `json:"cover_image_id"`
	Owner            Member    `json:"owner"`
}

// ChatRoomOnPrivateWithMember チャットルームを表す構造体。
type ChatRoomOnPrivateWithMember struct {
	ChatRoomID       uuid.UUID                               `json:"chat_room_id"`
	Name             string                                  `json:"name"`
	IsPrivate        bool                                    `json:"is_private"`
	OwnerID          UUID                                    `json:"owner_id"`
	FromOrganization bool                                    `json:"from_organization"`
	CoverImage       NullableEntity[ImageWithAttachableItem] `json:"cover_image"`
	Partner          ChatRoomBelongingMember                 `json:"partner"`
}

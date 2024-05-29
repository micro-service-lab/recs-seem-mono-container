package entity

import "github.com/google/uuid"

// ChatRoomCreateAction チャットルーム作成アクションを表す構造体。
type ChatRoomCreateAction struct {
	ChatRoomCreateActionID uuid.UUID `json:"chat_room_create_action_id"`
	ChatRoomActionID       uuid.UUID `json:"chat_room_action_id"`
	Name                   string    `json:"name"`
	CreatedBy              UUID      `json:"created_by"`
}

// ChatRoomCreateActionOnChatRoom チャットルーム作成アクションを表す構造体。
type ChatRoomCreateActionOnChatRoom struct {
	ChatRoomCreateActionID uuid.UUID                    `json:"chat_room_create_action_id"`
	ChatRoomActionID       uuid.UUID                    `json:"chat_room_action_id"`
	Name                   string                       `json:"name"`
	CreatedBy              NullableEntity[SimpleMember] `json:"created_by"`
}

// ChatRoomCreateActionOnChatRoomForQuery チャットルーム作成アクションを表す構造体(クエリー用)。
type ChatRoomCreateActionOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomCreateActionOnChatRoom
}

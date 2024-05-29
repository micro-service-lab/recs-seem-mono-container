package entity

import "github.com/google/uuid"

// ChatRoomUpdateNameAction チャットルーム名前更新アクションを表す構造体。
type ChatRoomUpdateNameAction struct {
	ChatRoomUpdateNameActionID uuid.UUID `json:"chat_room_update_name_action_id"`
	ChatRoomActionID           uuid.UUID `json:"chat_room_action_id"`
	Name                       string    `json:"name"`
	UpdatedBy                  UUID      `json:"updated_by"`
}

// ChatRoomUpdateNameActionOnChatRoom チャットルーム名前更新アクションを表す構造体。
type ChatRoomUpdateNameActionOnChatRoom struct {
	ChatRoomUpdateNameActionID uuid.UUID                    `json:"chat_room_update_name_action_id"`
	ChatRoomActionID           uuid.UUID                    `json:"chat_room_action_id"`
	Name                       string                       `json:"name"`
	UpdatedBy                  NullableEntity[SimpleMember] `json:"updated_by"`
}

// ChatRoomUpdateNameActionOnChatRoomForQuery チャットルーム名前更新アクションを表す構造体(クエリー用)。
type ChatRoomUpdateNameActionOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomUpdateNameActionOnChatRoom
}

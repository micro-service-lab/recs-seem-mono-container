package entity

import "github.com/google/uuid"

// ChatRoomCreateAction チャットルーム作成アクションを表す構造体。
type ChatRoomCreateAction struct {
	ChatRoomCreateActionID uuid.UUID `json:"chat_room_create_action_id"`
	ChatRoomActionID       uuid.UUID `json:"chat_room_action_id"`
	Name                   string    `json:"name"`
	CreatedBy              UUID      `json:"created_by"`
}

// ChatRoomCreateActionWithCreatedBy チャットルーム作成アクションを表す構造体。
type ChatRoomCreateActionWithCreatedBy struct {
	ChatRoomCreateActionID uuid.UUID                    `json:"chat_room_create_action_id"`
	ChatRoomActionID       uuid.UUID                    `json:"chat_room_action_id"`
	Name                   string                       `json:"name"`
	CreatedBy              NullableEntity[SimpleMember] `json:"created_by"`
}

// ChatRoomCreateActionWithCreatedByForQuery チャットルーム作成アクションを表す構造体(クエリー用)。
type ChatRoomCreateActionWithCreatedByForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomCreateActionWithCreatedBy
}

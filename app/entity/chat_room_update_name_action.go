package entity

import "github.com/google/uuid"

// ChatRoomUpdateNameAction チャットルーム名前更新アクションを表す構造体。
type ChatRoomUpdateNameAction struct {
	ChatRoomUpdateNameActionID uuid.UUID `json:"chat_room_update_name_action_id"`
	ChatRoomActionID           uuid.UUID `json:"chat_room_action_id"`
	Name                       string    `json:"name"`
	UpdatedBy                  UUID      `json:"updated_by"`
}

// ChatRoomUpdateNameActionWithUpdatedBy チャットルーム名前更新アクションを表す構造体。
type ChatRoomUpdateNameActionWithUpdatedBy struct {
	ChatRoomUpdateNameActionID uuid.UUID                    `json:"chat_room_update_name_action_id"`
	ChatRoomActionID           uuid.UUID                    `json:"chat_room_action_id"`
	Name                       string                       `json:"name"`
	UpdatedBy                  NullableEntity[SimpleMember] `json:"updated_by"`
}

// ChatRoomUpdateNameActionWithUpdatedByForQuery チャットルーム名前更新アクションを表す構造体(クエリー用)。
type ChatRoomUpdateNameActionWithUpdatedByForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomUpdateNameActionWithUpdatedBy
}

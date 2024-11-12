package entity

import "github.com/google/uuid"

// ChatRoomDeleteMessageAction チャットルームメッセージ削除アクションを表す構造体。
type ChatRoomDeleteMessageAction struct {
	ChatRoomDeleteMessageActionID uuid.UUID `json:"chat_room_delete_message_action_id"`
	ChatRoomActionID              uuid.UUID `json:"chat_room_action_id"`
	DeletedBy                     UUID      `json:"deleted_by"`
}

// ChatRoomDeleteMessageActionWithDeletedBy チャットルームメッセージ削除アクションを表す構造体。
type ChatRoomDeleteMessageActionWithDeletedBy struct {
	ChatRoomDeleteMessageActionID uuid.UUID                    `json:"chat_room_delete_message_action_id"`
	ChatRoomActionID              uuid.UUID                    `json:"chat_room_action_id"`
	DeletedBy                     NullableEntity[SimpleMember] `json:"deleted_by"`
}

// ChatRoomDeleteMessageActionWithDeletedByForQuery チャットルームメッセージ削除アクションを表す構造体(クエリー用)。
type ChatRoomDeleteMessageActionWithDeletedByForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomDeleteMessageActionWithDeletedBy
}

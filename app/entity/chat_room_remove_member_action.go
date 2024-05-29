package entity

import "github.com/google/uuid"

// ChatRoomRemoveMemberAction チャットルームメンバー強制退会アクションを表す構造体。
type ChatRoomRemoveMemberAction struct {
	ChatRoomRemoveMemberActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	ChatRoomActionID             uuid.UUID `json:"chat_room_action_id"`
	RemovedBy                    UUID      `json:"removed_by"`
}

// ChatRoomRemoveMemberActionOnChatRoom チャットルームメンバー強制退会アクションを表す構造体。
type ChatRoomRemoveMemberActionOnChatRoom struct {
	ChatRoomRemoveMemberActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	ChatRoomActionID             uuid.UUID                    `json:"chat_room_action_id"`
	RemovedBy                    NullableEntity[SimpleMember] `json:"removed_by"`
}

// ChatRoomRemoveMemberActionOnChatRoomForQuery チャットルームメンバー強制退会アクションを表す構造体(クエリー用)。
type ChatRoomRemoveMemberActionOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomRemoveMemberActionOnChatRoom
}

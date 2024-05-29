package entity

import "github.com/google/uuid"

// ChatRoomAddMemberAction チャットルームメンバー追加アクションを表す構造体。
type ChatRoomAddMemberAction struct {
	ChatRoomAddMemberActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	ChatRoomActionID          uuid.UUID `json:"chat_room_action_id"`
	AddedBy                   UUID      `json:"added_by"`
}

// ChatRoomAddMemberActionOnChatRoom チャットルームメンバー追加アクションを表す構造体。
type ChatRoomAddMemberActionOnChatRoom struct {
	ChatRoomAddMemberActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	ChatRoomActionID          uuid.UUID                    `json:"chat_room_action_id"`
	AddedBy                   NullableEntity[SimpleMember] `json:"added_by"`
}

// ChatRoomAddMemberActionOnChatRoomForQuery チャットルームメンバー追加アクションを表す構造体(クエリー用)。
type ChatRoomAddMemberActionOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomAddMemberActionOnChatRoom
}

package entity

import "github.com/google/uuid"

// ChatRoomWithdrawAction チャットルームメンバー退会アクションを表す構造体。
type ChatRoomWithdrawAction struct {
	ChatRoomWithdrawActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	ChatRoomActionID         uuid.UUID `json:"chat_room_action_id"`
	MemberID                 UUID      `json:"member_id"`
}

// ChatRoomWithdrawActionOnChatRoom チャットルームメンバー退会アクションを表す構造体。
type ChatRoomWithdrawActionOnChatRoom struct {
	ChatRoomWithdrawActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	ChatRoomActionID         uuid.UUID                    `json:"chat_room_action_id"`
	MemberID                 NullableEntity[SimpleMember] `json:"member_id"`
}

// ChatRoomWithdrawActionOnChatRoomForQuery チャットルームメンバー退会アクションを表す構造体(クエリー用)。
type ChatRoomWithdrawActionOnChatRoomForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomWithdrawActionOnChatRoom
}

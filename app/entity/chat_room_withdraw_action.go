package entity

import "github.com/google/uuid"

// ChatRoomWithdrawAction チャットルームメンバー退会アクションを表す構造体。
type ChatRoomWithdrawAction struct {
	ChatRoomWithdrawActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	ChatRoomActionID         uuid.UUID `json:"chat_room_action_id"`
	MemberID                 UUID      `json:"member_id"`
}

// ChatRoomWithdrawActionWithMember チャットルームメンバー退会アクションを表す構造体。
type ChatRoomWithdrawActionWithMember struct {
	ChatRoomWithdrawActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	ChatRoomActionID         uuid.UUID                    `json:"chat_room_action_id"`
	Member                   NullableEntity[SimpleMember] `json:"member_id"`
}

// ChatRoomWithdrawActionWithMemberForQuery チャットルームメンバー退会アクションを表す構造体(クエリー用)。
type ChatRoomWithdrawActionWithMemberForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomWithdrawActionWithMember
}

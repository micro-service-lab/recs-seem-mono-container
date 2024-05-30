package entity

import "github.com/google/uuid"

// ChatRoomAddedMember チャットルーム追加メンバーを表す構造体。
type ChatRoomAddedMember struct {
	ChatRoomAddMemberActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	MemberID                  UUID      `json:"member_id"`
}

// MemberOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバーを表す構造体。
type MemberOnChatRoomAddMemberAction struct {
	ChatRoomAddMemberActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	Member                    NullableEntity[SimpleMember] `json:"member"`
}

// MemberOnChatRoomAddMemberActionForQuery チャットルーム追加メンバーアクション上のメンバーを表す構造体(クエリー用)。
type MemberOnChatRoomAddMemberActionForQuery struct {
	Pkey Int `json:"-"`
	MemberOnChatRoomAddMemberAction
}

package entity

import "github.com/google/uuid"

// ChatRoomRemovedMember チャットルーム追放メンバーを表す構造体。
type ChatRoomRemovedMember struct {
	ChatRoomRemoveMemberActionID uuid.UUID `json:"chat_room_remove_member_action_id"`
	MemberID                     UUID      `json:"member_id"`
}

// MemberOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバーを表す構造体。
type MemberOnChatRoomRemoveMemberAction struct {
	ChatRoomRemoveMemberActionID uuid.UUID                    `json:"chat_room_remove_member_action_id"`
	Member                       NullableEntity[SimpleMember] `json:"member"`
}

// MemberOnChatRoomRemoveMemberActionForQuery チャットルーム追放メンバーアクション上のメンバーを表す構造体(クエリー用)。
type MemberOnChatRoomRemoveMemberActionForQuery struct {
	Pkey Int `json:"-"`
	MemberOnChatRoomRemoveMemberAction
}

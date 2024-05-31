package entity

import "github.com/google/uuid"

// ChatRoomRemoveMemberAction チャットルームメンバー強制退会アクションを表す構造体。
type ChatRoomRemoveMemberAction struct {
	ChatRoomRemoveMemberActionID uuid.UUID `json:"chat_room_remove_member_action_id"`
	ChatRoomActionID             uuid.UUID `json:"chat_room_action_id"`
	RemovedBy                    UUID      `json:"removed_by"`
}

// ChatRoomRemoveMemberActionWithRemovedBy チャットルームメンバー強制退会アクションを表す構造体。
type ChatRoomRemoveMemberActionWithRemovedBy struct {
	ChatRoomRemoveMemberActionID uuid.UUID                    `json:"chat_room_remove_member_action_id"`
	ChatRoomActionID             uuid.UUID                    `json:"chat_room_action_id"`
	RemovedBy                    NullableEntity[SimpleMember] `json:"removed_by"`
}

// ChatRoomRemoveMemberActionWithRemovedByForQuery チャットルームメンバー強制退会アクションを表す構造体(クエリー用)。
type ChatRoomRemoveMemberActionWithRemovedByForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomRemoveMemberActionWithRemovedBy
}

// ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers チャットルームメンバー強制退会アクションを表す構造体。
type ChatRoomRemoveMemberActionWithRemovedByAndRemoveMember struct {
	ChatRoomRemoveMemberActionID uuid.UUID                            `json:"chat_room_remove_member_action_id"`
	ChatRoomActionID             uuid.UUID                            `json:"chat_room_action_id"`
	RemovedBy                    NullableEntity[SimpleMember]         `json:"removed_by"`
	RemoveMembers                []MemberOnChatRoomRemoveMemberAction `json:"remove_members"`
}

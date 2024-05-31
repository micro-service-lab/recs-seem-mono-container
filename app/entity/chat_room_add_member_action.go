package entity

import "github.com/google/uuid"

// ChatRoomAddMemberAction チャットルームメンバー追加アクションを表す構造体。
type ChatRoomAddMemberAction struct {
	ChatRoomAddMemberActionID uuid.UUID `json:"chat_room_add_member_action_id"`
	ChatRoomActionID          uuid.UUID `json:"chat_room_action_id"`
	AddedBy                   UUID      `json:"added_by"`
}

// ChatRoomAddMemberActionWithAddedBy チャットルームメンバー追加アクションを表す構造体。
type ChatRoomAddMemberActionWithAddedBy struct {
	ChatRoomAddMemberActionID uuid.UUID                    `json:"chat_room_add_member_action_id"`
	ChatRoomActionID          uuid.UUID                    `json:"chat_room_action_id"`
	AddedBy                   NullableEntity[SimpleMember] `json:"added_by"`
}

// ChatRoomAddMemberActionWithAddedByForQuery チャットルームメンバー追加アクションを表す構造体(クエリー用)。
type ChatRoomAddMemberActionWithAddedByForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomAddMemberActionWithAddedBy
}

// ChatRoomAddMemberActionWithAddedByAndAddMembers チャットルームメンバー追加アクションを表す構造体。
type ChatRoomAddMemberActionWithAddedByAndAddMembers struct {
	ChatRoomAddMemberActionID uuid.UUID                         `json:"chat_room_add_member_action_id"`
	ChatRoomActionID          uuid.UUID                         `json:"chat_room_action_id"`
	AddedBy                   NullableEntity[SimpleMember]      `json:"added_by"`
	AddMembers                []MemberOnChatRoomAddMemberAction `json:"add_members"`
}

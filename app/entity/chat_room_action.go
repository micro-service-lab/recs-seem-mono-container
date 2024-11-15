package entity

import (
	"time"

	"github.com/google/uuid"
)

// ChatRoomAction チャットルームアクションを表す構造体。
type ChatRoomAction struct {
	ChatRoomActionID     uuid.UUID `json:"chat_room_action_id"`
	ChatRoomID           uuid.UUID `json:"chat_room_id"`
	ChatRoomActionTypeID uuid.UUID `json:"chat_room_action_type_id"`
	ActedAt              time.Time `json:"acted_at"`
}

// ChatRoomActionWithActionType チャットルームアクションを表す構造体。
type ChatRoomActionWithActionType struct {
	ChatRoomActionID   uuid.UUID          `json:"chat_room_action_id"`
	ChatRoomID         uuid.UUID          `json:"chat_room_id"`
	ActedAt            time.Time          `json:"acted_at"`
	ChatRoomActionType ChatRoomActionType `json:"chat_room_action_type"`
}

// ChatRoomActionWithDetail チャットルームアクションを表す構造体。
//
//nolint:lll
type ChatRoomActionWithDetail struct {
	ChatRoomActionID            uuid.UUID                                                `json:"chat_room_action_id"`
	ChatRoomID                  uuid.UUID                                                `json:"chat_room_id"`
	ChatRoomActionTypeID        uuid.UUID                                                `json:"chat_room_action_type_id"`
	ActedAt                     time.Time                                                `json:"acted_at"`
	ChatRoomCreateAction        NullableEntity[ChatRoomCreateActionWithCreatedBy]        `json:"chat_room_create_action,omitempty"`
	ChatRoomUpdateNameAction    NullableEntity[ChatRoomUpdateNameActionWithUpdatedBy]    `json:"chat_room_update_name_action,omitempty"`
	ChatRoomAddMemberAction     NullableEntity[ChatRoomAddMemberActionWithAddedBy]       `json:"chat_room_add_member_action,omitempty"`
	ChatRoomRemoveMemberAction  NullableEntity[ChatRoomRemoveMemberActionWithRemovedBy]  `json:"chat_room_remove_member_action,omitempty"`
	ChatRoomWithdrawAction      NullableEntity[ChatRoomWithdrawActionWithMember]         `json:"chat_room_withdraw_action,omitempty"`
	ChatRoomDeleteMessageAction NullableEntity[ChatRoomDeleteMessageActionWithDeletedBy] `json:"chat_room_delete_message_action,omitempty"`
	Message                     NullableEntity[MessageWithSender]                        `json:"message,omitempty"`
}

// ChatRoomActionWithDetailForQuery チャットルームアクションを表す構造体(クエリー用)。
type ChatRoomActionWithDetailForQuery struct {
	Pkey Int `json:"-"`
	ChatRoomActionWithDetail
}

// ChatRoomActionPractical チャットルームアクションを表す構造体。
//
//nolint:lll
type ChatRoomActionPractical struct {
	ChatRoomActionID            uuid.UUID                                                               `json:"chat_room_action_id"`
	ChatRoomID                  uuid.UUID                                                               `json:"chat_room_id"`
	ChatRoomActionTypeID        uuid.UUID                                                               `json:"chat_room_action_type_id"`
	ActedAt                     time.Time                                                               `json:"acted_at"`
	ChatRoomCreateAction        NullableEntity[ChatRoomCreateActionWithCreatedBy]                       `json:"chat_room_create_action,omitempty"`
	ChatRoomUpdateNameAction    NullableEntity[ChatRoomUpdateNameActionWithUpdatedBy]                   `json:"chat_room_update_name_action,omitempty"`
	ChatRoomAddMemberAction     NullableEntity[ChatRoomAddMemberActionWithAddedByAndAddMembers]         `json:"chat_room_add_member_action,omitempty"`
	ChatRoomRemoveMemberAction  NullableEntity[ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers] `json:"chat_room_remove_member_action,omitempty"`
	ChatRoomWithdrawAction      NullableEntity[ChatRoomWithdrawActionWithMember]                        `json:"chat_room_withdraw_action,omitempty"`
	ChatRoomDeleteMessageAction NullableEntity[ChatRoomDeleteMessageActionWithDeletedBy]                `json:"chat_room_delete_message_action,omitempty"`
	Message                     NullableEntity[MessageWithSenderAndReadReceiptCountAndAttachments]      `json:"message,omitempty"`
}

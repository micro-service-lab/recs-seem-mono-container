package ws

import (
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// EventType イベントの種別を表す型。
type EventType string

const (
	// EventTypeConnectingMembers クライアントが接続した際、接続しているメンバーを表すイベント。
	EventTypeConnectingMembers EventType = "connecting_members"
	// EventTypeConnected クライアントが接続したことを表すイベント。
	EventTypeConnected EventType = "connected"
	// EventTypeDisconnected クライアントが切断したことを表すイベント。
	EventTypeDisconnected EventType = "disconnected"
	// EventTypeChatRoomAddedMe チャットルームに自分が追加されたことを表すイベント。
	EventTypeChatRoomAddedMe EventType = "chat_room:added:me"
	// EventTypeChatRoomAddedMember チャットルームに自分以外のメンバーが追加されたことを表すイベント。
	EventTypeChatRoomAddedMember EventType = "chat_room:added:member"
	// EventTypeChatRoomRemovedMe チャットルームから自分が削除されたことを表すイベント。
	EventTypeChatRoomRemovedMe EventType = "chat_room:removed:me"
	// EventTypeChatRoomRemovedMember チャットルームから自分以外のメンバーが削除されたことを表すイベント。
	EventTypeChatRoomRemovedMember EventType = "chat_room:removed:member"
	// EventTypeChatRoomWithdrawnMember チャットルームから自分以外のメンバーが退室したことを表すイベント。
	EventTypeChatRoomWithdrawnMember EventType = "chat_room:withdrawn:member"
	// EventTypeChatRoomUpdatedName チャットルームの名前が更新されたことを表すイベント。
	EventTypeChatRoomUpdatedName EventType = "chat_room:updated:name"
	// EventTypeChatRoomDeletedMessage チャットルームのメッセージが削除されたことを表すイベント。
	EventTypeChatRoomDeletedMessage EventType = "chat_room:deleted:message"
	// EventTypeChatRoomEditedMessage チャットルームのメッセージが編集されたことを表すイベント。
	EventTypeChatRoomEditedMessage EventType = "chat_room:edited:message"
	// EventTypeChatRoomSentMessage チャットルームにメッセージが送信されたことを表すイベント。
	EventTypeChatRoomSentMessage EventType = "chat_room:sent:message"
	// EventTypeChatRoomReadMessage チャットルームのメッセージが既読になったことを表すイベント。
	EventTypeChatRoomReadMessage EventType = "chat_room:read:message"
	// EventTypeChatRoomDeleted チャットルームが削除されたことを表すイベント。
	EventTypeChatRoomDeleted EventType = "chat_room:deleted"
)

// ConnectingMembersEventData クライアントが接続した際、接続しているメンバーを表す構造体。
type ConnectingMembersEventData struct {
	ConnectingMembers []ConnectingMember `json:"connecting_members"`
}

// ConnectingMember クライアントが接続した際、接続しているメンバーを表す構造体。
type ConnectingMember struct {
	MemberID   uuid.UUID   `json:"member_id"`
	ConnectIDs []uuid.UUID `json:"connect_ids"`
}

// ConnectedEventData クライアントが接続した際のイベントデータを表す構造体。
type ConnectedEventData struct {
	AuthMemberID uuid.UUID `json:"auth_member_id"`
	ConnectID    uuid.UUID `json:"connect_id"`
}

// DisconnectedEventData クライアントが切断した際のイベントデータを表す構造体。
type DisconnectedEventData struct {
	AuthMemberID uuid.UUID `json:"auth_member_id"`
	ConnectID    uuid.UUID `json:"connect_id"`
}

// ChatRoomAddedMeEventData チャットルームに自分が追加された際のイベントデータを表す構造体。
type ChatRoomAddedMeEventData struct {
	ChatRoom entity.ChatRoom `json:"chat_room"`
}

// ChatRoomAddedMemberEventData チャットルームに自分以外のメンバーが追加された際のイベントデータを表す構造体。
type ChatRoomAddedMemberEventData struct {
	ChatRoomID           uuid.UUID                                              `json:"chat_room_id"`
	Action               entity.ChatRoomAddMemberActionWithAddedByAndAddMembers `json:"action"`
	ChatRoomActionID     uuid.UUID                                              `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                              `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                              `json:"acted_at"`
}

// ChatRoomRemovedMeEventData チャットルームから自分が削除された際のイベントデータを表す構造体。
type ChatRoomRemovedMeEventData struct {
	ChatRoomID           uuid.UUID                                                      `json:"chat_room_id"`
	Action               entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers `json:"action"`
	ChatRoomActionID     uuid.UUID                                                      `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                                      `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                                      `json:"acted_at"`
}

// ChatRoomRemovedMemberEventData チャットルームから自分以外のメンバーが削除された際のイベントデータを表す構造体。
type ChatRoomRemovedMemberEventData struct {
	ChatRoomID           uuid.UUID                                                      `json:"chat_room_id"`
	Action               entity.ChatRoomRemoveMemberActionWithRemovedByAndRemoveMembers `json:"action"`
	ChatRoomActionID     uuid.UUID                                                      `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                                      `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                                      `json:"acted_at"`
}

// ChatRoomWithdrawnMemberEventData チャットルームから自分以外のメンバーが退室した際のイベントデータを表す構造体。
type ChatRoomWithdrawnMemberEventData struct {
	ChatRoomID           uuid.UUID                               `json:"chat_room_id"`
	Action               entity.ChatRoomWithdrawActionWithMember `json:"action"`
	ChatRoomActionID     uuid.UUID                               `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                               `json:"chat_room_action_type_id"`
	ActedAt              time.Time                               `json:"acted_at"`
}

// ChatRoomUpdatedNameEventData チャットルームの名前が更新された際のイベントデータを表す構造体。
type ChatRoomUpdatedNameEventData struct {
	ChatRoomID           uuid.UUID                                    `json:"chat_room_id"`
	Action               entity.ChatRoomUpdateNameActionWithUpdatedBy `json:"action"`
	ChatRoomActionID     uuid.UUID                                    `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                    `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                    `json:"acted_at"`
}

// ChatRoomDeletedMessageEventData チャットルームのメッセージが削除された際のイベントデータを表す構造体。
type ChatRoomDeletedMessageEventData struct {
	ChatRoomID           uuid.UUID                                       `json:"chat_room_id"`
	Action               entity.ChatRoomDeleteMessageActionWithDeletedBy `json:"action"`
	ChatRoomActionID     uuid.UUID                                       `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                       `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                       `json:"acted_at"`
}

// ChatRoomEditedMessageEventData チャットルームのメッセージが編集された際のイベントデータを表す構造体。
type ChatRoomEditedMessageEventData struct {
	ChatRoomID uuid.UUID      `json:"chat_room_id"`
	Message    entity.Message `json:"message"`
}

// ChatRoomSentMessageEventData チャットルームにメッセージが送信された際のイベントデータを表す構造体。
type ChatRoomSentMessageEventData struct {
	ChatRoomID           uuid.UUID                                                 `json:"chat_room_id"`
	Action               entity.MessageWithSenderAndReadReceiptCountAndAttachments `json:"action"`
	ChatRoomActionID     uuid.UUID                                                 `json:"chat_room_action_id"`
	ChatRoomActionTypeID uuid.UUID                                                 `json:"chat_room_action_type_id"`
	ActedAt              time.Time                                                 `json:"acted_at"`
}

// ChatRoomReadMessageEventData チャットルームのメッセージが既読になった際のイベントデータを表す構造体。
type ChatRoomReadMessageEventData struct {
	ChatRoomID uuid.UUID   `json:"chat_room_id"`
	MessageIDs []uuid.UUID `json:"message_ids"`
}

// ChatRoomDeletedEventData チャットルームが削除された際のイベントデータを表す構造体。
type ChatRoomDeletedEventData struct {
	ChatRoom  entity.ChatRoom `json:"chat_room"`
	DeletedBy entity.Member   `json:"deleted_by"`
}

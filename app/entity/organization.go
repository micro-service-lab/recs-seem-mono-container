package entity

import (
	"github.com/google/uuid"
)

// Organization 組織を表す構造体。
type Organization struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Color          String    `json:"color"`
	Description    String    `json:"description"`
	IsPersonal     bool      `json:"is_personal"`
	IsWhole        bool      `json:"is_whole"`
	ChatRoomID     UUID      `json:"chat_room_id"`
}

// OrganizationWithDetail 組織を表す構造体。
type OrganizationWithDetail struct {
	OrganizationID uuid.UUID             `json:"organization_id"`
	Name           string                `json:"name"`
	Color          String                `json:"color"`
	Description    String                `json:"description"`
	IsPersonal     bool                  `json:"is_personal"`
	IsWhole        bool                  `json:"is_whole"`
	ChatRoomID     UUID                  `json:"chat_room_id"`
	Group          NullableEntity[Group] `json:"group,omitempty"`
	Grade          NullableEntity[Grade] `json:"grade,omitempty"`
}

// OrganizationWithDetailForQuery 組織を表す構造体(クエリー用)。
type OrganizationWithDetailForQuery struct {
	Pkey Int `json:"-"`
	OrganizationWithDetail
}

// OrganizationWithChatRoom 組織と専用チャットルームを表す構造体。
type OrganizationWithChatRoom struct {
	OrganizationID uuid.UUID              `json:"organization_id"`
	Name           string                 `json:"name"`
	Color          String                 `json:"color"`
	Description    String                 `json:"description"`
	IsPersonal     bool                   `json:"is_personal"`
	IsWhole        bool                   `json:"is_whole"`
	ChatRoom       ChatRoomWithCoverImage `json:"chat_room"`
}

// OrganizationWithChatRoomForQuery 組織と専用チャットルームを表す構造体(クエリー用)。
type OrganizationWithChatRoomForQuery struct {
	Pkey Int `json:"-"`
	OrganizationWithChatRoom
}

// OrganizationWithChatRoomAndDetail 組織と専用チャットルーム、詳細を表す構造体。
type OrganizationWithChatRoomAndDetail struct {
	OrganizationID uuid.UUID              `json:"organization_id"`
	Name           string                 `json:"name"`
	Color          String                 `json:"color"`
	Description    String                 `json:"description"`
	IsPersonal     bool                   `json:"is_personal"`
	IsWhole        bool                   `json:"is_whole"`
	ChatRoom       ChatRoomWithCoverImage `json:"chat_room"`
	Group          NullableEntity[Group]  `json:"group,omitempty"`
	Grade          NullableEntity[Grade]  `json:"grade,omitempty"`
}

// OrganizationWithChatRoomAndDetailForQuery 組織と専用チャットルーム、詳細を表す構造体(クエリー用)。
type OrganizationWithChatRoomAndDetailForQuery struct {
	Pkey Int `json:"-"`
	OrganizationWithChatRoomAndDetail
}

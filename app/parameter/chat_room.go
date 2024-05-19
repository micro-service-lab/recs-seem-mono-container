package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomParam チャットルーム作成のパラメータ。
type CreateChatRoomParam struct {
	Name             string
	IsPrivate        bool
	CoverImageID     entity.UUID
	OwnerID          entity.UUID
	FromOrganization bool
}

// UpdateChatRoomParams チャットルーム更新のパラメータ。
type UpdateChatRoomParams struct {
	Name         string
	IsPrivate    bool
	CoverImageID entity.UUID
	OwnerID      entity.UUID
}

// WhereChatRoomParam チャットルーム検索のパラメータ。
type WhereChatRoomParam struct {
	WhereInOwner            bool
	InOwner                 []uuid.UUID
	WhereIsPrivate          bool
	IsPrivate               bool
	WhereIsFromOrganization bool
	IsFromOrganization      bool
	WhereFromOrganizations  bool
	FromOrganizations       []uuid.UUID
}

// ChatRoomOrderMethod チャットルームの並び替え方法。
type ChatRoomOrderMethod string

// ParseChatRoomOrderMethod はチャットルームの並び替え方法をパースする。
func ParseChatRoomOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomOrderMethodDefault):
		return ChatRoomOrderMethodDefault, nil
	default:
		return ChatRoomOrderMethodDefault, nil
	}
}

const (
	// ChatRoomDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomOrderMethodDefault:
		return ChatRoomDefaultCursorKey
	default:
		return ChatRoomDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomOrderMethodDefault はデフォルト。
	ChatRoomOrderMethodDefault ChatRoomOrderMethod = "default"
)

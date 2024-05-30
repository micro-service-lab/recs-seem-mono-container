package parameter

import (
	"time"

	"github.com/google/uuid"
)

// CreateChatRoomActionParam チャットルームアクション作成のパラメータ。
type CreateChatRoomActionParam struct {
	ChatRoomID           uuid.UUID
	ChatRoomActionTypeID uuid.UUID
	ActedAt              time.Time
}

// WhereChatRoomActionParam チャットルームアクション検索のパラメータ。
type WhereChatRoomActionParam struct {
	WhereInChatRoomActionTypeIDs bool
	InChatRoomActionTypeIDs      []uuid.UUID
}

// ChatRoomActionOrderMethod チャットルームアクションの並び替え方法。
type ChatRoomActionOrderMethod string

// ParseChatRoomActionOrderMethod はチャットルームアクションの並び替え方法をパースする。
func ParseChatRoomActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomActionOrderMethodDefault):
		return ChatRoomActionOrderMethodDefault, nil
	case string(ChatRoomActionOrderMethodActedAt):
		return ChatRoomActionOrderMethodActedAt, nil
	case string(ChatRoomActionOrderMethodReverseActedAt):
		return ChatRoomActionOrderMethodReverseActedAt, nil
	default:
		return ChatRoomActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomActionDefaultCursorKey = "default"
	// ChatRoomActionActedAtCursorKey はアクション日時カーソルキー。
	ChatRoomActionActedAtCursorKey = "acted_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomActionOrderMethodDefault:
		return ChatRoomActionDefaultCursorKey
	case ChatRoomActionOrderMethodActedAt, ChatRoomActionOrderMethodReverseActedAt:
		return ChatRoomActionActedAtCursorKey
	default:
		return ChatRoomActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomActionOrderMethodDefault はデフォルト。
	ChatRoomActionOrderMethodDefault ChatRoomActionOrderMethod = "default"
	// ChatRoomActionOrderMethodActedAt はアクション日時順。
	ChatRoomActionOrderMethodActedAt ChatRoomActionOrderMethod = "acted_at"
	// ChatRoomActionOrderMethodReverseActedAt はアクション日時逆順。
	ChatRoomActionOrderMethodReverseActedAt ChatRoomActionOrderMethod = "r_acted_at"
)

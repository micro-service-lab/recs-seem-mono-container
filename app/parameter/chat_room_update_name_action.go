package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomUpdateNameActionParam チャットルーム名前更新アクション作成のパラメータ。
type CreateChatRoomUpdateNameActionParam struct {
	ChatRoomActionID uuid.UUID
	UpdatedBy        entity.UUID
	Name             string
}

// WhereChatRoomUpdateNameActionParam チャットルーム名前更新アクション検索のパラメータ。
type WhereChatRoomUpdateNameActionParam struct{}

// ChatRoomUpdateNameActionOrderMethod チャットルーム名前更新アクションの並び替え方法。
type ChatRoomUpdateNameActionOrderMethod string

// ParseChatRoomUpdateNameActionOrderMethod はチャットルーム名前更新アクションの並び替え方法をパースする。
func ParseChatRoomUpdateNameActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomUpdateNameActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomUpdateNameActionOrderMethodDefault):
		return ChatRoomUpdateNameActionOrderMethodDefault, nil
	default:
		return ChatRoomUpdateNameActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomUpdateNameActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomUpdateNameActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomUpdateNameActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomUpdateNameActionOrderMethodDefault:
		return ChatRoomUpdateNameActionDefaultCursorKey
	default:
		return ChatRoomUpdateNameActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomUpdateNameActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomUpdateNameActionOrderMethodDefault はデフォルト。
	ChatRoomUpdateNameActionOrderMethodDefault ChatRoomUpdateNameActionOrderMethod = "default"
)

package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomCreateActionParam チャットルーム作成アクション作成のパラメータ。
type CreateChatRoomCreateActionParam struct {
	ChatRoomActionID uuid.UUID
	CreatedBy        entity.UUID
	Name             string
}

// WhereChatRoomCreateActionParam チャットルーム作成アクション検索のパラメータ。
type WhereChatRoomCreateActionParam struct{}

// ChatRoomCreateActionOrderMethod チャットルーム作成アクションの並び替え方法。
type ChatRoomCreateActionOrderMethod string

// ParseChatRoomCreateActionOrderMethod はチャットルーム作成アクションの並び替え方法をパースする。
func ParseChatRoomCreateActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomCreateActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomCreateActionOrderMethodDefault):
		return ChatRoomCreateActionOrderMethodDefault, nil
	default:
		return ChatRoomCreateActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomCreateActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomCreateActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomCreateActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomCreateActionOrderMethodDefault:
		return ChatRoomCreateActionDefaultCursorKey
	default:
		return ChatRoomCreateActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomCreateActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomCreateActionOrderMethodDefault はデフォルト。
	ChatRoomCreateActionOrderMethodDefault ChatRoomCreateActionOrderMethod = "default"
)

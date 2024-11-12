package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomDeleteMessageActionParam チャットルームメッセージ削除アクション作成のパラメータ。
type CreateChatRoomDeleteMessageActionParam struct {
	ChatRoomActionID uuid.UUID
	DeletedBy        entity.UUID
}

// WhereChatRoomDeleteMessageActionParam チャットルームメッセージ削除アクション検索のパラメータ。
type WhereChatRoomDeleteMessageActionParam struct{}

// ChatRoomDeleteMessageActionOrderMethod チャットルームメッセージ削除アクションの並び替え方法。
type ChatRoomDeleteMessageActionOrderMethod string

// ParseChatRoomDeleteMessageActionOrderMethod はチャットルームメッセージ削除アクションの並び替え方法をパースする。
func ParseChatRoomDeleteMessageActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomDeleteMessageActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomDeleteMessageActionOrderMethodDefault):
		return ChatRoomDeleteMessageActionOrderMethodDefault, nil
	default:
		return ChatRoomDeleteMessageActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomDeleteMessageActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomDeleteMessageActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomDeleteMessageActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomDeleteMessageActionOrderMethodDefault:
		return ChatRoomDeleteMessageActionDefaultCursorKey
	default:
		return ChatRoomDeleteMessageActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomDeleteMessageActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomDeleteMessageActionOrderMethodDefault はデフォルト。
	ChatRoomDeleteMessageActionOrderMethodDefault ChatRoomDeleteMessageActionOrderMethod = "default"
)

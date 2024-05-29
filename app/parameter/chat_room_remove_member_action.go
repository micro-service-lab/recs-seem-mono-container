package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomRemoveMemberActionParam チャットルームメンバー強制退会アクション作成のパラメータ。
type CreateChatRoomRemoveMemberActionParam struct {
	ChatRoomActionID uuid.UUID
	RemovedBy        entity.UUID
}

// WhereChatRoomRemoveMemberActionParam チャットルームメンバー強制退会アクション検索のパラメータ。
type WhereChatRoomRemoveMemberActionParam struct{}

// ChatRoomRemoveMemberActionOrderMethod チャットルームメンバー強制退会アクションの並び替え方法。
type ChatRoomRemoveMemberActionOrderMethod string

// ParseChatRoomRemoveMemberActionOrderMethod はチャットルームメンバー強制退会アクションの並び替え方法をパースする。
func ParseChatRoomRemoveMemberActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomRemoveMemberActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomRemoveMemberActionOrderMethodDefault):
		return ChatRoomRemoveMemberActionOrderMethodDefault, nil
	default:
		return ChatRoomRemoveMemberActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomRemoveMemberActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomRemoveMemberActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomRemoveMemberActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomRemoveMemberActionOrderMethodDefault:
		return ChatRoomRemoveMemberActionDefaultCursorKey
	default:
		return ChatRoomRemoveMemberActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomRemoveMemberActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomRemoveMemberActionOrderMethodDefault はデフォルト。
	ChatRoomRemoveMemberActionOrderMethodDefault ChatRoomRemoveMemberActionOrderMethod = "default"
)

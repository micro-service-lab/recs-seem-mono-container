package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomAddMemberActionParam チャットルームメンバー追加アクション作成のパラメータ。
type CreateChatRoomAddMemberActionParam struct {
	ChatRoomActionID uuid.UUID
	AddedBy          entity.UUID
}

// WhereChatRoomAddMemberActionParam チャットルームメンバー追加アクション検索のパラメータ。
type WhereChatRoomAddMemberActionParam struct{}

// ChatRoomAddMemberActionOrderMethod チャットルームメンバー追加アクションの並び替え方法。
type ChatRoomAddMemberActionOrderMethod string

// ParseChatRoomAddMemberActionOrderMethod はチャットルームメンバー追加アクションの並び替え方法をパースする。
func ParseChatRoomAddMemberActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomAddMemberActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomAddMemberActionOrderMethodDefault):
		return ChatRoomAddMemberActionOrderMethodDefault, nil
	default:
		return ChatRoomAddMemberActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomAddMemberActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomAddMemberActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomAddMemberActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomAddMemberActionOrderMethodDefault:
		return ChatRoomAddMemberActionDefaultCursorKey
	default:
		return ChatRoomAddMemberActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomAddMemberActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomAddMemberActionOrderMethodDefault はデフォルト。
	ChatRoomAddMemberActionOrderMethodDefault ChatRoomAddMemberActionOrderMethod = "default"
)

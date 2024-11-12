package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomWithdrawActionParam チャットルームメンバー脱退アクション作成のパラメータ。
type CreateChatRoomWithdrawActionParam struct {
	ChatRoomActionID uuid.UUID
	MemberID         entity.UUID
}

// WhereChatRoomWithdrawActionParam チャットルームメンバー脱退アクション検索のパラメータ。
type WhereChatRoomWithdrawActionParam struct{}

// ChatRoomWithdrawActionOrderMethod チャットルームメンバー脱退アクションの並び替え方法。
type ChatRoomWithdrawActionOrderMethod string

// ParseChatRoomWithdrawActionOrderMethod はチャットルームメンバー脱退アクションの並び替え方法をパースする。
func ParseChatRoomWithdrawActionOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomWithdrawActionOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomWithdrawActionOrderMethodDefault):
		return ChatRoomWithdrawActionOrderMethodDefault, nil
	default:
		return ChatRoomWithdrawActionOrderMethodDefault, nil
	}
}

const (
	// ChatRoomWithdrawActionDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomWithdrawActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomWithdrawActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomWithdrawActionOrderMethodDefault:
		return ChatRoomWithdrawActionDefaultCursorKey
	default:
		return ChatRoomWithdrawActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomWithdrawActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomWithdrawActionOrderMethodDefault はデフォルト。
	ChatRoomWithdrawActionOrderMethodDefault ChatRoomWithdrawActionOrderMethod = "default"
)

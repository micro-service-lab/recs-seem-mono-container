package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomAddedMemberParam チャットルーム追加メンバーのパラメータ。
type CreateChatRoomAddedMemberParam struct {
	ChatRoomAddMemberActionID uuid.UUID
	MemberID                  entity.UUID
}

// WhereMemberOnChatRoomAddMemberActionParam アクション上のメンバー検索のパラメータ。
type WhereMemberOnChatRoomAddMemberActionParam struct{}

// MemberOnChatRoomAddMemberActionOrderMethod アクション上のメンバーの並び替え方法。
type MemberOnChatRoomAddMemberActionOrderMethod string

// ParseMemberOnChatRoomAddMemberActionOrderMethod はアクション上のメンバーの並び替え方法をパースする。
func ParseMemberOnChatRoomAddMemberActionOrderMethod(v string) (any, error) {
	if v == "" {
		return MemberOnChatRoomAddMemberActionOrderMethodDefault, nil
	}
	switch v {
	case string(MemberOnChatRoomAddMemberActionOrderMethodDefault):
		return MemberOnChatRoomAddMemberActionOrderMethodDefault, nil
	default:
		return MemberOnChatRoomAddMemberActionOrderMethodDefault, nil
	}
}

const (
	// MemberOnChatRoomAddMemberActionDefaultCursorKey はデフォルトカーソルキー。
	MemberOnChatRoomAddMemberActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MemberOnChatRoomAddMemberActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case MemberOnChatRoomAddMemberActionOrderMethodDefault:
		return MemberOnChatRoomAddMemberActionDefaultCursorKey
	default:
		return MemberOnChatRoomAddMemberActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MemberOnChatRoomAddMemberActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MemberOnChatRoomAddMemberActionOrderMethodDefault はデフォルト。
	MemberOnChatRoomAddMemberActionOrderMethodDefault MemberOnChatRoomAddMemberActionOrderMethod = "default"
)

package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateChatRoomRemovedMemberParam チャットルーム追放メンバーのパラメータ。
type CreateChatRoomRemovedMemberParam struct {
	ChatRoomRemoveMemberActionID uuid.UUID
	MemberID                     entity.UUID
}

// WhereMemberOnChatRoomRemoveMemberActionParam アクション上のメンバー検索のパラメータ。
type WhereMemberOnChatRoomRemoveMemberActionParam struct{}

// MemberOnChatRoomRemoveMemberActionOrderMethod アクション上のメンバーの並び替え方法。
type MemberOnChatRoomRemoveMemberActionOrderMethod string

// ParseMemberOnChatRoomRemoveMemberActionOrderMethod はアクション上のメンバーの並び替え方法をパースする。
func ParseMemberOnChatRoomRemoveMemberActionOrderMethod(v string) (any, error) {
	if v == "" {
		return MemberOnChatRoomRemoveMemberActionOrderMethodDefault, nil
	}
	switch v {
	case string(MemberOnChatRoomRemoveMemberActionOrderMethodDefault):
		return MemberOnChatRoomRemoveMemberActionOrderMethodDefault, nil
	default:
		return MemberOnChatRoomRemoveMemberActionOrderMethodDefault, nil
	}
}

const (
	// MemberOnChatRoomRemoveMemberActionDefaultCursorKey はデフォルトカーソルキー。
	MemberOnChatRoomRemoveMemberActionDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MemberOnChatRoomRemoveMemberActionOrderMethod) GetCursorKeyName() string {
	switch m {
	case MemberOnChatRoomRemoveMemberActionOrderMethodDefault:
		return MemberOnChatRoomRemoveMemberActionDefaultCursorKey
	default:
		return MemberOnChatRoomRemoveMemberActionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MemberOnChatRoomRemoveMemberActionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MemberOnChatRoomRemoveMemberActionOrderMethodDefault はデフォルト。
	MemberOnChatRoomRemoveMemberActionOrderMethodDefault MemberOnChatRoomRemoveMemberActionOrderMethod = "default"
)

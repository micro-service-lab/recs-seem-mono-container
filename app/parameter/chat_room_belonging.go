package parameter

import (
	"time"

	"github.com/google/uuid"
)

// BelongChatRoomParam チャットルーム所属のパラメータ。
type BelongChatRoomParam struct {
	MemberID   uuid.UUID
	ChatRoomID uuid.UUID
	AddedAt    time.Time
}

// WhereChatRoomOnMemberParam メンバー上のチャットルーム検索のパラメータ。
type WhereChatRoomOnMemberParam struct {
	WhereLikeName bool
	SearchName    string
}

// ChatRoomOnMemberOrderMethod メンバー上のチャットルームの並び替え方法。
type ChatRoomOnMemberOrderMethod string

// ParseChatRoomOnMemberOrderMethod はメンバー上のチャットルームの並び替え方法をパースする。
func ParseChatRoomOnMemberOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomOnMemberOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomOnMemberOrderMethodDefault):
		return ChatRoomOnMemberOrderMethodDefault, nil
	case string(ChatRoomOnMemberOrderMethodName):
		return ChatRoomOnMemberOrderMethodName, nil
	case string(ChatRoomOnMemberOrderMethodReverseName):
		return ChatRoomOnMemberOrderMethodReverseName, nil
	case string(ChatRoomOnMemberOrderMethodOldAdd):
		return ChatRoomOnMemberOrderMethodOldAdd, nil
	case string(ChatRoomOnMemberOrderMethodLateAdd):
		return ChatRoomOnMemberOrderMethodLateAdd, nil
	case string(ChatRoomOnMemberOrderMethodOldChat):
		return ChatRoomOnMemberOrderMethodOldChat, nil
	case string(ChatRoomOnMemberOrderMethodLateChat):
		return ChatRoomOnMemberOrderMethodLateChat, nil
	default:
		return ChatRoomOnMemberOrderMethodDefault, nil
	}
}

const (
	// ChatRoomOnMemberDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomOnMemberDefaultCursorKey = "default"
	// ChatRoomOnMemberNameCursorKey は名前カーソルキー。
	ChatRoomOnMemberNameCursorKey = "name"
	// ChatRoomOnMemberAddedAtCursorKey は追加日時カーソルキー。
	ChatRoomOnMemberAddedAtCursorKey = "added_at"
	// ChatRoomOnMemberLastChatAtCursorKey は最終チャット日時カーソルキー。
	ChatRoomOnMemberLastChatAtCursorKey = "last_chat_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomOnMemberOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomOnMemberOrderMethodDefault:
		return ChatRoomOnMemberDefaultCursorKey
	case ChatRoomOnMemberOrderMethodName, ChatRoomOnMemberOrderMethodReverseName:
		return ChatRoomOnMemberNameCursorKey
	case ChatRoomOnMemberOrderMethodOldAdd, ChatRoomOnMemberOrderMethodLateAdd:
		return ChatRoomOnMemberAddedAtCursorKey
	case ChatRoomOnMemberOrderMethodOldChat, ChatRoomOnMemberOrderMethodLateChat:
		return ChatRoomOnMemberLastChatAtCursorKey
	default:
		return ChatRoomOnMemberDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomOnMemberOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomOnMemberOrderMethodDefault はデフォルト。
	ChatRoomOnMemberOrderMethodDefault ChatRoomOnMemberOrderMethod = "default"
	// ChatRoomOnMemberOrderMethodName は名前順。
	ChatRoomOnMemberOrderMethodName ChatRoomOnMemberOrderMethod = "name"
	// ChatRoomOnMemberOrderMethodReverseName は名前逆順。
	ChatRoomOnMemberOrderMethodReverseName ChatRoomOnMemberOrderMethod = "r_name"
	// ChatRoomOnMemberOrderMethodOldAdd は追加古い順。
	ChatRoomOnMemberOrderMethodOldAdd ChatRoomOnMemberOrderMethod = "old_add"
	// ChatRoomOnMemberOrderMethodLateAdd は追加新しい順。
	ChatRoomOnMemberOrderMethodLateAdd ChatRoomOnMemberOrderMethod = "late_add"
	// ChatRoomOnMemberOrderMethodOldChat は最終チャット古い順。
	ChatRoomOnMemberOrderMethodOldChat ChatRoomOnMemberOrderMethod = "old_chat"
	// ChatRoomOnMemberOrderMethodLateChat は最終チャット新しい順。
	ChatRoomOnMemberOrderMethodLateChat ChatRoomOnMemberOrderMethod = "late_chat"
)

// WhereMemberOnChatRoomParam チャットルーム上のメンバー検索のパラメータ。
type WhereMemberOnChatRoomParam struct {
	WhereLikeName bool
	SearchName    string
}

// MemberOnChatRoomOrderMethod チャットルーム上のメンバーの並び替え方法。
type MemberOnChatRoomOrderMethod string

// ParseMemberOnChatRoomOrderMethod はチャットルーム上のメンバーの並び替え方法をパースする。
func ParseMemberOnChatRoomOrderMethod(v string) (any, error) {
	if v == "" {
		return MemberOnChatRoomOrderMethodDefault, nil
	}
	switch v {
	case string(MemberOnChatRoomOrderMethodDefault):
		return MemberOnChatRoomOrderMethodDefault, nil
	case string(MemberOnChatRoomOrderMethodName):
		return MemberOnChatRoomOrderMethodName, nil
	case string(MemberOnChatRoomOrderMethodReverseName):
		return MemberOnChatRoomOrderMethodReverseName, nil
	case string(MemberOnChatRoomOrderMethodOldAdd):
		return MemberOnChatRoomOrderMethodOldAdd, nil
	case string(MemberOnChatRoomOrderMethodLateAdd):
		return MemberOnChatRoomOrderMethodLateAdd, nil
	default:
		return MemberOnChatRoomOrderMethodDefault, nil
	}
}

const (
	// MemberOnChatRoomDefaultCursorKey はデフォルトカーソルキー。
	MemberOnChatRoomDefaultCursorKey = "default"
	// MemberOnChatRoomNameCursorKey は名前カーソルキー。
	MemberOnChatRoomNameCursorKey = "name"
	// MemberOnChatRoomAddedAtCursorKey は追加日時カーソルキー。
	MemberOnChatRoomAddedAtCursorKey = "added_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MemberOnChatRoomOrderMethod) GetCursorKeyName() string {
	switch m {
	case MemberOnChatRoomOrderMethodDefault:
		return MemberOnChatRoomDefaultCursorKey
	case MemberOnChatRoomOrderMethodName, MemberOnChatRoomOrderMethodReverseName:
		return MemberOnChatRoomNameCursorKey
	case MemberOnChatRoomOrderMethodOldAdd, MemberOnChatRoomOrderMethodLateAdd:
		return MemberOnChatRoomAddedAtCursorKey
	default:
		return MemberOnChatRoomDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MemberOnChatRoomOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MemberOnChatRoomOrderMethodDefault はデフォルト。
	MemberOnChatRoomOrderMethodDefault MemberOnChatRoomOrderMethod = "default"
	// MemberOnChatRoomOrderMethodName は名前順。
	MemberOnChatRoomOrderMethodName MemberOnChatRoomOrderMethod = "name"
	// MemberOnChatRoomOrderMethodReverseName は名前逆順。
	MemberOnChatRoomOrderMethodReverseName MemberOnChatRoomOrderMethod = "r_name"
	// MemberOnChatRoomOrderMethodOldAdd は追加古い順。
	MemberOnChatRoomOrderMethodOldAdd MemberOnChatRoomOrderMethod = "old_add"
	// MemberOnChatRoomOrderMethodLateAdd は追加新しい順。
	MemberOnChatRoomOrderMethodLateAdd MemberOnChatRoomOrderMethod = "late_add"
)

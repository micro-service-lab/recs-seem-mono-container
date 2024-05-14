package parameter

import "github.com/google/uuid"

// CreateChatRoomParam チャットルーム作成のパラメータ。
type CreateChatRoomParam struct {
	Name             string
	IsPrivate        bool
	CoverImageURL    string
	OwnerID          uuid.UUID
	FromOrganization bool
}

// UpdateChatRoomParams チャットルーム更新のパラメータ。
type UpdateChatRoomParams struct {
	Name          string
	IsPrivate     bool
	CoverImageURL string
	OwnerID       uuid.UUID
}

// WhereChatRoomParam チャットルーム検索のパラメータ。
type WhereChatRoomParam struct {
	WhereInOwner   bool
	InOwners       []uuid.UUID
	WhereLikeName  bool
	SearchName     string
	WhereIsPrivate bool
	IsPrivate      bool
}

// ChatRoomOrderMethod チャットルームの並び替え方法。
type ChatRoomOrderMethod string

// ParseChatRoomOrderMethod はチャットルームの並び替え方法をパースする。
func ParseChatRoomOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomOrderMethodDefault):
		return ChatRoomOrderMethodDefault, nil
	case string(ChatRoomOrderMethodName):
		return ChatRoomOrderMethodName, nil
	case string(ChatRoomOrderMethodReverseName):
		return ChatRoomOrderMethodReverseName, nil
	default:
		return ChatRoomOrderMethodDefault, nil
	}
}

const (
	// ChatRoomDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomDefaultCursorKey = "default"
	// ChatRoomNameCursorKey は名前カーソルキー。
	ChatRoomNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomOrderMethodDefault:
		return ChatRoomDefaultCursorKey
	case ChatRoomOrderMethodName:
		return ChatRoomNameCursorKey
	case ChatRoomOrderMethodReverseName:
		return ChatRoomNameCursorKey
	default:
		return ChatRoomDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomOrderMethodDefault はデフォルト。
	ChatRoomOrderMethodDefault ChatRoomOrderMethod = "default"
	// ChatRoomOrderMethodName は名前順。
	ChatRoomOrderMethodName ChatRoomOrderMethod = "name"
	// ChatRoomOrderMethodReverseName は名前逆順。
	ChatRoomOrderMethodReverseName ChatRoomOrderMethod = "r_name"
)

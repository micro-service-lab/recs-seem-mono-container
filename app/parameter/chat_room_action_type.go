package parameter

// CreateChatRoomActionTypeParam チャットルームアクションタイプ作成のパラメータ。
type CreateChatRoomActionTypeParam struct {
	Name string
	Key  string
}

// UpdateChatRoomActionTypeParams チャットルームアクションタイプ更新のパラメータ。
type UpdateChatRoomActionTypeParams struct {
	Name string
	Key  string
}

// UpdateChatRoomActionTypeByKeyParams チャットルームアクションタイプ更新のパラメータ。
type UpdateChatRoomActionTypeByKeyParams struct {
	Name string
}

// WhereChatRoomActionTypeParam チャットルームアクションタイプ検索のパラメータ。
type WhereChatRoomActionTypeParam struct {
	WhereLikeName bool
	SearchName    string
}

// ChatRoomActionTypeOrderMethod チャットルームアクションタイプの並び替え方法。
type ChatRoomActionTypeOrderMethod string

// ParseChatRoomActionTypeOrderMethod はチャットルームアクションタイプの並び替え方法をパースする。
func ParseChatRoomActionTypeOrderMethod(v string) (any, error) {
	if v == "" {
		return ChatRoomActionTypeOrderMethodDefault, nil
	}
	switch v {
	case string(ChatRoomActionTypeOrderMethodDefault):
		return ChatRoomActionTypeOrderMethodDefault, nil
	case string(ChatRoomActionTypeOrderMethodName):
		return ChatRoomActionTypeOrderMethodName, nil
	case string(ChatRoomActionTypeOrderMethodReverseName):
		return ChatRoomActionTypeOrderMethodReverseName, nil
	default:
		return ChatRoomActionTypeOrderMethodDefault, nil
	}
}

const (
	// ChatRoomActionTypeDefaultCursorKey はデフォルトカーソルキー。
	ChatRoomActionTypeDefaultCursorKey = "default"
	// ChatRoomActionTypeNameCursorKey は名前カーソルキー。
	ChatRoomActionTypeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ChatRoomActionTypeOrderMethod) GetCursorKeyName() string {
	switch m {
	case ChatRoomActionTypeOrderMethodDefault:
		return ChatRoomActionTypeDefaultCursorKey
	case ChatRoomActionTypeOrderMethodName:
		return ChatRoomActionTypeNameCursorKey
	case ChatRoomActionTypeOrderMethodReverseName:
		return ChatRoomActionTypeNameCursorKey
	default:
		return ChatRoomActionTypeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ChatRoomActionTypeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ChatRoomActionTypeOrderMethodDefault はデフォルト。
	ChatRoomActionTypeOrderMethodDefault ChatRoomActionTypeOrderMethod = "default"
	// ChatRoomActionTypeOrderMethodName は名前順。
	ChatRoomActionTypeOrderMethodName ChatRoomActionTypeOrderMethod = "name"
	// ChatRoomActionTypeOrderMethodReverseName は名前逆順。
	ChatRoomActionTypeOrderMethodReverseName ChatRoomActionTypeOrderMethod = "r_name"
)

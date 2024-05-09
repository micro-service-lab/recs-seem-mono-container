package parameter

// CreateEventTypeParam イベントタイプ作成のパラメータ。
type CreateEventTypeParam struct {
	Name  string
	Key   string
	Color string
}

// UpdateEventTypeParams イベントタイプ更新のパラメータ。
type UpdateEventTypeParams struct {
	Name  string
	Key   string
	Color string
}

// UpdateEventTypeByKeyParams イベントタイプ更新のパラメータ。
type UpdateEventTypeByKeyParams struct {
	Name  string
	Color string
}

// WhereEventTypeParam イベントタイプ検索のパラメータ。
type WhereEventTypeParam struct {
	WhereLikeName bool
	SearchName    string
}

// EventTypeOrderMethod イベントタイプの並び替え方法。
type EventTypeOrderMethod string

// ParseEventTypeOrderMethod はイベントタイプの並び替え方法をパースする。
func ParseEventTypeOrderMethod(v string) (any, error) {
	if v == "" {
		return EventTypeOrderMethodDefault, nil
	}
	switch v {
	case string(EventTypeOrderMethodDefault):
		return EventTypeOrderMethodDefault, nil
	case string(EventTypeOrderMethodName):
		return EventTypeOrderMethodName, nil
	case string(EventTypeOrderMethodReverseName):
		return EventTypeOrderMethodReverseName, nil
	default:
		return EventTypeOrderMethodDefault, nil
	}
}

const (
	// EventTypeDefaultCursorKey はデフォルトカーソルキー。
	EventTypeDefaultCursorKey = "default"
	// EventTypeNameCursorKey は名前カーソルキー。
	EventTypeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m EventTypeOrderMethod) GetCursorKeyName() string {
	switch m {
	case EventTypeOrderMethodDefault:
		return EventTypeDefaultCursorKey
	case EventTypeOrderMethodName:
		return EventTypeNameCursorKey
	case EventTypeOrderMethodReverseName:
		return EventTypeNameCursorKey
	default:
		return EventTypeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m EventTypeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// EventTypeOrderMethodDefault はデフォルト。
	EventTypeOrderMethodDefault EventTypeOrderMethod = "default"
	// EventTypeOrderMethodName は名前順。
	EventTypeOrderMethodName EventTypeOrderMethod = "name"
	// EventTypeOrderMethodReverseName は名前逆順。
	EventTypeOrderMethodReverseName EventTypeOrderMethod = "r_name"
)

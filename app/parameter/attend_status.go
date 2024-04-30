package parameter

// CreateAttendStatusParam 出席ステータス作成のパラメータ。
type CreateAttendStatusParam struct {
	Name string
	Key  string
}

// UpdateAttendStatusParams 出席ステータス更新のパラメータ。
type UpdateAttendStatusParams struct {
	Name string
	Key  string
}

// UpdateAttendStatusByKeyParams 出席ステータス更新のパラメータ。
type UpdateAttendStatusByKeyParams struct {
	Name string
}

// WhereAttendStatusParam 出席ステータス検索のパラメータ。
type WhereAttendStatusParam struct {
	WhereLikeName bool
	SearchName    string
}

// AttendStatusOrderMethod 出席ステータスの並び替え方法。
type AttendStatusOrderMethod string

// ParseAttendStatusOrderMethod は出席ステータスの並び替え方法をパースする。
func ParseAttendStatusOrderMethod(v string) (any, error) {
	if v == "" {
		return AttendStatusOrderMethodDefault, nil
	}
	switch v {
	case string(AttendStatusOrderMethodDefault):
		return AttendStatusOrderMethodDefault, nil
	case string(AttendStatusOrderMethodName):
		return AttendStatusOrderMethodName, nil
	case string(AttendStatusOrderMethodReverseName):
		return AttendStatusOrderMethodReverseName, nil
	default:
		return AttendStatusOrderMethodDefault, nil
	}
}

const (
	// AttendStatusDefaultCursorKey はデフォルトカーソルキー。
	AttendStatusDefaultCursorKey = "default"
	// AttendStatusNameCursorKey は名前カーソルキー。
	AttendStatusNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AttendStatusOrderMethod) GetCursorKeyName() string {
	switch m {
	case AttendStatusOrderMethodDefault:
		return AttendStatusDefaultCursorKey
	case AttendStatusOrderMethodName:
		return AttendStatusNameCursorKey
	case AttendStatusOrderMethodReverseName:
		return AttendStatusNameCursorKey
	default:
		return AttendStatusDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AttendStatusOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AttendStatusOrderMethodDefault はデフォルト。
	AttendStatusOrderMethodDefault AttendStatusOrderMethod = "default"
	// AttendStatusOrderMethodName は名前順。
	AttendStatusOrderMethodName AttendStatusOrderMethod = "name"
	// AttendStatusOrderMethodReverseName は名前逆順。
	AttendStatusOrderMethodReverseName AttendStatusOrderMethod = "r_name"
)

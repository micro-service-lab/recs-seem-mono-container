package parameter

// CreateAttendanceTypeParam 出席ステータス作成のパラメータ。
type CreateAttendanceTypeParam struct {
	Name  string
	Key   string
	Color string
}

// UpdateAttendanceTypeParams 出席ステータス更新のパラメータ。
type UpdateAttendanceTypeParams struct {
	Name  string
	Key   string
	Color string
}

// UpdateAttendanceTypeByKeyParams 出席ステータス更新のパラメータ。
type UpdateAttendanceTypeByKeyParams struct {
	Name  string
	Color string
}

// WhereAttendanceTypeParam 出席ステータス検索のパラメータ。
type WhereAttendanceTypeParam struct {
	WhereLikeName bool
	SearchName    string
}

// AttendanceTypeOrderMethod 出席ステータスの並び替え方法。
type AttendanceTypeOrderMethod string

// ParseAttendanceTypeOrderMethod は出席ステータスの並び替え方法をパースする。
func ParseAttendanceTypeOrderMethod(v string) (any, error) {
	if v == "" {
		return AttendanceTypeOrderMethodDefault, nil
	}
	switch v {
	case string(AttendanceTypeOrderMethodDefault):
		return AttendanceTypeOrderMethodDefault, nil
	case string(AttendanceTypeOrderMethodName):
		return AttendanceTypeOrderMethodName, nil
	case string(AttendanceTypeOrderMethodReverseName):
		return AttendanceTypeOrderMethodReverseName, nil
	default:
		return AttendanceTypeOrderMethodDefault, nil
	}
}

const (
	// AttendanceTypeDefaultCursorKey はデフォルトカーソルキー。
	AttendanceTypeDefaultCursorKey = "default"
	// AttendanceTypeNameCursorKey は名前カーソルキー。
	AttendanceTypeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AttendanceTypeOrderMethod) GetCursorKeyName() string {
	switch m {
	case AttendanceTypeOrderMethodDefault:
		return AttendanceTypeDefaultCursorKey
	case AttendanceTypeOrderMethodName:
		return AttendanceTypeNameCursorKey
	case AttendanceTypeOrderMethodReverseName:
		return AttendanceTypeNameCursorKey
	default:
		return AttendanceTypeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AttendanceTypeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AttendanceTypeOrderMethodDefault はデフォルト。
	AttendanceTypeOrderMethodDefault AttendanceTypeOrderMethod = "default"
	// AttendanceTypeOrderMethodName は名前順。
	AttendanceTypeOrderMethodName AttendanceTypeOrderMethod = "name"
	// AttendanceTypeOrderMethodReverseName は名前逆順。
	AttendanceTypeOrderMethodReverseName AttendanceTypeOrderMethod = "r_name"
)

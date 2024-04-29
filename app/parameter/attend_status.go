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

// ParseGetAttendStatusesOrderParam は出席ステータスの並び替え方法をパースする。
func ParseAttendStatusOrderMethod(v string) (any, error) {
	if v == "" {
		return AttendStatusOrderMethodName, nil
	}
	switch v {
	case string(AttendStatusOrderMethodName):
		return AttendStatusOrderMethodName, nil
	case string(AttendStatusOrderMethodReverseName):
		return AttendStatusOrderMethodReverseName, nil
	default:
		return AttendStatusOrderMethodName, nil
	}
}

const (
	// AttendStatusOrderMethodName は名前順。
	AttendStatusOrderMethodName AttendStatusOrderMethod = "name"
	// AttendStatusOrderMethodReverseName は名前逆順。
	AttendStatusOrderMethodReverseName AttendStatusOrderMethod = "r_name"
)

package parameter

// CreateRecordTypeParam 議事録タイプ作成のパラメータ。
type CreateRecordTypeParam struct {
	Name string
	Key  string
}

// UpdateRecordTypeParams 議事録タイプ更新のパラメータ。
type UpdateRecordTypeParams struct {
	Name string
	Key  string
}

// UpdateRecordTypeByKeyParams 議事録タイプ更新のパラメータ。
type UpdateRecordTypeByKeyParams struct {
	Name string
}

// WhereRecordTypeParam 議事録タイプ検索のパラメータ。
type WhereRecordTypeParam struct {
	WhereLikeName bool
	SearchName    string
}

// RecordTypeOrderMethod 議事録タイプの並び替え方法。
type RecordTypeOrderMethod string

// ParseRecordTypeOrderMethod は議事録タイプの並び替え方法をパースする。
func ParseRecordTypeOrderMethod(v string) (any, error) {
	if v == "" {
		return RecordTypeOrderMethodDefault, nil
	}
	switch v {
	case string(RecordTypeOrderMethodDefault):
		return RecordTypeOrderMethodDefault, nil
	case string(RecordTypeOrderMethodName):
		return RecordTypeOrderMethodName, nil
	case string(RecordTypeOrderMethodReverseName):
		return RecordTypeOrderMethodReverseName, nil
	default:
		return RecordTypeOrderMethodDefault, nil
	}
}

const (
	// RecordTypeDefaultCursorKey はデフォルトカーソルキー。
	RecordTypeDefaultCursorKey = "default"
	// RecordTypeNameCursorKey は名前カーソルキー。
	RecordTypeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m RecordTypeOrderMethod) GetCursorKeyName() string {
	switch m {
	case RecordTypeOrderMethodDefault:
		return RecordTypeDefaultCursorKey
	case RecordTypeOrderMethodName:
		return RecordTypeNameCursorKey
	case RecordTypeOrderMethodReverseName:
		return RecordTypeNameCursorKey
	default:
		return RecordTypeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m RecordTypeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// RecordTypeOrderMethodDefault はデフォルト。
	RecordTypeOrderMethodDefault RecordTypeOrderMethod = "default"
	// RecordTypeOrderMethodName は名前順。
	RecordTypeOrderMethodName RecordTypeOrderMethod = "name"
	// RecordTypeOrderMethodReverseName は名前逆順。
	RecordTypeOrderMethodReverseName RecordTypeOrderMethod = "r_name"
)

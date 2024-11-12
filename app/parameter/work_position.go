package parameter

// CreateWorkPositionParam ワークポジション作成のパラメータ。
type CreateWorkPositionParam struct {
	Name        string
	Description string
}

// UpdateWorkPositionParams ワークポジション更新のパラメータ。
type UpdateWorkPositionParams struct {
	Name        string
	Description string
}

// WhereWorkPositionParam ワークポジション検索のパラメータ。
type WhereWorkPositionParam struct {
	WhereLikeName bool
	SearchName    string
}

// WorkPositionOrderMethod ワークポジションの並び替え方法。
type WorkPositionOrderMethod string

// ParseWorkPositionOrderMethod はワークポジションの並び替え方法をパースする。
func ParseWorkPositionOrderMethod(v string) (any, error) {
	if v == "" {
		return WorkPositionOrderMethodDefault, nil
	}
	switch v {
	case string(WorkPositionOrderMethodDefault):
		return WorkPositionOrderMethodDefault, nil
	case string(WorkPositionOrderMethodName):
		return WorkPositionOrderMethodName, nil
	case string(WorkPositionOrderMethodReverseName):
		return WorkPositionOrderMethodReverseName, nil
	default:
		return WorkPositionOrderMethodDefault, nil
	}
}

const (
	// WorkPositionDefaultCursorKey はデフォルトカーソルキー。
	WorkPositionDefaultCursorKey = "default"
	// WorkPositionNameCursorKey は名前カーソルキー。
	WorkPositionNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m WorkPositionOrderMethod) GetCursorKeyName() string {
	switch m {
	case WorkPositionOrderMethodDefault:
		return WorkPositionDefaultCursorKey
	case WorkPositionOrderMethodName:
		return WorkPositionNameCursorKey
	case WorkPositionOrderMethodReverseName:
		return WorkPositionNameCursorKey
	default:
		return WorkPositionDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m WorkPositionOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// WorkPositionOrderMethodDefault はデフォルト。
	WorkPositionOrderMethodDefault WorkPositionOrderMethod = "default"
	// WorkPositionOrderMethodName は名前順。
	WorkPositionOrderMethodName WorkPositionOrderMethod = "name"
	// WorkPositionOrderMethodReverseName は名前逆順。
	WorkPositionOrderMethodReverseName WorkPositionOrderMethod = "r_name"
)

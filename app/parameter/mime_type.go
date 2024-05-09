package parameter

// CreateMimeTypeParam マイムタイプ作成のパラメータ。
type CreateMimeTypeParam struct {
	Name string
	Key  string
	Kind string
}

// UpdateMimeTypeParams マイムタイプ更新のパラメータ。
type UpdateMimeTypeParams struct {
	Name string
	Key  string
	Kind string
}

// UpdateMimeTypeByKeyParams マイムタイプ更新のパラメータ。
type UpdateMimeTypeByKeyParams struct {
	Name string
	Kind string
}

// WhereMimeTypeParam マイムタイプ検索のパラメータ。
type WhereMimeTypeParam struct {
	WhereLikeName bool
	SearchName    string
}

// MimeTypeOrderMethod マイムタイプの並び替え方法。
type MimeTypeOrderMethod string

// ParseMimeTypeOrderMethod はマイムタイプの並び替え方法をパースする。
func ParseMimeTypeOrderMethod(v string) (any, error) {
	if v == "" {
		return MimeTypeOrderMethodDefault, nil
	}
	switch v {
	case string(MimeTypeOrderMethodDefault):
		return MimeTypeOrderMethodDefault, nil
	case string(MimeTypeOrderMethodName):
		return MimeTypeOrderMethodName, nil
	case string(MimeTypeOrderMethodReverseName):
		return MimeTypeOrderMethodReverseName, nil
	default:
		return MimeTypeOrderMethodDefault, nil
	}
}

const (
	// MimeTypeDefaultCursorKey はデフォルトカーソルキー。
	MimeTypeDefaultCursorKey = "default"
	// MimeTypeNameCursorKey は名前カーソルキー。
	MimeTypeNameCursorKey = "name"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MimeTypeOrderMethod) GetCursorKeyName() string {
	switch m {
	case MimeTypeOrderMethodDefault:
		return MimeTypeDefaultCursorKey
	case MimeTypeOrderMethodName:
		return MimeTypeNameCursorKey
	case MimeTypeOrderMethodReverseName:
		return MimeTypeNameCursorKey
	default:
		return MimeTypeDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MimeTypeOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MimeTypeOrderMethodDefault はデフォルト。
	MimeTypeOrderMethodDefault MimeTypeOrderMethod = "default"
	// MimeTypeOrderMethodName は名前順。
	MimeTypeOrderMethodName MimeTypeOrderMethod = "name"
	// MimeTypeOrderMethodReverseName は名前逆順。
	MimeTypeOrderMethodReverseName MimeTypeOrderMethod = "r_name"
)

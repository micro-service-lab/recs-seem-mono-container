package parameter

import (
	"io"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateImageServiceParam 画像作成のパラメータ。
type CreateImageServiceParam struct {
	Origin io.Reader
	Alias  string
}

// CreateImageSpecifyFilenameServiceParam 画像作成のパラメータ。
type CreateImageSpecifyFilenameServiceParam struct {
	Origin   io.Reader
	Filename string
	Alias    string
}

// CreateImageFromOuterServiceParam 画像作成のパラメータ。
type CreateImageFromOuterServiceParam struct {
	URL        string
	Alias      string
	Size       entity.Float
	MimeTypeID uuid.UUID
	Height     entity.Float
	Width      entity.Float
}

// CreateImageParam 画像作成のパラメータ。
type CreateImageParam struct {
	Height           entity.Float
	Width            entity.Float
	AttachableItemID uuid.UUID
}

// WhereImageParam 画像検索のパラメータ。
type WhereImageParam struct{}

// ImageOrderMethod 画像の並び替え方法。
type ImageOrderMethod string

// ParseImageOrderMethod は画像の並び替え方法をパースする。
func ParseImageOrderMethod(v string) (any, error) {
	if v == "" {
		return ImageOrderMethodDefault, nil
	}
	switch v {
	case string(ImageOrderMethodDefault):
		return ImageOrderMethodDefault, nil
	default:
		return ImageOrderMethodDefault, nil
	}
}

const (
	// ImageDefaultCursorKey はデフォルトカーソルキー。
	ImageDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ImageOrderMethod) GetCursorKeyName() string {
	switch m {
	case ImageOrderMethodDefault:
		return ImageDefaultCursorKey
	default:
		return ImageDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ImageOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ImageOrderMethodDefault はデフォルト。
	ImageOrderMethodDefault ImageOrderMethod = "default"
)

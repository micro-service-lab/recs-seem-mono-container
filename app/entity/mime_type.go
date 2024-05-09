package entity

import "github.com/google/uuid"

// MimeType MIMEタイプを表す構造体。
type MimeType struct {
	MimeTypeID uuid.UUID `json:"mime_type_id"`
	Name       string    `json:"name"`
	Key        string    `json:"key"`
	Kind       string    `json:"kind"`
}

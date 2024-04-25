package entity

import "github.com/google/uuid"

// RecordType レコードタイプを表す構造体。
type RecordType struct {
	RecordTypeID uuid.UUID `json:"record_type_id"`
	Name         string    `json:"name"`
	Key          string    `json:"key"`
}

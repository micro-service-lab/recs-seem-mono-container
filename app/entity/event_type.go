package entity

import "github.com/google/uuid"

// EventType イベント種別を表す構造体。
type EventType struct {
	EventTypeID uuid.UUID `json:"event_type_id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Color       string    `json:"color"`
}

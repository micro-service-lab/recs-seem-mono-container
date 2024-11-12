package ws

import "github.com/google/uuid"

// Targets ターゲットを表す構造体。
type Targets struct {
	Members []uuid.UUID `json:"members"`
	All     bool        `json:"all"`
}

// Payload WebSocket のペイロードを表す構造体。
type Payload struct {
	EventType EventType `json:"event_type"`
	Data      any       `json:"data"`
	Targets   Targets   `json:"targets"`
}

// TypedPayload WebSocket ジェネリクスを使ったペイロードを表す構造体。
type TypedPayload[T any] struct {
	EventType EventType `json:"event_type"`
	Data      T         `json:"data"`
	Targets   Targets   `json:"targets"`
}

// ResponsePayload WebSocket のレスポンスペイロードを表す構造体。
type ResponsePayload struct {
	EventType EventType `json:"event_type"`
	Data      any       `json:"data"`
}

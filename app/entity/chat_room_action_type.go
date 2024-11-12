package entity

import "github.com/google/uuid"

// ChatRoomActionType チャットルームアクションタイプを表す構造体。
type ChatRoomActionType struct {
	ChatRoomActionTypeID uuid.UUID `json:"record_type_id"`
	Name                 string    `json:"name"`
	Key                  string    `json:"key"`
}

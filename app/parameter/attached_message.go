package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// AttachItemMessageParam 添付メッセージ作成のパラメータ。
type AttachItemMessageParam struct {
	MessageID        uuid.UUID
	AttachableItemID entity.UUID
}

// WhereAttachedItemOnMessageParam メッセージ上の添付アイテム検索のパラメータ。
type WhereAttachedItemOnMessageParam struct {
	WhereInMimeType bool
	InMimeTypes     []uuid.UUID
	WhereIsImage    bool
	WhereIsFile     bool
}

// AttachedItemOnMessageOrderMethod メッセージ上の添付アイテムの並び替え方法。
type AttachedItemOnMessageOrderMethod string

// ParseAttachedItemOnMessageOrderMethod はメッセージ上の添付アイテムの並び替え方法をパースする。
func ParseAttachedItemOnMessageOrderMethod(v string) (any, error) {
	if v == "" {
		return AttachedItemOnMessageOrderMethodDefault, nil
	}
	switch v {
	case string(AttachedItemOnMessageOrderMethodDefault):
		return AttachedItemOnMessageOrderMethodDefault, nil
	default:
		return AttachedItemOnMessageOrderMethodDefault, nil
	}
}

const (
	// AttachedItemOnMessageDefaultCursorKey はデフォルトカーソルキー。
	AttachedItemOnMessageDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AttachedItemOnMessageOrderMethod) GetCursorKeyName() string {
	switch m {
	case AttachedItemOnMessageOrderMethodDefault:
		return AttachedItemOnMessageDefaultCursorKey
	default:
		return AttachedItemOnMessageDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AttachedItemOnMessageOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AttachedItemOnMessageOrderMethodDefault はデフォルト。
	AttachedItemOnMessageOrderMethodDefault AttachedItemOnMessageOrderMethod = "default"
)

// WhereAttachedItemOnChatRoomParam チャットルーム上の添付アイテム検索のパラメータ。
type WhereAttachedItemOnChatRoomParam struct {
	WhereInMimeType bool
	InMimeTypes     []uuid.UUID
	WhereIsImage    bool
	WhereIsFile     bool
}

// AttachedItemOnChatRoomOrderMethod チャットルーム上の添付アイテムの並び替え方法。
type AttachedItemOnChatRoomOrderMethod string

// ParseAttachedItemOnChatRoomOrderMethod はチャットルーム上の添付アイテムの並び替え方法をパースする。
func ParseAttachedItemOnChatRoomOrderMethod(v string) (any, error) {
	if v == "" {
		return AttachedItemOnChatRoomOrderMethodDefault, nil
	}
	switch v {
	case string(AttachedItemOnChatRoomOrderMethodDefault):
		return AttachedItemOnChatRoomOrderMethodDefault, nil
	default:
		return AttachedItemOnChatRoomOrderMethodDefault, nil
	}
}

const (
	// AttachedItemOnChatRoomDefaultCursorKey はデフォルトカーソルキー。
	AttachedItemOnChatRoomDefaultCursorKey = "default"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m AttachedItemOnChatRoomOrderMethod) GetCursorKeyName() string {
	switch m {
	case AttachedItemOnChatRoomOrderMethodDefault:
		return AttachedItemOnChatRoomDefaultCursorKey
	default:
		return AttachedItemOnChatRoomDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m AttachedItemOnChatRoomOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// AttachedItemOnChatRoomOrderMethodDefault はデフォルト。
	AttachedItemOnChatRoomOrderMethodDefault AttachedItemOnChatRoomOrderMethod = "default"
)

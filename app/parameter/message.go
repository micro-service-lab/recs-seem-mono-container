package parameter

import (
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateMessageParam メッセージ作成のパラメータ。
type CreateMessageParam struct {
	ChatRoomID uuid.UUID
	SenderID   entity.UUID
	Body       string
	PostedAt   time.Time
}

// UpdateMessageParams メッセージ更新のパラメータ。
type UpdateMessageParams struct {
	Body         string
	LastEditedAt time.Time
}

// WhereMessageParam メッセージ検索のパラメータ。
type WhereMessageParam struct {
	WhereInChatRoom          bool
	InChatRoom               []uuid.UUID
	WhereInSender            bool
	InSender                 []uuid.UUID
	WhereLikeBody            bool
	SearchBody               string
	WhereEarlierPostedAt     bool
	EarlierPostedAt          time.Time
	WhereLaterPostedAt       bool
	LaterPostedAt            time.Time
	WhereEarlierLastEditedAt bool
	EarlierLastEditedAt      time.Time
	WhereLaterLastEditedAt   bool
	LaterLastEditedAt        time.Time
}

// MessageOrderMethod メッセージの並び替え方法。
type MessageOrderMethod string

// ParseMessageOrderMethod はメッセージの並び替え方法をパースする。
func ParseMessageOrderMethod(v string) (any, error) {
	if v == "" {
		return MessageOrderMethodDefault, nil
	}
	switch v {
	case string(MessageOrderMethodDefault):
		return MessageOrderMethodDefault, nil
	case string(MessageOrderMethodPostedAt):
		return MessageOrderMethodPostedAt, nil
	case string(MessageOrderMethodReversePostedAt):
		return MessageOrderMethodReversePostedAt, nil
	case string(MessageOrderMethodLastEditedAt):
		return MessageOrderMethodLastEditedAt, nil
	case string(MessageOrderMethodReverseLastEditedAt):
		return MessageOrderMethodReverseLastEditedAt, nil
	default:
		return MessageOrderMethodDefault, nil
	}
}

const (
	// MessageDefaultCursorKey はデフォルトカーソルキー。
	MessageDefaultCursorKey = "default"
	// MessagePostedAtCursorKey は投稿日時カーソルキー。
	MessagePostedAtCursorKey = "posted_at"
	// MessageLastEditedAtCursorKey は最終編集日時カーソルキー。
	MessageLastEditedAtCursorKey = "last_edited_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m MessageOrderMethod) GetCursorKeyName() string {
	switch m {
	case MessageOrderMethodDefault:
		return MessageDefaultCursorKey
	case MessageOrderMethodPostedAt, MessageOrderMethodReversePostedAt:
		return MessagePostedAtCursorKey
	case MessageOrderMethodLastEditedAt, MessageOrderMethodReverseLastEditedAt:
		return MessageLastEditedAtCursorKey
	default:
		return MessageDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m MessageOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// MessageOrderMethodDefault はデフォルト。
	MessageOrderMethodDefault MessageOrderMethod = "default"
	// MessageOrderMethodPostedAt は投稿日時順。
	MessageOrderMethodPostedAt MessageOrderMethod = "posted_at"
	// MessageOrderMethodReversePostedAt は投稿日時逆順。
	MessageOrderMethodReversePostedAt MessageOrderMethod = "r_posted_at"
	// MessageOrderMethodLastEditedAt は最終編集日時順。
	MessageOrderMethodLastEditedAt MessageOrderMethod = "last_edited_at"
	// MessageOrderMethodReverseLastEditedAt は最終編集日時逆順。
	MessageOrderMethodReverseLastEditedAt MessageOrderMethod = "r_last_edited_at"
)

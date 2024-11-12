package parameter

import (
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateReadReceiptParam 既読情報作成のパラメータ。
type CreateReadReceiptParam struct {
	MemberID  uuid.UUID
	MessageID uuid.UUID
	ReadAt    entity.Timestamptz
}

// ReadReceiptParam 既読情報のパラメータ。
type ReadReceiptParam struct {
	MemberID  uuid.UUID
	MessageID uuid.UUID
	ReadAt    entity.Timestamptz
}

// ReadReceiptsParam 複数既読情報のパラメータ。
type ReadReceiptsParam struct {
	MemberID   uuid.UUID
	MessageIDs []uuid.UUID
	ReadAt     entity.Timestamptz
}

// WhereReadableMemberOnMessageParam メッセージ上の既読情報検索のパラメータ。
type WhereReadableMemberOnMessageParam struct {
	WhereLikeName  bool
	SearchName     string
	WhereIsRead    bool
	WhereIsNotRead bool
}

// WhereExistsReadReceiptParam 既読情報の存在検索のパラメータ。
type WhereExistsReadReceiptParam struct {
	WhereIsRead    bool
	WhereIsNotRead bool
}

// WhereReadableMessageOnChatRoomAndMemberParam メッセージ上のメッセージの既読情報検索のパラメータ。
type WhereReadableMessageOnChatRoomAndMemberParam struct {
	WhereIsRead    bool
	WhereIsNotRead bool
}

// WhereReadsOnMessageParam メッセージ上の既読情報検索のパラメータ。
type WhereReadsOnMessageParam struct {
	WhereIsRead    bool
	WhereIsNotRead bool
}

// ReadableMemberOnMessageOrderMethod メンバー上のメッセージの並び替え方法。
type ReadableMemberOnMessageOrderMethod string

// ParseReadableMemberOnMessageOrderMethod はメンバー上のメッセージの並び替え方法をパースする。
func ParseReadableMemberOnMessageOrderMethod(v string) (any, error) {
	if v == "" {
		return ReadableMemberOnMessageOrderMethodDefault, nil
	}
	switch v {
	case string(ReadableMemberOnMessageOrderMethodDefault):
		return ReadableMemberOnMessageOrderMethodDefault, nil
	case string(ReadableMemberOnMessageOrderMethodName):
		return ReadableMemberOnMessageOrderMethodName, nil
	case string(ReadableMemberOnMessageOrderMethodReverseName):
		return ReadableMemberOnMessageOrderMethodReverseName, nil
	case string(ReadableMemberOnMessageOrderMethodReadAt):
		return ReadableMemberOnMessageOrderMethodReadAt, nil
	case string(ReadableMemberOnMessageOrderMethodReverseReadAt):
		return ReadableMemberOnMessageOrderMethodReverseReadAt, nil
	default:
		return ReadableMemberOnMessageOrderMethodDefault, nil
	}
}

const (
	// ReadableMemberOnMessageDefaultCursorKey はデフォルトカーソルキー。
	ReadableMemberOnMessageDefaultCursorKey = "default"
	// ReadableMemberOnMessageNameCursorKey は名前カーソルキー。
	ReadableMemberOnMessageNameCursorKey = "name"
	// ReadableMemberOnMessageReadAtCursorKey は既読日時カーソルキー。
	ReadableMemberOnMessageReadAtCursorKey = "read_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ReadableMemberOnMessageOrderMethod) GetCursorKeyName() string {
	switch m {
	case ReadableMemberOnMessageOrderMethodDefault:
		return ReadableMemberOnMessageDefaultCursorKey
	case ReadableMemberOnMessageOrderMethodName, ReadableMemberOnMessageOrderMethodReverseName:
		return ReadableMemberOnMessageNameCursorKey
	case ReadableMemberOnMessageOrderMethodReadAt, ReadableMemberOnMessageOrderMethodReverseReadAt:
		return ReadableMemberOnMessageReadAtCursorKey
	default:
		return ReadableMemberOnMessageDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ReadableMemberOnMessageOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ReadableMemberOnMessageOrderMethodDefault はデフォルト。
	ReadableMemberOnMessageOrderMethodDefault ReadableMemberOnMessageOrderMethod = "default"
	// ReadableMemberOnMessageOrderMethodName は名前順。
	ReadableMemberOnMessageOrderMethodName ReadableMemberOnMessageOrderMethod = "name"
	// ReadableMemberOnMessageOrderMethodReverseName は名前逆順。
	ReadableMemberOnMessageOrderMethodReverseName ReadableMemberOnMessageOrderMethod = "r_name"
	// ReadableMemberOnMessageOrderMethodReadAt は既読日時順。
	ReadableMemberOnMessageOrderMethodReadAt ReadableMemberOnMessageOrderMethod = "read_at"
	// ReadableMemberOnMessageOrderMethodReverseReadAt は既読日時逆順。
	ReadableMemberOnMessageOrderMethodReverseReadAt ReadableMemberOnMessageOrderMethod = "r_read_at"
)

// WhereReadableMessageOnMemberParam メンバー上のメッセージの既読情報検索のパラメータ。
type WhereReadableMessageOnMemberParam struct {
	WhereIsRead    bool
	WhereIsNotRead bool
}

// ReadableMessageOnMemberOrderMethod メッセージ上のメンバーの並び替え方法。
type ReadableMessageOnMemberOrderMethod string

// ParseReadableMessageOnMemberOrderMethod はメッセージ上のメンバーの並び替え方法をパースする。
func ParseReadableMessageOnMemberOrderMethod(v string) (any, error) {
	if v == "" {
		return ReadableMessageOnMemberOrderMethodDefault, nil
	}
	switch v {
	case string(ReadableMessageOnMemberOrderMethodDefault):
		return ReadableMessageOnMemberOrderMethodDefault, nil
	case string(ReadableMessageOnMemberOrderMethodReadAt):
		return ReadableMessageOnMemberOrderMethodReadAt, nil
	case string(ReadableMessageOnMemberOrderMethodReverseReadAt):
		return ReadableMessageOnMemberOrderMethodReverseReadAt, nil
	default:
		return ReadableMessageOnMemberOrderMethodDefault, nil
	}
}

const (
	// ReadableMessageOnMemberDefaultCursorKey はデフォルトカーソルキー。
	ReadableMessageOnMemberDefaultCursorKey = "default"
	// ReadableMessageOnMemberReadAtCursorKey は既読日時カーソルキー。
	ReadableMessageOnMemberReadAtCursorKey = "read_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ReadableMessageOnMemberOrderMethod) GetCursorKeyName() string {
	switch m {
	case ReadableMessageOnMemberOrderMethodDefault:
		return ReadableMessageOnMemberDefaultCursorKey
	case ReadableMessageOnMemberOrderMethodReadAt, ReadableMessageOnMemberOrderMethodReverseReadAt:
		return ReadableMessageOnMemberReadAtCursorKey
	default:
		return ReadableMessageOnMemberDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ReadableMessageOnMemberOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ReadableMessageOnMemberOrderMethodDefault はデフォルト。
	ReadableMessageOnMemberOrderMethodDefault ReadableMessageOnMemberOrderMethod = "default"
	// ReadableMessageOnMemberOrderMethodReadAt は既読日時順。
	ReadableMessageOnMemberOrderMethodReadAt ReadableMessageOnMemberOrderMethod = "read_at"
	// ReadableMessageOnMemberOrderMethodReverseReadAt は既読日時逆順。
	ReadableMessageOnMemberOrderMethodReverseReadAt ReadableMessageOnMemberOrderMethod = "r_read_at"
)

// ReadableMessageOnChatRoomAndMemberOrderMethod チャットルームとメンバー上の既読情報の並び替え方法。
type ReadableMessageOnChatRoomAndMemberOrderMethod string

// ParseReadableMessageOnChatRoomAndMemberOrderMethod はチャットルームとメンバー上の既読情報の並び替え方法をパースする。
func ParseReadableMessageOnChatRoomAndMemberOrderMethod(v string) (any, error) {
	if v == "" {
		return ReadableMessageOnChatRoomAndMemberOrderMethodDefault, nil
	}
	switch v {
	case string(ReadableMessageOnChatRoomAndMemberOrderMethodDefault):
		return ReadableMessageOnChatRoomAndMemberOrderMethodDefault, nil
	case string(ReadableMessageOnChatRoomAndMemberOrderMethodReadAt):
		return ReadableMessageOnChatRoomAndMemberOrderMethodReadAt, nil
	case string(ReadableMessageOnChatRoomAndMemberOrderMethodReverseReadAt):
		return ReadableMessageOnChatRoomAndMemberOrderMethodReverseReadAt, nil
	default:
		return ReadableMessageOnChatRoomAndMemberOrderMethodDefault, nil
	}
}

const (
	// ReadableMessageOnChatRoomAndMemberDefaultCursorKey はデフォルトカーソルキー。
	ReadableMessageOnChatRoomAndMemberDefaultCursorKey = "default"
	// ReadableMessageOnChatRoomAndMemberReadAtCursorKey は既読日時カーソルキー。
	ReadableMessageOnChatRoomAndMemberReadAtCursorKey = "read_at"
)

// GetCursorKeyName はカーソルキー名を取得する。
func (m ReadableMessageOnChatRoomAndMemberOrderMethod) GetCursorKeyName() string {
	switch m {
	case ReadableMessageOnChatRoomAndMemberOrderMethodDefault:
		return ReadableMessageOnChatRoomAndMemberDefaultCursorKey
	case ReadableMessageOnChatRoomAndMemberOrderMethodReadAt, ReadableMessageOnChatRoomAndMemberOrderMethodReverseReadAt:
		return ReadableMessageOnChatRoomAndMemberReadAtCursorKey
	default:
		return ReadableMessageOnChatRoomAndMemberDefaultCursorKey
	}
}

// GetStringValue は文字列を取得する。
func (m ReadableMessageOnChatRoomAndMemberOrderMethod) GetStringValue() string {
	return string(m)
}

const (
	// ReadableMessageOnChatRoomAndMemberOrderMethodDefault はデフォルト。
	ReadableMessageOnChatRoomAndMemberOrderMethodDefault ReadableMessageOnChatRoomAndMemberOrderMethod = "default"
	// ReadableMessageOnChatRoomAndMemberOrderMethodReadAt は既読日時順。
	ReadableMessageOnChatRoomAndMemberOrderMethodReadAt ReadableMessageOnChatRoomAndMemberOrderMethod = "read_at"
	// ReadableMessageOnChatRoomAndMemberOrderMethodReverseReadAt は既読日時逆順。
	ReadableMessageOnChatRoomAndMemberOrderMethodReverseReadAt ReadableMessageOnChatRoomAndMemberOrderMethod = "r_read_at"
)

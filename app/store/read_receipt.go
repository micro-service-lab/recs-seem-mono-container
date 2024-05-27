package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ReadReceipt 既読を表すインターフェース。
type ReadReceipt interface {
	// CountReadableMembersOnMessage メッセージ上のメンバー数を取得する。
	CountReadableMembersOnMessage(
		ctx context.Context, messageID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam) (int64, error)
	// CountReadableMembersOnMessageWithSd SD付きでメッセージ上のメンバー数を取得する。
	CountReadableMembersOnMessageWithSd(
		ctx context.Context, sd Sd, messageID uuid.UUID, where parameter.WhereReadableMemberOnMessageParam) (int64, error)
	// CountReadableMessagesOnChatRoomAndMember チャットルーム、メンバー上のメッセージ数を取得する。
	CountReadableMessagesOnChatRoomAndMember(
		ctx context.Context, chatRoomID, memberID uuid.UUID, where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) (int64, error)
	// CountReadableMessagesOnChatRoomAndMemberWithSd SD付きでチャットルーム、メンバー上のメッセージ数を取得する。
	CountReadableMessagesOnChatRoomAndMemberWithSd(
		ctx context.Context, sd Sd, chatRoomID, memberID uuid.UUID,
		where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) (int64, error)
	// CountReadableMessagesOnChatRoomsAndMember 複数のチャットルーム、メンバー上のメッセージ数を取得する。
	CountReadableMessagesOnChatRoomsAndMember(
		ctx context.Context, chatRoomIDs []uuid.UUID, memberID uuid.UUID,
		where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) ([]entity.ReadReceiptGroupByChatRoom, error)
	// CountReadableMessagesOnChatRoomsAndMemberWithSd SD付きで複数のチャットルーム、メンバー上のメッセージ数を取得する。
	CountReadableMessagesOnChatRoomsAndMemberWithSd(
		ctx context.Context, sd Sd, chatRoomIDs []uuid.UUID, memberID uuid.UUID,
		where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) ([]entity.ReadReceiptGroupByChatRoom, error)
	// CountReadableMessagesOnMember メンバー上のメッセージ数を取得する。
	CountReadableMessagesOnMember(
		ctx context.Context, memberID uuid.UUID, where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) (int64, error)
	// CountReadableMessagesOnMemberWithSd SD付きでメンバー上のメッセージ数を取得する。
	CountReadableMessagesOnMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID, where parameter.WhereReadableMessageOnChatRoomAndMemberParam,
	) (int64, error)
	// CreateReadReceipt 既読情報を作成する。
	CreateReadReceipt(ctx context.Context, param parameter.CreateReadReceiptParam) (entity.ReadReceipt, error)
	// CreateReadReceiptWithSd SD付きで既読情報を作成する。
	CreateReadReceiptWithSd(ctx context.Context, sd Sd, param parameter.CreateReadReceiptParam) (entity.ReadReceipt, error)
	// CreateReadReceipts 複数の既読情報を作成する。
	CreateReadReceipts(ctx context.Context, params []parameter.CreateReadReceiptParam) (int64, error)
	// CreateReadReceiptsWithSd SD付きで複数の既読情報を作成する。
	CreateReadReceiptsWithSd(ctx context.Context, sd Sd, params []parameter.CreateReadReceiptParam) (int64, error)
	// ReadReceipt 既読にする。
	ReadReceipt(ctx context.Context, param parameter.ReadReceiptParam) (entity.ReadReceipt, error)
	// ReadReceiptWithSd SD付きで既読にする。
	ReadReceiptWithSd(ctx context.Context, sd Sd, param parameter.ReadReceiptParam) (entity.ReadReceipt, error)
	// ReadReceipts 複数既読にする。
	ReadReceipts(ctx context.Context, param parameter.ReadReceiptsParam) (int64, error)
	// ReadReceiptsWithSd SD付きで複数既読にする。
	ReadReceiptsWithSd(ctx context.Context, sd Sd, param parameter.ReadReceiptsParam) (int64, error)
	// GetReadableMessagesOnMember メンバー上のメッセージを取得する。
	GetReadableMessagesOnMember(
		ctx context.Context,
		memberID uuid.UUID,
		where parameter.WhereReadableMessageOnMemberParam,
		order parameter.ReadableMessageOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ReadableMessageOnMember], error)
	// GetReadableMessagesOnMemberWithSd SD付きでメンバー上のメッセージを取得する。
	GetReadableMessagesOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberID uuid.UUID,
		where parameter.WhereReadableMessageOnMemberParam,
		order parameter.ReadableMessageOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ReadableMessageOnMember], error)
	// GetPluralReadableMessagesOnMember メンバー上の複数のメッセージを取得する。
	GetPluralReadableMessagesOnMember(
		ctx context.Context,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ReadableMessageOnMemberOrderMethod,
	) (ListResult[entity.ReadableMessageOnMember], error)
	// GetPluralReadableMessagesOnMemberWithSd SD付きでメンバー上の複数のメッセージを取得する。
	GetPluralReadableMessagesOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ReadableMessageOnMemberOrderMethod,
	) (ListResult[entity.ReadableMessageOnMember], error)
	// GetReadableMembersOnMessage メッセージ上のメンバーを取得する。
	GetReadableMembersOnMessage(
		ctx context.Context,
		memberID uuid.UUID,
		where parameter.WhereReadableMemberOnMessageParam,
		order parameter.ReadableMemberOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ReadableMemberOnMessage], error)
	// GetReadableMembersOnMessageWithSd SD付きでメッセージ上のメンバーを取得する。
	GetReadableMembersOnMessageWithSd(
		ctx context.Context,
		sd Sd,
		memberID uuid.UUID,
		where parameter.WhereReadableMemberOnMessageParam,
		order parameter.ReadableMemberOnMessageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ReadableMemberOnMessage], error)
	// GetPluralReadableMembersOnMessage メッセージ上の複数のメンバーを取得する。
	GetPluralReadableMembersOnMessage(
		ctx context.Context,
		messageIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ReadableMemberOnMessageOrderMethod,
	) (ListResult[entity.ReadableMemberOnMessage], error)
	// GetPluralReadableMembersOnMessageWithSd SD付きでメッセージ上の複数のメンバーを取得する。
	GetPluralReadableMembersOnMessageWithSd(
		ctx context.Context,
		sd Sd,
		messageIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ReadableMemberOnMessageOrderMethod,
	) (ListResult[entity.ReadableMemberOnMessage], error)
}
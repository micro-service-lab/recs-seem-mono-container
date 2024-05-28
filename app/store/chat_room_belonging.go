package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomBelonging チャットルーム所属を表すインターフェース。
type ChatRoomBelonging interface {
	// CountChatRoomsOnMember メンバー上のチャットルーム数を取得する。
	CountChatRoomsOnMember(
		ctx context.Context, memberID uuid.UUID, where parameter.WhereChatRoomOnMemberParam) (int64, error)
	// CountChatRoomsOnMemberWithSd SD付きでメンバー上のチャットルーム数を取得する。
	CountChatRoomsOnMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID, where parameter.WhereChatRoomOnMemberParam) (int64, error)
	// CountMembersOnChatRoom チャットルーム上のメンバー数を取得する。
	CountMembersOnChatRoom(
		ctx context.Context, chatRoomID uuid.UUID, where parameter.WhereMemberOnChatRoomParam) (int64, error)
	// CountMembersOnChatRoomWithSd SD付きでチャットルーム上のメンバー数を取得する。
	CountMembersOnChatRoomWithSd(
		ctx context.Context, sd Sd, chatRoomID uuid.UUID, where parameter.WhereMemberOnChatRoomParam) (int64, error)
	// BelongChatRoom メンバーをチャットルームに所属させる。
	BelongChatRoom(ctx context.Context, param parameter.BelongChatRoomParam) (entity.ChatRoomBelonging, error)
	// BelongChatRoomWithSd SD付きでメンバーをチャットルームに所属させる。
	BelongChatRoomWithSd(ctx context.Context, sd Sd, param parameter.BelongChatRoomParam) (entity.ChatRoomBelonging, error)
	// BelongChatRooms メンバーを複数のチャットルームに所属させる。
	BelongChatRooms(ctx context.Context, params []parameter.BelongChatRoomParam) (int64, error)
	// BelongChatRoomsWithSd SD付きでメンバーを複数のチャットルームに所属させる。
	BelongChatRoomsWithSd(ctx context.Context, sd Sd, params []parameter.BelongChatRoomParam) (int64, error)
	// DisbelongChatRoom メンバーをチャットルームから所属解除する。
	DisbelongChatRoom(ctx context.Context, memberID, chatRoomID uuid.UUID) (int64, error)
	// DisbelongChatRoomWithSd SD付きでメンバーをチャットルームから所属解除する。
	DisbelongChatRoomWithSd(ctx context.Context, sd Sd, memberID, chatRoomID uuid.UUID) (int64, error)
	// DisbelongChatRoomOnMember メンバー上のチャットルームから所属解除する。
	DisbelongChatRoomOnMember(ctx context.Context, memberID uuid.UUID) (int64, error)
	// DisbelongChatRoomOnMemberWithSd SD付きでメンバー上のチャットルームから所属解除する。
	DisbelongChatRoomOnMemberWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (int64, error)
	// DisbelongChatRoomOnMembers メンバー上の複数のチャットルームから所属解除する。
	DisbelongChatRoomOnMembers(ctx context.Context, memberIDs []uuid.UUID) (int64, error)
	// DisbelongChatRoomOnMembersWithSd SD付きでメンバー上の複数のチャットルームから所属解除する。
	DisbelongChatRoomOnMembersWithSd(ctx context.Context, sd Sd, memberIDs []uuid.UUID) (int64, error)
	// DisbelongChatRoomOnChatRoom チャットルーム上のメンバーから所属解除する。
	DisbelongChatRoomOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error)
	// DisbelongChatRoomOnChatRoomWithSd SD付きでチャットルーム上のメンバーから所属解除する。
	DisbelongChatRoomOnChatRoomWithSd(ctx context.Context, sd Sd, chatRoomID uuid.UUID) (int64, error)
	// DisbelongChatRoomOnChatRooms チャットルーム上の複数のメンバーから所属解除する。
	DisbelongChatRoomOnChatRooms(ctx context.Context, chatRoomIDs []uuid.UUID) (int64, error)
	// DisbelongChatRoomOnChatRoomsWithSd SD付きでチャットルーム上の複数のメンバーから所属解除する。
	DisbelongChatRoomOnChatRoomsWithSd(ctx context.Context, sd Sd, chatRoomIDs []uuid.UUID) (int64, error)
	// GetChatRoomsOnMember メンバー上のチャットルームを取得する。
	GetChatRoomsOnMember(
		ctx context.Context,
		memberID uuid.UUID,
		where parameter.WhereChatRoomOnMemberParam,
		order parameter.ChatRoomOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomOnMember], error)
	// GetChatRoomsOnMemberWithSd SD付きでメンバー上のチャットルームを取得する。
	GetChatRoomsOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberID uuid.UUID,
		where parameter.WhereChatRoomOnMemberParam,
		order parameter.ChatRoomOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoomOnMember], error)
	// GetMembersOnChatRoom チャットルーム上のメンバーを取得する。
	GetMembersOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		where parameter.WhereMemberOnChatRoomParam,
		order parameter.MemberOnChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoom], error)
	// GetMembersOnChatRoomWithSd SD付きでチャットルーム上のメンバーを取得する。
	GetMembersOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomID uuid.UUID,
		where parameter.WhereMemberOnChatRoomParam,
		order parameter.MemberOnChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoom], error)
	// GetPluralChatRoomsOnMember メンバー上の複数のチャットルームを取得する。
	GetPluralChatRoomsOnMember(
		ctx context.Context,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ChatRoomOnMemberOrderMethod,
	) (ListResult[entity.ChatRoomOnMember], error)
	// GetPluralChatRoomsOnMemberWithSd SD付きでメンバー上の複数のチャットルームを取得する。
	GetPluralChatRoomsOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.ChatRoomOnMemberOrderMethod,
	) (ListResult[entity.ChatRoomOnMember], error)
	// GetPluralMembersOnChatRoom チャットルーム上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoom(
		ctx context.Context,
		chatRoomIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.MemberOnChatRoomOrderMethod,
	) (ListResult[entity.MemberOnChatRoom], error)
	// GetPluralMembersOnChatRoomWithSd SD付きでチャットルーム上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.MemberOnChatRoomOrderMethod,
	) (ListResult[entity.MemberOnChatRoom], error)
}

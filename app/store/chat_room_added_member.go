package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomAddedMember チャットルーム追加メンバーを表すインターフェース。
type ChatRoomAddedMember interface {
	// CountMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバー数を取得する。
	CountMembersOnChatRoomAddMemberAction(
		ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomAddMemberActionParam) (int64, error)
	// CountMembersOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上のメンバー数を取得する。
	CountMembersOnChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomAddMemberActionParam,
	) (int64, error)
	// AddMemberToChatRoomAddMemberAction チャットルーム追加メンバーアクションにメンバーを追加する。
	AddMemberToChatRoomAddMemberAction(
		ctx context.Context, param parameter.CreateChatRoomAddedMemberParam) (entity.ChatRoomAddedMember, error)
	// AddMemberToChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクションにメンバーを追加する。
	AddMemberToChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomAddedMemberParam) (entity.ChatRoomAddedMember, error)
	// AddMembersToChatRoomAddMemberAction チャットルーム追加メンバーアクションにメンバーを複数追加する。
	AddMembersToChatRoomAddMemberAction(
		ctx context.Context, params []parameter.CreateChatRoomAddedMemberParam) (int64, error)
	// AddMembersToChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクションにメンバーを複数追加する。
	AddMembersToChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomAddedMemberParam) (int64, error)
	// DeleteChatRoomAddedMember チャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMember(
		ctx context.Context, chatRoomAddMemberActionID, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMemberWithSd SD付きでチャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMemberWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionID, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバーを削除する。
	DeleteChatRoomAddedMembersOnChatRoomAddMemberAction(
		ctx context.Context, chatRoomAddMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMemberOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上のメンバーを削除する。
	DeleteChatRoomAddedMemberOnChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnChatRoomAddMemberActions チャットルーム追加メンバーアクション上の複数のメンバーを削除する。
	DeleteChatRoomAddedMembersOnChatRoomAddMemberActions(
		ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnChatRoomAddMemberActionsWithSd SD付きでチャットルーム追加メンバーアクション上の複数のメンバーを削除する。
	DeleteChatRoomAddedMembersOnChatRoomAddMemberActionsWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnMember メンバー上のチャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMembersOnMember(
		ctx context.Context, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnMemberWithSd SD付きでメンバー上のチャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMembersOnMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnMembers メンバー上の複数のチャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMembersOnMembers(
		ctx context.Context, memberIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomAddedMembersOnMembersWithSd SD付きでメンバー上の複数のチャットルーム追加メンバーを削除する。
	DeleteChatRoomAddedMembersOnMembersWithSd(
		ctx context.Context, sd Sd, memberIDs []uuid.UUID) (int64, error)
	// GetMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバーを取得する。
	GetMembersOnChatRoomAddMemberAction(
		ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomAddMemberActionParam,
		order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoomAddMemberAction], error)
	// GetMembersOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上のメンバーを取得する。
	GetMembersOnChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomAddMemberActionParam,
		order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoomAddMemberAction], error)
	// GetPluralMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoomAddMemberAction(
		ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID,
		order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberOnChatRoomAddMemberAction], error)
	// GetPluralMembersOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoomAddMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomAddMemberActionIDs []uuid.UUID,
		order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberOnChatRoomAddMemberAction], error)
}

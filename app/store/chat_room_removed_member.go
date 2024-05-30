package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoomRemovedMember チャットルーム追放メンバーを表すインターフェース。
type ChatRoomRemovedMember interface {
	// CountMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバー数を取得する。
	CountMembersOnChatRoomRemoveMemberAction(
		ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomRemoveMemberActionParam) (int64, error)
	// CountMembersOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上のメンバー数を取得する。
	CountMembersOnChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
	) (int64, error)
	// RemoveMemberToChatRoomRemoveMemberAction チャットルーム追放メンバーアクションにメンバーを追放する。
	RemoveMemberToChatRoomRemoveMemberAction(
		ctx context.Context, param parameter.CreateChatRoomRemovedMemberParam) (entity.ChatRoomRemovedMember, error)
	// RemoveMemberToChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクションにメンバーを追放する。
	RemoveMemberToChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomRemovedMemberParam) (entity.ChatRoomRemovedMember, error)
	// RemoveMembersToChatRoomRemoveMemberAction チャットルーム追放メンバーアクションにメンバーを複数追放する。
	RemoveMembersToChatRoomRemoveMemberAction(
		ctx context.Context, params []parameter.CreateChatRoomRemovedMemberParam) (int64, error)
	// RemoveMembersToChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクションにメンバーを複数追放する。
	RemoveMembersToChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, params []parameter.CreateChatRoomRemovedMemberParam) (int64, error)
	// DeleteChatRoomRemovedMember チャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMember(
		ctx context.Context, chatRoomRemoveMemberActionID, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMemberWithSd SD付きでチャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMemberWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionID, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバーを削除する。
	DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(
		ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMemberOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上のメンバーを削除する。
	DeleteChatRoomRemovedMemberOnChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions チャットルーム追放メンバーアクション上の複数のメンバーを削除する。
	DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(
		ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActionsWithSd SD付きでチャットルーム追放メンバーアクション上の複数のメンバーを削除する。
	DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActionsWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnMember メンバー上のチャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMembersOnMember(
		ctx context.Context, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnMemberWithSd SD付きでメンバー上のチャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMembersOnMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnMembers メンバー上の複数のチャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMembersOnMembers(
		ctx context.Context, memberIDs []uuid.UUID) (int64, error)
	// DeleteChatRoomRemovedMembersOnMembersWithSd SD付きでメンバー上の複数のチャットルーム追放メンバーを削除する。
	DeleteChatRoomRemovedMembersOnMembersWithSd(
		ctx context.Context, sd Sd, memberIDs []uuid.UUID) (int64, error)
	// GetMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバーを取得する。
	GetMembersOnChatRoomRemoveMemberAction(
		ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
		order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoomRemoveMemberAction], error)
	// GetMembersOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上のメンバーを取得する。
	GetMembersOnChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionID uuid.UUID,
		where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
		order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnChatRoomRemoveMemberAction], error)
	// GetPluralMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoomRemoveMemberAction(
		ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID,
		order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberOnChatRoomRemoveMemberAction], error)
	// GetPluralMembersOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上の複数のメンバーを取得する。
	GetPluralMembersOnChatRoomRemoveMemberActionWithSd(
		ctx context.Context, sd Sd, chatRoomRemoveMemberActionIDs []uuid.UUID,
		order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberOnChatRoomRemoveMemberAction], error)
}

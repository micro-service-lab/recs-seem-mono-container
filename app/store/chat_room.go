package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// ChatRoom チャットルームを表すインターフェース。
type ChatRoom interface {
	// CountChatRooms チャットルーム数を取得する。
	CountChatRooms(ctx context.Context, where parameter.WhereChatRoomParam) (int64, error)
	// CountChatRoomsWithSd SD付きでチャットルーム数を取得する。
	CountChatRoomsWithSd(ctx context.Context, sd Sd, where parameter.WhereChatRoomParam) (int64, error)
	// CreateChatRoom チャットルームを作成する。
	CreateChatRoom(ctx context.Context, param parameter.CreateChatRoomParam) (entity.ChatRoom, error)
	// CreateChatRoomWithSd SD付きでチャットルームを作成する。
	CreateChatRoomWithSd(
		ctx context.Context, sd Sd, param parameter.CreateChatRoomParam) (entity.ChatRoom, error)
	// CreateChatRooms チャットルームを作成する。
	CreateChatRooms(ctx context.Context, params []parameter.CreateChatRoomParam) (int64, error)
	// CreateChatRoomsWithSd SD付きでチャットルームを作成する。
	CreateChatRoomsWithSd(ctx context.Context, sd Sd, params []parameter.CreateChatRoomParam) (int64, error)
	// DeleteChatRoom チャットルームを削除する。
	DeleteChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error)
	// DeleteChatRoomWithSd SD付きでチャットルームを削除する。
	DeleteChatRoomWithSd(ctx context.Context, sd Sd, chatRoomID uuid.UUID) (int64, error)
	// PluralDeleteChatRooms チャットルームを複数削除する。
	PluralDeleteChatRooms(ctx context.Context, chatRoomIDs []uuid.UUID) (int64, error)
	// PluralDeleteChatRoomsWithSd SD付きでチャットルームを複数削除する。
	PluralDeleteChatRoomsWithSd(ctx context.Context, sd Sd, chatRoomIDs []uuid.UUID) (int64, error)
	// FindChatRoomByID チャットルームを取得する。
	FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (entity.ChatRoom, error)
	// FindChatRoomByIDWithSd SD付きでチャットルームを取得する。
	FindChatRoomByIDWithSd(ctx context.Context, sd Sd, chatRoomID uuid.UUID) (entity.ChatRoom, error)
	// FindChatRoomOnPrivate チャットルームを取得する。
	FindChatRoomOnPrivate(
		ctx context.Context,
		ownerID uuid.UUID,
		memberID uuid.UUID,
	) (entity.ChatRoom, error)
	// FindChatRoomOnPrivateWithSd SD付きでチャットルームを取得する。
	FindChatRoomOnPrivateWithSd(
		ctx context.Context,
		sd Sd,
		ownerID uuid.UUID,
		memberID uuid.UUID,
	) (entity.ChatRoom, error)
	// GetChatRooms チャットルームを取得する。
	GetChatRooms(
		ctx context.Context,
		where parameter.WhereChatRoomParam,
		order parameter.ChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoom], error)
	// GetChatRoomsWithSd SD付きでチャットルームを取得する。
	GetChatRoomsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereChatRoomParam,
		order parameter.ChatRoomOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ChatRoom], error)
	// GetPluralChatRooms チャットルームを取得する。
	GetPluralChatRooms(
		ctx context.Context,
		chatRoomIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoom], error)
	// GetPluralChatRoomsWithSd SD付きでチャットルームを取得する。
	GetPluralChatRoomsWithSd(
		ctx context.Context,
		sd Sd,
		chatRoomIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.ChatRoom], error)
	// UpdateChatRoom チャットルームを更新する。
	UpdateChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		param parameter.UpdateChatRoomParams,
	) (entity.ChatRoom, error)
	// UpdateChatRoomWithSd SD付きでチャットルームを更新する。
	UpdateChatRoomWithSd(
		ctx context.Context, sd Sd, chatRoomID uuid.UUID,
		param parameter.UpdateChatRoomParams) (entity.ChatRoom, error)
}

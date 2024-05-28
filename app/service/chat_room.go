package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageChatRoom チャットルーム管理サービス。
type ManageChatRoom struct {
	DB      store.Store
	Storage storage.Storage
}

// FindChatRoomByID チャットルームをIDで取得する。
func (m *ManageChatRoom) FindChatRoomByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.ChatRoom, error) {
	e, err := m.DB.FindChatRoomByID(ctx, id)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room by id: %w", err)
	}
	return e, nil
}

// FindChatRoomByIDWithCoverImage チャットルームをIDで取得する。
func (m *ManageChatRoom) FindChatRoomByIDWithCoverImage(
	ctx context.Context,
	id uuid.UUID,
) (entity.ChatRoomWithCoverImage, error) {
	e, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, id)
	if err != nil {
		return entity.ChatRoomWithCoverImage{}, fmt.Errorf("failed to find chat room by id with cover image: %w", err)
	}
	return e, nil
}

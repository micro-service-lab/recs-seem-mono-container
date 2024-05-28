package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageAttachableItem 添付可能アイテム管理サービス。
type ManageAttachableItem struct {
	DB store.Store
}

// FindAttachableItemByID 添付可能アイテムをIDで取得する。
func (m *ManageAttachableItem) FindAttachableItemByID(
	ctx context.Context, id uuid.UUID,
) (entity.AttachableItemWithContent, error) {
	e, err := m.DB.FindAttachableItemByID(ctx, id)
	if err != nil {
		return entity.AttachableItemWithContent{}, fmt.Errorf("failed to find attend status by id: %w", err)
	}
	return e, nil
}

// FindAttachableItemByURL 添付可能アイテムをURLで取得する。
func (m *ManageAttachableItem) FindAttachableItemByURL(
	ctx context.Context, url string,
) (entity.AttachableItemWithContent, error) {
	e, err := m.DB.FindAttachableItemByURL(ctx, url)
	if err != nil {
		return entity.AttachableItemWithContent{}, fmt.Errorf("failed to find attend status by url: %w", err)
	}
	return e, nil
}

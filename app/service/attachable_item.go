package service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ManageAttachableItem 添付可能アイテム管理サービス。
type ManageAttachableItem struct {
	DB      store.Store
	Storage storage.Storage
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

// DownloadAttachableItem 添付可能アイテムをダウンロードする。
func (m *ManageAttachableItem) DownloadAttachableItem(
	ctx context.Context, _, id uuid.UUID,
) (io.ReadCloser, string, error) {
	e, err := m.DB.FindAttachableItemByID(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find attend status by id: %w", err)
	}
	if e.FromOuter {
		return nil, "", errhandle.NewCommonError(response.CannotDownloadOuterAttachableItem, nil)
	}
	key, err := m.Storage.GetKeyFromURL(ctx, e.URL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download attachable item: %w", err)
	}
	reader, err := m.Storage.GetObject(ctx, key)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download attachable item: %w", err)
	}

	return reader, e.Alias, nil
}

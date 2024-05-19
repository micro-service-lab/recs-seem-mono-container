package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// AttachableItem 添付可能アイテムを表すインターフェース。
type AttachableItem interface {
	// CountAttachableItems 添付可能アイテム数を取得する。
	CountAttachableItems(ctx context.Context, where parameter.WhereAttachableItemParam) (int64, error)
	// CountAttachableItemsWithSd SD付きで添付可能アイテム数を取得する。
	CountAttachableItemsWithSd(ctx context.Context, sd Sd, where parameter.WhereAttachableItemParam) (int64, error)
	// CreateAttachableItem 添付可能アイテムを作成する。
	CreateAttachableItem(ctx context.Context, param parameter.CreateAttachableItemParam) (entity.AttachableItem, error)
	// CreateAttachableItemWithSd SD付きで添付可能アイテムを作成する。
	CreateAttachableItemWithSd(
		ctx context.Context, sd Sd, param parameter.CreateAttachableItemParam) (entity.AttachableItem, error)
	// CreateAttachableItems 添付可能アイテムを作成する。
	CreateAttachableItems(ctx context.Context, params []parameter.CreateAttachableItemParam) (int64, error)
	// CreateAttachableItemsWithSd SD付きで添付可能アイテムを作成する。
	CreateAttachableItemsWithSd(ctx context.Context, sd Sd, params []parameter.CreateAttachableItemParam) (int64, error)
	// DeleteAttachableItem 添付可能アイテムを削除する。
	DeleteAttachableItem(ctx context.Context, attachableItemID uuid.UUID) (int64, error)
	// DeleteAttachableItemWithSd SD付きで添付可能アイテムを削除する。
	DeleteAttachableItemWithSd(ctx context.Context, sd Sd, attachableItemID uuid.UUID) (int64, error)
	// PluralDeleteAttachableItems 添付可能アイテムを複数削除する。
	PluralDeleteAttachableItems(ctx context.Context, attachableItemIDs []uuid.UUID) (int64, error)
	// PluralDeleteAttachableItemsWithSd SD付きで添付可能アイテムを複数削除する。
	PluralDeleteAttachableItemsWithSd(ctx context.Context, sd Sd, attachableItemIDs []uuid.UUID) (int64, error)
	// FindAttachableItemByID 添付可能アイテムを取得する。
	FindAttachableItemByID(ctx context.Context, attachableItemID uuid.UUID) (entity.AttachableItemWithContent, error)
	// FindAttachableItemByIDWithSd SD付きで添付可能アイテムを取得する。
	FindAttachableItemByIDWithSd(
		ctx context.Context, sd Sd, attachableItemID uuid.UUID) (entity.AttachableItemWithContent, error)
	// FindAttachableItemByIDWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
	FindAttachableItemByIDWithMimeType(
		ctx context.Context, attachableItemID uuid.UUID) (entity.AttachableItemWithMimeType, error)
	// FindAttachableItemByIDWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
	FindAttachableItemByIDWithMimeTypeWithSd(
		ctx context.Context, sd Sd, attachableItemID uuid.UUID) (entity.AttachableItemWithMimeType, error)
	// GetAttachableItems 添付可能アイテムを取得する。
	GetAttachableItems(
		ctx context.Context,
		where parameter.WhereAttachableItemParam,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachableItemWithContent], error)
	// GetAttachableItemsWithSd SD付きで添付可能アイテムを取得する。
	GetAttachableItemsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereAttachableItemParam,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachableItemWithContent], error)
	// GetAttachableItemsWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
	GetAttachableItemsWithMimeType(
		ctx context.Context,
		where parameter.WhereAttachableItemParam,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachableItemWithMimeType], error)
	// GetAttachableItemsWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
	GetAttachableItemsWithMimeTypeWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereAttachableItemParam,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttachableItemWithMimeType], error)
	// GetPluralAttachableItems 添付可能アイテムを取得する。
	GetPluralAttachableItems(
		ctx context.Context,
		attachableItemIDs []uuid.UUID,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachableItemWithContent], error)
	// GetPluralAttachableItemsWithSd SD付きで添付可能アイテムを取得する。
	GetPluralAttachableItemsWithSd(
		ctx context.Context,
		sd Sd,
		attachableItemIDs []uuid.UUID,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachableItemWithContent], error)
	// GetPluralAttachableItemsWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
	GetPluralAttachableItemsWithMimeType(
		ctx context.Context,
		attachableItemIDs []uuid.UUID,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachableItemWithMimeType], error)
	// GetPluralAttachableItemsWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
	GetPluralAttachableItemsWithMimeTypeWithSd(
		ctx context.Context,
		sd Sd,
		attachableItemIDs []uuid.UUID,
		order parameter.AttachableItemOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttachableItemWithMimeType], error)
	// UpdateAttachableItem 添付可能アイテムを更新する。
	UpdateAttachableItem(
		ctx context.Context,
		attachableItemID uuid.UUID,
		param parameter.UpdateAttachableItemParams,
	) (entity.AttachableItem, error)
	// UpdateAttachableItemWithSd SD付きで添付可能アイテムを更新する。
	UpdateAttachableItemWithSd(
		ctx context.Context, sd Sd, attachableItemID uuid.UUID,
		param parameter.UpdateAttachableItemParams) (entity.AttachableItem, error)
}

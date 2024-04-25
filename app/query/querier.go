// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package query

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CountAttachableItems(ctx context.Context, arg CountAttachableItemsParams) (int64, error)
	CreateAttachableItem(ctx context.Context, arg CreateAttachableItemParams) (AttachableItem, error)
	CreateAttachableItems(ctx context.Context, arg []CreateAttachableItemsParams) (int64, error)
	DeleteAttachableItem(ctx context.Context, attachableItemID uuid.UUID) error
	FindAttachableItemByID(ctx context.Context, attachableItemID uuid.UUID) (FindAttachableItemByIDRow, error)
	FindAttachableItemByIDWithMimeType(ctx context.Context, attachableItemID uuid.UUID) (FindAttachableItemByIDWithMimeTypeRow, error)
	GetAttachableItems(ctx context.Context, arg GetAttachableItemsParams) ([]GetAttachableItemsRow, error)
	GetAttachableItemsByMimeTypeIDWithMimeType(ctx context.Context, arg GetAttachableItemsByMimeTypeIDWithMimeTypeParams) ([]GetAttachableItemsByMimeTypeIDWithMimeTypeRow, error)
	GetAttachableItemsByMimeTypeIDWithMimeTypeUseKeysetPaginate(ctx context.Context, arg GetAttachableItemsByMimeTypeIDWithMimeTypeUseKeysetPaginateParams) ([]GetAttachableItemsByMimeTypeIDWithMimeTypeUseKeysetPaginateRow, error)
	GetAttachableItemsByMimeTypeIDWithMimeTypeUseNumberedPaginate(ctx context.Context, arg GetAttachableItemsByMimeTypeIDWithMimeTypeUseNumberedPaginateParams) ([]GetAttachableItemsByMimeTypeIDWithMimeTypeUseNumberedPaginateRow, error)
	GetAttachableItemsUseKeysetPaginate(ctx context.Context, arg GetAttachableItemsUseKeysetPaginateParams) ([]GetAttachableItemsUseKeysetPaginateRow, error)
	GetAttachableItemsUseNumberedPaginate(ctx context.Context, arg GetAttachableItemsUseNumberedPaginateParams) ([]GetAttachableItemsUseNumberedPaginateRow, error)
}

var _ Querier = (*Queries)(nil)
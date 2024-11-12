package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Image 画像を表すインターフェース。
type Image interface {
	// CountImages 画像数を取得する。
	CountImages(ctx context.Context, where parameter.WhereImageParam) (int64, error)
	// CountImagesWithSd SD付きで画像数を取得する。
	CountImagesWithSd(ctx context.Context, sd Sd, where parameter.WhereImageParam) (int64, error)
	// CreateImage 画像を作成する。
	CreateImage(ctx context.Context, param parameter.CreateImageParam) (entity.Image, error)
	// CreateImageWithSd SD付きで画像を作成する。
	CreateImageWithSd(
		ctx context.Context, sd Sd, param parameter.CreateImageParam) (entity.Image, error)
	// CreateImages 画像を作成する。
	CreateImages(ctx context.Context, params []parameter.CreateImageParam) (int64, error)
	// CreateImagesWithSd SD付きで画像を作成する。
	CreateImagesWithSd(ctx context.Context, sd Sd, params []parameter.CreateImageParam) (int64, error)
	// DeleteImage 画像を削除する。
	DeleteImage(ctx context.Context, imageID uuid.UUID) (int64, error)
	// DeleteImageWithSd SD付きで画像を削除する。
	DeleteImageWithSd(ctx context.Context, sd Sd, imageID uuid.UUID) (int64, error)
	// PluralDeleteImages 画像を複数削除する。
	PluralDeleteImages(ctx context.Context, imageIDs []uuid.UUID) (int64, error)
	// PluralDeleteImagesWithSd SD付きで画像を複数削除する。
	PluralDeleteImagesWithSd(ctx context.Context, sd Sd, imageIDs []uuid.UUID) (int64, error)
	// FindImageByID 画像を取得する。
	FindImageByID(ctx context.Context, imageID uuid.UUID) (entity.Image, error)
	// FindImageByIDWithSd SD付きで画像を取得する。
	FindImageByIDWithSd(ctx context.Context, sd Sd, imageID uuid.UUID) (entity.Image, error)
	// FindImageWithAttachableItem 画像を取得する。
	FindImageWithAttachableItem(ctx context.Context, imageID uuid.UUID) (entity.ImageWithAttachableItem, error)
	// FindImageWithAttachableItemWithSd SD付きで画像を取得する。
	FindImageWithAttachableItemWithSd(
		ctx context.Context, sd Sd, imageID uuid.UUID) (entity.ImageWithAttachableItem, error)
	// GetImages 画像を取得する。
	GetImages(
		ctx context.Context,
		where parameter.WhereImageParam,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Image], error)
	// GetImagesWithSd SD付きで画像を取得する。
	GetImagesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereImageParam,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Image], error)
	// GetPluralImages 画像を取得する。
	GetPluralImages(
		ctx context.Context,
		imageIDs []uuid.UUID,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Image], error)
	// GetPluralImagesWithSd SD付きで画像を取得する。
	GetPluralImagesWithSd(
		ctx context.Context,
		sd Sd,
		imageIDs []uuid.UUID,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Image], error)
	// GetImagesWithAttachableItem 画像を取得する。
	GetImagesWithAttachableItem(
		ctx context.Context,
		where parameter.WhereImageParam,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ImageWithAttachableItem], error)
	// GetImagesWithAttachableItemWithSd SD付きで画像を取得する。
	GetImagesWithAttachableItemWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereImageParam,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ImageWithAttachableItem], error)
	// GetPluralImagesWithAttachableItem 画像を取得する。
	GetPluralImagesWithAttachableItem(
		ctx context.Context,
		imageIDs []uuid.UUID,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ImageWithAttachableItem], error)
	// GetPluralImagesWithAttachableItemWithSd SD付きで画像を取得する。
	GetPluralImagesWithAttachableItemWithSd(
		ctx context.Context,
		sd Sd,
		imageIDs []uuid.UUID,
		order parameter.ImageOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ImageWithAttachableItem], error)
}

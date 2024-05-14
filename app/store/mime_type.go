package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// MimeType マイムタイプを表すインターフェース。
type MimeType interface {
	// CountMimeTypes マイムタイプ数を取得する。
	CountMimeTypes(ctx context.Context, where parameter.WhereMimeTypeParam) (int64, error)
	// CountMimeTypesWithSd SD付きでマイムタイプ数を取得する。
	CountMimeTypesWithSd(ctx context.Context, sd Sd, where parameter.WhereMimeTypeParam) (int64, error)
	// CreateMimeType マイムタイプを作成する。
	CreateMimeType(ctx context.Context, param parameter.CreateMimeTypeParam) (entity.MimeType, error)
	// CreateMimeTypeWithSd SD付きでマイムタイプを作成する。
	CreateMimeTypeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateMimeTypeParam) (entity.MimeType, error)
	// CreateMimeTypes マイムタイプを作成する。
	CreateMimeTypes(ctx context.Context, params []parameter.CreateMimeTypeParam) (int64, error)
	// CreateMimeTypesWithSd SD付きでマイムタイプを作成する。
	CreateMimeTypesWithSd(ctx context.Context, sd Sd, params []parameter.CreateMimeTypeParam) (int64, error)
	// DeleteMimeType マイムタイプを削除する。
	DeleteMimeType(ctx context.Context, mimeTypeID uuid.UUID) (int64, error)
	// DeleteMimeTypeWithSd SD付きでマイムタイプを削除する。
	DeleteMimeTypeWithSd(ctx context.Context, sd Sd, mimeTypeID uuid.UUID) (int64, error)
	// DeleteMimeTypeByKey マイムタイプを削除する。
	DeleteMimeTypeByKey(ctx context.Context, key string) (int64, error)
	// DeleteMimeTypeByKeyWithSd SD付きでマイムタイプを削除する。
	DeleteMimeTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeleteMimeTypes マイムタイプを複数削除する。
	PluralDeleteMimeTypes(ctx context.Context, mimeTypeIDs []uuid.UUID) (int64, error)
	// PluralDeleteMimeTypesWithSd SD付きでマイムタイプを複数削除する。
	PluralDeleteMimeTypesWithSd(ctx context.Context, sd Sd, mimeTypeIDs []uuid.UUID) (int64, error)
	// FindMimeTypeByID マイムタイプを取得する。
	FindMimeTypeByID(ctx context.Context, mimeTypeID uuid.UUID) (entity.MimeType, error)
	// FindMimeTypeByIDWithSd SD付きでマイムタイプを取得する。
	FindMimeTypeByIDWithSd(ctx context.Context, sd Sd, mimeTypeID uuid.UUID) (entity.MimeType, error)
	// FindMimeTypeByKey マイムタイプを取得する。
	FindMimeTypeByKey(ctx context.Context, key string) (entity.MimeType, error)
	// FindMimeTypeByKeyWithSd SD付きでマイムタイプを取得する。
	FindMimeTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.MimeType, error)
	// GetMimeTypes マイムタイプを取得する。
	GetMimeTypes(
		ctx context.Context,
		where parameter.WhereMimeTypeParam,
		order parameter.MimeTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MimeType], error)
	// GetMimeTypesWithSd SD付きでマイムタイプを取得する。
	GetMimeTypesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMimeTypeParam,
		order parameter.MimeTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MimeType], error)
	// GetPluralMimeTypes マイムタイプを取得する。
	GetPluralMimeTypes(
		ctx context.Context,
		mimeTypeIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.MimeType], error)
	// GetPluralMimeTypesWithSd SD付きでマイムタイプを取得する。
	GetPluralMimeTypesWithSd(
		ctx context.Context,
		sd Sd,
		mimeTypeIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.MimeType], error)
	// UpdateMimeType マイムタイプを更新する。
	UpdateMimeType(
		ctx context.Context,
		mimeTypeID uuid.UUID,
		param parameter.UpdateMimeTypeParams,
	) (entity.MimeType, error)
	// UpdateMimeTypeWithSd SD付きでマイムタイプを更新する。
	UpdateMimeTypeWithSd(
		ctx context.Context, sd Sd, mimeTypeID uuid.UUID,
		param parameter.UpdateMimeTypeParams) (entity.MimeType, error)
	// UpdateMimeTypeByKey マイムタイプを更新する。
	UpdateMimeTypeByKey(
		ctx context.Context, key string, param parameter.UpdateMimeTypeByKeyParams) (entity.MimeType, error)
	// UpdateMimeTypeByKeyWithSd SD付きでマイムタイプを更新する。
	UpdateMimeTypeByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdateMimeTypeByKeyParams,
	) (entity.MimeType, error)
}

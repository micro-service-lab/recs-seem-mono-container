package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// File ファイルを表すインターフェース。
type File interface {
	// CountFiles ファイル数を取得する。
	CountFiles(ctx context.Context, where parameter.WhereFileParam) (int64, error)
	// CountFilesWithSd SD付きでファイル数を取得する。
	CountFilesWithSd(ctx context.Context, sd Sd, where parameter.WhereFileParam) (int64, error)
	// CreateFile ファイルを作成する。
	CreateFile(ctx context.Context, param parameter.CreateFileParam) (entity.File, error)
	// CreateFileWithSd SD付きでファイルを作成する。
	CreateFileWithSd(
		ctx context.Context, sd Sd, param parameter.CreateFileParam) (entity.File, error)
	// CreateFiles ファイルを作成する。
	CreateFiles(ctx context.Context, params []parameter.CreateFileParam) (int64, error)
	// CreateFilesWithSd SD付きでファイルを作成する。
	CreateFilesWithSd(ctx context.Context, sd Sd, params []parameter.CreateFileParam) (int64, error)
	// DeleteFile ファイルを削除する。
	DeleteFile(ctx context.Context, fileID uuid.UUID) (int64, error)
	// DeleteFileWithSd SD付きでファイルを削除する。
	DeleteFileWithSd(ctx context.Context, sd Sd, fileID uuid.UUID) (int64, error)
	// PluralDeleteFiles ファイルを複数削除する。
	PluralDeleteFiles(ctx context.Context, fileIDs []uuid.UUID) (int64, error)
	// PluralDeleteFilesWithSd SD付きでファイルを複数削除する。
	PluralDeleteFilesWithSd(ctx context.Context, sd Sd, fileIDs []uuid.UUID) (int64, error)
	// FindFileByID ファイルを取得する。
	FindFileByID(ctx context.Context, fileID uuid.UUID) (entity.File, error)
	// FindFileByIDWithSd SD付きでファイルを取得する。
	FindFileByIDWithSd(ctx context.Context, sd Sd, fileID uuid.UUID) (entity.File, error)
	// FindFileWithAttachableItem ファイルを取得する。
	FindFileWithAttachableItem(ctx context.Context, fileID uuid.UUID) (entity.FileWithAttachableItem, error)
	// FindFileWithAttachableItemWithSd SD付きでファイルを取得する。
	FindFileWithAttachableItemWithSd(
		ctx context.Context, sd Sd, fileID uuid.UUID) (entity.FileWithAttachableItem, error)
	// GetFiles ファイルを取得する。
	GetFiles(
		ctx context.Context,
		where parameter.WhereFileParam,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.File], error)
	// GetFilesWithSd SD付きでファイルを取得する。
	GetFilesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereFileParam,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.File], error)
	// GetPluralFiles ファイルを取得する。
	GetPluralFiles(
		ctx context.Context,
		fileIDs []uuid.UUID,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.File], error)
	// GetPluralFilesWithSd SD付きでファイルを取得する。
	GetPluralFilesWithSd(
		ctx context.Context,
		sd Sd,
		fileIDs []uuid.UUID,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.File], error)
	// GetFilesWithAttachableItem ファイルを取得する。
	GetFilesWithAttachableItem(
		ctx context.Context,
		where parameter.WhereFileParam,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.FileWithAttachableItem], error)
	// GetFilesWithAttachableItemWithSd SD付きでファイルを取得する。
	GetFilesWithAttachableItemWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereFileParam,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.FileWithAttachableItem], error)
	// GetPluralFilesWithAttachableItem ファイルを取得する。
	GetPluralFilesWithAttachableItem(
		ctx context.Context,
		fileIDs []uuid.UUID,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.FileWithAttachableItem], error)
	// GetPluralFilesWithAttachableItemWithSd SD付きでファイルを取得する。
	GetPluralFilesWithAttachableItemWithSd(
		ctx context.Context,
		sd Sd,
		fileIDs []uuid.UUID,
		order parameter.FileOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.FileWithAttachableItem], error)
}

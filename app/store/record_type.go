package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// RecordType 議事録タイプを表すインターフェース。
type RecordType interface {
	// CountRecordTypes 議事録タイプ数を取得する。
	CountRecordTypes(ctx context.Context, where parameter.WhereRecordTypeParam) (int64, error)
	// CountRecordTypesWithSd SD付きで議事録タイプ数を取得する。
	CountRecordTypesWithSd(ctx context.Context, sd Sd, where parameter.WhereRecordTypeParam) (int64, error)
	// CreateRecordType 議事録タイプを作成する。
	CreateRecordType(ctx context.Context, param parameter.CreateRecordTypeParam) (entity.RecordType, error)
	// CreateRecordTypeWithSd SD付きで議事録タイプを作成する。
	CreateRecordTypeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateRecordTypeParam) (entity.RecordType, error)
	// CreateRecordTypes 議事録タイプを作成する。
	CreateRecordTypes(ctx context.Context, params []parameter.CreateRecordTypeParam) (int64, error)
	// CreateRecordTypesWithSd SD付きで議事録タイプを作成する。
	CreateRecordTypesWithSd(ctx context.Context, sd Sd, params []parameter.CreateRecordTypeParam) (int64, error)
	// DeleteRecordType 議事録タイプを削除する。
	DeleteRecordType(ctx context.Context, recordTypeID uuid.UUID) (int64, error)
	// DeleteRecordTypeWithSd SD付きで議事録タイプを削除する。
	DeleteRecordTypeWithSd(ctx context.Context, sd Sd, recordTypeID uuid.UUID) (int64, error)
	// DeleteRecordTypeByKey 議事録タイプを削除する。
	DeleteRecordTypeByKey(ctx context.Context, key string) (int64, error)
	// DeleteRecordTypeByKeyWithSd SD付きで議事録タイプを削除する。
	DeleteRecordTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeleteRecordTypes 議事録タイプを複数削除する。
	PluralDeleteRecordTypes(ctx context.Context, recordTypeIDs []uuid.UUID) (int64, error)
	// PluralDeleteRecordTypesWithSd SD付きで議事録タイプを複数削除する。
	PluralDeleteRecordTypesWithSd(ctx context.Context, sd Sd, recordTypeIDs []uuid.UUID) (int64, error)
	// FindRecordTypeByID 議事録タイプを取得する。
	FindRecordTypeByID(ctx context.Context, recordTypeID uuid.UUID) (entity.RecordType, error)
	// FindRecordTypeByIDWithSd SD付きで議事録タイプを取得する。
	FindRecordTypeByIDWithSd(ctx context.Context, sd Sd, recordTypeID uuid.UUID) (entity.RecordType, error)
	// FindRecordTypeByKey 議事録タイプを取得する。
	FindRecordTypeByKey(ctx context.Context, key string) (entity.RecordType, error)
	// FindRecordTypeByKeyWithSd SD付きで議事録タイプを取得する。
	FindRecordTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.RecordType, error)
	// GetRecordTypes 議事録タイプを取得する。
	GetRecordTypes(
		ctx context.Context,
		where parameter.WhereRecordTypeParam,
		order parameter.RecordTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.RecordType], error)
	// GetRecordTypesWithSd SD付きで議事録タイプを取得する。
	GetRecordTypesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereRecordTypeParam,
		order parameter.RecordTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.RecordType], error)
	// GetPluralRecordTypes 議事録タイプを取得する。
	GetPluralRecordTypes(
		ctx context.Context,
		RecordTypeIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.RecordType], error)
	// GetPluralRecordTypesWithSd SD付きで議事録タイプを取得する。
	GetPluralRecordTypesWithSd(
		ctx context.Context,
		sd Sd,
		RecordTypeIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.RecordType], error)
	// UpdateRecordType 議事録タイプを更新する。
	UpdateRecordType(
		ctx context.Context,
		recordTypeID uuid.UUID,
		param parameter.UpdateRecordTypeParams,
	) (entity.RecordType, error)
	// UpdateRecordTypeWithSd SD付きで議事録タイプを更新する。
	UpdateRecordTypeWithSd(
		ctx context.Context, sd Sd, recordTypeID uuid.UUID,
		param parameter.UpdateRecordTypeParams) (entity.RecordType, error)
	// UpdateRecordTypeByKey 議事録タイプを更新する。
	UpdateRecordTypeByKey(
		ctx context.Context, key string, param parameter.UpdateRecordTypeByKeyParams) (entity.RecordType, error)
	// UpdateRecordTypeByKeyWithSd SD付きで議事録タイプを更新する。
	UpdateRecordTypeByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdateRecordTypeByKeyParams,
	) (entity.RecordType, error)
}

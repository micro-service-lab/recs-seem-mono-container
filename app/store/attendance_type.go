package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// AttendanceType 出欠状況タイプを表すインターフェース。
type AttendanceType interface {
	// CountAttendanceTypes 出欠状況タイプ数を取得する。
	CountAttendanceTypes(ctx context.Context, where parameter.WhereAttendanceTypeParam) (int64, error)
	// CountAttendanceTypesWithSd SD付きで出欠状況タイプ数を取得する。
	CountAttendanceTypesWithSd(ctx context.Context, sd Sd, where parameter.WhereAttendanceTypeParam) (int64, error)
	// CreateAttendanceType 出欠状況タイプを作成する。
	CreateAttendanceType(ctx context.Context, param parameter.CreateAttendanceTypeParam) (entity.AttendanceType, error)
	// CreateAttendanceTypeWithSd SD付きで出欠状況タイプを作成する。
	CreateAttendanceTypeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateAttendanceTypeParam) (entity.AttendanceType, error)
	// CreateAttendanceTypes 出欠状況タイプを作成する。
	CreateAttendanceTypes(ctx context.Context, params []parameter.CreateAttendanceTypeParam) (int64, error)
	// CreateAttendanceTypesWithSd SD付きで出欠状況タイプを作成する。
	CreateAttendanceTypesWithSd(ctx context.Context, sd Sd, params []parameter.CreateAttendanceTypeParam) (int64, error)
	// DeleteAttendanceType 出欠状況タイプを削除する。
	DeleteAttendanceType(ctx context.Context, attendanceTypeID uuid.UUID) (int64, error)
	// DeleteAttendanceTypeWithSd SD付きで出欠状況タイプを削除する。
	DeleteAttendanceTypeWithSd(ctx context.Context, sd Sd, attendanceTypeID uuid.UUID) (int64, error)
	// DeleteAttendanceTypeByKey 出欠状況タイプを削除する。
	DeleteAttendanceTypeByKey(ctx context.Context, key string) (int64, error)
	// DeleteAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを削除する。
	DeleteAttendanceTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (int64, error)
	// PluralDeleteAttendanceTypes 出欠状況タイプを複数削除する。
	PluralDeleteAttendanceTypes(ctx context.Context, attendanceTypeIDs []uuid.UUID) (int64, error)
	// PluralDeleteAttendanceTypesWithSd SD付きで出欠状況タイプを複数削除する。
	PluralDeleteAttendanceTypesWithSd(ctx context.Context, sd Sd, attendanceTypeIDs []uuid.UUID) (int64, error)
	// FindAttendanceTypeByID 出欠状況タイプを取得する。
	FindAttendanceTypeByID(ctx context.Context, attendanceTypeID uuid.UUID) (entity.AttendanceType, error)
	// FindAttendanceTypeByIDWithSd SD付きで出欠状況タイプを取得する。
	FindAttendanceTypeByIDWithSd(ctx context.Context, sd Sd, attendanceTypeID uuid.UUID) (entity.AttendanceType, error)
	// FindAttendanceTypeByKey 出欠状況タイプを取得する。
	FindAttendanceTypeByKey(ctx context.Context, key string) (entity.AttendanceType, error)
	// FindAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを取得する。
	FindAttendanceTypeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.AttendanceType, error)
	// GetAttendanceTypes 出欠状況タイプを取得する。
	GetAttendanceTypes(
		ctx context.Context,
		where parameter.WhereAttendanceTypeParam,
		order parameter.AttendanceTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttendanceType], error)
	// GetAttendanceTypesWithSd SD付きで出欠状況タイプを取得する。
	GetAttendanceTypesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereAttendanceTypeParam,
		order parameter.AttendanceTypeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttendanceType], error)
	// GetPluralAttendanceTypes 出欠状況タイプを取得する。
	GetPluralAttendanceTypes(
		ctx context.Context,
		attendanceTypeIDs []uuid.UUID,
		order parameter.AttendanceTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttendanceType], error)
	// GetPluralAttendanceTypesWithSd SD付きで出欠状況タイプを取得する。
	GetPluralAttendanceTypesWithSd(
		ctx context.Context,
		sd Sd,
		attendanceTypeIDs []uuid.UUID,
		order parameter.AttendanceTypeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.AttendanceType], error)
	// UpdateAttendanceType 出欠状況タイプを更新する。
	UpdateAttendanceType(
		ctx context.Context,
		attendanceTypeID uuid.UUID,
		param parameter.UpdateAttendanceTypeParams,
	) (entity.AttendanceType, error)
	// UpdateAttendanceTypeWithSd SD付きで出欠状況タイプを更新する。
	UpdateAttendanceTypeWithSd(
		ctx context.Context, sd Sd, attendanceTypeID uuid.UUID,
		param parameter.UpdateAttendanceTypeParams) (entity.AttendanceType, error)
	// UpdateAttendanceTypeByKey 出欠状況タイプを更新する。
	UpdateAttendanceTypeByKey(
		ctx context.Context, key string, param parameter.UpdateAttendanceTypeByKeyParams) (entity.AttendanceType, error)
	// UpdateAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを更新する。
	UpdateAttendanceTypeByKeyWithSd(
		ctx context.Context,
		sd Sd,
		key string,
		param parameter.UpdateAttendanceTypeByKeyParams,
	) (entity.AttendanceType, error)
}

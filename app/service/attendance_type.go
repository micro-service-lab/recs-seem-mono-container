package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// AttendanceTypeKey 出欠状況タイプキー。
type AttendanceTypeKey string

const (
	// AttendanceTypeKeyEarlyLeave 早退。
	AttendanceTypeKeyEarlyLeave AttendanceTypeKey = "early_leave"
	// AttendanceTypeKeyLateArrival 遅刻。
	AttendanceTypeKeyLateArrival AttendanceTypeKey = "late_arrival"
	// AttendanceTypeKeyAbsence 欠席。
	AttendanceTypeKeyAbsence AttendanceTypeKey = "absence"
)

// AttendanceType 出欠状況タイプ。
type AttendanceType struct {
	Key   string
	Name  string
	Color string
}

// AttendanceTypes 出欠状況タイプ一覧。
var AttendanceTypes = []AttendanceType{
	{Key: "early_leave", Name: "早退", Color: "#ADFF66"},
	{Key: "late_arrival", Name: "遅刻", Color: "#FFB866"},
	{Key: "absence", Name: "欠席", Color: "#FF4D4D"},
}

// ManageAttendanceType 出欠状況タイプ管理サービス。
type ManageAttendanceType struct {
	DB store.Store
}

// CreateAttendanceType 出欠状況タイプを作成する。
func (m *ManageAttendanceType) CreateAttendanceType(
	ctx context.Context,
	name, key, color string,
) (entity.AttendanceType, error) {
	p := parameter.CreateAttendanceTypeParam{
		Name:  name,
		Key:   key,
		Color: color,
	}
	e, err := m.DB.CreateAttendanceType(ctx, p)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to create attendance type: %w", err)
	}
	return e, nil
}

// CreateAttendanceTypes 出欠状況タイプを複数作成する。
func (m *ManageAttendanceType) CreateAttendanceTypes(
	ctx context.Context, ps []parameter.CreateAttendanceTypeParam,
) (int64, error) {
	es, err := m.DB.CreateAttendanceTypes(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create attendance types: %w", err)
	}
	return es, nil
}

// UpdateAttendanceType 出欠状況タイプを更新する。
func (m *ManageAttendanceType) UpdateAttendanceType(
	ctx context.Context, id uuid.UUID, name, key, color string,
) (entity.AttendanceType, error) {
	p := parameter.UpdateAttendanceTypeParams{
		Name:  name,
		Key:   key,
		Color: color,
	}
	e, err := m.DB.UpdateAttendanceType(ctx, id, p)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	return e, nil
}

// DeleteAttendanceType 出欠状況タイプを削除する。
func (m *ManageAttendanceType) DeleteAttendanceType(ctx context.Context, id uuid.UUID) error {
	err := m.DB.DeleteAttendanceType(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return nil
}

// PluralDeleteAttendanceTypes 出欠状況タイプを複数削除する。
func (m *ManageAttendanceType) PluralDeleteAttendanceTypes(ctx context.Context, ids []uuid.UUID) error {
	err := m.DB.PluralDeleteAttendanceTypes(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to plural delete attendance types: %w", err)
	}
	return nil
}

// FindAttendanceTypeByID 出欠状況タイプをIDで取得する。
func (m *ManageAttendanceType) FindAttendanceTypeByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.AttendanceType, error) {
	e, err := m.DB.FindAttendanceTypeByID(ctx, id)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type by id: %w", err)
	}
	return e, nil
}

// FindAttendanceTypeByKey 出欠状況タイプをキーで取得する。
func (m *ManageAttendanceType) FindAttendanceTypeByKey(ctx context.Context, key string) (entity.AttendanceType, error) {
	e, err := m.DB.FindAttendanceTypeByKey(ctx, key)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type by key: %w", err)
	}
	return e, nil
}

// GetAttendanceTypes 出欠状況タイプを取得する。
func (m *ManageAttendanceType) GetAttendanceTypes(
	ctx context.Context,
	whereSearchName string,
	order parameter.AttendanceTypeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.AttendanceType], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereAttendanceTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetAttendanceTypes(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get attendance types: %w", err)
	}
	return r, nil
}

// GetAttendanceTypesCount 出欠状況タイプの数を取得する。
func (m *ManageAttendanceType) GetAttendanceTypesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereAttendanceTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountAttendanceTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get attendance types count: %w", err)
	}
	return c, nil
}

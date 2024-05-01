package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// AttendStatusKey 出席状況キー。
type AttendStatusKey string

const (
	// AttendStatusKeyPresent 出席。
	AttendStatusKeyPresent AttendStatusKey = "present"
	// AttendStatusKeyAbsent 欠席。
	AttendStatusKeyAbsent AttendStatusKey = "absent"
	// AttendStatusKeyTemporarilyAbsent 一時欠席。
	AttendStatusKeyTemporarilyAbsent AttendStatusKey = "temporarily_absent"
	// AttendStatusKeyGoHome 退室。
	AttendStatusKeyGoHome AttendStatusKey = "go_home"
	// AttendStatusKeyNotAttend 未出席。
	AttendStatusKeyNotAttend AttendStatusKey = "not_attend"
)

// AttendStatus 出席状況。
type AttendStatus struct {
	Key  string
	Name string
}

// AttendStatues 出席状況一覧。
var AttendStatues = []AttendStatus{
	{Key: "present", Name: "出席"},
	{Key: "absent", Name: "欠席"},
	{Key: "temporarily_absent", Name: "一時退席"},
	{Key: "go_home", Name: "退室"},
	{Key: "not_attend", Name: "未出席"},
}

// ManageAttendStatus 出席状況管理サービス。
type ManageAttendStatus struct {
	db store.Store
}

// CreateAttendStatus 出席状況を作成する。
func (m *ManageAttendStatus) CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error) {
	p := parameter.CreateAttendStatusParam{
		Name: name,
		Key:  key,
	}
	e, err := m.db.CreateAttendStatus(ctx, p)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to create attend status: %w", err)
	}
	return e, nil
}

// CreateAttendStatuses 出席状況を複数作成する。
func (m *ManageAttendStatus) CreateAttendStatuses(
	ctx context.Context, ps []parameter.CreateAttendStatusParam,
) (int64, error) {
	es, err := m.db.CreateAttendStatuses(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return es, nil
}

// UpdateAttendStatus 出席状況を更新する。
func (m *ManageAttendStatus) UpdateAttendStatus(
	ctx context.Context, id uuid.UUID, name, key string,
) (entity.AttendStatus, error) {
	p := parameter.UpdateAttendStatusParams{
		Name: name,
		Key:  key,
	}
	e, err := m.db.UpdateAttendStatus(ctx, id, p)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	return e, nil
}

// DeleteAttendStatus 出席状況を削除する。
func (m *ManageAttendStatus) DeleteAttendStatus(ctx context.Context, id uuid.UUID) error {
	err := m.db.DeleteAttendStatus(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

// FindAttendStatusByID 出席状況をIDで取得する。
func (m *ManageAttendStatus) FindAttendStatusByID(ctx context.Context, id uuid.UUID) (entity.AttendStatus, error) {
	e, err := m.db.FindAttendStatusByID(ctx, id)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status by id: %w", err)
	}
	return e, nil
}

// FindAttendStatusByKey 出席状況をキーで取得する。
func (m *ManageAttendStatus) FindAttendStatusByKey(ctx context.Context, key string) (entity.AttendStatus, error) {
	e, err := m.db.FindAttendStatusByKey(ctx, key)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status by key: %w", err)
	}
	return e, nil
}

// GetAttendStatuses 出席状況を取得する。
func (m *ManageAttendStatus) GetAttendStatuses(
	ctx context.Context,
	whereSearchName string,
	order parameter.AttendStatusOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.AttendStatus], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereAttendStatusParam{
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
	r, err := m.db.GetAttendStatuses(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	return r, nil
}

// GetAttendStatusesCount 出席状況の数を取得する。
func (m *ManageAttendStatus) GetAttendStatusesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereAttendStatusParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.db.CountAttendStatuses(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get attend statuses count: %w", err)
	}
	return c, nil
}

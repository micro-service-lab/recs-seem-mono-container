package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageAttendStatus 出席状況管理サービス。
type ManageAttendStatus struct {
	db store.Store
}

// CreateAttendStatus 出席状況を作成する。
func (m *ManageAttendStatus) CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error) {
	p := store.CreateAttendStatusParam{
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
func (m *ManageAttendStatus) CreateAttendStatuses(ctx context.Context, ps []store.CreateAttendStatusParam) (int64, error) {
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
	p := store.UpdateAttendStatusParams{
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
	order store.AttendStatusOrderMethod,
	pg Pagination,
	limit int64,
	cursor string,
	offset int64,
	withCount bool,
) (store.ListResult[entity.AttendStatus], error) {
	wc := store.WithCountParam{
		Valid: withCount,
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := store.WhereAttendStatusParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: offset},
			Limit:  entity.Int{Int64: limit},
		}
	case CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: cursor,
			Limit:  entity.Int{Int64: limit},
		}
	_:
	}
	r, err := m.db.GetAttendStatuses(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	return r, nil
}

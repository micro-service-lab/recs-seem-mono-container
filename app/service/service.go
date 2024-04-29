// Package service provides a application service.
package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// Manager is a manager for services.
type Manager struct {
	ManageAbsence
	ManageAttendStatus
}

// NewManager creates a new Manager.
func NewManager(db store.Store) *Manager {
	return &Manager{
		ManageAttendStatus: ManageAttendStatus{db: db},
		ManageAbsence:      ManageAbsence{db: db},
	}
}

var _ ManagerInterface = (*Manager)(nil)

// ManagerInterface is a interface for manager.
type ManagerInterface interface {
	AttendStatusManager
}

// AttendStatusManager is a interface for attend status service.
type AttendStatusManager interface {
	CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error)
	CreateAttendStatuses(ctx context.Context, ps []parameter.CreateAttendStatusParam) (int64, error)
	UpdateAttendStatus(ctx context.Context, id uuid.UUID, name, key string) (entity.AttendStatus, error)
	DeleteAttendStatus(ctx context.Context, id uuid.UUID) error
	FindAttendStatusByID(ctx context.Context, id uuid.UUID) (entity.AttendStatus, error)
	FindAttendStatusByKey(ctx context.Context, key string) (entity.AttendStatus, error)
	GetAttendStatuses(
		ctx context.Context,
		whereSearchName string,
		order parameter.AttendStatusOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.AttendStatus], error)
	GetAttendStatusesCount(ctx context.Context, whereSearchName string) (int64, error)
}

var _ ManagerInterface = (*Manager)(nil)

// Package service provides a application service.
package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// Pagination ページネーション。
type Pagination string

const (
	// NumberedPagination ページネーション。
	NumberedPagination Pagination = "numbered"
	// CursorPagination カーソルページネーション。
	CursorPagination Pagination = "cursor"
	// NonePagination ページネーションなし。
	NonePagination Pagination = "none"
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
	}
}

// ManagerInterface is a interface for manager.
type ManagerInterface interface {
	AttendStatusManager
}

// AttendStatusManager is a interface for attend status service.
type AttendStatusManager interface {
	CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error)
	CreateAttendStatuses(ctx context.Context, ps []store.CreateAttendStatusParam) (int64, error)
	UpdateAttendStatus(ctx context.Context, id uuid.UUID, name, key string) (entity.AttendStatus, error)
	DeleteAttendStatus(ctx context.Context, id uuid.UUID) error
	FindAttendStatusByID(ctx context.Context, id uuid.UUID) (entity.AttendStatus, error)
	FindAttendStatusByKey(ctx context.Context, key string) (entity.AttendStatus, error)
}

var _ ManagerInterface = (*Manager)(nil)

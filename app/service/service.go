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
	ManageAttendanceType
	ManageEventType
}

// NewManager creates a new Manager.
func NewManager(db store.Store) *Manager {
	return &Manager{
		ManageAttendStatus:   ManageAttendStatus{DB: db},
		ManageAbsence:        ManageAbsence{DB: db},
		ManageAttendanceType: ManageAttendanceType{DB: db},
		ManageEventType:      ManageEventType{DB: db},
	}
}

//go:generate moq -out service_mock.go . ManagerInterface

// ManagerInterface is a interface for manager.
type ManagerInterface interface {
	AttendStatusManager
	AttendanceTypeManager
	EventTypeManager
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

// AttendanceTypeManager is a interface for attendance type service.
type AttendanceTypeManager interface {
	CreateAttendanceType(ctx context.Context, name, key, color string) (entity.AttendanceType, error)
	CreateAttendanceTypes(ctx context.Context, ps []parameter.CreateAttendanceTypeParam) (int64, error)
	UpdateAttendanceType(ctx context.Context, id uuid.UUID, name, key, color string) (entity.AttendanceType, error)
	DeleteAttendanceType(ctx context.Context, id uuid.UUID) error
	FindAttendanceTypeByID(ctx context.Context, id uuid.UUID) (entity.AttendanceType, error)
	FindAttendanceTypeByKey(ctx context.Context, key string) (entity.AttendanceType, error)
	GetAttendanceTypes(
		ctx context.Context,
		whereSearchName string,
		order parameter.AttendanceTypeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.AttendanceType], error)
	GetAttendanceTypesCount(ctx context.Context, whereSearchName string) (int64, error)
}

// EventTypeManager is a interface for event type service.
type EventTypeManager interface {
	CreateEventType(ctx context.Context, name, key, color string) (entity.EventType, error)
	CreateEventTypes(ctx context.Context, ps []parameter.CreateEventTypeParam) (int64, error)
	UpdateEventType(ctx context.Context, id uuid.UUID, name, key, color string) (entity.EventType, error)
	DeleteEventType(ctx context.Context, id uuid.UUID) error
	FindEventTypeByID(ctx context.Context, id uuid.UUID) (entity.EventType, error)
	FindEventTypeByKey(ctx context.Context, key string) (entity.EventType, error)
	GetEventTypes(
		ctx context.Context,
		whereSearchName string,
		order parameter.EventTypeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.EventType], error)
	GetEventTypesCount(ctx context.Context, whereSearchName string) (int64, error)
}

var _ ManagerInterface = (*Manager)(nil)

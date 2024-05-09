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
	ManagePermissionCategory
	ManagePolicyCategory
	ManageMimeType
	ManageRecordType
}

// NewManager creates a new Manager.
func NewManager(db store.Store) *Manager {
	return &Manager{
		ManageAttendStatus:       ManageAttendStatus{DB: db},
		ManageAbsence:            ManageAbsence{DB: db},
		ManageAttendanceType:     ManageAttendanceType{DB: db},
		ManageEventType:          ManageEventType{DB: db},
		ManagePermissionCategory: ManagePermissionCategory{DB: db},
		ManagePolicyCategory:     ManagePolicyCategory{DB: db},
		ManageMimeType:           ManageMimeType{DB: db},
		ManageRecordType:         ManageRecordType{DB: db},
	}
}

//go:generate moq -out service_mock.go . ManagerInterface

// ManagerInterface is a interface for manager.
type ManagerInterface interface {
	AttendStatusManager
	AttendanceTypeManager
	EventTypeManager
	PermissionCategoryManager
	PolicyCategoryManager
	MimeTypeManager
	RecordTypeManager
}

// AttendStatusManager is a interface for attend status service.
type AttendStatusManager interface {
	CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error)
	CreateAttendStatuses(ctx context.Context, ps []parameter.CreateAttendStatusParam) (int64, error)
	UpdateAttendStatus(ctx context.Context, id uuid.UUID, name, key string) (entity.AttendStatus, error)
	DeleteAttendStatus(ctx context.Context, id uuid.UUID) error
	PluralDeleteAttendStatuses(ctx context.Context, ids []uuid.UUID) error
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
	PluralDeleteAttendanceTypes(ctx context.Context, ids []uuid.UUID) error
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
	PluralDeleteEventTypes(ctx context.Context, ids []uuid.UUID) error
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

// PermissionCategoryManager is a interface for event type service.
type PermissionCategoryManager interface {
	CreatePermissionCategory(ctx context.Context, name, key, description string) (entity.PermissionCategory, error)
	CreatePermissionCategories(ctx context.Context, ps []parameter.CreatePermissionCategoryParam) (int64, error)
	UpdatePermissionCategory(
		ctx context.Context, id uuid.UUID, name, key, description string) (entity.PermissionCategory, error)
	DeletePermissionCategory(ctx context.Context, id uuid.UUID) error
	PluralDeletePermissionCategories(ctx context.Context, ids []uuid.UUID) error
	FindPermissionCategoryByID(ctx context.Context, id uuid.UUID) (entity.PermissionCategory, error)
	FindPermissionCategoryByKey(ctx context.Context, key string) (entity.PermissionCategory, error)
	GetPermissionCategories(
		ctx context.Context,
		whereSearchName string,
		order parameter.PermissionCategoryOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.PermissionCategory], error)
	GetPermissionCategoriesCount(ctx context.Context, whereSearchName string) (int64, error)
}

// PolicyCategoryManager is a interface for event type service.
type PolicyCategoryManager interface {
	CreatePolicyCategory(ctx context.Context, name, key, description string) (entity.PolicyCategory, error)
	CreatePolicyCategories(ctx context.Context, ps []parameter.CreatePolicyCategoryParam) (int64, error)
	UpdatePolicyCategory(
		ctx context.Context, id uuid.UUID, name, key, description string) (entity.PolicyCategory, error)
	DeletePolicyCategory(ctx context.Context, id uuid.UUID) error
	PluralDeletePolicyCategories(ctx context.Context, ids []uuid.UUID) error
	FindPolicyCategoryByID(ctx context.Context, id uuid.UUID) (entity.PolicyCategory, error)
	FindPolicyCategoryByKey(ctx context.Context, key string) (entity.PolicyCategory, error)
	GetPolicyCategories(
		ctx context.Context,
		whereSearchName string,
		order parameter.PolicyCategoryOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.PolicyCategory], error)
	GetPolicyCategoriesCount(ctx context.Context, whereSearchName string) (int64, error)
}

// MimeTypeManager is a interface for mime type service.
type MimeTypeManager interface {
	CreateMimeType(ctx context.Context, name, key, kind string) (entity.MimeType, error)
	CreateMimeTypes(ctx context.Context, ps []parameter.CreateMimeTypeParam) (int64, error)
	UpdateMimeType(ctx context.Context, id uuid.UUID, name, key, kind string) (entity.MimeType, error)
	DeleteMimeType(ctx context.Context, id uuid.UUID) error
	PluralDeleteMimeTypes(ctx context.Context, ids []uuid.UUID) error
	FindMimeTypeByID(ctx context.Context, id uuid.UUID) (entity.MimeType, error)
	FindMimeTypeByKey(ctx context.Context, key string) (entity.MimeType, error)
	GetMimeTypes(
		ctx context.Context,
		whereSearchName string,
		order parameter.MimeTypeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MimeType], error)
	GetMimeTypesCount(ctx context.Context, whereSearchName string) (int64, error)
}

// RecordTypeManager is a interface for record type service.
type RecordTypeManager interface {
	CreateRecordType(ctx context.Context, name, key string) (entity.RecordType, error)
	CreateRecordTypes(ctx context.Context, ps []parameter.CreateRecordTypeParam) (int64, error)
	UpdateRecordType(ctx context.Context, id uuid.UUID, name, key string) (entity.RecordType, error)
	DeleteRecordType(ctx context.Context, id uuid.UUID) error
	PluralDeleteRecordTypes(ctx context.Context, ids []uuid.UUID) error
	FindRecordTypeByID(ctx context.Context, id uuid.UUID) (entity.RecordType, error)
	FindRecordTypeByKey(ctx context.Context, key string) (entity.RecordType, error)
	GetRecordTypes(
		ctx context.Context,
		whereSearchName string,
		order parameter.RecordTypeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.RecordType], error)
	GetRecordTypesCount(ctx context.Context, whereSearchName string) (int64, error)
}

var _ ManagerInterface = (*Manager)(nil)

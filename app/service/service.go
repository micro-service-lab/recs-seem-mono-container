// Package service provides a application service.
package service

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/hasher"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
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
	ManagePermission
	ManagePolicy
	ManageRole
	ManageRoleAssociation
	ManageOrganization
	ManageImage
	ManageFile
	ManageAttachableItem
	ManageGrade
	ManageGroup
	ManageChatRoom
	ManageMember
	ManageStudent
	ManageProfessor
}

// NewManager creates a new Manager.
func NewManager(
	db store.Store,
	_ i18n.Translation,
	stg storage.Storage,
	h hasher.Hash,
	clk clock.Clock,
) *Manager {
	return &Manager{
		ManageAttendStatus:       ManageAttendStatus{DB: db},
		ManageAbsence:            ManageAbsence{DB: db},
		ManageAttendanceType:     ManageAttendanceType{DB: db},
		ManageEventType:          ManageEventType{DB: db},
		ManagePermissionCategory: ManagePermissionCategory{DB: db},
		ManagePolicyCategory:     ManagePolicyCategory{DB: db},
		ManageMimeType:           ManageMimeType{DB: db},
		ManageRecordType:         ManageRecordType{DB: db},
		ManagePermission:         ManagePermission{DB: db},
		ManagePolicy:             ManagePolicy{DB: db},
		ManageRole:               ManageRole{DB: db},
		ManageRoleAssociation:    ManageRoleAssociation{DB: db},
		ManageOrganization:       ManageOrganization{DB: db},
		ManageImage:              ManageImage{DB: db, Storage: stg},
		ManageFile:               ManageFile{DB: db, Storage: stg},
		ManageAttachableItem:     ManageAttachableItem{DB: db},
		ManageGrade:              ManageGrade{DB: db},
		ManageGroup:              ManageGroup{DB: db},
		ManageChatRoom:           ManageChatRoom{DB: db, Storage: stg},
		ManageMember:             ManageMember{DB: db, Hash: h, Clocker: clk, Storage: stg},
		ManageStudent:            ManageStudent{DB: db, Hash: h, Clocker: clk, Storage: stg},
		ManageProfessor:          ManageProfessor{DB: db, Hash: h, Clocker: clk, Storage: stg},
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
	PermissionManager
	PolicyManager
	RoleManager
	RoleAssociationManager
	OrganizationManager
	ImageManager
	FileManager
	AttachableItemManager
	GradeManager
	GroupManager
	ChatRoomManager
	MemberManager
	StudentManager
	ProfessorManager
}

// AttendStatusManager is a interface for attend status service.
type AttendStatusManager interface {
	CreateAttendStatus(ctx context.Context, name, key string) (entity.AttendStatus, error)
	CreateAttendStatuses(ctx context.Context, ps []parameter.CreateAttendStatusParam) (int64, error)
	UpdateAttendStatus(ctx context.Context, id uuid.UUID, name, key string) (entity.AttendStatus, error)
	DeleteAttendStatus(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteAttendStatuses(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeleteAttendanceType(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteAttendanceTypes(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeleteEventType(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteEventTypes(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeletePermissionCategory(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeletePermissionCategories(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeletePolicyCategory(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeletePolicyCategories(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeleteMimeType(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteMimeTypes(ctx context.Context, ids []uuid.UUID) (int64, error)
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
	DeleteRecordType(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteRecordTypes(ctx context.Context, ids []uuid.UUID) (int64, error)
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

// PermissionManager is a interface for event type service.
type PermissionManager interface {
	CreatePermission(ctx context.Context, name, key, description string, categoryID uuid.UUID) (entity.Permission, error)
	CreatePermissions(ctx context.Context, ps []parameter.CreatePermissionParam) (int64, error)
	UpdatePermission(
		ctx context.Context, id uuid.UUID, name, key, description string, categoryID uuid.UUID) (entity.Permission, error)
	DeletePermission(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeletePermissions(ctx context.Context, ids []uuid.UUID) (int64, error)
	FindPermissionByID(ctx context.Context, id uuid.UUID) (entity.Permission, error)
	FindPermissionByIDWithCategory(
		ctx context.Context,
		id uuid.UUID,
	) (entity.PermissionWithCategory, error)
	FindPermissionByKey(ctx context.Context, key string) (entity.Permission, error)
	FindPermissionByKeyWithCategory(
		ctx context.Context, key string,
	) (entity.PermissionWithCategory, error)
	GetPermissions(
		ctx context.Context,
		whereSearchName string,
		whereInCategories []uuid.UUID,
		order parameter.PermissionOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Permission], error)
	GetPermissionsWithCategory(
		ctx context.Context,
		whereSearchName string,
		whereInCategories []uuid.UUID,
		order parameter.PermissionOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.PermissionWithCategory], error)
	GetPermissionsCount(ctx context.Context, whereSearchName string, whereInCategories []uuid.UUID) (int64, error)
}

// PolicyManager is a interface for event type service.
type PolicyManager interface {
	CreatePolicy(ctx context.Context, name, key, description string, categoryID uuid.UUID) (entity.Policy, error)
	CreatePolicies(ctx context.Context, ps []parameter.CreatePolicyParam) (int64, error)
	UpdatePolicy(
		ctx context.Context, id uuid.UUID, name, key, description string, categoryID uuid.UUID) (entity.Policy, error)
	DeletePolicy(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeletePolicies(ctx context.Context, ids []uuid.UUID) (int64, error)
	FindPolicyByID(ctx context.Context, id uuid.UUID) (entity.Policy, error)
	FindPolicyByIDWithCategory(
		ctx context.Context,
		id uuid.UUID,
	) (entity.PolicyWithCategory, error)
	FindPolicyByKey(ctx context.Context, key string) (entity.Policy, error)
	FindPolicyByKeyWithCategory(
		ctx context.Context, key string,
	) (entity.PolicyWithCategory, error)
	GetPolicies(
		ctx context.Context,
		whereSearchName string,
		whereInCategories []uuid.UUID,
		order parameter.PolicyOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Policy], error)
	GetPoliciesWithCategory(
		ctx context.Context,
		whereSearchName string,
		whereInCategories []uuid.UUID,
		order parameter.PolicyOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.PolicyWithCategory], error)
	GetPoliciesCount(ctx context.Context, whereSearchName string, whereInCategories []uuid.UUID) (int64, error)
}

// RoleManager is a interface for event type service.
type RoleManager interface {
	CreateRole(ctx context.Context, name, description string) (entity.Role, error)
	CreateRoles(ctx context.Context, ps []parameter.CreateRoleParam) (int64, error)
	UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (entity.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteRoles(ctx context.Context, ids []uuid.UUID) (int64, error)
	FindRoleByID(ctx context.Context, id uuid.UUID) (entity.Role, error)
	GetRoles(
		ctx context.Context,
		whereSearchName string,
		order parameter.RoleOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Role], error)
	GetRolesCount(ctx context.Context, whereSearchName string) (int64, error)
}

// RoleAssociationManager is a interface for role association service.
type RoleAssociationManager interface {
	AssociateRoles(
		ctx context.Context, params []parameter.AssociationRoleParam,
	) (int64, error)
	DisassociateRoleOnPolicy(
		ctx context.Context, policyID uuid.UUID,
	) (int64, error)
	DisassociateRoleOnPolicies(
		ctx context.Context, policyIDs []uuid.UUID,
	) (int64, error)
	PluralDisassociateRoleOnPolicy(
		ctx context.Context, policyID uuid.UUID, roleIDs []uuid.UUID,
	) (int64, error)
	DisassociatePolicyOnRole(
		ctx context.Context, roleID uuid.UUID,
	) (int64, error)
	DisassociatePolicyOnRoles(
		ctx context.Context, roleIDs []uuid.UUID,
	) (int64, error)
	PluralDisassociatePolicyOnRole(
		ctx context.Context, roleID uuid.UUID, policyIDs []uuid.UUID,
	) (int64, error)
	GetRolesOnPolicy(
		ctx context.Context, policyID uuid.UUID,
		whereSearchName string,
		order parameter.RoleOnPolicyOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.RoleOnPolicy], error)
	GetPoliciesOnRole(
		ctx context.Context, roleID uuid.UUID,
		whereSearchName string,
		order parameter.PolicyOnRoleOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.PolicyOnRole], error)
	GetRolesOnPolicyCount(
		ctx context.Context, policyID uuid.UUID,
		whereSearchName string,
	) (int64, error)
	GetPoliciesOnRoleCount(
		ctx context.Context, roleID uuid.UUID,
		whereSearchName string,
	) (int64, error)
}

// OrganizationManager is a interface for organization service.
type OrganizationManager interface {
	CreateWholeOrganization(
		ctx context.Context,
		name string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Organization, err error)
	DeleteWholeOrganization(ctx context.Context) (c int64, err error)
	UpdateWholeOrganization(
		ctx context.Context,
		name string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Organization, err error)
	FindWholeOrganization(ctx context.Context) (entity.Organization, error)
	CreateOrganization(
		ctx context.Context, name string, description, color entity.String,
	) (e entity.Organization, err error)
	UpdateOrganization(
		ctx context.Context, id uuid.UUID, name string, description, color entity.String,
	) (e entity.Organization, err error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) (int64, error)
	FindOrganizationByID(
		ctx context.Context,
		id uuid.UUID,
	) (entity.Organization, error)
	FindOrganizationWithChatRoom(
		ctx context.Context,
		id uuid.UUID,
	) (entity.OrganizationWithChatRoom, error)
	FindOrganizationWithDetail(
		ctx context.Context,
		id uuid.UUID,
	) (entity.OrganizationWithDetail, error)
	FindOrganizationWithChatRoomAndDetail(
		ctx context.Context,
		id uuid.UUID,
	) (entity.OrganizationWithChatRoomAndDetail, error)
	GetOrganizations(
		ctx context.Context,
		whereSearchName string,
		whereOrganizationType parameter.WhereOrganizationType,
		wherePersonalMemberID uuid.UUID,
		order parameter.OrganizationOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Organization], error)
	GetOrganizationsWithDetail(
		ctx context.Context,
		whereSearchName string,
		whereOrganizationType parameter.WhereOrganizationType,
		wherePersonalMemberID uuid.UUID,
		order parameter.OrganizationOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.OrganizationWithDetail], error)
	GetOrganizationsWithChatRoom(
		ctx context.Context,
		whereSearchName string,
		whereOrganizationType parameter.WhereOrganizationType,
		wherePersonalMemberID uuid.UUID,
		order parameter.OrganizationOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.OrganizationWithChatRoom], error)
	GetOrganizationsWithChatRoomAndDetail(
		ctx context.Context,
		whereSearchName string,
		whereOrganizationType parameter.WhereOrganizationType,
		wherePersonalMemberID uuid.UUID,
		order parameter.OrganizationOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error)
	GetOrganizationsCount(
		ctx context.Context,
		whereSearchName string,
		whereOrganizationType parameter.WhereOrganizationType,
		wherePersonalMemberID uuid.UUID,
	) (int64, error)
}

// ImageManager is a interface for image service.
type ImageManager interface {
	CreateImage(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
	) (entity.Image, error)
	CreateImages(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageServiceParam,
	) ([]entity.Image, error)
	CreateImageSpecifyFilename(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
		filename string,
	) (entity.Image, error)
	CreateImagesSpecifyFilename(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageSpecifyFilenameServiceParam,
	) ([]entity.Image, error)
	CreateImageFromOuter(
		ctx context.Context,
		url,
		alias string,
		size entity.Float,
		ownerID entity.UUID,
		mimeTypeID uuid.UUID,
		height, width entity.Float,
	) (entity.Image, error)
	CreateImagesFromOuter(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageFromOuterServiceParam,
	) ([]entity.Image, error)
	DeleteImage(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteImages(
		ctx context.Context, ids []uuid.UUID,
	) (int64, error)
	GetImages(
		ctx context.Context,
		order parameter.ImageOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Image], error)
	GetImagesCount(
		ctx context.Context,
	) (int64, error)
}

// FileManager is a interface for file service.
type FileManager interface {
	CreateFile(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
	) (entity.File, error)
	CreateFiles(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileServiceParam,
	) ([]entity.File, error)
	CreateFileSpecifyFilename(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
		filename string,
	) (entity.File, error)
	CreateFilesSpecifyFilename(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileSpecifyFilenameServiceParam,
	) ([]entity.File, error)
	CreateFileFromOuter(
		ctx context.Context,
		url,
		alias string,
		size entity.Float,
		ownerID entity.UUID,
		mimeTypeID uuid.UUID,
	) (entity.File, error)
	CreateFilesFromOuter(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileFromOuterServiceParam,
	) ([]entity.File, error)
	DeleteFile(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteFiles(
		ctx context.Context, ids []uuid.UUID,
	) (int64, error)
	GetFiles(
		ctx context.Context,
		order parameter.FileOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.File], error)
	GetFilesCount(
		ctx context.Context,
	) (int64, error)
}

// GradeManager is a interface for grade service.
type GradeManager interface {
	CreateGrade(
		ctx context.Context,
		name, key string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Grade, err error)
	CreateGrades(
		ctx context.Context, ps []parameter.CreateGradeServiceParam,
	) (c int64, err error)
	DeleteGrade(ctx context.Context, id uuid.UUID) (c int64, err error)
	PluralDeleteGrades(
		ctx context.Context, ids []uuid.UUID,
	) (c int64, err error)
	UpdateGrade(
		ctx context.Context,
		id uuid.UUID,
		name string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Grade, err error)
	GetGrades(
		ctx context.Context,
		order parameter.GradeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Grade], error)
	GetGradesWithOrganization(
		ctx context.Context,
		order parameter.GradeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.GradeWithOrganization], error)
	GetGradesCount(
		ctx context.Context,
	) (int64, error)
}

// GroupManager is a interface for group service.
type GroupManager interface {
	CreateGroup(
		ctx context.Context,
		name, key string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Group, err error)
	CreateGroups(
		ctx context.Context, ps []parameter.CreateGroupServiceParam,
	) (c int64, err error)
	DeleteGroup(ctx context.Context, id uuid.UUID) (c int64, err error)
	PluralDeleteGroups(
		ctx context.Context, ids []uuid.UUID,
	) (c int64, err error)
	UpdateGroup(
		ctx context.Context,
		id uuid.UUID,
		name string,
		description, color entity.String,
		coverImageID entity.UUID,
	) (e entity.Group, err error)
	GetGroups(
		ctx context.Context,
		order parameter.GroupOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Group], error)
	GetGroupsWithOrganization(
		ctx context.Context,
		order parameter.GroupOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.GroupWithOrganization], error)
	GetGroupsCount(
		ctx context.Context,
	) (int64, error)
}

// ChatRoomManager is a interface for chat room service.
type ChatRoomManager interface {
	FindChatRoomByID(
		ctx context.Context,
		id uuid.UUID,
	) (entity.ChatRoom, error)
	FindChatRoomByIDWithCoverImage(
		ctx context.Context,
		id uuid.UUID,
	) (entity.ChatRoomWithCoverImage, error)
}

// MemberManager is a interface for member service.
type MemberManager interface {
	UpdateMember(
		ctx context.Context,
		id uuid.UUID,
		email,
		name string,
		firstName,
		lastName entity.String,
		profileImageID entity.UUID,
	) (e entity.Member, err error)
	DeleteMember(ctx context.Context, id uuid.UUID) (int64, error)
	UpdateMemberPassword(
		ctx context.Context,
		id uuid.UUID,
		rawPassword string,
	) (entity.Member, error)
	UpdateMemberRole(
		ctx context.Context,
		id uuid.UUID,
		roleID entity.UUID,
	) (e entity.Member, err error)
	UpdateMemberLoginID(
		ctx context.Context,
		id uuid.UUID,
		loginID string,
	) (e entity.Member, err error)
}

// StudentManager is a interface for student service.
type StudentManager interface {
	CreateStudent(
		ctx context.Context,
		loginID,
		rawPassword,
		email,
		name string,
		firstName,
		lastName entity.String,
		gradeID,
		groupID uuid.UUID,
		roleID entity.UUID,
	) (e entity.Student, err error)
	DeleteStudent(ctx context.Context, id uuid.UUID) (c int64, err error)
	UpdateStudentGrade(
		ctx context.Context,
		id uuid.UUID,
		gradeID uuid.UUID,
	) (e entity.Student, err error)
	UpdateStudentGroup(
		ctx context.Context,
		id uuid.UUID,
		groupID uuid.UUID,
	) (e entity.Student, err error)
}

// ProfessorManager is a interface for professor service.
type ProfessorManager interface {
	CreateProfessor(
		ctx context.Context,
		loginID,
		rawPassword,
		email,
		name string,
		firstName,
		lastName entity.String,
		roleID entity.UUID,
	) (e entity.Professor, err error)
	DeleteProfessor(ctx context.Context, id uuid.UUID) (c int64, err error)
}

// AttachableItemManager is a interface for attachable item service.
type AttachableItemManager interface {
	FindAttachableItemByID(ctx context.Context, id uuid.UUID) (entity.AttachableItemWithContent, error)
	FindAttachableItemByURL(ctx context.Context, url string) (entity.AttachableItemWithContent, error)
}

var _ ManagerInterface = (*Manager)(nil)

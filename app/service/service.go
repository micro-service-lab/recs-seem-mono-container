// Package service provides a application service.
package service

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/hasher"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
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
	ManageChatRoomActionType
	ManageAuth
	ManageChatRoomAction
	ManageChatRoomBelonging
	ManageMembership
	ManageMessage
	ManageAttachedMessage
	ManageReadReceipt
}

// NewManager creates a new Manager.
func NewManager(
	db store.Store,
	_ i18n.Translation,
	stg storage.Storage,
	h hasher.Hash,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
	cfg config.Config,
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
		ManageOrganization:       ManageOrganization{DB: db, Clocker: clk, Storage: stg},
		ManageImage:              ManageImage{DB: db, Storage: stg},
		ManageFile:               ManageFile{DB: db, Storage: stg},
		ManageAttachableItem:     ManageAttachableItem{DB: db},
		ManageGrade:              ManageGrade{DB: db, Clocker: clk, Storage: stg},
		ManageGroup:              ManageGroup{DB: db, Clocker: clk, Storage: stg},
		ManageChatRoom:           ManageChatRoom{DB: db, Storage: stg, Clocker: clk},
		ManageMember:             ManageMember{DB: db, Hash: h, Clocker: clk, Storage: stg},
		ManageStudent:            ManageStudent{DB: db, Hash: h, Clocker: clk, Storage: stg},
		ManageProfessor:          ManageProfessor{DB: db, Hash: h, Clocker: clk, Storage: stg},
		ManageChatRoomActionType: ManageChatRoomActionType{DB: db},
		ManageAuth: ManageAuth{
			DB: db, Hash: h, Auth: auth, SessionManager: ssm, Clocker: clk, Config: cfg,
		},
		ManageChatRoomAction:    ManageChatRoomAction{DB: db},
		ManageChatRoomBelonging: ManageChatRoomBelonging{DB: db, Clocker: clk},
		ManageMembership:        ManageMembership{DB: db, Clocker: clk},
		ManageMessage:           ManageMessage{DB: db, Clocker: clk, Storage: stg},
		ManageAttachedMessage:   ManageAttachedMessage{DB: db, Clocker: clk, Storage: stg},
		ManageReadReceipt:       ManageReadReceipt{DB: db, Clocker: clk},
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
	ChatRoomActionTypeManager
	AuthManager
	ChatRoomActionManager
	ChatRoomBelongingManager
	MembershipManager
	MessageManager
	AttachedMessageManager
	ReadReceiptManager
}

// AuthManager is a interface for auth service.
type AuthManager interface {
	Login(ctx context.Context, loginID, password string) (entity.AuthJwt, error)
	RefreshToken(ctx context.Context, refreshToken string) (entity.AuthJwt, error)
	Logout(ctx context.Context, memberID uuid.UUID) error
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

// ChatRoomActionTypeManager is a interface for chat room action type service.
type ChatRoomActionTypeManager interface {
	CreateChatRoomActionType(ctx context.Context, name, key string) (entity.ChatRoomActionType, error)
	CreateChatRoomActionTypes(ctx context.Context, ps []parameter.CreateChatRoomActionTypeParam) (int64, error)
	UpdateChatRoomActionType(ctx context.Context, id uuid.UUID, name, key string) (entity.ChatRoomActionType, error)
	DeleteChatRoomActionType(ctx context.Context, id uuid.UUID) (int64, error)
	PluralDeleteChatRoomActionTypes(ctx context.Context, ids []uuid.UUID) (int64, error)
	FindChatRoomActionTypeByID(ctx context.Context, id uuid.UUID) (entity.ChatRoomActionType, error)
	FindChatRoomActionTypeByKey(ctx context.Context, key string) (entity.ChatRoomActionType, error)
	GetChatRoomActionTypes(
		ctx context.Context,
		whereSearchName string,
		order parameter.ChatRoomActionTypeOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.ChatRoomActionType], error)
	GetChatRoomActionTypesCount(ctx context.Context, whereSearchName string) (int64, error)
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
		ctx context.Context,
		name string,
		description, color entity.String,
		ownerID uuid.UUID,
		members []uuid.UUID,
		withChatRoom bool,
		chatRoomCoverImageID entity.UUID,
	) (e entity.Organization, err error)
	UpdateOrganization(
		ctx context.Context,
		id uuid.UUID,
		name string,
		description, color entity.String,
		ownerID uuid.UUID,
	) (e entity.Organization, err error)
	DeleteOrganization(
		ctx context.Context,
		id uuid.UUID,
		ownerID uuid.UUID,
	) (c int64, err error)
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
	) (entity.ImageWithAttachableItem, error)
	CreateImages(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageServiceParam,
	) ([]entity.ImageWithAttachableItem, error)
	CreateImageSpecifyFilename(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
		filename string,
	) (entity.ImageWithAttachableItem, error)
	CreateImagesSpecifyFilename(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageSpecifyFilenameServiceParam,
	) ([]entity.ImageWithAttachableItem, error)
	CreateImageFromOuter(
		ctx context.Context,
		url,
		alias string,
		size entity.Float,
		ownerID entity.UUID,
		mimeTypeID uuid.UUID,
		height, width entity.Float,
	) (entity.ImageWithAttachableItem, error)
	CreateImagesFromOuter(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateImageFromOuterServiceParam,
	) ([]entity.ImageWithAttachableItem, error)
	DeleteImage(ctx context.Context, id uuid.UUID, ownerID entity.UUID) (int64, error)
	PluralDeleteImages(
		ctx context.Context, ids []uuid.UUID, ownerID entity.UUID,
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
	) (entity.FileWithAttachableItem, error)
	CreateFiles(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileServiceParam,
	) ([]entity.FileWithAttachableItem, error)
	CreateFileSpecifyFilename(
		ctx context.Context,
		origin io.Reader,
		alias string,
		ownerID entity.UUID,
		filename string,
	) (entity.FileWithAttachableItem, error)
	CreateFilesSpecifyFilename(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileSpecifyFilenameServiceParam,
	) ([]entity.FileWithAttachableItem, error)
	CreateFileFromOuter(
		ctx context.Context,
		url,
		alias string,
		size entity.Float,
		ownerID entity.UUID,
		mimeTypeID uuid.UUID,
	) (entity.FileWithAttachableItem, error)
	CreateFilesFromOuter(
		ctx context.Context,
		ownerID entity.UUID,
		params []parameter.CreateFileFromOuterServiceParam,
	) ([]entity.FileWithAttachableItem, error)
	DeleteFile(ctx context.Context, id uuid.UUID, ownerID entity.UUID) (int64, error)
	PluralDeleteFiles(
		ctx context.Context, ids []uuid.UUID, ownerID entity.UUID,
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
	FindPrivateChatRoom(
		ctx context.Context,
		ownerID,
		memberID uuid.UUID,
	) (entity.ChatRoom, error)
	CreateChatRoom(
		ctx context.Context,
		name string,
		coverImageID entity.UUID,
		ownerID uuid.UUID,
		members []uuid.UUID,
	) (e entity.ChatRoom, err error)
	CreatePrivateChatRoom(
		ctx context.Context,
		ownerID uuid.UUID,
		memberID uuid.UUID,
	) (e entity.ChatRoom, err error)
	UpdateChatRoom(
		ctx context.Context,
		id uuid.UUID,
		name string,
		coverImageID entity.UUID,
		ownerID uuid.UUID,
	) (e entity.ChatRoom, err error)
	DeleteChatRoom(
		ctx context.Context,
		id uuid.UUID,
		ownerID uuid.UUID,
	) (c int64, err error)
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
	FindMemberByID(
		ctx context.Context,
		id uuid.UUID,
	) (entity.Member, error)
	FindAuthMemberByID(
		ctx context.Context,
		id uuid.UUID,
	) (entity.AuthMember, error)
	GetMembers(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.Member], error)
	GetMembersWithAttendStatus(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithAttendStatus], error)
	GetMembersWithDetail(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithDetail], error)
	GetMembersWithCrew(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithCrew], error)
	GetMembersWithProfileImage(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithProfileImage], error)
	GetMembersWithPersonalOrganization(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithPersonalOrganization], error)
	GetMembersWithRole(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithRole], error)
	GetMembersWithCrewAndProfileImageAndAttendStatus(
		ctx context.Context,
		whereSearchName string,
		whereHasInPolicies []uuid.UUID,
		whereInAttendStatuses []uuid.UUID,
		whereInGrades []uuid.UUID,
		whereInGroups []uuid.UUID,
		order parameter.MemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.MemberWithCrewAndProfileImageAndAttendStatus], error)
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
	) (e entity.StudentWithMember, err error)
	UpdateStudentGroup(
		ctx context.Context,
		id uuid.UUID,
		groupID uuid.UUID,
	) (e entity.StudentWithMember, err error)
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

// ChatRoomActionManager is a interface for chat room action service.
type ChatRoomActionManager interface {
	GetChatRoomActionsOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		whereInTypes []uuid.UUID,
		order parameter.ChatRoomActionOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (store.ListResult[entity.ChatRoomActionPractical], error)
}

// ChatRoomBelongingManager is a interface for chat room belonging service.
type ChatRoomBelongingManager interface {
	BelongMembersOnChatRoom(
		ctx context.Context,
		chatRoomID,
		ownerID uuid.UUID,
		memberIDs []uuid.UUID,
	) (e int64, err error)
	RemoveMembersFromChatRoom(
		ctx context.Context,
		chatRoomID,
		ownerID uuid.UUID,
		memberIDs []uuid.UUID,
	) (e int64, err error)
	WithdrawMemberFromChatRoom(
		ctx context.Context,
		chatRoomID,
		memberID uuid.UUID,
	) (e int64, err error)
	GetChatRoomsOnMember(
		ctx context.Context,
		memberID uuid.UUID,
		whereSearchName string,
		order parameter.ChatRoomOnMemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (es store.ListResult[entity.PracticalChatRoomOnMember], err error)
	GetChatRoomsOnMemberCount(
		ctx context.Context, memberID uuid.UUID,
		whereSearchName string,
	) (es int64, err error)
	GetMembersOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		whereSearchName string,
		order parameter.MemberOnChatRoomOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (es store.ListResult[entity.MemberOnChatRoom], err error)
	GetMembersOnChatRoomCount(
		ctx context.Context, chatRoomID uuid.UUID,
		whereSearchName string,
	) (es int64, err error)
}

// MembershipManager is a interface for membership service.
type MembershipManager interface {
	BelongMembersOnOrganization(
		ctx context.Context,
		organizationID,
		ownerID uuid.UUID,
		memberIDs []uuid.UUID,
	) (e int64, err error)
	RemoveMembersFromOrganization(
		ctx context.Context,
		organizationID,
		ownerID uuid.UUID,
		memberIDs []uuid.UUID,
	) (e int64, err error)
	WithdrawMemberFromOrganization(
		ctx context.Context,
		organizationID,
		memberID uuid.UUID,
	) (e int64, err error)
	GetOrganizationsOnMember(
		ctx context.Context,
		memberID uuid.UUID,
		whereSearchName string,
		order parameter.OrganizationOnMemberOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (es store.ListResult[entity.OrganizationOnMember], err error)
	GetOrganizationsOnMemberCount(
		ctx context.Context, memberID uuid.UUID,
		whereSearchName string,
	) (es int64, err error)
	GetMembersOnOrganization(
		ctx context.Context,
		chatRoomID uuid.UUID,
		whereSearchName string,
		order parameter.MemberOnOrganizationOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (es store.ListResult[entity.MemberOnOrganization], err error)
	GetMembersOnOrganizationCount(
		ctx context.Context, chatRoomID uuid.UUID,
		whereSearchName string,
	) (es int64, err error)
}

// AttachableItemManager is a interface for attachable item service.
type AttachableItemManager interface {
	FindAttachableItemByID(ctx context.Context, id uuid.UUID) (entity.AttachableItemWithContent, error)
	FindAttachableItemByURL(ctx context.Context, url string) (entity.AttachableItemWithContent, error)
}

// MessageManager is a interface for message service.
type MessageManager interface {
	CreateMessage(
		ctx context.Context,
		senderID, chatRoomID uuid.UUID,
		content string,
		attachments []uuid.UUID,
	) (e entity.Message, err error)
	CreateMessageOnPrivateRoom(
		ctx context.Context,
		senderID, receiverID uuid.UUID,
		content string,
		attachments []uuid.UUID,
	) (e entity.Message, err error)
	DeleteMessage(
		ctx context.Context,
		chatRoomID,
		ownerID, messageID uuid.UUID,
	) (e int64, err error)
	ForceDeleteMessages(
		ctx context.Context,
		messageIDs []uuid.UUID,
	) (e int64, err error)
	DeleteMessagesBefore(
		ctx context.Context,
		chatRoomIDs []uuid.UUID,
		earlierPostedAt time.Time,
	) (e int64, err error)
	DeleteMessagesBeforeAll(
		ctx context.Context,
		earlierPostedAt time.Time,
	) (e int64, err error)
	EditMessage(
		ctx context.Context,
		chatRoomID,
		ownerID, messageID uuid.UUID,
		content string,
	) (e entity.Message, err error)
	GetMessagesOnChatRoom(
		ctx context.Context,
		chatRoomID uuid.UUID,
		whereInSenders []uuid.UUID,
		whereSearchBody string,
		whereEarlierPostedAt time.Time,
		whereLaterPostedAt time.Time,
		whereEarlierLastEditedAt time.Time,
		whereLaterLastEditedAt time.Time,
		order parameter.MessageOrderMethod,
		pg parameter.Pagination,
		limit parameter.Limit,
		cursor parameter.Cursor,
		offset parameter.Offset,
		withCount parameter.WithCount,
	) (e store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments], err error)
}

// AttachedMessageManager is a interface for manage attached message service.
type AttachedMessageManager interface {
	AttachItemsOnMessage(
		ctx context.Context,
		chatRoomID,
		messageID, ownerID uuid.UUID,
		attachments []uuid.UUID,
	) (e int64, err error)
	DetachItemsOnMessage(
		ctx context.Context,
		chatRoomID,
		messageID, ownerID uuid.UUID,
		attachments []uuid.UUID,
	) (e int64, err error)
}

// ReadReceiptManager is a interface for read receipt service.
type ReadReceiptManager interface {
	CountUnreadReceiptsOnMember(
		ctx context.Context,
		memberID uuid.UUID,
	) (int64, error)
	ReadMessage(
		ctx context.Context,
		chatRoomID,
		memberID, messageID uuid.UUID,
	) (read bool, err error)
	ReadMessagesOnMember(
		ctx context.Context,
		memberID uuid.UUID,
	) (e int64, err error)
	ReadMessagesOnChatRoomAndMember(
		ctx context.Context,
		chatRoomID, memberID uuid.UUID,
	) (e int64, err error)
}

var _ ManagerInterface = (*Manager)(nil)

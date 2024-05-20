package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Organization オーガナイゼーションを表すインターフェース。
type Organization interface {
	// CountOrganizations オーガナイゼーション数を取得する。
	CountOrganizations(ctx context.Context, where parameter.WhereOrganizationParam) (int64, error)
	// CountOrganizationsWithSd SD付きでオーガナイゼーション数を取得する。
	CountOrganizationsWithSd(ctx context.Context, sd Sd, where parameter.WhereOrganizationParam) (int64, error)
	// CreateOrganization オーガナイゼーションを作成する。
	CreateOrganization(ctx context.Context, param parameter.CreateOrganizationParam) (entity.Organization, error)
	// CreateOrganizationWithSd SD付きでオーガナイゼーションを作成する。
	CreateOrganizationWithSd(
		ctx context.Context, sd Sd, param parameter.CreateOrganizationParam) (entity.Organization, error)
	// CreateOrganizations オーガナイゼーションを作成する。
	CreateOrganizations(ctx context.Context, params []parameter.CreateOrganizationParam) (int64, error)
	// CreateOrganizationsWithSd SD付きでオーガナイゼーションを作成する。
	CreateOrganizationsWithSd(ctx context.Context, sd Sd, params []parameter.CreateOrganizationParam) (int64, error)
	// DeleteOrganization オーガナイゼーションを削除する。
	DeleteOrganization(ctx context.Context, organizationID uuid.UUID) (int64, error)
	// DeleteOrganizationWithSd SD付きでオーガナイゼーションを削除する。
	DeleteOrganizationWithSd(ctx context.Context, sd Sd, organizationID uuid.UUID) (int64, error)
	// PluralDeleteOrganizations オーガナイゼーションを複数削除する。
	PluralDeleteOrganizations(ctx context.Context, organizationIDs []uuid.UUID) (int64, error)
	// PluralDeleteOrganizationsWithSd SD付きでオーガナイゼーションを複数削除する。
	PluralDeleteOrganizationsWithSd(ctx context.Context, sd Sd, organizationIDs []uuid.UUID) (int64, error)
	// FindOrganizationByID オーガナイゼーションを取得する。
	FindOrganizationByID(ctx context.Context, organizationID uuid.UUID) (entity.Organization, error)
	// FindOrganizationByIDWithSd SD付きでオーガナイゼーションを取得する。
	FindOrganizationByIDWithSd(ctx context.Context, sd Sd, organizationID uuid.UUID) (entity.Organization, error)
	// FindOrganizationWithChatRoom オーガナイゼーションを取得する。
	FindOrganizationWithChatRoom(ctx context.Context, organizationID uuid.UUID) (entity.OrganizationWithChatRoom, error)
	// FindOrganizationWithChatRoomWithSd SD付きでオーガナイゼーションを取得する。
	FindOrganizationWithChatRoomWithSd(
		ctx context.Context, sd Sd, organizationID uuid.UUID) (entity.OrganizationWithChatRoom, error)
	// FindOrganizationWithDetail オーガナイゼーションを取得する。
	FindOrganizationWithDetail(ctx context.Context, organizationID uuid.UUID) (entity.OrganizationWithDetail, error)
	// FindOrganizationWithDetailWithSd SD付きでオーガナイゼーションを取得する。
	FindOrganizationWithDetailWithSd(
		ctx context.Context, sd Sd, organizationID uuid.UUID) (entity.OrganizationWithDetail, error)
	// FindOrganizationWithChatRoomAndDetail オーガナイゼーションを取得する。
	FindOrganizationWithChatRoomAndDetail(
		ctx context.Context, organizationID uuid.UUID) (entity.OrganizationWithChatRoomAndDetail, error)
	// FindOrganizationWithChatRoomAndDetailWithSd SD付きでオーガナイゼーションを取得する。
	FindOrganizationWithChatRoomAndDetailWithSd(
		ctx context.Context, sd Sd, organizationID uuid.UUID) (entity.OrganizationWithChatRoomAndDetail, error)
	// FindWholeOrganization 全体オーガナイゼーションを取得する。
	FindWholeOrganization(ctx context.Context) (entity.Organization, error)
	// FindWholeOrganizationWithSd SD付きで全体オーガナイゼーションを取得する。
	FindWholeOrganizationWithSd(ctx context.Context, sd Sd) (entity.Organization, error)
	// FindPersonalOrganization 個人オーガナイゼーションを取得する。
	FindPersonalOrganization(ctx context.Context, memberID uuid.UUID) (entity.Organization, error)
	// FindPersonalOrganizationWithSd SD付きで個人オーガナイゼーションを取得する。
	FindPersonalOrganizationWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.Organization, error)
	// GetOrganizations オーガナイゼーションを取得する。
	GetOrganizations(
		ctx context.Context,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Organization], error)
	// GetOrganizationsWithSd SD付きでオーガナイゼーションを取得する。
	GetOrganizationsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Organization], error)
	// GetPluralOrganizations オーガナイゼーションを取得する。
	GetPluralOrganizations(
		ctx context.Context,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Organization], error)
	// GetPluralOrganizationsWithSd SD付きでオーガナイゼーションを取得する。
	GetPluralOrganizationsWithSd(
		ctx context.Context,
		sd Sd,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Organization], error)
	// GetOrganizationsWithChatRoom オーガナイゼーションを取得する。
	GetOrganizationsWithChatRoom(
		ctx context.Context,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithChatRoom], error)
	// GetOrganizationsWithChatRoomWithSd SD付きでオーガナイゼーションを取得する。
	GetOrganizationsWithChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithChatRoom], error)
	// GetPluralOrganizationsWithChatRoom オーガナイゼーションを取得する。
	GetPluralOrganizationsWithChatRoom(
		ctx context.Context,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithChatRoom], error)
	// GetPluralOrganizationsWithChatRoomWithSd SD付きでオーガナイゼーションを取得する。
	GetPluralOrganizationsWithChatRoomWithSd(
		ctx context.Context,
		sd Sd,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithChatRoom], error)
	// GetOrganizationsWithDetail オーガナイゼーションを取得する。
	GetOrganizationsWithDetail(
		ctx context.Context,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithDetail], error)
	// GetOrganizationsWithDetailWithSd SD付きでオーガナイゼーションを取得する。
	GetOrganizationsWithDetailWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithDetail], error)
	// GetPluralOrganizationsWithDetail オーガナイゼーションを取得する。
	GetPluralOrganizationsWithDetail(
		ctx context.Context,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithDetail], error)
	// GetPluralOrganizationsWithDetailWithSd SD付きでオーガナイゼーションを取得する。
	GetPluralOrganizationsWithDetailWithSd(
		ctx context.Context,
		sd Sd,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithDetail], error)
	// GetOrganizationsWithChatRoomAndDetail オーガナイゼーションを取得する。
	GetOrganizationsWithChatRoomAndDetail(
		ctx context.Context,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithChatRoomAndDetail], error)
	// GetOrganizationsWithChatRoomAndDetailWithSd SD付きでオーガナイゼーションを取得する。
	GetOrganizationsWithChatRoomAndDetailWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereOrganizationParam,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationWithChatRoomAndDetail], error)
	// GetPluralOrganizationsWithChatRoomAndDetail オーガナイゼーションを取得する。
	GetPluralOrganizationsWithChatRoomAndDetail(
		ctx context.Context,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithChatRoomAndDetail], error)
	// GetPluralOrganizationsWithChatRoomAndDetailWithSd SD付きでオーガナイゼーションを取得する。
	GetPluralOrganizationsWithChatRoomAndDetailWithSd(
		ctx context.Context,
		sd Sd,
		organizationIDs []uuid.UUID,
		order parameter.OrganizationOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.OrganizationWithChatRoomAndDetail], error)
	// UpdateOrganization オーガナイゼーションを更新する。
	UpdateOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
		param parameter.UpdateOrganizationParams,
	) (entity.Organization, error)
	// UpdateOrganizationWithSd SD付きでオーガナイゼーションを更新する。
	UpdateOrganizationWithSd(
		ctx context.Context, sd Sd, organizationID uuid.UUID,
		param parameter.UpdateOrganizationParams) (entity.Organization, error)
}

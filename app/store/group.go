package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Group グループを表すインターフェース。
type Group interface {
	// CountGroups グループ数を取得する。
	CountGroups(ctx context.Context, where parameter.WhereGroupParam) (int64, error)
	// CountGroupsWithSd SD付きでグループ数を取得する。
	CountGroupsWithSd(ctx context.Context, sd Sd, where parameter.WhereGroupParam) (int64, error)
	// CreateGroup グループを作成する。
	CreateGroup(ctx context.Context, param parameter.CreateGroupParam) (entity.Group, error)
	// CreateGroupWithSd SD付きでグループを作成する。
	CreateGroupWithSd(
		ctx context.Context, sd Sd, param parameter.CreateGroupParam) (entity.Group, error)
	// CreateGroups グループを作成する。
	CreateGroups(ctx context.Context, params []parameter.CreateGroupParam) (int64, error)
	// CreateGroupsWithSd SD付きでグループを作成する。
	CreateGroupsWithSd(ctx context.Context, sd Sd, params []parameter.CreateGroupParam) (int64, error)
	// DeleteGroup グループを削除する。
	DeleteGroup(ctx context.Context, groupID uuid.UUID) (int64, error)
	// DeleteGroupWithSd SD付きでグループを削除する。
	DeleteGroupWithSd(ctx context.Context, sd Sd, groupID uuid.UUID) (int64, error)
	// PluralDeleteGroups グループを複数削除する。
	PluralDeleteGroups(ctx context.Context, groupIDs []uuid.UUID) (int64, error)
	// PluralDeleteGroupsWithSd SD付きでグループを複数削除する。
	PluralDeleteGroupsWithSd(ctx context.Context, sd Sd, groupIDs []uuid.UUID) (int64, error)
	// FindGroupByID グループを取得する。
	FindGroupByID(ctx context.Context, groupID uuid.UUID) (entity.Group, error)
	// FindGroupByIDWithSd SD付きでグループを取得する。
	FindGroupByIDWithSd(ctx context.Context, sd Sd, groupID uuid.UUID) (entity.Group, error)
	// FindGroupWithOrganization グループを取得する。
	FindGroupWithOrganization(ctx context.Context, groupID uuid.UUID) (entity.GroupWithOrganization, error)
	// FindGroupWithOrganizationWithSd SD付きでグループを取得する。
	FindGroupWithOrganizationWithSd(
		ctx context.Context, sd Sd, groupID uuid.UUID) (entity.GroupWithOrganization, error)
	// GetGroups グループを取得する。
	GetGroups(
		ctx context.Context,
		where parameter.WhereGroupParam,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Group], error)
	// GetGroupsWithSd SD付きでグループを取得する。
	GetGroupsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereGroupParam,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Group], error)
	// GetPluralGroups グループを取得する。
	GetPluralGroups(
		ctx context.Context,
		groupIDs []uuid.UUID,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Group], error)
	// GetPluralGroupsWithSd SD付きでグループを取得する。
	GetPluralGroupsWithSd(
		ctx context.Context,
		sd Sd,
		groupIDs []uuid.UUID,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Group], error)
	// GetGroupsWithOrganization グループを取得する。
	GetGroupsWithOrganization(
		ctx context.Context,
		where parameter.WhereGroupParam,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.GroupWithOrganization], error)
	// GetGroupsWithOrganizationWithSd SD付きでグループを取得する。
	GetGroupsWithOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereGroupParam,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.GroupWithOrganization], error)
	// GetPluralGroupsWithOrganization グループを取得する。
	GetPluralGroupsWithOrganization(
		ctx context.Context,
		groupIDs []uuid.UUID,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.GroupWithOrganization], error)
	// GetPluralGroupsWithOrganizationWithSd SD付きでグループを取得する。
	GetPluralGroupsWithOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		groupIDs []uuid.UUID,
		order parameter.GroupOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.GroupWithOrganization], error)
}

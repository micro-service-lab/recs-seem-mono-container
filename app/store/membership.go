package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Membership チャットルームを表すインターフェース。
type Membership interface {
	// CountOrganizationsOnMember メンバー上のチャットルーム数を取得する。
	CountOrganizationsOnMember(
		ctx context.Context, memberID uuid.UUID, where parameter.WhereOrganizationOnMemberParam) (int64, error)
	// CountOrganizationsOnMemberWithSd SD付きでメンバー上のチャットルーム数を取得する。
	CountOrganizationsOnMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID, where parameter.WhereOrganizationOnMemberParam) (int64, error)
	// CountMembersOnOrganization チャットルーム上のメンバー数を取得する。
	CountMembersOnOrganization(
		ctx context.Context, organizationID uuid.UUID, where parameter.WhereMemberOnOrganizationParam) (int64, error)
	// CountMembersOnOrganizationWithSd SD付きでチャットルーム上のメンバー数を取得する。
	CountMembersOnOrganizationWithSd(
		ctx context.Context, sd Sd, organizationID uuid.UUID, where parameter.WhereMemberOnOrganizationParam) (int64, error)
	// BelongOrganization メンバーをチャットルームに所属させる。
	BelongOrganization(ctx context.Context, param parameter.BelongOrganizationParam) (entity.Membership, error)
	// BelongOrganizationWithSd SD付きでメンバーをチャットルームに所属させる。
	BelongOrganizationWithSd(
		ctx context.Context, sd Sd, param parameter.BelongOrganizationParam) (entity.Membership, error)
	// BelongOrganizations メンバーを複数のチャットルームに所属させる。
	BelongOrganizations(ctx context.Context, params []parameter.BelongOrganizationParam) (int64, error)
	// BelongOrganizationsWithSd SD付きでメンバーを複数のチャットルームに所属させる。
	BelongOrganizationsWithSd(ctx context.Context, sd Sd, params []parameter.BelongOrganizationParam) (int64, error)
	// DisbelongOrganization メンバーをチャットルームから所属解除する。
	DisbelongOrganization(ctx context.Context, memberID, organizationID uuid.UUID) (int64, error)
	// DisbelongOrganizationWithSd SD付きでメンバーをチャットルームから所属解除する。
	DisbelongOrganizationWithSd(ctx context.Context, sd Sd, memberID, organizationID uuid.UUID) (int64, error)
	// DisbelongOrganizationOnMember メンバー上のチャットルームから所属解除する。
	DisbelongOrganizationOnMember(ctx context.Context, memberID uuid.UUID) (int64, error)
	// DisbelongOrganizationOnMemberWithSd SD付きでメンバー上のチャットルームから所属解除する。
	DisbelongOrganizationOnMemberWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (int64, error)
	// DisbelongOrganizationOnMembers メンバー上の複数のチャットルームから所属解除する。
	DisbelongOrganizationOnMembers(ctx context.Context, memberIDs []uuid.UUID) (int64, error)
	// DisbelongOrganizationOnMembersWithSd SD付きでメンバー上の複数のチャットルームから所属解除する。
	DisbelongOrganizationOnMembersWithSd(ctx context.Context, sd Sd, memberIDs []uuid.UUID) (int64, error)
	// DisbelongOrganizationOnOrganization チャットルーム上のメンバーから所属解除する。
	DisbelongOrganizationOnOrganization(ctx context.Context, organizationID uuid.UUID) (int64, error)
	// DisbelongOrganizationOnOrganizationWithSd SD付きでチャットルーム上のメンバーから所属解除する。
	DisbelongOrganizationOnOrganizationWithSd(ctx context.Context, sd Sd, organizationID uuid.UUID) (int64, error)
	// DisbelongOrganizationOnOrganizations チャットルーム上の複数のメンバーから所属解除する。
	DisbelongOrganizationOnOrganizations(ctx context.Context, organizationIDs []uuid.UUID) (int64, error)
	// DisbelongOrganizationOnOrganizationsWithSd SD付きでチャットルーム上の複数のメンバーから所属解除する。
	DisbelongOrganizationOnOrganizationsWithSd(ctx context.Context, sd Sd, organizationIDs []uuid.UUID) (int64, error)
	// GetOrganizationsOnMember メンバー上のチャットルームを取得する。
	GetOrganizationsOnMember(
		ctx context.Context,
		memberID uuid.UUID,
		where parameter.WhereOrganizationOnMemberParam,
		order parameter.OrganizationOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationOnMember], error)
	// GetOrganizationsOnMemberWithSd SD付きでメンバー上のチャットルームを取得する。
	GetOrganizationsOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberID uuid.UUID,
		where parameter.WhereOrganizationOnMemberParam,
		order parameter.OrganizationOnMemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.OrganizationOnMember], error)
	// GetMembersOnOrganization チャットルーム上のメンバーを取得する。
	GetMembersOnOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
		where parameter.WhereMemberOnOrganizationParam,
		order parameter.MemberOnOrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnOrganization], error)
	// GetMembersOnOrganizationWithSd SD付きでチャットルーム上のメンバーを取得する。
	GetMembersOnOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		organizationID uuid.UUID,
		where parameter.WhereMemberOnOrganizationParam,
		order parameter.MemberOnOrganizationOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberOnOrganization], error)
	// GetPluralOrganizationsOnMember メンバー上の複数のチャットルームを取得する。
	GetPluralOrganizationsOnMember(
		ctx context.Context,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.OrganizationOnMemberOrderMethod,
	) (ListResult[entity.OrganizationOnMember], error)
	// GetPluralOrganizationsOnMemberWithSd SD付きでメンバー上の複数のチャットルームを取得する。
	GetPluralOrganizationsOnMemberWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.OrganizationOnMemberOrderMethod,
	) (ListResult[entity.OrganizationOnMember], error)
	// GetPluralMembersOnOrganization チャットルーム上の複数のメンバーを取得する。
	GetPluralMembersOnOrganization(
		ctx context.Context,
		organizationIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.MemberOnOrganizationOrderMethod,
	) (ListResult[entity.MemberOnOrganization], error)
	// GetPluralMembersOnOrganizationWithSd SD付きでチャットルーム上の複数のメンバーを取得する。
	GetPluralMembersOnOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		organizationIDs []uuid.UUID,
		np NumberedPaginationParam,
		order parameter.MemberOnOrganizationOrderMethod,
	) (ListResult[entity.MemberOnOrganization], error)
}

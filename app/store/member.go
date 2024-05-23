package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Member メンバーを表すインターフェース。
type Member interface {
	// CountMembers メンバー数を取得する。
	CountMembers(ctx context.Context, where parameter.WhereMemberParam) (int64, error)
	// CountMembersWithSd SD付きでメンバー数を取得する。
	CountMembersWithSd(ctx context.Context, sd Sd, where parameter.WhereMemberParam) (int64, error)
	// CreateMember メンバーを作成する。
	CreateMember(ctx context.Context, param parameter.CreateMemberParam) (entity.Member, error)
	// CreateMemberWithSd SD付きでメンバーを作成する。
	CreateMemberWithSd(
		ctx context.Context, sd Sd, param parameter.CreateMemberParam) (entity.Member, error)
	// CreateMembers メンバーを作成する。
	CreateMembers(ctx context.Context, params []parameter.CreateMemberParam) (int64, error)
	// CreateMembersWithSd SD付きでメンバーを作成する。
	CreateMembersWithSd(ctx context.Context, sd Sd, params []parameter.CreateMemberParam) (int64, error)
	// DeleteMember メンバーを削除する。
	DeleteMember(ctx context.Context, memberID uuid.UUID) (int64, error)
	// DeleteMemberWithSd SD付きでメンバーを削除する。
	DeleteMemberWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (int64, error)
	// PluralDeleteMembers メンバーを複数削除する。
	PluralDeleteMembers(ctx context.Context, memberIDs []uuid.UUID) (int64, error)
	// PluralDeleteMembersWithSd SD付きでメンバーを複数削除する。
	PluralDeleteMembersWithSd(ctx context.Context, sd Sd, memberIDs []uuid.UUID) (int64, error)
	// FindMemberByID メンバーを取得する。
	FindMemberByID(ctx context.Context, memberID uuid.UUID) (entity.Member, error)
	// FindMemberByIDWithSd SD付きでメンバーを取得する。
	FindMemberByIDWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.Member, error)
	// FindMemberCredentialsByID メンバーの認証情報を取得する。
	FindMemberCredentialsByID(ctx context.Context, memberID uuid.UUID) (entity.MemberCredentials, error)
	// FindMemberCredentialsByIDWithSd SD付きでメンバーの認証情報を取得する。
	FindMemberCredentialsByIDWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberCredentials, error)
	// FindMemberByLoginID メンバーを取得する。
	FindMemberByLoginID(ctx context.Context, loginID string) (entity.Member, error)
	// FindMemberByLoginIDWithSd SD付きでメンバーを取得する。
	FindMemberByLoginIDWithSd(ctx context.Context, sd Sd, loginID string) (entity.Member, error)
	// FindMemberCredentialsByLoginID メンバーの認証情報を取得する。
	FindMemberCredentialsByLoginID(ctx context.Context, loginID string) (entity.MemberCredentials, error)
	// FindMemberCredentialsByLoginIDWithSd SD付きでメンバーの認証情報を取得する。
	FindMemberCredentialsByLoginIDWithSd(ctx context.Context, sd Sd, loginID string) (entity.MemberCredentials, error)
	// FindMemberWithAttendStatus メンバーを取得する。
	FindMemberWithAttendStatus(ctx context.Context, memberID uuid.UUID) (entity.MemberWithAttendStatus, error)
	// FindMemberWithAttendStatusWithSd SD付きでメンバーを取得する。
	FindMemberWithAttendStatusWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithAttendStatus, error)
	// FindMemberWithProfileImage メンバーを取得する。
	FindMemberWithProfileImage(ctx context.Context, memberID uuid.UUID) (entity.MemberWithProfileImage, error)
	// FindMemberWithProfileImageWithSd SD付きでメンバーを取得する。
	FindMemberWithProfileImageWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithProfileImage, error)
	// FindMemberWithCrew メンバーを取得する。
	FindMemberWithCrew(ctx context.Context, memberID uuid.UUID) (entity.MemberWithCrew, error)
	// FindMemberWithCrewWithSd SD付きでメンバーを取得する。
	FindMemberWithCrewWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithCrew, error)
	// FindMemberWithPersonalOrganization メンバーを取得する。
	FindMemberWithPersonalOrganization(
		ctx context.Context, memberID uuid.UUID) (entity.MemberWithPersonalOrganization, error)
	// FindMemberWithPersonalOrganizationWithSd SD付きでメンバーを取得する。
	FindMemberWithPersonalOrganizationWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithPersonalOrganization, error)
	// FindMemberWithRole メンバーを取得する。
	FindMemberWithRole(ctx context.Context, memberID uuid.UUID) (entity.MemberWithRole, error)
	// FindMemberWithRoleWithSd SD付きでメンバーを取得する。
	FindMemberWithRoleWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithRole, error)
	// FindMemberWithDetail メンバーを取得する。
	FindMemberWithDetail(ctx context.Context, memberID uuid.UUID) (entity.MemberWithDetail, error)
	// FindMemberWithDetailWithSd SD付きでメンバーを取得する。
	FindMemberWithDetailWithSd(ctx context.Context, sd Sd, memberID uuid.UUID) (entity.MemberWithDetail, error)
	// GetMembers メンバーを取得する。
	GetMembers(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Member], error)
	// GetMembersWithSd SD付きでメンバーを取得する。
	GetMembersWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Member], error)
	// GetPluralMembers メンバーを取得する。
	GetPluralMembers(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Member], error)
	// GetPluralMembersWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Member], error)
	// GetMembersWithAttendStatus メンバーを取得する。
	GetMembersWithAttendStatus(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithAttendStatus], error)
	// GetMembersWithAttendStatusWithSd SD付きでメンバーを取得する。
	GetMembersWithAttendStatusWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithAttendStatus], error)
	// GetPluralMembersWithAttendStatus メンバーを取得する。
	GetPluralMembersWithAttendStatus(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithAttendStatus], error)
	// GetPluralMembersWithAttendStatusWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithAttendStatusWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithAttendStatus], error)
	// GetMembersWithDetail メンバーを取得する。
	GetMembersWithDetail(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithDetail], error)
	// GetMembersWithDetailWithSd SD付きでメンバーを取得する。
	GetMembersWithDetailWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithDetail], error)
	// GetPluralMembersWithDetail メンバーを取得する。
	GetPluralMembersWithDetail(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithDetail], error)
	// GetPluralMembersWithDetailWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithDetailWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithDetail], error)
	// GetMembersWithProfileImage メンバーを取得する。
	GetMembersWithProfileImage(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithProfileImage], error)
	// GetMembersWithProfileImageWithSd SD付きでメンバーを取得する。
	GetMembersWithProfileImageWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithProfileImage], error)
	// GetPluralMembersWithProfileImage メンバーを取得する。
	GetPluralMembersWithProfileImage(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithProfileImage], error)
	// GetPluralMembersWithProfileImageWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithProfileImageWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithProfileImage], error)
	// GetMembersWithCrew メンバーを取得する。
	GetMembersWithCrew(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithCrew], error)
	// GetMembersWithCrewWithSd SD付きでメンバーを取得する。
	GetMembersWithCrewWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithCrew], error)
	// GetPluralMembersWithCrew メンバーを取得する。
	GetPluralMembersWithCrew(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithCrew], error)
	// GetPluralMembersWithCrewWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithCrewWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithCrew], error)
	// GetMembersWithRole メンバーを取得する。
	GetMembersWithRole(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithRole], error)
	// GetMembersWithRoleWithSd SD付きでメンバーを取得する。
	GetMembersWithRoleWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithRole], error)
	// GetPluralMembersWithRole メンバーを取得する。
	GetPluralMembersWithRole(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithRole], error)
	// GetPluralMembersWithRoleWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithRoleWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithRole], error)
	// GetMembersWithPersonalOrganization メンバーを取得する。
	GetMembersWithPersonalOrganization(
		ctx context.Context,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithPersonalOrganization], error)
	// GetMembersWithPersonalOrganizationWithSd SD付きでメンバーを取得する。
	GetMembersWithPersonalOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereMemberParam,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.MemberWithPersonalOrganization], error)
	// GetPluralMembersWithPersonalOrganization メンバーを取得する。
	GetPluralMembersWithPersonalOrganization(
		ctx context.Context,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithPersonalOrganization], error)
	// GetPluralMembersWithPersonalOrganizationWithSd SD付きでメンバーを取得する。
	GetPluralMembersWithPersonalOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		memberIDs []uuid.UUID,
		order parameter.MemberOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.MemberWithPersonalOrganization], error)
	// UpdateMember メンバーを更新する。
	UpdateMember(
		ctx context.Context,
		memberID uuid.UUID,
		param parameter.UpdateMemberParams,
	) (entity.Member, error)
	// UpdateMemberWithSd SD付きでメンバーを更新する。
	UpdateMemberWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		param parameter.UpdateMemberParams) (entity.Member, error)
	// UpdateMemberAttendStatus メンバーの出席状況を更新する。
	UpdateMemberAttendStatus(
		ctx context.Context,
		memberID uuid.UUID,
		attendStatusID uuid.UUID,
	) (entity.Member, error)
	// UpdateMemberAttendStatusWithSd SD付きでメンバーの出席状況を更新する。
	UpdateMemberAttendStatusWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		attendStatusID uuid.UUID) (entity.Member, error)
	// UpdateMemberGrade メンバーの学年を更新する。
	UpdateMemberGrade(
		ctx context.Context,
		memberID uuid.UUID,
		gradeID uuid.UUID,
	) (entity.Member, error)
	// UpdateMemberGradeWithSd SD付きでメンバーの学年を更新する。
	UpdateMemberGradeWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		gradeID uuid.UUID) (entity.Member, error)
	// UpdateMemberGroup メンバーのグループを更新する。
	UpdateMemberGroup(
		ctx context.Context,
		memberID uuid.UUID,
		groupID uuid.UUID,
	) (entity.Member, error)
	// UpdateMemberGroupWithSd SD付きでメンバーのグループを更新する。
	UpdateMemberGroupWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		groupID uuid.UUID) (entity.Member, error)
	// UpdateMemberLoginID メンバーのログインIDを更新する。
	UpdateMemberLoginID(
		ctx context.Context,
		memberID uuid.UUID,
		loginID string,
	) (entity.Member, error)
	// UpdateMemberLoginIDWithSd SD付きでメンバーのログインIDを更新する。
	UpdateMemberLoginIDWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		loginID string) (entity.Member, error)
	// UpdateMemberPassword メンバーのパスワードを更新する。
	UpdateMemberPassword(
		ctx context.Context,
		memberID uuid.UUID,
		password string,
	) (entity.Member, error)
	// UpdateMemberPasswordWithSd SD付きでメンバーのパスワードを更新する。
	UpdateMemberPasswordWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		password string) (entity.Member, error)
	// UpdateMemberRole メンバーの権限を更新する。
	UpdateMemberRole(
		ctx context.Context,
		memberID uuid.UUID,
		roleID entity.UUID,
	) (entity.Member, error)
	// UpdateMemberRoleWithSd SD付きでメンバーの権限を更新する。
	UpdateMemberRoleWithSd(
		ctx context.Context, sd Sd, memberID uuid.UUID,
		roleID entity.UUID) (entity.Member, error)
}

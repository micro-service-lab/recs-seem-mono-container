package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convMember(e query.Member) entity.Member {
	return entity.Member{
		MemberID:               e.MemberID,
		Email:                  e.Email,
		Name:                   e.Name,
		FirstName:              e.FirstName,
		LastName:               e.LastName,
		AttendStatusID:         e.AttendStatusID,
		ProfileImageID:         entity.UUID(e.ProfileImageID),
		GradeID:                e.GradeID,
		GroupID:                e.GroupID,
		PersonalOrganizationID: e.PersonalOrganizationID,
		RoleID:                 entity.UUID(e.RoleID),
	}
}

func convMemberCredentials(e query.Member) entity.MemberCredentials {
	return entity.MemberCredentials{
		MemberID: e.MemberID,
		LoginID:  e.LoginID,
		Password: e.Password,
	}
}

func convMemberWithDetail(e query.FindMemberByIDWithDetailRow) entity.MemberWithDetail {
	return entity.MemberWithDetail{
		MemberID:               e.MemberID,
		Email:                  e.Email,
		Name:                   e.Name,
		FirstName:              e.FirstName,
		LastName:               e.LastName,
		AttendStatusID:         e.AttendStatusID,
		ProfileImageID:         entity.UUID(e.ProfileImageID),
		GradeID:                e.GradeID,
		GroupID:                e.GroupID,
		PersonalOrganizationID: e.PersonalOrganizationID,
		RoleID:                 entity.UUID(e.RoleID),
		Student: entity.NullableEntity[entity.Student]{
			Valid: e.StudentID.Valid,
			Entity: entity.Student{
				StudentID: e.StudentID.Bytes,
				MemberID:  e.MemberID,
			},
		},
		Professor: entity.NullableEntity[entity.Professor]{
			Valid: e.ProfessorID.Valid,
			Entity: entity.Professor{
				ProfessorID: e.ProfessorID.Bytes,
				MemberID:    e.MemberID,
			},
		},
	}
}

func convMemberWithAttendStatus(e query.FindMemberByIDWithAttendStatusRow) entity.MemberWithAttendStatus {
	return entity.MemberWithAttendStatus{
		MemberID:  e.MemberID,
		Email:     e.Email,
		Name:      e.Name,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		AttendStatus: entity.AttendStatus{
			AttendStatusID: e.AttendStatusID,
			Name:           e.AttendStatusName.String,
			Key:            e.AttendStatusKey.String,
		},
		ProfileImageID:         entity.UUID(e.ProfileImageID),
		GradeID:                e.GradeID,
		GroupID:                e.GroupID,
		PersonalOrganizationID: e.PersonalOrganizationID,
		RoleID:                 entity.UUID(e.RoleID),
	}
}

func convMemberWithProfileImage(e query.FindMemberByIDWithProfileImageRow) entity.MemberWithProfileImage {
	return entity.MemberWithProfileImage{
		MemberID:       e.MemberID,
		Email:          e.Email,
		Name:           e.Name,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		AttendStatusID: e.AttendStatusID,
		ProfileImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
			Valid: e.ProfileImageID.Valid,
			Entity: entity.ImageWithAttachableItem{
				ImageID: e.ProfileImageID.Bytes,
				AttachableItem: entity.AttachableItem{
					AttachableItemID: e.ProfileImageAttachableItemID.Bytes,
					OwnerID:          entity.UUID(e.ProfileImageOwnerID),
					FromOuter:        e.ProfileImageFromOuter.Bool,
					URL:              e.ProfileImageUrl.String,
					Alias:            e.ProfileImageAlias.String,
					Size:             entity.Float(e.ProfileImageSize),
					MimeTypeID:       e.ProfileImageMimeTypeID.Bytes,
				},
			},
		},
		GradeID:                e.GradeID,
		GroupID:                e.GroupID,
		PersonalOrganizationID: e.PersonalOrganizationID,
		RoleID:                 entity.UUID(e.RoleID),
	}
}

func convMemberWithCrew(e query.FindMemberByIDWithCrewRow) entity.MemberWithCrew {
	return entity.MemberWithCrew{
		MemberID:       e.MemberID,
		Email:          e.Email,
		Name:           e.Name,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		AttendStatusID: e.AttendStatusID,
		ProfileImageID: entity.UUID(e.ProfileImageID),
		Grade: entity.GradeWithOrganization{
			GradeID: e.GradeID,
			Key:     e.GradeKey.String,
			Organization: entity.Organization{
				OrganizationID: e.GradeOrganizationID.Bytes,
				Name:           e.GradeOrganizationName.String,
				Color:          entity.String(e.GradeOrganizationColor),
				Description:    entity.String(e.GradeOrganizationDescription),
				IsPersonal:     e.GradeOrganizationIsPersonal.Bool,
				IsWhole:        e.GradeOrganizationIsWhole.Bool,
				ChatRoomID:     entity.UUID(e.GradeOrganizationChatRoomID),
			},
		},
		Group: entity.GroupWithOrganization{
			GroupID: e.GroupID,
			Key:     e.GroupKey.String,
			Organization: entity.Organization{
				OrganizationID: e.GroupOrganizationID.Bytes,
				Name:           e.GroupOrganizationName.String,
				Color:          entity.String(e.GroupOrganizationColor),
				Description:    entity.String(e.GroupOrganizationDescription),
				IsPersonal:     e.GroupOrganizationIsPersonal.Bool,
				IsWhole:        e.GroupOrganizationIsWhole.Bool,
				ChatRoomID:     entity.UUID(e.GroupOrganizationChatRoomID),
			},
		},
		PersonalOrganizationID: e.PersonalOrganizationID,
		RoleID:                 entity.UUID(e.RoleID),
	}
}

func convMemberWithRole(e query.FindMemberByIDWithRoleRow) entity.MemberWithRole {
	return entity.MemberWithRole{
		MemberID:               e.MemberID,
		Email:                  e.Email,
		Name:                   e.Name,
		FirstName:              e.FirstName,
		LastName:               e.LastName,
		AttendStatusID:         e.AttendStatusID,
		ProfileImageID:         entity.UUID(e.ProfileImageID),
		GradeID:                e.GradeID,
		GroupID:                e.GroupID,
		PersonalOrganizationID: e.PersonalOrganizationID,
		Role: entity.NullableEntity[entity.Role]{
			Valid: e.RoleID.Valid,
			Entity: entity.Role{
				RoleID:      e.RoleID.Bytes,
				Name:        e.RoleName.String,
				Description: e.RoleDescription.String,
			},
		},
	}
}

func convMemberWithPersonalOrganization(
	e query.FindMemberByIDWithPersonalOrganizationRow,
) entity.MemberWithPersonalOrganization {
	return entity.MemberWithPersonalOrganization{
		MemberID:       e.MemberID,
		Email:          e.Email,
		Name:           e.Name,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		AttendStatusID: e.AttendStatusID,
		ProfileImageID: entity.UUID(e.ProfileImageID),
		GradeID:        e.GradeID,
		GroupID:        e.GroupID,
		PersonalOrganization: entity.Organization{
			OrganizationID: e.PersonalOrganizationID,
			Name:           e.OrganizationName.String,
			Color:          entity.String(e.OrganizationColor),
			Description:    entity.String(e.OrganizationDescription),
			IsPersonal:     e.OrganizationIsPersonal.Bool,
			IsWhole:        e.OrganizationIsWhole.Bool,
			ChatRoomID:     entity.UUID(e.OrganizationChatRoomID),
		},
		RoleID: entity.UUID(e.RoleID),
	}
}

// countMembers はメンバー数を取得する内部関数です。
func countMembers(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
) (int64, error) {
	p := query.CountMembersParams{
		WhereLikeName:      where.WhereLikeName,
		SearchName:         where.SearchName,
		WhereHasPolicy:     where.WhereHasPolicy,
		HasPolicyIds:       where.HasPolicyIDs,
		WhenInAttendStatus: where.WhenInAttendStatus,
		InAttendStatusIds:  where.InAttendStatusIDs,
		WhenInGrade:        where.WhenInGrade,
		InGradeIds:         where.InGradeIDs,
		WhenInGroup:        where.WhenInGroup,
		InGroupIds:         where.InGroupIDs,
	}
	c, err := qtx.CountMembers(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count members: %w", err)
	}
	return c, nil
}

// CountMembers はメンバー数を取得します。
func (a *PgAdapter) CountMembers(ctx context.Context, where parameter.WhereMemberParam) (int64, error) {
	return countMembers(ctx, a.query, where)
}

// CountMembersWithSd はSD付きでメンバー数を取得します。
func (a *PgAdapter) CountMembersWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMembers(ctx, qtx, where)
}

// createMember はメンバーを作成する内部関数です。
func createMember(
	ctx context.Context, qtx *query.Queries, param parameter.CreateMemberParam, now time.Time,
) (entity.Member, error) {
	p := query.CreateMemberParams{
		LoginID:                param.LoginID,
		Password:               param.Password,
		Email:                  param.Email,
		Name:                   param.Name,
		FirstName:              param.FirstName,
		LastName:               param.LastName,
		AttendStatusID:         param.AttendStatusID,
		GradeID:                param.GradeID,
		GroupID:                param.GroupID,
		ProfileImageID:         pgtype.UUID(param.ProfileImageID),
		RoleID:                 pgtype.UUID(param.RoleID),
		PersonalOrganizationID: param.PersonalOrganizationID,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	e, err := qtx.CreateMember(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Member{}, errhandle.NewModelDuplicatedError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to create member: %w", err)
	}
	return convMember(e), nil
}

// CreateMember はメンバーを作成します。
func (a *PgAdapter) CreateMember(
	ctx context.Context, param parameter.CreateMemberParam,
) (entity.Member, error) {
	return createMember(ctx, a.query, param, a.clocker.Now())
}

// CreateMemberWithSd はSD付きでメンバーを作成します。
func (a *PgAdapter) CreateMemberWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateMemberParam,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return createMember(ctx, qtx, param, a.clocker.Now())
}

// createMembers は複数のメンバーを作成する内部関数です。
func createMembers(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateMemberParam, now time.Time,
) (int64, error) {
	param := make([]query.CreateMembersParams, len(params))
	for i, p := range params {
		param[i] = query.CreateMembersParams{
			LoginID:                p.LoginID,
			Password:               p.Password,
			Email:                  p.Email,
			Name:                   p.Name,
			FirstName:              p.FirstName,
			LastName:               p.LastName,
			AttendStatusID:         p.AttendStatusID,
			GradeID:                p.GradeID,
			GroupID:                p.GroupID,
			ProfileImageID:         pgtype.UUID(p.ProfileImageID),
			RoleID:                 pgtype.UUID(p.RoleID),
			PersonalOrganizationID: p.PersonalOrganizationID,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
	}
	n, err := qtx.CreateMembers(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("member")
		}
		return 0, fmt.Errorf("failed to create members: %w", err)
	}
	return n, nil
}

// CreateMembers は複数のメンバーを作成します。
func (a *PgAdapter) CreateMembers(
	ctx context.Context, params []parameter.CreateMemberParam,
) (int64, error) {
	return createMembers(ctx, a.query, params, a.clocker.Now())
}

// CreateMembersWithSd はSD付きで複数のメンバーを作成します。
func (a *PgAdapter) CreateMembersWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createMembers(ctx, qtx, params, a.clocker.Now())
}

// deleteMember はメンバーを削除する内部関数です。
func deleteMember(ctx context.Context, qtx *query.Queries, memberID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteMember(ctx, memberID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete member: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("member")
	}
	return c, nil
}

// DeleteMember はメンバーを削除します。
func (a *PgAdapter) DeleteMember(ctx context.Context, memberID uuid.UUID) (int64, error) {
	return deleteMember(ctx, a.query, memberID)
}

// DeleteMemberWithSd はSD付きでメンバーを削除します。
func (a *PgAdapter) DeleteMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteMember(ctx, qtx, memberID)
}

// pluralDeleteMembers は複数のメンバーを削除する内部関数です。
func pluralDeleteMembers(ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteMembers(ctx, memberIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete members: %w", err)
	}
	if c != int64(len(memberIDs)) {
		return 0, errhandle.NewModelNotFoundError("member")
	}
	return c, nil
}

// PluralDeleteMembers は複数のメンバーを削除します。
func (a *PgAdapter) PluralDeleteMembers(ctx context.Context, memberIDs []uuid.UUID) (int64, error) {
	return pluralDeleteMembers(ctx, a.query, memberIDs)
}

// PluralDeleteMembersWithSd はSD付きで複数のメンバーを削除します。
func (a *PgAdapter) PluralDeleteMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteMembers(ctx, qtx, memberIDs)
}

// findMemberByID はメンバーをIDで取得する内部関数です。
func findMemberByID(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.Member, error) {
	e, err := qtx.FindMemberByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to find member: %w", err)
	}
	return convMember(e), nil
}

// FindMemberByID はメンバーをIDで取得します。
func (a *PgAdapter) FindMemberByID(ctx context.Context, memberID uuid.UUID) (entity.Member, error) {
	return findMemberByID(ctx, a.query, memberID)
}

// FindMemberByIDWithSd はSD付きでメンバーをIDで取得します。
func (a *PgAdapter) FindMemberByIDWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return findMemberByID(ctx, qtx, memberID)
}

func findMemberCredentialsByID(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberCredentials, error) {
	e, err := qtx.FindMemberByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberCredentials{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberCredentials{}, fmt.Errorf("failed to find member: %w", err)
	}
	return convMemberCredentials(e), nil
}

// FindMemberCredentialsByID はメンバーの認証情報をIDで取得します。
func (a *PgAdapter) FindMemberCredentialsByID(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberCredentials, error) {
	return findMemberCredentialsByID(ctx, a.query, memberID)
}

// FindMemberCredentialsByIDWithSd はSD付きでメンバーの認証情報をIDで取得します。
func (a *PgAdapter) FindMemberCredentialsByIDWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberCredentials, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberCredentials{}, store.ErrNotFoundDescriptor
	}
	return findMemberCredentialsByID(ctx, qtx, memberID)
}

func findMemberByLoginID(
	ctx context.Context, qtx *query.Queries, loginID string,
) (entity.Member, error) {
	e, err := qtx.FindMemberByLoginID(ctx, loginID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to find member by login ID: %w", err)
	}
	return convMember(e), nil
}

// FindMemberByLoginID はメンバーをログインIDで取得します。
func (a *PgAdapter) FindMemberByLoginID(
	ctx context.Context, loginID string,
) (entity.Member, error) {
	return findMemberByLoginID(ctx, a.query, loginID)
}

// FindMemberByLoginIDWithSd はSD付きでメンバーをログインIDで取得します。
func (a *PgAdapter) FindMemberByLoginIDWithSd(
	ctx context.Context, sd store.Sd, loginID string,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return findMemberByLoginID(ctx, qtx, loginID)
}

func findMemberCredentialsByLoginID(
	ctx context.Context, qtx *query.Queries, loginID string,
) (entity.MemberCredentials, error) {
	e, err := qtx.FindMemberByLoginID(ctx, loginID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberCredentials{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberCredentials{}, fmt.Errorf("failed to find member by login ID: %w", err)
	}
	return convMemberCredentials(e), nil
}

// FindMemberCredentialsByLoginID はメンバーの認証情報をログインIDで取得します。
func (a *PgAdapter) FindMemberCredentialsByLoginID(
	ctx context.Context, loginID string,
) (entity.MemberCredentials, error) {
	return findMemberCredentialsByLoginID(ctx, a.query, loginID)
}

// FindMemberCredentialsByLoginIDWithSd はSD付きでメンバーの認証情報をログインIDで取得します。
func (a *PgAdapter) FindMemberCredentialsByLoginIDWithSd(
	ctx context.Context, sd store.Sd, loginID string,
) (entity.MemberCredentials, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberCredentials{}, store.ErrNotFoundDescriptor
	}
	return findMemberCredentialsByLoginID(ctx, qtx, loginID)
}

// findMemberWithAttendStatus はメンバーと出席状況を取得する内部関数です。
func findMemberWithAttendStatus(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithAttendStatus, error) {
	e, err := qtx.FindMemberByIDWithAttendStatus(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithAttendStatus{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithAttendStatus{}, fmt.Errorf("failed to find member with attend status: %w", err)
	}
	return convMemberWithAttendStatus(e), nil
}

// FindMemberWithAttendStatus はメンバーと出席状況を取得します。
func (a *PgAdapter) FindMemberWithAttendStatus(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithAttendStatus, error) {
	return findMemberWithAttendStatus(ctx, a.query, memberID)
}

// FindMemberWithAttendStatusWithSd はSD付きでメンバーと出席状況を取得します。
func (a *PgAdapter) FindMemberWithAttendStatusWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithAttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithAttendStatus{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithAttendStatus(ctx, qtx, memberID)
}

// findMemberWithProfileImage はメンバーとプロフィール画像を取得する内部関数です。
func findMemberWithProfileImage(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithProfileImage, error) {
	e, err := qtx.FindMemberByIDWithProfileImage(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithProfileImage{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithProfileImage{}, fmt.Errorf("failed to find member with profile image: %w", err)
	}
	return convMemberWithProfileImage(e), nil
}

// FindMemberWithProfileImage はメンバーとプロフィール画像を取得します。
func (a *PgAdapter) FindMemberWithProfileImage(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithProfileImage, error) {
	return findMemberWithProfileImage(ctx, a.query, memberID)
}

// FindMemberWithProfileImageWithSd はSD付きでメンバーとプロフィール画像を取得します。
func (a *PgAdapter) FindMemberWithProfileImageWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithProfileImage, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithProfileImage{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithProfileImage(ctx, qtx, memberID)
}

// findMemberWithDetail はメンバーと詳細情報を取得する内部関数です。
func findMemberWithDetail(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithDetail, error) {
	e, err := qtx.FindMemberByIDWithDetail(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithDetail{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithDetail{}, fmt.Errorf("failed to find member with detail: %w", err)
	}
	return convMemberWithDetail(e), nil
}

// FindMemberWithDetail はメンバーと詳細情報を取得します。
func (a *PgAdapter) FindMemberWithDetail(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithDetail, error) {
	return findMemberWithDetail(ctx, a.query, memberID)
}

// FindMemberWithDetailWithSd はSD付きでメンバーと詳細情報を取得します。
func (a *PgAdapter) FindMemberWithDetailWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithDetail, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithDetail{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithDetail(ctx, qtx, memberID)
}

// findMemberWithCrew はメンバーとクルーを取得する内部関数です。
func findMemberWithCrew(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithCrew, error) {
	e, err := qtx.FindMemberByIDWithCrew(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithCrew{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithCrew{}, fmt.Errorf("failed to find member with crew: %w", err)
	}
	return convMemberWithCrew(e), nil
}

// FindMemberWithCrew はメンバーとクルーを取得します。
func (a *PgAdapter) FindMemberWithCrew(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithCrew, error) {
	return findMemberWithCrew(ctx, a.query, memberID)
}

// FindMemberWithCrewWithSd はSD付きでメンバーとクルーを取得します。
func (a *PgAdapter) FindMemberWithCrewWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithCrew, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithCrew{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithCrew(ctx, qtx, memberID)
}

// findMemberWithRole はメンバーとロールを取得する内部関数です。
func findMemberWithRole(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithRole, error) {
	e, err := qtx.FindMemberByIDWithRole(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithRole{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithRole{}, fmt.Errorf("failed to find member with crew: %w", err)
	}
	return convMemberWithRole(e), nil
}

// FindMemberWithRole はメンバーとロールを取得します。
func (a *PgAdapter) FindMemberWithRole(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithRole, error) {
	return findMemberWithRole(ctx, a.query, memberID)
}

// FindMemberWithRoleWithSd はSD付きでメンバーとロールを取得します。
func (a *PgAdapter) FindMemberWithRoleWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithRole, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithRole{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithRole(ctx, qtx, memberID)
}

func findMemberWithPersonalOrganization(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.MemberWithPersonalOrganization, error) {
	e, err := qtx.FindMemberByIDWithPersonalOrganization(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MemberWithPersonalOrganization{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.MemberWithPersonalOrganization{},
			fmt.Errorf("failed to find member with personal organization: %w", err)
	}
	return convMemberWithPersonalOrganization(e), nil
}

// FindMemberWithPersonalOrganization はメンバーと個人組織を取得します。
func (a *PgAdapter) FindMemberWithPersonalOrganization(
	ctx context.Context, memberID uuid.UUID,
) (entity.MemberWithPersonalOrganization, error) {
	return findMemberWithPersonalOrganization(ctx, a.query, memberID)
}

// FindMemberWithPersonalOrganizationWithSd はSD付きでメンバーと個人組織を取得します。
func (a *PgAdapter) FindMemberWithPersonalOrganizationWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.MemberWithPersonalOrganization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MemberWithPersonalOrganization{}, store.ErrNotFoundDescriptor
	}
	return findMemberWithPersonalOrganization(ctx, qtx, memberID)
}

func convCountMembersParams(p parameter.WhereMemberParam) query.CountMembersParams {
	return query.CountMembersParams{
		WhereLikeName:      p.WhereLikeName,
		SearchName:         p.SearchName,
		WhereHasPolicy:     p.WhereHasPolicy,
		HasPolicyIds:       p.HasPolicyIDs,
		WhenInAttendStatus: p.WhenInAttendStatus,
		InAttendStatusIds:  p.InAttendStatusIDs,
		WhenInGrade:        p.WhenInGrade,
		InGradeIds:         p.InGradeIDs,
		WhenInGroup:        p.WhenInGroup,
		InGroupIds:         p.InGroupIDs,
	}
}

func convGetMembersParams(p parameter.WhereMemberParam, orderMethod string) query.GetMembersParams {
	return query.GetMembersParams{
		WhereLikeName:      p.WhereLikeName,
		SearchName:         p.SearchName,
		WhereHasPolicy:     p.WhereHasPolicy,
		HasPolicyIds:       p.HasPolicyIDs,
		WhenInAttendStatus: p.WhenInAttendStatus,
		InAttendStatusIds:  p.InAttendStatusIDs,
		WhenInGrade:        p.WhenInGrade,
		InGradeIds:         p.InGradeIDs,
		WhenInGroup:        p.WhenInGroup,
		InGroupIds:         p.InGroupIDs,
		OrderMethod:        orderMethod,
	}
}

func convGetMembersUseKeysetPaginateParams(p parameter.WhereMemberParam,
	subCursor, orderMethod string,
	limit int32, cursorDir string, cursor int32, subCursorValue any,
) query.GetMembersUseKeysetPaginateParams {
	var nameCursor string
	var ok bool
	switch subCursor {
	case parameter.MemberNameCursorKey:
		nameCursor, ok = subCursorValue.(string)
		if !ok {
			nameCursor = ""
		}
	}
	return query.GetMembersUseKeysetPaginateParams{
		WhereLikeName:      p.WhereLikeName,
		SearchName:         p.SearchName,
		WhereHasPolicy:     p.WhereHasPolicy,
		HasPolicyIds:       p.HasPolicyIDs,
		WhenInAttendStatus: p.WhenInAttendStatus,
		InAttendStatusIds:  p.InAttendStatusIDs,
		WhenInGrade:        p.WhenInGrade,
		InGradeIds:         p.InGradeIDs,
		WhenInGroup:        p.WhenInGroup,
		InGroupIds:         p.InGroupIDs,
		OrderMethod:        orderMethod,
		Limit:              limit,
		CursorDirection:    cursorDir,
		Cursor:             cursor,
		NameCursor:         nameCursor,
	}
}

func convGetMembersUseNumberedPaginateParams(
	p parameter.WhereMemberParam, orderMethod string, limit, offset int32,
) query.GetMembersUseNumberedPaginateParams {
	return query.GetMembersUseNumberedPaginateParams{
		WhereLikeName:      p.WhereLikeName,
		SearchName:         p.SearchName,
		WhereHasPolicy:     p.WhereHasPolicy,
		HasPolicyIds:       p.HasPolicyIDs,
		WhenInAttendStatus: p.WhenInAttendStatus,
		InAttendStatusIds:  p.InAttendStatusIDs,
		WhenInGrade:        p.WhenInGrade,
		InGradeIds:         p.InGradeIDs,
		WhenInGroup:        p.WhenInGroup,
		InGroupIds:         p.InGroupIDs,
		OrderMethod:        orderMethod,
		Limit:              limit,
		Offset:             offset,
	}
}

// getMembers はメンバーを取得する内部関数です。
func getMembers(
	ctx context.Context,
	qtx *query.Queries,
	where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Member], error) {
	eConvFunc := func(e query.Member) (entity.Member, error) {
		return convMember(e), nil
	}
	runCFunc := func() (int64, error) {
		p := convCountMembersParams(where)
		r, err := qtx.CountMembers(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Member, error) {
		p := convGetMembersParams(where, orderMethod)
		r, err := qtx.GetMembers(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Member{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Member, error) {
		r, err := qtx.GetMembersUseKeysetPaginate(ctx, convGetMembersUseKeysetPaginateParams(
			where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Member, error) {
		r, err := qtx.GetMembersUseNumberedPaginate(ctx, convGetMembersUseNumberedPaginateParams(
			where, orderMethod, limit, offset,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Member) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.MMembersPkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.MMembersPkey), e.Name
		}
		return entity.Int(e.MMembersPkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.Member]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembers はメンバーを取得します。
func (a *PgAdapter) GetMembers(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Member], error) {
	return getMembers(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithSd はSD付きでメンバーを取得します。
func (a *PgAdapter) GetMembersWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Member], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Member]{}, store.ErrNotFoundDescriptor
	}
	return getMembers(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembers は複数のメンバーを取得する内部関数です。
func getPluralMembers(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Member], error) {
	var e []query.Member
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembers(ctx, query.GetPluralMembersParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralMembersUseNumberedPaginate(ctx, query.GetPluralMembersUseNumberedPaginateParams{
			MemberIds:   memberIDs,
			Offset:      int32(np.Offset.Int64),
			Limit:       int32(np.Limit.Int64),
			OrderMethod: orderMethod.GetStringValue(),
		})
	}
	if err != nil {
		return store.ListResult[entity.Member]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.Member, len(e))
	for i, v := range e {
		entities[i] = convMember(v)
	}
	return store.ListResult[entity.Member]{Data: entities}, nil
}

// GetPluralMembers は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembers(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Member], error) {
	return getPluralMembers(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Member], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Member]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembers(ctx, qtx, memberIDs, order, np)
}

func getMembersWithAttendStatus(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	eConvFunc := func(e entity.MemberWithAttendStatusForQuery) (entity.MemberWithAttendStatus, error) {
		return e.MemberWithAttendStatus, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithAttendStatusForQuery, error) {
		r, err := qtx.GetMembersWithAttendStatus(
			ctx, query.GetMembersWithAttendStatusParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithAttendStatusForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithAttendStatusForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithAttendStatusForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithAttendStatus: convMemberWithAttendStatus(query.FindMemberByIDWithAttendStatusRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithAttendStatusForQuery, error) {
		r, err := qtx.GetMembersWithAttendStatusUseKeysetPaginate(
			ctx, query.GetMembersWithAttendStatusUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithAttendStatusForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithAttendStatusForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithAttendStatus: convMemberWithAttendStatus(query.FindMemberByIDWithAttendStatusRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithAttendStatusForQuery, error) {
		r, err := qtx.GetMembersWithAttendStatusUseNumberedPaginate(
			ctx, query.GetMembersWithAttendStatusUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithAttendStatusForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithAttendStatusForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithAttendStatus: convMemberWithAttendStatus(query.FindMemberByIDWithAttendStatusRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithAttendStatusForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithAttendStatus]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithAttendStatus はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithAttendStatus(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	return getMembersWithAttendStatus(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithAttendStatusWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithAttendStatusWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithAttendStatus]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithAttendStatus(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithAttendStatus は複数のメンバーを取得する内部関数です。
func getPluralMembersWithAttendStatus(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	var e []query.GetPluralMembersWithAttendStatusRow
	var te []query.GetPluralMembersWithAttendStatusUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithAttendStatus(ctx, query.GetPluralMembersWithAttendStatusParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithAttendStatusUseNumberedPaginate(
			ctx, query.GetPluralMembersWithAttendStatusUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithAttendStatusRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithAttendStatusRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithAttendStatus]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithAttendStatus, len(e))
	for i, v := range e {
		entities[i] = convMemberWithAttendStatus(query.FindMemberByIDWithAttendStatusRow(v))
	}
	return store.ListResult[entity.MemberWithAttendStatus]{Data: entities}, nil
}

// GetPluralMembersWithAttendStatus は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithAttendStatus(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	return getPluralMembersWithAttendStatus(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithAttendStatusWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithAttendStatusWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithAttendStatus], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithAttendStatus]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithAttendStatus(ctx, qtx, memberIDs, order, np)
}

func getMembersWithProfileImage(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	eConvFunc := func(e entity.MemberWithProfileImageForQuery) (entity.MemberWithProfileImage, error) {
		return e.MemberWithProfileImage, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithProfileImageForQuery, error) {
		r, err := qtx.GetMembersWithProfileImage(
			ctx, query.GetMembersWithProfileImageParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithProfileImageForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithProfileImageForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithProfileImageForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithProfileImage: convMemberWithProfileImage(query.FindMemberByIDWithProfileImageRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithProfileImageForQuery, error) {
		r, err := qtx.GetMembersWithProfileImageUseKeysetPaginate(
			ctx, query.GetMembersWithProfileImageUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithProfileImageForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithProfileImageForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithProfileImage: convMemberWithProfileImage(query.FindMemberByIDWithProfileImageRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithProfileImageForQuery, error) {
		r, err := qtx.GetMembersWithProfileImageUseNumberedPaginate(
			ctx, query.GetMembersWithProfileImageUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithProfileImageForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithProfileImageForQuery{
				Pkey:                   entity.Int(v.MMembersPkey),
				MemberWithProfileImage: convMemberWithProfileImage(query.FindMemberByIDWithProfileImageRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithProfileImageForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithProfileImage]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithProfileImage はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithProfileImage(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	return getMembersWithProfileImage(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithProfileImageWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithProfileImageWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithProfileImage]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithProfileImage(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithProfileImage は複数のメンバーを取得する内部関数です。
func getPluralMembersWithProfileImage(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	var e []query.GetPluralMembersWithProfileImageRow
	var te []query.GetPluralMembersWithProfileImageUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithProfileImage(ctx, query.GetPluralMembersWithProfileImageParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithProfileImageUseNumberedPaginate(
			ctx, query.GetPluralMembersWithProfileImageUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithProfileImageRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithProfileImageRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithProfileImage]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithProfileImage, len(e))
	for i, v := range e {
		entities[i] = convMemberWithProfileImage(query.FindMemberByIDWithProfileImageRow(v))
	}
	return store.ListResult[entity.MemberWithProfileImage]{Data: entities}, nil
}

// GetPluralMembersWithProfileImage は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithProfileImage(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	return getPluralMembersWithProfileImage(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithProfileImageWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithProfileImageWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithProfileImage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithProfileImage]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithProfileImage(ctx, qtx, memberIDs, order, np)
}

func getMembersWithCrew(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	eConvFunc := func(e entity.MemberWithCrewForQuery) (entity.MemberWithCrew, error) {
		return e.MemberWithCrew, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithCrewForQuery, error) {
		r, err := qtx.GetMembersWithCrew(ctx, query.GetMembersWithCrewParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithCrewForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithCrewForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithCrewForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithCrew: convMemberWithCrew(query.FindMemberByIDWithCrewRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithCrewForQuery, error) {
		r, err := qtx.GetMembersWithCrewUseKeysetPaginate(
			ctx, query.GetMembersWithCrewUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithCrewForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithCrewForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithCrew: convMemberWithCrew(query.FindMemberByIDWithCrewRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithCrewForQuery, error) {
		r, err := qtx.GetMembersWithCrewUseNumberedPaginate(
			ctx, query.GetMembersWithCrewUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithCrewForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithCrewForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithCrew: convMemberWithCrew(query.FindMemberByIDWithCrewRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithCrewForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithCrew]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithCrew はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithCrew(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	return getMembersWithCrew(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithCrewWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithCrewWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithCrew]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithCrew(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithCrew は複数のメンバーを取得する内部関数です。
func getPluralMembersWithCrew(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	var e []query.GetPluralMembersWithCrewRow
	var te []query.GetPluralMembersWithCrewUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithCrew(ctx, query.GetPluralMembersWithCrewParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithCrewUseNumberedPaginate(
			ctx, query.GetPluralMembersWithCrewUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithCrewRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithCrewRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithCrew]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithCrew, len(e))
	for i, v := range e {
		entities[i] = convMemberWithCrew(query.FindMemberByIDWithCrewRow(v))
	}
	return store.ListResult[entity.MemberWithCrew]{Data: entities}, nil
}

// GetPluralMembersWithCrew は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithCrew(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	return getPluralMembersWithCrew(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithCrewWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithCrewWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithCrew], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithCrew]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithCrew(ctx, qtx, memberIDs, order, np)
}

func getMembersWithDetail(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	eConvFunc := func(e entity.MemberWithDetailForQuery) (entity.MemberWithDetail, error) {
		return e.MemberWithDetail, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithDetailForQuery, error) {
		r, err := qtx.GetMembersWithDetail(ctx, query.GetMembersWithDetailParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithDetailForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithDetailForQuery{
				Pkey:             entity.Int(v.MMembersPkey),
				MemberWithDetail: convMemberWithDetail(query.FindMemberByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithDetailForQuery, error) {
		r, err := qtx.GetMembersWithDetailUseKeysetPaginate(
			ctx, query.GetMembersWithDetailUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithDetailForQuery{
				Pkey:             entity.Int(v.MMembersPkey),
				MemberWithDetail: convMemberWithDetail(query.FindMemberByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithDetailForQuery, error) {
		r, err := qtx.GetMembersWithDetailUseNumberedPaginate(
			ctx, query.GetMembersWithDetailUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithDetailForQuery{
				Pkey:             entity.Int(v.MMembersPkey),
				MemberWithDetail: convMemberWithDetail(query.FindMemberByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithDetailForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithDetail]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithDetail はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithDetail(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	return getMembersWithDetail(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithDetailWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithDetailWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithDetail(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithDetail は複数のメンバーを取得する内部関数です。
func getPluralMembersWithDetail(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	var e []query.GetPluralMembersWithDetailRow
	var te []query.GetPluralMembersWithDetailUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithDetail(ctx, query.GetPluralMembersWithDetailParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithDetailUseNumberedPaginate(
			ctx, query.GetPluralMembersWithDetailUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithDetailRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithDetailRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithDetail]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithDetail, len(e))
	for i, v := range e {
		entities[i] = convMemberWithDetail(query.FindMemberByIDWithDetailRow(v))
	}
	return store.ListResult[entity.MemberWithDetail]{Data: entities}, nil
}

// GetPluralMembersWithDetail は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithDetail(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	return getPluralMembersWithDetail(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithDetailWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithDetailWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithDetail(ctx, qtx, memberIDs, order, np)
}

func getMembersWithRole(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithRole], error) {
	eConvFunc := func(e entity.MemberWithRoleForQuery) (entity.MemberWithRole, error) {
		return e.MemberWithRole, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithRoleForQuery, error) {
		r, err := qtx.GetMembersWithRole(ctx, query.GetMembersWithRoleParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithRoleForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithRoleForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithRoleForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithRole: convMemberWithRole(query.FindMemberByIDWithRoleRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithRoleForQuery, error) {
		r, err := qtx.GetMembersWithRoleUseKeysetPaginate(
			ctx, query.GetMembersWithRoleUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithRoleForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithRoleForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithRole: convMemberWithRole(query.FindMemberByIDWithRoleRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithRoleForQuery, error) {
		r, err := qtx.GetMembersWithRoleUseNumberedPaginate(
			ctx, query.GetMembersWithRoleUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithRoleForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithRoleForQuery{
				Pkey:           entity.Int(v.MMembersPkey),
				MemberWithRole: convMemberWithRole(query.FindMemberByIDWithRoleRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithRoleForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithRole]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithRole はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithRole(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithRole], error) {
	return getMembersWithRole(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithRoleWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithRoleWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithRole], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithRole]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithRole(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithRole は複数のメンバーを取得する内部関数です。
func getPluralMembersWithRole(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithRole], error) {
	var e []query.GetPluralMembersWithRoleRow
	var te []query.GetPluralMembersWithRoleUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithRole(ctx, query.GetPluralMembersWithRoleParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithRoleUseNumberedPaginate(
			ctx, query.GetPluralMembersWithRoleUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithRoleRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithRoleRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithRole]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithRole, len(e))
	for i, v := range e {
		entities[i] = convMemberWithRole(query.FindMemberByIDWithRoleRow(v))
	}
	return store.ListResult[entity.MemberWithRole]{Data: entities}, nil
}

// GetPluralMembersWithRole は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithRole(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithRole], error) {
	return getPluralMembersWithRole(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithRoleWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithRoleWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithRole], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithRole]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithRole(ctx, qtx, memberIDs, order, np)
}

func getMembersWithPersonalOrganization(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	eConvFunc := func(e entity.MemberWithPersonalOrganizationForQuery) (entity.MemberWithPersonalOrganization, error) {
		return e.MemberWithPersonalOrganization, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembers(ctx, convCountMembersParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count members: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberWithPersonalOrganizationForQuery, error) {
		r, err := qtx.GetMembersWithPersonalOrganization(
			ctx, query.GetMembersWithPersonalOrganizationParams(convGetMembersParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberWithPersonalOrganizationForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithPersonalOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithPersonalOrganizationForQuery{
				Pkey: entity.Int(v.MMembersPkey),
				MemberWithPersonalOrganization: convMemberWithPersonalOrganization(
					query.FindMemberByIDWithPersonalOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberWithPersonalOrganizationForQuery, error) {
		r, err := qtx.GetMembersWithPersonalOrganizationUseKeysetPaginate(
			ctx, query.GetMembersWithPersonalOrganizationUseKeysetPaginateParams(convGetMembersUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithPersonalOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithPersonalOrganizationForQuery{
				Pkey: entity.Int(v.MMembersPkey),
				MemberWithPersonalOrganization: convMemberWithPersonalOrganization(
					query.FindMemberByIDWithPersonalOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberWithPersonalOrganizationForQuery, error) {
		r, err := qtx.GetMembersWithPersonalOrganizationUseNumberedPaginate(
			ctx, query.GetMembersWithPersonalOrganizationUseNumberedPaginateParams(convGetMembersUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
		e := make([]entity.MemberWithPersonalOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MemberWithPersonalOrganizationForQuery{
				Pkey: entity.Int(v.MMembersPkey),
				MemberWithPersonalOrganization: convMemberWithPersonalOrganization(
					query.FindMemberByIDWithPersonalOrganizationRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MemberWithPersonalOrganizationForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberNameCursorKey:
			return entity.Int(e.Pkey), e.Name
		}
		return entity.Int(e.Pkey), nil
	}

	res, err := store.RunListQuery(
		ctx,
		order,
		np,
		cp,
		wc,
		eConvFunc,
		runCFunc,
		runQFunc,
		runQCPFunc,
		runQNPFunc,
		selector,
	)
	if err != nil {
		return store.ListResult[entity.MemberWithPersonalOrganization]{}, fmt.Errorf("failed to get members: %w", err)
	}
	return res, nil
}

// GetMembersWithPersonalOrganization はメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithPersonalOrganization(
	ctx context.Context, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	return getMembersWithPersonalOrganization(ctx, a.query, where, order, np, cp, wc)
}

// GetMembersWithPersonalOrganizationWithSd はSD付きでメンバーとチャットルームを取得します。
func (a *PgAdapter) GetMembersWithPersonalOrganizationWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMemberParam,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithPersonalOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getMembersWithPersonalOrganization(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMembersWithPersonalOrganization は複数のメンバーを取得する内部関数です。
func getPluralMembersWithPersonalOrganization(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	var e []query.GetPluralMembersWithPersonalOrganizationRow
	var te []query.GetPluralMembersWithPersonalOrganizationUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersWithPersonalOrganization(ctx, query.GetPluralMembersWithPersonalOrganizationParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMembersWithPersonalOrganizationUseNumberedPaginate(ctx,
			query.GetPluralMembersWithPersonalOrganizationUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMembersWithPersonalOrganizationRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMembersWithPersonalOrganizationRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberWithPersonalOrganization]{}, fmt.Errorf("failed to get members: %w", err)
	}
	entities := make([]entity.MemberWithPersonalOrganization, len(e))
	for i, v := range e {
		entities[i] = convMemberWithPersonalOrganization(query.FindMemberByIDWithPersonalOrganizationRow(v))
	}
	return store.ListResult[entity.MemberWithPersonalOrganization]{Data: entities}, nil
}

// GetPluralMembersWithPersonalOrganization は複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithPersonalOrganization(
	ctx context.Context, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	return getPluralMembersWithPersonalOrganization(ctx, a.query, memberIDs, order, np)
}

// GetPluralMembersWithPersonalOrganizationWithSd はSD付きで複数のメンバーを取得します。
func (a *PgAdapter) GetPluralMembersWithPersonalOrganizationWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	order parameter.MemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberWithPersonalOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberWithPersonalOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersWithPersonalOrganization(ctx, qtx, memberIDs, order, np)
}

// updateMember はメンバーを更新する内部関数です。
func updateMember(
	ctx context.Context, qtx *query.Queries,
	memberID uuid.UUID, param parameter.UpdateMemberParams, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberParams{
		MemberID:       memberID,
		Email:          param.Email,
		Name:           param.Name,
		FirstName:      param.FirstName,
		LastName:       param.LastName,
		ProfileImageID: pgtype.UUID(param.ProfileImageID),
		UpdatedAt:      now,
	}
	e, err := qtx.UpdateMember(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMember はメンバーを更新します。
func (a *PgAdapter) UpdateMember(
	ctx context.Context, memberID uuid.UUID, param parameter.UpdateMemberParams,
) (entity.Member, error) {
	return updateMember(ctx, a.query, memberID, param, a.clocker.Now())
}

// UpdateMemberWithSd はSD付きでメンバーを更新します。
func (a *PgAdapter) UpdateMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, param parameter.UpdateMemberParams,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMember(ctx, qtx, memberID, param, a.clocker.Now())
}

func updateMemberAttendStatus(
	ctx context.Context, qtx *query.Queries,
	memberID, attendStatusID uuid.UUID, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberAttendStatusParams{
		MemberID:       memberID,
		AttendStatusID: attendStatusID,
		UpdatedAt:      now,
	}
	e, err := qtx.UpdateMemberAttendStatus(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member attend status: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberAttendStatus はメンバーの出席状況を更新します。
func (a *PgAdapter) UpdateMemberAttendStatus(
	ctx context.Context, memberID, attendStatusID uuid.UUID,
) (entity.Member, error) {
	return updateMemberAttendStatus(ctx, a.query, memberID, attendStatusID, a.clocker.Now())
}

// UpdateMemberAttendStatusWithSd はSD付きでメンバーの出席状況を更新します。
func (a *PgAdapter) UpdateMemberAttendStatusWithSd(
	ctx context.Context, sd store.Sd, memberID, attendStatusID uuid.UUID,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberAttendStatus(ctx, qtx, memberID, attendStatusID, a.clocker.Now())
}

func updateMemberRole(
	ctx context.Context, qtx *query.Queries,
	memberID uuid.UUID, roleID entity.UUID, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberRoleParams{
		MemberID:  memberID,
		RoleID:    pgtype.UUID(roleID),
		UpdatedAt: now,
	}
	e, err := qtx.UpdateMemberRole(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member role: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberRole はメンバーのロールを更新します。
func (a *PgAdapter) UpdateMemberRole(
	ctx context.Context, memberID uuid.UUID, roleID entity.UUID,
) (entity.Member, error) {
	return updateMemberRole(ctx, a.query, memberID, roleID, a.clocker.Now())
}

// UpdateMemberRoleWithSd はSD付きでメンバーのロールを更新します。
func (a *PgAdapter) UpdateMemberRoleWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, roleID entity.UUID,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberRole(ctx, qtx, memberID, roleID, a.clocker.Now())
}

// updateMemberGroup はメンバーのグループを更新する内部関数です。
func updateMemberGroup(
	ctx context.Context, qtx *query.Queries,
	memberID, groupID uuid.UUID, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberGroupParams{
		MemberID:  memberID,
		GroupID:   groupID,
		UpdatedAt: now,
	}
	e, err := qtx.UpdateMemberGroup(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member group: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberGroup はメンバーのグループを更新します。
func (a *PgAdapter) UpdateMemberGroup(
	ctx context.Context, memberID, groupID uuid.UUID,
) (entity.Member, error) {
	return updateMemberGroup(ctx, a.query, memberID, groupID, a.clocker.Now())
}

// UpdateMemberGroupWithSd はSD付きでメンバーのグループを更新します。
func (a *PgAdapter) UpdateMemberGroupWithSd(
	ctx context.Context, sd store.Sd, memberID, groupID uuid.UUID,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberGroup(ctx, qtx, memberID, groupID, a.clocker.Now())
}

// updateMemberGrade はメンバーの学年を更新する内部関数です。
func updateMemberGrade(
	ctx context.Context, qtx *query.Queries,
	memberID, gradeID uuid.UUID, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberGradeParams{
		MemberID:  memberID,
		GradeID:   gradeID,
		UpdatedAt: now,
	}
	e, err := qtx.UpdateMemberGrade(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member grade: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberGrade はメンバーの学年を更新します。
func (a *PgAdapter) UpdateMemberGrade(
	ctx context.Context, memberID, gradeID uuid.UUID,
) (entity.Member, error) {
	return updateMemberGrade(ctx, a.query, memberID, gradeID, a.clocker.Now())
}

// UpdateMemberGradeWithSd はSD付きでメンバーの学年を更新します。
func (a *PgAdapter) UpdateMemberGradeWithSd(
	ctx context.Context, sd store.Sd, memberID, gradeID uuid.UUID,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberGrade(ctx, qtx, memberID, gradeID, a.clocker.Now())
}

func updateMemberLoginID(
	ctx context.Context, qtx *query.Queries,
	memberID uuid.UUID, loginID string, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberLoginIDParams{
		MemberID:  memberID,
		LoginID:   loginID,
		UpdatedAt: now,
	}
	e, err := qtx.UpdateMemberLoginID(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member login id: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberLoginID はメンバーのログインIDを更新します。
func (a *PgAdapter) UpdateMemberLoginID(
	ctx context.Context, memberID uuid.UUID, loginID string,
) (entity.Member, error) {
	return updateMemberLoginID(ctx, a.query, memberID, loginID, a.clocker.Now())
}

// UpdateMemberLoginIDWithSd はSD付きでメンバーのログインIDを更新します。
func (a *PgAdapter) UpdateMemberLoginIDWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, loginID string,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberLoginID(ctx, qtx, memberID, loginID, a.clocker.Now())
}

func updateMemberPassword(
	ctx context.Context, qtx *query.Queries,
	memberID uuid.UUID, password string, now time.Time,
) (entity.Member, error) {
	p := query.UpdateMemberPasswordParams{
		MemberID:  memberID,
		Password:  password,
		UpdatedAt: now,
	}
	e, err := qtx.UpdateMemberPassword(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Member{}, errhandle.NewModelNotFoundError("member")
		}
		return entity.Member{}, fmt.Errorf("failed to update member password: %w", err)
	}
	return convMember(query.Member(e)), nil
}

// UpdateMemberPassword はメンバーのパスワードを更新します。
func (a *PgAdapter) UpdateMemberPassword(
	ctx context.Context, memberID uuid.UUID, password string,
) (entity.Member, error) {
	return updateMemberPassword(ctx, a.query, memberID, password, a.clocker.Now())
}

// UpdateMemberPasswordWithSd はSD付きでメンバーのパスワードを更新します。
func (a *PgAdapter) UpdateMemberPasswordWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, password string,
) (entity.Member, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Member{}, store.ErrNotFoundDescriptor
	}
	return updateMemberPassword(ctx, qtx, memberID, password, a.clocker.Now())
}

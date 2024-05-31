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

func convOrganizationOnMember(r query.GetOrganizationsOnMemberRow) entity.OrganizationOnMember {
	return entity.OrganizationOnMember{
		Organization: entity.Organization{
			OrganizationID: r.OrganizationID,
			Name:           r.OrganizationName.String,
			Color:          entity.String(r.OrganizationColor),
			Description:    entity.String(r.OrganizationDescription),
			IsPersonal:     r.OrganizationIsPersonal.Bool,
			IsWhole:        r.OrganizationIsWhole.Bool,
			ChatRoomID:     entity.UUID(r.OrganizationChatRoomID),
		},
		WorkPositionID: entity.UUID(r.WorkPositionID),
		AddedAt:        r.AddedAt,
	}
}

func convMemberOnOrganization(r query.GetMembersOnOrganizationRow) entity.MemberOnOrganization {
	var profileImg entity.NullableEntity[entity.ImageWithAttachableItem]
	if r.MemberProfileImageID.Valid {
		profileImg = entity.NullableEntity[entity.ImageWithAttachableItem]{
			Valid: true,
			Entity: entity.ImageWithAttachableItem{
				ImageID: r.MemberProfileImageID.Bytes,
				Height:  entity.Float(r.MemberProfileImageHeight),
				Width:   entity.Float(r.MemberProfileImageWidth),
				AttachableItem: entity.AttachableItem{
					AttachableItemID: r.MemberProfileImageAttachableItemID.Bytes,
					OwnerID:          entity.UUID(r.MemberProfileImageOwnerID),
					FromOuter:        r.MemberProfileImageFromOuter.Bool,
					URL:              r.MemberProfileImageUrl.String,
					Alias:            r.MemberProfileImageAlias.String,
					Size:             entity.Float(r.MemberProfileImageSize),
					MimeTypeID:       r.MemberProfileImageMimeTypeID.Bytes,
				},
			},
		}
	}
	return entity.MemberOnOrganization{
		Member: entity.MemberCard{
			MemberID:     r.MemberID,
			Name:         r.MemberName.String,
			FirstName:    entity.String(r.MemberFirstName),
			LastName:     entity.String(r.MemberLastName),
			Email:        r.MemberEmail.String,
			ProfileImage: profileImg,
			GradeID:      r.MemberGradeID.Bytes,
			GroupID:      r.MemberGroupID.Bytes,
		},
		WorkPositionID: entity.UUID(r.WorkPositionID),
		AddedAt:        r.AddedAt,
	}
}

func countOrganizationsOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	where parameter.WhereOrganizationOnMemberParam,
) (int64, error) {
	p := query.CountOrganizationsOnMemberParams{
		MemberID:      memberID,
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountOrganizationsOnMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count memberships on member: %w", err)
	}
	return c, nil
}

// CountOrganizationsOnMember メンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountOrganizationsOnMember(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereOrganizationOnMemberParam,
) (int64, error) {
	return countOrganizationsOnMember(ctx, a.query, memberID, where)
}

// CountOrganizationsOnMemberWithSd SD付きでメンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountOrganizationsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, where parameter.WhereOrganizationOnMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countOrganizationsOnMember(ctx, qtx, memberID, where)
}

func countMembersOnOrganization(
	ctx context.Context,
	qtx *query.Queries,
	organizationID uuid.UUID,
	where parameter.WhereMemberOnOrganizationParam,
) (int64, error) {
	p := query.CountMembersOnOrganizationParams{
		OrganizationID: organizationID,
		WhereLikeName:  where.WhereLikeName,
		SearchName:     where.SearchName,
	}
	c, err := qtx.CountMembersOnOrganization(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count members on membership: %w", err)
	}
	return c, nil
}

// CountMembersOnOrganization チャットルーム上のメンバー数を取得する。
func (a *PgAdapter) CountMembersOnOrganization(
	ctx context.Context, organizationID uuid.UUID, where parameter.WhereMemberOnOrganizationParam,
) (int64, error) {
	return countMembersOnOrganization(ctx, a.query, organizationID, where)
}

// CountMembersOnOrganizationWithSd SD付きでチャットルーム上のメンバー数を取得する。
func (a *PgAdapter) CountMembersOnOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID, where parameter.WhereMemberOnOrganizationParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMembersOnOrganization(ctx, qtx, organizationID, where)
}

func belongOrganization(
	ctx context.Context,
	qtx *query.Queries,
	param parameter.BelongOrganizationParam,
) (entity.Membership, error) {
	p := query.CreateMembershipParams{
		MemberID:       param.MemberID,
		OrganizationID: param.OrganizationID,
		WorkPositionID: pgtype.UUID(param.WorkPositionID),
		AddedAt:        param.AddedAt,
	}
	b, err := qtx.CreateMembership(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Membership{}, errhandle.NewModelNotFoundError("membership")
		}
		return entity.Membership{}, fmt.Errorf("failed to belong membership: %w", err)
	}
	return entity.Membership{
		MemberID:       b.MemberID,
		OrganizationID: b.OrganizationID,
		WorkPositionID: entity.UUID(b.WorkPositionID),
		AddedAt:        b.AddedAt,
	}, nil
}

// BelongOrganization メンバーをチャットルームに所属させる。
func (a *PgAdapter) BelongOrganization(
	ctx context.Context, param parameter.BelongOrganizationParam,
) (entity.Membership, error) {
	return belongOrganization(ctx, a.query, param)
}

// BelongOrganizationWithSd SD付きでメンバーをチャットルームに所属させる。
func (a *PgAdapter) BelongOrganizationWithSd(
	ctx context.Context, sd store.Sd, param parameter.BelongOrganizationParam,
) (entity.Membership, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Membership{}, store.ErrNotFoundDescriptor
	}
	return belongOrganization(ctx, qtx, param)
}

func belongOrganizations(
	ctx context.Context,
	qtx *query.Queries,
	params []parameter.BelongOrganizationParam,
) (int64, error) {
	param := make([]query.CreateMembershipsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateMembershipsParams{
			MemberID:       p.MemberID,
			OrganizationID: p.OrganizationID,
			WorkPositionID: pgtype.UUID(p.WorkPositionID),
			AddedAt:        p.AddedAt,
		}
	}
	b, err := qtx.CreateMemberships(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelNotFoundError("membership")
		}
		return 0, fmt.Errorf("failed to belong memberships: %w", err)
	}
	return b, nil
}

// BelongOrganizations メンバーを複数のチャットルームに所属させる。
func (a *PgAdapter) BelongOrganizations(
	ctx context.Context, params []parameter.BelongOrganizationParam,
) (int64, error) {
	return belongOrganizations(ctx, a.query, params)
}

// BelongOrganizationsWithSd SD付きでメンバーを複数のチャットルームに所属させる。
func (a *PgAdapter) BelongOrganizationsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.BelongOrganizationParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return belongOrganizations(ctx, qtx, params)
}

func disbelongOrganization(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	organizationID uuid.UUID,
) (int64, error) {
	p := query.DeleteMembershipParams{
		MemberID:       memberID,
		OrganizationID: organizationID,
	}
	b, err := qtx.DeleteMembership(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong membership: %w", err)
	}
	if b != 1 {
		return 0, errhandle.NewModelNotFoundError("membership")
	}
	return b, nil
}

// DisbelongOrganization メンバーをチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganization(
	ctx context.Context, memberID, organizationID uuid.UUID,
) (int64, error) {
	return disbelongOrganization(ctx, a.query, memberID, organizationID)
}

// DisbelongOrganizationWithSd SD付きでメンバーをチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganizationWithSd(
	ctx context.Context, sd store.Sd, memberID, organizationID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongOrganization(ctx, qtx, memberID, organizationID)
}

func disbelongOrganizationOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteMembershipsOnMember(ctx, memberID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong membership on member: %w", err)
	}
	return b, nil
}

// DisbelongOrganizationOnMember メンバー上のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnMember(ctx context.Context, memberID uuid.UUID) (int64, error) {
	return disbelongOrganizationOnMember(ctx, a.query, memberID)
}

// DisbelongOrganizationOnMemberWithSd SD付きでメンバー上のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongOrganizationOnMember(ctx, qtx, memberID)
}

func disbelongOrganizationOnMembers(
	ctx context.Context,
	qtx *query.Queries,
	memberIDs []uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteMembershipsOnMembers(ctx, memberIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong membership on members: %w", err)
	}
	return b, nil
}

// DisbelongOrganizationOnMembers メンバー上の複数のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnMembers(
	ctx context.Context, memberIDs []uuid.UUID,
) (int64, error) {
	return disbelongOrganizationOnMembers(ctx, a.query, memberIDs)
}

// DisbelongOrganizationOnMembersWithSd SD付きでメンバー上の複数のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongOrganizationOnMembers(ctx, qtx, memberIDs)
}

func disbelongOrganizationOnOrganization(
	ctx context.Context,
	qtx *query.Queries,
	organizationID uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteMembershipsOnOrganization(ctx, organizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong membership on organization: %w", err)
	}
	return b, nil
}

// DisbelongOrganizationOnOrganization チャットルーム上のメンバーから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnOrganization(
	ctx context.Context, organizationID uuid.UUID,
) (int64, error) {
	return disbelongOrganizationOnOrganization(ctx, a.query, organizationID)
}

// DisbelongOrganizationOnOrganizationWithSd SD付きでチャットルーム上のメンバーから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongOrganizationOnOrganization(ctx, qtx, organizationID)
}

func disbelongOrganizationOnOrganizations(
	ctx context.Context,
	qtx *query.Queries,
	organizationIDs []uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteMembershipsOnOrganizations(ctx, organizationIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong membership on organizations: %w", err)
	}
	return b, nil
}

// DisbelongOrganizationOnOrganizations チャットルーム上の複数のメンバーから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnOrganizations(
	ctx context.Context, organizationIDs []uuid.UUID,
) (int64, error) {
	return disbelongOrganizationOnOrganizations(ctx, a.query, organizationIDs)
}

// DisbelongOrganizationOnOrganizationsWithSd SD付きでチャットルーム上の複数のメンバーから所属解除する。
func (a *PgAdapter) DisbelongOrganizationOnOrganizationsWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongOrganizationOnOrganizations(ctx, qtx, organizationIDs)
}

func getOrganizationsOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	where parameter.WhereOrganizationOnMemberParam,
	order parameter.OrganizationOnMemberOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.OrganizationOnMember], error) {
	eConvFunc := func(e entity.OrganizationOnMemberForQuery) (entity.OrganizationOnMember, error) {
		return e.OrganizationOnMember, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountOrganizationsOnMemberParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountOrganizationsOnMember(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count memberships on member: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.OrganizationOnMemberForQuery, error) {
		p := query.GetOrganizationsOnMemberParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetOrganizationsOnMember(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.OrganizationOnMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get memberships on member: %w", err)
		}
		fq := make([]entity.OrganizationOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.OrganizationOnMemberForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				OrganizationOnMember: convOrganizationOnMember(e),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.OrganizationOnMemberForQuery, error) {
		var addCursor time.Time
		var nameCursor string
		var ok bool
		var err error
		switch subCursor {
		case parameter.OrganizationOnMemberNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		case parameter.OrganizationOnMemberAddedAtCursorKey:
			cv, ok := subCursorValue.(string)
			addCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				addCursor = time.Time{}
			}
		}
		p := query.GetOrganizationsOnMemberUseKeysetPaginateParams{
			MemberID:        memberID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
			NameCursor:      nameCursor,
			AddCursor:       addCursor,
		}
		r, err := qtx.GetOrganizationsOnMemberUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get memberships on member: %w", err)
		}
		fq := make([]entity.OrganizationOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.OrganizationOnMemberForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				OrganizationOnMember: convOrganizationOnMember(query.GetOrganizationsOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.OrganizationOnMemberForQuery, error) {
		p := query.GetOrganizationsOnMemberUseNumberedPaginateParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Limit:         limit,
			Offset:        offset,
		}
		r, err := qtx.GetOrganizationsOnMemberUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get memberships on member: %w", err)
		}
		fq := make([]entity.OrganizationOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.OrganizationOnMemberForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				OrganizationOnMember: convOrganizationOnMember(query.GetOrganizationsOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.OrganizationOnMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.OrganizationOnMemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.OrganizationOnMemberNameCursorKey:
			return entity.Int(e.Pkey), e.OrganizationOnMember.Organization.Name
		case parameter.OrganizationOnMemberAddedAtCursorKey:
			return entity.Int(e.Pkey), e.OrganizationOnMember.AddedAt
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
		return store.ListResult[entity.OrganizationOnMember]{}, fmt.Errorf("failed to get memberships on member: %w", err)
	}
	return res, nil
}

// GetOrganizationsOnMember メンバー上のチャットルームを取得する。
func (a *PgAdapter) GetOrganizationsOnMember(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereOrganizationOnMemberParam,
	order parameter.OrganizationOnMemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationOnMember], error) {
	return getOrganizationsOnMember(ctx, a.query, memberID, where, order, np, cp, wc)
}

// GetOrganizationsOnMemberWithSd SD付きでメンバー上のチャットルームを取得する。
func (a *PgAdapter) GetOrganizationsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
	where parameter.WhereOrganizationOnMemberParam, order parameter.OrganizationOnMemberOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getOrganizationsOnMember(ctx, qtx, memberID, where, order, np, cp, wc)
}

func getMembersOnOrganization(
	ctx context.Context,
	qtx *query.Queries,
	organizationID uuid.UUID,
	where parameter.WhereMemberOnOrganizationParam,
	order parameter.MemberOnOrganizationOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnOrganization], error) {
	eConvFunc := func(e entity.MemberOnOrganizationForQuery) (entity.MemberOnOrganization, error) {
		return e.MemberOnOrganization, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountMembersOnOrganizationParams{
			OrganizationID: organizationID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
		}
		r, err := qtx.CountMembersOnOrganization(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count members on membership: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberOnOrganizationForQuery, error) {
		p := query.GetMembersOnOrganizationParams{
			OrganizationID: organizationID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
			OrderMethod:    orderMethod,
		}
		r, err := qtx.GetMembersOnOrganization(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberOnOrganizationForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members on membership: %w", err)
		}
		fq := make([]entity.MemberOnOrganizationForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnOrganizationForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				MemberOnOrganization: convMemberOnOrganization(e),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberOnOrganizationForQuery, error) {
		var addCursor time.Time
		var nameCursor string
		var ok bool
		var err error
		switch subCursor {
		case parameter.MemberOnOrganizationNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		case parameter.MemberOnOrganizationAddedAtCursorKey:
			cv, ok := subCursorValue.(string)
			addCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				addCursor = time.Time{}
			}
		}
		p := query.GetMembersOnOrganizationUseKeysetPaginateParams{
			OrganizationID:  organizationID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
			AddedAtCursor:   addCursor,
		}
		r, err := qtx.GetMembersOnOrganizationUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on membership: %w", err)
		}
		fq := make([]entity.MemberOnOrganizationForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnOrganizationForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				MemberOnOrganization: convMemberOnOrganization(query.GetMembersOnOrganizationRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberOnOrganizationForQuery, error) {
		p := query.GetMembersOnOrganizationUseNumberedPaginateParams{
			OrganizationID: organizationID,
			WhereLikeName:  where.WhereLikeName,
			SearchName:     where.SearchName,
			OrderMethod:    orderMethod,
			Offset:         offset,
			Limit:          limit,
		}
		r, err := qtx.GetMembersOnOrganizationUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on membership: %w", err)
		}
		fq := make([]entity.MemberOnOrganizationForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnOrganizationForQuery{
				Pkey:                 entity.Int(e.MMembershipsPkey),
				MemberOnOrganization: convMemberOnOrganization(query.GetMembersOnOrganizationRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.MemberOnOrganizationForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberOnOrganizationDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberOnOrganizationNameCursorKey:
			return entity.Int(e.Pkey), e.MemberOnOrganization.Member.Name
		case parameter.MemberOnOrganizationAddedAtCursorKey:
			return entity.Int(e.Pkey), e.MemberOnOrganization.AddedAt
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
		return store.ListResult[entity.MemberOnOrganization]{}, fmt.Errorf("failed to get members on membership: %w", err)
	}
	return res, nil
}

// GetMembersOnOrganization チャットルーム上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnOrganization(
	ctx context.Context, organizationID uuid.UUID, where parameter.WhereMemberOnOrganizationParam,
	order parameter.MemberOnOrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberOnOrganization], error) {
	return getMembersOnOrganization(ctx, a.query, organizationID, where, order, np, cp, wc)
}

// GetMembersOnOrganizationWithSd SD付きでチャットルーム上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
	where parameter.WhereMemberOnOrganizationParam, order parameter.MemberOnOrganizationOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberOnOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getMembersOnOrganization(ctx, qtx, organizationID, where, order, np, cp, wc)
}

func getPluralOrganizationsOnMember(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.OrganizationOnMemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationOnMember], error) {
	var e []query.GetPluralOrganizationsOnMemberRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralOrganizationsOnMember(ctx, query.GetPluralOrganizationsOnMemberParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralOrganizationsOnMemberUseNumberedPaginateRow
		qe, err = qtx.GetPluralOrganizationsOnMemberUseNumberedPaginate(
			ctx, query.GetPluralOrganizationsOnMemberUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Limit:       int32(np.Limit.Int64),
				Offset:      int32(np.Offset.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralOrganizationsOnMemberRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralOrganizationsOnMemberRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.OrganizationOnMember]{},
			fmt.Errorf("failed to get memberships on member: %w", err)
	}
	entities := make([]entity.OrganizationOnMember, len(e))
	for i, v := range e {
		entities[i] = convOrganizationOnMember(query.GetOrganizationsOnMemberRow(v))
	}
	return store.ListResult[entity.OrganizationOnMember]{Data: entities}, nil
}

// GetPluralOrganizationsOnMember メンバー上の複数のチャットルームを取得する。
func (a *PgAdapter) GetPluralOrganizationsOnMember(
	ctx context.Context, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.OrganizationOnMemberOrderMethod,
) (store.ListResult[entity.OrganizationOnMember], error) {
	return getPluralOrganizationsOnMember(ctx, a.query, memberIDs, order, np)
}

// GetPluralOrganizationsOnMemberWithSd SD付きでメンバー上の複数のチャットルームを取得する。
func (a *PgAdapter) GetPluralOrganizationsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.OrganizationOnMemberOrderMethod,
) (store.ListResult[entity.OrganizationOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralOrganizationsOnMember(ctx, qtx, memberIDs, order, np)
}

func getPluralMembersOnOrganization(
	ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnOrganizationOrderMethod,
) (store.ListResult[entity.MemberOnOrganization], error) {
	var e []query.GetPluralMembersOnOrganizationRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersOnOrganization(ctx, query.GetPluralMembersOnOrganizationParams{
			OrganizationIds: organizationIDs,
			OrderMethod:     order.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralMembersOnOrganizationUseNumberedPaginateRow
		qe, err = qtx.GetPluralMembersOnOrganizationUseNumberedPaginate(
			ctx, query.GetPluralMembersOnOrganizationUseNumberedPaginateParams{
				OrganizationIds: organizationIDs,
				Limit:           int32(np.Limit.Int64),
				Offset:          int32(np.Offset.Int64),
				OrderMethod:     order.GetStringValue(),
			})
		e = make([]query.GetPluralMembersOnOrganizationRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralMembersOnOrganizationRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberOnOrganization]{},
			fmt.Errorf("failed to get members on membership: %w", err)
	}
	entities := make([]entity.MemberOnOrganization, len(e))
	for i, v := range e {
		entities[i] = convMemberOnOrganization(query.GetMembersOnOrganizationRow(v))
	}
	return store.ListResult[entity.MemberOnOrganization]{Data: entities}, nil
}

// GetPluralMembersOnOrganization チャットルーム上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnOrganization(
	ctx context.Context, organizationIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnOrganizationOrderMethod,
) (store.ListResult[entity.MemberOnOrganization], error) {
	return getPluralMembersOnOrganization(ctx, a.query, organizationIDs, np, order)
}

// GetPluralMembersOnOrganizationWithSd SD付きでチャットルーム上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnOrganizationOrderMethod,
) (store.ListResult[entity.MemberOnOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersOnOrganization(ctx, qtx, organizationIDs, np, order)
}

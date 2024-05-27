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

func convOrganizationWithChatRoom(e query.FindOrganizationByIDWithChatRoomRow) entity.OrganizationWithChatRoom {
	return entity.OrganizationWithChatRoom{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoom: entity.ChatRoomWithCoverImage{
			ChatRoomID:       e.ChatRoomID.Bytes,
			Name:             e.ChatRoomName.String,
			IsPrivate:        e.ChatRoomIsPrivate.Bool,
			FromOrganization: e.ChatRoomFromOrganization.Bool,
			CoverImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
				Valid: e.ChatRoomCoverImageID.Valid,
				Entity: entity.ImageWithAttachableItem{
					ImageID: e.ChatRoomCoverImageID.Bytes,
					Height:  entity.Float(e.ChatRoomCoverImageHeight),
					Width:   entity.Float(e.ChatRoomCoverImageWidth),
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.ChatRoomCoverImageAttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.ChatRoomCoverImageOwnerID),
						FromOuter:        e.ChatRoomCoverImageFromOuter.Bool,
						URL:              e.ChatRoomCoverImageUrl.String,
						Alias:            e.ChatRoomCoverImageAlias.String,
						Size:             entity.Float(e.ChatRoomCoverImageSize),
						MimeTypeID:       e.ChatRoomCoverImageMimeTypeID.Bytes,
					},
				},
			},
			OwnerID: entity.UUID(e.ChatRoomOwnerID),
		},
	}
}

func convOrganizationWithDetail(e query.FindOrganizationByIDWithDetailRow) entity.OrganizationWithDetail {
	return entity.OrganizationWithDetail{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
		Group: entity.NullableEntity[entity.Group]{
			Valid: e.GroupID.Valid,
			Entity: entity.Group{
				GroupID:        e.GroupID.Bytes,
				Key:            e.GroupKey.String,
				OrganizationID: e.OrganizationID,
			},
		},
		Grade: entity.NullableEntity[entity.Grade]{
			Valid: e.GradeID.Valid,
			Entity: entity.Grade{
				GradeID:        e.GradeID.Bytes,
				Key:            e.GradeKey.String,
				OrganizationID: e.OrganizationID,
			},
		},
	}
}

func convOrganizationWithChatRoomAndDetail(
	e query.FindOrganizationByIDWithChatRoomAndDetailRow,
) entity.OrganizationWithChatRoomAndDetail {
	return entity.OrganizationWithChatRoomAndDetail{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoom: entity.ChatRoomWithCoverImage{
			ChatRoomID:       e.ChatRoomID.Bytes,
			Name:             e.ChatRoomName.String,
			IsPrivate:        e.ChatRoomIsPrivate.Bool,
			FromOrganization: e.ChatRoomFromOrganization.Bool,
			CoverImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
				Valid: e.ChatRoomCoverImageID.Valid,
				Entity: entity.ImageWithAttachableItem{
					ImageID: e.ChatRoomCoverImageID.Bytes,
					Height:  entity.Float(e.ChatRoomCoverImageHeight),
					Width:   entity.Float(e.ChatRoomCoverImageWidth),
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.ChatRoomCoverImageAttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.ChatRoomCoverImageOwnerID),
						FromOuter:        e.ChatRoomCoverImageFromOuter.Bool,
						URL:              e.ChatRoomCoverImageUrl.String,
						Alias:            e.ChatRoomCoverImageAlias.String,
						Size:             entity.Float(e.ChatRoomCoverImageSize),
						MimeTypeID:       e.ChatRoomCoverImageMimeTypeID.Bytes,
					},
				},
			},
			OwnerID: entity.UUID(e.ChatRoomOwnerID),
		},
		Group: entity.NullableEntity[entity.Group]{
			Valid: e.GroupID.Valid,
			Entity: entity.Group{
				GroupID:        e.GroupID.Bytes,
				Key:            e.GroupKey.String,
				OrganizationID: e.OrganizationID,
			},
		},
		Grade: entity.NullableEntity[entity.Grade]{
			Valid: e.GradeID.Valid,
			Entity: entity.Grade{
				GradeID:        e.GradeID.Bytes,
				Key:            e.GradeKey.String,
				OrganizationID: e.OrganizationID,
			},
		},
	}
}

// countOrganizations はオーガナイゼーション数を取得する内部関数です。
func countOrganizations(
	ctx context.Context, qtx *query.Queries, where parameter.WhereOrganizationParam,
) (int64, error) {
	p := query.CountOrganizationsParams{
		WhereLikeName:    where.WhereLikeName,
		SearchName:       where.SearchName,
		WhereIsWhole:     where.WhereIsWhole,
		IsWhole:          where.IsWhole,
		WhereIsPersonal:  where.WhereIsPersonal,
		IsPersonal:       where.IsPersonal,
		PersonalMemberID: where.PersonalMemberID,
		WhereIsGroup:     where.WhereIsGroup,
		WhereIsGrade:     where.WhereIsGrade,
	}
	c, err := qtx.CountOrganizations(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count organizations: %w", err)
	}
	return c, nil
}

// CountOrganizations はオーガナイゼーション数を取得します。
func (a *PgAdapter) CountOrganizations(ctx context.Context, where parameter.WhereOrganizationParam) (int64, error) {
	return countOrganizations(ctx, a.query, where)
}

// CountOrganizationsWithSd はSD付きでオーガナイゼーション数を取得します。
func (a *PgAdapter) CountOrganizationsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereOrganizationParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countOrganizations(ctx, qtx, where)
}

// createOrganization はオーガナイゼーションを作成する内部関数です。
func createOrganization(
	ctx context.Context, qtx *query.Queries, param parameter.CreateOrganizationParam, now time.Time,
) (entity.Organization, error) {
	p := query.CreateOrganizationParams{
		Name:        param.Name,
		Description: pgtype.Text(param.Description),
		Color:       pgtype.Text(param.Color),
		IsPersonal:  param.IsPersonal,
		IsWhole:     param.IsWhole,
		ChatRoomID:  pgtype.UUID(param.ChatRoomID),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	e, err := qtx.CreateOrganization(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Organization{}, errhandle.NewModelDuplicatedError("organization")
		}
		return entity.Organization{}, fmt.Errorf("failed to create organization: %w", err)
	}
	entity := entity.Organization{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
	}
	return entity, nil
}

// CreateOrganization はオーガナイゼーションを作成します。
func (a *PgAdapter) CreateOrganization(
	ctx context.Context, param parameter.CreateOrganizationParam,
) (entity.Organization, error) {
	return createOrganization(ctx, a.query, param, a.clocker.Now())
}

// CreateOrganizationWithSd はSD付きでオーガナイゼーションを作成します。
func (a *PgAdapter) CreateOrganizationWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateOrganizationParam,
) (entity.Organization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Organization{}, store.ErrNotFoundDescriptor
	}
	return createOrganization(ctx, qtx, param, a.clocker.Now())
}

// createOrganizations は複数のオーガナイゼーションを作成する内部関数です。
func createOrganizations(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateOrganizationParam, now time.Time,
) (int64, error) {
	param := make([]query.CreateOrganizationsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateOrganizationsParams{
			Name:        p.Name,
			Description: pgtype.Text(p.Description),
			Color:       pgtype.Text(p.Color),
			IsPersonal:  p.IsPersonal,
			IsWhole:     p.IsWhole,
			ChatRoomID:  pgtype.UUID(p.ChatRoomID),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
	n, err := qtx.CreateOrganizations(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("organization")
		}
		return 0, fmt.Errorf("failed to create organizations: %w", err)
	}
	return n, nil
}

// CreateOrganizations は複数のオーガナイゼーションを作成します。
func (a *PgAdapter) CreateOrganizations(
	ctx context.Context, params []parameter.CreateOrganizationParam,
) (int64, error) {
	return createOrganizations(ctx, a.query, params, a.clocker.Now())
}

// CreateOrganizationsWithSd はSD付きで複数のオーガナイゼーションを作成します。
func (a *PgAdapter) CreateOrganizationsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateOrganizationParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createOrganizations(ctx, qtx, params, a.clocker.Now())
}

// deleteOrganization はオーガナイゼーションを削除する内部関数です。
func deleteOrganization(ctx context.Context, qtx *query.Queries, organizationID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteOrganization(ctx, organizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("organization")
	}
	return c, nil
}

// DeleteOrganization はオーガナイゼーションを削除します。
func (a *PgAdapter) DeleteOrganization(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	return deleteOrganization(ctx, a.query, organizationID)
}

// DeleteOrganizationWithSd はSD付きでオーガナイゼーションを削除します。
func (a *PgAdapter) DeleteOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteOrganization(ctx, qtx, organizationID)
}

// pluralDeleteOrganizations は複数のオーガナイゼーションを削除する内部関数です。
func pluralDeleteOrganizations(ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteOrganizations(ctx, organizationIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organizations: %w", err)
	}
	if c != int64(len(organizationIDs)) {
		return 0, errhandle.NewModelNotFoundError("organization")
	}
	return c, nil
}

// PluralDeleteOrganizations は複数のオーガナイゼーションを削除します。
func (a *PgAdapter) PluralDeleteOrganizations(ctx context.Context, organizationIDs []uuid.UUID) (int64, error) {
	return pluralDeleteOrganizations(ctx, a.query, organizationIDs)
}

// PluralDeleteOrganizationsWithSd はSD付きで複数のオーガナイゼーションを削除します。
func (a *PgAdapter) PluralDeleteOrganizationsWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteOrganizations(ctx, qtx, organizationIDs)
}

// findOrganizationByID はオーガナイゼーションをIDで取得する内部関数です。
func findOrganizationByID(
	ctx context.Context, qtx *query.Queries, organizationID uuid.UUID,
) (entity.Organization, error) {
	e, err := qtx.FindOrganizationByID(ctx, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Organization{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.Organization{}, fmt.Errorf("failed to find organization: %w", err)
	}
	entity := entity.Organization{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
	}
	return entity, nil
}

// FindOrganizationByID はオーガナイゼーションをIDで取得します。
func (a *PgAdapter) FindOrganizationByID(ctx context.Context, organizationID uuid.UUID) (entity.Organization, error) {
	return findOrganizationByID(ctx, a.query, organizationID)
}

// FindOrganizationByIDWithSd はSD付きでオーガナイゼーションをIDで取得します。
func (a *PgAdapter) FindOrganizationByIDWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (entity.Organization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Organization{}, store.ErrNotFoundDescriptor
	}
	return findOrganizationByID(ctx, qtx, organizationID)
}

// findOrganizationWithChatRoom はオーガナイゼーションとチャットルームを取得する内部関数です。
func findOrganizationWithChatRoom(
	ctx context.Context, qtx *query.Queries, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoom, error) {
	e, err := qtx.FindOrganizationByIDWithChatRoom(ctx, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.OrganizationWithChatRoom{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.OrganizationWithChatRoom{}, fmt.Errorf("failed to find organization with chat room: %w", err)
	}
	return convOrganizationWithChatRoom(e), nil
}

// FindOrganizationWithChatRoom はオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) FindOrganizationWithChatRoom(
	ctx context.Context, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoom, error) {
	return findOrganizationWithChatRoom(ctx, a.query, organizationID)
}

// FindOrganizationWithChatRoomWithSd はSD付きでオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) FindOrganizationWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.OrganizationWithChatRoom{}, store.ErrNotFoundDescriptor
	}
	return findOrganizationWithChatRoom(ctx, qtx, organizationID)
}

func findOrganizationWithDetail(
	ctx context.Context, qtx *query.Queries, organizationID uuid.UUID,
) (entity.OrganizationWithDetail, error) {
	e, err := qtx.FindOrganizationByIDWithDetail(ctx, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.OrganizationWithDetail{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.OrganizationWithDetail{}, fmt.Errorf("failed to find organization with detail: %w", err)
	}
	return convOrganizationWithDetail(e), nil
}

// FindOrganizationWithDetail はオーガナイゼーションと詳細を取得します。
func (a *PgAdapter) FindOrganizationWithDetail(
	ctx context.Context, organizationID uuid.UUID,
) (entity.OrganizationWithDetail, error) {
	return findOrganizationWithDetail(ctx, a.query, organizationID)
}

// FindOrganizationWithDetailWithSd はSD付きでオーガナイゼーションと詳細を取得します。
func (a *PgAdapter) FindOrganizationWithDetailWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (entity.OrganizationWithDetail, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.OrganizationWithDetail{}, store.ErrNotFoundDescriptor
	}
	return findOrganizationWithDetail(ctx, qtx, organizationID)
}

func findOrganizationWithChatRoomAndDetail(
	ctx context.Context, qtx *query.Queries, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoomAndDetail, error) {
	e, err := qtx.FindOrganizationByIDWithChatRoomAndDetail(ctx, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.OrganizationWithChatRoomAndDetail{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.OrganizationWithChatRoomAndDetail{},
			fmt.Errorf("failed to find organization with chat room and detail: %w", err)
	}
	return convOrganizationWithChatRoomAndDetail(e), nil
}

// FindOrganizationWithChatRoomAndDetail はオーガナイゼーションとチャットルーム、詳細を取得します。
func (a *PgAdapter) FindOrganizationWithChatRoomAndDetail(
	ctx context.Context, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoomAndDetail, error) {
	return findOrganizationWithChatRoomAndDetail(ctx, a.query, organizationID)
}

// FindOrganizationWithChatRoomAndDetailWithSd はSD付きでオーガナイゼーションとチャットルーム、詳細を取得します。
func (a *PgAdapter) FindOrganizationWithChatRoomAndDetailWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID,
) (entity.OrganizationWithChatRoomAndDetail, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.OrganizationWithChatRoomAndDetail{}, store.ErrNotFoundDescriptor
	}
	return findOrganizationWithChatRoomAndDetail(ctx, qtx, organizationID)
}

// getOrganizations はオーガナイゼーションを取得する内部関数です。
func getOrganizations(
	ctx context.Context,
	qtx *query.Queries,
	where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Organization], error) {
	eConvFunc := func(e query.Organization) (entity.Organization, error) {
		return entity.Organization{
			OrganizationID: e.OrganizationID,
			Name:           e.Name,
			Description:    entity.String(e.Description),
			Color:          entity.String(e.Color),
			IsPersonal:     e.IsPersonal,
			IsWhole:        e.IsWhole,
			ChatRoomID:     entity.UUID(e.ChatRoomID),
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountOrganizationsParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
		}
		r, err := qtx.CountOrganizations(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count organizations: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Organization, error) {
		p := query.GetOrganizationsParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizations(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Organization{}, nil
			}
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Organization, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.OrganizationNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetOrganizationsUseKeysetPaginateParams{
			Limit:            limit,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			CursorDirection:  cursorDir,
			Cursor:           cursor,
			NameCursor:       nameCursor,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Organization, error) {
		p := query.GetOrganizationsUseNumberedPaginateParams{
			Limit:            limit,
			Offset:           offset,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Organization) (entity.Int, any) {
		switch subCursor {
		case parameter.OrganizationDefaultCursorKey:
			return entity.Int(e.MOrganizationsPkey), nil
		case parameter.OrganizationNameCursorKey:
			return entity.Int(e.MOrganizationsPkey), e.Name
		}
		return entity.Int(e.MOrganizationsPkey), nil
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
		return store.ListResult[entity.Organization]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	return res, nil
}

// GetOrganizations はオーガナイゼーションを取得します。
func (a *PgAdapter) GetOrganizations(
	ctx context.Context, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Organization], error) {
	return getOrganizations(ctx, a.query, where, order, np, cp, wc)
}

// GetOrganizationsWithSd はSD付きでオーガナイゼーションを取得します。
func (a *PgAdapter) GetOrganizationsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Organization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Organization]{}, store.ErrNotFoundDescriptor
	}
	return getOrganizations(ctx, qtx, where, order, np, cp, wc)
}

func findWholeOrganization(ctx context.Context, qtx *query.Queries) (entity.Organization, error) {
	e, err := qtx.FindWholeOrganization(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Organization{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.Organization{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	entity := entity.Organization{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
	}
	return entity, nil
}

// FindWholeOrganization は全体オーガナイゼーションを取得します。
func (a *PgAdapter) FindWholeOrganization(ctx context.Context) (entity.Organization, error) {
	return findWholeOrganization(ctx, a.query)
}

// FindWholeOrganizationWithSd はSD付きで全体オーガナイゼーションを取得します。
func (a *PgAdapter) FindWholeOrganizationWithSd(ctx context.Context, sd store.Sd) (entity.Organization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Organization{}, store.ErrNotFoundDescriptor
	}
	return findWholeOrganization(ctx, qtx)
}

func findPersonalOrganization(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (entity.Organization, error) {
	e, err := qtx.FindPersonalOrganization(ctx, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Organization{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.Organization{}, fmt.Errorf("failed to find personal organization: %w", err)
	}
	entity := entity.Organization{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
	}
	return entity, nil
}

// FindPersonalOrganization は個人オーガナイゼーションを取得します。
func (a *PgAdapter) FindPersonalOrganization(ctx context.Context, memberID uuid.UUID) (entity.Organization, error) {
	return findPersonalOrganization(ctx, a.query, memberID)
}

// FindPersonalOrganizationWithSd はSD付きで個人オーガナイゼーションを取得します。
func (a *PgAdapter) FindPersonalOrganizationWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (entity.Organization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Organization{}, store.ErrNotFoundDescriptor
	}
	return findPersonalOrganization(ctx, qtx, memberID)
}

// getPluralOrganizations は複数のオーガナイゼーションを取得する内部関数です。
func getPluralOrganizations(
	ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID,
	orderMethod parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Organization], error) {
	var e []query.Organization
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralOrganizations(ctx, query.GetPluralOrganizationsParams{
			OrganizationIds: organizationIDs,
			OrderMethod:     orderMethod.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralOrganizationsUseNumberedPaginate(ctx, query.GetPluralOrganizationsUseNumberedPaginateParams{
			OrganizationIds: organizationIDs,
			Offset:          int32(np.Offset.Int64),
			Limit:           int32(np.Limit.Int64),
			OrderMethod:     orderMethod.GetStringValue(),
		})
	}
	if err != nil {
		return store.ListResult[entity.Organization]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	entities := make([]entity.Organization, len(e))
	for i, v := range e {
		entities[i] = entity.Organization{
			OrganizationID: v.OrganizationID,
			Name:           v.Name,
			Description:    entity.String(v.Description),
			Color:          entity.String(v.Color),
			IsPersonal:     v.IsPersonal,
			IsWhole:        v.IsWhole,
			ChatRoomID:     entity.UUID(v.ChatRoomID),
		}
	}
	return store.ListResult[entity.Organization]{Data: entities}, nil
}

// GetPluralOrganizations は複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizations(
	ctx context.Context, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Organization], error) {
	return getPluralOrganizations(ctx, a.query, organizationIDs, order, np)
}

// GetPluralOrganizationsWithSd はSD付きで複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Organization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Organization]{}, store.ErrNotFoundDescriptor
	}
	return getPluralOrganizations(ctx, qtx, organizationIDs, order, np)
}

func getOrganizationsWithChatRoom(
	ctx context.Context, qtx *query.Queries, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	eConvFunc := func(e entity.OrganizationWithChatRoomForQuery) (entity.OrganizationWithChatRoom, error) {
		return e.OrganizationWithChatRoom, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountOrganizationsParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
		}
		r, err := qtx.CountOrganizations(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count organizations: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.OrganizationWithChatRoomForQuery, error) {
		p := query.GetOrganizationsWithChatRoomParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoom(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.OrganizationWithChatRoomForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomForQuery{
				Pkey:                     entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoom: convOrganizationWithChatRoom(query.FindOrganizationByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.OrganizationWithChatRoomForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.OrganizationNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetOrganizationsWithChatRoomUseKeysetPaginateParams{
			Limit:            limit,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			CursorDirection:  cursorDir,
			Cursor:           cursor,
			NameCursor:       nameCursor,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomForQuery{
				Pkey:                     entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoom: convOrganizationWithChatRoom(query.FindOrganizationByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.OrganizationWithChatRoomForQuery, error) {
		p := query.GetOrganizationsWithChatRoomUseNumberedPaginateParams{
			Limit:            limit,
			Offset:           offset,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomForQuery{
				Pkey:                     entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoom: convOrganizationWithChatRoom(query.FindOrganizationByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.OrganizationWithChatRoomForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.OrganizationDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.OrganizationNameCursorKey:
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
		return store.ListResult[entity.OrganizationWithChatRoom]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	return res, nil
}

// GetOrganizationsWithChatRoom はオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) GetOrganizationsWithChatRoom(
	ctx context.Context, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	return getOrganizationsWithChatRoom(ctx, a.query, where, order, np, cp, wc)
}

// GetOrganizationsWithChatRoomWithSd はSD付きでオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) GetOrganizationsWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getOrganizationsWithChatRoom(ctx, qtx, where, order, np, cp, wc)
}

// getPluralOrganizationsWithChatRoom は複数のオーガナイゼーションを取得する内部関数です。
func getPluralOrganizationsWithChatRoom(
	ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID,
	orderMethod parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	var e []query.GetPluralOrganizationsWithChatRoomRow
	var te []query.GetPluralOrganizationsWithChatRoomUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralOrganizationsWithChatRoom(ctx, query.GetPluralOrganizationsWithChatRoomParams{
			OrganizationIds: organizationIDs,
			OrderMethod:     orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralOrganizationsWithChatRoomUseNumberedPaginate(
			ctx, query.GetPluralOrganizationsWithChatRoomUseNumberedPaginateParams{
				OrganizationIds: organizationIDs,
				Offset:          int32(np.Offset.Int64),
				Limit:           int32(np.Limit.Int64),
				OrderMethod:     orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralOrganizationsWithChatRoomRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralOrganizationsWithChatRoomRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.OrganizationWithChatRoom]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	entities := make([]entity.OrganizationWithChatRoom, len(e))
	for i, v := range e {
		entities[i] = convOrganizationWithChatRoom(query.FindOrganizationByIDWithChatRoomRow(v))
	}
	return store.ListResult[entity.OrganizationWithChatRoom]{Data: entities}, nil
}

// GetPluralOrganizationsWithChatRoom は複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithChatRoom(
	ctx context.Context, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	return getPluralOrganizationsWithChatRoom(ctx, a.query, organizationIDs, order, np)
}

// GetPluralOrganizationsWithChatRoomWithSd はSD付きで複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getPluralOrganizationsWithChatRoom(ctx, qtx, organizationIDs, order, np)
}

func getOrganizationsWithDetail(
	ctx context.Context, qtx *query.Queries, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	eConvFunc := func(e entity.OrganizationWithDetailForQuery) (entity.OrganizationWithDetail, error) {
		return e.OrganizationWithDetail, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountOrganizationsParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
		}
		r, err := qtx.CountOrganizations(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count organizations: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.OrganizationWithDetailForQuery, error) {
		p := query.GetOrganizationsWithDetailParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithDetail(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.OrganizationWithDetailForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithDetailForQuery{
				Pkey:                   entity.Int(v.MOrganizationsPkey),
				OrganizationWithDetail: convOrganizationWithDetail(query.FindOrganizationByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.OrganizationWithDetailForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.OrganizationNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetOrganizationsWithDetailUseKeysetPaginateParams{
			Limit:            limit,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			CursorDirection:  cursorDir,
			Cursor:           cursor,
			NameCursor:       nameCursor,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithDetailUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithDetailForQuery{
				Pkey:                   entity.Int(v.MOrganizationsPkey),
				OrganizationWithDetail: convOrganizationWithDetail(query.FindOrganizationByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.OrganizationWithDetailForQuery, error) {
		p := query.GetOrganizationsWithDetailUseNumberedPaginateParams{
			Limit:            limit,
			Offset:           offset,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithDetailUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithDetailForQuery{
				Pkey:                   entity.Int(v.MOrganizationsPkey),
				OrganizationWithDetail: convOrganizationWithDetail(query.FindOrganizationByIDWithDetailRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.OrganizationWithDetailForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.OrganizationDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.OrganizationNameCursorKey:
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
		return store.ListResult[entity.OrganizationWithDetail]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	return res, nil
}

// GetOrganizationsWithDetail はオーガナイゼーションと詳細を取得します。
func (a *PgAdapter) GetOrganizationsWithDetail(
	ctx context.Context, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	return getOrganizationsWithDetail(ctx, a.query, where, order, np, cp, wc)
}

// GetOrganizationsWithDetailWithSd はSD付きでオーガナイゼーションと詳細を取得します。
func (a *PgAdapter) GetOrganizationsWithDetailWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getOrganizationsWithDetail(ctx, qtx, where, order, np, cp, wc)
}

// getPluralOrganizationsWithDetail は複数のオーガナイゼーションを取得する内部関数です。
func getPluralOrganizationsWithDetail(
	ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID,
	orderMethod parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	var e []query.GetPluralOrganizationsWithDetailRow
	var te []query.GetPluralOrganizationsWithDetailUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralOrganizationsWithDetail(ctx, query.GetPluralOrganizationsWithDetailParams{
			OrganizationIds: organizationIDs,
			OrderMethod:     orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralOrganizationsWithDetailUseNumberedPaginate(
			ctx, query.GetPluralOrganizationsWithDetailUseNumberedPaginateParams{
				OrganizationIds: organizationIDs,
				Offset:          int32(np.Offset.Int64),
				Limit:           int32(np.Limit.Int64),
				OrderMethod:     orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralOrganizationsWithDetailRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralOrganizationsWithDetailRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.OrganizationWithDetail]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	entities := make([]entity.OrganizationWithDetail, len(e))
	for i, v := range e {
		entities[i] = convOrganizationWithDetail(query.FindOrganizationByIDWithDetailRow(v))
	}
	return store.ListResult[entity.OrganizationWithDetail]{Data: entities}, nil
}

// GetPluralOrganizationsWithDetail は複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithDetail(
	ctx context.Context, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	return getPluralOrganizationsWithDetail(ctx, a.query, organizationIDs, order, np)
}

// GetPluralOrganizationsWithDetailWithSd はSD付きで複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithDetailWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getPluralOrganizationsWithDetail(ctx, qtx, organizationIDs, order, np)
}

func getOrganizationsWithChatRoomAndDetail(
	ctx context.Context, qtx *query.Queries, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	eConvFunc := func(
		e entity.OrganizationWithChatRoomAndDetailForQuery,
	) (entity.OrganizationWithChatRoomAndDetail, error) {
		return e.OrganizationWithChatRoomAndDetail, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountOrganizationsParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
		}
		r, err := qtx.CountOrganizations(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count organizations: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.OrganizationWithChatRoomAndDetailForQuery, error) {
		p := query.GetOrganizationsWithChatRoomAndDetailParams{
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoomAndDetail(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.OrganizationWithChatRoomAndDetailForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomAndDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomAndDetailForQuery{
				Pkey: entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoomAndDetail: convOrganizationWithChatRoomAndDetail(
					query.FindOrganizationByIDWithChatRoomAndDetailRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.OrganizationWithChatRoomAndDetailForQuery, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.OrganizationNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetOrganizationsWithChatRoomAndDetailUseKeysetPaginateParams{
			Limit:            limit,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			CursorDirection:  cursorDir,
			Cursor:           cursor,
			NameCursor:       nameCursor,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoomAndDetailUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomAndDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomAndDetailForQuery{
				Pkey: entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoomAndDetail: convOrganizationWithChatRoomAndDetail(
					query.FindOrganizationByIDWithChatRoomAndDetailRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(
		orderMethod string, limit, offset int32,
	) ([]entity.OrganizationWithChatRoomAndDetailForQuery, error) {
		p := query.GetOrganizationsWithChatRoomAndDetailUseNumberedPaginateParams{
			Limit:            limit,
			Offset:           offset,
			WhereLikeName:    where.WhereLikeName,
			SearchName:       where.SearchName,
			WhereIsWhole:     where.WhereIsWhole,
			IsWhole:          where.IsWhole,
			WhereIsPersonal:  where.WhereIsPersonal,
			IsPersonal:       where.IsPersonal,
			PersonalMemberID: where.PersonalMemberID,
			WhereIsGroup:     where.WhereIsGroup,
			WhereIsGrade:     where.WhereIsGrade,
			OrderMethod:      orderMethod,
		}
		r, err := qtx.GetOrganizationsWithChatRoomAndDetailUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}
		e := make([]entity.OrganizationWithChatRoomAndDetailForQuery, len(r))
		for i, v := range r {
			e[i] = entity.OrganizationWithChatRoomAndDetailForQuery{
				Pkey: entity.Int(v.MOrganizationsPkey),
				OrganizationWithChatRoomAndDetail: convOrganizationWithChatRoomAndDetail(
					query.FindOrganizationByIDWithChatRoomAndDetailRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.OrganizationWithChatRoomAndDetailForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.OrganizationDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.OrganizationNameCursorKey:
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
		return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{},
			fmt.Errorf("failed to get organizations: %w", err)
	}
	return res, nil
}

// GetOrganizationsWithChatRoomAndDetail はオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) GetOrganizationsWithChatRoomAndDetail(
	ctx context.Context, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	return getOrganizationsWithChatRoomAndDetail(ctx, a.query, where, order, np, cp, wc)
}

// GetOrganizationsWithChatRoomAndDetailWithSd はSD付きでオーガナイゼーションとチャットルームを取得します。
func (a *PgAdapter) GetOrganizationsWithChatRoomAndDetailWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereOrganizationParam,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{}, store.ErrNotFoundDescriptor
	}
	return getOrganizationsWithChatRoomAndDetail(ctx, qtx, where, order, np, cp, wc)
}

// getPluralOrganizationsWithChatRoomAndDetail は複数のオーガナイゼーションを取得する内部関数です。
func getPluralOrganizationsWithChatRoomAndDetail(
	ctx context.Context, qtx *query.Queries, organizationIDs []uuid.UUID,
	orderMethod parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	var e []query.GetPluralOrganizationsWithChatRoomAndDetailRow
	var te []query.GetPluralOrganizationsWithChatRoomAndDetailUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralOrganizationsWithChatRoomAndDetail(ctx, query.GetPluralOrganizationsWithChatRoomAndDetailParams{
			OrganizationIds: organizationIDs,
			OrderMethod:     orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralOrganizationsWithChatRoomAndDetailUseNumberedPaginate(
			ctx, query.GetPluralOrganizationsWithChatRoomAndDetailUseNumberedPaginateParams{
				OrganizationIds: organizationIDs,
				Offset:          int32(np.Offset.Int64),
				Limit:           int32(np.Limit.Int64),
				OrderMethod:     orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralOrganizationsWithChatRoomAndDetailRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralOrganizationsWithChatRoomAndDetailRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{},
			fmt.Errorf("failed to get organizations: %w", err)
	}
	entities := make([]entity.OrganizationWithChatRoomAndDetail, len(e))
	for i, v := range e {
		entities[i] = convOrganizationWithChatRoomAndDetail(query.FindOrganizationByIDWithChatRoomAndDetailRow(v))
	}
	return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{Data: entities}, nil
}

// GetPluralOrganizationsWithChatRoomAndDetail は複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithChatRoomAndDetail(
	ctx context.Context, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	return getPluralOrganizationsWithChatRoomAndDetail(ctx, a.query, organizationIDs, order, np)
}

// GetPluralOrganizationsWithChatRoomAndDetailWithSd はSD付きで複数のオーガナイゼーションを取得します。
func (a *PgAdapter) GetPluralOrganizationsWithChatRoomAndDetailWithSd(
	ctx context.Context, sd store.Sd, organizationIDs []uuid.UUID,
	order parameter.OrganizationOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{}, store.ErrNotFoundDescriptor
	}
	return getPluralOrganizationsWithChatRoomAndDetail(ctx, qtx, organizationIDs, order, np)
}

// updateOrganization はオーガナイゼーションを更新する内部関数です。
func updateOrganization(
	ctx context.Context, qtx *query.Queries,
	organizationID uuid.UUID, param parameter.UpdateOrganizationParams, now time.Time,
) (entity.Organization, error) {
	p := query.UpdateOrganizationParams{
		OrganizationID: organizationID,
		Name:           param.Name,
		Description:    pgtype.Text(param.Description),
		Color:          pgtype.Text(param.Color),
		UpdatedAt:      now,
	}
	e, err := qtx.UpdateOrganization(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Organization{}, errhandle.NewModelNotFoundError("organization")
		}
		return entity.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}
	entity := entity.Organization{
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    entity.String(e.Description),
		Color:          entity.String(e.Color),
		IsPersonal:     e.IsPersonal,
		IsWhole:        e.IsWhole,
		ChatRoomID:     entity.UUID(e.ChatRoomID),
	}
	return entity, nil
}

// UpdateOrganization はオーガナイゼーションを更新します。
func (a *PgAdapter) UpdateOrganization(
	ctx context.Context, organizationID uuid.UUID, param parameter.UpdateOrganizationParams,
) (entity.Organization, error) {
	return updateOrganization(ctx, a.query, organizationID, param, a.clocker.Now())
}

// UpdateOrganizationWithSd はSD付きでオーガナイゼーションを更新します。
func (a *PgAdapter) UpdateOrganizationWithSd(
	ctx context.Context, sd store.Sd, organizationID uuid.UUID, param parameter.UpdateOrganizationParams,
) (entity.Organization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Organization{}, store.ErrNotFoundDescriptor
	}
	return updateOrganization(ctx, qtx, organizationID, param, a.clocker.Now())
}

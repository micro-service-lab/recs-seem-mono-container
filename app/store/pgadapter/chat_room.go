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

func convChatRoomWithCoverImage(e query.FindChatRoomByIDWithCoverImageRow) entity.ChatRoomWithCoverImage {
	var image entity.NullableEntity[entity.ImageWithAttachableItem]
	if e.CoverImageID.Valid {
		image = entity.NullableEntity[entity.ImageWithAttachableItem]{
			Valid: true,
			Entity: entity.ImageWithAttachableItem{
				ImageID: e.CoverImageID.Bytes,
				Height:  entity.Float(e.CoverImageHeight),
				Width:   entity.Float(e.CoverImageWidth),
				AttachableItem: entity.AttachableItem{
					AttachableItemID: e.CoverImageAttachableItemID.Bytes,
					OwnerID:          entity.UUID(e.CoverImageOwnerID),
					FromOuter:        e.CoverImageFromOuter.Bool,
					URL:              e.CoverImageUrl.String,
					Alias:            e.CoverImageAlias.String,
					Size:             entity.Float(e.CoverImageSize),
					MimeTypeID:       e.CoverImageMimeTypeID.Bytes,
					ImageID:          entity.UUID(e.CoverImageID),
				},
			},
		}
	}
	return entity.ChatRoomWithCoverImage{
		ChatRoomID:       e.ChatRoomID,
		Name:             e.Name,
		IsPrivate:        e.IsPrivate,
		FromOrganization: e.FromOrganization,
		OwnerID:          entity.UUID(e.OwnerID),
		CoverImage:       image,
	}
}

func countChatRooms(ctx context.Context, qtx *query.Queries, where parameter.WhereChatRoomParam) (int64, error) {
	p := query.CountChatRoomsParams{
		WhereInOwner:            where.WhereInOwner,
		InOwner:                 where.InOwner,
		WhereIsPrivate:          where.WhereIsPrivate,
		IsPrivate:               where.IsPrivate,
		WhereIsFromOrganization: where.WhereIsFromOrganization,
		IsFromOrganization:      where.IsFromOrganization,
		WhereFromOrganizations:  where.WhereFromOrganizations,
		FromOrganizations:       where.FromOrganizations,
	}
	c, err := qtx.CountChatRooms(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat rooms: %w", err)
	}
	return c, nil
}

// CountChatRooms チャットルーム数を取得する。
func (a *PgAdapter) CountChatRooms(ctx context.Context, where parameter.WhereChatRoomParam) (int64, error) {
	return countChatRooms(ctx, a.query, where)
}

// CountChatRoomsWithSd SD付きでチャットルーム数を取得する。
func (a *PgAdapter) CountChatRoomsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRooms(ctx, qtx, where)
}

func createChatRoom(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomParam, now time.Time,
) (entity.ChatRoom, error) {
	p := query.CreateChatRoomParams{
		Name:             param.Name,
		IsPrivate:        param.IsPrivate,
		CoverImageID:     pgtype.UUID(param.CoverImageID),
		OwnerID:          pgtype.UUID(param.OwnerID),
		FromOrganization: param.FromOrganization,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	e, err := qtx.CreateChatRoom(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoom{}, errhandle.NewModelDuplicatedError("chat room")
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	entity := entity.ChatRoom{
		ChatRoomID:       e.ChatRoomID,
		Name:             e.Name,
		IsPrivate:        e.IsPrivate,
		FromOrganization: e.FromOrganization,
		CoverImageID:     entity.UUID(e.CoverImageID),
		OwnerID:          entity.UUID(e.OwnerID),
	}
	return entity, nil
}

// CreateChatRoom チャットルームを作成する。
func (a *PgAdapter) CreateChatRoom(ctx context.Context, param parameter.CreateChatRoomParam) (entity.ChatRoom, error) {
	return createChatRoom(ctx, a.query, param, a.clocker.Now())
}

// CreateChatRoomWithSd SD付きでチャットルームを作成する。
func (a *PgAdapter) CreateChatRoomWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomParam,
) (entity.ChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoom{}, store.ErrNotFoundDescriptor
	}
	return createChatRoom(ctx, qtx, param, a.clocker.Now())
}

func createChatRooms(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomParam, now time.Time,
) (int64, error) {
	param := make([]query.CreateChatRoomsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomsParams{
			Name:             p.Name,
			IsPrivate:        p.IsPrivate,
			CoverImageID:     pgtype.UUID(p.CoverImageID),
			OwnerID:          pgtype.UUID(p.OwnerID),
			FromOrganization: p.FromOrganization,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
	}
	n, err := qtx.CreateChatRooms(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room")
		}
		return 0, fmt.Errorf("failed to create chat rooms: %w", err)
	}
	return n, nil
}

// CreateChatRooms チャットルームを作成する。
func (a *PgAdapter) CreateChatRooms(
	ctx context.Context, params []parameter.CreateChatRoomParam,
) (int64, error) {
	return createChatRooms(ctx, a.query, params, a.clocker.Now())
}

// CreateChatRoomsWithSd SD付きでチャットルームを作成する。
func (a *PgAdapter) CreateChatRoomsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRooms(ctx, qtx, params, a.clocker.Now())
}

func deleteChatRoom(ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteChatRoom(ctx, chatRoomID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room")
	}
	return c, nil
}

// DeleteChatRoom チャットルームを削除する。
func (a *PgAdapter) DeleteChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error) {
	return deleteChatRoom(ctx, a.query, chatRoomID)
}

// DeleteChatRoomWithSd SD付きでチャットルームを削除する。
func (a *PgAdapter) DeleteChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoom(ctx, qtx, chatRoomID)
}

func pluralDeleteChatRooms(ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteChatRooms(ctx, chatRoomIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat rooms: %w", err)
	}
	if c != int64(len(chatRoomIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room")
	}
	return c, nil
}

// PluralDeleteChatRooms チャットルームを複数削除する。
func (a *PgAdapter) PluralDeleteChatRooms(ctx context.Context, chatRoomIDs []uuid.UUID) (int64, error) {
	return pluralDeleteChatRooms(ctx, a.query, chatRoomIDs)
}

// PluralDeleteChatRoomsWithSd SD付きでチャットルームを複数削除する。
func (a *PgAdapter) PluralDeleteChatRoomsWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRooms(ctx, qtx, chatRoomIDs)
}

func findChatRoomByID(ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID) (entity.ChatRoom, error) {
	e, err := qtx.FindChatRoomByID(ctx, chatRoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError("chat room")
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room: %w", err)
	}
	entity := entity.ChatRoom{
		ChatRoomID:       e.ChatRoomID,
		Name:             e.Name,
		IsPrivate:        e.IsPrivate,
		FromOrganization: e.FromOrganization,
		CoverImageID:     entity.UUID(e.CoverImageID),
		OwnerID:          entity.UUID(e.OwnerID),
	}
	return entity, nil
}

// FindChatRoomByID チャットルームを取得する。
func (a *PgAdapter) FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (entity.ChatRoom, error) {
	return findChatRoomByID(ctx, a.query, chatRoomID)
}

// FindChatRoomByIDWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) FindChatRoomByIDWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
) (entity.ChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoom{}, store.ErrNotFoundDescriptor
	}
	return findChatRoomByID(ctx, qtx, chatRoomID)
}

func findChatRoomByIDWithCoverImage(
	ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID,
) (entity.ChatRoomWithCoverImage, error) {
	e, err := qtx.FindChatRoomByIDWithCoverImage(ctx, chatRoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomWithCoverImage{}, errhandle.NewModelNotFoundError("chat room")
		}
		return entity.ChatRoomWithCoverImage{}, fmt.Errorf("failed to find chat room: %w", err)
	}
	return convChatRoomWithCoverImage(e), nil
}

// FindChatRoomByIDWithCoverImage チャットルームを取得する。
func (a *PgAdapter) FindChatRoomByIDWithCoverImage(
	ctx context.Context, chatRoomID uuid.UUID,
) (entity.ChatRoomWithCoverImage, error) {
	return findChatRoomByIDWithCoverImage(ctx, a.query, chatRoomID)
}

// FindChatRoomByIDWithCoverImageWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) FindChatRoomByIDWithCoverImageWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
) (entity.ChatRoomWithCoverImage, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomWithCoverImage{}, store.ErrNotFoundDescriptor
	}
	return findChatRoomByIDWithCoverImage(ctx, qtx, chatRoomID)
}

func findChatRoomOnPrivate(
	ctx context.Context, qtx *query.Queries, ownerID, memberID uuid.UUID,
) (entity.ChatRoom, error) {
	p := query.FindChatRoomOnPrivateParams{
		OwnerID:  ownerID,
		MemberID: memberID,
	}
	e, err := qtx.FindChatRoomOnPrivate(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError("chat room")
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room: %w", err)
	}
	entity := entity.ChatRoom{
		ChatRoomID:       e.ChatRoomID,
		Name:             e.Name,
		IsPrivate:        e.IsPrivate,
		FromOrganization: e.FromOrganization,
		CoverImageID:     entity.UUID(e.CoverImageID),
		OwnerID:          entity.UUID(e.OwnerID),
	}
	return entity, nil
}

// FindChatRoomOnPrivate チャットルームを取得する。
func (a *PgAdapter) FindChatRoomOnPrivate(
	ctx context.Context, ownerID, memberID uuid.UUID,
) (entity.ChatRoom, error) {
	return findChatRoomOnPrivate(ctx, a.query, ownerID, memberID)
}

// FindChatRoomOnPrivateWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) FindChatRoomOnPrivateWithSd(
	ctx context.Context, sd store.Sd, ownerID, memberID uuid.UUID,
) (entity.ChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoom{}, store.ErrNotFoundDescriptor
	}
	return findChatRoomOnPrivate(ctx, qtx, ownerID, memberID)
}

func getChatRooms(
	ctx context.Context,
	qtx *query.Queries,
	where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoom], error) {
	eConvFunc := func(e query.ChatRoom) (entity.ChatRoom, error) {
		return entity.ChatRoom{
			ChatRoomID:       e.ChatRoomID,
			Name:             e.Name,
			IsPrivate:        e.IsPrivate,
			FromOrganization: e.FromOrganization,
			CoverImageID:     entity.UUID(e.CoverImageID),
			OwnerID:          entity.UUID(e.OwnerID),
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountChatRoomsParams{
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.CountChatRooms(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat rooms: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.ChatRoom, error) {
		p := query.GetChatRoomsParams{
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.GetChatRooms(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.ChatRoom{}, nil
			}
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.ChatRoom, error) {
		p := query.GetChatRoomsUseKeysetPaginateParams{
			Limit:                   limit,
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
			CursorDirection:         cursorDir,
			Cursor:                  cursor,
		}
		r, err := qtx.GetChatRoomsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.ChatRoom, error) {
		p := query.GetChatRoomsUseNumberedPaginateParams{
			Limit:                   limit,
			Offset:                  offset,
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.GetChatRoomsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.ChatRoom) (entity.Int, any) {
		switch subCursor {
		case parameter.PermissionDefaultCursorKey:
			return entity.Int(e.MChatRoomsPkey), nil
		}
		return entity.Int(e.MChatRoomsPkey), nil
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
		return store.ListResult[entity.ChatRoom]{}, fmt.Errorf("failed to get chat rooms: %w", err)
	}
	return res, nil
}

// GetChatRooms チャットルームを取得する。
func (a *PgAdapter) GetChatRooms(
	ctx context.Context, where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoom], error) {
	return getChatRooms(ctx, a.query, where, order, np, cp, wc)
}

// GetChatRoomsWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) GetChatRoomsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getChatRooms(ctx, qtx, where, order, np, cp, wc)
}

func getPluralChatRooms(
	ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID,
	_ parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoom], error) {
	var e []query.ChatRoom
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRooms(ctx, chatRoomIDs)
	} else {
		e, err = qtx.GetPluralChatRoomsUseNumberedPaginate(ctx, query.GetPluralChatRoomsUseNumberedPaginateParams{
			ChatRoomIds: chatRoomIDs,
			Offset:      int32(np.Offset.Int64),
			Limit:       int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.ChatRoom]{}, fmt.Errorf("failed to get chat rooms: %w", err)
	}
	entities := make([]entity.ChatRoom, len(e))
	for i, v := range e {
		entities[i] = entity.ChatRoom{
			ChatRoomID:       v.ChatRoomID,
			Name:             v.Name,
			IsPrivate:        v.IsPrivate,
			FromOrganization: v.FromOrganization,
			CoverImageID:     entity.UUID(v.CoverImageID),
			OwnerID:          entity.UUID(v.OwnerID),
		}
	}
	return store.ListResult[entity.ChatRoom]{Data: entities}, nil
}

// GetPluralChatRooms チャットルームを取得する。
func (a *PgAdapter) GetPluralChatRooms(
	ctx context.Context, chatRoomIDs []uuid.UUID,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoom], error) {
	return getPluralChatRooms(ctx, a.query, chatRoomIDs, order, np)
}

// GetPluralChatRoomsWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) GetPluralChatRoomsWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRooms(ctx, qtx, chatRoomIDs, order, np)
}

func getChatRoomsWithCoverImage(
	ctx context.Context,
	qtx *query.Queries,
	where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	eConvFunc := func(e query.FindChatRoomByIDWithCoverImageRow) (entity.ChatRoomWithCoverImage, error) {
		return convChatRoomWithCoverImage(e), nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountChatRoomsParams{
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.CountChatRooms(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat rooms: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.FindChatRoomByIDWithCoverImageRow, error) {
		p := query.GetChatRoomsWithCoverImageParams{
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.GetChatRoomsWithCoverImage(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.FindChatRoomByIDWithCoverImageRow{}, nil
			}
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		var es []query.FindChatRoomByIDWithCoverImageRow
		for _, v := range r {
			es = append(es, query.FindChatRoomByIDWithCoverImageRow(v))
		}
		return es, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.FindChatRoomByIDWithCoverImageRow, error) {
		p := query.GetChatRoomsWithCoverImageUseKeysetPaginateParams{
			Limit:                   int32(limit),
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
			CursorDirection:         cursorDir,
			Cursor:                  int32(cursor),
		}
		r, err := qtx.GetChatRoomsWithCoverImageUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		var es []query.FindChatRoomByIDWithCoverImageRow
		for _, v := range r {
			es = append(es, query.FindChatRoomByIDWithCoverImageRow(v))
		}
		return es, nil
	}
	runQNPFunc := func(
		_ string, limit, offset int32,
	) ([]query.FindChatRoomByIDWithCoverImageRow, error) {
		p := query.GetChatRoomsWithCoverImageUseNumberedPaginateParams{
			Limit:                   int32(limit),
			Offset:                  int32(offset),
			WhereInOwner:            where.WhereInOwner,
			InOwner:                 where.InOwner,
			WhereIsPrivate:          where.WhereIsPrivate,
			IsPrivate:               where.IsPrivate,
			WhereIsFromOrganization: where.WhereIsFromOrganization,
			IsFromOrganization:      where.IsFromOrganization,
			WhereFromOrganizations:  where.WhereFromOrganizations,
			FromOrganizations:       where.FromOrganizations,
		}
		r, err := qtx.GetChatRoomsWithCoverImageUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms: %w", err)
		}
		var es []query.FindChatRoomByIDWithCoverImageRow
		for _, v := range r {
			es = append(es, query.FindChatRoomByIDWithCoverImageRow(v))
		}
		return es, nil
	}
	selector := func(subCursor string, e query.FindChatRoomByIDWithCoverImageRow) (entity.Int, any) {
		switch subCursor {
		case parameter.PermissionDefaultCursorKey:
			return entity.Int(e.MChatRoomsPkey), nil
		}
		return entity.Int(e.MChatRoomsPkey), nil
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
		return store.ListResult[entity.ChatRoomWithCoverImage]{}, fmt.Errorf("failed to get chat rooms: %w", err)
	}
	return res, nil
}

// GetChatRoomsWithCoverImage チャットルームを取得する。
func (a *PgAdapter) GetChatRoomsWithCoverImage(
	ctx context.Context, where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	return getChatRoomsWithCoverImage(ctx, a.query, where, order, np, cp, wc)
}

// GetChatRoomsWithCoverImageWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) GetChatRoomsWithCoverImageWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomParam,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomWithCoverImage]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomsWithCoverImage(ctx, qtx, where, order, np, cp, wc)
}

func getPluralChatRoomsWithCoverImage(
	ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID,
	_ parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	var e []query.FindChatRoomByIDWithCoverImageRow
	var err error
	if !np.Valid {
		var ge []query.GetPluralChatRoomsWithCoverImageRow
		ge, err = qtx.GetPluralChatRoomsWithCoverImage(ctx, chatRoomIDs)
		e = make([]query.FindChatRoomByIDWithCoverImageRow, len(ge))
		for i, v := range ge {
			e[i] = query.FindChatRoomByIDWithCoverImageRow(v)
		}
	} else {
		var ge []query.GetPluralChatRoomsWithCoverImageUseNumberedPaginateRow
		ge, err = qtx.GetPluralChatRoomsWithCoverImageUseNumberedPaginate(
			ctx, query.GetPluralChatRoomsWithCoverImageUseNumberedPaginateParams{
				ChatRoomIds: chatRoomIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
			})
		e = make([]query.FindChatRoomByIDWithCoverImageRow, len(ge))
		for i, v := range ge {
			e[i] = query.FindChatRoomByIDWithCoverImageRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomWithCoverImage]{}, fmt.Errorf("failed to get chat rooms: %w", err)
	}
	entities := make([]entity.ChatRoomWithCoverImage, len(e))
	for i, v := range e {
		entities[i] = convChatRoomWithCoverImage(v)
	}
	return store.ListResult[entity.ChatRoomWithCoverImage]{Data: entities}, nil
}

// GetPluralChatRoomsWithCoverImage チャットルームを取得する。
func (a *PgAdapter) GetPluralChatRoomsWithCoverImage(
	ctx context.Context, chatRoomIDs []uuid.UUID,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	return getPluralChatRoomsWithCoverImage(ctx, a.query, chatRoomIDs, order, np)
}

// GetPluralChatRoomsWithCoverImageWithSd SD付きでチャットルームを取得する。
func (a *PgAdapter) GetPluralChatRoomsWithCoverImageWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
	order parameter.ChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithCoverImage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomWithCoverImage]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomsWithCoverImage(ctx, qtx, chatRoomIDs, order, np)
}

func updateChatRoom(
	ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID, param parameter.UpdateChatRoomParams, now time.Time,
) (entity.ChatRoom, error) {
	p := query.UpdateChatRoomParams{
		ChatRoomID:   chatRoomID,
		Name:         param.Name,
		CoverImageID: pgtype.UUID(param.CoverImageID),
		UpdatedAt:    now,
	}
	e, err := qtx.UpdateChatRoom(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError("chat room")
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to update chat room: %w", err)
	}
	entity := entity.ChatRoom{
		ChatRoomID:       e.ChatRoomID,
		Name:             e.Name,
		IsPrivate:        e.IsPrivate,
		FromOrganization: e.FromOrganization,
		CoverImageID:     entity.UUID(e.CoverImageID),
		OwnerID:          entity.UUID(e.OwnerID),
	}
	return entity, nil
}

// UpdateChatRoom チャットルームを更新する。
func (a *PgAdapter) UpdateChatRoom(
	ctx context.Context, chatRoomID uuid.UUID, param parameter.UpdateChatRoomParams,
) (entity.ChatRoom, error) {
	return updateChatRoom(ctx, a.query, chatRoomID, param, a.clocker.Now())
}

// UpdateChatRoomWithSd SD付きでチャットルームを更新する。
func (a *PgAdapter) UpdateChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID, param parameter.UpdateChatRoomParams,
) (entity.ChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoom{}, store.ErrNotFoundDescriptor
	}
	return updateChatRoom(ctx, qtx, chatRoomID, param, a.clocker.Now())
}

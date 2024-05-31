package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convChatRoomOnMember(r query.GetChatRoomsOnMemberRow) entity.ChatRoomOnMember {
	var coverImg entity.NullableEntity[entity.ImageWithAttachableItem]
	if r.ChatRoomCoverImageID.Valid {
		coverImg = entity.NullableEntity[entity.ImageWithAttachableItem]{
			Valid: true,
			Entity: entity.ImageWithAttachableItem{
				ImageID: r.ChatRoomCoverImageID.Bytes,
				Height:  entity.Float(r.ChatRoomCoverImageHeight),
				Width:   entity.Float(r.ChatRoomCoverImageWidth),
				AttachableItem: entity.AttachableItem{
					AttachableItemID: r.ChatRoomCoverImageAttachableItemID.Bytes,
					OwnerID:          entity.UUID(r.ChatRoomCoverImageOwnerID),
					FromOuter:        r.ChatRoomCoverImageFromOuter.Bool,
					URL:              r.ChatRoomCoverImageUrl.String,
					Alias:            r.ChatRoomCoverImageAlias.String,
					Size:             entity.Float(r.ChatRoomCoverImageSize),
					MimeTypeID:       r.ChatRoomCoverImageMimeTypeID.Bytes,
				},
			},
		}
	}
	var latestMessage entity.NullableEntity[entity.MessageCard]
	if r.ChatRoomLatestMessageID != uuid.Nil {
		latestMessage = entity.NullableEntity[entity.MessageCard]{
			Valid: true,
			Entity: entity.MessageCard{
				MessageID: r.ChatRoomLatestMessageID,
				Body:      r.ChatRoomLatestMessageBody,
				PostedAt:  r.ChatRoomLatestMessagePostedAt,
			},
		}
	}
	return entity.ChatRoomOnMember{
		ChatRoom: entity.PracticalChatRoom{
			ChatRoomID:       r.ChatRoomID,
			Name:             r.ChatRoomName.String,
			IsPrivate:        r.ChatRoomIsPrivate.Bool,
			FromOrganization: r.ChatRoomFromOrganization.Bool,
			CoverImage:       coverImg,
			LatestMessage:    latestMessage,
			OwnerID:          entity.UUID(r.ChatRoomOwnerID),
		},
		AddedAt: r.AddedAt,
	}
}

func convMemberOnChatRoom(r query.GetMembersOnChatRoomRow) entity.MemberOnChatRoom {
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
	return entity.MemberOnChatRoom{
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
		AddedAt: r.AddedAt,
	}
}

func countChatRoomsOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	where parameter.WhereChatRoomOnMemberParam,
) (int64, error) {
	p := query.CountChatRoomsOnMemberParams{
		MemberID:      memberID,
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountChatRoomsOnMember(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat rooms on member: %w", err)
	}
	return c, nil
}

// CountChatRoomsOnMember メンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountChatRoomsOnMember(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereChatRoomOnMemberParam,
) (int64, error) {
	return countChatRoomsOnMember(ctx, a.query, memberID, where)
}

// CountChatRoomsOnMemberWithSd SD付きでメンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountChatRoomsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID, where parameter.WhereChatRoomOnMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomsOnMember(ctx, qtx, memberID, where)
}

func countMembersOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	where parameter.WhereMemberOnChatRoomParam,
) (int64, error) {
	p := query.CountMembersOnChatRoomParams{
		ChatRoomID:    chatRoomID,
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountMembersOnChatRoom(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count members on chat room: %w", err)
	}
	return c, nil
}

// CountMembersOnChatRoom チャットルーム上のメンバー数を取得する。
func (a *PgAdapter) CountMembersOnChatRoom(
	ctx context.Context, chatRoomID uuid.UUID, where parameter.WhereMemberOnChatRoomParam,
) (int64, error) {
	return countMembersOnChatRoom(ctx, a.query, chatRoomID, where)
}

// CountMembersOnChatRoomWithSd SD付きでチャットルーム上のメンバー数を取得する。
func (a *PgAdapter) CountMembersOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID, where parameter.WhereMemberOnChatRoomParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMembersOnChatRoom(ctx, qtx, chatRoomID, where)
}

func belongChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	param parameter.BelongChatRoomParam,
) (entity.ChatRoomBelonging, error) {
	p := query.CreateChatRoomBelongingParams{
		MemberID:   param.MemberID,
		ChatRoomID: param.ChatRoomID,
		AddedAt:    param.AddedAt,
	}
	b, err := qtx.CreateChatRoomBelonging(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomBelonging{}, errhandle.NewModelNotFoundError("chat room belonging")
		}
		return entity.ChatRoomBelonging{}, fmt.Errorf("failed to belong chat room: %w", err)
	}
	return entity.ChatRoomBelonging{
		MemberID:   b.MemberID,
		ChatRoomID: b.ChatRoomID,
		AddedAt:    b.AddedAt,
	}, nil
}

// BelongChatRoom メンバーをチャットルームに所属させる。
func (a *PgAdapter) BelongChatRoom(
	ctx context.Context, param parameter.BelongChatRoomParam,
) (entity.ChatRoomBelonging, error) {
	return belongChatRoom(ctx, a.query, param)
}

// BelongChatRoomWithSd SD付きでメンバーをチャットルームに所属させる。
func (a *PgAdapter) BelongChatRoomWithSd(
	ctx context.Context, sd store.Sd, param parameter.BelongChatRoomParam,
) (entity.ChatRoomBelonging, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomBelonging{}, store.ErrNotFoundDescriptor
	}
	return belongChatRoom(ctx, qtx, param)
}

func belongChatRooms(
	ctx context.Context,
	qtx *query.Queries,
	params []parameter.BelongChatRoomParam,
) (int64, error) {
	param := make([]query.CreateChatRoomBelongingsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomBelongingsParams{
			MemberID:   p.MemberID,
			ChatRoomID: p.ChatRoomID,
			AddedAt:    p.AddedAt,
		}
	}
	b, err := qtx.CreateChatRoomBelongings(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelNotFoundError("chat room belonging")
		}
		return 0, fmt.Errorf("failed to belong chat rooms: %w", err)
	}
	return b, nil
}

// BelongChatRooms メンバーを複数のチャットルームに所属させる。
func (a *PgAdapter) BelongChatRooms(
	ctx context.Context, params []parameter.BelongChatRoomParam,
) (int64, error) {
	return belongChatRooms(ctx, a.query, params)
}

// BelongChatRoomsWithSd SD付きでメンバーを複数のチャットルームに所属させる。
func (a *PgAdapter) BelongChatRoomsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.BelongChatRoomParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return belongChatRooms(ctx, qtx, params)
}

func disbelongChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	chatRoomID uuid.UUID,
) (int64, error) {
	p := query.DeleteChatRoomBelongingParams{
		MemberID:   memberID,
		ChatRoomID: chatRoomID,
	}
	b, err := qtx.DeleteChatRoomBelonging(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room: %w", err)
	}
	if b != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room belonging")
	}
	return b, nil
}

// DisbelongChatRoom メンバーをチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoom(
	ctx context.Context, memberID, chatRoomID uuid.UUID,
) (int64, error) {
	return disbelongChatRoom(ctx, a.query, memberID, chatRoomID)
}

// DisbelongChatRoomWithSd SD付きでメンバーをチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoomWithSd(
	ctx context.Context, sd store.Sd, memberID, chatRoomID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongChatRoom(ctx, qtx, memberID, chatRoomID)
}

func disbelongChatRoomOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteChatRoomBelongingsOnMember(ctx, memberID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on member: %w", err)
	}
	return b, nil
}

// DisbelongChatRoomOnMember メンバー上のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnMember(ctx context.Context, memberID uuid.UUID) (int64, error) {
	return disbelongChatRoomOnMember(ctx, a.query, memberID)
}

// DisbelongChatRoomOnMemberWithSd SD付きでメンバー上のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongChatRoomOnMember(ctx, qtx, memberID)
}

func disbelongChatRoomOnMembers(
	ctx context.Context,
	qtx *query.Queries,
	memberIDs []uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteChatRoomBelongingsOnMembers(ctx, memberIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on members: %w", err)
	}
	return b, nil
}

// DisbelongChatRoomOnMembers メンバー上の複数のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnMembers(
	ctx context.Context, memberIDs []uuid.UUID,
) (int64, error) {
	return disbelongChatRoomOnMembers(ctx, a.query, memberIDs)
}

// DisbelongChatRoomOnMembersWithSd SD付きでメンバー上の複数のチャットルームから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongChatRoomOnMembers(ctx, qtx, memberIDs)
}

func disbelongChatRoomOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteChatRoomBelongingsOnChatRoom(ctx, chatRoomID)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on chat room: %w", err)
	}
	return b, nil
}

// DisbelongChatRoomOnChatRoom チャットルーム上のメンバーから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error) {
	return disbelongChatRoomOnChatRoom(ctx, a.query, chatRoomID)
}

// DisbelongChatRoomOnChatRoomWithSd SD付きでチャットルーム上のメンバーから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongChatRoomOnChatRoom(ctx, qtx, chatRoomID)
}

func disbelongChatRoomOnChatRooms(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomIDs []uuid.UUID,
) (int64, error) {
	b, err := qtx.DeleteChatRoomBelongingsOnChatRooms(ctx, chatRoomIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to disbelong chat room on chat rooms: %w", err)
	}
	return b, nil
}

// DisbelongChatRoomOnChatRooms チャットルーム上の複数のメンバーから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnChatRooms(
	ctx context.Context, chatRoomIDs []uuid.UUID,
) (int64, error) {
	return disbelongChatRoomOnChatRooms(ctx, a.query, chatRoomIDs)
}

// DisbelongChatRoomOnChatRoomsWithSd SD付きでチャットルーム上の複数のメンバーから所属解除する。
func (a *PgAdapter) DisbelongChatRoomOnChatRoomsWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return disbelongChatRoomOnChatRooms(ctx, qtx, chatRoomIDs)
}

func getChatRoomsOnMember(
	ctx context.Context,
	qtx *query.Queries,
	memberID uuid.UUID,
	where parameter.WhereChatRoomOnMemberParam,
	order parameter.ChatRoomOnMemberOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	eConvFunc := func(e entity.ChatRoomOnMemberForQuery) (entity.ChatRoomOnMember, error) {
		return e.ChatRoomOnMember, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountChatRoomsOnMemberParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountChatRoomsOnMember(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat rooms on member: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.ChatRoomOnMemberForQuery, error) {
		p := query.GetChatRoomsOnMemberParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetChatRoomsOnMember(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomOnMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat rooms on member: %w", err)
		}
		fq := make([]entity.ChatRoomOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ChatRoomOnMemberForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				ChatRoomOnMember: convChatRoomOnMember(e),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.ChatRoomOnMemberForQuery, error) {
		var addCursor time.Time
		var lastChatCursor time.Time
		var nameCursor string
		var ok bool
		var err error
		switch subCursor {
		case parameter.ChatRoomOnMemberNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		case parameter.ChatRoomOnMemberAddedAtCursorKey:
			cv, ok := subCursorValue.(string)
			addCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				addCursor = time.Time{}
			}
		case parameter.ChatRoomOnMemberLastChatAtCursorKey:
			cv, ok := subCursorValue.(string)
			lastChatCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				lastChatCursor = time.Time{}
			}
		}
		p := query.GetChatRoomsOnMemberUseKeysetPaginateParams{
			MemberID:        memberID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
			NameCursor:      nameCursor,
			AddCursor:       addCursor,
			ChatCursor:      lastChatCursor,
		}
		r, err := qtx.GetChatRoomsOnMemberUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms on member: %w", err)
		}
		fq := make([]entity.ChatRoomOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ChatRoomOnMemberForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				ChatRoomOnMember: convChatRoomOnMember(query.GetChatRoomsOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.ChatRoomOnMemberForQuery, error) {
		p := query.GetChatRoomsOnMemberUseNumberedPaginateParams{
			MemberID:      memberID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Limit:         limit,
			Offset:        offset,
		}
		r, err := qtx.GetChatRoomsOnMemberUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat rooms on member: %w", err)
		}
		fq := make([]entity.ChatRoomOnMemberForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.ChatRoomOnMemberForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				ChatRoomOnMember: convChatRoomOnMember(query.GetChatRoomsOnMemberRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.ChatRoomOnMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomOnMemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.ChatRoomOnMemberNameCursorKey:
			return entity.Int(e.Pkey), e.ChatRoomOnMember.ChatRoom.Name
		case parameter.ChatRoomOnMemberAddedAtCursorKey:
			return entity.Int(e.Pkey), e.ChatRoomOnMember.AddedAt
		case parameter.ChatRoomOnMemberLastChatAtCursorKey:
			return entity.Int(e.Pkey), e.ChatRoomOnMember.ChatRoom.LatestMessage.Entity.PostedAt
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
		return store.ListResult[entity.ChatRoomOnMember]{}, fmt.Errorf("failed to get chat rooms on member: %w", err)
	}
	return res, nil
}

// GetChatRoomsOnMember メンバー上のチャットルームを取得する。
func (a *PgAdapter) GetChatRoomsOnMember(
	ctx context.Context, memberID uuid.UUID, where parameter.WhereChatRoomOnMemberParam,
	order parameter.ChatRoomOnMemberOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	return getChatRoomsOnMember(ctx, a.query, memberID, where, order, np, cp, wc)
}

// GetChatRoomsOnMemberWithSd SD付きでメンバー上のチャットルームを取得する。
func (a *PgAdapter) GetChatRoomsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
	where parameter.WhereChatRoomOnMemberParam, order parameter.ChatRoomOnMemberOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomsOnMember(ctx, qtx, memberID, where, order, np, cp, wc)
}

func getMembersOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	where parameter.WhereMemberOnChatRoomParam,
	order parameter.MemberOnChatRoomOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	eConvFunc := func(e entity.MemberOnChatRoomForQuery) (entity.MemberOnChatRoom, error) {
		return e.MemberOnChatRoom, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountMembersOnChatRoomParams{
			ChatRoomID:    chatRoomID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountMembersOnChatRoom(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count members on chat room: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MemberOnChatRoomForQuery, error) {
		p := query.GetMembersOnChatRoomParams{
			ChatRoomID:    chatRoomID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetMembersOnChatRoom(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberOnChatRoomForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members on chat room: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnChatRoomForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				MemberOnChatRoom: convMemberOnChatRoom(e),
			}
		}
		return fq, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MemberOnChatRoomForQuery, error) {
		var addCursor time.Time
		var nameCursor string
		var ok bool
		var err error
		switch subCursor {
		case parameter.MemberOnChatRoomNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		case parameter.MemberOnChatRoomAddedAtCursorKey:
			cv, ok := subCursorValue.(string)
			addCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				addCursor = time.Time{}
			}
		}
		p := query.GetMembersOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
			AddedAtCursor:   addCursor,
		}
		r, err := qtx.GetMembersOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnChatRoomForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				MemberOnChatRoom: convMemberOnChatRoom(query.GetMembersOnChatRoomRow(e)),
			}
		}
		return fq, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MemberOnChatRoomForQuery, error) {
		p := query.GetMembersOnChatRoomUseNumberedPaginateParams{
			ChatRoomID:    chatRoomID,
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetMembersOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.MemberOnChatRoomForQuery{
				Pkey:             entity.Int(e.MChatRoomBelongingsPkey),
				MemberOnChatRoom: convMemberOnChatRoom(query.GetMembersOnChatRoomRow(e)),
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.MemberOnChatRoomForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MemberOnChatRoomDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MemberOnChatRoomNameCursorKey:
			return entity.Int(e.Pkey), e.MemberOnChatRoom.Member.Name
		case parameter.MemberOnChatRoomAddedAtCursorKey:
			return entity.Int(e.Pkey), e.MemberOnChatRoom.AddedAt
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
		return store.ListResult[entity.MemberOnChatRoom]{}, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	return res, nil
}

// GetMembersOnChatRoom チャットルーム上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoom(
	ctx context.Context, chatRoomID uuid.UUID, where parameter.WhereMemberOnChatRoomParam,
	order parameter.MemberOnChatRoomOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	return getMembersOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetMembersOnChatRoomWithSd SD付きでチャットルーム上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
	where parameter.WhereMemberOnChatRoomParam, order parameter.MemberOnChatRoomOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getMembersOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

func getPluralChatRoomsOnMember(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
	orderMethod parameter.ChatRoomOnMemberOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	var e []query.GetPluralChatRoomsOnMemberRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomsOnMember(ctx, query.GetPluralChatRoomsOnMemberParams{
			MemberIds:   memberIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralChatRoomsOnMemberUseNumberedPaginateRow
		qe, err = qtx.GetPluralChatRoomsOnMemberUseNumberedPaginate(
			ctx, query.GetPluralChatRoomsOnMemberUseNumberedPaginateParams{
				MemberIds:   memberIDs,
				Limit:       int32(np.Limit.Int64),
				Offset:      int32(np.Offset.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralChatRoomsOnMemberRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralChatRoomsOnMemberRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomOnMember]{},
			fmt.Errorf("failed to get chat rooms on member: %w", err)
	}
	entities := make([]entity.ChatRoomOnMember, len(e))
	for i, v := range e {
		entities[i] = convChatRoomOnMember(query.GetChatRoomsOnMemberRow(v))
	}
	return store.ListResult[entity.ChatRoomOnMember]{Data: entities}, nil
}

// GetPluralChatRoomsOnMember メンバー上の複数のチャットルームを取得する。
func (a *PgAdapter) GetPluralChatRoomsOnMember(
	ctx context.Context, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.ChatRoomOnMemberOrderMethod,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	return getPluralChatRoomsOnMember(ctx, a.query, memberIDs, order, np)
}

// GetPluralChatRoomsOnMemberWithSd SD付きでメンバー上の複数のチャットルームを取得する。
func (a *PgAdapter) GetPluralChatRoomsOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.ChatRoomOnMemberOrderMethod,
) (store.ListResult[entity.ChatRoomOnMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomOnMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomsOnMember(ctx, qtx, memberIDs, order, np)
}

func getPluralMembersOnChatRoom(
	ctx context.Context, qtx *query.Queries, chatRoomIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnChatRoomOrderMethod,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	var e []query.GetPluralMembersOnChatRoomRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersOnChatRoom(ctx, query.GetPluralMembersOnChatRoomParams{
			ChatRoomIds: chatRoomIDs,
			OrderMethod: order.GetStringValue(),
		})
	} else {
		var qe []query.GetPluralMembersOnChatRoomUseNumberedPaginateRow
		qe, err = qtx.GetPluralMembersOnChatRoomUseNumberedPaginate(
			ctx, query.GetPluralMembersOnChatRoomUseNumberedPaginateParams{
				ChatRoomIds: chatRoomIDs,
				Limit:       int32(np.Limit.Int64),
				Offset:      int32(np.Offset.Int64),
				OrderMethod: order.GetStringValue(),
			})
		e = make([]query.GetPluralMembersOnChatRoomRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralMembersOnChatRoomRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberOnChatRoom]{},
			fmt.Errorf("failed to get members on chat room: %w", err)
	}
	entities := make([]entity.MemberOnChatRoom, len(e))
	for i, v := range e {
		entities[i] = convMemberOnChatRoom(query.GetMembersOnChatRoomRow(v))
	}
	return store.ListResult[entity.MemberOnChatRoom]{Data: entities}, nil
}

// GetPluralMembersOnChatRoom チャットルーム上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoom(
	ctx context.Context, chatRoomIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnChatRoomOrderMethod,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	return getPluralMembersOnChatRoom(ctx, a.query, chatRoomIDs, np, order)
}

// GetPluralMembersOnChatRoomWithSd SD付きでチャットルーム上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomIDs []uuid.UUID,
	np store.NumberedPaginationParam, order parameter.MemberOnChatRoomOrderMethod,
) (store.ListResult[entity.MemberOnChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersOnChatRoom(ctx, qtx, chatRoomIDs, np, order)
}

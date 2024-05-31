package pgadapter

import (
	"context"
	"errors"
	"fmt"

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

func convMembersOnChatRoomAddMemberAction(
	r query.GetMembersOnChatRoomAddMemberActionRow,
) entity.MemberOnChatRoomAddMemberActionForQuery {
	var member entity.NullableEntity[entity.SimpleMember]
	if r.MemberID.Valid {
		member = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       r.MemberID.Bytes,
				Name:           r.MemberName.String,
				FirstName:      entity.String(r.MemberFirstName),
				LastName:       entity.String(r.MemberLastName),
				Email:          r.MemberEmail.String,
				ProfileImageID: entity.UUID(r.MemberProfileImageID),
				GradeID:        r.MemberGradeID.Bytes,
				GroupID:        r.MemberGroupID.Bytes,
			},
		}
	}
	return entity.MemberOnChatRoomAddMemberActionForQuery{
		Pkey: entity.Int(r.TChatRoomAddedMembersPkey),
		MemberOnChatRoomAddMemberAction: entity.MemberOnChatRoomAddMemberAction{
			ChatRoomAddMemberActionID: r.ChatRoomAddMemberActionID,
			Member:                    member,
		},
	}
}

func countMembersOnChatRoomAddMemberAction(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomAddMemberActionID uuid.UUID,
	_ parameter.WhereMemberOnChatRoomAddMemberActionParam,
) (int64, error) {
	c, err := qtx.CountMembersOnChatRoomAddMemberAction(ctx, chatRoomAddMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to count members on chat room add member action: %w", err)
	}
	return c, nil
}

// CountMembersOnChatRoomAddMemberAction メンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountMembersOnChatRoomAddMemberAction(
	ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomAddMemberActionParam,
) (int64, error) {
	return countMembersOnChatRoomAddMemberAction(ctx, a.query, chatRoomAddMemberActionID, where)
}

// CountMembersOnChatRoomAddMemberActionWithSd SD付きでメンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountMembersOnChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomAddMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMembersOnChatRoomAddMemberAction(ctx, qtx, chatRoomAddMemberActionID, where)
}

func addMemberToChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomAddedMemberParam,
) (entity.ChatRoomAddedMember, error) {
	p := query.CreateChatRoomAddedMemberParams{
		MemberID:                  pgtype.UUID(param.MemberID),
		ChatRoomAddMemberActionID: param.ChatRoomAddMemberActionID,
	}
	b, err := qtx.CreateChatRoomAddedMember(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomAddedMember{}, errhandle.NewModelNotFoundError("chat room add member action")
		}
		return entity.ChatRoomAddedMember{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
	}
	return entity.ChatRoomAddedMember{
		ChatRoomAddMemberActionID: b.ChatRoomAddMemberActionID,
		MemberID:                  entity.UUID(b.MemberID),
	}, nil
}

// AddMemberToChatRoomAddMemberAction チャットルーム追加メンバーアクションにメンバーを追加する。
func (a *PgAdapter) AddMemberToChatRoomAddMemberAction(
	ctx context.Context, param parameter.CreateChatRoomAddedMemberParam,
) (entity.ChatRoomAddedMember, error) {
	return addMemberToChatRoomAddMemberAction(ctx, a.query, param)
}

// AddMemberToChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクションにメンバーを追加する。
func (a *PgAdapter) AddMemberToChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomAddedMemberParam,
) (entity.ChatRoomAddedMember, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomAddedMember{}, store.ErrNotFoundDescriptor
	}
	return addMemberToChatRoomAddMemberAction(ctx, qtx, param)
}

func addMembersToChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomAddedMemberParam,
) (int64, error) {
	ps := make([]query.CreateChatRoomAddedMembersParams, len(params))
	for i, param := range params {
		ps[i] = query.CreateChatRoomAddedMembersParams{
			MemberID:                  pgtype.UUID(param.MemberID),
			ChatRoomAddMemberActionID: param.ChatRoomAddMemberActionID,
		}
	}
	c, err := qtx.CreateChatRoomAddedMembers(ctx, ps)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelNotFoundError("chat room add member action")
		}
		return 0, fmt.Errorf("failed to add members to chat room add member action: %w", err)
	}
	return c, nil
}

// AddMembersToChatRoomAddMemberAction チャットルーム追加メンバーアクションにメンバーを複数追加する。
func (a *PgAdapter) AddMembersToChatRoomAddMemberAction(
	ctx context.Context, params []parameter.CreateChatRoomAddedMemberParam,
) (int64, error) {
	return addMembersToChatRoomAddMemberAction(ctx, a.query, params)
}

// AddMembersToChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクションにメンバーを複数追加する。
func (a *PgAdapter) AddMembersToChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomAddedMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return addMembersToChatRoomAddMemberAction(ctx, qtx, params)
}

func deleteChatRoomAddedMember(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionID, memberID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddedMember(ctx, query.DeleteChatRoomAddedMemberParams{
		ChatRoomAddMemberActionID: chatRoomAddMemberActionID,
		MemberID: pgtype.UUID{
			Bytes: memberID,
			Valid: true,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room added member: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room added member")
	}
	return c, nil
}

// DeleteChatRoomAddedMember チャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMember(
	ctx context.Context, chatRoomAddMemberActionID, memberID uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddedMember(ctx, a.query, chatRoomAddMemberActionID, memberID)
}

// DeleteChatRoomAddedMemberWithSd SD付きでチャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMemberWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionID, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddedMember(ctx, qtx, chatRoomAddMemberActionID, memberID)
}

func deleteChatRoomAddedMembersOnChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddedMembersOnChatRoomAddMemberAction(ctx, chatRoomAddMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room added members on chat room add member action: %w", err)
	}
	return c, nil
}

// DeleteChatRoomAddedMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnChatRoomAddMemberAction(
	ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddedMembersOnChatRoomAddMemberAction(ctx, a.query, chatRoomAddMemberActionID)
}

// DeleteChatRoomAddedMemberOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMemberOnChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddedMembersOnChatRoomAddMemberAction(ctx, qtx, chatRoomAddMemberActionID)
}

func deleteChatRoomAddedMembersOnChatRoomAddMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddedMembersOnChatRoomAddMemberActions(ctx, chatRoomAddMemberActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room added members on chat room add member actions: %w", err)
	}
	return c, nil
}

// DeleteChatRoomAddedMembersOnChatRoomAddMemberActions チャットルーム追加メンバーアクション上の複数のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnChatRoomAddMemberActions(
	ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddedMembersOnChatRoomAddMemberActions(ctx, a.query, chatRoomAddMemberActionIDs)
}

// DeleteChatRoomAddedMembersOnChatRoomAddMemberActionsWithSd SD付きでチャットルーム追加メンバーアクション上の複数のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnChatRoomAddMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddedMembersOnChatRoomAddMemberActions(ctx, qtx, chatRoomAddMemberActionIDs)
}

func deleteChatRoomAddedMembersOnMember(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddedMembersOnMember(ctx, pgtype.UUID{Bytes: memberID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room added members on member: %w", err)
	}
	return c, nil
}

// DeleteChatRoomAddedMembersOnMember メンバー上のチャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnMember(
	ctx context.Context, memberID uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddedMembersOnMember(ctx, a.query, memberID)
}

// DeleteChatRoomAddedMembersOnMemberWithSd SD付きでメンバー上のチャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddedMembersOnMember(ctx, qtx, memberID)
}

func deleteChatRoomAddedMembersOnMembers(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddedMembersOnMembers(ctx, memberIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room added members on members: %w", err)
	}
	return c, nil
}

// DeleteChatRoomAddedMembersOnMembers メンバー上の複数のチャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnMembers(
	ctx context.Context, memberIDs []uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddedMembersOnMembers(ctx, a.query, memberIDs)
}

// DeleteChatRoomAddedMembersOnMembersWithSd SD付きでメンバー上の複数のチャットルーム追加メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomAddedMembersOnMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddedMembersOnMembers(ctx, qtx, memberIDs)
}

func getMembersOnChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionID uuid.UUID,
	_ parameter.WhereMemberOnChatRoomAddMemberActionParam,
	order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	eConvFunc := func(e entity.MemberOnChatRoomAddMemberActionForQuery) (entity.MemberOnChatRoomAddMemberAction, error) {
		return e.MemberOnChatRoomAddMemberAction, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembersOnChatRoomAddMemberAction(ctx, chatRoomAddMemberActionID)
		if err != nil {
			return 0, fmt.Errorf("failed to count members on chat room add member action: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.MemberOnChatRoomAddMemberActionForQuery, error) {
		r, err := qtx.GetMembersOnChatRoomAddMemberAction(ctx, chatRoomAddMemberActionID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberOnChatRoomAddMemberActionForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members on chat room add member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomAddMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomAddMemberAction(e)
		}
		return fq, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.MemberOnChatRoomAddMemberActionForQuery, error) {
		p := query.GetMembersOnChatRoomAddMemberActionUseKeysetPaginateParams{
			ChatRoomAddMemberActionID: chatRoomAddMemberActionID,
			CursorDirection:           cursorDir,
			Cursor:                    cursor,
			Limit:                     limit,
		}
		r, err := qtx.GetMembersOnChatRoomAddMemberActionUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room add member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomAddMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomAddMemberAction(query.GetMembersOnChatRoomAddMemberActionRow(e))
		}
		return fq, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.MemberOnChatRoomAddMemberActionForQuery, error) {
		p := query.GetMembersOnChatRoomAddMemberActionUseNumberedPaginateParams{
			ChatRoomAddMemberActionID: chatRoomAddMemberActionID,
			Limit:                     limit,
			Offset:                    offset,
		}
		r, err := qtx.GetMembersOnChatRoomAddMemberActionUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room add member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomAddMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomAddMemberAction(query.GetMembersOnChatRoomAddMemberActionRow(e))
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.MemberOnChatRoomAddMemberActionForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomOnMemberDefaultCursorKey:
			return entity.Int(e.Pkey), nil
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
		return store.ListResult[entity.MemberOnChatRoomAddMemberAction]{},
			fmt.Errorf("failed to get members on chat room add member action: %w", err)
	}
	return res, nil
}

// GetMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoomAddMemberAction(
	ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomAddMemberActionParam,
	order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	return getMembersOnChatRoomAddMemberAction(ctx, a.query, chatRoomAddMemberActionID, where, order, np, cp, wc)
}

// GetMembersOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomAddMemberActionParam,
	order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoomAddMemberAction]{}, store.ErrNotFoundDescriptor
	}
	return getMembersOnChatRoomAddMemberAction(ctx, qtx, chatRoomAddMemberActionID, where, order, np, cp, wc)
}

func getPluralMembersOnChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionIDs []uuid.UUID,
	_ parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	var e []query.GetPluralMembersOnChatRoomAddMemberActionRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersOnChatRoomAddMemberAction(ctx, chatRoomAddMemberActionIDs)
	} else {
		var qe []query.GetPluralMembersOnChatRoomAddMemberActionUseNumberedPaginateRow
		qe, err = qtx.GetPluralMembersOnChatRoomAddMemberActionUseNumberedPaginate(
			ctx, query.GetPluralMembersOnChatRoomAddMemberActionUseNumberedPaginateParams{
				Limit:                      int32(np.Limit.Int64),
				Offset:                     int32(np.Offset.Int64),
				ChatRoomAddMemberActionIds: chatRoomAddMemberActionIDs,
			})
		e = make([]query.GetPluralMembersOnChatRoomAddMemberActionRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralMembersOnChatRoomAddMemberActionRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberOnChatRoomAddMemberAction]{},
			fmt.Errorf("failed to get plural members on chat room add member action: %w", err)
	}
	entities := make([]entity.MemberOnChatRoomAddMemberAction, len(e))
	for i, v := range e {
		entities[i] = convMembersOnChatRoomAddMemberAction(
			query.GetMembersOnChatRoomAddMemberActionRow(v)).MemberOnChatRoomAddMemberAction
	}
	return store.ListResult[entity.MemberOnChatRoomAddMemberAction]{Data: entities}, nil
}

// GetPluralMembersOnChatRoomAddMemberAction チャットルーム追加メンバーアクション上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoomAddMemberAction(
	ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID,
	order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	return getPluralMembersOnChatRoomAddMemberAction(ctx, a.query, chatRoomAddMemberActionIDs, order, np)
}

// GetPluralMembersOnChatRoomAddMemberActionWithSd SD付きでチャットルーム追加メンバーアクション上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionIDs []uuid.UUID,
	order parameter.MemberOnChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomAddMemberAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoomAddMemberAction]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersOnChatRoomAddMemberAction(ctx, qtx, chatRoomAddMemberActionIDs, order, np)
}

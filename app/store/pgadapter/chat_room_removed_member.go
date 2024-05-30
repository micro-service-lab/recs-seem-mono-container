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

func convMembersOnChatRoomRemoveMemberAction(
	r query.GetMembersOnChatRoomRemoveMemberActionRow,
) entity.MemberOnChatRoomRemoveMemberActionForQuery {
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
			},
		}
	}
	return entity.MemberOnChatRoomRemoveMemberActionForQuery{
		Pkey: entity.Int(r.TChatRoomRemovedMembersPkey),
		MemberOnChatRoomRemoveMemberAction: entity.MemberOnChatRoomRemoveMemberAction{
			ChatRoomRemoveMemberActionID: r.ChatRoomRemoveMemberActionID,
			Member:                       member,
		},
	}
}

func countMembersOnChatRoomRemoveMemberAction(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomRemoveMemberActionID uuid.UUID,
	_ parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
) (int64, error) {
	c, err := qtx.CountMembersOnChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to count members on chat room remove member action: %w", err)
	}
	return c, nil
}

// CountMembersOnChatRoomRemoveMemberAction メンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
) (int64, error) {
	return countMembersOnChatRoomRemoveMemberAction(ctx, a.query, chatRoomRemoveMemberActionID, where)
}

// CountMembersOnChatRoomRemoveMemberActionWithSd SD付きでメンバー上のチャットルーム数を取得する。
func (a *PgAdapter) CountMembersOnChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMembersOnChatRoomRemoveMemberAction(ctx, qtx, chatRoomRemoveMemberActionID, where)
}

func removeMemberToChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomRemovedMemberParam,
) (entity.ChatRoomRemovedMember, error) {
	p := query.CreateChatRoomRemovedMemberParams{
		MemberID:                     pgtype.UUID(param.MemberID),
		ChatRoomRemoveMemberActionID: param.ChatRoomRemoveMemberActionID,
	}
	b, err := qtx.CreateChatRoomRemovedMember(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomRemovedMember{}, errhandle.NewModelNotFoundError("chat room remove member action")
		}
		return entity.ChatRoomRemovedMember{},
			fmt.Errorf("failed to remove member to chat room remove member action: %w", err)
	}
	return entity.ChatRoomRemovedMember{
		ChatRoomRemoveMemberActionID: b.ChatRoomRemoveMemberActionID,
		MemberID:                     entity.UUID(b.MemberID),
	}, nil
}

// RemoveMemberToChatRoomRemoveMemberAction チャットルーム追放メンバーアクションにメンバーを追放する。
func (a *PgAdapter) RemoveMemberToChatRoomRemoveMemberAction(
	ctx context.Context, param parameter.CreateChatRoomRemovedMemberParam,
) (entity.ChatRoomRemovedMember, error) {
	return removeMemberToChatRoomRemoveMemberAction(ctx, a.query, param)
}

// RemoveMemberToChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクションにメンバーを追放する。
func (a *PgAdapter) RemoveMemberToChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomRemovedMemberParam,
) (entity.ChatRoomRemovedMember, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomRemovedMember{}, store.ErrNotFoundDescriptor
	}
	return removeMemberToChatRoomRemoveMemberAction(ctx, qtx, param)
}

func removeMembersToChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomRemovedMemberParam,
) (int64, error) {
	ps := make([]query.CreateChatRoomRemovedMembersParams, len(params))
	for i, param := range params {
		ps[i] = query.CreateChatRoomRemovedMembersParams{
			MemberID:                     pgtype.UUID(param.MemberID),
			ChatRoomRemoveMemberActionID: param.ChatRoomRemoveMemberActionID,
		}
	}
	c, err := qtx.CreateChatRoomRemovedMembers(ctx, ps)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelNotFoundError("chat room remove member action")
		}
		return 0, fmt.Errorf("failed to remove members to chat room remove member action: %w", err)
	}
	return c, nil
}

// RemoveMembersToChatRoomRemoveMemberAction チャットルーム追放メンバーアクションにメンバーを複数追放する。
func (a *PgAdapter) RemoveMembersToChatRoomRemoveMemberAction(
	ctx context.Context, params []parameter.CreateChatRoomRemovedMemberParam,
) (int64, error) {
	return removeMembersToChatRoomRemoveMemberAction(ctx, a.query, params)
}

// RemoveMembersToChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクションにメンバーを複数追放する。
func (a *PgAdapter) RemoveMembersToChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomRemovedMemberParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return removeMembersToChatRoomRemoveMemberAction(ctx, qtx, params)
}

func deleteChatRoomRemovedMember(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionID, memberID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemovedMember(ctx, query.DeleteChatRoomRemovedMemberParams{
		ChatRoomRemoveMemberActionID: chatRoomRemoveMemberActionID,
		MemberID: pgtype.UUID{
			Bytes: memberID,
			Valid: true,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room removed member: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room removed member")
	}
	return c, nil
}

// DeleteChatRoomRemovedMember チャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMember(
	ctx context.Context, chatRoomRemoveMemberActionID, memberID uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemovedMember(ctx, a.query, chatRoomRemoveMemberActionID, memberID)
}

// DeleteChatRoomRemovedMemberWithSd SD付きでチャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMemberWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionID, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemovedMember(ctx, qtx, chatRoomRemoveMemberActionID, memberID)
}

func deleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room removed members on chat room remove member action: %w", err)
	}
	return c, nil
}

// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(ctx, a.query, chatRoomRemoveMemberActionID)
}

// DeleteChatRoomRemovedMemberOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMemberOnChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(ctx, qtx, chatRoomRemoveMemberActionID)
}

func deleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(ctx, chatRoomRemoveMemberActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room removed members on chat room remove member actions: %w", err)
	}
	return c, nil
}

// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions チャットルーム追放メンバーアクション上の複数のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(
	ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(ctx, a.query, chatRoomRemoveMemberActionIDs)
}

// DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActionsWithSd SD付きでチャットルーム追放メンバーアクション上の複数のメンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(ctx, qtx, chatRoomRemoveMemberActionIDs)
}

func deleteChatRoomRemovedMembersOnMember(
	ctx context.Context, qtx *query.Queries, memberID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemovedMembersOnMember(ctx, pgtype.UUID{Bytes: memberID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room removed members on member: %w", err)
	}
	return c, nil
}

// DeleteChatRoomRemovedMembersOnMember メンバー上のチャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnMember(
	ctx context.Context, memberID uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemovedMembersOnMember(ctx, a.query, memberID)
}

// DeleteChatRoomRemovedMembersOnMemberWithSd SD付きでメンバー上のチャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnMemberWithSd(
	ctx context.Context, sd store.Sd, memberID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemovedMembersOnMember(ctx, qtx, memberID)
}

func deleteChatRoomRemovedMembersOnMembers(
	ctx context.Context, qtx *query.Queries, memberIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemovedMembersOnMembers(ctx, memberIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room removed members on members: %w", err)
	}
	return c, nil
}

// DeleteChatRoomRemovedMembersOnMembers メンバー上の複数のチャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnMembers(
	ctx context.Context, memberIDs []uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemovedMembersOnMembers(ctx, a.query, memberIDs)
}

// DeleteChatRoomRemovedMembersOnMembersWithSd SD付きでメンバー上の複数のチャットルーム追放メンバーを削除する。
func (a *PgAdapter) DeleteChatRoomRemovedMembersOnMembersWithSd(
	ctx context.Context, sd store.Sd, memberIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemovedMembersOnMembers(ctx, qtx, memberIDs)
}

func getMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionID uuid.UUID,
	_ parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
	order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	eConvFunc := func(
		e entity.MemberOnChatRoomRemoveMemberActionForQuery,
	) (entity.MemberOnChatRoomRemoveMemberAction, error) {
		return e.MemberOnChatRoomRemoveMemberAction, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMembersOnChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionID)
		if err != nil {
			return 0, fmt.Errorf("failed to count members on chat room remove member action: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.MemberOnChatRoomRemoveMemberActionForQuery, error) {
		r, err := qtx.GetMembersOnChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MemberOnChatRoomRemoveMemberActionForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get members on chat room remove member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomRemoveMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomRemoveMemberAction(e)
		}
		return fq, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.MemberOnChatRoomRemoveMemberActionForQuery, error) {
		p := query.GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateParams{
			ChatRoomRemoveMemberActionID: chatRoomRemoveMemberActionID,
			CursorDirection:              cursorDir,
			Cursor:                       cursor,
			Limit:                        limit,
		}
		r, err := qtx.GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room remove member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomRemoveMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomRemoveMemberAction(query.GetMembersOnChatRoomRemoveMemberActionRow(e))
		}
		return fq, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.MemberOnChatRoomRemoveMemberActionForQuery, error) {
		p := query.GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams{
			ChatRoomRemoveMemberActionID: chatRoomRemoveMemberActionID,
			Limit:                        limit,
			Offset:                       offset,
		}
		r, err := qtx.GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get members on chat room remove member action: %w", err)
		}
		fq := make([]entity.MemberOnChatRoomRemoveMemberActionForQuery, len(r))
		for i, e := range r {
			fq[i] = convMembersOnChatRoomRemoveMemberAction(query.GetMembersOnChatRoomRemoveMemberActionRow(e))
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.MemberOnChatRoomRemoveMemberActionForQuery) (entity.Int, any) {
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
		return store.ListResult[entity.MemberOnChatRoomRemoveMemberAction]{},
			fmt.Errorf("failed to get members on chat room remove member action: %w", err)
	}
	return res, nil
}

// GetMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
	order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	return getMembersOnChatRoomRemoveMemberAction(ctx, a.query, chatRoomRemoveMemberActionID, where, order, np, cp, wc)
}

// GetMembersOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上のメンバーを取得する。
func (a *PgAdapter) GetMembersOnChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionID uuid.UUID,
	where parameter.WhereMemberOnChatRoomRemoveMemberActionParam,
	order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoomRemoveMemberAction]{}, store.ErrNotFoundDescriptor
	}
	return getMembersOnChatRoomRemoveMemberAction(ctx, qtx, chatRoomRemoveMemberActionID, where, order, np, cp, wc)
}

func getPluralMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionIDs []uuid.UUID,
	_ parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	var e []query.GetPluralMembersOnChatRoomRemoveMemberActionRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMembersOnChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionIDs)
	} else {
		var qe []query.GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow
		qe, err = qtx.GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate(
			ctx, query.GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams{
				Limit:                         int32(np.Limit.Int64),
				Offset:                        int32(np.Offset.Int64),
				ChatRoomRemoveMemberActionIds: chatRoomRemoveMemberActionIDs,
			})
		e = make([]query.GetPluralMembersOnChatRoomRemoveMemberActionRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralMembersOnChatRoomRemoveMemberActionRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MemberOnChatRoomRemoveMemberAction]{},
			fmt.Errorf("failed to get plural members on chat room remove member action: %w", err)
	}
	entities := make([]entity.MemberOnChatRoomRemoveMemberAction, len(e))
	for i, v := range e {
		entities[i] = convMembersOnChatRoomRemoveMemberAction(
			query.GetMembersOnChatRoomRemoveMemberActionRow(v)).MemberOnChatRoomRemoveMemberAction
	}
	return store.ListResult[entity.MemberOnChatRoomRemoveMemberAction]{Data: entities}, nil
}

// GetPluralMembersOnChatRoomRemoveMemberAction チャットルーム追放メンバーアクション上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoomRemoveMemberAction(
	ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID,
	order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	return getPluralMembersOnChatRoomRemoveMemberAction(ctx, a.query, chatRoomRemoveMemberActionIDs, order, np)
}

// GetPluralMembersOnChatRoomRemoveMemberActionWithSd SD付きでチャットルーム追放メンバーアクション上の複数のメンバーを取得する。
func (a *PgAdapter) GetPluralMembersOnChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionIDs []uuid.UUID,
	order parameter.MemberOnChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.MemberOnChatRoomRemoveMemberAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MemberOnChatRoomRemoveMemberAction]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMembersOnChatRoomRemoveMemberAction(ctx, qtx, chatRoomRemoveMemberActionIDs, order, np)
}

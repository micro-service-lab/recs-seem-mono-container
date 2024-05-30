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

func convChatRoomWithdrawActionOnChatRoom(
	e query.GetChatRoomWithdrawActionsOnChatRoomRow,
) entity.ChatRoomWithdrawActionWithMemberForQuery {
	var member entity.NullableEntity[entity.SimpleMember]
	if e.MemberID.Valid {
		member = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.MemberID.Bytes,
				Name:           e.WithdrawMemberName.String,
				Email:          e.WithdrawMemberEmail.String,
				FirstName:      entity.String(e.WithdrawMemberFirstName),
				LastName:       entity.String(e.WithdrawMemberLastName),
				ProfileImageID: entity.UUID(e.WithdrawMemberProfileImageID),
			},
		}
	}
	return entity.ChatRoomWithdrawActionWithMemberForQuery{
		Pkey: entity.Int(e.TChatRoomWithdrawActionsPkey),
		ChatRoomWithdrawActionWithMember: entity.ChatRoomWithdrawActionWithMember{
			ChatRoomWithdrawActionID: e.ChatRoomWithdrawActionID,
			ChatRoomActionID:         e.ChatRoomActionID,
			Member:                   member,
		},
	}
}

// countChatRoomWithdrawActions はチャットルームメンバー脱退アクション数を取得する内部関数です。
func countChatRoomWithdrawActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomWithdrawActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomWithdrawActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room withdraw actions: %w", err)
	}
	return c, nil
}

// CountChatRoomWithdrawActions はチャットルームメンバー脱退アクション数を取得します。
func (a *PgAdapter) CountChatRoomWithdrawActions(
	ctx context.Context, where parameter.WhereChatRoomWithdrawActionParam,
) (int64, error) {
	return countChatRoomWithdrawActions(ctx, a.query, where)
}

// CountChatRoomWithdrawActionsWithSd はSD付きでチャットルームメンバー脱退アクション数を取得します。
func (a *PgAdapter) CountChatRoomWithdrawActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomWithdrawActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomWithdrawActions(ctx, qtx, where)
}

// createChatRoomWithdrawAction はチャットルームメンバー脱退アクションを作成する内部関数です。
func createChatRoomWithdrawAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomWithdrawActionParam,
) (entity.ChatRoomWithdrawAction, error) {
	e, err := qtx.CreateChatRoomWithdrawAction(ctx, query.CreateChatRoomWithdrawActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		MemberID:         pgtype.UUID(param.MemberID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomWithdrawAction{}, errhandle.NewModelDuplicatedError("chat room withdraw action")
		}
		return entity.ChatRoomWithdrawAction{}, fmt.Errorf("failed to create chat room withdraw action: %w", err)
	}
	entity := entity.ChatRoomWithdrawAction{
		ChatRoomWithdrawActionID: e.ChatRoomWithdrawActionID,
		ChatRoomActionID:         e.ChatRoomActionID,
		MemberID:                 entity.UUID(e.MemberID),
	}
	return entity, nil
}

// CreateChatRoomWithdrawAction はチャットルームメンバー脱退アクションを作成します。
func (a *PgAdapter) CreateChatRoomWithdrawAction(
	ctx context.Context, param parameter.CreateChatRoomWithdrawActionParam,
) (entity.ChatRoomWithdrawAction, error) {
	return createChatRoomWithdrawAction(ctx, a.query, param)
}

// CreateChatRoomWithdrawActionWithSd はSD付きでチャットルームメンバー脱退アクションを作成します。
func (a *PgAdapter) CreateChatRoomWithdrawActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomWithdrawActionParam,
) (entity.ChatRoomWithdrawAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomWithdrawAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomWithdrawAction(ctx, qtx, param)
}

// createChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを作成する内部関数です。
func createChatRoomWithdrawActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomWithdrawActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomWithdrawActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomWithdrawActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			MemberID:         pgtype.UUID(p.MemberID),
		}
	}
	n, err := qtx.CreateChatRoomWithdrawActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room withdraw action")
		}
		return 0, fmt.Errorf("failed to create chat room withdraw actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを作成します。
func (a *PgAdapter) CreateChatRoomWithdrawActions(
	ctx context.Context, params []parameter.CreateChatRoomWithdrawActionParam,
) (int64, error) {
	return createChatRoomWithdrawActions(ctx, a.query, params)
}

// CreateChatRoomWithdrawActionsWithSd はSD付きで複数のチャットルームメンバー脱退アクションを作成します。
func (a *PgAdapter) CreateChatRoomWithdrawActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomWithdrawActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomWithdrawActions(ctx, qtx, params)
}

// deleteChatRoomWithdrawAction はチャットルームメンバー脱退アクションを削除する内部関数です。
func deleteChatRoomWithdrawAction(
	ctx context.Context, qtx *query.Queries, chatRoomWithdrawActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomWithdrawAction(ctx, chatRoomWithdrawActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room withdraw action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room withdraw action")
	}
	return c, nil
}

// DeleteChatRoomWithdrawAction はチャットルームメンバー脱退アクションを削除します。
func (a *PgAdapter) DeleteChatRoomWithdrawAction(
	ctx context.Context, chatRoomWithdrawActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomWithdrawAction(ctx, a.query, chatRoomWithdrawActionID)
}

// DeleteChatRoomWithdrawActionWithSd はSD付きでチャットルームメンバー脱退アクションを削除します。
func (a *PgAdapter) DeleteChatRoomWithdrawActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomWithdrawActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomWithdrawAction(ctx, qtx, chatRoomWithdrawActionID)
}

// pluralDeleteChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを削除する内部関数です。
func pluralDeleteChatRoomWithdrawActions(
	ctx context.Context, qtx *query.Queries, chatRoomWithdrawActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomWithdrawActions(ctx, chatRoomWithdrawActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room withdraw actions: %w", err)
	}
	if c != int64(len(chatRoomWithdrawActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room withdraw action")
	}
	return c, nil
}

// PluralDeleteChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomWithdrawActions(
	ctx context.Context, chatRoomWithdrawActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomWithdrawActions(ctx, a.query, chatRoomWithdrawActionIDs)
}

// PluralDeleteChatRoomWithdrawActionsWithSd はSD付きで複数のチャットルームメンバー脱退アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomWithdrawActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomWithdrawActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomWithdrawActions(ctx, qtx, chatRoomWithdrawActionIDs)
}

// getChatRoomWithdrawActions はチャットルームメンバー脱退アクションを取得する内部関数です。
func getChatRoomWithdrawActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomWithdrawActionParam,
	order parameter.ChatRoomWithdrawActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	eConvFunc := func(
		e entity.ChatRoomWithdrawActionWithMemberForQuery,
	) (entity.ChatRoomWithdrawActionWithMember, error) {
		return e.ChatRoomWithdrawActionWithMember, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomWithdrawActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room withdraw actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomWithdrawActionWithMemberForQuery, error) {
		r, err := qtx.GetChatRoomWithdrawActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomWithdrawActionWithMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room withdraw actions: %w", err)
		}
		e := make([]entity.ChatRoomWithdrawActionWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomWithdrawActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomWithdrawActionWithMemberForQuery, error) {
		p := query.GetChatRoomWithdrawActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomWithdrawActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room withdraw actions: %w", err)
		}
		e := make([]entity.ChatRoomWithdrawActionWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomWithdrawActionOnChatRoom(query.GetChatRoomWithdrawActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomWithdrawActionWithMemberForQuery, error) {
		p := query.GetChatRoomWithdrawActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomWithdrawActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room withdraw actions: %w", err)
		}
		e := make([]entity.ChatRoomWithdrawActionWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomWithdrawActionOnChatRoom(query.GetChatRoomWithdrawActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomWithdrawActionWithMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomWithdrawActionDefaultCursorKey:
			return e.Pkey, nil
		}
		return e.Pkey, nil
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
		return store.ListResult[entity.ChatRoomWithdrawActionWithMember]{},
			fmt.Errorf("failed to get chat room withdraw actions: %w", err)
	}
	return res, nil
}

// GetChatRoomWithdrawActionsOnChatRoom はチャットルームメンバー脱退アクションを取得します。
func (a *PgAdapter) GetChatRoomWithdrawActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomWithdrawActionParam,
	order parameter.ChatRoomWithdrawActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	return getChatRoomWithdrawActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomWithdrawActionsOnChatRoomWithSd はSD付きでチャットルームメンバー脱退アクションを取得します。
func (a *PgAdapter) GetChatRoomWithdrawActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomWithdrawActionParam,
	order parameter.ChatRoomWithdrawActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomWithdrawActionWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomWithdrawActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを取得する内部関数です。
func getPluralChatRoomWithdrawActions(
	ctx context.Context, qtx *query.Queries, chatRoomWithdrawActionIDs []uuid.UUID,
	_ parameter.ChatRoomWithdrawActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	var e []query.GetPluralChatRoomWithdrawActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomWithdrawActions(ctx, chatRoomWithdrawActionIDs)
	} else {
		var ne []query.GetPluralChatRoomWithdrawActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomWithdrawActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomWithdrawActionsUseNumberedPaginateParams{
				Limit:                     int32(np.Limit.Int64),
				Offset:                    int32(np.Offset.Int64),
				ChatRoomWithdrawActionIds: chatRoomWithdrawActionIDs,
			})
		e = make([]query.GetPluralChatRoomWithdrawActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomWithdrawActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomWithdrawActionWithMember]{},
			fmt.Errorf("failed to get chat room withdraw actions: %w", err)
	}
	entities := make([]entity.ChatRoomWithdrawActionWithMember, len(e))
	for i, v := range e {
		entities[i] = convChatRoomWithdrawActionOnChatRoom(
			query.GetChatRoomWithdrawActionsOnChatRoomRow(v)).ChatRoomWithdrawActionWithMember
	}
	return store.ListResult[entity.ChatRoomWithdrawActionWithMember]{Data: entities}, nil
}

// GetPluralChatRoomWithdrawActions は複数のチャットルームメンバー脱退アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomWithdrawActions(
	ctx context.Context, chatRoomWithdrawActionIDs []uuid.UUID,
	order parameter.ChatRoomWithdrawActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	return getPluralChatRoomWithdrawActions(ctx, a.query, chatRoomWithdrawActionIDs, order, np)
}

// GetPluralChatRoomWithdrawActionsWithSd はSD付きで複数のチャットルームメンバー脱退アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomWithdrawActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomWithdrawActionIDs []uuid.UUID,
	order parameter.ChatRoomWithdrawActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomWithdrawActionWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomWithdrawActionWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomWithdrawActions(ctx, qtx, chatRoomWithdrawActionIDs, order, np)
}

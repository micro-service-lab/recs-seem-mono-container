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

func convChatRoomRemoveMemberActionOnChatRoom(
	e query.GetChatRoomRemoveMemberActionsOnChatRoomRow,
) entity.ChatRoomRemoveMemberActionWithRemovedByForQuery {
	var removedBy entity.NullableEntity[entity.SimpleMember]
	if e.RemovedBy.Valid {
		removedBy = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.RemovedBy.Bytes,
				Name:           e.RemoveMemberName.String,
				Email:          e.RemoveMemberEmail.String,
				FirstName:      entity.String(e.RemoveMemberFirstName),
				LastName:       entity.String(e.RemoveMemberLastName),
				ProfileImageID: entity.UUID(e.RemoveMemberProfileImageID),
				GradeID:        e.RemoveMemberGradeID.Bytes,
				GroupID:        e.RemoveMemberGroupID.Bytes,
			},
		}
	}
	return entity.ChatRoomRemoveMemberActionWithRemovedByForQuery{
		Pkey: entity.Int(e.TChatRoomRemoveMemberActionsPkey),
		ChatRoomRemoveMemberActionWithRemovedBy: entity.ChatRoomRemoveMemberActionWithRemovedBy{
			ChatRoomRemoveMemberActionID: e.ChatRoomRemoveMemberActionID,
			ChatRoomActionID:             e.ChatRoomActionID,
			RemovedBy:                    removedBy,
		},
	}
}

// countChatRoomRemoveMemberActions はチャットルームメンバー追放アクション数を取得する内部関数です。
func countChatRoomRemoveMemberActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomRemoveMemberActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomRemoveMemberActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room remove member actions: %w", err)
	}
	return c, nil
}

// CountChatRoomRemoveMemberActions はチャットルームメンバー追放アクション数を取得します。
func (a *PgAdapter) CountChatRoomRemoveMemberActions(
	ctx context.Context, where parameter.WhereChatRoomRemoveMemberActionParam,
) (int64, error) {
	return countChatRoomRemoveMemberActions(ctx, a.query, where)
}

// CountChatRoomRemoveMemberActionsWithSd はSD付きでチャットルームメンバー追放アクション数を取得します。
func (a *PgAdapter) CountChatRoomRemoveMemberActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomRemoveMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomRemoveMemberActions(ctx, qtx, where)
}

// createChatRoomRemoveMemberAction はチャットルームメンバー追放アクションを作成する内部関数です。
func createChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomRemoveMemberActionParam,
) (entity.ChatRoomRemoveMemberAction, error) {
	e, err := qtx.CreateChatRoomRemoveMemberAction(ctx, query.CreateChatRoomRemoveMemberActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		RemovedBy:        pgtype.UUID(param.RemovedBy),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomRemoveMemberAction{}, errhandle.NewModelDuplicatedError("chat room remove member action")
		}
		return entity.ChatRoomRemoveMemberAction{}, fmt.Errorf("failed to create chat room remove member action: %w", err)
	}
	entity := entity.ChatRoomRemoveMemberAction{
		ChatRoomRemoveMemberActionID: e.ChatRoomRemoveMemberActionID,
		ChatRoomActionID:             e.ChatRoomActionID,
		RemovedBy:                    entity.UUID(e.RemovedBy),
	}
	return entity, nil
}

// CreateChatRoomRemoveMemberAction はチャットルームメンバー追放アクションを作成します。
func (a *PgAdapter) CreateChatRoomRemoveMemberAction(
	ctx context.Context, param parameter.CreateChatRoomRemoveMemberActionParam,
) (entity.ChatRoomRemoveMemberAction, error) {
	return createChatRoomRemoveMemberAction(ctx, a.query, param)
}

// CreateChatRoomRemoveMemberActionWithSd はSD付きでチャットルームメンバー追放アクションを作成します。
func (a *PgAdapter) CreateChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomRemoveMemberActionParam,
) (entity.ChatRoomRemoveMemberAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomRemoveMemberAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomRemoveMemberAction(ctx, qtx, param)
}

// createChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを作成する内部関数です。
func createChatRoomRemoveMemberActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomRemoveMemberActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomRemoveMemberActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomRemoveMemberActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			RemovedBy:        pgtype.UUID(p.RemovedBy),
		}
	}
	n, err := qtx.CreateChatRoomRemoveMemberActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room remove member action")
		}
		return 0, fmt.Errorf("failed to create chat room remove member actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを作成します。
func (a *PgAdapter) CreateChatRoomRemoveMemberActions(
	ctx context.Context, params []parameter.CreateChatRoomRemoveMemberActionParam,
) (int64, error) {
	return createChatRoomRemoveMemberActions(ctx, a.query, params)
}

// CreateChatRoomRemoveMemberActionsWithSd はSD付きで複数のチャットルームメンバー追放アクションを作成します。
func (a *PgAdapter) CreateChatRoomRemoveMemberActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomRemoveMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomRemoveMemberActions(ctx, qtx, params)
}

// deleteChatRoomRemoveMemberAction はチャットルームメンバー追放アクションを削除する内部関数です。
func deleteChatRoomRemoveMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomRemoveMemberAction(ctx, chatRoomRemoveMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room remove member action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room remove member action")
	}
	return c, nil
}

// DeleteChatRoomRemoveMemberAction はチャットルームメンバー追放アクションを削除します。
func (a *PgAdapter) DeleteChatRoomRemoveMemberAction(
	ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomRemoveMemberAction(ctx, a.query, chatRoomRemoveMemberActionID)
}

// DeleteChatRoomRemoveMemberActionWithSd はSD付きでチャットルームメンバー追放アクションを削除します。
func (a *PgAdapter) DeleteChatRoomRemoveMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomRemoveMemberAction(ctx, qtx, chatRoomRemoveMemberActionID)
}

// pluralDeleteChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを削除する内部関数です。
func pluralDeleteChatRoomRemoveMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomRemoveMemberActions(ctx, chatRoomRemoveMemberActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room remove member actions: %w", err)
	}
	if c != int64(len(chatRoomRemoveMemberActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room remove member action")
	}
	return c, nil
}

// PluralDeleteChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomRemoveMemberActions(
	ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomRemoveMemberActions(ctx, a.query, chatRoomRemoveMemberActionIDs)
}

// PluralDeleteChatRoomRemoveMemberActionsWithSd はSD付きで複数のチャットルームメンバー追放アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomRemoveMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomRemoveMemberActions(ctx, qtx, chatRoomRemoveMemberActionIDs)
}

// getChatRoomRemoveMemberActions はチャットルームメンバー追放アクションを取得する内部関数です。
func getChatRoomRemoveMemberActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomRemoveMemberActionParam,
	order parameter.ChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	eConvFunc := func(
		e entity.ChatRoomRemoveMemberActionWithRemovedByForQuery,
	) (entity.ChatRoomRemoveMemberActionWithRemovedBy, error) {
		return e.ChatRoomRemoveMemberActionWithRemovedBy, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomRemoveMemberActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room remove member actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, error) {
		r, err := qtx.GetChatRoomRemoveMemberActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomRemoveMemberActionWithRemovedByForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room remove member actions: %w", err)
		}
		e := make([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomRemoveMemberActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, error) {
		p := query.GetChatRoomRemoveMemberActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomRemoveMemberActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room remove member actions: %w", err)
		}
		e := make([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomRemoveMemberActionOnChatRoom(query.GetChatRoomRemoveMemberActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, error) {
		p := query.GetChatRoomRemoveMemberActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomRemoveMemberActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room remove member actions: %w", err)
		}
		e := make([]entity.ChatRoomRemoveMemberActionWithRemovedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomRemoveMemberActionOnChatRoom(query.GetChatRoomRemoveMemberActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomRemoveMemberActionWithRemovedByForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomRemoveMemberActionDefaultCursorKey:
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
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{},
			fmt.Errorf("failed to get chat room remove member actions: %w", err)
	}
	return res, nil
}

// GetChatRoomRemoveMemberActionsOnChatRoom はチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetChatRoomRemoveMemberActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomRemoveMemberActionParam,
	order parameter.ChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	return getChatRoomRemoveMemberActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomRemoveMemberActionsOnChatRoomWithSd はSD付きでチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetChatRoomRemoveMemberActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomRemoveMemberActionParam,
	order parameter.ChatRoomRemoveMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomRemoveMemberActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを取得する内部関数です。
func getPluralChatRoomRemoveMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomRemoveMemberActionIDs []uuid.UUID,
	_ parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	var e []query.GetPluralChatRoomRemoveMemberActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomRemoveMemberActions(ctx, chatRoomRemoveMemberActionIDs)
	} else {
		var ne []query.GetPluralChatRoomRemoveMemberActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomRemoveMemberActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomRemoveMemberActionsUseNumberedPaginateParams{
				Limit:                         int32(np.Limit.Int64),
				Offset:                        int32(np.Offset.Int64),
				ChatRoomRemoveMemberActionIds: chatRoomRemoveMemberActionIDs,
			})
		e = make([]query.GetPluralChatRoomRemoveMemberActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomRemoveMemberActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{},
			fmt.Errorf("failed to get chat room remove member actions: %w", err)
	}
	entities := make([]entity.ChatRoomRemoveMemberActionWithRemovedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomRemoveMemberActionOnChatRoom(
			query.GetChatRoomRemoveMemberActionsOnChatRoomRow(v)).ChatRoomRemoveMemberActionWithRemovedBy
	}
	return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{Data: entities}, nil
}

// GetPluralChatRoomRemoveMemberActions は複数のチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomRemoveMemberActions(
	ctx context.Context, chatRoomRemoveMemberActionIDs []uuid.UUID,
	order parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	return getPluralChatRoomRemoveMemberActions(ctx, a.query, chatRoomRemoveMemberActionIDs, order, np)
}

// GetPluralChatRoomRemoveMemberActionsWithSd はSD付きで複数のチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomRemoveMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomRemoveMemberActionIDs []uuid.UUID,
	order parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomRemoveMemberActions(ctx, qtx, chatRoomRemoveMemberActionIDs, order, np)
}

// getPluralChatRoomRemoveMemberActionsByChatRoomActionIDs はチャットルームメンバー追放アクションを取得する内部関数です。
func getPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	_ parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	var e []query.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(ctx, chatRoomActionIDs)
	} else {
		var ne []query.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
			})
		e = make([]query.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{},
			fmt.Errorf("failed to get chat room remove member actions: %w", err)
	}
	entities := make([]entity.ChatRoomRemoveMemberActionWithRemovedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomRemoveMemberActionOnChatRoom(
			query.GetChatRoomRemoveMemberActionsOnChatRoomRow(v)).ChatRoomRemoveMemberActionWithRemovedBy
	}
	return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{Data: entities}, nil
}

// GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDs はチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	return getPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(ctx, a.query, chatRoomActionIDs, order, np)
}

// GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsWithSd はSD付きでチャットルームメンバー追放アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomRemoveMemberActionsByChatRoomActionIDsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomRemoveMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomRemoveMemberActionWithRemovedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomRemoveMemberActionsByChatRoomActionIDs(ctx, qtx, chatRoomActionIDs, order, np)
}

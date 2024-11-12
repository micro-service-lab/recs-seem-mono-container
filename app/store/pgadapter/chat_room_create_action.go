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

func convChatRoomCreateActionOnChatRoom(
	e query.GetChatRoomCreateActionsOnChatRoomRow,
) entity.ChatRoomCreateActionWithCreatedByForQuery {
	var createdBy entity.NullableEntity[entity.SimpleMember]
	if e.CreatedBy.Valid {
		createdBy = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.CreatedBy.Bytes,
				Name:           e.CreateMemberName.String,
				Email:          e.CreateMemberEmail.String,
				FirstName:      entity.String(e.CreateMemberFirstName),
				LastName:       entity.String(e.CreateMemberLastName),
				ProfileImageID: entity.UUID(e.CreateMemberProfileImageID),
				GradeID:        e.CreateMemberGradeID.Bytes,
				GroupID:        e.CreateMemberGroupID.Bytes,
			},
		}
	}
	return entity.ChatRoomCreateActionWithCreatedByForQuery{
		Pkey: entity.Int(e.TChatRoomCreateActionsPkey),
		ChatRoomCreateActionWithCreatedBy: entity.ChatRoomCreateActionWithCreatedBy{
			ChatRoomCreateActionID: e.ChatRoomCreateActionID,
			ChatRoomActionID:       e.ChatRoomActionID,
			Name:                   entity.String(e.Name),
			CreatedBy:              createdBy,
		},
	}
}

// countChatRoomCreateActions はチャットルーム作成アクション数を取得する内部関数です。
func countChatRoomCreateActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomCreateActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomCreateActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room create actions: %w", err)
	}
	return c, nil
}

// CountChatRoomCreateActions はチャットルーム作成アクション数を取得します。
func (a *PgAdapter) CountChatRoomCreateActions(
	ctx context.Context, where parameter.WhereChatRoomCreateActionParam,
) (int64, error) {
	return countChatRoomCreateActions(ctx, a.query, where)
}

// CountChatRoomCreateActionsWithSd はSD付きでチャットルーム作成アクション数を取得します。
func (a *PgAdapter) CountChatRoomCreateActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomCreateActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomCreateActions(ctx, qtx, where)
}

// createChatRoomCreateAction はチャットルーム作成アクションを作成する内部関数です。
func createChatRoomCreateAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomCreateActionParam,
) (entity.ChatRoomCreateAction, error) {
	e, err := qtx.CreateChatRoomCreateAction(ctx, query.CreateChatRoomCreateActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		CreatedBy:        pgtype.UUID(param.CreatedBy),
		Name:             pgtype.Text(param.Name),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomCreateAction{}, errhandle.NewModelDuplicatedError("chat room create action")
		}
		return entity.ChatRoomCreateAction{}, fmt.Errorf("failed to create chat room create action: %w", err)
	}
	entity := entity.ChatRoomCreateAction{
		ChatRoomCreateActionID: e.ChatRoomCreateActionID,
		ChatRoomActionID:       e.ChatRoomActionID,
		Name:                   entity.String(e.Name),
		CreatedBy:              entity.UUID(e.CreatedBy),
	}
	return entity, nil
}

// CreateChatRoomCreateAction はチャットルーム作成アクションを作成します。
func (a *PgAdapter) CreateChatRoomCreateAction(
	ctx context.Context, param parameter.CreateChatRoomCreateActionParam,
) (entity.ChatRoomCreateAction, error) {
	return createChatRoomCreateAction(ctx, a.query, param)
}

// CreateChatRoomCreateActionWithSd はSD付きでチャットルーム作成アクションを作成します。
func (a *PgAdapter) CreateChatRoomCreateActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomCreateActionParam,
) (entity.ChatRoomCreateAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomCreateAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomCreateAction(ctx, qtx, param)
}

// createChatRoomCreateActions は複数のチャットルーム作成アクションを作成する内部関数です。
func createChatRoomCreateActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomCreateActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomCreateActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomCreateActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			CreatedBy:        pgtype.UUID(p.CreatedBy),
			Name:             pgtype.Text(p.Name),
		}
	}
	n, err := qtx.CreateChatRoomCreateActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room create action")
		}
		return 0, fmt.Errorf("failed to create chat room create actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomCreateActions は複数のチャットルーム作成アクションを作成します。
func (a *PgAdapter) CreateChatRoomCreateActions(
	ctx context.Context, params []parameter.CreateChatRoomCreateActionParam,
) (int64, error) {
	return createChatRoomCreateActions(ctx, a.query, params)
}

// CreateChatRoomCreateActionsWithSd はSD付きで複数のチャットルーム作成アクションを作成します。
func (a *PgAdapter) CreateChatRoomCreateActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomCreateActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomCreateActions(ctx, qtx, params)
}

// deleteChatRoomCreateAction はチャットルーム作成アクションを削除する内部関数です。
func deleteChatRoomCreateAction(
	ctx context.Context, qtx *query.Queries, chatRoomCreateActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomCreateAction(ctx, chatRoomCreateActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room create action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room create action")
	}
	return c, nil
}

// DeleteChatRoomCreateAction はチャットルーム作成アクションを削除します。
func (a *PgAdapter) DeleteChatRoomCreateAction(
	ctx context.Context, chatRoomCreateActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomCreateAction(ctx, a.query, chatRoomCreateActionID)
}

// DeleteChatRoomCreateActionWithSd はSD付きでチャットルーム作成アクションを削除します。
func (a *PgAdapter) DeleteChatRoomCreateActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomCreateActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomCreateAction(ctx, qtx, chatRoomCreateActionID)
}

// pluralDeleteChatRoomCreateActions は複数のチャットルーム作成アクションを削除する内部関数です。
func pluralDeleteChatRoomCreateActions(
	ctx context.Context, qtx *query.Queries, chatRoomCreateActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomCreateActions(ctx, chatRoomCreateActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room create actions: %w", err)
	}
	if c != int64(len(chatRoomCreateActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room create action")
	}
	return c, nil
}

// PluralDeleteChatRoomCreateActions は複数のチャットルーム作成アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomCreateActions(
	ctx context.Context, chatRoomCreateActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomCreateActions(ctx, a.query, chatRoomCreateActionIDs)
}

// PluralDeleteChatRoomCreateActionsWithSd はSD付きで複数のチャットルーム作成アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomCreateActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomCreateActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomCreateActions(ctx, qtx, chatRoomCreateActionIDs)
}

// getChatRoomCreateActions はチャットルーム作成アクションを取得する内部関数です。
func getChatRoomCreateActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomCreateActionParam,
	order parameter.ChatRoomCreateActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	eConvFunc := func(
		e entity.ChatRoomCreateActionWithCreatedByForQuery,
	) (entity.ChatRoomCreateActionWithCreatedBy, error) {
		return e.ChatRoomCreateActionWithCreatedBy, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomCreateActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room create actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomCreateActionWithCreatedByForQuery, error) {
		r, err := qtx.GetChatRoomCreateActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomCreateActionWithCreatedByForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room create actions: %w", err)
		}
		e := make([]entity.ChatRoomCreateActionWithCreatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomCreateActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomCreateActionWithCreatedByForQuery, error) {
		p := query.GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomCreateActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room create actions: %w", err)
		}
		e := make([]entity.ChatRoomCreateActionWithCreatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomCreateActionOnChatRoom(query.GetChatRoomCreateActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomCreateActionWithCreatedByForQuery, error) {
		p := query.GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomCreateActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room create actions: %w", err)
		}
		e := make([]entity.ChatRoomCreateActionWithCreatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomCreateActionOnChatRoom(query.GetChatRoomCreateActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomCreateActionWithCreatedByForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomCreateActionDefaultCursorKey:
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
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{},
			fmt.Errorf("failed to get chat room create actions: %w", err)
	}
	return res, nil
}

// GetChatRoomCreateActionsOnChatRoom はチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetChatRoomCreateActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomCreateActionParam,
	order parameter.ChatRoomCreateActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	return getChatRoomCreateActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomCreateActionsOnChatRoomWithSd はSD付きでチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetChatRoomCreateActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomCreateActionParam,
	order parameter.ChatRoomCreateActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomCreateActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomCreateActions は複数のチャットルーム作成アクションを取得する内部関数です。
func getPluralChatRoomCreateActions(
	ctx context.Context, qtx *query.Queries, chatRoomCreateActionIDs []uuid.UUID,
	_ parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	var e []query.GetPluralChatRoomCreateActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomCreateActions(ctx, chatRoomCreateActionIDs)
	} else {
		var ne []query.GetPluralChatRoomCreateActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomCreateActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomCreateActionsUseNumberedPaginateParams{
				Limit:                   int32(np.Limit.Int64),
				Offset:                  int32(np.Offset.Int64),
				ChatRoomCreateActionIds: chatRoomCreateActionIDs,
			})
		e = make([]query.GetPluralChatRoomCreateActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomCreateActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{},
			fmt.Errorf("failed to get chat room create actions: %w", err)
	}
	entities := make([]entity.ChatRoomCreateActionWithCreatedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomCreateActionOnChatRoom(
			query.GetChatRoomCreateActionsOnChatRoomRow(v)).ChatRoomCreateActionWithCreatedBy
	}
	return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{Data: entities}, nil
}

// GetPluralChatRoomCreateActions は複数のチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomCreateActions(
	ctx context.Context, chatRoomCreateActionIDs []uuid.UUID,
	order parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	return getPluralChatRoomCreateActions(ctx, a.query, chatRoomCreateActionIDs, order, np)
}

// GetPluralChatRoomCreateActionsWithSd はSD付きで複数のチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomCreateActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomCreateActionIDs []uuid.UUID,
	order parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomCreateActions(ctx, qtx, chatRoomCreateActionIDs, order, np)
}

// getPluralChatRoomCreateActionsByChatRoomActionIDs はチャットルーム作成アクションを取得する内部関数です。
func getPluralChatRoomCreateActionsByChatRoomActionIDs(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	_ parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	var e []query.GetPluralChatRoomCreateActionsByChatRoomActionIDsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomCreateActionsByChatRoomActionIDs(ctx, chatRoomActionIDs)
	} else {
		var ne []query.GetPluralChatRoomCreateActionsByChatRoomActionIDsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomCreateActionsByChatRoomActionIDsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomCreateActionsByChatRoomActionIDsUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
			})
		e = make([]query.GetPluralChatRoomCreateActionsByChatRoomActionIDsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomCreateActionsByChatRoomActionIDsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{},
			fmt.Errorf("failed to get chat room create actions: %w", err)
	}
	entities := make([]entity.ChatRoomCreateActionWithCreatedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomCreateActionOnChatRoom(
			query.GetChatRoomCreateActionsOnChatRoomRow(v)).ChatRoomCreateActionWithCreatedBy
	}
	return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{Data: entities}, nil
}

// GetPluralChatRoomCreateActionsByChatRoomActionIDs はチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomCreateActionsByChatRoomActionIDs(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	return getPluralChatRoomCreateActionsByChatRoomActionIDs(ctx, a.query, chatRoomActionIDs, order, np)
}

// GetPluralChatRoomCreateActionsByChatRoomActionIDsWithSd はSD付きでチャットルーム作成アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomCreateActionsByChatRoomActionIDsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomCreateActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomCreateActionWithCreatedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomCreateActionWithCreatedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomCreateActionsByChatRoomActionIDs(ctx, qtx, chatRoomActionIDs, order, np)
}

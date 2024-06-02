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

func convChatRoomAddMemberActionOnChatRoom(
	e query.GetChatRoomAddMemberActionsOnChatRoomRow,
) entity.ChatRoomAddMemberActionWithAddedByForQuery {
	var addedBy entity.NullableEntity[entity.SimpleMember]
	if e.AddedBy.Valid {
		addedBy = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.AddedBy.Bytes,
				Name:           e.AddMemberName.String,
				Email:          e.AddMemberEmail.String,
				FirstName:      entity.String(e.AddMemberFirstName),
				LastName:       entity.String(e.AddMemberLastName),
				ProfileImageID: entity.UUID(e.AddMemberProfileImageID),
				GradeID:        e.AddMemberGradeID.Bytes,
				GroupID:        e.AddMemberGroupID.Bytes,
			},
		}
	}
	return entity.ChatRoomAddMemberActionWithAddedByForQuery{
		Pkey: entity.Int(e.TChatRoomAddMemberActionsPkey),
		ChatRoomAddMemberActionWithAddedBy: entity.ChatRoomAddMemberActionWithAddedBy{
			ChatRoomAddMemberActionID: e.ChatRoomAddMemberActionID,
			ChatRoomActionID:          e.ChatRoomActionID,
			AddedBy:                   addedBy,
		},
	}
}

// countChatRoomAddMemberActions はチャットルームメンバー追加アクション数を取得する内部関数です。
func countChatRoomAddMemberActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomAddMemberActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomAddMemberActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room add member actions: %w", err)
	}
	return c, nil
}

// CountChatRoomAddMemberActions はチャットルームメンバー追加アクション数を取得します。
func (a *PgAdapter) CountChatRoomAddMemberActions(
	ctx context.Context, where parameter.WhereChatRoomAddMemberActionParam,
) (int64, error) {
	return countChatRoomAddMemberActions(ctx, a.query, where)
}

// CountChatRoomAddMemberActionsWithSd はSD付きでチャットルームメンバー追加アクション数を取得します。
func (a *PgAdapter) CountChatRoomAddMemberActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomAddMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomAddMemberActions(ctx, qtx, where)
}

// createChatRoomAddMemberAction はチャットルームメンバー追加アクションを作成する内部関数です。
func createChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomAddMemberActionParam,
) (entity.ChatRoomAddMemberAction, error) {
	e, err := qtx.CreateChatRoomAddMemberAction(ctx, query.CreateChatRoomAddMemberActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		AddedBy:          pgtype.UUID(param.AddedBy),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomAddMemberAction{}, errhandle.NewModelDuplicatedError("chat room add member action")
		}
		return entity.ChatRoomAddMemberAction{}, fmt.Errorf("failed to create chat room add member action: %w", err)
	}
	entity := entity.ChatRoomAddMemberAction{
		ChatRoomAddMemberActionID: e.ChatRoomAddMemberActionID,
		ChatRoomActionID:          e.ChatRoomActionID,
		AddedBy:                   entity.UUID(e.AddedBy),
	}
	return entity, nil
}

// CreateChatRoomAddMemberAction はチャットルームメンバー追加アクションを作成します。
func (a *PgAdapter) CreateChatRoomAddMemberAction(
	ctx context.Context, param parameter.CreateChatRoomAddMemberActionParam,
) (entity.ChatRoomAddMemberAction, error) {
	return createChatRoomAddMemberAction(ctx, a.query, param)
}

// CreateChatRoomAddMemberActionWithSd はSD付きでチャットルームメンバー追加アクションを作成します。
func (a *PgAdapter) CreateChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomAddMemberActionParam,
) (entity.ChatRoomAddMemberAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomAddMemberAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomAddMemberAction(ctx, qtx, param)
}

// createChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを作成する内部関数です。
func createChatRoomAddMemberActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomAddMemberActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomAddMemberActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomAddMemberActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			AddedBy:          pgtype.UUID(p.AddedBy),
		}
	}
	n, err := qtx.CreateChatRoomAddMemberActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room add member action")
		}
		return 0, fmt.Errorf("failed to create chat room add member actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを作成します。
func (a *PgAdapter) CreateChatRoomAddMemberActions(
	ctx context.Context, params []parameter.CreateChatRoomAddMemberActionParam,
) (int64, error) {
	return createChatRoomAddMemberActions(ctx, a.query, params)
}

// CreateChatRoomAddMemberActionsWithSd はSD付きで複数のチャットルームメンバー追加アクションを作成します。
func (a *PgAdapter) CreateChatRoomAddMemberActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomAddMemberActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomAddMemberActions(ctx, qtx, params)
}

// deleteChatRoomAddMemberAction はチャットルームメンバー追加アクションを削除する内部関数です。
func deleteChatRoomAddMemberAction(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAddMemberAction(ctx, chatRoomAddMemberActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room add member action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room add member action")
	}
	return c, nil
}

// DeleteChatRoomAddMemberAction はチャットルームメンバー追加アクションを削除します。
func (a *PgAdapter) DeleteChatRoomAddMemberAction(
	ctx context.Context, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomAddMemberAction(ctx, a.query, chatRoomAddMemberActionID)
}

// DeleteChatRoomAddMemberActionWithSd はSD付きでチャットルームメンバー追加アクションを削除します。
func (a *PgAdapter) DeleteChatRoomAddMemberActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAddMemberAction(ctx, qtx, chatRoomAddMemberActionID)
}

// pluralDeleteChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを削除する内部関数です。
func pluralDeleteChatRoomAddMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomAddMemberActions(ctx, chatRoomAddMemberActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room add member actions: %w", err)
	}
	if c != int64(len(chatRoomAddMemberActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room add member action")
	}
	return c, nil
}

// PluralDeleteChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomAddMemberActions(
	ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomAddMemberActions(ctx, a.query, chatRoomAddMemberActionIDs)
}

// PluralDeleteChatRoomAddMemberActionsWithSd はSD付きで複数のチャットルームメンバー追加アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomAddMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomAddMemberActions(ctx, qtx, chatRoomAddMemberActionIDs)
}

// getChatRoomAddMemberActions はチャットルームメンバー追加アクションを取得する内部関数です。
func getChatRoomAddMemberActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomAddMemberActionParam,
	order parameter.ChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	eConvFunc := func(
		e entity.ChatRoomAddMemberActionWithAddedByForQuery,
	) (entity.ChatRoomAddMemberActionWithAddedBy, error) {
		return e.ChatRoomAddMemberActionWithAddedBy, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomAddMemberActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room add member actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomAddMemberActionWithAddedByForQuery, error) {
		r, err := qtx.GetChatRoomAddMemberActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomAddMemberActionWithAddedByForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room add member actions: %w", err)
		}
		e := make([]entity.ChatRoomAddMemberActionWithAddedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomAddMemberActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomAddMemberActionWithAddedByForQuery, error) {
		p := query.GetChatRoomAddMemberActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomAddMemberActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room add member actions: %w", err)
		}
		e := make([]entity.ChatRoomAddMemberActionWithAddedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomAddMemberActionOnChatRoom(query.GetChatRoomAddMemberActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomAddMemberActionWithAddedByForQuery, error) {
		p := query.GetChatRoomAddMemberActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomAddMemberActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room add member actions: %w", err)
		}
		e := make([]entity.ChatRoomAddMemberActionWithAddedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomAddMemberActionOnChatRoom(query.GetChatRoomAddMemberActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomAddMemberActionWithAddedByForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomAddMemberActionDefaultCursorKey:
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
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{},
			fmt.Errorf("failed to get chat room add member actions: %w", err)
	}
	return res, nil
}

// GetChatRoomAddMemberActionsOnChatRoom はチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetChatRoomAddMemberActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomAddMemberActionParam,
	order parameter.ChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	return getChatRoomAddMemberActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomAddMemberActionsOnChatRoomWithSd はSD付きでチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetChatRoomAddMemberActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomAddMemberActionParam,
	order parameter.ChatRoomAddMemberActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomAddMemberActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを取得する内部関数です。
func getPluralChatRoomAddMemberActions(
	ctx context.Context, qtx *query.Queries, chatRoomAddMemberActionIDs []uuid.UUID,
	_ parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	var e []query.GetPluralChatRoomAddMemberActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomAddMemberActions(ctx, chatRoomAddMemberActionIDs)
	} else {
		var ne []query.GetPluralChatRoomAddMemberActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomAddMemberActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomAddMemberActionsUseNumberedPaginateParams{
				Limit:                      int32(np.Limit.Int64),
				Offset:                     int32(np.Offset.Int64),
				ChatRoomAddMemberActionIds: chatRoomAddMemberActionIDs,
			})
		e = make([]query.GetPluralChatRoomAddMemberActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomAddMemberActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{},
			fmt.Errorf("failed to get chat room add member actions: %w", err)
	}
	entities := make([]entity.ChatRoomAddMemberActionWithAddedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomAddMemberActionOnChatRoom(
			query.GetChatRoomAddMemberActionsOnChatRoomRow(v)).ChatRoomAddMemberActionWithAddedBy
	}
	return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{Data: entities}, nil
}

// GetPluralChatRoomAddMemberActions は複数のチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomAddMemberActions(
	ctx context.Context, chatRoomAddMemberActionIDs []uuid.UUID,
	order parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	return getPluralChatRoomAddMemberActions(ctx, a.query, chatRoomAddMemberActionIDs, order, np)
}

// GetPluralChatRoomAddMemberActionsWithSd はSD付きで複数のチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomAddMemberActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomAddMemberActionIDs []uuid.UUID,
	order parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomAddMemberActions(ctx, qtx, chatRoomAddMemberActionIDs, order, np)
}

// getPluralChatRoomAddMemberActionsByChatRoomActionIDs はチャットルームメンバー追加アクションを取得する内部関数です。
func getPluralChatRoomAddMemberActionsByChatRoomActionIDs(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	_ parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	var e []query.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomAddMemberActionsByChatRoomActionIDs(ctx, chatRoomActionIDs)
	} else {
		var ne []query.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
			})
		e = make([]query.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomAddMemberActionsByChatRoomActionIDsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{},
			fmt.Errorf("failed to get chat room add member actions: %w", err)
	}
	entities := make([]entity.ChatRoomAddMemberActionWithAddedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomAddMemberActionOnChatRoom(
			query.GetChatRoomAddMemberActionsOnChatRoomRow(v)).ChatRoomAddMemberActionWithAddedBy
	}
	return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{Data: entities}, nil
}

// GetPluralChatRoomAddMemberActionsByChatRoomActionIDs はチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomAddMemberActionsByChatRoomActionIDs(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	return getPluralChatRoomAddMemberActionsByChatRoomActionIDs(ctx, a.query, chatRoomActionIDs, order, np)
}

// GetPluralChatRoomAddMemberActionsByChatRoomIDsWithSd はSD付きでチャットルームメンバー追加アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomAddMemberActionsByChatRoomActionIDsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomAddMemberActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomAddMemberActionWithAddedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomAddMemberActionsByChatRoomActionIDs(ctx, qtx, chatRoomActionIDs, order, np)
}

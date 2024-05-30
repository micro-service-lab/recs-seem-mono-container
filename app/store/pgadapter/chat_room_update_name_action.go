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

func convChatRoomUpdateNameActionOnChatRoom(
	e query.GetChatRoomUpdateNameActionsOnChatRoomRow,
) entity.ChatRoomUpdateNameActionWithUpdatedByForQuery {
	var updatedBy entity.NullableEntity[entity.SimpleMember]
	if e.UpdatedBy.Valid {
		updatedBy = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.UpdatedBy.Bytes,
				Name:           e.UpdateMemberName.String,
				Email:          e.UpdateMemberEmail.String,
				FirstName:      entity.String(e.UpdateMemberFirstName),
				LastName:       entity.String(e.UpdateMemberLastName),
				ProfileImageID: entity.UUID(e.UpdateMemberProfileImageID),
			},
		}
	}
	return entity.ChatRoomUpdateNameActionWithUpdatedByForQuery{
		Pkey: entity.Int(e.TChatRoomUpdateNameActionsPkey),
		ChatRoomUpdateNameActionWithUpdatedBy: entity.ChatRoomUpdateNameActionWithUpdatedBy{
			ChatRoomUpdateNameActionID: e.ChatRoomUpdateNameActionID,
			ChatRoomActionID:           e.ChatRoomActionID,
			Name:                       e.Name,
			UpdatedBy:                  updatedBy,
		},
	}
}

// countChatRoomUpdateNameActions はチャットルーム名前更新アクション数を取得する内部関数です。
func countChatRoomUpdateNameActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomUpdateNameActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomUpdateNameActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room update name actions: %w", err)
	}
	return c, nil
}

// CountChatRoomUpdateNameActions はチャットルーム名前更新アクション数を取得します。
func (a *PgAdapter) CountChatRoomUpdateNameActions(
	ctx context.Context, where parameter.WhereChatRoomUpdateNameActionParam,
) (int64, error) {
	return countChatRoomUpdateNameActions(ctx, a.query, where)
}

// CountChatRoomUpdateNameActionsWithSd はSD付きでチャットルーム名前更新アクション数を取得します。
func (a *PgAdapter) CountChatRoomUpdateNameActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomUpdateNameActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomUpdateNameActions(ctx, qtx, where)
}

// createChatRoomUpdateNameAction はチャットルーム名前更新アクションを作成する内部関数です。
func createChatRoomUpdateNameAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomUpdateNameActionParam,
) (entity.ChatRoomUpdateNameAction, error) {
	e, err := qtx.CreateChatRoomUpdateNameAction(ctx, query.CreateChatRoomUpdateNameActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		UpdatedBy:        pgtype.UUID(param.UpdatedBy),
		Name:             param.Name,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomUpdateNameAction{}, errhandle.NewModelDuplicatedError("chat room update name action")
		}
		return entity.ChatRoomUpdateNameAction{}, fmt.Errorf("failed to create chat room update name action: %w", err)
	}
	entity := entity.ChatRoomUpdateNameAction{
		ChatRoomUpdateNameActionID: e.ChatRoomUpdateNameActionID,
		ChatRoomActionID:           e.ChatRoomActionID,
		Name:                       e.Name,
		UpdatedBy:                  entity.UUID(e.UpdatedBy),
	}
	return entity, nil
}

// CreateChatRoomUpdateNameAction はチャットルーム名前更新アクションを作成します。
func (a *PgAdapter) CreateChatRoomUpdateNameAction(
	ctx context.Context, param parameter.CreateChatRoomUpdateNameActionParam,
) (entity.ChatRoomUpdateNameAction, error) {
	return createChatRoomUpdateNameAction(ctx, a.query, param)
}

// CreateChatRoomUpdateNameActionWithSd はSD付きでチャットルーム名前更新アクションを作成します。
func (a *PgAdapter) CreateChatRoomUpdateNameActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomUpdateNameActionParam,
) (entity.ChatRoomUpdateNameAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomUpdateNameAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomUpdateNameAction(ctx, qtx, param)
}

// createChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを作成する内部関数です。
func createChatRoomUpdateNameActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomUpdateNameActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomUpdateNameActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomUpdateNameActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			UpdatedBy:        pgtype.UUID(p.UpdatedBy),
			Name:             p.Name,
		}
	}
	n, err := qtx.CreateChatRoomUpdateNameActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room update name action")
		}
		return 0, fmt.Errorf("failed to create chat room update name actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを作成します。
func (a *PgAdapter) CreateChatRoomUpdateNameActions(
	ctx context.Context, params []parameter.CreateChatRoomUpdateNameActionParam,
) (int64, error) {
	return createChatRoomUpdateNameActions(ctx, a.query, params)
}

// CreateChatRoomUpdateNameActionsWithSd はSD付きで複数のチャットルーム名前更新アクションを作成します。
func (a *PgAdapter) CreateChatRoomUpdateNameActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomUpdateNameActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomUpdateNameActions(ctx, qtx, params)
}

// deleteChatRoomUpdateNameAction はチャットルーム名前更新アクションを削除する内部関数です。
func deleteChatRoomUpdateNameAction(
	ctx context.Context, qtx *query.Queries, chatRoomUpdateNameActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomUpdateNameAction(ctx, chatRoomUpdateNameActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room update name action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room update name action")
	}
	return c, nil
}

// DeleteChatRoomUpdateNameAction はチャットルーム名前更新アクションを削除します。
func (a *PgAdapter) DeleteChatRoomUpdateNameAction(
	ctx context.Context, chatRoomUpdateNameActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomUpdateNameAction(ctx, a.query, chatRoomUpdateNameActionID)
}

// DeleteChatRoomUpdateNameActionWithSd はSD付きでチャットルーム名前更新アクションを削除します。
func (a *PgAdapter) DeleteChatRoomUpdateNameActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomUpdateNameActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomUpdateNameAction(ctx, qtx, chatRoomUpdateNameActionID)
}

// pluralDeleteChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを削除する内部関数です。
func pluralDeleteChatRoomUpdateNameActions(
	ctx context.Context, qtx *query.Queries, chatRoomUpdateNameActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomUpdateNameActions(ctx, chatRoomUpdateNameActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room update name actions: %w", err)
	}
	if c != int64(len(chatRoomUpdateNameActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room update name action")
	}
	return c, nil
}

// PluralDeleteChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomUpdateNameActions(
	ctx context.Context, chatRoomUpdateNameActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomUpdateNameActions(ctx, a.query, chatRoomUpdateNameActionIDs)
}

// PluralDeleteChatRoomUpdateNameActionsWithSd はSD付きで複数のチャットルーム名前更新アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomUpdateNameActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomUpdateNameActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomUpdateNameActions(ctx, qtx, chatRoomUpdateNameActionIDs)
}

// getChatRoomUpdateNameActions はチャットルーム名前更新アクションを取得する内部関数です。
func getChatRoomUpdateNameActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomUpdateNameActionParam,
	order parameter.ChatRoomUpdateNameActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	eConvFunc := func(
		e entity.ChatRoomUpdateNameActionWithUpdatedByForQuery,
	) (entity.ChatRoomUpdateNameActionWithUpdatedBy, error) {
		return e.ChatRoomUpdateNameActionWithUpdatedBy, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomUpdateNameActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room update name actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, error) {
		r, err := qtx.GetChatRoomUpdateNameActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomUpdateNameActionWithUpdatedByForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room update name actions: %w", err)
		}
		e := make([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomUpdateNameActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, error) {
		p := query.GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room update name actions: %w", err)
		}
		e := make([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomUpdateNameActionOnChatRoom(query.GetChatRoomUpdateNameActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, error) {
		p := query.GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room update name actions: %w", err)
		}
		e := make([]entity.ChatRoomUpdateNameActionWithUpdatedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomUpdateNameActionOnChatRoom(query.GetChatRoomUpdateNameActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomUpdateNameActionWithUpdatedByForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomUpdateNameActionDefaultCursorKey:
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
		return store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy]{},
			fmt.Errorf("failed to get chat room update name actions: %w", err)
	}
	return res, nil
}

// GetChatRoomUpdateNameActionsOnChatRoom はチャットルーム名前更新アクションを取得します。
func (a *PgAdapter) GetChatRoomUpdateNameActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomUpdateNameActionParam,
	order parameter.ChatRoomUpdateNameActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	return getChatRoomUpdateNameActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomUpdateNameActionsOnChatRoomWithSd はSD付きでチャットルーム名前更新アクションを取得します。
func (a *PgAdapter) GetChatRoomUpdateNameActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomUpdateNameActionParam,
	order parameter.ChatRoomUpdateNameActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomUpdateNameActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを取得する内部関数です。
func getPluralChatRoomUpdateNameActions(
	ctx context.Context, qtx *query.Queries, chatRoomUpdateNameActionIDs []uuid.UUID,
	_ parameter.ChatRoomUpdateNameActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	var e []query.GetPluralChatRoomUpdateNameActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomUpdateNameActions(ctx, chatRoomUpdateNameActionIDs)
	} else {
		var ne []query.GetPluralChatRoomUpdateNameActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomUpdateNameActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomUpdateNameActionsUseNumberedPaginateParams{
				Limit:                       int32(np.Limit.Int64),
				Offset:                      int32(np.Offset.Int64),
				ChatRoomUpdateNameActionIds: chatRoomUpdateNameActionIDs,
			})
		e = make([]query.GetPluralChatRoomUpdateNameActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomUpdateNameActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy]{},
			fmt.Errorf("failed to get chat room update name actions: %w", err)
	}
	entities := make([]entity.ChatRoomUpdateNameActionWithUpdatedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomUpdateNameActionOnChatRoom(
			query.GetChatRoomUpdateNameActionsOnChatRoomRow(v)).ChatRoomUpdateNameActionWithUpdatedBy
	}
	return store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy]{Data: entities}, nil
}

// GetPluralChatRoomUpdateNameActions は複数のチャットルーム名前更新アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomUpdateNameActions(
	ctx context.Context, chatRoomUpdateNameActionIDs []uuid.UUID,
	order parameter.ChatRoomUpdateNameActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	return getPluralChatRoomUpdateNameActions(ctx, a.query, chatRoomUpdateNameActionIDs, order, np)
}

// GetPluralChatRoomUpdateNameActionsWithSd はSD付きで複数のチャットルーム名前更新アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomUpdateNameActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomUpdateNameActionIDs []uuid.UUID,
	order parameter.ChatRoomUpdateNameActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomUpdateNameActionWithUpdatedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomUpdateNameActions(ctx, qtx, chatRoomUpdateNameActionIDs, order, np)
}

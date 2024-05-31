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

func convChatRoomDeleteMessageActionOnChatRoom(
	e query.GetChatRoomDeleteMessageActionsOnChatRoomRow,
) entity.ChatRoomDeleteMessageActionWithDeletedByForQuery {
	var addedBy entity.NullableEntity[entity.SimpleMember]
	if e.DeletedBy.Valid {
		addedBy = entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       e.DeletedBy.Bytes,
				Name:           e.DeleteMessageMemberName.String,
				Email:          e.DeleteMessageMemberEmail.String,
				FirstName:      entity.String(e.DeleteMessageMemberFirstName),
				LastName:       entity.String(e.DeleteMessageMemberLastName),
				ProfileImageID: entity.UUID(e.DeleteMessageMemberProfileImageID),
				GradeID:        e.DeleteMessageMemberGradeID.Bytes,
				GroupID:        e.DeleteMessageMemberGroupID.Bytes,
			},
		}
	}
	return entity.ChatRoomDeleteMessageActionWithDeletedByForQuery{
		Pkey: entity.Int(e.TChatRoomDeleteMessageActionsPkey),
		ChatRoomDeleteMessageActionWithDeletedBy: entity.ChatRoomDeleteMessageActionWithDeletedBy{
			ChatRoomDeleteMessageActionID: e.ChatRoomDeleteMessageActionID,
			ChatRoomActionID:              e.ChatRoomActionID,
			DeletedBy:                     addedBy,
		},
	}
}

// countChatRoomDeleteMessageActions はチャットルームメッセージ削除アクション数を取得する内部関数です。
func countChatRoomDeleteMessageActions(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereChatRoomDeleteMessageActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomDeleteMessageActions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room delete message actions: %w", err)
	}
	return c, nil
}

// CountChatRoomDeleteMessageActions はチャットルームメッセージ削除アクション数を取得します。
func (a *PgAdapter) CountChatRoomDeleteMessageActions(
	ctx context.Context, where parameter.WhereChatRoomDeleteMessageActionParam,
) (int64, error) {
	return countChatRoomDeleteMessageActions(ctx, a.query, where)
}

// CountChatRoomDeleteMessageActionsWithSd はSD付きでチャットルームメッセージ削除アクション数を取得します。
func (a *PgAdapter) CountChatRoomDeleteMessageActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomDeleteMessageActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomDeleteMessageActions(ctx, qtx, where)
}

// createChatRoomDeleteMessageAction はチャットルームメッセージ削除アクションを作成する内部関数です。
func createChatRoomDeleteMessageAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomDeleteMessageActionParam,
) (entity.ChatRoomDeleteMessageAction, error) {
	e, err := qtx.CreateChatRoomDeleteMessageAction(ctx, query.CreateChatRoomDeleteMessageActionParams{
		ChatRoomActionID: param.ChatRoomActionID,
		DeletedBy:        pgtype.UUID(param.DeletedBy),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomDeleteMessageAction{}, errhandle.NewModelDuplicatedError("chat room delete message action")
		}
		return entity.ChatRoomDeleteMessageAction{}, fmt.Errorf("failed to create chat room delete message action: %w", err)
	}
	entity := entity.ChatRoomDeleteMessageAction{
		ChatRoomDeleteMessageActionID: e.ChatRoomDeleteMessageActionID,
		ChatRoomActionID:              e.ChatRoomActionID,
		DeletedBy:                     entity.UUID(e.DeletedBy),
	}
	return entity, nil
}

// CreateChatRoomDeleteMessageAction はチャットルームメッセージ削除アクションを作成します。
func (a *PgAdapter) CreateChatRoomDeleteMessageAction(
	ctx context.Context, param parameter.CreateChatRoomDeleteMessageActionParam,
) (entity.ChatRoomDeleteMessageAction, error) {
	return createChatRoomDeleteMessageAction(ctx, a.query, param)
}

// CreateChatRoomDeleteMessageActionWithSd はSD付きでチャットルームメッセージ削除アクションを作成します。
func (a *PgAdapter) CreateChatRoomDeleteMessageActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomDeleteMessageActionParam,
) (entity.ChatRoomDeleteMessageAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomDeleteMessageAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomDeleteMessageAction(ctx, qtx, param)
}

// createChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを作成する内部関数です。
func createChatRoomDeleteMessageActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomDeleteMessageActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomDeleteMessageActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomDeleteMessageActionsParams{
			ChatRoomActionID: p.ChatRoomActionID,
			DeletedBy:        pgtype.UUID(p.DeletedBy),
		}
	}
	n, err := qtx.CreateChatRoomDeleteMessageActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room delete message action")
		}
		return 0, fmt.Errorf("failed to create chat room delete message actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを作成します。
func (a *PgAdapter) CreateChatRoomDeleteMessageActions(
	ctx context.Context, params []parameter.CreateChatRoomDeleteMessageActionParam,
) (int64, error) {
	return createChatRoomDeleteMessageActions(ctx, a.query, params)
}

// CreateChatRoomDeleteMessageActionsWithSd はSD付きで複数のチャットルームメッセージ削除アクションを作成します。
func (a *PgAdapter) CreateChatRoomDeleteMessageActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomDeleteMessageActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomDeleteMessageActions(ctx, qtx, params)
}

// deleteChatRoomDeleteMessageAction はチャットルームメッセージ削除アクションを削除する内部関数です。
func deleteChatRoomDeleteMessageAction(
	ctx context.Context, qtx *query.Queries, chatRoomDeleteMessageActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomDeleteMessageAction(ctx, chatRoomDeleteMessageActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room delete message action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room delete message action")
	}
	return c, nil
}

// DeleteChatRoomDeleteMessageAction はチャットルームメッセージ削除アクションを削除します。
func (a *PgAdapter) DeleteChatRoomDeleteMessageAction(
	ctx context.Context, chatRoomDeleteMessageActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomDeleteMessageAction(ctx, a.query, chatRoomDeleteMessageActionID)
}

// DeleteChatRoomDeleteMessageActionWithSd はSD付きでチャットルームメッセージ削除アクションを削除します。
func (a *PgAdapter) DeleteChatRoomDeleteMessageActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomDeleteMessageActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomDeleteMessageAction(ctx, qtx, chatRoomDeleteMessageActionID)
}

// pluralDeleteChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを削除する内部関数です。
func pluralDeleteChatRoomDeleteMessageActions(
	ctx context.Context, qtx *query.Queries, chatRoomDeleteMessageActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomDeleteMessageActions(ctx, chatRoomDeleteMessageActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room delete message actions: %w", err)
	}
	if c != int64(len(chatRoomDeleteMessageActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room delete message action")
	}
	return c, nil
}

// PluralDeleteChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomDeleteMessageActions(
	ctx context.Context, chatRoomDeleteMessageActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomDeleteMessageActions(ctx, a.query, chatRoomDeleteMessageActionIDs)
}

// PluralDeleteChatRoomDeleteMessageActionsWithSd はSD付きで複数のチャットルームメッセージ削除アクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomDeleteMessageActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomDeleteMessageActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomDeleteMessageActions(ctx, qtx, chatRoomDeleteMessageActionIDs)
}

// getChatRoomDeleteMessageActions はチャットルームメッセージ削除アクションを取得する内部関数です。
func getChatRoomDeleteMessageActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	_ parameter.WhereChatRoomDeleteMessageActionParam,
	order parameter.ChatRoomDeleteMessageActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	eConvFunc := func(
		e entity.ChatRoomDeleteMessageActionWithDeletedByForQuery,
	) (entity.ChatRoomDeleteMessageActionWithDeletedBy, error) {
		return e.ChatRoomDeleteMessageActionWithDeletedBy, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomDeleteMessageActions(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room delete message actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, error) {
		r, err := qtx.GetChatRoomDeleteMessageActionsOnChatRoom(ctx, chatRoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomDeleteMessageActionWithDeletedByForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room delete message actions: %w", err)
		}
		e := make([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomDeleteMessageActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, error) {
		p := query.GetChatRoomDeleteMessageActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetChatRoomDeleteMessageActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room delete message actions: %w", err)
		}
		e := make([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomDeleteMessageActionOnChatRoom(query.GetChatRoomDeleteMessageActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, error) {
		p := query.GetChatRoomDeleteMessageActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID: chatRoomID,
			Limit:      limit,
			Offset:     offset,
		}
		r, err := qtx.GetChatRoomDeleteMessageActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room delete message actions: %w", err)
		}
		e := make([]entity.ChatRoomDeleteMessageActionWithDeletedByForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomDeleteMessageActionOnChatRoom(query.GetChatRoomDeleteMessageActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomDeleteMessageActionWithDeletedByForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomDeleteMessageActionDefaultCursorKey:
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
		return store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy]{},
			fmt.Errorf("failed to get chat room delete message actions: %w", err)
	}
	return res, nil
}

// GetChatRoomDeleteMessageActionsOnChatRoom はチャットルームメッセージ削除アクションを取得します。
func (a *PgAdapter) GetChatRoomDeleteMessageActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomDeleteMessageActionParam,
	order parameter.ChatRoomDeleteMessageActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	return getChatRoomDeleteMessageActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomDeleteMessageActionsOnChatRoomWithSd はSD付きでチャットルームメッセージ削除アクションを取得します。
func (a *PgAdapter) GetChatRoomDeleteMessageActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomDeleteMessageActionParam,
	order parameter.ChatRoomDeleteMessageActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomDeleteMessageActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを取得する内部関数です。
func getPluralChatRoomDeleteMessageActions(
	ctx context.Context, qtx *query.Queries, chatRoomDeleteMessageActionIDs []uuid.UUID,
	_ parameter.ChatRoomDeleteMessageActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	var e []query.GetPluralChatRoomDeleteMessageActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomDeleteMessageActions(ctx, chatRoomDeleteMessageActionIDs)
	} else {
		var ne []query.GetPluralChatRoomDeleteMessageActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomDeleteMessageActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomDeleteMessageActionsUseNumberedPaginateParams{
				Limit:                          int32(np.Limit.Int64),
				Offset:                         int32(np.Offset.Int64),
				ChatRoomDeleteMessageActionIds: chatRoomDeleteMessageActionIDs,
			})
		e = make([]query.GetPluralChatRoomDeleteMessageActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomDeleteMessageActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy]{},
			fmt.Errorf("failed to get chat room delete message actions: %w", err)
	}
	entities := make([]entity.ChatRoomDeleteMessageActionWithDeletedBy, len(e))
	for i, v := range e {
		entities[i] = convChatRoomDeleteMessageActionOnChatRoom(
			query.GetChatRoomDeleteMessageActionsOnChatRoomRow(v)).ChatRoomDeleteMessageActionWithDeletedBy
	}
	return store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy]{Data: entities}, nil
}

// GetPluralChatRoomDeleteMessageActions は複数のチャットルームメッセージ削除アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomDeleteMessageActions(
	ctx context.Context, chatRoomDeleteMessageActionIDs []uuid.UUID,
	order parameter.ChatRoomDeleteMessageActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	return getPluralChatRoomDeleteMessageActions(ctx, a.query, chatRoomDeleteMessageActionIDs, order, np)
}

// GetPluralChatRoomDeleteMessageActionsWithSd はSD付きで複数のチャットルームメッセージ削除アクションを取得します。
func (a *PgAdapter) GetPluralChatRoomDeleteMessageActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomDeleteMessageActionIDs []uuid.UUID,
	order parameter.ChatRoomDeleteMessageActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomDeleteMessageActionWithDeletedBy]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomDeleteMessageActions(ctx, qtx, chatRoomDeleteMessageActionIDs, order, np)
}

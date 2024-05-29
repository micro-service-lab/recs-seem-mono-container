package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func countChatRoomActionTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereChatRoomActionTypeParam,
) (int64, error) {
	p := query.CountChatRoomActionTypesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountChatRoomActionTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room action type: %w", err)
	}
	return c, nil
}

// CountChatRoomActionTypes チャットルームアクションタイプ数を取得する。
func (a *PgAdapter) CountChatRoomActionTypes(
	ctx context.Context, where parameter.WhereChatRoomActionTypeParam,
) (int64, error) {
	return countChatRoomActionTypes(ctx, a.query, where)
}

// CountChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプ数を取得する。
func (a *PgAdapter) CountChatRoomActionTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomActionTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomActionTypes(ctx, qtx, where)
}

func createChatRoomActionType(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomActionTypeParam,
) (entity.ChatRoomActionType, error) {
	p := query.CreateChatRoomActionTypeParams{
		Name: param.Name,
		Key:  param.Key,
	}
	e, err := qtx.CreateChatRoomActionType(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomActionType{}, errhandle.NewModelDuplicatedError("chat room action type")
		}
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to create chat room action type: %w", err)
	}
	entity := entity.ChatRoomActionType{
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		Name:                 e.Name,
		Key:                  e.Key,
	}
	return entity, nil
}

// CreateChatRoomActionType チャットルームアクションタイプを作成する。
func (a *PgAdapter) CreateChatRoomActionType(
	ctx context.Context, param parameter.CreateChatRoomActionTypeParam,
) (entity.ChatRoomActionType, error) {
	return createChatRoomActionType(ctx, a.query, param)
}

// CreateChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを作成する。
func (a *PgAdapter) CreateChatRoomActionTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomActionTypeParam,
) (entity.ChatRoomActionType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomActionType{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomActionType(ctx, qtx, param)
}

func createChatRoomActionTypes(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomActionTypeParam,
) (int64, error) {
	p := make([]query.CreateChatRoomActionTypesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateChatRoomActionTypesParams{
			Name: param.Name,
			Key:  param.Key,
		}
	}
	c, err := qtx.CreateChatRoomActionTypes(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room action type")
		}
		return 0, fmt.Errorf("failed to create chat room action types: %w", err)
	}
	return c, nil
}

// CreateChatRoomActionTypes チャットルームアクションタイプを作成する。
func (a *PgAdapter) CreateChatRoomActionTypes(
	ctx context.Context, params []parameter.CreateChatRoomActionTypeParam,
) (int64, error) {
	return createChatRoomActionTypes(ctx, a.query, params)
}

// CreateChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを作成する。
func (a *PgAdapter) CreateChatRoomActionTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomActionTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomActionTypes(ctx, qtx, params)
}

func deleteChatRoomActionType(ctx context.Context, qtx *query.Queries, chatRoomActionTypeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteChatRoomActionType(ctx, chatRoomActionTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room action type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room action type")
	}
	return c, nil
}

// DeleteChatRoomActionType チャットルームアクションタイプを削除する。
func (a *PgAdapter) DeleteChatRoomActionType(ctx context.Context, chatRoomActionTypeID uuid.UUID) (int64, error) {
	return deleteChatRoomActionType(ctx, a.query, chatRoomActionTypeID)
}

// DeleteChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを削除する。
func (a *PgAdapter) DeleteChatRoomActionTypeWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionTypeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomActionType(ctx, qtx, chatRoomActionTypeID)
}

func deleteChatRoomActionTypeByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteChatRoomActionTypeByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room action type: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room action type")
	}
	return c, nil
}

// DeleteChatRoomActionTypeByKey チャットルームアクションタイプを削除する。
func (a *PgAdapter) DeleteChatRoomActionTypeByKey(ctx context.Context, key string) (int64, error) {
	return deleteChatRoomActionTypeByKey(ctx, a.query, key)
}

// DeleteChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを削除する。
func (a *PgAdapter) DeleteChatRoomActionTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomActionTypeByKey(ctx, qtx, key)
}

func pluralDeleteChatRoomActionTypes(
	ctx context.Context, qtx *query.Queries, chatRoomActionTypeIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomActionTypes(ctx, chatRoomActionTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete chat room action types: %w", err)
	}
	if c != int64(len(chatRoomActionTypeIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room action type")
	}
	return c, nil
}

// PluralDeleteChatRoomActionTypes チャットルームアクションタイプを複数削除する。
func (a *PgAdapter) PluralDeleteChatRoomActionTypes(
	ctx context.Context, chatRoomActionTypeIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomActionTypes(ctx, a.query, chatRoomActionTypeIDs)
}

// PluralDeleteChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを複数削除する。
func (a *PgAdapter) PluralDeleteChatRoomActionTypesWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionTypeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomActionTypes(ctx, qtx, chatRoomActionTypeIDs)
}

func findChatRoomActionTypeByID(
	ctx context.Context, qtx *query.Queries, chatRoomActionTypeID uuid.UUID,
) (entity.ChatRoomActionType, error) {
	e, err := qtx.FindChatRoomActionTypeByID(ctx, chatRoomActionTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomActionType{}, errhandle.NewModelNotFoundError("chat room action type")
		}
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	entity := entity.ChatRoomActionType{
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		Name:                 e.Name,
		Key:                  e.Key,
	}
	return entity, nil
}

// FindChatRoomActionTypeByID チャットルームアクションタイプを取得する。
func (a *PgAdapter) FindChatRoomActionTypeByID(
	ctx context.Context, chatRoomActionTypeID uuid.UUID,
) (entity.ChatRoomActionType, error) {
	return findChatRoomActionTypeByID(ctx, a.query, chatRoomActionTypeID)
}

// FindChatRoomActionTypeByIDWithSd SD付きでチャットルームアクションタイプを取得する。
func (a *PgAdapter) FindChatRoomActionTypeByIDWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionTypeID uuid.UUID,
) (entity.ChatRoomActionType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomActionType{}, store.ErrNotFoundDescriptor
	}
	return findChatRoomActionTypeByID(ctx, qtx, chatRoomActionTypeID)
}

func findChatRoomActionTypeByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.ChatRoomActionType, error) {
	e, err := qtx.FindChatRoomActionTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomActionType{}, errhandle.NewModelNotFoundError("chat room action type")
		}
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	entity := entity.ChatRoomActionType{
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		Name:                 e.Name,
		Key:                  e.Key,
	}
	return entity, nil
}

// FindChatRoomActionTypeByKey チャットルームアクションタイプを取得する。
func (a *PgAdapter) FindChatRoomActionTypeByKey(ctx context.Context, key string) (entity.ChatRoomActionType, error) {
	return findChatRoomActionTypeByKey(ctx, a.query, key)
}

// FindChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを取得する。
func (a *PgAdapter) FindChatRoomActionTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.ChatRoomActionType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomActionType{}, store.ErrNotFoundDescriptor
	}
	return findChatRoomActionTypeByKey(ctx, qtx, key)
}

// ChatRoomActionTypeCursor is a cursor for ChatRoomActionType.
type ChatRoomActionTypeCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getChatRoomActionTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereChatRoomActionTypeParam,
	order parameter.ChatRoomActionTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	eConvFunc := func(e query.ChatRoomActionType) (entity.ChatRoomActionType, error) {
		return entity.ChatRoomActionType{
			ChatRoomActionTypeID: e.ChatRoomActionTypeID,
			Name:                 e.Name,
			Key:                  e.Key,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountChatRoomActionTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountChatRoomActionTypes(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room action types: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.ChatRoomActionType, error) {
		p := query.GetChatRoomActionTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetChatRoomActionTypes(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.ChatRoomActionType{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room action types: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.ChatRoomActionType, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.ChatRoomActionTypeNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetChatRoomActionTypesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetChatRoomActionTypesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room action types: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.ChatRoomActionType, error) {
		p := query.GetChatRoomActionTypesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetChatRoomActionTypesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room action types: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.ChatRoomActionType) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomActionTypeDefaultCursorKey:
			return entity.Int(e.MChatRoomActionTypesPkey), nil
		case parameter.ChatRoomActionTypeNameCursorKey:
			return entity.Int(e.MChatRoomActionTypesPkey), e.Name
		}
		return entity.Int(e.MChatRoomActionTypesPkey), nil
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
		return store.ListResult[entity.ChatRoomActionType]{}, fmt.Errorf("failed to get chat room action types: %w", err)
	}
	return res, nil
}

// GetChatRoomActionTypes チャットルームアクションタイプを取得する。
func (a *PgAdapter) GetChatRoomActionTypes(
	ctx context.Context,
	where parameter.WhereChatRoomActionTypeParam,
	order parameter.ChatRoomActionTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	return getChatRoomActionTypes(ctx, a.query, where, order, np, cp, wc)
}

// GetChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを取得する。
func (a *PgAdapter) GetChatRoomActionTypesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereChatRoomActionTypeParam,
	order parameter.ChatRoomActionTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomActionType]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomActionTypes(ctx, qtx, where, order, np, cp, wc)
}

func getPluralChatRoomActionTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID,
	order parameter.ChatRoomActionTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	var e []query.ChatRoomActionType
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomActionTypes(ctx, query.GetPluralChatRoomActionTypesParams{
			ChatRoomActionTypeIds: ids,
			OrderMethod:           order.GetStringValue(),
		})
	} else {
		p := query.GetPluralChatRoomActionTypesUseNumberedPaginateParams{
			ChatRoomActionTypeIds: ids,
			Offset:                int32(np.Offset.Int64),
			Limit:                 int32(np.Limit.Int64),
		}
		e, err = qtx.GetPluralChatRoomActionTypesUseNumberedPaginate(ctx, p)
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomActionType]{},
			fmt.Errorf("failed to get plural chat room action types: %w", err)
	}
	entities := make([]entity.ChatRoomActionType, len(e))
	for i, v := range e {
		entities[i] = entity.ChatRoomActionType{
			ChatRoomActionTypeID: v.ChatRoomActionTypeID,
			Name:                 v.Name,
			Key:                  v.Key,
		}
	}
	return store.ListResult[entity.ChatRoomActionType]{Data: entities}, nil
}

// GetPluralChatRoomActionTypes チャットルームアクションタイプを取得する。
func (a *PgAdapter) GetPluralChatRoomActionTypes(
	ctx context.Context, ids []uuid.UUID, order parameter.ChatRoomActionTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	return getPluralChatRoomActionTypes(ctx, a.query, ids, order, np)
}

// GetPluralChatRoomActionTypesWithSd SD付きでチャットルームアクションタイプを取得する。
func (a *PgAdapter) GetPluralChatRoomActionTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID,
	order parameter.ChatRoomActionTypeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomActionType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomActionTypes(ctx, qtx, ids, order, np)
}

func updateChatRoomActionType(
	ctx context.Context, qtx *query.Queries,
	chatRoomActionTypeID uuid.UUID, param parameter.UpdateChatRoomActionTypeParams,
) (entity.ChatRoomActionType, error) {
	p := query.UpdateChatRoomActionTypeParams{
		ChatRoomActionTypeID: chatRoomActionTypeID,
		Name:                 param.Name,
		Key:                  param.Key,
	}
	e, err := qtx.UpdateChatRoomActionType(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomActionType{}, errhandle.NewModelNotFoundError("chat room action type")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomActionType{}, errhandle.NewModelDuplicatedError("chat room action type")
		}
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to update chat room action type: %w", err)
	}
	entity := entity.ChatRoomActionType{
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		Name:                 e.Name,
		Key:                  e.Key,
	}
	return entity, nil
}

// UpdateChatRoomActionType チャットルームアクションタイプを更新する。
func (a *PgAdapter) UpdateChatRoomActionType(
	ctx context.Context, chatRoomActionTypeID uuid.UUID, param parameter.UpdateChatRoomActionTypeParams,
) (entity.ChatRoomActionType, error) {
	return updateChatRoomActionType(ctx, a.query, chatRoomActionTypeID, param)
}

// UpdateChatRoomActionTypeWithSd SD付きでチャットルームアクションタイプを更新する。
func (a *PgAdapter) UpdateChatRoomActionTypeWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionTypeID uuid.UUID, param parameter.UpdateChatRoomActionTypeParams,
) (entity.ChatRoomActionType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomActionType{}, store.ErrNotFoundDescriptor
	}
	return updateChatRoomActionType(ctx, qtx, chatRoomActionTypeID, param)
}

func updateChatRoomActionTypeByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateChatRoomActionTypeByKeyParams,
) (entity.ChatRoomActionType, error) {
	p := query.UpdateChatRoomActionTypeByKeyParams{
		Key:  key,
		Name: param.Name,
	}
	e, err := qtx.UpdateChatRoomActionTypeByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomActionType{}, errhandle.NewModelNotFoundError("chat room action type")
		}
		return entity.ChatRoomActionType{}, fmt.Errorf("failed to update chat room action type: %w", err)
	}
	entity := entity.ChatRoomActionType{
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		Name:                 e.Name,
		Key:                  e.Key,
	}
	return entity, nil
}

// UpdateChatRoomActionTypeByKey チャットルームアクションタイプを更新する。
func (a *PgAdapter) UpdateChatRoomActionTypeByKey(
	ctx context.Context, key string, param parameter.UpdateChatRoomActionTypeByKeyParams,
) (entity.ChatRoomActionType, error) {
	return updateChatRoomActionTypeByKey(ctx, a.query, key, param)
}

// UpdateChatRoomActionTypeByKeyWithSd SD付きでチャットルームアクションタイプを更新する。
func (a *PgAdapter) UpdateChatRoomActionTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateChatRoomActionTypeByKeyParams,
) (entity.ChatRoomActionType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomActionType{}, store.ErrNotFoundDescriptor
	}
	return updateChatRoomActionTypeByKey(ctx, qtx, key, param)
}

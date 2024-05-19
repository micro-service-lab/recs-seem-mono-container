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

func countAttendStatuses(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttendStatusParam,
) (int64, error) {
	p := query.CountAttendStatusesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountAttendStatuses(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count attend status: %w", err)
	}
	return c, nil
}

// CountAttendStatuses 出席ステータス数を取得する。
func (a *PgAdapter) CountAttendStatuses(ctx context.Context, where parameter.WhereAttendStatusParam) (int64, error) {
	return countAttendStatuses(ctx, a.query, where)
}

// CountAttendStatusesWithSd SD付きで出席ステータス数を取得する。
func (a *PgAdapter) CountAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttendStatusParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countAttendStatuses(ctx, qtx, where)
}

func createAttendStatus(
	ctx context.Context, qtx *query.Queries, param parameter.CreateAttendStatusParam,
) (entity.AttendStatus, error) {
	p := query.CreateAttendStatusParams{
		Name: param.Name,
		Key:  param.Key,
	}
	e, err := qtx.CreateAttendStatus(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.AttendStatus{}, errhandle.NewModelDuplicatedError("attend status")
		}
		return entity.AttendStatus{}, fmt.Errorf("failed to create attend status: %w", err)
	}
	entity := entity.AttendStatus{
		AttendStatusID: e.AttendStatusID,
		Name:           e.Name,
		Key:            e.Key,
	}
	return entity, nil
}

// CreateAttendStatus 出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatus(
	ctx context.Context, param parameter.CreateAttendStatusParam,
) (entity.AttendStatus, error) {
	return createAttendStatus(ctx, a.query, param)
}

// CreateAttendStatusWithSd SD付きで出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatusWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAttendStatusParam,
) (entity.AttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	return createAttendStatus(ctx, qtx, param)
}

func createAttendStatuses(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateAttendStatusParam,
) (int64, error) {
	p := make([]query.CreateAttendStatusesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateAttendStatusesParams{
			Name: param.Name,
			Key:  param.Key,
		}
	}
	c, err := qtx.CreateAttendStatuses(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("attend status")
		}
		return 0, fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return c, nil
}

// CreateAttendStatuses 出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatuses(
	ctx context.Context, params []parameter.CreateAttendStatusParam,
) (int64, error) {
	return createAttendStatuses(ctx, a.query, params)
}

// CreateAttendStatusesWithSd SD付きで出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAttendStatusParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createAttendStatuses(ctx, qtx, params)
}

func deleteAttendStatus(ctx context.Context, qtx *query.Queries, attendStatusID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAttendStatus(ctx, attendStatusID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attend status: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("attend status")
	}
	return c, nil
}

// DeleteAttendStatus 出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatus(ctx context.Context, attendStatusID uuid.UUID) (int64, error) {
	return deleteAttendStatus(ctx, a.query, attendStatusID)
}

// DeleteAttendStatusWithSd SD付きで出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteAttendStatus(ctx, qtx, attendStatusID)
}

func deleteAttendStatusByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteAttendStatusByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attend status: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("attend status")
	}
	return c, nil
}

// DeleteAttendStatusByKey 出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusByKey(ctx context.Context, key string) (int64, error) {
	return deleteAttendStatusByKey(ctx, a.query, key)
}

// DeleteAttendStatusByKeyWithSd SD付きで出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteAttendStatusByKey(ctx, qtx, key)
}

func pluralDeleteAttendStatuses(ctx context.Context, qtx *query.Queries, attendStatusIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteAttendStatuses(ctx, attendStatusIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attend statuses: %w", err)
	}
	if c != int64(len(attendStatusIDs)) {
		return 0, errhandle.NewModelNotFoundError("attend status")
	}
	return c, nil
}

// PluralDeleteAttendStatuses 出席ステータスを複数削除する。
func (a *PgAdapter) PluralDeleteAttendStatuses(ctx context.Context, attendStatusIDs []uuid.UUID) (int64, error) {
	return pluralDeleteAttendStatuses(ctx, a.query, attendStatusIDs)
}

// PluralDeleteAttendStatusesWithSd SD付きで出席ステータスを複数削除する。
func (a *PgAdapter) PluralDeleteAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, attendStatusIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteAttendStatuses(ctx, qtx, attendStatusIDs)
}

func findAttendStatusByID(
	ctx context.Context, qtx *query.Queries, attendStatusID uuid.UUID,
) (entity.AttendStatus, error) {
	e, err := qtx.FindAttendStatusByID(ctx, attendStatusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendStatus{}, errhandle.NewModelNotFoundError("attend status")
		}
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	entity := entity.AttendStatus{
		AttendStatusID: e.AttendStatusID,
		Name:           e.Name,
		Key:            e.Key,
	}
	return entity, nil
}

// FindAttendStatusByID 出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByID(
	ctx context.Context, attendStatusID uuid.UUID,
) (entity.AttendStatus, error) {
	return findAttendStatusByID(ctx, a.query, attendStatusID)
}

// FindAttendStatusByIDWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByIDWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID,
) (entity.AttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	return findAttendStatusByID(ctx, qtx, attendStatusID)
}

func findAttendStatusByKey(ctx context.Context, qtx *query.Queries, key string) (entity.AttendStatus, error) {
	e, err := qtx.FindAttendStatusByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendStatus{}, errhandle.NewModelNotFoundError("attend status")
		}
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	entity := entity.AttendStatus{
		AttendStatusID: e.AttendStatusID,
		Name:           e.Name,
		Key:            e.Key,
	}
	return entity, nil
}

// FindAttendStatusByKey 出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByKey(ctx context.Context, key string) (entity.AttendStatus, error) {
	return findAttendStatusByKey(ctx, a.query, key)
}

// FindAttendStatusByKeyWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.AttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	return findAttendStatusByKey(ctx, qtx, key)
}

// AttendStatusCursor is a cursor for AttendStatus.
type AttendStatusCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getAttendStatuses(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttendStatusParam,
	order parameter.AttendStatusOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendStatus], error) {
	eConvFunc := func(e query.AttendStatus) (entity.AttendStatus, error) {
		return entity.AttendStatus{
			AttendStatusID: e.AttendStatusID,
			Name:           e.Name,
			Key:            e.Key,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttendStatusesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountAttendStatuses(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attend statuses: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.AttendStatus, error) {
		p := query.GetAttendStatusesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetAttendStatuses(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attend statuses: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.AttendStatus, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.AttendStatusNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetAttendStatusesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetAttendStatusesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attend statuses: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.AttendStatus, error) {
		p := query.GetAttendStatusesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetAttendStatusesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attend statuses: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.AttendStatus) (entity.Int, any) {
		switch subCursor {
		case parameter.AttendStatusDefaultCursorKey:
			return entity.Int(e.MAttendStatusesPkey), nil
		case parameter.AttendStatusNameCursorKey:
			return entity.Int(e.MAttendStatusesPkey), e.Name
		}
		return entity.Int(e.MAttendStatusesPkey), nil
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
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	return res, nil
}

// GetAttendStatuses 出席ステータスを取得する。
func (a *PgAdapter) GetAttendStatuses(
	ctx context.Context,
	where parameter.WhereAttendStatusParam,
	order parameter.AttendStatusOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendStatus], error) {
	return getAttendStatuses(ctx, a.query, where, order, np, cp, wc)
}

// GetAttendStatusesWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) GetAttendStatusesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereAttendStatusParam,
	order parameter.AttendStatusOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendStatus], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendStatus]{}, store.ErrNotFoundDescriptor
	}
	return getAttendStatuses(ctx, qtx, where, order, np, cp, wc)
}

func getPluralAttendStatuses(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID,
	order parameter.AttendStatusOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	var e []query.AttendStatus
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralAttendStatuses(ctx, query.GetPluralAttendStatusesParams{
			AttendStatusIds: ids,
			OrderMethod:     order.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralAttendStatusesUseNumberedPaginate(ctx, query.GetPluralAttendStatusesUseNumberedPaginateParams{
			AttendStatusIds: ids,
			Offset:          int32(np.Offset.Int64),
			Limit:           int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get plural attend statuses: %w", err)
	}
	entities := make([]entity.AttendStatus, len(e))
	for i, v := range e {
		entities[i] = entity.AttendStatus{
			AttendStatusID: v.AttendStatusID,
			Name:           v.Name,
			Key:            v.Key,
		}
	}
	return store.ListResult[entity.AttendStatus]{Data: entities}, nil
}

// GetPluralAttendStatuses 出席ステータスを取得する。
func (a *PgAdapter) GetPluralAttendStatuses(
	ctx context.Context, ids []uuid.UUID, order parameter.AttendStatusOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	return getPluralAttendStatuses(ctx, a.query, ids, order, np)
}

// GetPluralAttendStatusesWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) GetPluralAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID,
	order parameter.AttendStatusOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendStatus]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAttendStatuses(ctx, qtx, ids, order, np)
}

func updateAttendStatus(
	ctx context.Context, qtx *query.Queries, attendStatusID uuid.UUID, param parameter.UpdateAttendStatusParams,
) (entity.AttendStatus, error) {
	p := query.UpdateAttendStatusParams{
		AttendStatusID: attendStatusID,
		Name:           param.Name,
		Key:            param.Key,
	}
	e, err := qtx.UpdateAttendStatus(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendStatus{}, errhandle.NewModelNotFoundError("attend status")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.AttendStatus{}, errhandle.NewModelDuplicatedError("attend status")
		}
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	entity := entity.AttendStatus{
		AttendStatusID: e.AttendStatusID,
		Name:           e.Name,
		Key:            e.Key,
	}
	return entity, nil
}

// UpdateAttendStatus 出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatus(
	ctx context.Context, attendStatusID uuid.UUID, param parameter.UpdateAttendStatusParams,
) (entity.AttendStatus, error) {
	return updateAttendStatus(ctx, a.query, attendStatusID, param)
}

// UpdateAttendStatusWithSd SD付きで出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatusWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID, param parameter.UpdateAttendStatusParams,
) (entity.AttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	return updateAttendStatus(ctx, qtx, attendStatusID, param)
}

func updateAttendStatusByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateAttendStatusByKeyParams,
) (entity.AttendStatus, error) {
	p := query.UpdateAttendStatusByKeyParams{
		Key:  key,
		Name: param.Name,
	}
	e, err := qtx.UpdateAttendStatusByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendStatus{}, errhandle.NewModelNotFoundError("attend status")
		}
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	entity := entity.AttendStatus{
		AttendStatusID: e.AttendStatusID,
		Name:           e.Name,
		Key:            e.Key,
	}
	return entity, nil
}

// UpdateAttendStatusByKey 出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatusByKey(
	ctx context.Context, key string, param parameter.UpdateAttendStatusByKeyParams,
) (entity.AttendStatus, error) {
	return updateAttendStatusByKey(ctx, a.query, key, param)
}

// UpdateAttendStatusByKeyWithSd SD付きで出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateAttendStatusByKeyParams,
) (entity.AttendStatus, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	return updateAttendStatusByKey(ctx, qtx, key, param)
}

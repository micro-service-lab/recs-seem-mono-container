package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
)

func countAttendanceTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttendanceTypeParam,
) (int64, error) {
	p := query.CountAttendanceTypesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
	}
	c, err := qtx.CountAttendanceTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendance type: %w", err)
	}
	return c, nil
}

// CountAttendanceTypes 出欠状況タイプ数を取得する。
func (a *PgAdapter) CountAttendanceTypes(ctx context.Context, where parameter.WhereAttendanceTypeParam) (int64, error) {
	c, err := countAttendanceTypes(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendance type: %w", err)
	}
	return c, nil
}

// CountAttendanceTypesWithSd SD付きで出欠状況タイプ数を取得する。
func (a *PgAdapter) CountAttendanceTypesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttendanceTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countAttendanceTypes(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendance type: %w", err)
	}
	return c, nil
}

func createAttendanceType(
	ctx context.Context, qtx *query.Queries, param parameter.CreateAttendanceTypeParam,
) (entity.AttendanceType, error) {
	p := query.CreateAttendanceTypeParams{
		Name:  param.Name,
		Key:   param.Key,
		Color: param.Color,
	}
	e, err := qtx.CreateAttendanceType(ctx, p)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to create attendance type: %w", err)
	}
	entity := entity.AttendanceType{
		AttendanceTypeID: e.AttendanceTypeID,
		Name:             e.Name,
		Key:              e.Key,
		Color:            e.Color,
	}
	return entity, nil
}

// CreateAttendanceType 出欠状況タイプを作成する。
func (a *PgAdapter) CreateAttendanceType(
	ctx context.Context, param parameter.CreateAttendanceTypeParam,
) (entity.AttendanceType, error) {
	e, err := createAttendanceType(ctx, a.query, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to create attendance type: %w", err)
	}
	return e, nil
}

// CreateAttendanceTypeWithSd SD付きで出欠状況タイプを作成する。
func (a *PgAdapter) CreateAttendanceTypeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAttendanceTypeParam,
) (entity.AttendanceType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendanceType{}, store.ErrNotFoundDescriptor
	}
	e, err := createAttendanceType(ctx, qtx, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to create attendance type: %w", err)
	}
	return e, nil
}

func createAttendanceTypes(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateAttendanceTypeParam,
) (int64, error) {
	p := make([]query.CreateAttendanceTypesParams, len(params))
	for i, param := range params {
		p[i] = query.CreateAttendanceTypesParams{
			Name:  param.Name,
			Key:   param.Key,
			Color: param.Color,
		}
	}
	c, err := qtx.CreateAttendanceTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create attendance types: %w", err)
	}
	return c, nil
}

// CreateAttendanceTypes 出欠状況タイプを作成する。
func (a *PgAdapter) CreateAttendanceTypes(
	ctx context.Context, params []parameter.CreateAttendanceTypeParam,
) (int64, error) {
	c, err := createAttendanceTypes(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create attendance types: %w", err)
	}
	return c, nil
}

// CreateAttendanceTypesWithSd SD付きで出欠状況タイプを作成する。
func (a *PgAdapter) CreateAttendanceTypesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAttendanceTypeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createAttendanceTypes(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create attendance types: %w", err)
	}
	return c, nil
}

func deleteAttendanceType(ctx context.Context, qtx *query.Queries, attendanceTypeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAttendanceType(ctx, attendanceTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

// DeleteAttendanceType 出欠状況タイプを削除する。
func (a *PgAdapter) DeleteAttendanceType(ctx context.Context, attendanceTypeID uuid.UUID) (int64, error) {
	c, err := deleteAttendanceType(ctx, a.query, attendanceTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

// DeleteAttendanceTypeWithSd SD付きで出欠状況タイプを削除する。
func (a *PgAdapter) DeleteAttendanceTypeWithSd(
	ctx context.Context, sd store.Sd, attendanceTypeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := deleteAttendanceType(ctx, qtx, attendanceTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

func deleteAttendanceTypeByKey(ctx context.Context, qtx *query.Queries, key string) (int64, error) {
	c, err := qtx.DeleteAttendanceTypeByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

// DeleteAttendanceTypeByKey 出欠状況タイプを削除する。
func (a *PgAdapter) DeleteAttendanceTypeByKey(ctx context.Context, key string) (int64, error) {
	c, err := deleteAttendanceTypeByKey(ctx, a.query, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

// DeleteAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを削除する。
func (a *PgAdapter) DeleteAttendanceTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := deleteAttendanceTypeByKey(ctx, qtx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attendance type: %w", err)
	}
	return c, nil
}

func pluralDeleteAttendanceTypes(ctx context.Context, qtx *query.Queries, attendanceTypeIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteAttendanceTypes(ctx, attendanceTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete attendance types: %w", err)
	}
	return c, nil
}

// PluralDeleteAttendanceTypes 出欠状況タイプを複数削除する。
func (a *PgAdapter) PluralDeleteAttendanceTypes(ctx context.Context, attendanceTypeIDs []uuid.UUID) (int64, error) {
	c, err := pluralDeleteAttendanceTypes(ctx, a.query, attendanceTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete attendance types: %w", err)
	}
	return c, nil
}

// PluralDeleteAttendanceTypesWithSd SD付きで出欠状況タイプを複数削除する。
func (a *PgAdapter) PluralDeleteAttendanceTypesWithSd(
	ctx context.Context, sd store.Sd, attendanceTypeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := pluralDeleteAttendanceTypes(ctx, qtx, attendanceTypeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete attendance types: %w", err)
	}
	return c, nil
}

func findAttendanceTypeByID(
	ctx context.Context, qtx *query.Queries, attendanceTypeID uuid.UUID,
) (entity.AttendanceType, error) {
	e, err := qtx.FindAttendanceTypeByID(ctx, attendanceTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendanceType{}, errhandle.NewModelNotFoundError("attendance type")
		}
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	entity := entity.AttendanceType{
		AttendanceTypeID: e.AttendanceTypeID,
		Name:             e.Name,
		Key:              e.Key,
		Color:            e.Color,
	}
	return entity, nil
}

// FindAttendanceTypeByID 出欠状況タイプを取得する。
func (a *PgAdapter) FindAttendanceTypeByID(
	ctx context.Context, attendanceTypeID uuid.UUID,
) (entity.AttendanceType, error) {
	e, err := findAttendanceTypeByID(ctx, a.query, attendanceTypeID)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	return e, nil
}

// FindAttendanceTypeByIDWithSd SD付きで出欠状況タイプを取得する。
func (a *PgAdapter) FindAttendanceTypeByIDWithSd(
	ctx context.Context, sd store.Sd, attendanceTypeID uuid.UUID,
) (entity.AttendanceType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendanceType{}, store.ErrNotFoundDescriptor
	}
	e, err := findAttendanceTypeByID(ctx, qtx, attendanceTypeID)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	return e, nil
}

func findAttendanceTypeByKey(ctx context.Context, qtx *query.Queries, key string) (entity.AttendanceType, error) {
	e, err := qtx.FindAttendanceTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendanceType{}, errhandle.NewModelNotFoundError("attendance type")
		}
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	entity := entity.AttendanceType{
		AttendanceTypeID: e.AttendanceTypeID,
		Name:             e.Name,
		Key:              e.Key,
		Color:            e.Color,
	}
	return entity, nil
}

// FindAttendanceTypeByKey 出欠状況タイプを取得する。
func (a *PgAdapter) FindAttendanceTypeByKey(ctx context.Context, key string) (entity.AttendanceType, error) {
	e, err := findAttendanceTypeByKey(ctx, a.query, key)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	return e, nil
}

// FindAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを取得する。
func (a *PgAdapter) FindAttendanceTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.AttendanceType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendanceType{}, store.ErrNotFoundDescriptor
	}
	e, err := findAttendanceTypeByKey(ctx, qtx, key)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to find attendance type: %w", err)
	}
	return e, nil
}

// AttendanceTypeCursor is a cursor for AttendanceType.
type AttendanceTypeCursor struct {
	CursorID         int32
	NameCursor       string
	CursorPointsNext bool
}

func getAttendanceTypes(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttendanceTypeParam,
	order parameter.AttendanceTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendanceType], error) {
	eConvFunc := func(e query.AttendanceType) (entity.AttendanceType, error) {
		return entity.AttendanceType{
			AttendanceTypeID: e.AttendanceTypeID,
			Name:             e.Name,
			Key:              e.Key,
			Color:            e.Color,
		}, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttendanceTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		r, err := qtx.CountAttendanceTypes(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attendance types: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.AttendanceType, error) {
		p := query.GetAttendanceTypesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
		}
		r, err := qtx.GetAttendanceTypes(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.AttendanceType{}, nil
			}
			return nil, fmt.Errorf("failed to get attendance types: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.AttendanceType, error) {
		var nameCursor string
		var ok bool
		switch subCursor {
		case parameter.AttendanceTypeNameCursorKey:
			nameCursor, ok = subCursorValue.(string)
			if !ok {
				nameCursor = ""
			}
		}
		p := query.GetAttendanceTypesUseKeysetPaginateParams{
			WhereLikeName:   where.WhereLikeName,
			SearchName:      where.SearchName,
			OrderMethod:     orderMethod,
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			NameCursor:      nameCursor,
		}
		r, err := qtx.GetAttendanceTypesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attendance types: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.AttendanceType, error) {
		p := query.GetAttendanceTypesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   orderMethod,
			Offset:        offset,
			Limit:         limit,
		}
		r, err := qtx.GetAttendanceTypesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attendance types: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.AttendanceType) (entity.Int, any) {
		switch subCursor {
		case parameter.AttendanceTypeDefaultCursorKey:
			return entity.Int(e.MAttendanceTypesPkey), nil
		case parameter.AttendanceTypeNameCursorKey:
			return entity.Int(e.MAttendanceTypesPkey), e.Name
		}
		return entity.Int(e.MAttendanceTypesPkey), nil
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
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get attendance types: %w", err)
	}
	return res, nil
}

// GetAttendanceTypes 出欠状況タイプを取得する。
func (a *PgAdapter) GetAttendanceTypes(
	ctx context.Context,
	where parameter.WhereAttendanceTypeParam,
	order parameter.AttendanceTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendanceType], error) {
	r, err := getAttendanceTypes(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get attendance types: %w", err)
	}
	return r, nil
}

// GetAttendanceTypesWithSd SD付きで出欠状況タイプを取得する。
func (a *PgAdapter) GetAttendanceTypesWithSd(
	ctx context.Context,
	sd store.Sd,
	where parameter.WhereAttendanceTypeParam,
	order parameter.AttendanceTypeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttendanceType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendanceType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getAttendanceTypes(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get attendance types: %w", err)
	}
	return r, nil
}

func getPluralAttendanceTypes(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendanceType], error) {
	p := query.GetPluralAttendanceTypesParams{
		AttendanceTypeIds: ids,
		Offset:            int32(np.Offset.Int64),
		Limit:             int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralAttendanceTypes(ctx, p)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get plural attendance types: %w", err)
	}
	entities := make([]entity.AttendanceType, len(e))
	for i, v := range e {
		entities[i] = entity.AttendanceType{
			AttendanceTypeID: v.AttendanceTypeID,
			Name:             v.Name,
			Key:              v.Key,
			Color:            v.Color,
		}
	}
	return store.ListResult[entity.AttendanceType]{Data: entities}, nil
}

// GetPluralAttendanceTypes 出欠状況タイプを取得する。
func (a *PgAdapter) GetPluralAttendanceTypes(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendanceType], error) {
	r, err := getPluralAttendanceTypes(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get plural attendance types: %w", err)
	}
	return r, nil
}

// GetPluralAttendanceTypesWithSd SD付きで出欠状況タイプを取得する。
func (a *PgAdapter) GetPluralAttendanceTypesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendanceType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendanceType]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralAttendanceTypes(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.AttendanceType]{}, fmt.Errorf("failed to get plural attendance types: %w", err)
	}
	return r, nil
}

func updateAttendanceType(
	ctx context.Context, qtx *query.Queries, attendanceTypeID uuid.UUID, param parameter.UpdateAttendanceTypeParams,
) (entity.AttendanceType, error) {
	p := query.UpdateAttendanceTypeParams{
		AttendanceTypeID: attendanceTypeID,
		Name:             param.Name,
		Key:              param.Key,
		Color:            param.Color,
	}
	e, err := qtx.UpdateAttendanceType(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendanceType{}, errhandle.NewModelNotFoundError("attendance type")
		}
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	entity := entity.AttendanceType{
		AttendanceTypeID: e.AttendanceTypeID,
		Name:             e.Name,
		Key:              e.Key,
		Color:            e.Color,
	}
	return entity, nil
}

// UpdateAttendanceType 出欠状況タイプを更新する。
func (a *PgAdapter) UpdateAttendanceType(
	ctx context.Context, attendanceTypeID uuid.UUID, param parameter.UpdateAttendanceTypeParams,
) (entity.AttendanceType, error) {
	e, err := updateAttendanceType(ctx, a.query, attendanceTypeID, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	return e, nil
}

// UpdateAttendanceTypeWithSd SD付きで出欠状況タイプを更新する。
func (a *PgAdapter) UpdateAttendanceTypeWithSd(
	ctx context.Context, sd store.Sd, attendanceTypeID uuid.UUID, param parameter.UpdateAttendanceTypeParams,
) (entity.AttendanceType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendanceType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateAttendanceType(ctx, qtx, attendanceTypeID, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	return e, nil
}

func updateAttendanceTypeByKey(
	ctx context.Context, qtx *query.Queries, key string, param parameter.UpdateAttendanceTypeByKeyParams,
) (entity.AttendanceType, error) {
	p := query.UpdateAttendanceTypeByKeyParams{
		Key:   key,
		Name:  param.Name,
		Color: param.Color,
	}
	e, err := qtx.UpdateAttendanceTypeByKey(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttendanceType{}, errhandle.NewModelNotFoundError("attendance type")
		}
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	entity := entity.AttendanceType{
		AttendanceTypeID: e.AttendanceTypeID,
		Name:             e.Name,
		Key:              e.Key,
		Color:            e.Color,
	}
	return entity, nil
}

// UpdateAttendanceTypeByKey 出欠状況タイプを更新する。
func (a *PgAdapter) UpdateAttendanceTypeByKey(
	ctx context.Context, key string, param parameter.UpdateAttendanceTypeByKeyParams,
) (entity.AttendanceType, error) {
	e, err := updateAttendanceTypeByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	return e, nil
}

// UpdateAttendanceTypeByKeyWithSd SD付きで出欠状況タイプを更新する。
func (a *PgAdapter) UpdateAttendanceTypeByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateAttendanceTypeByKeyParams,
) (entity.AttendanceType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendanceType{}, store.ErrNotFoundDescriptor
	}
	e, err := updateAttendanceTypeByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.AttendanceType{}, fmt.Errorf("failed to update attendance type: %w", err)
	}
	return e, nil
}

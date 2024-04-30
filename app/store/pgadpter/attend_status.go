package pgadapter

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
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
	c, err := countAttendStatuses(ctx, a.query, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count attend status: %w", err)
	}
	return c, nil
}

// CountAttendStatusesWithSd SD付きで出席ステータス数を取得する。
func (a *PgAdapter) CountAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttendStatusParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countAttendStatuses(ctx, qtx, where)
	if err != nil {
		return 0, fmt.Errorf("failed to count attend status: %w", err)
	}
	return c, nil
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
	e, err := createAttendStatus(ctx, a.query, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to create attend status: %w", err)
	}
	return e, nil
}

// CreateAttendStatusWithSd SD付きで出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatusWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAttendStatusParam,
) (entity.AttendStatus, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	e, err := createAttendStatus(ctx, qtx, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to create attend status: %w", err)
	}
	return e, nil
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
		return 0, fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return c, nil
}

// CreateAttendStatuses 出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatuses(
	ctx context.Context, params []parameter.CreateAttendStatusParam,
) (int64, error) {
	c, err := createAttendStatuses(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return c, nil
}

// CreateAttendStatusesWithSd SD付きで出席ステータスを作成する。
func (a *PgAdapter) CreateAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAttendStatusParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := createAttendStatuses(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create attend statuses: %w", err)
	}
	return c, nil
}

func deleteAttendStatus(ctx context.Context, qtx *query.Queries, attendStatusID uuid.UUID) error {
	err := qtx.DeleteAttendStatus(ctx, attendStatusID)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

// DeleteAttendStatus 出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatus(ctx context.Context, attendStatusID uuid.UUID) error {
	err := deleteAttendStatus(ctx, a.query, attendStatusID)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

// DeleteAttendStatusWithSd SD付きで出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteAttendStatus(ctx, qtx, attendStatusID)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

func deleteAttendStatusByKey(ctx context.Context, qtx *query.Queries, key string) error {
	err := qtx.DeleteAttendStatusByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

// DeleteAttendStatusByKey 出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusByKey(ctx context.Context, key string) error {
	err := deleteAttendStatusByKey(ctx, a.query, key)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

// DeleteAttendStatusByKeyWithSd SD付きで出席ステータスを削除する。
func (a *PgAdapter) DeleteAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteAttendStatusByKey(ctx, qtx, key)
	if err != nil {
		return fmt.Errorf("failed to delete attend status: %w", err)
	}
	return nil
}

func findAttendStatusByID(
	ctx context.Context, qtx *query.Queries, attendStatusID uuid.UUID,
) (entity.AttendStatus, error) {
	e, err := qtx.FindAttendStatusByID(ctx, attendStatusID)
	if err != nil {
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
	e, err := findAttendStatusByID(ctx, a.query, attendStatusID)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	return e, nil
}

// FindAttendStatusByIDWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByIDWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID,
) (entity.AttendStatus, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	e, err := findAttendStatusByID(ctx, qtx, attendStatusID)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	return e, nil
}

func findAttendStatusByKey(ctx context.Context, qtx *query.Queries, key string) (entity.AttendStatus, error) {
	e, err := qtx.FindAttendStatusByKey(ctx, key)
	if err != nil {
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
	e, err := findAttendStatusByKey(ctx, a.query, key)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	return e, nil
}

// FindAttendStatusByKeyWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) FindAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.AttendStatus, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	e, err := findAttendStatusByKey(ctx, qtx, key)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to find attend status: %w", err)
	}
	return e, nil
}

const (
	attendStatusNameCursorKey = "name"
)

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
	fmt.Println("order", order)
	var withCount int64
	if wc.Valid {
		var err error
		p := query.CountAttendStatusesParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
		}
		withCount, err = qtx.CountAttendStatuses(ctx, p)
		if err != nil {
			return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to count attend status: %w", err)
		}
	}
	wcAtr := store.WithCountAttribute{
		Count: withCount,
		Valid: wc.Valid,
	}
	if np.Valid {
		p := query.GetAttendStatusesUseNumberedPaginateParams{
			WhereLikeName: where.WhereLikeName,
			SearchName:    where.SearchName,
			OrderMethod:   string(order),
			Offset:        int32(np.Offset.Int64),
			Limit:         int32(np.Limit.Int64),
		}
		e, err := qtx.GetAttendStatusesUseNumberedPaginate(ctx, p)
		if err != nil {
			return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
		}
		entities := make([]entity.AttendStatus, len(e))
		for i, v := range e {
			entities[i] = entity.AttendStatus{
				AttendStatusID: v.AttendStatusID,
				Name:           v.Name,
				Key:            v.Key,
			}
		}
		return store.ListResult[entity.AttendStatus]{Data: entities, WithCount: wcAtr}, nil
	} else if cp.Valid {
		var err error
		isFirst := cp.Cursor == "" // 初回のリクエストかどうか
		pointsNext := false        // ページネーションの方向(true: 次データ, false: 前データ)
		limit := cp.Limit.Int64
		var e []query.AttendStatus // 結果リストデータ
		var subCursor string
		// サブカーソルに何を使用するか
		switch string(order) {
		case string(parameter.AttendStatusOrderMethodName):
			subCursor = attendStatusNameCursorKey
		case string(parameter.AttendStatusOrderMethodReverseName):
			subCursor = attendStatusNameCursorKey
		}
		// カーソルのデコード+チェック
		var decodedCursor store.Cursor
		var nameCursor string
		var ok bool
		if !isFirst {
			cursorCheck := func(cur string) bool {
				decodedCursor, err = store.DecodeCursor(cur)
				if err != nil {
					return false
				}
				if decodedCursor.SubCursorName != subCursor {
					return false
				}

				switch subCursor {
				case attendStatusNameCursorKey:
					nameCursor, ok = decodedCursor.SubCursor.(string)
				}
				if !ok {
					return false
				}
				return true
			}
			if !cursorCheck(cp.Cursor) {
				isFirst = true
				cp.Cursor = ""
			}
		}

		if !isFirst {
			// 今回の指定カーソルの方向を引き継ぐ
			pointsNext = decodedCursor.CursorPointsNext == true
			var cursorDir string
			if pointsNext {
				cursorDir = "next"
			} else {
				cursorDir = "prev"
			}
			ID := decodedCursor.CursorID
			p := query.GetAttendStatusesUseKeysetPaginateParams{
				WhereLikeName:   where.WhereLikeName,
				SearchName:      where.SearchName,
				OrderMethod:     string(order),
				Limit:           int32(limit) + 1,
				CursorDirection: cursorDir,
				NameCursor:      nameCursor,
				Cursor:          int32(ID),
			}
			e, err = qtx.GetAttendStatusesUseKeysetPaginate(ctx, p)
			if err != nil {
				return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
			}
		} else {
			p := query.GetAttendStatusesUseNumberedPaginateParams{
				WhereLikeName: where.WhereLikeName,
				SearchName:    where.SearchName,
				OrderMethod:   string(order),
				Offset:        0,
				Limit:         int32(limit) + 1,
			}
			var err error
			e, err = qtx.GetAttendStatusesUseNumberedPaginate(ctx, p)
			if err != nil {
				return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
			}
		}
		if len(e) == 0 {
			return store.ListResult[entity.AttendStatus]{Data: []entity.AttendStatus{}, WithCount: wcAtr}, nil
		}
		hasPagination := len(e) > int(limit)
		if hasPagination {
			e = e[:limit]
		}
		eLen := len(e)
		var firstValue, lastValue any
		lastIndex := eLen - 1
		if lastIndex < 0 {
			lastIndex = 0
		}
		switch subCursor {
		case attendStatusNameCursorKey:
			firstValue = e[0].Name
			lastValue = e[lastIndex].Name
		}
		firstData := store.CursorData{
			ID:    entity.Int(e[0].MAttendStatusesPkey),
			Name:  subCursor,
			Value: firstValue,
		}
		lastData := store.CursorData{
			ID:    entity.Int(e[lastIndex].MAttendStatusesPkey),
			Name:  subCursor,
			Value: lastValue,
		}
		var pageInfo store.CursorPaginationAttribute
		if pointsNext || isFirst {
			pageInfo = store.CalculatePagination(isFirst, hasPagination, pointsNext, firstData, lastData)
		} else {
			pageInfo = store.CalculatePagination(isFirst, hasPagination, pointsNext, lastData, firstData)
		}

		entities := make([]entity.AttendStatus, len(e))
		for i, v := range e {
			entities[i] = entity.AttendStatus{
				AttendStatusID: v.AttendStatusID,
				Name:           v.Name,
				Key:            v.Key,
			}
		}
		return store.ListResult[entity.AttendStatus]{Data: entities, CursorPagination: pageInfo, WithCount: wcAtr}, nil
	}
	p := query.GetAttendStatusesParams{
		WhereLikeName: where.WhereLikeName,
		SearchName:    where.SearchName,
		OrderMethod:   string(order),
	}
	e, err := qtx.GetAttendStatuses(ctx, p)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	entities := make([]entity.AttendStatus, len(e))
	for i, v := range e {
		entities[i] = entity.AttendStatus{
			AttendStatusID: v.AttendStatusID,
			Name:           v.Name,
			Key:            v.Key,
		}
	}
	return store.ListResult[entity.AttendStatus]{Data: entities, WithCount: wcAtr}, nil
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
	r, err := getAttendStatuses(ctx, a.query, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	return r, nil
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
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendStatus]{}, store.ErrNotFoundDescriptor
	}
	r, err := getAttendStatuses(ctx, qtx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get attend statuses: %w", err)
	}
	return r, nil
}

func getPluralAttendStatuses(
	ctx context.Context, qtx *query.Queries, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	uids := make([]uuid.UUID, len(ids))
	for i, v := range ids {
		uids[i] = v
	}
	p := query.GetPluralAttendStatusesParams{
		AttendStatusIds: uids,
		Offset:          int32(np.Offset.Int64),
		Limit:           int32(np.Limit.Int64),
	}
	e, err := qtx.GetPluralAttendStatuses(ctx, p)
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
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	r, err := getPluralAttendStatuses(ctx, a.query, ids, np)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get plural attend statuses: %w", err)
	}
	return r, nil
}

// GetPluralAttendStatusesWithSd SD付きで出席ステータスを取得する。
func (a *PgAdapter) GetPluralAttendStatusesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttendStatus], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttendStatus]{}, store.ErrNotFoundDescriptor
	}
	r, err := getPluralAttendStatuses(ctx, qtx, ids, np)
	if err != nil {
		return store.ListResult[entity.AttendStatus]{}, fmt.Errorf("failed to get plural attend statuses: %w", err)
	}
	return r, nil
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
	e, err := updateAttendStatus(ctx, a.query, attendStatusID, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	return e, nil
}

// UpdateAttendStatusWithSd SD付きで出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatusWithSd(
	ctx context.Context, sd store.Sd, attendStatusID uuid.UUID, param parameter.UpdateAttendStatusParams,
) (entity.AttendStatus, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	e, err := updateAttendStatus(ctx, qtx, attendStatusID, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	return e, nil
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
	e, err := updateAttendStatusByKey(ctx, a.query, key, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	return e, nil
}

// UpdateAttendStatusByKeyWithSd SD付きで出席ステータスを更新する。
func (a *PgAdapter) UpdateAttendStatusByKeyWithSd(
	ctx context.Context, sd store.Sd, key string, param parameter.UpdateAttendStatusByKeyParams,
) (entity.AttendStatus, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttendStatus{}, store.ErrNotFoundDescriptor
	}
	e, err := updateAttendStatusByKey(ctx, qtx, key, param)
	if err != nil {
		return entity.AttendStatus{}, fmt.Errorf("failed to update attend status: %w", err)
	}
	return e, nil
}

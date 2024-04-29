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

func countAbsences(ctx context.Context, qtx *query.Queries) (int64, error) {
	c, err := qtx.CountAbsences(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count absences: %w", err)
	}
	return c, nil
}

// CountAbsences 欠席数を取得する。
func (a *PgAdapter) CountAbsences(ctx context.Context) (int64, error) {
	c, err := countAbsences(ctx, a.query)
	if err != nil {
		return 0, fmt.Errorf("failed to count absences: %w", err)
	}
	return c, nil
}

// CountAbsencesWithSd SD付きで欠席数を取得する。
func (a *PgAdapter) CountAbsencesWithSd(ctx context.Context, sd store.Sd) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	c, err := countAbsences(ctx, qtx)
	if err != nil {
		return 0, fmt.Errorf("failed to count absences: %w", err)
	}
	return c, nil
}

func createAbsence(ctx context.Context, qtx *query.Queries, param parameter.CreateAbsenceParam) (entity.Absence, error) {
	e, err := qtx.CreateAbsence(ctx, param.AttendanceID)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to create absence: %w", err)
	}
	entity := entity.Absence{
		AbsenceID:    e.AbsenceID,
		AttendanceID: e.AttendanceID,
	}
	return entity, nil
}

// CreateAbsence 欠席を作成する。
func (a *PgAdapter) CreateAbsence(ctx context.Context, param parameter.CreateAbsenceParam) (entity.Absence, error) {
	e, err := createAbsence(ctx, a.query, param)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to create absence: %w", err)
	}
	return e, nil
}

// CreateAbsenceWithSd SD付きで欠席を作成する。
func (a *PgAdapter) CreateAbsenceWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAbsenceParam,
) (entity.Absence, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Absence{}, store.ErrNotFoundDescriptor
	}
	e, err := createAbsence(ctx, qtx, param)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to create absence: %w", err)
	}
	return e, nil
}

func createAbsences(ctx context.Context, qtx *query.Queries, params []parameter.CreateAbsenceParam) (int64, error) {
	aIDs := make([]uuid.UUID, len(params))
	for i, p := range params {
		aIDs[i] = p.AttendanceID
	}
	e, err := qtx.CreateAbsences(ctx, aIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to create absences: %w", err)
	}
	return e, nil
}

// CreateAbsences 欠席を作成する。
func (a *PgAdapter) CreateAbsences(ctx context.Context, params []parameter.CreateAbsenceParam) (int64, error) {
	e, err := createAbsences(ctx, a.query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create absences: %w", err)
	}
	return e, nil
}

// CreateAbsencesWithSd SD付きで欠席を作成する。
func (a *PgAdapter) CreateAbsencesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAbsenceParam,
) (int64, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	e, err := createAbsences(ctx, qtx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create absences: %w", err)
	}
	return e, nil
}

func deleteAbsence(ctx context.Context, qtx *query.Queries, absenceID uuid.UUID) error {
	err := qtx.DeleteAbsence(ctx, absenceID)
	if err != nil {
		return fmt.Errorf("failed to delete absence: %w", err)
	}
	return nil
}

// DeleteAbsence 欠席を削除する。
func (a *PgAdapter) DeleteAbsence(ctx context.Context, absenceID uuid.UUID) error {
	err := deleteAbsence(ctx, a.query, absenceID)
	if err != nil {
		return fmt.Errorf("failed to delete absence: %w", err)
	}
	return nil
}

// DeleteAbsenceWithSd SD付きで欠席を削除する。
func (a *PgAdapter) DeleteAbsenceWithSd(ctx context.Context, sd store.Sd, absenceID uuid.UUID) error {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ErrNotFoundDescriptor
	}
	err := deleteAbsence(ctx, qtx, absenceID)
	if err != nil {
		return fmt.Errorf("failed to delete absence: %w", err)
	}
	return nil
}

func findAbsenceByID(ctx context.Context, qtx *query.Queries, absenceID uuid.UUID) (entity.Absence, error) {
	e, err := qtx.FindAbsenceByID(ctx, absenceID)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to find absence: %w", err)
	}
	entity := entity.Absence{
		AbsenceID:    e.AbsenceID,
		AttendanceID: e.AttendanceID,
	}
	return entity, nil
}

// FindAbsenceByID 欠席を取得する。
func (a *PgAdapter) FindAbsenceByID(ctx context.Context, absenceID uuid.UUID) (entity.Absence, error) {
	e, err := findAbsenceByID(ctx, a.query, absenceID)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to find absence: %w", err)
	}
	return e, nil
}

// FindAbsenceByIDWithSd SD付きで欠席を取得する。
func (a *PgAdapter) FindAbsenceByIDWithSd(
	ctx context.Context, sd store.Sd, absenceID uuid.UUID,
) (entity.Absence, error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Absence{}, store.ErrNotFoundDescriptor
	}
	e, err := findAbsenceByID(ctx, qtx, absenceID)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to find absence: %w", err)
	}
	return e, nil
}

func getAbsencesList(
	ctx context.Context,
	qtx *query.Queries,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	var withCount int64
	if wc.Valid {
		var err error
		withCount, err = qtx.CountAbsences(ctx)
		if err != nil {
			return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to count absences: %w", err)
		}
	}
	wcAtr := store.WithCountAttribute{
		Count: withCount,
	}
	if np.Valid {
		p := query.GetAbsencesUseNumberedPaginateParams{
			Offset: int32(np.Offset.Int64),
			Limit:  int32(np.Limit.Int64),
		}
		ql, err := qtx.GetAbsencesUseNumberedPaginate(ctx, p)
		if err != nil {
			return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get absences: %w", err)
		}
		entities := make([]entity.Absence, len(ql))
		for i, e := range ql {
			entities[i] = entity.Absence{
				AbsenceID:    e.AbsenceID,
				AttendanceID: e.AttendanceID,
			}
		}
		return store.ListResult[entity.Absence]{Data: entities, WithCount: wcAtr}, nil
	} else if cp.Valid {
		isFirst := cp.Cursor == ""
		pointsNext := false
		limit := cp.Limit.Int64
		var ql []query.Absence

		if !isFirst {
			decodedCursor, err := store.DecodeCursor(cp.Cursor)
			if err != nil {
				return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to decode cursor: %w", err)
			}
			pointsNext = decodedCursor.CursorPointsNext == true
			ID := decodedCursor.CursorID
			p := query.GetAbsencesUseKeysetPaginateParams{
				Limit:           int32(limit) + 1,
				CursorDirection: "next",
				Cursor:          int32(ID),
			}
			ql, err = qtx.GetAbsencesUseKeysetPaginate(ctx, p)
			if err != nil {
				return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get absences: %w", err)
			}
		} else {
			p := query.GetAbsencesUseNumberedPaginateParams{
				Offset: 0,
				Limit:  int32(limit) + 1,
			}
			var err error
			ql, err = qtx.GetAbsencesUseNumberedPaginate(ctx, p)
			if err != nil {
				return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get absences: %w", err)
			}
		}
		hasPagination := len(ql) > int(limit)
		if hasPagination {
			ql = ql[:limit]
		}
		firstData := store.CursorData{
			ID: entity.Int(ql[0].TAbsencesPkey),
		}
		lastData := store.CursorData{
			ID: entity.Int(ql[limit-1].TAbsencesPkey),
		}
		pageInfo := store.CalculatePagination(isFirst, hasPagination, pointsNext, firstData, lastData)

		entities := make([]entity.Absence, len(ql))
		for i, e := range ql {
			entities[i] = entity.Absence{
				AbsenceID:    e.AbsenceID,
				AttendanceID: e.AttendanceID,
			}
		}
		return store.ListResult[entity.Absence]{Data: entities, CursorPagination: pageInfo, WithCount: wcAtr}, nil
	}
	ql, err := qtx.GetAbsences(ctx)
	if err != nil {
		return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get absences: %w", err)
	}
	entities := make([]entity.Absence, len(ql))
	for i, e := range ql {
		entities[i] = entity.Absence{
			AbsenceID:    e.AbsenceID,
			AttendanceID: e.AttendanceID,
		}
	}
	return store.ListResult[entity.Absence]{Data: entities, WithCount: wcAtr}, nil
}

// GetAbsences 欠席を取得する。
func (a *PgAdapter) GetAbsences(
	ctx context.Context, np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	return getAbsencesList(ctx, a.query, np, cp, wc)
}

// GetAbsencesWithSd SD付きで欠席を取得する。
func (a *PgAdapter) GetAbsencesWithSd(
	ctx context.Context,
	sd store.Sd,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Absence]{}, store.ErrNotFoundDescriptor
	}
	return getAbsencesList(ctx, qtx, np, cp, wc)
}

func getPluralAbsences(
	ctx context.Context,
	qtx *query.Queries,
	ids []uuid.UUID,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.Absence], error) {
	aIDs := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		aIDs[i] = id
	}
	p := query.GetPluralAbsencesParams{
		AbsenceIds: aIDs,
		Limit:      int32(np.Limit.Int64),
		Offset:     int32(np.Offset.Int64),
	}
	ql, err := qtx.GetPluralAbsences(ctx, p)
	if err != nil {
		return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get plural absences: %w", err)
	}
	entities := make([]entity.Absence, len(ql))
	for i, e := range ql {
		entities[i] = entity.Absence{
			AbsenceID:    e.AbsenceID,
			AttendanceID: e.AttendanceID,
		}
	}
	return store.ListResult[entity.Absence]{Data: entities}, nil
}

// GetPluralAbsences 複数の欠席を取得する。
func (a *PgAdapter) GetPluralAbsences(
	ctx context.Context, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Absence], error) {
	return getPluralAbsences(ctx, a.query, ids, np)
}

// GetPluralAbsencesWithSd SD付きで複数の欠席を取得する。
func (a *PgAdapter) GetPluralAbsencesWithSd(
	ctx context.Context, sd store.Sd, ids []uuid.UUID, np store.NumberedPaginationParam,
) (store.ListResult[entity.Absence], error) {
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Absence]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAbsences(ctx, qtx, ids, np)
}

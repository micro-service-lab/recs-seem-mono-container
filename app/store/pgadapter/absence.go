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

func countAbsences(ctx context.Context, qtx *query.Queries) (int64, error) {
	c, err := qtx.CountAbsences(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count absences: %w", err)
	}
	return c, nil
}

// CountAbsences 欠席数を取得する。
func (a *PgAdapter) CountAbsences(ctx context.Context) (int64, error) {
	return countAbsences(ctx, a.query)
}

// CountAbsencesWithSd SD付きで欠席数を取得する。
func (a *PgAdapter) CountAbsencesWithSd(ctx context.Context, sd store.Sd) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countAbsences(ctx, qtx)
}

func createAbsence(
	ctx context.Context, qtx *query.Queries, param parameter.CreateAbsenceParam,
) (entity.Absence, error) {
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
	return createAbsence(ctx, a.query, param)
}

// CreateAbsenceWithSd SD付きで欠席を作成する。
func (a *PgAdapter) CreateAbsenceWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAbsenceParam,
) (entity.Absence, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Absence{}, store.ErrNotFoundDescriptor
	}
	return createAbsence(ctx, qtx, param)
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
	return createAbsences(ctx, a.query, params)
}

// CreateAbsencesWithSd SD付きで欠席を作成する。
func (a *PgAdapter) CreateAbsencesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAbsenceParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createAbsences(ctx, qtx, params)
}

func deleteAbsence(ctx context.Context, qtx *query.Queries, absenceID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAbsence(ctx, absenceID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete absence: %w", err)
	}
	return c, nil
}

// DeleteAbsence 欠席を削除する。
func (a *PgAdapter) DeleteAbsence(ctx context.Context, absenceID uuid.UUID) (int64, error) {
	return deleteAbsence(ctx, a.query, absenceID)
}

// DeleteAbsenceWithSd SD付きで欠席を削除する。
func (a *PgAdapter) DeleteAbsenceWithSd(ctx context.Context, sd store.Sd, absenceID uuid.UUID) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteAbsence(ctx, qtx, absenceID)
}

func pluralDeleteAbsences(ctx context.Context, qtx *query.Queries, absenceIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteAbsences(ctx, absenceIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete absences: %w", err)
	}
	return c, nil
}

// PluralDeleteAbsences 欠席を複数削除する。
func (a *PgAdapter) PluralDeleteAbsences(ctx context.Context, absenceIDs []uuid.UUID) (int64, error) {
	return pluralDeleteAbsences(ctx, a.query, absenceIDs)
}

// PluralDeleteAbsencesWithSd SD付きで欠席を複数削除する。
func (a *PgAdapter) PluralDeleteAbsencesWithSd(
	ctx context.Context, sd store.Sd, absenceIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteAbsences(ctx, qtx, absenceIDs)
}

func findAbsenceByID(ctx context.Context, qtx *query.Queries, absenceID uuid.UUID) (entity.Absence, error) {
	e, err := qtx.FindAbsenceByID(ctx, absenceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Absence{}, errhandle.NewModelNotFoundError("absence")
		}
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
	return findAbsenceByID(ctx, a.query, absenceID)
}

// FindAbsenceByIDWithSd SD付きで欠席を取得する。
func (a *PgAdapter) FindAbsenceByIDWithSd(
	ctx context.Context, sd store.Sd, absenceID uuid.UUID,
) (entity.Absence, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Absence{}, store.ErrNotFoundDescriptor
	}
	return findAbsenceByID(ctx, qtx, absenceID)
}

func getAbsencesList(
	ctx context.Context,
	qtx *query.Queries,
	order parameter.AbsenceOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	eConvFunc := func(e query.Absence) (entity.Absence, error) {
		return entity.Absence{
			AbsenceID:    e.AbsenceID,
			AttendanceID: e.AttendanceID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountAbsences(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count absences: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Absence, error) {
		r, err := qtx.GetAbsences(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get absences: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(
		_, _ string, limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Absence, error) {
		p := query.GetAbsencesUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetAbsencesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get absences: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(
		_ string, limit, offset int32,
	) ([]query.Absence, error) {
		p := query.GetAbsencesUseNumberedPaginateParams{
			Offset: offset,
			Limit:  limit,
		}
		r, err := qtx.GetAbsencesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get absences: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Absence) (entity.Int, any) {
		switch subCursor {
		case parameter.AbsenceDefaultCursorKey:
			return entity.Int(e.TAbsencesPkey), nil
		}
		return entity.Int(e.TAbsencesPkey), nil
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
		return store.ListResult[entity.Absence]{}, fmt.Errorf("failed to get absences list: %w", err)
	}
	return res, nil
}

// GetAbsences 欠席を取得する。
func (a *PgAdapter) GetAbsences(
	ctx context.Context, order parameter.AbsenceOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	return getAbsencesList(ctx, a.query, order, np, cp, wc)
}

// GetAbsencesWithSd SD付きで欠席を取得する。
func (a *PgAdapter) GetAbsencesWithSd(
	ctx context.Context,
	sd store.Sd,
	order parameter.AbsenceOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Absence], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Absence]{}, store.ErrNotFoundDescriptor
	}
	return getAbsencesList(ctx, qtx, order, np, cp, wc)
}

func getPluralAbsences(
	ctx context.Context,
	qtx *query.Queries,
	ids []uuid.UUID,
	np store.NumberedPaginationParam,
) (store.ListResult[entity.Absence], error) {
	var ql []query.Absence
	var err error
	if !np.Valid {
		ql, err = qtx.GetPluralAbsences(ctx, ids)
	} else {
		p := query.GetPluralAbsencesUseNumberedPaginateParams{
			AbsenceIds: ids,
			Limit:      int32(np.Limit.Int64),
			Offset:     int32(np.Offset.Int64),
		}
		ql, err = qtx.GetPluralAbsencesUseNumberedPaginate(ctx, p)
	}
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
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Absence]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAbsences(ctx, qtx, ids, np)
}

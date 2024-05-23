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

func convGradeWithOrganization(e query.FindGradeByIDWithOrganizationRow) entity.GradeWithOrganization {
	return entity.GradeWithOrganization{
		GradeID: e.GradeID,
		Key:     e.Key,
		Organization: entity.Organization{
			OrganizationID: e.OrganizationID,
			Name:           e.OrganizationName.String,
			Description:    entity.String(e.OrganizationDescription),
			Color:          entity.String(e.OrganizationColor),
			IsPersonal:     e.OrganizationIsPersonal.Bool,
			IsWhole:        e.OrganizationIsWhole.Bool,
			ChatRoomID:     entity.UUID(e.OrganizationChatRoomID),
		},
	}
}

// countGrades は学年数を取得する内部関数です。
func countGrades(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereGradeParam,
) (int64, error) {
	c, err := qtx.CountGrades(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count grades: %w", err)
	}
	return c, nil
}

// CountGrades は学年数を取得します。
func (a *PgAdapter) CountGrades(ctx context.Context, where parameter.WhereGradeParam) (int64, error) {
	return countGrades(ctx, a.query, where)
}

// CountGradesWithSd はSD付きで学年数を取得します。
func (a *PgAdapter) CountGradesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGradeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countGrades(ctx, qtx, where)
}

// createGrade は学年を作成する内部関数です。
func createGrade(
	ctx context.Context, qtx *query.Queries, param parameter.CreateGradeParam,
) (entity.Grade, error) {
	p := query.CreateGradeParams{
		Key:            param.Key,
		OrganizationID: param.OrganizationID,
	}
	e, err := qtx.CreateGrade(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Grade{}, errhandle.NewModelDuplicatedError("grade")
		}
		return entity.Grade{}, fmt.Errorf("failed to create grade: %w", err)
	}
	entity := entity.Grade{
		GradeID:        e.GradeID,
		Key:            e.Key,
		OrganizationID: e.OrganizationID,
	}
	return entity, nil
}

// CreateGrade は学年を作成します。
func (a *PgAdapter) CreateGrade(
	ctx context.Context, param parameter.CreateGradeParam,
) (entity.Grade, error) {
	return createGrade(ctx, a.query, param)
}

// CreateGradeWithSd はSD付きで学年を作成します。
func (a *PgAdapter) CreateGradeWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateGradeParam,
) (entity.Grade, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Grade{}, store.ErrNotFoundDescriptor
	}
	return createGrade(ctx, qtx, param)
}

// createGrades は複数の学年を作成する内部関数です。
func createGrades(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateGradeParam,
) (int64, error) {
	param := make([]query.CreateGradesParams, len(params))
	for i, p := range params {
		param[i] = query.CreateGradesParams{
			Key:            p.Key,
			OrganizationID: p.OrganizationID,
		}
	}
	n, err := qtx.CreateGrades(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("grade")
		}
		return 0, fmt.Errorf("failed to create grades: %w", err)
	}
	return n, nil
}

// CreateGrades は複数の学年を作成します。
func (a *PgAdapter) CreateGrades(
	ctx context.Context, params []parameter.CreateGradeParam,
) (int64, error) {
	return createGrades(ctx, a.query, params)
}

// CreateGradesWithSd はSD付きで複数の学年を作成します。
func (a *PgAdapter) CreateGradesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateGradeParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createGrades(ctx, qtx, params)
}

// deleteGrade は学年を削除する内部関数です。
func deleteGrade(ctx context.Context, qtx *query.Queries, gradeID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteGrade(ctx, gradeID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete grade: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("grade")
	}
	return c, nil
}

// DeleteGrade は学年を削除します。
func (a *PgAdapter) DeleteGrade(ctx context.Context, gradeID uuid.UUID) (int64, error) {
	return deleteGrade(ctx, a.query, gradeID)
}

// DeleteGradeWithSd はSD付きで学年を削除します。
func (a *PgAdapter) DeleteGradeWithSd(
	ctx context.Context, sd store.Sd, gradeID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteGrade(ctx, qtx, gradeID)
}

// pluralDeleteGrades は複数の学年を削除する内部関数です。
func pluralDeleteGrades(ctx context.Context, qtx *query.Queries, gradeIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteGrades(ctx, gradeIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete grades: %w", err)
	}
	if c != int64(len(gradeIDs)) {
		return 0, errhandle.NewModelNotFoundError("grade")
	}
	return c, nil
}

// PluralDeleteGrades は複数の学年を削除します。
func (a *PgAdapter) PluralDeleteGrades(ctx context.Context, gradeIDs []uuid.UUID) (int64, error) {
	return pluralDeleteGrades(ctx, a.query, gradeIDs)
}

// PluralDeleteGradesWithSd はSD付きで複数の学年を削除します。
func (a *PgAdapter) PluralDeleteGradesWithSd(
	ctx context.Context, sd store.Sd, gradeIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteGrades(ctx, qtx, gradeIDs)
}

// findGradeByID は学年をIDで取得する内部関数です。
func findGradeByID(
	ctx context.Context, qtx *query.Queries, gradeID uuid.UUID,
) (entity.Grade, error) {
	e, err := qtx.FindGradeByID(ctx, gradeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Grade{}, errhandle.NewModelNotFoundError("grade")
		}
		return entity.Grade{}, fmt.Errorf("failed to find grade: %w", err)
	}
	entity := entity.Grade{
		GradeID:        e.GradeID,
		Key:            e.Key,
		OrganizationID: e.OrganizationID,
	}
	return entity, nil
}

// FindGradeByID は学年をIDで取得します。
func (a *PgAdapter) FindGradeByID(ctx context.Context, gradeID uuid.UUID) (entity.Grade, error) {
	return findGradeByID(ctx, a.query, gradeID)
}

// FindGradeByIDWithSd はSD付きで学年をIDで取得します。
func (a *PgAdapter) FindGradeByIDWithSd(
	ctx context.Context, sd store.Sd, gradeID uuid.UUID,
) (entity.Grade, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Grade{}, store.ErrNotFoundDescriptor
	}
	return findGradeByID(ctx, qtx, gradeID)
}

// findGradeWithOrganization は学年とオーガナイゼーションを取得する内部関数です。
func findGradeWithOrganization(
	ctx context.Context, qtx *query.Queries, gradeID uuid.UUID,
) (entity.GradeWithOrganization, error) {
	e, err := qtx.FindGradeByIDWithOrganization(ctx, gradeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.GradeWithOrganization{}, errhandle.NewModelNotFoundError("grade")
		}
		return entity.GradeWithOrganization{}, fmt.Errorf("failed to find grade with organization: %w", err)
	}
	return convGradeWithOrganization(e), nil
}

// FindGradeWithOrganization は学年とオーガナイゼーションを取得します。
func (a *PgAdapter) FindGradeWithOrganization(
	ctx context.Context, gradeID uuid.UUID,
) (entity.GradeWithOrganization, error) {
	return findGradeWithOrganization(ctx, a.query, gradeID)
}

// FindGradeWithOrganizationWithSd はSD付きで学年とオーガナイゼーションを取得します。
func (a *PgAdapter) FindGradeWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, gradeID uuid.UUID,
) (entity.GradeWithOrganization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.GradeWithOrganization{}, store.ErrNotFoundDescriptor
	}
	return findGradeWithOrganization(ctx, qtx, gradeID)
}

// getGrades は学年を取得する内部関数です。
func getGrades(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereGradeParam,
	order parameter.GradeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Grade], error) {
	eConvFunc := func(e query.Grade) (entity.Grade, error) {
		return entity.Grade{
			GradeID:        e.GradeID,
			Key:            e.Key,
			OrganizationID: e.OrganizationID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountGrades(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count grades: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Grade, error) {
		r, err := qtx.GetGrades(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Grade{}, nil
			}
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Grade, error) {
		p := query.GetGradesUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetGradesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.Grade, error) {
		p := query.GetGradesUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetGradesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Grade) (entity.Int, any) {
		switch subCursor {
		case parameter.GradeDefaultCursorKey:
			return entity.Int(e.MGradesPkey), nil
		}
		return entity.Int(e.MGradesPkey), nil
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
		return store.ListResult[entity.Grade]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	return res, nil
}

// GetGrades は学年を取得します。
func (a *PgAdapter) GetGrades(
	ctx context.Context, where parameter.WhereGradeParam,
	order parameter.GradeOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Grade], error) {
	return getGrades(ctx, a.query, where, order, np, cp, wc)
}

// GetGradesWithSd はSD付きで学年を取得します。
func (a *PgAdapter) GetGradesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGradeParam,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Grade], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Grade]{}, store.ErrNotFoundDescriptor
	}
	return getGrades(ctx, qtx, where, order, np, cp, wc)
}

// getPluralGrades は複数の学年を取得する内部関数です。
func getPluralGrades(
	ctx context.Context, qtx *query.Queries, gradeIDs []uuid.UUID,
	_ parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Grade], error) {
	var e []query.Grade
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralGrades(ctx, gradeIDs)
	} else {
		e, err = qtx.GetPluralGradesUseNumberedPaginate(ctx, query.GetPluralGradesUseNumberedPaginateParams{
			GradeIds: gradeIDs,
			Offset:   int32(np.Offset.Int64),
			Limit:    int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Grade]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	entities := make([]entity.Grade, len(e))
	for i, v := range e {
		entities[i] = entity.Grade{
			GradeID:        v.GradeID,
			Key:            v.Key,
			OrganizationID: v.OrganizationID,
		}
	}
	return store.ListResult[entity.Grade]{Data: entities}, nil
}

// GetPluralGrades は複数の学年を取得します。
func (a *PgAdapter) GetPluralGrades(
	ctx context.Context, gradeIDs []uuid.UUID,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Grade], error) {
	return getPluralGrades(ctx, a.query, gradeIDs, order, np)
}

// GetPluralGradesWithSd はSD付きで複数の学年を取得します。
func (a *PgAdapter) GetPluralGradesWithSd(
	ctx context.Context, sd store.Sd, gradeIDs []uuid.UUID,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Grade], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Grade]{}, store.ErrNotFoundDescriptor
	}
	return getPluralGrades(ctx, qtx, gradeIDs, order, np)
}

func getGradesWithOrganization(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereGradeParam,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	eConvFunc := func(e entity.GradeWithOrganizationForQuery) (entity.GradeWithOrganization, error) {
		return e.GradeWithOrganization, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountGrades(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count grades: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.GradeWithOrganizationForQuery, error) {
		r, err := qtx.GetGradesWithOrganization(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.GradeWithOrganizationForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		e := make([]entity.GradeWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GradeWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGradesPkey),
				GradeWithOrganization: convGradeWithOrganization(query.FindGradeByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.GradeWithOrganizationForQuery, error) {
		p := query.GetGradesWithOrganizationUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetGradesWithOrganizationUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		e := make([]entity.GradeWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GradeWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGradesPkey),
				GradeWithOrganization: convGradeWithOrganization(query.FindGradeByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.GradeWithOrganizationForQuery, error) {
		p := query.GetGradesWithOrganizationUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetGradesWithOrganizationUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get grades: %w", err)
		}
		e := make([]entity.GradeWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GradeWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGradesPkey),
				GradeWithOrganization: convGradeWithOrganization(query.FindGradeByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.GradeWithOrganizationForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.GradeDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		}
		return entity.Int(e.Pkey), nil
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
		return store.ListResult[entity.GradeWithOrganization]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	return res, nil
}

// GetGradesWithOrganization は学年とオーガナイゼーションを取得します。
func (a *PgAdapter) GetGradesWithOrganization(
	ctx context.Context, where parameter.WhereGradeParam,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	return getGradesWithOrganization(ctx, a.query, where, order, np, cp, wc)
}

// GetGradesWithOrganizationWithSd はSD付きで学年とオーガナイゼーションを取得します。
func (a *PgAdapter) GetGradesWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGradeParam,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.GradeWithOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getGradesWithOrganization(ctx, qtx, where, order, np, cp, wc)
}

// getPluralGradesWithOrganization は複数の学年を取得する内部関数です。
func getPluralGradesWithOrganization(
	ctx context.Context, qtx *query.Queries, gradeIDs []uuid.UUID,
	_ parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	var e []query.GetPluralGradesWithOrganizationRow
	var te []query.GetPluralGradesWithOrganizationUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralGradesWithOrganization(ctx, gradeIDs)
	} else {
		te, err = qtx.GetPluralGradesWithOrganizationUseNumberedPaginate(
			ctx, query.GetPluralGradesWithOrganizationUseNumberedPaginateParams{
				GradeIds: gradeIDs,
				Offset:   int32(np.Offset.Int64),
				Limit:    int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralGradesWithOrganizationRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralGradesWithOrganizationRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.GradeWithOrganization]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	entities := make([]entity.GradeWithOrganization, len(e))
	for i, v := range e {
		entities[i] = convGradeWithOrganization(query.FindGradeByIDWithOrganizationRow(v))
	}
	return store.ListResult[entity.GradeWithOrganization]{Data: entities}, nil
}

// GetPluralGradesWithOrganization は複数の学年を取得します。
func (a *PgAdapter) GetPluralGradesWithOrganization(
	ctx context.Context, gradeIDs []uuid.UUID,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	return getPluralGradesWithOrganization(ctx, a.query, gradeIDs, order, np)
}

// GetPluralGradesWithOrganizationWithSd はSD付きで複数の学年を取得します。
func (a *PgAdapter) GetPluralGradesWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, gradeIDs []uuid.UUID,
	order parameter.GradeOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GradeWithOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.GradeWithOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getPluralGradesWithOrganization(ctx, qtx, gradeIDs, order, np)
}

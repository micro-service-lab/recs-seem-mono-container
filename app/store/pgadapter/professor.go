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

func convProfessorWithMember(e query.FindProfessorByIDWithMemberRow) entity.ProfessorWithMember {
	return entity.ProfessorWithMember{
		ProfessorID: e.ProfessorID,
		Member: entity.MemberCard{
			MemberID:  e.MemberID,
			Name:      e.MemberName.String,
			FirstName: entity.String(e.MemberFirstName),
			LastName:  entity.String(e.MemberLastName),
			Email:     e.MemberEmail.String,
			ProfileImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
				Valid: e.MemberProfileImageID.Valid,
				Entity: entity.ImageWithAttachableItem{
					ImageID: e.MemberProfileImageID.Bytes,
					Height:  entity.Float(e.MemberProfileImageHeight),
					Width:   entity.Float(e.MemberProfileImageWidth),
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.MemberProfileImageAttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.MemberProfileImageOwnerID),
						FromOuter:        e.MemberProfileImageFromOuter.Bool,
						URL:              e.MemberProfileImageUrl.String,
						Alias:            e.MemberProfileImageAlias.String,
						Size:             entity.Float(e.MemberProfileImageSize),
						MimeTypeID:       e.MemberProfileImageMimeTypeID.Bytes,
					},
				},
			},
		},
	}
}

// countProfessors は教授数を取得する内部関数です。
func countProfessors(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereProfessorParam,
) (int64, error) {
	c, err := qtx.CountProfessors(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count professors: %w", err)
	}
	return c, nil
}

// CountProfessors は教授数を取得します。
func (a *PgAdapter) CountProfessors(ctx context.Context, where parameter.WhereProfessorParam) (int64, error) {
	return countProfessors(ctx, a.query, where)
}

// CountProfessorsWithSd はSD付きで教授数を取得します。
func (a *PgAdapter) CountProfessorsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereProfessorParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countProfessors(ctx, qtx, where)
}

// createProfessor は教授を作成する内部関数です。
func createProfessor(
	ctx context.Context, qtx *query.Queries, param parameter.CreateProfessorParam,
) (entity.Professor, error) {
	e, err := qtx.CreateProfessor(ctx, param.MemberID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Professor{}, errhandle.NewModelDuplicatedError("professor")
		}
		return entity.Professor{}, fmt.Errorf("failed to create professor: %w", err)
	}
	entity := entity.Professor{
		ProfessorID: e.ProfessorID,
		MemberID:    e.MemberID,
	}
	return entity, nil
}

// CreateProfessor は教授を作成します。
func (a *PgAdapter) CreateProfessor(
	ctx context.Context, param parameter.CreateProfessorParam,
) (entity.Professor, error) {
	return createProfessor(ctx, a.query, param)
}

// CreateProfessorWithSd はSD付きで教授を作成します。
func (a *PgAdapter) CreateProfessorWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateProfessorParam,
) (entity.Professor, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Professor{}, store.ErrNotFoundDescriptor
	}
	return createProfessor(ctx, qtx, param)
}

// createProfessors は複数の教授を作成する内部関数です。
func createProfessors(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateProfessorParam,
) (int64, error) {
	param := make([]uuid.UUID, len(params))
	for i, p := range params {
		param[i] = p.MemberID
	}
	n, err := qtx.CreateProfessors(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("professor")
		}
		return 0, fmt.Errorf("failed to create professors: %w", err)
	}
	return n, nil
}

// CreateProfessors は複数の教授を作成します。
func (a *PgAdapter) CreateProfessors(
	ctx context.Context, params []parameter.CreateProfessorParam,
) (int64, error) {
	return createProfessors(ctx, a.query, params)
}

// CreateProfessorsWithSd はSD付きで複数の教授を作成します。
func (a *PgAdapter) CreateProfessorsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateProfessorParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createProfessors(ctx, qtx, params)
}

// deleteProfessor は教授を削除する内部関数です。
func deleteProfessor(ctx context.Context, qtx *query.Queries, professorID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteProfessor(ctx, professorID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete professor: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("professor")
	}
	return c, nil
}

// DeleteProfessor は教授を削除します。
func (a *PgAdapter) DeleteProfessor(ctx context.Context, professorID uuid.UUID) (int64, error) {
	return deleteProfessor(ctx, a.query, professorID)
}

// DeleteProfessorWithSd はSD付きで教授を削除します。
func (a *PgAdapter) DeleteProfessorWithSd(
	ctx context.Context, sd store.Sd, professorID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteProfessor(ctx, qtx, professorID)
}

// pluralDeleteProfessors は複数の教授を削除する内部関数です。
func pluralDeleteProfessors(ctx context.Context, qtx *query.Queries, professorIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteProfessors(ctx, professorIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete professors: %w", err)
	}
	if c != int64(len(professorIDs)) {
		return 0, errhandle.NewModelNotFoundError("professor")
	}
	return c, nil
}

// PluralDeleteProfessors は複数の教授を削除します。
func (a *PgAdapter) PluralDeleteProfessors(ctx context.Context, professorIDs []uuid.UUID) (int64, error) {
	return pluralDeleteProfessors(ctx, a.query, professorIDs)
}

// PluralDeleteProfessorsWithSd はSD付きで複数の教授を削除します。
func (a *PgAdapter) PluralDeleteProfessorsWithSd(
	ctx context.Context, sd store.Sd, professorIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteProfessors(ctx, qtx, professorIDs)
}

// findProfessorByID は教授をIDで取得する内部関数です。
func findProfessorByID(
	ctx context.Context, qtx *query.Queries, professorID uuid.UUID,
) (entity.Professor, error) {
	e, err := qtx.FindProfessorByID(ctx, professorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Professor{}, errhandle.NewModelNotFoundError("professor")
		}
		return entity.Professor{}, fmt.Errorf("failed to find professor: %w", err)
	}
	entity := entity.Professor{
		ProfessorID: e.ProfessorID,
		MemberID:    e.MemberID,
	}
	return entity, nil
}

// FindProfessorByID は教授をIDで取得します。
func (a *PgAdapter) FindProfessorByID(ctx context.Context, professorID uuid.UUID) (entity.Professor, error) {
	return findProfessorByID(ctx, a.query, professorID)
}

// FindProfessorByIDWithSd はSD付きで教授をIDで取得します。
func (a *PgAdapter) FindProfessorByIDWithSd(
	ctx context.Context, sd store.Sd, professorID uuid.UUID,
) (entity.Professor, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Professor{}, store.ErrNotFoundDescriptor
	}
	return findProfessorByID(ctx, qtx, professorID)
}

// findProfessorWithMember は教授とオーガナイゼーションを取得する内部関数です。
func findProfessorWithMember(
	ctx context.Context, qtx *query.Queries, professorID uuid.UUID,
) (entity.ProfessorWithMember, error) {
	e, err := qtx.FindProfessorByIDWithMember(ctx, professorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ProfessorWithMember{}, errhandle.NewModelNotFoundError("professor")
		}
		return entity.ProfessorWithMember{}, fmt.Errorf("failed to find professor with member: %w", err)
	}
	return convProfessorWithMember(e), nil
}

// FindProfessorWithMember は教授とオーガナイゼーションを取得します。
func (a *PgAdapter) FindProfessorWithMember(
	ctx context.Context, professorID uuid.UUID,
) (entity.ProfessorWithMember, error) {
	return findProfessorWithMember(ctx, a.query, professorID)
}

// FindProfessorWithMemberWithSd はSD付きで教授とオーガナイゼーションを取得します。
func (a *PgAdapter) FindProfessorWithMemberWithSd(
	ctx context.Context, sd store.Sd, professorID uuid.UUID,
) (entity.ProfessorWithMember, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ProfessorWithMember{}, store.ErrNotFoundDescriptor
	}
	return findProfessorWithMember(ctx, qtx, professorID)
}

// getProfessors は教授を取得する内部関数です。
func getProfessors(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Professor], error) {
	eConvFunc := func(e query.Professor) (entity.Professor, error) {
		return entity.Professor{
			ProfessorID: e.ProfessorID,
			MemberID:    e.MemberID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountProfessors(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count professors: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Professor, error) {
		r, err := qtx.GetProfessors(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Professor{}, nil
			}
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Professor, error) {
		p := query.GetProfessorsUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetProfessorsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.Professor, error) {
		p := query.GetProfessorsUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetProfessorsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Professor) (entity.Int, any) {
		switch subCursor {
		case parameter.ProfessorDefaultCursorKey:
			return entity.Int(e.MProfessorsPkey), nil
		}
		return entity.Int(e.MProfessorsPkey), nil
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
		return store.ListResult[entity.Professor]{}, fmt.Errorf("failed to get professors: %w", err)
	}
	return res, nil
}

// GetProfessors は教授を取得します。
func (a *PgAdapter) GetProfessors(
	ctx context.Context, where parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Professor], error) {
	return getProfessors(ctx, a.query, where, order, np, cp, wc)
}

// GetProfessorsWithSd はSD付きで教授を取得します。
func (a *PgAdapter) GetProfessorsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Professor], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Professor]{}, store.ErrNotFoundDescriptor
	}
	return getProfessors(ctx, qtx, where, order, np, cp, wc)
}

// getPluralProfessors は複数の教授を取得する内部関数です。
func getPluralProfessors(
	ctx context.Context, qtx *query.Queries, professorIDs []uuid.UUID,
	_ parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Professor], error) {
	var e []query.Professor
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralProfessors(ctx, professorIDs)
	} else {
		e, err = qtx.GetPluralProfessorsUseNumberedPaginate(ctx, query.GetPluralProfessorsUseNumberedPaginateParams{
			ProfessorIds: professorIDs,
			Offset:       int32(np.Offset.Int64),
			Limit:        int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Professor]{}, fmt.Errorf("failed to get professors: %w", err)
	}
	entities := make([]entity.Professor, len(e))
	for i, v := range e {
		entities[i] = entity.Professor{
			ProfessorID: v.ProfessorID,
			MemberID:    v.MemberID,
		}
	}
	return store.ListResult[entity.Professor]{Data: entities}, nil
}

// GetPluralProfessors は複数の教授を取得します。
func (a *PgAdapter) GetPluralProfessors(
	ctx context.Context, professorIDs []uuid.UUID,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Professor], error) {
	return getPluralProfessors(ctx, a.query, professorIDs, order, np)
}

// GetPluralProfessorsWithSd はSD付きで複数の教授を取得します。
func (a *PgAdapter) GetPluralProfessorsWithSd(
	ctx context.Context, sd store.Sd, professorIDs []uuid.UUID,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Professor], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Professor]{}, store.ErrNotFoundDescriptor
	}
	return getPluralProfessors(ctx, qtx, professorIDs, order, np)
}

func getProfessorsWithMember(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	eConvFunc := func(e entity.ProfessorWithMemberForQuery) (entity.ProfessorWithMember, error) {
		return e.ProfessorWithMember, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountProfessors(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count professors: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ProfessorWithMemberForQuery, error) {
		r, err := qtx.GetProfessorsWithMember(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ProfessorWithMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		e := make([]entity.ProfessorWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ProfessorWithMemberForQuery{
				Pkey:                entity.Int(v.MProfessorsPkey),
				ProfessorWithMember: convProfessorWithMember(query.FindProfessorByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ProfessorWithMemberForQuery, error) {
		p := query.GetProfessorsWithMemberUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetProfessorsWithMemberUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		e := make([]entity.ProfessorWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ProfessorWithMemberForQuery{
				Pkey:                entity.Int(v.MProfessorsPkey),
				ProfessorWithMember: convProfessorWithMember(query.FindProfessorByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ProfessorWithMemberForQuery, error) {
		p := query.GetProfessorsWithMemberUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetProfessorsWithMemberUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get professors: %w", err)
		}
		e := make([]entity.ProfessorWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ProfessorWithMemberForQuery{
				Pkey:                entity.Int(v.MProfessorsPkey),
				ProfessorWithMember: convProfessorWithMember(query.FindProfessorByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ProfessorWithMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ProfessorDefaultCursorKey:
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
		return store.ListResult[entity.ProfessorWithMember]{}, fmt.Errorf("failed to get professors: %w", err)
	}
	return res, nil
}

// GetProfessorsWithMember は教授とオーガナイゼーションを取得します。
func (a *PgAdapter) GetProfessorsWithMember(
	ctx context.Context, where parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	return getProfessorsWithMember(ctx, a.query, where, order, np, cp, wc)
}

// GetProfessorsWithMemberWithSd はSD付きで教授とオーガナイゼーションを取得します。
func (a *PgAdapter) GetProfessorsWithMemberWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereProfessorParam,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ProfessorWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getProfessorsWithMember(ctx, qtx, where, order, np, cp, wc)
}

// getPluralProfessorsWithMember は複数の教授を取得する内部関数です。
func getPluralProfessorsWithMember(
	ctx context.Context, qtx *query.Queries, professorIDs []uuid.UUID,
	_ parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	var e []query.GetPluralProfessorsWithMemberRow
	var te []query.GetPluralProfessorsWithMemberUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralProfessorsWithMember(ctx, professorIDs)
	} else {
		te, err = qtx.GetPluralProfessorsWithMemberUseNumberedPaginate(
			ctx, query.GetPluralProfessorsWithMemberUseNumberedPaginateParams{
				ProfessorIds: professorIDs,
				Offset:       int32(np.Offset.Int64),
				Limit:        int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralProfessorsWithMemberRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralProfessorsWithMemberRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ProfessorWithMember]{}, fmt.Errorf("failed to get professors: %w", err)
	}
	entities := make([]entity.ProfessorWithMember, len(e))
	for i, v := range e {
		entities[i] = convProfessorWithMember(query.FindProfessorByIDWithMemberRow(v))
	}
	return store.ListResult[entity.ProfessorWithMember]{Data: entities}, nil
}

// GetPluralProfessorsWithMember は複数の教授を取得します。
func (a *PgAdapter) GetPluralProfessorsWithMember(
	ctx context.Context, professorIDs []uuid.UUID,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	return getPluralProfessorsWithMember(ctx, a.query, professorIDs, order, np)
}

// GetPluralProfessorsWithMemberWithSd はSD付きで複数の教授を取得します。
func (a *PgAdapter) GetPluralProfessorsWithMemberWithSd(
	ctx context.Context, sd store.Sd, professorIDs []uuid.UUID,
	order parameter.ProfessorOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ProfessorWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ProfessorWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralProfessorsWithMember(ctx, qtx, professorIDs, order, np)
}

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

func convStudentWithMember(e query.FindStudentByIDWithMemberRow) entity.StudentWithMember {
	return entity.StudentWithMember{
		StudentID: e.StudentID,
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

// countStudents は生徒数を取得する内部関数です。
func countStudents(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereStudentParam,
) (int64, error) {
	c, err := qtx.CountStudents(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count students: %w", err)
	}
	return c, nil
}

// CountStudents は生徒数を取得します。
func (a *PgAdapter) CountStudents(ctx context.Context, where parameter.WhereStudentParam) (int64, error) {
	return countStudents(ctx, a.query, where)
}

// CountStudentsWithSd はSD付きで生徒数を取得します。
func (a *PgAdapter) CountStudentsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereStudentParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countStudents(ctx, qtx, where)
}

// createStudent は生徒を作成する内部関数です。
func createStudent(
	ctx context.Context, qtx *query.Queries, param parameter.CreateStudentParam,
) (entity.Student, error) {
	e, err := qtx.CreateStudent(ctx, param.MemberID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Student{}, errhandle.NewModelDuplicatedError("student")
		}
		return entity.Student{}, fmt.Errorf("failed to create student: %w", err)
	}
	entity := entity.Student{
		StudentID: e.StudentID,
		MemberID:  e.MemberID,
	}
	return entity, nil
}

// CreateStudent は生徒を作成します。
func (a *PgAdapter) CreateStudent(
	ctx context.Context, param parameter.CreateStudentParam,
) (entity.Student, error) {
	return createStudent(ctx, a.query, param)
}

// CreateStudentWithSd はSD付きで生徒を作成します。
func (a *PgAdapter) CreateStudentWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateStudentParam,
) (entity.Student, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Student{}, store.ErrNotFoundDescriptor
	}
	return createStudent(ctx, qtx, param)
}

// createStudents は複数の生徒を作成する内部関数です。
func createStudents(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateStudentParam,
) (int64, error) {
	param := make([]uuid.UUID, len(params))
	for i, p := range params {
		param[i] = p.MemberID
	}
	n, err := qtx.CreateStudents(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("student")
		}
		return 0, fmt.Errorf("failed to create students: %w", err)
	}
	return n, nil
}

// CreateStudents は複数の生徒を作成します。
func (a *PgAdapter) CreateStudents(
	ctx context.Context, params []parameter.CreateStudentParam,
) (int64, error) {
	return createStudents(ctx, a.query, params)
}

// CreateStudentsWithSd はSD付きで複数の生徒を作成します。
func (a *PgAdapter) CreateStudentsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateStudentParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createStudents(ctx, qtx, params)
}

// deleteStudent は生徒を削除する内部関数です。
func deleteStudent(ctx context.Context, qtx *query.Queries, studentID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteStudent(ctx, studentID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete student: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("student")
	}
	return c, nil
}

// DeleteStudent は生徒を削除します。
func (a *PgAdapter) DeleteStudent(ctx context.Context, studentID uuid.UUID) (int64, error) {
	return deleteStudent(ctx, a.query, studentID)
}

// DeleteStudentWithSd はSD付きで生徒を削除します。
func (a *PgAdapter) DeleteStudentWithSd(
	ctx context.Context, sd store.Sd, studentID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteStudent(ctx, qtx, studentID)
}

// pluralDeleteStudents は複数の生徒を削除する内部関数です。
func pluralDeleteStudents(ctx context.Context, qtx *query.Queries, studentIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteStudents(ctx, studentIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete students: %w", err)
	}
	if c != int64(len(studentIDs)) {
		return 0, errhandle.NewModelNotFoundError("student")
	}
	return c, nil
}

// PluralDeleteStudents は複数の生徒を削除します。
func (a *PgAdapter) PluralDeleteStudents(ctx context.Context, studentIDs []uuid.UUID) (int64, error) {
	return pluralDeleteStudents(ctx, a.query, studentIDs)
}

// PluralDeleteStudentsWithSd はSD付きで複数の生徒を削除します。
func (a *PgAdapter) PluralDeleteStudentsWithSd(
	ctx context.Context, sd store.Sd, studentIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteStudents(ctx, qtx, studentIDs)
}

// findStudentByID は生徒をIDで取得する内部関数です。
func findStudentByID(
	ctx context.Context, qtx *query.Queries, studentID uuid.UUID,
) (entity.Student, error) {
	e, err := qtx.FindStudentByID(ctx, studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Student{}, errhandle.NewModelNotFoundError("student")
		}
		return entity.Student{}, fmt.Errorf("failed to find student: %w", err)
	}
	entity := entity.Student{
		StudentID: e.StudentID,
		MemberID:  e.MemberID,
	}
	return entity, nil
}

// FindStudentByID は生徒をIDで取得します。
func (a *PgAdapter) FindStudentByID(ctx context.Context, studentID uuid.UUID) (entity.Student, error) {
	return findStudentByID(ctx, a.query, studentID)
}

// FindStudentByIDWithSd はSD付きで生徒をIDで取得します。
func (a *PgAdapter) FindStudentByIDWithSd(
	ctx context.Context, sd store.Sd, studentID uuid.UUID,
) (entity.Student, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Student{}, store.ErrNotFoundDescriptor
	}
	return findStudentByID(ctx, qtx, studentID)
}

// findStudentWithMember は生徒とオーガナイゼーションを取得する内部関数です。
func findStudentWithMember(
	ctx context.Context, qtx *query.Queries, studentID uuid.UUID,
) (entity.StudentWithMember, error) {
	e, err := qtx.FindStudentByIDWithMember(ctx, studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.StudentWithMember{}, errhandle.NewModelNotFoundError("student")
		}
		return entity.StudentWithMember{}, fmt.Errorf("failed to find student with member: %w", err)
	}
	return convStudentWithMember(e), nil
}

// FindStudentWithMember は生徒とオーガナイゼーションを取得します。
func (a *PgAdapter) FindStudentWithMember(
	ctx context.Context, studentID uuid.UUID,
) (entity.StudentWithMember, error) {
	return findStudentWithMember(ctx, a.query, studentID)
}

// FindStudentWithMemberWithSd はSD付きで生徒とオーガナイゼーションを取得します。
func (a *PgAdapter) FindStudentWithMemberWithSd(
	ctx context.Context, sd store.Sd, studentID uuid.UUID,
) (entity.StudentWithMember, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.StudentWithMember{}, store.ErrNotFoundDescriptor
	}
	return findStudentWithMember(ctx, qtx, studentID)
}

// getStudents は生徒を取得する内部関数です。
func getStudents(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereStudentParam,
	order parameter.StudentOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Student], error) {
	eConvFunc := func(e query.Student) (entity.Student, error) {
		return entity.Student{
			StudentID: e.StudentID,
			MemberID:  e.MemberID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountStudents(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count students: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Student, error) {
		r, err := qtx.GetStudents(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Student{}, nil
			}
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Student, error) {
		p := query.GetStudentsUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetStudentsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.Student, error) {
		p := query.GetStudentsUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetStudentsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Student) (entity.Int, any) {
		switch subCursor {
		case parameter.StudentDefaultCursorKey:
			return entity.Int(e.MStudentsPkey), nil
		}
		return entity.Int(e.MStudentsPkey), nil
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
		return store.ListResult[entity.Student]{}, fmt.Errorf("failed to get students: %w", err)
	}
	return res, nil
}

// GetStudents は生徒を取得します。
func (a *PgAdapter) GetStudents(
	ctx context.Context, where parameter.WhereStudentParam,
	order parameter.StudentOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Student], error) {
	return getStudents(ctx, a.query, where, order, np, cp, wc)
}

// GetStudentsWithSd はSD付きで生徒を取得します。
func (a *PgAdapter) GetStudentsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereStudentParam,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Student], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Student]{}, store.ErrNotFoundDescriptor
	}
	return getStudents(ctx, qtx, where, order, np, cp, wc)
}

// getPluralStudents は複数の生徒を取得する内部関数です。
func getPluralStudents(
	ctx context.Context, qtx *query.Queries, studentIDs []uuid.UUID,
	_ parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Student], error) {
	var e []query.Student
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralStudents(ctx, studentIDs)
	} else {
		e, err = qtx.GetPluralStudentsUseNumberedPaginate(ctx, query.GetPluralStudentsUseNumberedPaginateParams{
			StudentIds: studentIDs,
			Offset:     int32(np.Offset.Int64),
			Limit:      int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Student]{}, fmt.Errorf("failed to get students: %w", err)
	}
	entities := make([]entity.Student, len(e))
	for i, v := range e {
		entities[i] = entity.Student{
			StudentID: v.StudentID,
			MemberID:  v.MemberID,
		}
	}
	return store.ListResult[entity.Student]{Data: entities}, nil
}

// GetPluralStudents は複数の生徒を取得します。
func (a *PgAdapter) GetPluralStudents(
	ctx context.Context, studentIDs []uuid.UUID,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Student], error) {
	return getPluralStudents(ctx, a.query, studentIDs, order, np)
}

// GetPluralStudentsWithSd はSD付きで複数の生徒を取得します。
func (a *PgAdapter) GetPluralStudentsWithSd(
	ctx context.Context, sd store.Sd, studentIDs []uuid.UUID,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Student], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Student]{}, store.ErrNotFoundDescriptor
	}
	return getPluralStudents(ctx, qtx, studentIDs, order, np)
}

func getStudentsWithMember(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereStudentParam,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.StudentWithMember], error) {
	eConvFunc := func(e entity.StudentWithMemberForQuery) (entity.StudentWithMember, error) {
		return e.StudentWithMember, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountStudents(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count students: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.StudentWithMemberForQuery, error) {
		r, err := qtx.GetStudentsWithMember(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.StudentWithMemberForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		e := make([]entity.StudentWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.StudentWithMemberForQuery{
				Pkey:              entity.Int(v.MStudentsPkey),
				StudentWithMember: convStudentWithMember(query.FindStudentByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.StudentWithMemberForQuery, error) {
		p := query.GetStudentsWithMemberUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetStudentsWithMemberUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		e := make([]entity.StudentWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.StudentWithMemberForQuery{
				Pkey:              entity.Int(v.MStudentsPkey),
				StudentWithMember: convStudentWithMember(query.FindStudentByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.StudentWithMemberForQuery, error) {
		p := query.GetStudentsWithMemberUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetStudentsWithMemberUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get students: %w", err)
		}
		e := make([]entity.StudentWithMemberForQuery, len(r))
		for i, v := range r {
			e[i] = entity.StudentWithMemberForQuery{
				Pkey:              entity.Int(v.MStudentsPkey),
				StudentWithMember: convStudentWithMember(query.FindStudentByIDWithMemberRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.StudentWithMemberForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.StudentDefaultCursorKey:
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
		return store.ListResult[entity.StudentWithMember]{}, fmt.Errorf("failed to get students: %w", err)
	}
	return res, nil
}

// GetStudentsWithMember は生徒とオーガナイゼーションを取得します。
func (a *PgAdapter) GetStudentsWithMember(
	ctx context.Context, where parameter.WhereStudentParam,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.StudentWithMember], error) {
	return getStudentsWithMember(ctx, a.query, where, order, np, cp, wc)
}

// GetStudentsWithMemberWithSd はSD付きで生徒とオーガナイゼーションを取得します。
func (a *PgAdapter) GetStudentsWithMemberWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereStudentParam,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.StudentWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.StudentWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getStudentsWithMember(ctx, qtx, where, order, np, cp, wc)
}

// getPluralStudentsWithMember は複数の生徒を取得する内部関数です。
func getPluralStudentsWithMember(
	ctx context.Context, qtx *query.Queries, studentIDs []uuid.UUID,
	_ parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.StudentWithMember], error) {
	var e []query.GetPluralStudentsWithMemberRow
	var te []query.GetPluralStudentsWithMemberUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralStudentsWithMember(ctx, studentIDs)
	} else {
		te, err = qtx.GetPluralStudentsWithMemberUseNumberedPaginate(
			ctx, query.GetPluralStudentsWithMemberUseNumberedPaginateParams{
				StudentIds: studentIDs,
				Offset:     int32(np.Offset.Int64),
				Limit:      int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralStudentsWithMemberRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralStudentsWithMemberRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.StudentWithMember]{}, fmt.Errorf("failed to get students: %w", err)
	}
	entities := make([]entity.StudentWithMember, len(e))
	for i, v := range e {
		entities[i] = convStudentWithMember(query.FindStudentByIDWithMemberRow(v))
	}
	return store.ListResult[entity.StudentWithMember]{Data: entities}, nil
}

// GetPluralStudentsWithMember は複数の生徒を取得します。
func (a *PgAdapter) GetPluralStudentsWithMember(
	ctx context.Context, studentIDs []uuid.UUID,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.StudentWithMember], error) {
	return getPluralStudentsWithMember(ctx, a.query, studentIDs, order, np)
}

// GetPluralStudentsWithMemberWithSd はSD付きで複数の生徒を取得します。
func (a *PgAdapter) GetPluralStudentsWithMemberWithSd(
	ctx context.Context, sd store.Sd, studentIDs []uuid.UUID,
	order parameter.StudentOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.StudentWithMember], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.StudentWithMember]{}, store.ErrNotFoundDescriptor
	}
	return getPluralStudentsWithMember(ctx, qtx, studentIDs, order, np)
}

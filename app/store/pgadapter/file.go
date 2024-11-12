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

func convFileWithAttachableItem(e query.FindFileByIDWithAttachableItemRow) entity.FileWithAttachableItem {
	return entity.FileWithAttachableItem{
		FileID: e.FileID,
		AttachableItem: entity.AttachableItem{
			AttachableItemID: e.AttachableItemID,
			OwnerID:          entity.UUID(e.OwnerID),
			FromOuter:        e.FromOuter.Bool,
			URL:              e.Url.String,
			Alias:            e.Alias.String,
			Size:             entity.Float(e.Size),
			MimeTypeID:       e.MimeTypeID.Bytes,
		},
	}
}

// countFiles はファイル数を取得する内部関数です。
func countFiles(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereFileParam,
) (int64, error) {
	c, err := qtx.CountFiles(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count files: %w", err)
	}
	return c, nil
}

// CountFiles はファイル数を取得します。
func (a *PgAdapter) CountFiles(ctx context.Context, where parameter.WhereFileParam) (int64, error) {
	return countFiles(ctx, a.query, where)
}

// CountFilesWithSd はSD付きでファイル数を取得します。
func (a *PgAdapter) CountFilesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereFileParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countFiles(ctx, qtx, where)
}

// createFile はファイルを作成する内部関数です。
func createFile(
	ctx context.Context, qtx *query.Queries, param parameter.CreateFileParam,
) (entity.File, error) {
	e, err := qtx.CreateFile(ctx, param.AttachableItemID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.File{}, errhandle.NewModelDuplicatedError("file")
		}
		return entity.File{}, fmt.Errorf("failed to create file: %w", err)
	}
	entity := entity.File{
		FileID:           e.FileID,
		AttachableItemID: e.AttachableItemID,
	}
	return entity, nil
}

// CreateFile はファイルを作成します。
func (a *PgAdapter) CreateFile(
	ctx context.Context, param parameter.CreateFileParam,
) (entity.File, error) {
	return createFile(ctx, a.query, param)
}

// CreateFileWithSd はSD付きでファイルを作成します。
func (a *PgAdapter) CreateFileWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateFileParam,
) (entity.File, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.File{}, store.ErrNotFoundDescriptor
	}
	return createFile(ctx, qtx, param)
}

// createFiles は複数のファイルを作成する内部関数です。
func createFiles(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateFileParam,
) (int64, error) {
	param := make([]uuid.UUID, len(params))
	for i, p := range params {
		param[i] = p.AttachableItemID
	}
	n, err := qtx.CreateFiles(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("file")
		}
		return 0, fmt.Errorf("failed to create files: %w", err)
	}
	return n, nil
}

// CreateFiles は複数のファイルを作成します。
func (a *PgAdapter) CreateFiles(
	ctx context.Context, params []parameter.CreateFileParam,
) (int64, error) {
	return createFiles(ctx, a.query, params)
}

// CreateFilesWithSd はSD付きで複数のファイルを作成します。
func (a *PgAdapter) CreateFilesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateFileParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createFiles(ctx, qtx, params)
}

// deleteFile はファイルを削除する内部関数です。
func deleteFile(ctx context.Context, qtx *query.Queries, fileID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteFile(ctx, fileID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete file: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("file")
	}
	return c, nil
}

// DeleteFile はファイルを削除します。
func (a *PgAdapter) DeleteFile(ctx context.Context, fileID uuid.UUID) (int64, error) {
	return deleteFile(ctx, a.query, fileID)
}

// DeleteFileWithSd はSD付きでファイルを削除します。
func (a *PgAdapter) DeleteFileWithSd(
	ctx context.Context, sd store.Sd, fileID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteFile(ctx, qtx, fileID)
}

// pluralDeleteFiles は複数のファイルを削除する内部関数です。
func pluralDeleteFiles(ctx context.Context, qtx *query.Queries, fileIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteFiles(ctx, fileIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete files: %w", err)
	}
	if c != int64(len(fileIDs)) {
		return 0, errhandle.NewModelNotFoundError("file")
	}
	return c, nil
}

// PluralDeleteFiles は複数のファイルを削除します。
func (a *PgAdapter) PluralDeleteFiles(ctx context.Context, fileIDs []uuid.UUID) (int64, error) {
	return pluralDeleteFiles(ctx, a.query, fileIDs)
}

// PluralDeleteFilesWithSd はSD付きで複数のファイルを削除します。
func (a *PgAdapter) PluralDeleteFilesWithSd(
	ctx context.Context, sd store.Sd, fileIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteFiles(ctx, qtx, fileIDs)
}

// findFileByID はファイルをIDで取得する内部関数です。
func findFileByID(
	ctx context.Context, qtx *query.Queries, fileID uuid.UUID,
) (entity.File, error) {
	e, err := qtx.FindFileByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.File{}, errhandle.NewModelNotFoundError("file")
		}
		return entity.File{}, fmt.Errorf("failed to find file: %w", err)
	}
	entity := entity.File{
		FileID:           e.FileID,
		AttachableItemID: e.AttachableItemID,
	}
	return entity, nil
}

// FindFileByID はファイルをIDで取得します。
func (a *PgAdapter) FindFileByID(ctx context.Context, fileID uuid.UUID) (entity.File, error) {
	return findFileByID(ctx, a.query, fileID)
}

// FindFileByIDWithSd はSD付きでファイルをIDで取得します。
func (a *PgAdapter) FindFileByIDWithSd(
	ctx context.Context, sd store.Sd, fileID uuid.UUID,
) (entity.File, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.File{}, store.ErrNotFoundDescriptor
	}
	return findFileByID(ctx, qtx, fileID)
}

// findFileWithAttachableItem はファイルとオーガナイゼーションを取得する内部関数です。
func findFileWithAttachableItem(
	ctx context.Context, qtx *query.Queries, fileID uuid.UUID,
) (entity.FileWithAttachableItem, error) {
	e, err := qtx.FindFileByIDWithAttachableItem(ctx, fileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.FileWithAttachableItem{}, errhandle.NewModelNotFoundError("file")
		}
		return entity.FileWithAttachableItem{}, fmt.Errorf("failed to find file with member: %w", err)
	}
	return convFileWithAttachableItem(e), nil
}

// FindFileWithAttachableItem はファイルとオーガナイゼーションを取得します。
func (a *PgAdapter) FindFileWithAttachableItem(
	ctx context.Context, fileID uuid.UUID,
) (entity.FileWithAttachableItem, error) {
	return findFileWithAttachableItem(ctx, a.query, fileID)
}

// FindFileWithAttachableItemWithSd はSD付きでファイルとオーガナイゼーションを取得します。
func (a *PgAdapter) FindFileWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, fileID uuid.UUID,
) (entity.FileWithAttachableItem, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.FileWithAttachableItem{}, store.ErrNotFoundDescriptor
	}
	return findFileWithAttachableItem(ctx, qtx, fileID)
}

// getFiles はファイルを取得する内部関数です。
func getFiles(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereFileParam,
	order parameter.FileOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.File], error) {
	eConvFunc := func(e query.File) (entity.File, error) {
		return entity.File{
			FileID:           e.FileID,
			AttachableItemID: e.AttachableItemID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountFiles(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count files: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.File, error) {
		r, err := qtx.GetFiles(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.File{}, nil
			}
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.File, error) {
		p := query.GetFilesUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetFilesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.File, error) {
		p := query.GetFilesUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetFilesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.File) (entity.Int, any) {
		switch subCursor {
		case parameter.FileDefaultCursorKey:
			return entity.Int(e.TFilesPkey), nil
		}
		return entity.Int(e.TFilesPkey), nil
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
		return store.ListResult[entity.File]{}, fmt.Errorf("failed to get files: %w", err)
	}
	return res, nil
}

// GetFiles はファイルを取得します。
func (a *PgAdapter) GetFiles(
	ctx context.Context, where parameter.WhereFileParam,
	order parameter.FileOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.File], error) {
	return getFiles(ctx, a.query, where, order, np, cp, wc)
}

// GetFilesWithSd はSD付きでファイルを取得します。
func (a *PgAdapter) GetFilesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereFileParam,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.File], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.File]{}, store.ErrNotFoundDescriptor
	}
	return getFiles(ctx, qtx, where, order, np, cp, wc)
}

// getPluralFiles は複数のファイルを取得する内部関数です。
func getPluralFiles(
	ctx context.Context, qtx *query.Queries, fileIDs []uuid.UUID,
	_ parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.File], error) {
	var e []query.File
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralFiles(ctx, fileIDs)
	} else {
		e, err = qtx.GetPluralFilesUseNumberedPaginate(ctx, query.GetPluralFilesUseNumberedPaginateParams{
			FileIds: fileIDs,
			Offset:  int32(np.Offset.Int64),
			Limit:   int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.File]{}, fmt.Errorf("failed to get files: %w", err)
	}
	entities := make([]entity.File, len(e))
	for i, v := range e {
		entities[i] = entity.File{
			FileID:           v.FileID,
			AttachableItemID: v.AttachableItemID,
		}
	}
	return store.ListResult[entity.File]{Data: entities}, nil
}

// GetPluralFiles は複数のファイルを取得します。
func (a *PgAdapter) GetPluralFiles(
	ctx context.Context, fileIDs []uuid.UUID,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.File], error) {
	return getPluralFiles(ctx, a.query, fileIDs, order, np)
}

// GetPluralFilesWithSd はSD付きで複数のファイルを取得します。
func (a *PgAdapter) GetPluralFilesWithSd(
	ctx context.Context, sd store.Sd, fileIDs []uuid.UUID,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.File], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.File]{}, store.ErrNotFoundDescriptor
	}
	return getPluralFiles(ctx, qtx, fileIDs, order, np)
}

func getFilesWithAttachableItem(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereFileParam,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	eConvFunc := func(e entity.FileWithAttachableItemForQuery) (entity.FileWithAttachableItem, error) {
		return e.FileWithAttachableItem, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountFiles(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count files: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.FileWithAttachableItemForQuery, error) {
		r, err := qtx.GetFilesWithAttachableItem(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.FileWithAttachableItemForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		e := make([]entity.FileWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.FileWithAttachableItemForQuery{
				Pkey:                   entity.Int(v.TFilesPkey),
				FileWithAttachableItem: convFileWithAttachableItem(query.FindFileByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.FileWithAttachableItemForQuery, error) {
		p := query.GetFilesWithAttachableItemUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetFilesWithAttachableItemUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		e := make([]entity.FileWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.FileWithAttachableItemForQuery{
				Pkey:                   entity.Int(v.TFilesPkey),
				FileWithAttachableItem: convFileWithAttachableItem(query.FindFileByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.FileWithAttachableItemForQuery, error) {
		p := query.GetFilesWithAttachableItemUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetFilesWithAttachableItemUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get files: %w", err)
		}
		e := make([]entity.FileWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.FileWithAttachableItemForQuery{
				Pkey:                   entity.Int(v.TFilesPkey),
				FileWithAttachableItem: convFileWithAttachableItem(query.FindFileByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.FileWithAttachableItemForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.FileDefaultCursorKey:
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
		return store.ListResult[entity.FileWithAttachableItem]{}, fmt.Errorf("failed to get files: %w", err)
	}
	return res, nil
}

// GetFilesWithAttachableItem はファイルとオーガナイゼーションを取得します。
func (a *PgAdapter) GetFilesWithAttachableItem(
	ctx context.Context, where parameter.WhereFileParam,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	return getFilesWithAttachableItem(ctx, a.query, where, order, np, cp, wc)
}

// GetFilesWithAttachableItemWithSd はSD付きでファイルとオーガナイゼーションを取得します。
func (a *PgAdapter) GetFilesWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereFileParam,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.FileWithAttachableItem]{}, store.ErrNotFoundDescriptor
	}
	return getFilesWithAttachableItem(ctx, qtx, where, order, np, cp, wc)
}

// getPluralFilesWithAttachableItem は複数のファイルを取得する内部関数です。
func getPluralFilesWithAttachableItem(
	ctx context.Context, qtx *query.Queries, fileIDs []uuid.UUID,
	_ parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	var e []query.GetPluralFilesWithAttachableItemRow
	var te []query.GetPluralFilesWithAttachableItemUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralFilesWithAttachableItem(ctx, fileIDs)
	} else {
		te, err = qtx.GetPluralFilesWithAttachableItemUseNumberedPaginate(
			ctx, query.GetPluralFilesWithAttachableItemUseNumberedPaginateParams{
				FileIds: fileIDs,
				Offset:  int32(np.Offset.Int64),
				Limit:   int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralFilesWithAttachableItemRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralFilesWithAttachableItemRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.FileWithAttachableItem]{}, fmt.Errorf("failed to get files: %w", err)
	}
	entities := make([]entity.FileWithAttachableItem, len(e))
	for i, v := range e {
		entities[i] = convFileWithAttachableItem(query.FindFileByIDWithAttachableItemRow(v))
	}
	return store.ListResult[entity.FileWithAttachableItem]{Data: entities}, nil
}

// GetPluralFilesWithAttachableItem は複数のファイルを取得します。
func (a *PgAdapter) GetPluralFilesWithAttachableItem(
	ctx context.Context, fileIDs []uuid.UUID,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	return getPluralFilesWithAttachableItem(ctx, a.query, fileIDs, order, np)
}

// GetPluralFilesWithAttachableItemWithSd はSD付きで複数のファイルを取得します。
func (a *PgAdapter) GetPluralFilesWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, fileIDs []uuid.UUID,
	order parameter.FileOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.FileWithAttachableItem], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.FileWithAttachableItem]{}, store.ErrNotFoundDescriptor
	}
	return getPluralFilesWithAttachableItem(ctx, qtx, fileIDs, order, np)
}

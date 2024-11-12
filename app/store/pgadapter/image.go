package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convImageWithAttachableItem(e query.FindImageByIDWithAttachableItemRow) entity.ImageWithAttachableItem {
	return entity.ImageWithAttachableItem{
		ImageID: e.ImageID,
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

// countImages は画像数を取得する内部関数です。
func countImages(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereImageParam,
) (int64, error) {
	c, err := qtx.CountImages(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count images: %w", err)
	}
	return c, nil
}

// CountImages は画像数を取得します。
func (a *PgAdapter) CountImages(ctx context.Context, where parameter.WhereImageParam) (int64, error) {
	return countImages(ctx, a.query, where)
}

// CountImagesWithSd はSD付きで画像数を取得します。
func (a *PgAdapter) CountImagesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereImageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countImages(ctx, qtx, where)
}

// createImage は画像を作成する内部関数です。
func createImage(
	ctx context.Context, qtx *query.Queries, param parameter.CreateImageParam,
) (entity.Image, error) {
	e, err := qtx.CreateImage(ctx, query.CreateImageParams{
		AttachableItemID: param.AttachableItemID,
		Height:           pgtype.Float8(param.Height),
		Width:            pgtype.Float8(param.Width),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Image{}, errhandle.NewModelDuplicatedError("image")
		}
		return entity.Image{}, fmt.Errorf("failed to create image: %w", err)
	}
	entity := entity.Image{
		ImageID:          e.ImageID,
		AttachableItemID: e.AttachableItemID,
	}
	return entity, nil
}

// CreateImage は画像を作成します。
func (a *PgAdapter) CreateImage(
	ctx context.Context, param parameter.CreateImageParam,
) (entity.Image, error) {
	return createImage(ctx, a.query, param)
}

// CreateImageWithSd はSD付きで画像を作成します。
func (a *PgAdapter) CreateImageWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateImageParam,
) (entity.Image, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Image{}, store.ErrNotFoundDescriptor
	}
	return createImage(ctx, qtx, param)
}

// createImages は複数の画像を作成する内部関数です。
func createImages(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateImageParam,
) (int64, error) {
	param := make([]query.CreateImagesParams, len(params))
	for i, p := range params {
		param[i] = query.CreateImagesParams{
			AttachableItemID: p.AttachableItemID,
			Height:           pgtype.Float8(p.Height),
			Width:            pgtype.Float8(p.Width),
		}
	}
	n, err := qtx.CreateImages(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("image")
		}
		return 0, fmt.Errorf("failed to create images: %w", err)
	}
	return n, nil
}

// CreateImages は複数の画像を作成します。
func (a *PgAdapter) CreateImages(
	ctx context.Context, params []parameter.CreateImageParam,
) (int64, error) {
	return createImages(ctx, a.query, params)
}

// CreateImagesWithSd はSD付きで複数の画像を作成します。
func (a *PgAdapter) CreateImagesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateImageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createImages(ctx, qtx, params)
}

// deleteImage は画像を削除する内部関数です。
func deleteImage(ctx context.Context, qtx *query.Queries, imageID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteImage(ctx, imageID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete image: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("image")
	}
	return c, nil
}

// DeleteImage は画像を削除します。
func (a *PgAdapter) DeleteImage(ctx context.Context, imageID uuid.UUID) (int64, error) {
	return deleteImage(ctx, a.query, imageID)
}

// DeleteImageWithSd はSD付きで画像を削除します。
func (a *PgAdapter) DeleteImageWithSd(
	ctx context.Context, sd store.Sd, imageID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteImage(ctx, qtx, imageID)
}

// pluralDeleteImages は複数の画像を削除する内部関数です。
func pluralDeleteImages(ctx context.Context, qtx *query.Queries, imageIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteImages(ctx, imageIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete images: %w", err)
	}
	if c != int64(len(imageIDs)) {
		return 0, errhandle.NewModelNotFoundError("image")
	}
	return c, nil
}

// PluralDeleteImages は複数の画像を削除します。
func (a *PgAdapter) PluralDeleteImages(ctx context.Context, imageIDs []uuid.UUID) (int64, error) {
	return pluralDeleteImages(ctx, a.query, imageIDs)
}

// PluralDeleteImagesWithSd はSD付きで複数の画像を削除します。
func (a *PgAdapter) PluralDeleteImagesWithSd(
	ctx context.Context, sd store.Sd, imageIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteImages(ctx, qtx, imageIDs)
}

// findImageByID は画像をIDで取得する内部関数です。
func findImageByID(
	ctx context.Context, qtx *query.Queries, imageID uuid.UUID,
) (entity.Image, error) {
	e, err := qtx.FindImageByID(ctx, imageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Image{}, errhandle.NewModelNotFoundError("image")
		}
		return entity.Image{}, fmt.Errorf("failed to find image: %w", err)
	}
	entity := entity.Image{
		ImageID:          e.ImageID,
		AttachableItemID: e.AttachableItemID,
	}
	return entity, nil
}

// FindImageByID は画像をIDで取得します。
func (a *PgAdapter) FindImageByID(ctx context.Context, imageID uuid.UUID) (entity.Image, error) {
	return findImageByID(ctx, a.query, imageID)
}

// FindImageByIDWithSd はSD付きで画像をIDで取得します。
func (a *PgAdapter) FindImageByIDWithSd(
	ctx context.Context, sd store.Sd, imageID uuid.UUID,
) (entity.Image, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Image{}, store.ErrNotFoundDescriptor
	}
	return findImageByID(ctx, qtx, imageID)
}

// findImageWithAttachableItem は画像とオーガナイゼーションを取得する内部関数です。
func findImageWithAttachableItem(
	ctx context.Context, qtx *query.Queries, imageID uuid.UUID,
) (entity.ImageWithAttachableItem, error) {
	e, err := qtx.FindImageByIDWithAttachableItem(ctx, imageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("image")
		}
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to find image with member: %w", err)
	}
	return convImageWithAttachableItem(e), nil
}

// FindImageWithAttachableItem は画像とオーガナイゼーションを取得します。
func (a *PgAdapter) FindImageWithAttachableItem(
	ctx context.Context, imageID uuid.UUID,
) (entity.ImageWithAttachableItem, error) {
	return findImageWithAttachableItem(ctx, a.query, imageID)
}

// FindImageWithAttachableItemWithSd はSD付きで画像とオーガナイゼーションを取得します。
func (a *PgAdapter) FindImageWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, imageID uuid.UUID,
) (entity.ImageWithAttachableItem, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ImageWithAttachableItem{}, store.ErrNotFoundDescriptor
	}
	return findImageWithAttachableItem(ctx, qtx, imageID)
}

// getImages は画像を取得する内部関数です。
func getImages(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereImageParam,
	order parameter.ImageOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Image], error) {
	eConvFunc := func(e query.Image) (entity.Image, error) {
		return entity.Image{
			ImageID:          e.ImageID,
			AttachableItemID: e.AttachableItemID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountImages(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count images: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Image, error) {
		r, err := qtx.GetImages(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Image{}, nil
			}
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Image, error) {
		p := query.GetImagesUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetImagesUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.Image, error) {
		p := query.GetImagesUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetImagesUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Image) (entity.Int, any) {
		switch subCursor {
		case parameter.ImageDefaultCursorKey:
			return entity.Int(e.TImagesPkey), nil
		}
		return entity.Int(e.TImagesPkey), nil
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
		return store.ListResult[entity.Image]{}, fmt.Errorf("failed to get images: %w", err)
	}
	return res, nil
}

// GetImages は画像を取得します。
func (a *PgAdapter) GetImages(
	ctx context.Context, where parameter.WhereImageParam,
	order parameter.ImageOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Image], error) {
	return getImages(ctx, a.query, where, order, np, cp, wc)
}

// GetImagesWithSd はSD付きで画像を取得します。
func (a *PgAdapter) GetImagesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereImageParam,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Image], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Image]{}, store.ErrNotFoundDescriptor
	}
	return getImages(ctx, qtx, where, order, np, cp, wc)
}

// getPluralImages は複数の画像を取得する内部関数です。
func getPluralImages(
	ctx context.Context, qtx *query.Queries, imageIDs []uuid.UUID,
	_ parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Image], error) {
	var e []query.Image
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralImages(ctx, imageIDs)
	} else {
		e, err = qtx.GetPluralImagesUseNumberedPaginate(ctx, query.GetPluralImagesUseNumberedPaginateParams{
			ImageIds: imageIDs,
			Offset:   int32(np.Offset.Int64),
			Limit:    int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Image]{}, fmt.Errorf("failed to get images: %w", err)
	}
	entities := make([]entity.Image, len(e))
	for i, v := range e {
		entities[i] = entity.Image{
			ImageID:          v.ImageID,
			AttachableItemID: v.AttachableItemID,
		}
	}
	return store.ListResult[entity.Image]{Data: entities}, nil
}

// GetPluralImages は複数の画像を取得します。
func (a *PgAdapter) GetPluralImages(
	ctx context.Context, imageIDs []uuid.UUID,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Image], error) {
	return getPluralImages(ctx, a.query, imageIDs, order, np)
}

// GetPluralImagesWithSd はSD付きで複数の画像を取得します。
func (a *PgAdapter) GetPluralImagesWithSd(
	ctx context.Context, sd store.Sd, imageIDs []uuid.UUID,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Image], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Image]{}, store.ErrNotFoundDescriptor
	}
	return getPluralImages(ctx, qtx, imageIDs, order, np)
}

func getImagesWithAttachableItem(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereImageParam,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	eConvFunc := func(e entity.ImageWithAttachableItemForQuery) (entity.ImageWithAttachableItem, error) {
		return e.ImageWithAttachableItem, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountImages(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count images: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.ImageWithAttachableItemForQuery, error) {
		r, err := qtx.GetImagesWithAttachableItem(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ImageWithAttachableItemForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		e := make([]entity.ImageWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ImageWithAttachableItemForQuery{
				Pkey:                    entity.Int(v.TImagesPkey),
				ImageWithAttachableItem: convImageWithAttachableItem(query.FindImageByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.ImageWithAttachableItemForQuery, error) {
		p := query.GetImagesWithAttachableItemUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetImagesWithAttachableItemUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		e := make([]entity.ImageWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ImageWithAttachableItemForQuery{
				Pkey:                    entity.Int(v.TImagesPkey),
				ImageWithAttachableItem: convImageWithAttachableItem(query.FindImageByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.ImageWithAttachableItemForQuery, error) {
		p := query.GetImagesWithAttachableItemUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetImagesWithAttachableItemUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}
		e := make([]entity.ImageWithAttachableItemForQuery, len(r))
		for i, v := range r {
			e[i] = entity.ImageWithAttachableItemForQuery{
				Pkey:                    entity.Int(v.TImagesPkey),
				ImageWithAttachableItem: convImageWithAttachableItem(query.FindImageByIDWithAttachableItemRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ImageWithAttachableItemForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ImageDefaultCursorKey:
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
		return store.ListResult[entity.ImageWithAttachableItem]{}, fmt.Errorf("failed to get images: %w", err)
	}
	return res, nil
}

// GetImagesWithAttachableItem は画像とオーガナイゼーションを取得します。
func (a *PgAdapter) GetImagesWithAttachableItem(
	ctx context.Context, where parameter.WhereImageParam,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	return getImagesWithAttachableItem(ctx, a.query, where, order, np, cp, wc)
}

// GetImagesWithAttachableItemWithSd はSD付きで画像とオーガナイゼーションを取得します。
func (a *PgAdapter) GetImagesWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereImageParam,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ImageWithAttachableItem]{}, store.ErrNotFoundDescriptor
	}
	return getImagesWithAttachableItem(ctx, qtx, where, order, np, cp, wc)
}

// getPluralImagesWithAttachableItem は複数の画像を取得する内部関数です。
func getPluralImagesWithAttachableItem(
	ctx context.Context, qtx *query.Queries, imageIDs []uuid.UUID,
	_ parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	var e []query.GetPluralImagesWithAttachableItemRow
	var te []query.GetPluralImagesWithAttachableItemUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralImagesWithAttachableItem(ctx, imageIDs)
	} else {
		te, err = qtx.GetPluralImagesWithAttachableItemUseNumberedPaginate(
			ctx, query.GetPluralImagesWithAttachableItemUseNumberedPaginateParams{
				ImageIds: imageIDs,
				Offset:   int32(np.Offset.Int64),
				Limit:    int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralImagesWithAttachableItemRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralImagesWithAttachableItemRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ImageWithAttachableItem]{}, fmt.Errorf("failed to get images: %w", err)
	}
	entities := make([]entity.ImageWithAttachableItem, len(e))
	for i, v := range e {
		entities[i] = convImageWithAttachableItem(query.FindImageByIDWithAttachableItemRow(v))
	}
	return store.ListResult[entity.ImageWithAttachableItem]{Data: entities}, nil
}

// GetPluralImagesWithAttachableItem は複数の画像を取得します。
func (a *PgAdapter) GetPluralImagesWithAttachableItem(
	ctx context.Context, imageIDs []uuid.UUID,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	return getPluralImagesWithAttachableItem(ctx, a.query, imageIDs, order, np)
}

// GetPluralImagesWithAttachableItemWithSd はSD付きで複数の画像を取得します。
func (a *PgAdapter) GetPluralImagesWithAttachableItemWithSd(
	ctx context.Context, sd store.Sd, imageIDs []uuid.UUID,
	order parameter.ImageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ImageWithAttachableItem], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ImageWithAttachableItem]{}, store.ErrNotFoundDescriptor
	}
	return getPluralImagesWithAttachableItem(ctx, qtx, imageIDs, order, np)
}

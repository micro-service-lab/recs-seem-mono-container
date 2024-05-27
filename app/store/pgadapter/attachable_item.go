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

func attachableItemConv(
	attachableItem query.FindAttachableItemByIDRow,
) entity.AttachableItemWithContent {
	var image entity.NullableEntity[entity.Image]
	if attachableItem.ImageID.Valid {
		image = entity.NullableEntity[entity.Image]{
			Entity: entity.Image{
				ImageID:          attachableItem.ImageID.Bytes,
				Height:           entity.Float(attachableItem.ImageHeight),
				Width:            entity.Float(attachableItem.ImageWidth),
				AttachableItemID: attachableItem.AttachableItemID,
			},
			Valid: true,
		}
	}
	var file entity.NullableEntity[entity.File]
	if attachableItem.FileID.Valid {
		file = entity.NullableEntity[entity.File]{
			Entity: entity.File{
				FileID:           attachableItem.FileID.Bytes,
				AttachableItemID: attachableItem.AttachableItemID,
			},
			Valid: true,
		}
	}
	entity := entity.AttachableItemWithContent{
		AttachableItemID: attachableItem.AttachableItemID,
		OwnerID:          entity.UUID(attachableItem.OwnerID),
		FromOuter:        attachableItem.FromOuter,
		URL:              attachableItem.Url,
		Alias:            attachableItem.Alias,
		Size:             entity.Float(attachableItem.Size),
		MimeTypeID:       attachableItem.MimeTypeID,
		Image:            image,
		File:             file,
	}
	return entity
}

func attachableItemConvWithMimeType(
	attachableItem query.FindAttachableItemByIDWithMimeTypeRow,
) entity.AttachableItemWithMimeType {
	var image entity.NullableEntity[entity.Image]
	if attachableItem.ImageID.Valid {
		image = entity.NullableEntity[entity.Image]{
			Entity: entity.Image{
				ImageID:          attachableItem.ImageID.Bytes,
				Height:           entity.Float(attachableItem.ImageHeight),
				Width:            entity.Float(attachableItem.ImageWidth),
				AttachableItemID: attachableItem.AttachableItemID,
			},
			Valid: true,
		}
	}
	var file entity.NullableEntity[entity.File]
	if attachableItem.FileID.Valid {
		file = entity.NullableEntity[entity.File]{
			Entity: entity.File{
				FileID:           attachableItem.FileID.Bytes,
				AttachableItemID: attachableItem.AttachableItemID,
			},
			Valid: true,
		}
	}
	mimeType := entity.NullableEntity[entity.MimeType]{
		Entity: entity.MimeType{
			MimeTypeID: attachableItem.MimeTypeID,
			Name:       attachableItem.MimeTypeName.String,
			Key:        attachableItem.MimeTypeKey.String,
			Kind:       attachableItem.MimeTypeKind.String,
		},
		Valid: true,
	}
	entity := entity.AttachableItemWithMimeType{
		AttachableItemID: attachableItem.AttachableItemID,
		OwnerID:          entity.UUID(attachableItem.OwnerID),
		FromOuter:        attachableItem.FromOuter,
		URL:              attachableItem.Url,
		Alias:            attachableItem.Alias,
		Size:             entity.Float(attachableItem.Size),
		MimeType:         mimeType,
		Image:            image,
		File:             file,
	}
	return entity
}

func countAttachableItems(ctx context.Context,
	qtx *query.Queries, where parameter.WhereAttachableItemParam,
) (int64, error) {
	p := query.CountAttachableItemsParams{
		WhereInMimeTypeIds: where.WhereInMimeType,
		InMimeTypeIds:      where.InMimeTypes,
		WhereInOwnerIds:    where.WhereInOwner,
		InOwnerIds:         where.InOwners,
	}
	c, err := qtx.CountAttachableItems(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count attachable items: %w", err)
	}
	return c, nil
}

// CountAttachableItems 添付可能アイテム数を取得する。
func (a *PgAdapter) CountAttachableItems(
	ctx context.Context, where parameter.WhereAttachableItemParam,
) (int64, error) {
	return countAttachableItems(ctx, a.query, where)
}

// CountAttachableItemsWithSd SD付きで添付可能アイテム数を取得する。
func (a *PgAdapter) CountAttachableItemsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttachableItemParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countAttachableItems(ctx, qtx, where)
}

func createAttachableItem(
	ctx context.Context, qtx *query.Queries, param parameter.CreateAttachableItemParam,
) (entity.AttachableItem, error) {
	p := query.CreateAttachableItemParams{
		Url:        param.URL,
		Alias:      param.Alias,
		Size:       pgtype.Float8(param.Size),
		OwnerID:    pgtype.UUID(param.OwnerID),
		FromOuter:  param.FromOuter,
		MimeTypeID: param.MimeTypeID,
	}
	attachableItem, err := qtx.CreateAttachableItem(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.AttachableItem{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItem{}, fmt.Errorf("failed to create attachable item: %w", err)
	}
	entity := entity.AttachableItem{
		AttachableItemID: attachableItem.AttachableItemID,
		OwnerID:          entity.UUID(attachableItem.OwnerID),
		FromOuter:        attachableItem.FromOuter,
		URL:              attachableItem.Url,
		Alias:            attachableItem.Alias,
		Size:             entity.Float(attachableItem.Size),
		MimeTypeID:       attachableItem.MimeTypeID,
	}
	return entity, nil
}

// CreateAttachableItem 添付可能アイテムを作成する。
func (a *PgAdapter) CreateAttachableItem(
	ctx context.Context, param parameter.CreateAttachableItemParam,
) (entity.AttachableItem, error) {
	return createAttachableItem(ctx, a.query, param)
}

// CreateAttachableItemWithSd SD付きで添付可能アイテムを作成する。
func (a *PgAdapter) CreateAttachableItemWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateAttachableItemParam,
) (entity.AttachableItem, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItem{}, store.ErrNotFoundDescriptor
	}
	return createAttachableItem(ctx, qtx, param)
}

func createAttachableItems(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateAttachableItemParam,
) (int64, error) {
	var p []query.CreateAttachableItemsParams
	for _, param := range params {
		p = append(p, query.CreateAttachableItemsParams{
			Url:        param.URL,
			Size:       pgtype.Float8(param.Size),
			Alias:      param.Alias,
			OwnerID:    pgtype.UUID(param.OwnerID),
			FromOuter:  param.FromOuter,
			MimeTypeID: param.MimeTypeID,
		})
	}
	n, err := qtx.CreateAttachableItems(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("attachable item")
		}
		return 0, fmt.Errorf("failed to create attachable items: %w", err)
	}
	return n, nil
}

// CreateAttachableItems 添付可能アイテムを作成する。
func (a *PgAdapter) CreateAttachableItems(
	ctx context.Context, params []parameter.CreateAttachableItemParam,
) (int64, error) {
	return createAttachableItems(ctx, a.query, params)
}

// CreateAttachableItemsWithSd SD付きで添付可能アイテムを作成する。
func (a *PgAdapter) CreateAttachableItemsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateAttachableItemParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createAttachableItems(ctx, qtx, params)
}

func deleteAttachableItem(
	ctx context.Context, qtx *query.Queries, attachableItemID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteAttachableItem(ctx, attachableItemID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable item: %w", err)
	}
	if c != 1 {
		return c, errhandle.NewModelNotFoundError("attachable item")
	}
	return c, nil
}

// DeleteAttachableItem 添付可能アイテムを削除する。
func (a *PgAdapter) DeleteAttachableItem(ctx context.Context, attachableItemID uuid.UUID) (int64, error) {
	return deleteAttachableItem(ctx, a.query, attachableItemID)
}

// DeleteAttachableItemWithSd SD付きで添付可能アイテムを削除する。
func (a *PgAdapter) DeleteAttachableItemWithSd(
	ctx context.Context, sd store.Sd, attachableItemID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteAttachableItem(ctx, qtx, attachableItemID)
}

func pluralDeleteAttachableItems(
	ctx context.Context, qtx *query.Queries, attachableItemIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteAttachableItems(ctx, attachableItemIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable items: %w", err)
	}
	if c != int64(len(attachableItemIDs)) {
		return c, errhandle.NewModelNotFoundError("attachable item")
	}
	return c, nil
}

// PluralDeleteAttachableItems 添付可能アイテムを複数削除する。
func (a *PgAdapter) PluralDeleteAttachableItems(
	ctx context.Context, attachableItemIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteAttachableItems(ctx, a.query, attachableItemIDs)
}

// PluralDeleteAttachableItemsWithSd SD付きで添付可能アイテムを複数削除する。
func (a *PgAdapter) PluralDeleteAttachableItemsWithSd(
	ctx context.Context, sd store.Sd, attachableItemIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteAttachableItems(ctx, qtx, attachableItemIDs)
}

func findAttachableItemByID(
	ctx context.Context, qtx *query.Queries, attachableItemID uuid.UUID,
) (entity.AttachableItemWithContent, error) {
	attachableItem, err := qtx.FindAttachableItemByID(ctx, attachableItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttachableItemWithContent{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItemWithContent{}, fmt.Errorf("failed to find attachable item: %w", err)
	}
	return attachableItemConv(attachableItem), nil
}

// FindAttachableItemByID 添付可能アイテムを取得する。
func (a *PgAdapter) FindAttachableItemByID(
	ctx context.Context, attachableItemID uuid.UUID,
) (entity.AttachableItemWithContent, error) {
	return findAttachableItemByID(ctx, a.query, attachableItemID)
}

// FindAttachableItemByIDWithSd SD付きで添付可能アイテムを取得する。
func (a *PgAdapter) FindAttachableItemByIDWithSd(
	ctx context.Context, sd store.Sd, attachableItemID uuid.UUID,
) (entity.AttachableItemWithContent, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItemWithContent{}, store.ErrNotFoundDescriptor
	}
	return findAttachableItemByID(ctx, qtx, attachableItemID)
}

func findAttachableItemByIDWithMimeType(
	ctx context.Context, qtx *query.Queries, attachableItemID uuid.UUID,
) (entity.AttachableItemWithMimeType, error) {
	attachableItem, err := qtx.FindAttachableItemByIDWithMimeType(ctx, attachableItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttachableItemWithMimeType{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItemWithMimeType{}, fmt.Errorf("failed to find attachable item: %w", err)
	}
	return attachableItemConvWithMimeType(attachableItem), nil
}

// FindAttachableItemByIDWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) FindAttachableItemByIDWithMimeType(
	ctx context.Context, attachableItemID uuid.UUID,
) (entity.AttachableItemWithMimeType, error) {
	return findAttachableItemByIDWithMimeType(ctx, a.query, attachableItemID)
}

// FindAttachableItemByIDWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) FindAttachableItemByIDWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, attachableItemID uuid.UUID,
) (entity.AttachableItemWithMimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItemWithMimeType{}, store.ErrNotFoundDescriptor
	}
	return findAttachableItemByIDWithMimeType(ctx, qtx, attachableItemID)
}

func findAttachableItemByURL(
	ctx context.Context, qtx *query.Queries, url string,
) (entity.AttachableItemWithContent, error) {
	attachableItem, err := qtx.FindAttachableItemByURL(ctx, url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttachableItemWithContent{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItemWithContent{}, fmt.Errorf("failed to find attachable item: %w", err)
	}
	return attachableItemConv(query.FindAttachableItemByIDRow(attachableItem)), nil
}

// FindAttachableItemByURL 添付可能アイテムを取得する。
func (a *PgAdapter) FindAttachableItemByURL(
	ctx context.Context, url string,
) (entity.AttachableItemWithContent, error) {
	return findAttachableItemByURL(ctx, a.query, url)
}

// FindAttachableItemByURLWithSd SD付きで添付可能アイテムを取得する。
func (a *PgAdapter) FindAttachableItemByURLWithSd(
	ctx context.Context, sd store.Sd, url string,
) (entity.AttachableItemWithContent, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItemWithContent{}, store.ErrNotFoundDescriptor
	}
	return findAttachableItemByURL(ctx, qtx, url)
}

func findAttachableItemByURLWithMimeType(
	ctx context.Context, qtx *query.Queries, url string,
) (entity.AttachableItemWithMimeType, error) {
	attachableItem, err := qtx.FindAttachableItemByURLWithMimeType(ctx, url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttachableItemWithMimeType{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItemWithMimeType{}, fmt.Errorf("failed to find attachable item: %w", err)
	}
	return attachableItemConvWithMimeType(query.FindAttachableItemByIDWithMimeTypeRow(attachableItem)), nil
}

// FindAttachableItemByURLWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) FindAttachableItemByURLWithMimeType(
	ctx context.Context, url string,
) (entity.AttachableItemWithMimeType, error) {
	return findAttachableItemByURLWithMimeType(ctx, a.query, url)
}

// FindAttachableItemByURLWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) FindAttachableItemByURLWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, url string,
) (entity.AttachableItemWithMimeType, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItemWithMimeType{}, store.ErrNotFoundDescriptor
	}
	return findAttachableItemByURLWithMimeType(ctx, qtx, url)
}

func getAttachableItems(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	eConvFunc := func(e entity.AttachableItemWithContentForQuery) (entity.AttachableItemWithContent, error) {
		return e.AttachableItemWithContent, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttachableItemsParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
		}
		r, err := qtx.CountAttachableItems(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attachable items: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.AttachableItemWithContentForQuery, error) {
		p := query.GetAttachableItemsParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
		}
		r, err := qtx.GetAttachableItems(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithContentForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConv(query.FindAttachableItemByIDRow(e))
			enity[i] = entity.AttachableItemWithContentForQuery{
				Pkey:                      entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithContent: convE,
			}
		}
		return enity, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.AttachableItemWithContentForQuery, error) {
		p := query.GetAttachableItemsUseKeysetPaginateParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
			Limit:              limit,
			CursorDirection:    cursorDir,
			Cursor:             cursor,
		}
		r, err := qtx.GetAttachableItemsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithContentForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConv(query.FindAttachableItemByIDRow(e))
			enity[i] = entity.AttachableItemWithContentForQuery{
				Pkey:                      entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithContent: convE,
			}
		}
		return enity, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.AttachableItemWithContentForQuery, error) {
		p := query.GetAttachableItemsUseNumberedPaginateParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
			Limit:              limit,
			Offset:             offset,
		}
		r, err := qtx.GetAttachableItemsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithContentForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConv(query.FindAttachableItemByIDRow(e))
			enity[i] = entity.AttachableItemWithContentForQuery{
				Pkey:                      entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithContent: convE,
			}
		}
		return enity, nil
	}
	selector := func(subCursor string, e entity.AttachableItemWithContentForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.AttendStatusDefaultCursorKey:
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
		return store.ListResult[entity.AttachableItemWithContent]{}, fmt.Errorf("failed to get attachable items: %w", err)
	}
	return res, nil
}

// GetAttachableItems 添付可能アイテムを取得する。
func (a *PgAdapter) GetAttachableItems(
	ctx context.Context, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	return getAttachableItems(ctx, a.query, where, order, np, cp, wc)
}

// GetAttachableItemsWithSd SD付きで添付可能アイテムを取得する。
func (a *PgAdapter) GetAttachableItemsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachableItemWithContent]{}, store.ErrNotFoundDescriptor
	}
	return getAttachableItems(ctx, qtx, where, order, np, cp, wc)
}

func getAttachableItemsWithMimeType(
	ctx context.Context, qtx *query.Queries, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	eConvFunc := func(e entity.AttachableItemWithMimeTypeForQuery) (entity.AttachableItemWithMimeType, error) {
		return e.AttachableItemWithMimeType, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttachableItemsParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
		}
		r, err := qtx.CountAttachableItems(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attachable items: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.AttachableItemWithMimeTypeForQuery, error) {
		p := query.GetAttachableItemsWithMimeTypeParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
		}
		r, err := qtx.GetAttachableItemsWithMimeType(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithMimeTypeForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConvWithMimeType(query.FindAttachableItemByIDWithMimeTypeRow(e))
			enity[i] = entity.AttachableItemWithMimeTypeForQuery{
				Pkey:                       entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithMimeType: convE,
			}
		}
		return enity, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.AttachableItemWithMimeTypeForQuery, error) {
		p := query.GetAttachableItemsWithMimeTypeUseKeysetPaginateParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
			Limit:              limit,
			CursorDirection:    cursorDir,
			Cursor:             cursor,
		}
		r, err := qtx.GetAttachableItemsWithMimeTypeUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithMimeTypeForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConvWithMimeType(query.FindAttachableItemByIDWithMimeTypeRow(e))
			enity[i] = entity.AttachableItemWithMimeTypeForQuery{
				Pkey:                       entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithMimeType: convE,
			}
		}
		return enity, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.AttachableItemWithMimeTypeForQuery, error) {
		p := query.GetAttachableItemsWithMimeTypeUseNumberedPaginateParams{
			WhereInMimeTypeIds: where.WhereInMimeType,
			InMimeTypeIds:      where.InMimeTypes,
			WhereInOwnerIds:    where.WhereInOwner,
			InOwnerIds:         where.InOwners,
			Limit:              limit,
			Offset:             offset,
		}
		r, err := qtx.GetAttachableItemsWithMimeTypeUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attachable items: %w", err)
		}
		enity := make([]entity.AttachableItemWithMimeTypeForQuery, len(r))
		for i, e := range r {
			convE := attachableItemConvWithMimeType(query.FindAttachableItemByIDWithMimeTypeRow(e))
			enity[i] = entity.AttachableItemWithMimeTypeForQuery{
				Pkey:                       entity.Int(e.TAttachableItemsPkey),
				AttachableItemWithMimeType: convE,
			}
		}
		return enity, nil
	}
	selector := func(subCursor string, e entity.AttachableItemWithMimeTypeForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.AttendStatusDefaultCursorKey:
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
		return store.ListResult[entity.AttachableItemWithMimeType]{}, fmt.Errorf("failed to get attachable items: %w", err)
	}
	return res, nil
}

// GetAttachableItemsWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetAttachableItemsWithMimeType(
	ctx context.Context, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	return getAttachableItemsWithMimeType(ctx, a.query, where, order, np, cp, wc)
}

// GetAttachableItemsWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetAttachableItemsWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereAttachableItemParam,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachableItemWithMimeType]{}, store.ErrNotFoundDescriptor
	}
	return getAttachableItemsWithMimeType(ctx, qtx, where, order, np, cp, wc)
}

func getPluralAttachableItems(
	ctx context.Context, qtx *query.Queries, attachableItemIDs []uuid.UUID,
	_ parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	var e []query.GetPluralAttachableItemsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralAttachableItems(ctx, attachableItemIDs)
	} else {
		var qe []query.GetPluralAttachableItemsUseNumberedPaginateRow
		qe, err = qtx.GetPluralAttachableItemsUseNumberedPaginate(
			ctx, query.GetPluralAttachableItemsUseNumberedPaginateParams{
				AttachableItemIds: attachableItemIDs,
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
			})
		e = make([]query.GetPluralAttachableItemsRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralAttachableItemsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.AttachableItemWithContent]{}, fmt.Errorf("failed to get attachable items: %w", err)
	}
	entities := make([]entity.AttachableItemWithContent, len(e))
	for i, v := range e {
		entities[i] = attachableItemConv(query.FindAttachableItemByIDRow(v))
	}
	return store.ListResult[entity.AttachableItemWithContent]{Data: entities}, nil
}

// GetPluralAttachableItems 添付可能アイテムを取得する。
func (a *PgAdapter) GetPluralAttachableItems(
	ctx context.Context, attachableItemIDs []uuid.UUID,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	return getPluralAttachableItems(ctx, a.query, attachableItemIDs, order, np)
}

// GetPluralAttachableItemsWithSd SD付きで添付可能アイテムを取得する。
func (a *PgAdapter) GetPluralAttachableItemsWithSd(
	ctx context.Context, sd store.Sd, attachableItemIDs []uuid.UUID,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithContent], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachableItemWithContent]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAttachableItems(ctx, qtx, attachableItemIDs, order, np)
}

func getPluralAttachableItemsWithMimeType(
	ctx context.Context, qtx *query.Queries, attachableItemIDs []uuid.UUID,
	_ parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	var e []query.GetPluralAttachableItemsWithMimeTypeRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralAttachableItemsWithMimeType(ctx, attachableItemIDs)
	} else {
		var qe []query.GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateRow
		qe, err = qtx.GetPluralAttachableItemsWithMimeTypeUseNumberedPaginate(
			ctx, query.GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateParams{
				AttachableItemIds: attachableItemIDs,
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
			})
		e = make([]query.GetPluralAttachableItemsWithMimeTypeRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralAttachableItemsWithMimeTypeRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.AttachableItemWithMimeType]{}, fmt.Errorf("failed to get attachable items: %w", err)
	}
	entities := make([]entity.AttachableItemWithMimeType, len(e))
	for i, v := range e {
		entities[i] = attachableItemConvWithMimeType(query.FindAttachableItemByIDWithMimeTypeRow(v))
	}
	return store.ListResult[entity.AttachableItemWithMimeType]{Data: entities}, nil
}

// GetPluralAttachableItemsWithMimeType 添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetPluralAttachableItemsWithMimeType(
	ctx context.Context, attachableItemIDs []uuid.UUID,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	return getPluralAttachableItemsWithMimeType(ctx, a.query, attachableItemIDs, order, np)
}

// GetPluralAttachableItemsWithMimeTypeWithSd SD付きで添付可能アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetPluralAttachableItemsWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, attachableItemIDs []uuid.UUID,
	order parameter.AttachableItemOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachableItemWithMimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachableItemWithMimeType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAttachableItemsWithMimeType(ctx, qtx, attachableItemIDs, order, np)
}

func updateAttachableItem(
	ctx context.Context, qtx *query.Queries, attachableItemID uuid.UUID, param parameter.UpdateAttachableItemParams,
) (entity.AttachableItem, error) {
	p := query.UpdateAttachableItemParams{
		AttachableItemID: attachableItemID,
		Url:              param.URL,
		Size:             pgtype.Float8(param.Size),
		Alias:            param.Alias,
		MimeTypeID:       param.MimeTypeID,
	}
	attachableItem, err := qtx.UpdateAttachableItem(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AttachableItem{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachableItem{}, fmt.Errorf("failed to update attachable item: %w", err)
	}
	entity := entity.AttachableItem{
		AttachableItemID: attachableItem.AttachableItemID,
		OwnerID:          entity.UUID(attachableItem.OwnerID),
		FromOuter:        attachableItem.FromOuter,
		URL:              attachableItem.Url,
		Alias:            attachableItem.Alias,
		Size:             entity.Float(attachableItem.Size),
		MimeTypeID:       attachableItem.MimeTypeID,
	}
	return entity, nil
}

// UpdateAttachableItem 添付可能アイテムを更新する。
func (a *PgAdapter) UpdateAttachableItem(
	ctx context.Context, attachableItemID uuid.UUID,
	param parameter.UpdateAttachableItemParams,
) (entity.AttachableItem, error) {
	return updateAttachableItem(ctx, a.query, attachableItemID, param)
}

// UpdateAttachableItemWithSd SD付きで添付可能アイテムを更新する。
func (a *PgAdapter) UpdateAttachableItemWithSd(
	ctx context.Context, sd store.Sd, attachableItemID uuid.UUID,
	param parameter.UpdateAttachableItemParams,
) (entity.AttachableItem, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachableItem{}, store.ErrNotFoundDescriptor
	}
	return updateAttachableItem(ctx, qtx, attachableItemID, param)
}

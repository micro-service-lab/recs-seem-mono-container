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

func countAttachedItemsOnChatRoom(ctx context.Context,
	qtx *query.Queries, chatRoomID uuid.UUID, where parameter.WhereAttachedItemOnChatRoomParam,
) (int64, error) {
	p := query.CountAttachedItemsOnChatRoomParams{
		ChatRoomID:      chatRoomID,
		WhereInMimeType: where.WhereInMimeType,
		InMimeTypes:     where.InMimeTypes,
		WhereIsImage:    where.WhereIsImage,
		WhereIsFile:     where.WhereIsFile,
	}
	c, err := qtx.CountAttachedItemsOnChatRoom(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count attached items on chat room: %w", err)
	}
	return c, nil
}

// CountAttachedItemsOnChatRoom チャットルームに関連付けられた添付アイテム数を取得する。
func (a *PgAdapter) CountAttachedItemsOnChatRoom(
	ctx context.Context, chatRoomID uuid.UUID, where parameter.WhereAttachedItemOnChatRoomParam,
) (int64, error) {
	return countAttachedItemsOnChatRoom(ctx, a.query, chatRoomID, where)
}

// CountAttachedItemsOnChatRoomWithSd SD付きでチャットルームに関連付けられた添付アイテム数を取得する。
func (a *PgAdapter) CountAttachedItemsOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID, where parameter.WhereAttachedItemOnChatRoomParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countAttachedItemsOnChatRoom(ctx, qtx, chatRoomID, where)
}

func countAttachedItemsOnMessage(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
) (int64, error) {
	p := query.CountAttachedItemsOnMessageParams{
		MessageID:       messageID,
		WhereInMimeType: where.WhereInMimeType,
		InMimeTypes:     where.InMimeTypes,
		WhereIsImage:    where.WhereIsImage,
		WhereIsFile:     where.WhereIsFile,
	}
	c, err := qtx.CountAttachedItemsOnMessage(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count attached items on message: %w", err)
	}
	return c, nil
}

// CountAttachedItemsOnMessage メッセージに関連付けられた添付アイテム数を取得する。
func (a *PgAdapter) CountAttachedItemsOnMessage(
	ctx context.Context, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
) (int64, error) {
	return countAttachedItemsOnMessage(ctx, a.query, messageID, where)
}

// CountAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテム数を取得する。
func (a *PgAdapter) CountAttachedItemsOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countAttachedItemsOnMessage(ctx, qtx, messageID, where)
}

func attachItemOnMessage(
	ctx context.Context, qtx *query.Queries, param parameter.AttachItemMessageParam,
) (entity.AttachedMessage, error) {
	p := query.CreateAttachedMessageParams{
		MessageID:        param.MessageID,
		AttachableItemID: pgtype.UUID(param.AttachableItemID),
	}
	attachedMessage, err := qtx.CreateAttachedMessage(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.AttachedMessage{}, errhandle.NewModelNotFoundError("attachable item")
		}
		return entity.AttachedMessage{}, fmt.Errorf("failed to attach item on message: %w", err)
	}
	entity := entity.AttachedMessage{
		AttachedMessageID: attachedMessage.AttachedMessageID,
		MessageID:         attachedMessage.MessageID,
		AttachableItemID:  entity.UUID(attachedMessage.AttachableItemID),
	}
	return entity, nil
}

// AttacheItemOnMessage 添付アイテムを関連付ける。
func (a *PgAdapter) AttacheItemOnMessage(
	ctx context.Context, param parameter.AttachItemMessageParam,
) (entity.AttachedMessage, error) {
	return attachItemOnMessage(ctx, a.query, param)
}

// AttacheItemOnMessageWithSd SD付きで添付アイテムを関連付ける。
func (a *PgAdapter) AttacheItemOnMessageWithSd(
	ctx context.Context, sd store.Sd, param parameter.AttachItemMessageParam,
) (entity.AttachedMessage, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.AttachedMessage{}, store.ErrNotFoundDescriptor
	}
	return attachItemOnMessage(ctx, qtx, param)
}

func attachItemsOnMessages(
	ctx context.Context, qtx *query.Queries, params []parameter.AttachItemMessageParam,
) (int64, error) {
	var p []query.CreateAttachedMessagesParams
	for _, param := range params {
		p = append(p, query.CreateAttachedMessagesParams{
			MessageID:        param.MessageID,
			AttachableItemID: pgtype.UUID(param.AttachableItemID),
		})
	}
	attachedMessages, err := qtx.CreateAttachedMessages(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelNotFoundError("attachable item")
		}
		return 0, fmt.Errorf("failed to attach items on messages: %w", err)
	}
	return attachedMessages, nil
}

// AttacheItemsOnMessages 添付アイテムを複数関連付ける。
func (a *PgAdapter) AttacheItemsOnMessages(
	ctx context.Context, params []parameter.AttachItemMessageParam,
) (int64, error) {
	return attachItemsOnMessages(ctx, a.query, params)
}

// AttacheItemsOnMessagesWithSd SD付きで添付アイテムを複数関連付ける。
func (a *PgAdapter) AttacheItemsOnMessagesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.AttachItemMessageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return attachItemsOnMessages(ctx, qtx, params)
}

func detachAttachedMessage(ctx context.Context, qtx *query.Queries, attachedMessageID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAttachedMessage(ctx, attachedMessageID)
	if err != nil {
		return 0, fmt.Errorf("failed to detach item on message: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("attached item")
	}
	return c, nil
}

// DetachAttachedMessage 添付アイテムを関連付けを解除する。
func (a *PgAdapter) DetachAttachedMessage(ctx context.Context, attachedMessageID uuid.UUID) (int64, error) {
	return detachAttachedMessage(ctx, a.query, attachedMessageID)
}

// DetachAttachedMessageWithSd SD付きで添付アイテムを関連付けを解除する。
func (a *PgAdapter) DetachAttachedMessageWithSd(
	ctx context.Context, sd store.Sd, attachedMessageID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return detachAttachedMessage(ctx, qtx, attachedMessageID)
}

func detachItemsOnMessage(ctx context.Context, qtx *query.Queries, messageID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAttachedMessagesOnMessage(ctx, messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to detach items on message: %w", err)
	}
	return c, nil
}

// DetachItemsOnMessage メッセージに関連付けられた添付アイテムを解除する。
func (a *PgAdapter) DetachItemsOnMessage(ctx context.Context, messageID uuid.UUID) (int64, error) {
	return detachItemsOnMessage(ctx, a.query, messageID)
}

// DetachItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテムを解除する。
func (a *PgAdapter) DetachItemsOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return detachItemsOnMessage(ctx, qtx, messageID)
}

func detachItemsOnMessages(ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID) (int64, error) {
	c, err := qtx.DeleteAttachedMessagesOnMessages(ctx, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to detach items on messages: %w", err)
	}
	return c, nil
}

// DetachItemsOnMessages メッセージに関連付けられた添付アイテムを複数解除する。
func (a *PgAdapter) DetachItemsOnMessages(ctx context.Context, messageIDs []uuid.UUID) (int64, error) {
	return detachItemsOnMessages(ctx, a.query, messageIDs)
}

// DetachItemsOnMessagesWithSd SD付きでメッセージに関連付けられた添付アイテムを複数解除する。
func (a *PgAdapter) DetachItemsOnMessagesWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return detachItemsOnMessages(ctx, qtx, messageIDs)
}

func getAttachedItemsOnChatRoom(
	ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID,
	where parameter.WhereAttachedItemOnChatRoomParam,
	order parameter.AttachedItemOnChatRoomOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnChatRoom], error) {
	eConvFunc := func(e entity.AttachedItemOnChatRoomForQuery) (entity.AttachedItemOnChatRoom, error) {
		return e.AttachedItemOnChatRoom, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttachedItemsOnChatRoomParams{
			ChatRoomID:      chatRoomID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.CountAttachedItemsOnChatRoom(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attached items on chat room: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.AttachedItemOnChatRoomForQuery, error) {
		p := query.GetAttachedItemsOnChatRoomParams{
			ChatRoomID:      chatRoomID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.GetAttachedItemsOnChatRoom(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.AttachedItemOnChatRoomForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get attached items on chat room: %w", err)
		}
		fq := make([]entity.AttachedItemOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnChatRoomForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnChatRoom: entity.AttachedItemOnChatRoom{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
					Message: entity.Message{
						MessageID:    e.MessageID,
						ChatRoomID:   chatRoomID,
						Body:         e.MessageBody.String,
						SenderID:     entity.UUID(e.MessageSenderID),
						PostedAt:     e.MessagePostedAt.Time,
						LastEditedAt: e.MessageLastEditedAt.Time,
					},
				},
			}
		}
		return fq, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.AttachedItemOnChatRoomForQuery, error) {
		p := query.GetAttachedItemsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:      chatRoomID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
		}
		r, err := qtx.GetAttachedItemsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on chat room: %w", err)
		}
		fq := make([]entity.AttachedItemOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnChatRoomForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnChatRoom: entity.AttachedItemOnChatRoom{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
					Message: entity.Message{
						MessageID:    e.MessageID,
						ChatRoomID:   chatRoomID,
						Body:         e.MessageBody.String,
						SenderID:     entity.UUID(e.MessageSenderID),
						PostedAt:     e.MessagePostedAt.Time,
						LastEditedAt: e.MessageLastEditedAt.Time,
					},
				},
			}
		}
		return fq, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.AttachedItemOnChatRoomForQuery, error) {
		p := query.GetAttachedItemsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID:      chatRoomID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			Limit:           limit,
			Offset:          offset,
		}
		r, err := qtx.GetAttachedItemsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on chat room: %w", err)
		}
		fq := make([]entity.AttachedItemOnChatRoomForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnChatRoomForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnChatRoom: entity.AttachedItemOnChatRoom{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
					Message: entity.Message{
						MessageID:    e.MessageID,
						ChatRoomID:   chatRoomID,
						Body:         e.MessageBody.String,
						SenderID:     entity.UUID(e.MessageSenderID),
						PostedAt:     e.MessagePostedAt.Time,
						LastEditedAt: e.MessageLastEditedAt.Time,
					},
				},
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.AttachedItemOnChatRoomForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.AttachedItemOnChatRoomDefaultCursorKey:
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
		return store.ListResult[entity.AttachedItemOnChatRoom]{},
			fmt.Errorf("failed to get attached items on chat room: %w", err)
	}
	return res, nil
}

// GetAttachedItemsOnChatRoom チャットルームに関連付けられた添付アイテムを取得する。
func (a *PgAdapter) GetAttachedItemsOnChatRoom(
	ctx context.Context, chatRoomID uuid.UUID,
	where parameter.WhereAttachedItemOnChatRoomParam, order parameter.AttachedItemOnChatRoomOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnChatRoom], error) {
	return getAttachedItemsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetAttachedItemsOnChatRoomWithSd SD付きでチャットルームに関連付けられた添付アイテムを取得する。
func (a *PgAdapter) GetAttachedItemsOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
	where parameter.WhereAttachedItemOnChatRoomParam, order parameter.AttachedItemOnChatRoomOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachedItemOnChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getAttachedItemsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

func getAttachedItemsOnMessage(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID,
	where parameter.WhereAttachedItemOnMessageParam, order parameter.AttachedItemOnMessageOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	eConvFunc := func(e entity.AttachedItemOnMessageForQuery) (entity.AttachedItemOnMessage, error) {
		return e.AttachedItemOnMessage, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttachedItemsOnMessageParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.CountAttachedItemsOnMessage(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attached items on message: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.AttachedItemOnMessageForQuery, error) {
		p := query.GetAttachedItemsOnMessageParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.GetAttachedItemsOnMessage(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.AttachedItemOnMessageForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnMessageForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnMessage: entity.AttachedItemOnMessage{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.AttachedItemOnMessageForQuery, error) {
		p := query.GetAttachedItemsOnMessageUseKeysetPaginateParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
		}
		r, err := qtx.GetAttachedItemsOnMessageUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnMessageForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnMessage: entity.AttachedItemOnMessage{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.AttachedItemOnMessageForQuery, error) {
		p := query.GetAttachedItemsOnMessageUseNumberedPaginateParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			Limit:           limit,
			Offset:          offset,
		}
		r, err := qtx.GetAttachedItemsOnMessageUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageForQuery, len(r))
		for i, e := range r {
			fq[i] = entity.AttachedItemOnMessageForQuery{
				Pkey: entity.Int(e.TAttachedMessagesPkey),
				AttachedItemOnMessage: entity.AttachedItemOnMessage{
					AttachedMessageID: e.AttachedMessageID,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.AttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.AttachedItemOwnerID),
						FromOuter:        e.AttachedItemFromOuter.Bool,
						URL:              e.AttachedItemUrl.String,
						Size:             entity.Float(e.AttachedItemSize),
						MimeTypeID:       e.AttachedItemMimeTypeID.Bytes,
					},
				},
			}
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.AttachedItemOnMessageForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.AttachedItemOnMessageDefaultCursorKey:
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
		return store.ListResult[entity.AttachedItemOnMessage]{},
			fmt.Errorf("failed to get attached items on message: %w", err)
	}
	return res, nil
}

// GetAttachedItemsOnMessage メッセージに関連付けられた添付アイテムを取得する。
func (a *PgAdapter) GetAttachedItemsOnMessage(
	ctx context.Context, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
	order parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	return getAttachedItemsOnMessage(ctx, a.query, messageID, where, order, np, cp, wc)
}

// GetAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた添付アイテムを取得する。
func (a *PgAdapter) GetAttachedItemsOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
	where parameter.WhereAttachedItemOnMessageParam, order parameter.AttachedItemOnMessageOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachedItemOnMessage]{}, store.ErrNotFoundDescriptor
	}
	return getAttachedItemsOnMessage(ctx, qtx, messageID, where, order, np, cp, wc)
}

func getPluralAttachedItemsOnMessage(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	_ parameter.AttachedItemOnChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	var e []query.GetPluralAttachedItemsOnMessageRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralAttachedItemsOnMessage(ctx, messageIDs)
	} else {
		var qe []query.GetPluralAttachedItemsOnMessageUseNumberedPaginateRow
		qe, err = qtx.GetPluralAttachedItemsOnMessageUseNumberedPaginate(
			ctx, query.GetPluralAttachedItemsOnMessageUseNumberedPaginateParams{
				MessageIds: messageIDs,
				Limit:      int32(np.Limit.Int64),
				Offset:     int32(np.Offset.Int64),
			})
		e = make([]query.GetPluralAttachedItemsOnMessageRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralAttachedItemsOnMessageRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.AttachedItemOnMessage]{},
			fmt.Errorf("failed to get attached items on message: %w", err)
	}
	entities := make([]entity.AttachedItemOnMessage, len(e))
	for i, v := range e {
		entities[i] = entity.AttachedItemOnMessage{
			AttachedMessageID: v.AttachedMessageID,
			AttachableItem: entity.AttachableItem{
				AttachableItemID: v.AttachableItemID.Bytes,
				OwnerID:          entity.UUID(v.AttachedItemOwnerID),
				FromOuter:        v.AttachedItemFromOuter.Bool,
				URL:              v.AttachedItemUrl.String,
				Size:             entity.Float(v.AttachedItemSize),
				MimeTypeID:       v.AttachedItemMimeTypeID.Bytes,
			},
		}
	}
	return store.ListResult[entity.AttachedItemOnMessage]{Data: entities}, nil
}

// GetPluralAttachedItemsOnMessage メッセージに関連付けられた複数の添付アイテムを取得する。
func (a *PgAdapter) GetPluralAttachedItemsOnMessage(
	ctx context.Context, messageIDs []uuid.UUID,
	order parameter.AttachedItemOnChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	return getPluralAttachedItemsOnMessage(ctx, a.query, messageIDs, order, np)
}

// GetPluralAttachedItemsOnMessageWithSd SD付きでメッセージに関連付けられた複数の添付アイテムを取得する。
func (a *PgAdapter) GetPluralAttachedItemsOnMessageWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	order parameter.AttachedItemOnChatRoomOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessage], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachedItemOnMessage]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAttachedItemsOnMessage(ctx, qtx, messageIDs, order, np)
}

func getAttachedItemsOnMessageWithMimeType(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
	order parameter.AttachedItemOnMessageOrderMethod,
	np store.NumberedPaginationParam, cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	eConvFunc := func(
		e entity.AttachedItemOnMessageWithMimeTypeForQuery,
	) (entity.AttachedItemOnMessageWithMimeType, error) {
		return e.AttachedItemOnMessageWithMimeType, nil
	}
	runCFunc := func() (int64, error) {
		p := query.CountAttachedItemsOnMessageParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.CountAttachedItemsOnMessage(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count attached items on message: %w", err)
		}
		return r, nil
	}
	conv := func(e query.GetAttachedItemsOnMessageWithMimeTypeRow) entity.AttachedItemOnMessageWithMimeTypeForQuery {
		mimeType := entity.MimeType{
			MimeTypeID: e.AttachedItemMimeTypeID.Bytes,
			Name:       e.MimeTypeName.String,
			Key:        e.MimeTypeKey.String,
			Kind:       e.MimeTypeKind.String,
		}
		return entity.AttachedItemOnMessageWithMimeTypeForQuery{
			Pkey: entity.Int(e.TAttachedMessagesPkey),
			AttachedItemOnMessageWithMimeType: entity.AttachedItemOnMessageWithMimeType{
				AttachedMessageID: e.AttachedMessageID,
				AttachableItem: entity.AttachableItemWithMimeType{
					AttachableItemID: e.AttachableItemID.Bytes,
					OwnerID:          entity.UUID(e.AttachedItemOwnerID),
					FromOuter:        e.AttachedItemFromOuter.Bool,
					URL:              e.AttachedItemUrl.String,
					Size:             entity.Float(e.AttachedItemSize),
					MimeType:         entity.NullableEntity[entity.MimeType]{Entity: mimeType},
				},
			},
		}
	}
	runQFunc := func(_ string) ([]entity.AttachedItemOnMessageWithMimeTypeForQuery, error) {
		p := query.GetAttachedItemsOnMessageWithMimeTypeParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
		}
		r, err := qtx.GetAttachedItemsOnMessageWithMimeType(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.AttachedItemOnMessageWithMimeTypeForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageWithMimeTypeForQuery, len(r))
		for i, e := range r {
			fq[i] = conv(e)
		}
		return fq, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.AttachedItemOnMessageWithMimeTypeForQuery, error) {
		p := query.GetAttachedItemsOnMessageWithMimeTypeUseKeysetPaginateParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			CursorDirection: cursorDir,
			Cursor:          cursor,
			Limit:           limit,
		}
		r, err := qtx.GetAttachedItemsOnMessageWithMimeTypeUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageWithMimeTypeForQuery, len(r))
		for i, e := range r {
			fq[i] = conv(query.GetAttachedItemsOnMessageWithMimeTypeRow(e))
		}
		return fq, nil
	}
	runQNPFunc := func(
		_ string, limit, offset int32,
	) ([]entity.AttachedItemOnMessageWithMimeTypeForQuery, error) {
		p := query.GetAttachedItemsOnMessageWithMimeTypeUseNumberedPaginateParams{
			MessageID:       messageID,
			WhereInMimeType: where.WhereInMimeType,
			InMimeTypes:     where.InMimeTypes,
			WhereIsImage:    where.WhereIsImage,
			WhereIsFile:     where.WhereIsFile,
			Limit:           limit,
			Offset:          offset,
		}
		r, err := qtx.GetAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		fq := make([]entity.AttachedItemOnMessageWithMimeTypeForQuery, len(r))
		for i, e := range r {
			fq[i] = conv(query.GetAttachedItemsOnMessageWithMimeTypeRow(e))
		}
		return fq, nil
	}
	selector := func(subCursor string, e entity.AttachedItemOnMessageWithMimeTypeForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.AttachedItemOnMessageDefaultCursorKey:
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
		return store.ListResult[entity.AttachedItemOnMessageWithMimeType]{},
			fmt.Errorf("failed to get attached items on message: %w", err)
	}
	return res, nil
}

// GetAttachedItemsOnMessageWithMimeType メッセージに関連付けられた添付アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetAttachedItemsOnMessageWithMimeType(
	ctx context.Context, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
	order parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	return getAttachedItemsOnMessageWithMimeType(ctx, a.query, messageID, where, order, np, cp, wc)
}

// GetAttachedItemsOnMessageWithMimeTypeWithSd SD付きでメッセージに関連付けられた添付アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetAttachedItemsOnMessageWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID, where parameter.WhereAttachedItemOnMessageParam,
	order parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachedItemOnMessageWithMimeType]{}, store.ErrNotFoundDescriptor
	}
	return getAttachedItemsOnMessageWithMimeType(ctx, qtx, messageID, where, order, np, cp, wc)
}

func getPluralAttachedItemsOnMessageWithMimeType(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	_ parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	var e []query.GetPluralAttachedItemsOnMessageWithMimeTypeRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralAttachedItemsOnMessageWithMimeType(ctx, messageIDs)
	} else {
		var qe []query.GetPluralAttachedItemsOnMessageWithMimeTypeUseNumberedPaginateRow
		qe, err = qtx.GetPluralAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate(
			ctx, query.GetPluralAttachedItemsOnMessageWithMimeTypeUseNumberedPaginateParams{
				MessageIds: messageIDs,
				Limit:      int32(np.Limit.Int64),
				Offset:     int32(np.Offset.Int64),
			})
		e = make([]query.GetPluralAttachedItemsOnMessageWithMimeTypeRow, len(qe))
		for i, v := range qe {
			e[i] = query.GetPluralAttachedItemsOnMessageWithMimeTypeRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.AttachedItemOnMessageWithMimeType]{},
			fmt.Errorf("failed to get attached items on message: %w", err)
	}
	entities := make([]entity.AttachedItemOnMessageWithMimeType, len(e))
	for i, v := range e {
		entities[i] = entity.AttachedItemOnMessageWithMimeType{
			AttachedMessageID: v.AttachedMessageID,
			AttachableItem: entity.AttachableItemWithMimeType{
				AttachableItemID: v.AttachableItemID.Bytes,
				OwnerID:          entity.UUID(v.AttachedItemOwnerID),
				FromOuter:        v.AttachedItemFromOuter.Bool,
				URL:              v.AttachedItemUrl.String,
				Size:             entity.Float(v.AttachedItemSize),
				MimeType: entity.NullableEntity[entity.MimeType]{Entity: entity.MimeType{
					MimeTypeID: v.AttachedItemMimeTypeID.Bytes,
					Name:       v.MimeTypeName.String,
					Key:        v.MimeTypeKey.String,
					Kind:       v.MimeTypeKind.String,
				}},
			},
		}
	}
	return store.ListResult[entity.AttachedItemOnMessageWithMimeType]{Data: entities}, nil
}

// GetPluralAttachedItemsOnMessageWithMimeType メッセージに関連付けられた複数の添付アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetPluralAttachedItemsOnMessageWithMimeType(
	ctx context.Context, messageIDs []uuid.UUID,
	order parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	return getPluralAttachedItemsOnMessageWithMimeType(ctx, a.query, messageIDs, order, np)
}

// GetPluralAttachedItemsOnMessageWithMimeTypeWithSd SD付きでメッセージに関連付けられた複数の添付アイテムとそのマイムタイプを取得する。
func (a *PgAdapter) GetPluralAttachedItemsOnMessageWithMimeTypeWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	order parameter.AttachedItemOnMessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.AttachedItemOnMessageWithMimeType], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.AttachedItemOnMessageWithMimeType]{}, store.ErrNotFoundDescriptor
	}
	return getPluralAttachedItemsOnMessageWithMimeType(ctx, qtx, messageIDs, order, np)
}

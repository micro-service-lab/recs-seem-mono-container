package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func convMessage(e query.Message) entity.Message {
	return entity.Message{
		MessageID:    e.MessageID,
		ChatRoomID:   e.ChatRoomID,
		SenderID:     entity.UUID(e.SenderID),
		Body:         e.Body,
		PostedAt:     e.PostedAt,
		LastEditedAt: e.LastEditedAt,
	}
}

func convMessageWithChatRoom(e query.FindMessageByIDWithChatRoomRow) entity.MessageWithChatRoom {
	return entity.MessageWithChatRoom{
		MessageID: e.MessageID,
		ChatRoom: entity.ChatRoomWithCoverImage{
			ChatRoomID:       e.ChatRoomID,
			Name:             e.ChatRoomName.String,
			IsPrivate:        e.ChatRoomIsPrivate.Bool,
			FromOrganization: e.ChatRoomFromOrganization.Bool,
			CoverImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
				Valid: e.ChatRoomCoverImageID.Valid,
				Entity: entity.ImageWithAttachableItem{
					ImageID: e.ChatRoomCoverImageID.Bytes,
					AttachableItem: entity.AttachableItem{
						AttachableItemID: e.ChatRoomCoverImageAttachableItemID.Bytes,
						OwnerID:          entity.UUID(e.ChatRoomCoverImageOwnerID),
						FromOuter:        e.ChatRoomCoverImageFromOuter.Bool,
						URL:              e.ChatRoomCoverImageUrl.String,
						Size:             entity.Float(e.ChatRoomCoverImageSize),
						MimeTypeID:       e.ChatRoomCoverImageMimeTypeID.Bytes,
					},
				},
			},
			OwnerID: entity.UUID(e.ChatRoomOwnerID),
		},
		SenderID:     entity.UUID(e.SenderID),
		Body:         e.Body,
		PostedAt:     e.PostedAt,
		LastEditedAt: e.LastEditedAt,
	}
}

func convMessageWithSender(e query.FindMessageByIDWithSenderRow) entity.MessageWithSender {
	return entity.MessageWithSender{
		MessageID:  e.MessageID,
		ChatRoomID: e.ChatRoomID,
		Sender: entity.NullableEntity[entity.MemberCard]{
			Valid: e.SenderID.Valid,
			Entity: entity.MemberCard{
				MemberID:  e.SenderID.Bytes,
				Name:      e.MemberName.String,
				FirstName: e.MemberFirstName.String,
				LastName:  e.MemberLastName.String,
				Email:     e.MemberEmail.String,
				ProfileImage: entity.NullableEntity[entity.ImageWithAttachableItem]{
					Valid: e.MemberProfileImageID.Valid,
					Entity: entity.ImageWithAttachableItem{
						ImageID: e.MemberProfileImageID.Bytes,
						AttachableItem: entity.AttachableItem{
							AttachableItemID: e.MemberProfileImageAttachableItemID.Bytes,
							OwnerID:          entity.UUID(e.MemberProfileImageOwnerID),
							FromOuter:        e.MemberProfileImageFromOuter.Bool,
							URL:              e.MemberProfileImageUrl.String,
							Size:             entity.Float(e.MemberProfileImageSize),
							MimeTypeID:       e.MemberProfileImageMimeTypeID.Bytes,
						},
					},
				},
			},
		},
		Body:         e.Body,
		PostedAt:     e.PostedAt,
		LastEditedAt: e.LastEditedAt,
	}
}

// countMessages はメッセージ数を取得する内部関数です。
func countMessages(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMessageParam,
) (int64, error) {
	p := query.CountMessagesParams{
		WhereInChatRoom:          where.WhereInChatRoom,
		InChatRoom:               where.InChatRoom,
		WhereInSender:            where.WhereInSender,
		InSender:                 where.InSender,
		WhereLikeBody:            where.WhereLikeBody,
		SearchBody:               where.SearchBody,
		WhereEarlierPostedAt:     where.WhereEarlierPostedAt,
		EarlierPostedAt:          where.EarlierPostedAt,
		WhereLaterPostedAt:       where.WhereLaterPostedAt,
		LaterPostedAt:            where.LaterPostedAt,
		WhereEarlierLastEditedAt: where.WhereEarlierLastEditedAt,
		EarlierLastEditedAt:      where.EarlierLastEditedAt,
		WhereLaterLastEditedAt:   where.WhereLaterLastEditedAt,
		LaterLastEditedAt:        where.LaterLastEditedAt,
	}
	c, err := qtx.CountMessages(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}
	return c, nil
}

// CountMessages はメッセージ数を取得します。
func (a *PgAdapter) CountMessages(ctx context.Context, where parameter.WhereMessageParam) (int64, error) {
	return countMessages(ctx, a.query, where)
}

// CountMessagesWithSd はSD付きでメッセージ数を取得します。
func (a *PgAdapter) CountMessagesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMessageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countMessages(ctx, qtx, where)
}

// createMessage はメッセージを作成する内部関数です。
func createMessage(
	ctx context.Context, qtx *query.Queries, param parameter.CreateMessageParam,
) (entity.Message, error) {
	p := query.CreateMessageParams{
		ChatRoomID:   param.ChatRoomID,
		SenderID:     pgtype.UUID(param.SenderID),
		Body:         param.Body,
		PostedAt:     param.PostedAt,
		LastEditedAt: param.PostedAt,
	}
	e, err := qtx.CreateMessage(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Message{}, errhandle.NewModelDuplicatedError("message")
		}
		return entity.Message{}, fmt.Errorf("failed to create message: %w", err)
	}
	return convMessage(e), nil
}

// CreateMessage はメッセージを作成します。
func (a *PgAdapter) CreateMessage(
	ctx context.Context, param parameter.CreateMessageParam,
) (entity.Message, error) {
	return createMessage(ctx, a.query, param)
}

// CreateMessageWithSd はSD付きでメッセージを作成します。
func (a *PgAdapter) CreateMessageWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateMessageParam,
) (entity.Message, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Message{}, store.ErrNotFoundDescriptor
	}
	return createMessage(ctx, qtx, param)
}

// createMessages は複数のメッセージを作成する内部関数です。
func createMessages(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateMessageParam,
) (int64, error) {
	param := make([]query.CreateMessagesParams, len(params))
	for i, p := range params {
		param[i] = query.CreateMessagesParams{
			ChatRoomID:   p.ChatRoomID,
			SenderID:     pgtype.UUID(p.SenderID),
			Body:         p.Body,
			PostedAt:     p.PostedAt,
			LastEditedAt: p.PostedAt,
		}
	}
	n, err := qtx.CreateMessages(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("message")
		}
		return 0, fmt.Errorf("failed to create messages: %w", err)
	}
	return n, nil
}

// CreateMessages は複数のメッセージを作成します。
func (a *PgAdapter) CreateMessages(
	ctx context.Context, params []parameter.CreateMessageParam,
) (int64, error) {
	return createMessages(ctx, a.query, params)
}

// CreateMessagesWithSd はSD付きで複数のメッセージを作成します。
func (a *PgAdapter) CreateMessagesWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateMessageParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createMessages(ctx, qtx, params)
}

// deleteMessage はメッセージを削除する内部関数です。
func deleteMessage(ctx context.Context, qtx *query.Queries, messageID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteMessage(ctx, messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete message: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("message")
	}
	return c, nil
}

// DeleteMessage はメッセージを削除します。
func (a *PgAdapter) DeleteMessage(ctx context.Context, messageID uuid.UUID) (int64, error) {
	return deleteMessage(ctx, a.query, messageID)
}

// DeleteMessageWithSd はSD付きでメッセージを削除します。
func (a *PgAdapter) DeleteMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteMessage(ctx, qtx, messageID)
}

func deleteMessagesOnChatRoom(ctx context.Context, qtx *query.Queries, chatRoomID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteMessagesOnChatRoom(ctx, chatRoomID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete messages on chat room: %w", err)
	}
	return c, nil
}

// DeleteMessagesOnChatRoom はチャットルーム内のメッセージを削除します。
func (a *PgAdapter) DeleteMessagesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error) {
	return deleteMessagesOnChatRoom(ctx, a.query, chatRoomID)
}

// DeleteMessagesOnChatRoomWithSd はSD付きでチャットルーム内のメッセージを削除します。
func (a *PgAdapter) DeleteMessagesOnChatRoomWithSd(
	ctx context.Context, sd store.Sd, chatRoomID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteMessagesOnChatRoom(ctx, qtx, chatRoomID)
}

// pluralDeleteMessages は複数のメッセージを削除する内部関数です。
func pluralDeleteMessages(ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteMessages(ctx, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete messages: %w", err)
	}
	if c != int64(len(messageIDs)) {
		return 0, errhandle.NewModelNotFoundError("message")
	}
	return c, nil
}

// PluralDeleteMessages は複数のメッセージを削除します。
func (a *PgAdapter) PluralDeleteMessages(ctx context.Context, messageIDs []uuid.UUID) (int64, error) {
	return pluralDeleteMessages(ctx, a.query, messageIDs)
}

// PluralDeleteMessagesWithSd はSD付きで複数のメッセージを削除します。
func (a *PgAdapter) PluralDeleteMessagesWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteMessages(ctx, qtx, messageIDs)
}

// findMessageByID はメッセージをIDで取得する内部関数です。
func findMessageByID(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID,
) (entity.Message, error) {
	e, err := qtx.FindMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Message{}, errhandle.NewModelNotFoundError("message")
		}
		return entity.Message{}, fmt.Errorf("failed to find message: %w", err)
	}
	return convMessage(e), nil
}

// FindMessageByID はメッセージをIDで取得します。
func (a *PgAdapter) FindMessageByID(ctx context.Context, messageID uuid.UUID) (entity.Message, error) {
	return findMessageByID(ctx, a.query, messageID)
}

// FindMessageByIDWithSd はSD付きでメッセージをIDで取得します。
func (a *PgAdapter) FindMessageByIDWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
) (entity.Message, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Message{}, store.ErrNotFoundDescriptor
	}
	return findMessageByID(ctx, qtx, messageID)
}

// findMessageWithChatRoom はメッセージと出席状況を取得する内部関数です。
func findMessageWithChatRoom(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID,
) (entity.MessageWithChatRoom, error) {
	e, err := qtx.FindMessageByIDWithChatRoom(ctx, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MessageWithChatRoom{}, errhandle.NewModelNotFoundError("message")
		}
		return entity.MessageWithChatRoom{}, fmt.Errorf("failed to find message with chat room: %w", err)
	}
	return convMessageWithChatRoom(e), nil
}

// FindMessageWithChatRoom はメッセージと出席状況を取得します。
func (a *PgAdapter) FindMessageWithChatRoom(
	ctx context.Context, messageID uuid.UUID,
) (entity.MessageWithChatRoom, error) {
	return findMessageWithChatRoom(ctx, a.query, messageID)
}

// FindMessageWithChatRoomWithSd はSD付きでメッセージと出席状況を取得します。
func (a *PgAdapter) FindMessageWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
) (entity.MessageWithChatRoom, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MessageWithChatRoom{}, store.ErrNotFoundDescriptor
	}
	return findMessageWithChatRoom(ctx, qtx, messageID)
}

// findMessageWithSender はメッセージとプロフィール画像を取得する内部関数です。
func findMessageWithSender(
	ctx context.Context, qtx *query.Queries, messageID uuid.UUID,
) (entity.MessageWithSender, error) {
	e, err := qtx.FindMessageByIDWithSender(ctx, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.MessageWithSender{}, errhandle.NewModelNotFoundError("message")
		}
		return entity.MessageWithSender{}, fmt.Errorf("failed to find message with sender: %w", err)
	}
	return convMessageWithSender(e), nil
}

// FindMessageWithSender はメッセージとプロフィール画像を取得します。
func (a *PgAdapter) FindMessageWithSender(
	ctx context.Context, messageID uuid.UUID,
) (entity.MessageWithSender, error) {
	return findMessageWithSender(ctx, a.query, messageID)
}

// FindMessageWithSenderWithSd はSD付きでメッセージとプロフィール画像を取得します。
func (a *PgAdapter) FindMessageWithSenderWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID,
) (entity.MessageWithSender, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.MessageWithSender{}, store.ErrNotFoundDescriptor
	}
	return findMessageWithSender(ctx, qtx, messageID)
}

func convCountMessagesParams(p parameter.WhereMessageParam) query.CountMessagesParams {
	return query.CountMessagesParams{
		WhereInChatRoom:          p.WhereInChatRoom,
		InChatRoom:               p.InChatRoom,
		WhereInSender:            p.WhereInSender,
		InSender:                 p.InSender,
		WhereLikeBody:            p.WhereLikeBody,
		SearchBody:               p.SearchBody,
		WhereEarlierPostedAt:     p.WhereEarlierPostedAt,
		EarlierPostedAt:          p.EarlierPostedAt,
		WhereLaterPostedAt:       p.WhereLaterPostedAt,
		LaterPostedAt:            p.LaterPostedAt,
		WhereEarlierLastEditedAt: p.WhereEarlierLastEditedAt,
		EarlierLastEditedAt:      p.EarlierLastEditedAt,
		WhereLaterLastEditedAt:   p.WhereLaterLastEditedAt,
		LaterLastEditedAt:        p.LaterLastEditedAt,
	}
}

func convGetMessagesParams(p parameter.WhereMessageParam, orderMethod string) query.GetMessagesParams {
	return query.GetMessagesParams{
		WhereInChatRoom:          p.WhereInChatRoom,
		InChatRoom:               p.InChatRoom,
		WhereInSender:            p.WhereInSender,
		InSender:                 p.InSender,
		WhereLikeBody:            p.WhereLikeBody,
		SearchBody:               p.SearchBody,
		WhereEarlierPostedAt:     p.WhereEarlierPostedAt,
		EarlierPostedAt:          p.EarlierPostedAt,
		WhereLaterPostedAt:       p.WhereLaterPostedAt,
		LaterPostedAt:            p.LaterPostedAt,
		WhereEarlierLastEditedAt: p.WhereEarlierLastEditedAt,
		EarlierLastEditedAt:      p.EarlierLastEditedAt,
		WhereLaterLastEditedAt:   p.WhereLaterLastEditedAt,
		LaterLastEditedAt:        p.LaterLastEditedAt,
		OrderMethod:              orderMethod,
	}
}

func convGetMessagesUseKeysetPaginateParams(p parameter.WhereMessageParam,
	subCursor, orderMethod string,
	limit int32, cursorDir string, cursor int32, subCursorValue any,
) query.GetMessagesUseKeysetPaginateParams {
	var postCursor, editCursor time.Time
	var err error
	switch subCursor {
	case parameter.MessagePostedAtCursorKey:
		cv, ok := subCursorValue.(string)
		postCursor, err = time.Parse(time.RFC3339, cv)
		if !ok || err != nil {
			postCursor = time.Time{}
		}
	case parameter.MessageLastEditedAtCursorKey:
		cv, ok := subCursorValue.(string)
		editCursor, err = time.Parse(time.RFC3339, cv)
		if !ok || err != nil {
			editCursor = time.Time{}
		}
	}
	return query.GetMessagesUseKeysetPaginateParams{
		WhereInChatRoom:          p.WhereInChatRoom,
		InChatRoom:               p.InChatRoom,
		WhereInSender:            p.WhereInSender,
		InSender:                 p.InSender,
		WhereLikeBody:            p.WhereLikeBody,
		SearchBody:               p.SearchBody,
		WhereEarlierPostedAt:     p.WhereEarlierPostedAt,
		EarlierPostedAt:          p.EarlierPostedAt,
		WhereLaterPostedAt:       p.WhereLaterPostedAt,
		LaterPostedAt:            p.LaterPostedAt,
		WhereEarlierLastEditedAt: p.WhereEarlierLastEditedAt,
		EarlierLastEditedAt:      p.EarlierLastEditedAt,
		WhereLaterLastEditedAt:   p.WhereLaterLastEditedAt,
		LaterLastEditedAt:        p.LaterLastEditedAt,
		OrderMethod:              orderMethod,
		Limit:                    limit,
		CursorDirection:          cursorDir,
		Cursor:                   cursor,
		PostedAtCursor:           postCursor,
		LastEditedAtCursor:       editCursor,
	}
}

func convGetMessagesUseNumberedPaginateParams(
	p parameter.WhereMessageParam, orderMethod string, limit, offset int32,
) query.GetMessagesUseNumberedPaginateParams {
	return query.GetMessagesUseNumberedPaginateParams{
		WhereInChatRoom:          p.WhereInChatRoom,
		InChatRoom:               p.InChatRoom,
		WhereInSender:            p.WhereInSender,
		InSender:                 p.InSender,
		WhereLikeBody:            p.WhereLikeBody,
		SearchBody:               p.SearchBody,
		WhereEarlierPostedAt:     p.WhereEarlierPostedAt,
		EarlierPostedAt:          p.EarlierPostedAt,
		WhereLaterPostedAt:       p.WhereLaterPostedAt,
		LaterPostedAt:            p.LaterPostedAt,
		WhereEarlierLastEditedAt: p.WhereEarlierLastEditedAt,
		EarlierLastEditedAt:      p.EarlierLastEditedAt,
		WhereLaterLastEditedAt:   p.WhereLaterLastEditedAt,
		LaterLastEditedAt:        p.LaterLastEditedAt,
		OrderMethod:              orderMethod,
		Limit:                    limit,
		Offset:                   offset,
	}
}

// getMessages はメッセージを取得する内部関数です。
func getMessages(
	ctx context.Context,
	qtx *query.Queries,
	where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Message], error) {
	eConvFunc := func(e query.Message) (entity.Message, error) {
		return convMessage(e), nil
	}
	runCFunc := func() (int64, error) {
		p := convCountMessagesParams(where)
		r, err := qtx.CountMessages(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to count messages: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.Message, error) {
		p := convGetMessagesParams(where, orderMethod)
		r, err := qtx.GetMessages(ctx, p)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Message{}, nil
			}
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.Message, error) {
		r, err := qtx.GetMessagesUseKeysetPaginate(ctx, convGetMessagesUseKeysetPaginateParams(
			where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.Message, error) {
		r, err := qtx.GetMessagesUseNumberedPaginate(ctx, convGetMessagesUseNumberedPaginateParams(
			where, orderMethod, limit, offset,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Message) (entity.Int, any) {
		switch subCursor {
		case parameter.MessageDefaultCursorKey:
			return entity.Int(e.TMessagesPkey), nil
		case parameter.MessagePostedAtCursorKey:
			return entity.Int(e.TMessagesPkey), e.PostedAt
		case parameter.MessageLastEditedAtCursorKey:
			return entity.Int(e.TMessagesPkey), e.LastEditedAt
		}
		return entity.Int(e.TMessagesPkey), nil
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
		return store.ListResult[entity.Message]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	return res, nil
}

// GetMessages はメッセージを取得します。
func (a *PgAdapter) GetMessages(
	ctx context.Context, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Message], error) {
	return getMessages(ctx, a.query, where, order, np, cp, wc)
}

// GetMessagesWithSd はSD付きでメッセージを取得します。
func (a *PgAdapter) GetMessagesWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Message], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Message]{}, store.ErrNotFoundDescriptor
	}
	return getMessages(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMessages は複数のメッセージを取得する内部関数です。
func getPluralMessages(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	orderMethod parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Message], error) {
	var e []query.Message
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMessages(ctx, query.GetPluralMessagesParams{
			MessageIds:  messageIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralMessagesUseNumberedPaginate(ctx, query.GetPluralMessagesUseNumberedPaginateParams{
			MessageIds:  messageIDs,
			Offset:      int32(np.Offset.Int64),
			Limit:       int32(np.Limit.Int64),
			OrderMethod: orderMethod.GetStringValue(),
		})
	}
	if err != nil {
		return store.ListResult[entity.Message]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	entities := make([]entity.Message, len(e))
	for i, v := range e {
		entities[i] = convMessage(v)
	}
	return store.ListResult[entity.Message]{Data: entities}, nil
}

// GetPluralMessages は複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessages(
	ctx context.Context, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Message], error) {
	return getPluralMessages(ctx, a.query, messageIDs, order, np)
}

// GetPluralMessagesWithSd はSD付きで複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessagesWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Message], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Message]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMessages(ctx, qtx, messageIDs, order, np)
}

func getMessagesWithChatRoom(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	eConvFunc := func(e entity.MessageWithChatRoomForQuery) (entity.MessageWithChatRoom, error) {
		return e.MessageWithChatRoom, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMessages(ctx, convCountMessagesParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count messages: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MessageWithChatRoomForQuery, error) {
		r, err := qtx.GetMessagesWithChatRoom(
			ctx, query.GetMessagesWithChatRoomParams(convGetMessagesParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MessageWithChatRoomForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithChatRoomForQuery{
				Pkey:                entity.Int(v.TMessagesPkey),
				MessageWithChatRoom: convMessageWithChatRoom(query.FindMessageByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MessageWithChatRoomForQuery, error) {
		r, err := qtx.GetMessagesWithChatRoomUseKeysetPaginate(
			ctx, query.GetMessagesWithChatRoomUseKeysetPaginateParams(convGetMessagesUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithChatRoomForQuery{
				Pkey:                entity.Int(v.TMessagesPkey),
				MessageWithChatRoom: convMessageWithChatRoom(query.FindMessageByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MessageWithChatRoomForQuery, error) {
		r, err := qtx.GetMessagesWithChatRoomUseNumberedPaginate(
			ctx, query.GetMessagesWithChatRoomUseNumberedPaginateParams(convGetMessagesUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithChatRoomForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithChatRoomForQuery{
				Pkey:                entity.Int(v.TMessagesPkey),
				MessageWithChatRoom: convMessageWithChatRoom(query.FindMessageByIDWithChatRoomRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MessageWithChatRoomForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MessageDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MessagePostedAtCursorKey:
			return entity.Int(e.Pkey), e.PostedAt
		case parameter.MessageLastEditedAtCursorKey:
			return entity.Int(e.Pkey), e.LastEditedAt
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
		return store.ListResult[entity.MessageWithChatRoom]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	return res, nil
}

// GetMessagesWithChatRoom はメッセージとチャットルームを取得します。
func (a *PgAdapter) GetMessagesWithChatRoom(
	ctx context.Context, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	return getMessagesWithChatRoom(ctx, a.query, where, order, np, cp, wc)
}

// GetMessagesWithChatRoomWithSd はSD付きでメッセージとチャットルームを取得します。
func (a *PgAdapter) GetMessagesWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MessageWithChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getMessagesWithChatRoom(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMessagesWithChatRoom は複数のメッセージを取得する内部関数です。
func getPluralMessagesWithChatRoom(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	orderMethod parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	var e []query.GetPluralMessagesWithChatRoomRow
	var te []query.GetPluralMessagesWithChatRoomUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMessagesWithChatRoom(ctx, query.GetPluralMessagesWithChatRoomParams{
			MessageIds:  messageIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMessagesWithChatRoomUseNumberedPaginate(
			ctx, query.GetPluralMessagesWithChatRoomUseNumberedPaginateParams{
				MessageIds:  messageIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMessagesWithChatRoomRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMessagesWithChatRoomRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MessageWithChatRoom]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	entities := make([]entity.MessageWithChatRoom, len(e))
	for i, v := range e {
		entities[i] = convMessageWithChatRoom(query.FindMessageByIDWithChatRoomRow(v))
	}
	return store.ListResult[entity.MessageWithChatRoom]{Data: entities}, nil
}

// GetPluralMessagesWithChatRoom は複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessagesWithChatRoom(
	ctx context.Context, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	return getPluralMessagesWithChatRoom(ctx, a.query, messageIDs, order, np)
}

// GetPluralMessagesWithChatRoomWithSd はSD付きで複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessagesWithChatRoomWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithChatRoom], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MessageWithChatRoom]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMessagesWithChatRoom(ctx, qtx, messageIDs, order, np)
}

func getMessagesWithSender(
	ctx context.Context, qtx *query.Queries, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithSender], error) {
	eConvFunc := func(e entity.MessageWithSenderForQuery) (entity.MessageWithSender, error) {
		return e.MessageWithSender, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountMessages(ctx, convCountMessagesParams(where))
		if err != nil {
			return 0, fmt.Errorf("failed to count messages: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.MessageWithSenderForQuery, error) {
		r, err := qtx.GetMessagesWithSender(
			ctx, query.GetMessagesWithSenderParams(convGetMessagesParams(where, orderMethod)))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.MessageWithSenderForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithSenderForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithSenderForQuery{
				Pkey:              entity.Int(v.TMessagesPkey),
				MessageWithSender: convMessageWithSender(query.FindMessageByIDWithSenderRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.MessageWithSenderForQuery, error) {
		r, err := qtx.GetMessagesWithSenderUseKeysetPaginate(
			ctx, query.GetMessagesWithSenderUseKeysetPaginateParams(convGetMessagesUseKeysetPaginateParams(
				where, subCursor, orderMethod, limit, cursorDir, cursor, subCursorValue,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithSenderForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithSenderForQuery{
				Pkey:              entity.Int(v.TMessagesPkey),
				MessageWithSender: convMessageWithSender(query.FindMessageByIDWithSenderRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.MessageWithSenderForQuery, error) {
		r, err := qtx.GetMessagesWithSenderUseNumberedPaginate(
			ctx, query.GetMessagesWithSenderUseNumberedPaginateParams(convGetMessagesUseNumberedPaginateParams(
				where, orderMethod, limit, offset,
			)))
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}
		e := make([]entity.MessageWithSenderForQuery, len(r))
		for i, v := range r {
			e[i] = entity.MessageWithSenderForQuery{
				Pkey:              entity.Int(v.TMessagesPkey),
				MessageWithSender: convMessageWithSender(query.FindMessageByIDWithSenderRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.MessageWithSenderForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.MessageDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		case parameter.MessagePostedAtCursorKey:
			return entity.Int(e.Pkey), e.PostedAt
		case parameter.MessageLastEditedAtCursorKey:
			return entity.Int(e.Pkey), e.LastEditedAt
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
		return store.ListResult[entity.MessageWithSender]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	return res, nil
}

// GetMessagesWithSender はメッセージとチャットルームを取得します。
func (a *PgAdapter) GetMessagesWithSender(
	ctx context.Context, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithSender], error) {
	return getMessagesWithSender(ctx, a.query, where, order, np, cp, wc)
}

// GetMessagesWithSenderWithSd はSD付きでメッセージとチャットルームを取得します。
func (a *PgAdapter) GetMessagesWithSenderWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereMessageParam,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.MessageWithSender], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MessageWithSender]{}, store.ErrNotFoundDescriptor
	}
	return getMessagesWithSender(ctx, qtx, where, order, np, cp, wc)
}

// getPluralMessagesWithSender は複数のメッセージを取得する内部関数です。
func getPluralMessagesWithSender(
	ctx context.Context, qtx *query.Queries, messageIDs []uuid.UUID,
	orderMethod parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithSender], error) {
	var e []query.GetPluralMessagesWithSenderRow
	var te []query.GetPluralMessagesWithSenderUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralMessagesWithSender(ctx, query.GetPluralMessagesWithSenderParams{
			MessageIds:  messageIDs,
			OrderMethod: orderMethod.GetStringValue(),
		})
	} else {
		te, err = qtx.GetPluralMessagesWithSenderUseNumberedPaginate(
			ctx, query.GetPluralMessagesWithSenderUseNumberedPaginateParams{
				MessageIds:  messageIDs,
				Offset:      int32(np.Offset.Int64),
				Limit:       int32(np.Limit.Int64),
				OrderMethod: orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralMessagesWithSenderRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralMessagesWithSenderRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.MessageWithSender]{}, fmt.Errorf("failed to get messages: %w", err)
	}
	entities := make([]entity.MessageWithSender, len(e))
	for i, v := range e {
		entities[i] = convMessageWithSender(query.FindMessageByIDWithSenderRow(v))
	}
	return store.ListResult[entity.MessageWithSender]{Data: entities}, nil
}

// GetPluralMessagesWithSender は複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessagesWithSender(
	ctx context.Context, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithSender], error) {
	return getPluralMessagesWithSender(ctx, a.query, messageIDs, order, np)
}

// GetPluralMessagesWithSenderWithSd はSD付きで複数のメッセージを取得します。
func (a *PgAdapter) GetPluralMessagesWithSenderWithSd(
	ctx context.Context, sd store.Sd, messageIDs []uuid.UUID,
	order parameter.MessageOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.MessageWithSender], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.MessageWithSender]{}, store.ErrNotFoundDescriptor
	}
	return getPluralMessagesWithSender(ctx, qtx, messageIDs, order, np)
}

func updateMessage(
	ctx context.Context, qtx *query.Queries,
	messageID uuid.UUID, param parameter.UpdateMessageParams,
) (entity.Message, error) {
	p := query.UpdateMessageParams{
		MessageID:    messageID,
		Body:         param.Body,
		LastEditedAt: param.LastEditedAt,
	}
	e, err := qtx.UpdateMessage(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Message{}, errhandle.NewModelNotFoundError("message")
		}
		return entity.Message{}, fmt.Errorf("failed to update message: %w", err)
	}
	return convMessage(query.Message(e)), nil
}

// UpdateMessage はメッセージを更新します。
func (a *PgAdapter) UpdateMessage(
	ctx context.Context, messageID uuid.UUID, param parameter.UpdateMessageParams,
) (entity.Message, error) {
	return updateMessage(ctx, a.query, messageID, param)
}

// UpdateMessageWithSd はSD付きでメッセージを更新します。
func (a *PgAdapter) UpdateMessageWithSd(
	ctx context.Context, sd store.Sd, messageID uuid.UUID, param parameter.UpdateMessageParams,
) (entity.Message, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Message{}, store.ErrNotFoundDescriptor
	}
	return updateMessage(ctx, qtx, messageID, param)
}

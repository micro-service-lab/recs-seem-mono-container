package pgadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convChatRoomActionOnChatRoom(
	e query.GetChatRoomActionsOnChatRoomRow,
) entity.ChatRoomActionWithDetailForQuery {
	var createAction entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]
	var updateNameAction entity.NullableEntity[entity.ChatRoomUpdateNameActionWithUpdatedBy]
	var addMemberAction entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedBy]
	var removeMemberAction entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedBy]
	var withdrawAction entity.NullableEntity[entity.ChatRoomWithdrawActionWithMember]
	var message entity.NullableEntity[entity.MessageWithSender]

	if e.ChatRoomCreateActionID.Valid {
		var createdBy entity.NullableEntity[entity.SimpleMember]
		if e.CreateMemberID.Valid {
			createdBy = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.CreateMemberID.Bytes,
					Name:           e.CreateMemberName.String,
					Email:          e.CreateMemberEmail.String,
					FirstName:      entity.String(e.CreateMemberFirstName),
					LastName:       entity.String(e.CreateMemberLastName),
					ProfileImageID: entity.UUID(e.CreateMemberProfileImageID),
				},
			}
		}
		createAction = entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]{
			Valid: true,
			Entity: entity.ChatRoomCreateActionWithCreatedBy{
				ChatRoomCreateActionID: e.ChatRoomCreateActionID.Bytes,
				ChatRoomActionID:       e.ChatRoomActionID,
				Name:                   e.CreateName.String,
				CreatedBy:              createdBy,
			},
		}
	}

	if e.ChatRoomUpdateNameActionID.Valid {
		var updatedBy entity.NullableEntity[entity.SimpleMember]
		if e.UpdateMemberID.Valid {
			updatedBy = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.UpdateMemberID.Bytes,
					Name:           e.UpdateMemberName.String,
					Email:          e.UpdateMemberEmail.String,
					FirstName:      entity.String(e.UpdateMemberFirstName),
					LastName:       entity.String(e.UpdateMemberLastName),
					ProfileImageID: entity.UUID(e.UpdateMemberProfileImageID),
				},
			}
		}
		updateNameAction = entity.NullableEntity[entity.ChatRoomUpdateNameActionWithUpdatedBy]{
			Valid: true,
			Entity: entity.ChatRoomUpdateNameActionWithUpdatedBy{
				ChatRoomUpdateNameActionID: e.ChatRoomUpdateNameActionID.Bytes,
				ChatRoomActionID:           e.ChatRoomActionID,
				Name:                       e.UpdateName.String,
				UpdatedBy:                  updatedBy,
			},
		}
	}

	if e.ChatRoomAddMemberActionID.Valid {
		var addedBy entity.NullableEntity[entity.SimpleMember]
		if e.AddMemberID.Valid {
			addedBy = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.AddMemberID.Bytes,
					Name:           e.AddMemberName.String,
					Email:          e.AddMemberEmail.String,
					FirstName:      entity.String(e.AddMemberFirstName),
					LastName:       entity.String(e.AddMemberLastName),
					ProfileImageID: entity.UUID(e.AddMemberProfileImageID),
				},
			}
		}
		addMemberAction = entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedBy]{
			Valid: true,
			Entity: entity.ChatRoomAddMemberActionWithAddedBy{
				ChatRoomAddMemberActionID: e.ChatRoomAddMemberActionID.Bytes,
				ChatRoomActionID:          e.ChatRoomActionID,
				AddedBy:                   addedBy,
			},
		}
	}

	if e.ChatRoomRemoveMemberActionID.Valid {
		var removedBy entity.NullableEntity[entity.SimpleMember]
		if e.RemoveMemberID.Valid {
			removedBy = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.RemoveMemberID.Bytes,
					Name:           e.RemoveMemberName.String,
					Email:          e.RemoveMemberEmail.String,
					FirstName:      entity.String(e.RemoveMemberFirstName),
					LastName:       entity.String(e.RemoveMemberLastName),
					ProfileImageID: entity.UUID(e.RemoveMemberProfileImageID),
				},
			}
		}
		removeMemberAction = entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedBy]{
			Valid: true,
			Entity: entity.ChatRoomRemoveMemberActionWithRemovedBy{
				ChatRoomRemoveMemberActionID: e.ChatRoomRemoveMemberActionID.Bytes,
				ChatRoomActionID:             e.ChatRoomActionID,
				RemovedBy:                    removedBy,
			},
		}
	}

	if e.ChatRoomWithdrawActionID.Valid {
		var member entity.NullableEntity[entity.SimpleMember]
		if e.WithdrawMemberID.Valid {
			member = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.WithdrawMemberID.Bytes,
					Name:           e.WithdrawMemberName.String,
					Email:          e.WithdrawMemberEmail.String,
					FirstName:      entity.String(e.WithdrawMemberFirstName),
					LastName:       entity.String(e.WithdrawMemberLastName),
					ProfileImageID: entity.UUID(e.WithdrawMemberProfileImageID),
				},
			}
		}
		withdrawAction = entity.NullableEntity[entity.ChatRoomWithdrawActionWithMember]{
			Valid: true,
			Entity: entity.ChatRoomWithdrawActionWithMember{
				ChatRoomWithdrawActionID: e.ChatRoomWithdrawActionID.Bytes,
				ChatRoomActionID:         e.ChatRoomActionID,
				Member:                   member,
			},
		}
	}

	if e.MessageID.Valid {
		var sender entity.NullableEntity[entity.MemberCard]
		if e.MessageSenderID.Valid {
			var profileImage entity.NullableEntity[entity.ImageWithAttachableItem]
			if e.MessageSenderProfileImageID.Valid {
				profileImage = entity.NullableEntity[entity.ImageWithAttachableItem]{
					Valid: true,
					Entity: entity.ImageWithAttachableItem{
						ImageID: e.MessageSenderProfileImageID.Bytes,
						AttachableItem: entity.AttachableItem{
							AttachableItemID: e.MessageSenderProfileImageAttachableItemID.Bytes,
							OwnerID:          entity.UUID(e.MessageSenderProfileImageOwnerID),
							FromOuter:        e.MessageSenderProfileImageFromOuter.Bool,
							URL:              e.MessageSenderProfileImageUrl.String,
							Alias:            e.MessageSenderProfileImageAlias.String,
							Size:             entity.Float(e.MessageSenderProfileImageSize),
							MimeTypeID:       e.MessageSenderProfileImageMimeTypeID.Bytes,
						},
					},
				}
			}
			sender = entity.NullableEntity[entity.MemberCard]{
				Valid: true,
				Entity: entity.MemberCard{
					MemberID:     e.MessageSenderID.Bytes,
					Name:         e.MessageSenderName.String,
					Email:        e.MessageSenderEmail.String,
					FirstName:    entity.String(e.MessageSenderFirstName),
					LastName:     entity.String(e.MessageSenderLastName),
					ProfileImage: profileImage,
				},
			}
		}
		message = entity.NullableEntity[entity.MessageWithSender]{
			Valid: true,
			Entity: entity.MessageWithSender{
				MessageID:        e.MessageID.Bytes,
				ChatRoomActionID: e.ChatRoomActionID,
				Sender:           sender,
				Body:             e.MessageBody.String,
				PostedAt:         e.MessagePostedAt.Time,
				LastEditedAt:     e.MessageLastEditedAt.Time,
			},
		}
	}

	return entity.ChatRoomActionWithDetailForQuery{
		Pkey: entity.Int(e.TChatRoomActionsPkey),
		ChatRoomActionWithDetail: entity.ChatRoomActionWithDetail{
			ChatRoomActionID:           e.ChatRoomActionID,
			ChatRoomID:                 e.ChatRoomID,
			ChatRoomActionTypeID:       e.ChatRoomActionTypeID,
			ActedAt:                    e.ActedAt,
			ChatRoomCreateAction:       createAction,
			ChatRoomUpdateNameAction:   updateNameAction,
			ChatRoomAddMemberAction:    addMemberAction,
			ChatRoomRemoveMemberAction: removeMemberAction,
			ChatRoomWithdrawAction:     withdrawAction,
			Message:                    message,
		},
	}
}

// countChatRoomActions はチャットルームアクション数を取得する内部関数です。
func countChatRoomActions(
	ctx context.Context, qtx *query.Queries, where parameter.WhereChatRoomActionParam,
) (int64, error) {
	c, err := qtx.CountChatRoomActions(ctx, query.CountChatRoomActionsParams{
		WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
		InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count chat room actions: %w", err)
	}
	return c, nil
}

// CountChatRoomActions はチャットルームアクション数を取得します。
func (a *PgAdapter) CountChatRoomActions(
	ctx context.Context, where parameter.WhereChatRoomActionParam,
) (int64, error) {
	return countChatRoomActions(ctx, a.query, where)
}

// CountChatRoomActionsWithSd はSD付きでチャットルームアクション数を取得します。
func (a *PgAdapter) CountChatRoomActionsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereChatRoomActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countChatRoomActions(ctx, qtx, where)
}

// createChatRoomAction はチャットルームアクションを作成する内部関数です。
func createChatRoomAction(
	ctx context.Context, qtx *query.Queries, param parameter.CreateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	e, err := qtx.CreateChatRoomAction(ctx, query.CreateChatRoomActionParams{
		ChatRoomID:           param.ChatRoomID,
		ChatRoomActionTypeID: param.ChatRoomActionTypeID,
		ActedAt:              param.ActedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.ChatRoomAction{}, errhandle.NewModelDuplicatedError("chat room action")
		}
		return entity.ChatRoomAction{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	entity := entity.ChatRoomAction{
		ChatRoomActionID:     e.ChatRoomActionID,
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		ActedAt:              e.ActedAt,
	}
	return entity, nil
}

// CreateChatRoomAction はチャットルームアクションを作成します。
func (a *PgAdapter) CreateChatRoomAction(
	ctx context.Context, param parameter.CreateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	return createChatRoomAction(ctx, a.query, param)
}

// CreateChatRoomActionWithSd はSD付きでチャットルームアクションを作成します。
func (a *PgAdapter) CreateChatRoomActionWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomAction{}, store.ErrNotFoundDescriptor
	}
	return createChatRoomAction(ctx, qtx, param)
}

// createChatRoomActions は複数のチャットルームアクションを作成する内部関数です。
func createChatRoomActions(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateChatRoomActionParam,
) (int64, error) {
	param := make([]query.CreateChatRoomActionsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateChatRoomActionsParams{
			ChatRoomID:           p.ChatRoomID,
			ChatRoomActionTypeID: p.ChatRoomActionTypeID,
			ActedAt:              p.ActedAt,
		}
	}
	n, err := qtx.CreateChatRoomActions(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("chat room action")
		}
		return 0, fmt.Errorf("failed to create chat room actions: %w", err)
	}
	return n, nil
}

// CreateChatRoomActions は複数のチャットルームアクションを作成します。
func (a *PgAdapter) CreateChatRoomActions(
	ctx context.Context, params []parameter.CreateChatRoomActionParam,
) (int64, error) {
	return createChatRoomActions(ctx, a.query, params)
}

// CreateChatRoomActionsWithSd はSD付きで複数のチャットルームアクションを作成します。
func (a *PgAdapter) CreateChatRoomActionsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateChatRoomActionParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createChatRoomActions(ctx, qtx, params)
}

// deleteChatRoomAction はチャットルームアクションを削除する内部関数です。
func deleteChatRoomAction(
	ctx context.Context, qtx *query.Queries, chatRoomActionID uuid.UUID,
) (int64, error) {
	c, err := qtx.DeleteChatRoomAction(ctx, chatRoomActionID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room action: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("chat room action")
	}
	return c, nil
}

// DeleteChatRoomAction はチャットルームアクションを削除します。
func (a *PgAdapter) DeleteChatRoomAction(
	ctx context.Context, chatRoomActionID uuid.UUID,
) (int64, error) {
	return deleteChatRoomAction(ctx, a.query, chatRoomActionID)
}

// DeleteChatRoomActionWithSd はSD付きでチャットルームアクションを削除します。
func (a *PgAdapter) DeleteChatRoomActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteChatRoomAction(ctx, qtx, chatRoomActionID)
}

// pluralDeleteChatRoomActions は複数のチャットルームアクションを削除する内部関数です。
func pluralDeleteChatRoomActions(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
) (int64, error) {
	c, err := qtx.PluralDeleteChatRoomActions(ctx, chatRoomActionIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room actions: %w", err)
	}
	if c != int64(len(chatRoomActionIDs)) {
		return 0, errhandle.NewModelNotFoundError("chat room action")
	}
	return c, nil
}

// PluralDeleteChatRoomActions は複数のチャットルームアクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomActions(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
) (int64, error) {
	return pluralDeleteChatRoomActions(ctx, a.query, chatRoomActionIDs)
}

// PluralDeleteChatRoomActionsWithSd はSD付きで複数のチャットルームアクションを削除します。
func (a *PgAdapter) PluralDeleteChatRoomActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteChatRoomActions(ctx, qtx, chatRoomActionIDs)
}

// getChatRoomActions はチャットルームアクションを取得する内部関数です。
func getChatRoomActionsOnChatRoom(
	ctx context.Context,
	qtx *query.Queries,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomActionParam,
	order parameter.ChatRoomActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	eConvFunc := func(
		e entity.ChatRoomActionWithDetailForQuery,
	) (entity.ChatRoomActionWithDetail, error) {
		return e.ChatRoomActionWithDetail, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomActions(ctx, query.CountChatRoomActionsParams{
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.ChatRoomActionWithDetailForQuery, error) {
		r, err := qtx.GetChatRoomActionsOnChatRoom(ctx, query.GetChatRoomActionsOnChatRoomParams{
			ChatRoomID:                   chatRoomID,
			OrderMethod:                  orderMethod,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.ChatRoomActionWithDetailForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		e := make([]entity.ChatRoomActionWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomActionOnChatRoom(v)
		}
		return e, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]entity.ChatRoomActionWithDetailForQuery, error) {
		var actCursor time.Time
		var err error
		switch subCursor {
		case parameter.ChatRoomActionActedAtCursorKey:
			cv, ok := subCursorValue.(string)
			actCursor, err = time.Parse(time.RFC3339, cv)
			if !ok || err != nil {
				actCursor = time.Time{}
			}
		}
		p := query.GetChatRoomActionsOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:                   chatRoomID,
			Limit:                        limit,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
			CursorDirection:              cursorDir,
			OrderMethod:                  orderMethod,
			ActedAtCursor:                actCursor,
			Cursor:                       cursor,
		}
		r, err := qtx.GetChatRoomActionsOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		e := make([]entity.ChatRoomActionWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomActionOnChatRoom(query.GetChatRoomActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.ChatRoomActionWithDetailForQuery, error) {
		p := query.GetChatRoomActionsOnChatRoomUseNumberedPaginateParams{
			ChatRoomID:                   chatRoomID,
			Limit:                        limit,
			Offset:                       offset,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
			OrderMethod:                  orderMethod,
		}
		r, err := qtx.GetChatRoomActionsOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		e := make([]entity.ChatRoomActionWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomActionOnChatRoom(query.GetChatRoomActionsOnChatRoomRow(v))
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.ChatRoomActionWithDetailForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomActionDefaultCursorKey:
			return e.Pkey, nil
		case parameter.ChatRoomActionActedAtCursorKey:
			return entity.Int(e.Pkey), e.ActedAt
		}
		return e.Pkey, nil
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
		return store.ListResult[entity.ChatRoomActionWithDetail]{},
			fmt.Errorf("failed to get chat room actions: %w", err)
	}
	return res, nil
}

// GetChatRoomActionsOnChatRoom はチャットルームアクションを取得します。
func (a *PgAdapter) GetChatRoomActionsOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomActionParam,
	order parameter.ChatRoomActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	return getChatRoomActionsOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomActionsOnChatRoomWithSd はSD付きでチャットルームアクションを取得します。
func (a *PgAdapter) GetChatRoomActionsOnChatRoomWithSd(
	ctx context.Context,
	sd store.Sd,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomActionParam,
	order parameter.ChatRoomActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomActionWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomActions は複数のチャットルームアクションを取得する内部関数です。
func getPluralChatRoomActions(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	orderMethod parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	var e []query.GetPluralChatRoomActionsRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomActions(ctx, query.GetPluralChatRoomActionsParams{
			ChatRoomActionIds: chatRoomActionIDs,
			OrderMethod:       orderMethod.GetStringValue(),
		})
	} else {
		var ne []query.GetPluralChatRoomActionsUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomActionsUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
				OrderMethod:       orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralChatRoomActionsRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomActionsRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomActionWithDetail]{},
			fmt.Errorf("failed to get chat room actions: %w", err)
	}
	entities := make([]entity.ChatRoomActionWithDetail, len(e))
	for i, v := range e {
		entities[i] = convChatRoomActionOnChatRoom(
			query.GetChatRoomActionsOnChatRoomRow(v)).ChatRoomActionWithDetail
	}
	return store.ListResult[entity.ChatRoomActionWithDetail]{Data: entities}, nil
}

// GetPluralChatRoomActions は複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActions(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	return getPluralChatRoomActions(ctx, a.query, chatRoomActionIDs, order, np)
}

// GetPluralChatRoomActionsWithSd はSD付きで複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomActionWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomActions(ctx, qtx, chatRoomActionIDs, order, np)
}

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
	e query.ChatRoomAction,
) entity.ChatRoomAction {
	return entity.ChatRoomAction{
		ChatRoomActionID:     e.ChatRoomActionID,
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		ActedAt:              e.ActedAt,
	}
}

func convChatRoomActionWithDetailOnChatRoom(
	e query.GetChatRoomActionsWithDetailOnChatRoomRow,
) entity.ChatRoomActionWithDetailForQuery {
	var createAction entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]
	var updateNameAction entity.NullableEntity[entity.ChatRoomUpdateNameActionWithUpdatedBy]
	var addMemberAction entity.NullableEntity[entity.ChatRoomAddMemberActionWithAddedBy]
	var removeMemberAction entity.NullableEntity[entity.ChatRoomRemoveMemberActionWithRemovedBy]
	var withdrawAction entity.NullableEntity[entity.ChatRoomWithdrawActionWithMember]
	var deleteMessageAction entity.NullableEntity[entity.ChatRoomDeleteMessageActionWithDeletedBy]
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
					GradeID:        e.CreateMemberGradeID.Bytes,
					GroupID:        e.CreateMemberGroupID.Bytes,
				},
			}
		}
		createAction = entity.NullableEntity[entity.ChatRoomCreateActionWithCreatedBy]{
			Valid: true,
			Entity: entity.ChatRoomCreateActionWithCreatedBy{
				ChatRoomCreateActionID: e.ChatRoomCreateActionID.Bytes,
				ChatRoomActionID:       e.ChatRoomActionID,
				Name:                   entity.String(e.CreateName),
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
					GradeID:        e.UpdateMemberGradeID.Bytes,
					GroupID:        e.UpdateMemberGroupID.Bytes,
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
					GradeID:        e.AddMemberGradeID.Bytes,
					GroupID:        e.AddMemberGroupID.Bytes,
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
					GradeID:        e.RemoveMemberGradeID.Bytes,
					GroupID:        e.RemoveMemberGroupID.Bytes,
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
					GradeID:        e.WithdrawMemberGradeID.Bytes,
					GroupID:        e.WithdrawMemberGroupID.Bytes,
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

	if e.ChatRoomDeleteMessageActionID.Valid {
		var deletedBy entity.NullableEntity[entity.SimpleMember]
		if e.DeleteMessageMemberID.Valid {
			deletedBy = entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       e.DeleteMessageMemberID.Bytes,
					Name:           e.DeleteMessageMemberName.String,
					Email:          e.DeleteMessageMemberEmail.String,
					FirstName:      entity.String(e.DeleteMessageMemberFirstName),
					LastName:       entity.String(e.DeleteMessageMemberLastName),
					ProfileImageID: entity.UUID(e.DeleteMessageMemberProfileImageID),
					GradeID:        e.DeleteMessageMemberGradeID.Bytes,
					GroupID:        e.DeleteMessageMemberGroupID.Bytes,
				},
			}
		}
		deleteMessageAction = entity.NullableEntity[entity.ChatRoomDeleteMessageActionWithDeletedBy]{
			Valid: true,
			Entity: entity.ChatRoomDeleteMessageActionWithDeletedBy{
				ChatRoomDeleteMessageActionID: e.ChatRoomDeleteMessageActionID.Bytes,
				ChatRoomActionID:              e.ChatRoomActionID,
				DeletedBy:                     deletedBy,
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
					GradeID:      e.MessageSenderGradeID.Bytes,
					GroupID:      e.MessageSenderGroupID.Bytes,
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
			ChatRoomActionID:            e.ChatRoomActionID,
			ChatRoomID:                  e.ChatRoomID,
			ChatRoomActionTypeID:        e.ChatRoomActionTypeID,
			ActedAt:                     e.ActedAt,
			ChatRoomCreateAction:        createAction,
			ChatRoomUpdateNameAction:    updateNameAction,
			ChatRoomAddMemberAction:     addMemberAction,
			ChatRoomRemoveMemberAction:  removeMemberAction,
			ChatRoomWithdrawAction:      withdrawAction,
			ChatRoomDeleteMessageAction: deleteMessageAction,
			Message:                     message,
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

func updateChatRoomAction(
	ctx context.Context, qtx *query.Queries, chatRoomActionID uuid.UUID, param parameter.UpdateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	e, err := qtx.UpdateChatRoomAction(ctx, query.UpdateChatRoomActionParams{
		ChatRoomActionID:     chatRoomActionID,
		ChatRoomActionTypeID: param.ChatRoomActionTypeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ChatRoomAction{}, errhandle.NewModelNotFoundError("chat room action")
		}
		return entity.ChatRoomAction{}, fmt.Errorf("failed to update chat room action: %w", err)
	}
	entity := entity.ChatRoomAction{
		ChatRoomActionID:     e.ChatRoomActionID,
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: e.ChatRoomActionTypeID,
		ActedAt:              e.ActedAt,
	}
	return entity, nil
}

// UpdateChatRoomAction はチャットルームアクションを更新します。
func (a *PgAdapter) UpdateChatRoomAction(
	ctx context.Context, chatRoomActionID uuid.UUID, param parameter.UpdateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	return updateChatRoomAction(ctx, a.query, chatRoomActionID, param)
}

// UpdateChatRoomActionWithSd はSD付きでチャットルームアクションを更新します。
func (a *PgAdapter) UpdateChatRoomActionWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionID uuid.UUID, param parameter.UpdateChatRoomActionParam,
) (entity.ChatRoomAction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.ChatRoomAction{}, store.ErrNotFoundDescriptor
	}
	return updateChatRoomAction(ctx, qtx, chatRoomActionID, param)
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
) (store.ListResult[entity.ChatRoomAction], error) {
	eConvFunc := func(
		e query.ChatRoomAction,
	) (entity.ChatRoomAction, error) {
		return convChatRoomActionOnChatRoom(e), nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountChatRoomActionsOnChatRoom(ctx, query.CountChatRoomActionsOnChatRoomParams{
			ChatRoomID:                   chatRoomID,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]query.ChatRoomAction, error) {
		r, err := qtx.GetChatRoomActionsOnChatRoom(ctx, query.GetChatRoomActionsOnChatRoomParams{
			ChatRoomID:                   chatRoomID,
			OrderMethod:                  orderMethod,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.ChatRoomAction{}, nil
			}
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(subCursor, orderMethod string,
		limit int32, cursorDir string, cursor int32, subCursorValue any,
	) ([]query.ChatRoomAction, error) {
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
		return r, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]query.ChatRoomAction, error) {
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
		return r, nil
	}
	selector := func(subCursor string, e query.ChatRoomAction) (entity.Int, any) {
		switch subCursor {
		case parameter.ChatRoomActionDefaultCursorKey:
			return entity.Int(e.TChatRoomActionsPkey), nil
		case parameter.ChatRoomActionActedAtCursorKey:
			return entity.Int(e.TChatRoomActionsPkey), e.ActedAt
		}
		return entity.Int(e.TChatRoomActionsPkey), nil
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
		return store.ListResult[entity.ChatRoomAction]{},
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
) (store.ListResult[entity.ChatRoomAction], error) {
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
) (store.ListResult[entity.ChatRoomAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomAction]{}, store.ErrNotFoundDescriptor
	}
	return getChatRoomActionsOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomActions は複数のチャットルームアクションを取得する内部関数です。
func getPluralChatRoomActions(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	orderMethod parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAction], error) {
	var e []query.ChatRoomAction
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomActions(ctx, query.GetPluralChatRoomActionsParams{
			ChatRoomActionIds: chatRoomActionIDs,
			OrderMethod:       orderMethod.GetStringValue(),
		})
	} else {
		e, err = qtx.GetPluralChatRoomActionsUseNumberedPaginate(
			ctx, query.GetPluralChatRoomActionsUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
				OrderMethod:       orderMethod.GetStringValue(),
			})
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomAction]{},
			fmt.Errorf("failed to get chat room actions: %w", err)
	}
	entities := make([]entity.ChatRoomAction, len(e))
	for i, v := range e {
		entities[i] = convChatRoomActionOnChatRoom(v)
	}
	return store.ListResult[entity.ChatRoomAction]{Data: entities}, nil
}

// GetPluralChatRoomActions は複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActions(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	orderMethod parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAction], error) {
	return getPluralChatRoomActions(ctx, a.query, chatRoomActionIDs, orderMethod, np)
}

// GetPluralChatRoomActionsWithSd はSD付きで複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActionsWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	orderMethod parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomAction], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomAction]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomActions(ctx, qtx, chatRoomActionIDs, orderMethod, np)
}

// getChatRoomActionsWithDetailOnChatRoom はチャットルームアクション詳細を取得する内部関数です。
func getChatRoomActionsWithDetailOnChatRoom(
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
		r, err := qtx.CountChatRoomActionsOnChatRoom(ctx, query.CountChatRoomActionsOnChatRoomParams{
			ChatRoomID:                   chatRoomID,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to count chat room actions: %w", err)
		}
		return r, nil
	}
	runQFunc := func(orderMethod string) ([]entity.ChatRoomActionWithDetailForQuery, error) {
		r, err := qtx.GetChatRoomActionsWithDetailOnChatRoom(ctx, query.GetChatRoomActionsWithDetailOnChatRoomParams{
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
			e[i] = convChatRoomActionWithDetailOnChatRoom(v)
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
		p := query.GetChatRoomActionsWithDetailOnChatRoomUseKeysetPaginateParams{
			ChatRoomID:                   chatRoomID,
			Limit:                        limit,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
			CursorDirection:              cursorDir,
			OrderMethod:                  orderMethod,
			ActedAtCursor:                actCursor,
			Cursor:                       cursor,
		}
		r, err := qtx.GetChatRoomActionsWithDetailOnChatRoomUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		e := make([]entity.ChatRoomActionWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomActionWithDetailOnChatRoom(query.GetChatRoomActionsWithDetailOnChatRoomRow(v))
		}
		return e, nil
	}
	runQNPFunc := func(orderMethod string, limit, offset int32) ([]entity.ChatRoomActionWithDetailForQuery, error) {
		p := query.GetChatRoomActionsWithDetailOnChatRoomUseNumberedPaginateParams{
			ChatRoomID:                   chatRoomID,
			Limit:                        limit,
			Offset:                       offset,
			WhereInChatRoomActionTypeIds: where.WhereInChatRoomActionTypeIDs,
			InChatRoomActionTypeIds:      where.InChatRoomActionTypeIDs,
			OrderMethod:                  orderMethod,
		}
		r, err := qtx.GetChatRoomActionsWithDetailOnChatRoomUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat room actions: %w", err)
		}
		e := make([]entity.ChatRoomActionWithDetailForQuery, len(r))
		for i, v := range r {
			e[i] = convChatRoomActionWithDetailOnChatRoom(query.GetChatRoomActionsWithDetailOnChatRoomRow(v))
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

// GetChatRoomActionsWithDetailOnChatRoom はチャットルームアクションを取得します。
func (a *PgAdapter) GetChatRoomActionsWithDetailOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	where parameter.WhereChatRoomActionParam,
	order parameter.ChatRoomActionOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	return getChatRoomActionsWithDetailOnChatRoom(ctx, a.query, chatRoomID, where, order, np, cp, wc)
}

// GetChatRoomActionsWithDetailOnChatRoomWithSd はSD付きでチャットルームアクションを取得します。
func (a *PgAdapter) GetChatRoomActionsWithDetailOnChatRoomWithSd(
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
	return getChatRoomActionsWithDetailOnChatRoom(ctx, qtx, chatRoomID, where, order, np, cp, wc)
}

// getPluralChatRoomActionsWithDetail は複数のチャットルームアクションを取得する内部関数です。
func getPluralChatRoomActionsWithDetail(
	ctx context.Context, qtx *query.Queries, chatRoomActionIDs []uuid.UUID,
	orderMethod parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	var e []query.GetPluralChatRoomActionsWithDetailRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralChatRoomActionsWithDetail(ctx, query.GetPluralChatRoomActionsWithDetailParams{
			ChatRoomActionIds: chatRoomActionIDs,
			OrderMethod:       orderMethod.GetStringValue(),
		})
	} else {
		var ne []query.GetPluralChatRoomActionsWithDetailUseNumberedPaginateRow
		ne, err = qtx.GetPluralChatRoomActionsWithDetailUseNumberedPaginate(
			ctx, query.GetPluralChatRoomActionsWithDetailUseNumberedPaginateParams{
				Limit:             int32(np.Limit.Int64),
				Offset:            int32(np.Offset.Int64),
				ChatRoomActionIds: chatRoomActionIDs,
				OrderMethod:       orderMethod.GetStringValue(),
			})
		e = make([]query.GetPluralChatRoomActionsWithDetailRow, len(ne))
		for i, v := range ne {
			e[i] = query.GetPluralChatRoomActionsWithDetailRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.ChatRoomActionWithDetail]{},
			fmt.Errorf("failed to get chat room actions: %w", err)
	}
	entities := make([]entity.ChatRoomActionWithDetail, len(e))
	for i, v := range e {
		entities[i] = convChatRoomActionWithDetailOnChatRoom(
			query.GetChatRoomActionsWithDetailOnChatRoomRow(v)).ChatRoomActionWithDetail
	}
	return store.ListResult[entity.ChatRoomActionWithDetail]{Data: entities}, nil
}

// GetPluralChatRoomActionsWithDetail は複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActionsWithDetail(
	ctx context.Context, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	return getPluralChatRoomActionsWithDetail(ctx, a.query, chatRoomActionIDs, order, np)
}

// GetPluralChatRoomActionsWithDetailWithSd はSD付きで複数のチャットルームアクションを取得します。
func (a *PgAdapter) GetPluralChatRoomActionsWithDetailWithSd(
	ctx context.Context, sd store.Sd, chatRoomActionIDs []uuid.UUID,
	order parameter.ChatRoomActionOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.ChatRoomActionWithDetail], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.ChatRoomActionWithDetail]{}, store.ErrNotFoundDescriptor
	}
	return getPluralChatRoomActionsWithDetail(ctx, qtx, chatRoomActionIDs, order, np)
}

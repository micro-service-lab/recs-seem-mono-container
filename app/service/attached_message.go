package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageAttachedMessage 添付付きメッセージ管理サービス。
type ManageAttachedMessage struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
}

// AttachItemsOnMessage メッセージにアイテムを添付する。
func (m *ManageAttachedMessage) AttachItemsOnMessage(
	ctx context.Context,
	chatRoomID,
	messageID, ownerID uuid.UUID,
	attachments []uuid.UUID,
) (e int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	msg, err := m.DB.FindMessageWithChatRoomActionWithSd(ctx, sd, messageID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(MessageTargetMessages)
		}
		return 0, fmt.Errorf("failed to find message: %w", err)
	}
	if msg.SenderID.Bytes != ownerID {
		return 0, errhandle.NewCommonError(response.NotMessageOwner, nil)
	}
	if msg.ChatRoomAction.ChatRoomID != chatRoomID {
		return 0, errhandle.NewCommonError(response.NotMatchChatRoomMessage, nil)
	}

	if len(attachments) > 0 {
		ai, err := m.DB.GetPluralAttachableItemsWithSd(
			ctx,
			sd,
			attachments,
			parameter.AttachableItemOrderMethodDefault,
			store.NumberedPaginationParam{},
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get attachable items: %w", err)
		}
		if len(ai.Data) != len(attachments) {
			return 0, errhandle.NewModelNotFoundError(MessageTargetAttachments)
		}

		for _, v := range ai.Data {
			if !v.OwnerID.Valid {
				return 0, errhandle.NewCommonError(response.CannotAttachSystemFile, nil)
			}
			if v.OwnerID.Bytes != ownerID {
				return 0, errhandle.NewCommonError(response.NotFileOwner, nil)
			}
		}
	}

	p := make([]parameter.AttachItemMessageParam, 0, len(attachments))
	for _, a := range attachments {
		p = append(p, parameter.AttachItemMessageParam{
			MessageID:        messageID,
			AttachableItemID: entity.UUID{Valid: true, Bytes: a},
		})
	}

	if e, err = m.DB.AttacheItemsOnMessagesWithSd(
		ctx,
		sd,
		p,
	); err != nil {
		var ufe errhandle.ModelDuplicatedError
		if errors.As(err, &ufe) {
			return 0, errhandle.NewModelDuplicatedError(MessageTargetAttacheMessage)
		}
		return 0, fmt.Errorf("failed to attach items on message: %w", err)
	}

	return e, nil
}

// DetachItemsOnMessage メッセージからアイテムを添付解除する。
func (m *ManageAttachedMessage) DetachItemsOnMessage(
	ctx context.Context,
	chatRoomID,
	messageID, ownerID uuid.UUID,
	attachments []uuid.UUID,
) (e int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	msg, err := m.DB.FindMessageWithChatRoomActionWithSd(ctx, sd, messageID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(MessageTargetMessages)
		}
		return 0, fmt.Errorf("failed to find message: %w", err)
	}
	if msg.SenderID.Bytes != ownerID {
		return 0, errhandle.NewCommonError(response.NotMessageOwner, nil)
	}
	if msg.ChatRoomAction.ChatRoomID != chatRoomID {
		return 0, errhandle.NewCommonError(response.NotMatchChatRoomMessage, nil)
	}

	var ai store.ListResult[entity.AttachableItemWithContent]
	if len(attachments) > 0 {
		ai, err = m.DB.GetPluralAttachableItemsWithSd(
			ctx,
			sd,
			attachments,
			parameter.AttachableItemOrderMethodDefault,
			store.NumberedPaginationParam{},
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get attachable items: %w", err)
		}
		if len(ai.Data) != len(attachments) {
			return 0, errhandle.NewModelNotFoundError(MessageTargetAttachments)
		}
	}

	if e, err = m.DB.PluralDetachItemsOnMessageWithSd(
		ctx,
		sd,
		messageID,
		attachments,
	); err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return 0, errhandle.NewModelNotFoundError(MessageTargetAttacheMessage)
		}
		return 0, fmt.Errorf("failed to detach items on message: %w", err)
	}

	if len(ai.Data) > 0 {
		var imageIDs []uuid.UUID
		var fileIDs []uuid.UUID

		for _, v := range ai.Data {
			if v.Image.Valid {
				imageIDs = append(imageIDs, v.Image.Entity.ImageID)
			} else if v.File.Valid {
				fileIDs = append(fileIDs, v.File.Entity.FileID)
			}
		}

		if len(imageIDs) > 0 {
			defer func(imageIDs []uuid.UUID) {
				if err == nil {
					_, err = pluralDeleteImages(ctx, sd, m.DB, m.Storage, imageIDs, entity.UUID{}, true)
				}
			}(imageIDs)
		}

		if len(fileIDs) > 0 {
			defer func(fileIDs []uuid.UUID) {
				if err == nil {
					_, err = pluralDeleteFiles(ctx, sd, m.DB, m.Storage, fileIDs, entity.UUID{}, true)
				}
			}(fileIDs)
		}
	}

	return e, nil
}

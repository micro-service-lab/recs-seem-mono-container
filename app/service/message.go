package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageMessage メッセージ管理サービス。
type ManageMessage struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
	WsHub   ws.HubInterface
}

// CreateMessage メッセージを作成する。
func (m *ManageMessage) CreateMessage(
	ctx context.Context,
	senderID, chatRoomID uuid.UUID,
	content string,
	attachments []uuid.UUID,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	now := m.Clocker.Now()

	sender, err := m.DB.FindMemberWithProfileImageWithSd(ctx, sd, senderID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetSender)
		}
		return entity.Message{}, fmt.Errorf("failed to find member: %w", err)
	}
	chatRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, chatRoomID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetChatRoom)
		}
		return entity.Message{}, fmt.Errorf("failed to find chat room: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	readableMemberIDs := make([]uuid.UUID, 0, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == senderID {
			belong = true
		} else {
			readableMemberIDs = append(readableMemberIDs, v.Member.MemberID)
		}
	}
	if !belong {
		return entity.Message{}, errhandle.NewCommonError(response.NotChatRoomMember, nil)
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
			return entity.Message{}, fmt.Errorf("failed to get plural attachable items: %w", err)
		}
		if len(ai.Data) != len(attachments) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetAttachments)
		}

		for _, v := range ai.Data {
			if !v.OwnerID.Valid {
				return entity.Message{}, errhandle.NewCommonError(response.CannotAttachSystemFile, nil)
			}
			if v.OwnerID.Bytes != senderID {
				return entity.Message{}, errhandle.NewCommonError(response.NotFileOwner, nil)
			}
		}
	}

	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyMessage))
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           chatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	e, err = m.DB.CreateMessageWithSd(
		ctx,
		sd,
		parameter.CreateMessageParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			SenderID:         entity.UUID{Valid: true, Bytes: senderID},
			Body:             content,
			PostedAt:         now,
		},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	var attachmentsData []entity.AttachedItemOnMessage
	if len(attachments) > 0 {
		aiParams := make([]parameter.AttachItemMessageParam, len(attachments))
		for i, v := range attachments {
			aiParams[i] = parameter.AttachItemMessageParam{
				MessageID:        e.MessageID,
				AttachableItemID: entity.UUID{Valid: true, Bytes: v},
			}
		}
		if _, err = m.DB.AttacheItemsOnMessagesWithSd(
			ctx,
			sd,
			aiParams,
		); err != nil {
			return entity.Message{}, fmt.Errorf("failed to attache items on messages: %w", err)
		}
		attachmentsEntity, err := m.DB.GetAttachedItemsOnMessageWithSd(
			ctx,
			sd,
			e.MessageID,
			parameter.WhereAttachedItemOnMessageParam{},
			parameter.AttachedItemOnMessageOrderMethodDefault,
			store.NumberedPaginationParam{},
			store.CursorPaginationParam{},
			store.WithCountParam{},
		)
		if err != nil {
			return entity.Message{}, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		attachmentsData = attachmentsEntity.Data
	}
	if len(readableMemberIDs) > 0 {
		rrParams := make([]parameter.CreateReadReceiptParam, len(readableMemberIDs))
		for i, v := range readableMemberIDs {
			rrParams[i] = parameter.CreateReadReceiptParam{
				MessageID: e.MessageID,
				MemberID:  v,
			}
		}
		if _, err = m.DB.CreateReadReceiptsWithSd(
			ctx,
			sd,
			rrParams,
		); err != nil {
			return entity.Message{}, fmt.Errorf("failed to create read receipts: %w", err)
		}
	}

	belongMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for i, v := range belongMembers.Data {
		belongMemberIDs[i] = v.Member.MemberID
	}

	msg := entity.MessageWithSenderAndReadReceiptCountAndAttachments{
		MessageID:        e.MessageID,
		ChatRoomActionID: e.ChatRoomActionID,
		Sender: entity.NullableEntity[entity.MemberCard]{
			Valid: true,
			Entity: entity.MemberCard{
				MemberID:     sender.MemberID,
				Name:         sender.Name,
				FirstName:    sender.FirstName,
				LastName:     sender.LastName,
				Email:        sender.Email,
				ProfileImage: sender.ProfileImage,
				GradeID:      sender.GradeID,
				GroupID:      sender.GroupID,
			},
		},
		Body:             e.Body,
		PostedAt:         e.PostedAt,
		LastEditedAt:     e.LastEditedAt,
		ReadReceiptCount: 0,
		Attachments:      attachmentsData,
	}

	defer func(
		room entity.ChatRoom, belongMemberIDs []uuid.UUID,
		msg entity.MessageWithSenderAndReadReceiptCountAndAttachments,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomSentMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomSentMessageEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               msg,
				ChatRoomActionID:     cra.ChatRoomActionID,
				ChatRoomActionTypeID: cra.ChatRoomActionTypeID,
				ActedAt:              cra.ActedAt,
			})
		}
	}(chatRoom, belongMemberIDs, msg)

	return e, nil
}

// CreateMessageOnPrivateRoom 個人チャットルームにメッセージを作成する。
func (m *ManageMessage) CreateMessageOnPrivateRoom(
	ctx context.Context,
	senderID, receiverID uuid.UUID,
	content string,
	attachments []uuid.UUID,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	now := m.Clocker.Now()
	if senderID == receiverID {
		return entity.Message{}, errhandle.NewCommonError(response.NotCreateMessageToSelf, nil)
	}
	sender, err := m.DB.FindMemberWithProfileImageWithSd(ctx, sd, senderID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetSender)
		}
		return entity.Message{}, fmt.Errorf("failed to find member: %w", err)
	}
	receiver, err := m.DB.FindMemberByIDWithSd(ctx, sd, receiverID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetReceiver)
		}
		return entity.Message{}, fmt.Errorf("failed to find member: %w", err)
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
			return entity.Message{}, fmt.Errorf("failed to get plural attachable items: %w", err)
		}
		if len(ai.Data) != len(attachments) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetAttachments)
		}

		for _, v := range ai.Data {
			if !v.OwnerID.Valid {
				return entity.Message{}, errhandle.NewCommonError(response.CannotAttachSystemFile, nil)
			}
			if v.OwnerID.Bytes != senderID {
				return entity.Message{}, errhandle.NewCommonError(response.NotFileOwner, nil)
			}
		}
	}

	cr, err := m.DB.FindChatRoomOnPrivateWithSd(ctx, sd, senderID, receiverID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			cr, err = createPrivateChatRoom(ctx, sd, now, m.DB, entity.Member{
				MemberID:               senderID,
				Name:                   sender.Name,
				FirstName:              sender.FirstName,
				LastName:               sender.LastName,
				AttendStatusID:         sender.AttendStatusID,
				Email:                  sender.Email,
				ProfileImageID:         entity.UUID{Valid: sender.ProfileImage.Valid, Bytes: sender.ProfileImage.Entity.ImageID},
				GradeID:                sender.GradeID,
				GroupID:                sender.GroupID,
				PersonalOrganizationID: sender.PersonalOrganizationID,
				RoleID:                 sender.RoleID,
			}, receiver)
			if err != nil {
				return entity.Message{}, fmt.Errorf("failed to create private chat room: %w", err)
			}

			defer func(room entity.ChatRoom, membersIDs []uuid.UUID) {
				if err == nil {
					m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
						Members: membersIDs,
					}, ws.ChatRoomAddedMeEventData{
						ChatRoom: room,
					})
				}
			}(cr, []uuid.UUID{senderID, receiverID})
		} else {
			return entity.Message{}, fmt.Errorf("failed to find chat room: %w", err)
		}
	}

	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyMessage))
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           cr.ChatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	e, err = m.DB.CreateMessageWithSd(
		ctx,
		sd,
		parameter.CreateMessageParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			SenderID:         entity.UUID{Valid: true, Bytes: senderID},
			Body:             content,
			PostedAt:         now,
		},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	var attachmentsData []entity.AttachedItemOnMessage
	if len(attachments) > 0 {
		aiParams := make([]parameter.AttachItemMessageParam, len(attachments))
		for i, v := range attachments {
			aiParams[i] = parameter.AttachItemMessageParam{
				MessageID:        e.MessageID,
				AttachableItemID: entity.UUID{Valid: true, Bytes: v},
			}
		}
		if _, err = m.DB.AttacheItemsOnMessagesWithSd(
			ctx,
			sd,
			aiParams,
		); err != nil {
			return entity.Message{}, fmt.Errorf("failed to attache items on messages: %w", err)
		}

		attachmentsEntity, err := m.DB.GetAttachedItemsOnMessageWithSd(
			ctx,
			sd,
			e.MessageID,
			parameter.WhereAttachedItemOnMessageParam{},
			parameter.AttachedItemOnMessageOrderMethodDefault,
			store.NumberedPaginationParam{},
			store.CursorPaginationParam{},
			store.WithCountParam{},
		)
		if err != nil {
			return entity.Message{}, fmt.Errorf("failed to get attached items on message: %w", err)
		}
		attachmentsData = attachmentsEntity.Data
	}

	if _, err = m.DB.CreateReadReceiptWithSd(
		ctx,
		sd,
		parameter.CreateReadReceiptParam{
			MessageID: e.MessageID,
			MemberID:  receiverID,
		},
	); err != nil {
		return entity.Message{}, fmt.Errorf("failed to create read receipt: %w", err)
	}

	msg := entity.MessageWithSenderAndReadReceiptCountAndAttachments{
		MessageID:        e.MessageID,
		ChatRoomActionID: e.ChatRoomActionID,
		Sender: entity.NullableEntity[entity.MemberCard]{
			Valid: true,
			Entity: entity.MemberCard{
				MemberID:     sender.MemberID,
				Name:         sender.Name,
				FirstName:    sender.FirstName,
				LastName:     sender.LastName,
				Email:        sender.Email,
				ProfileImage: sender.ProfileImage,
				GradeID:      sender.GradeID,
				GroupID:      sender.GroupID,
			},
		},
		Body:             e.Body,
		PostedAt:         e.PostedAt,
		LastEditedAt:     e.LastEditedAt,
		ReadReceiptCount: 0,
		Attachments:      attachmentsData,
	}

	defer func(
		room entity.ChatRoom, belongMemberIDs []uuid.UUID,
		msg entity.MessageWithSenderAndReadReceiptCountAndAttachments,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomSentMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomSentMessageEventData{
				ChatRoomID:           room.ChatRoomID,
				Action:               msg,
				ChatRoomActionID:     cra.ChatRoomActionID,
				ChatRoomActionTypeID: cra.ChatRoomActionTypeID,
				ActedAt:              cra.ActedAt,
			})
		}
	}(cr, []uuid.UUID{senderID, receiverID}, msg)

	return e, nil
}

// DeleteMessage メッセージを削除する。
func (m *ManageMessage) DeleteMessage(
	ctx context.Context,
	chatRoomID,
	ownerID, messageID uuid.UUID,
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
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		return 0, fmt.Errorf("failed to find member: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	belongMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == ownerID {
			belong = true
		}
		belongMemberIDs = append(belongMemberIDs, v.Member.MemberID)
	}
	if !belong {
		return 0, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}

	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyDeleteMessage))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.UpdateChatRoomActionWithSd(
		ctx, sd, msg.ChatRoomAction.ChatRoomActionID, parameter.UpdateChatRoomActionParam{
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		})
	if err != nil {
		return 0, fmt.Errorf("failed to update chat room action: %w", err)
	}
	deletedAction, err := m.DB.CreateChatRoomDeleteMessageActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomDeleteMessageActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			DeletedBy:        entity.UUID{Valid: true, Bytes: ownerID},
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat room delete message action: %w", err)
	}

	ai, err := m.DB.GetAttachedItemsOnMessageWithSd(
		ctx,
		sd,
		messageID,
		parameter.WhereAttachedItemOnMessageParam{},
		parameter.AttachedItemOnMessageOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get attached items on message: %w", err)
	}

	// readReceiptsはカスケード削除される
	if e, err = m.DB.DeleteMessageWithSd(ctx, sd, messageID); err != nil {
		return 0, fmt.Errorf("failed to delete message: %w", err)
	}

	var imageIDs []uuid.UUID
	var fileIDs []uuid.UUID
	for _, v := range ai.Data {
		if v.AttachableItem.ImageID.Valid {
			imageIDs = append(imageIDs, v.AttachableItem.ImageID.Bytes)
		} else if v.AttachableItem.FileID.Valid {
			fileIDs = append(fileIDs, v.AttachableItem.FileID.Bytes)
		}
	}

	if len(imageIDs) > 0 {
		_, err = pluralDeleteImages(ctx, sd, m.DB, m.Storage, imageIDs, entity.UUID{}, true)
		if err != nil {
			return 0, fmt.Errorf("failed to plural delete images: %w", err)
		}
	}
	if len(fileIDs) > 0 {
		_, err = pluralDeleteFiles(ctx, sd, m.DB, m.Storage, fileIDs, entity.UUID{}, true)
		if err != nil {
			return 0, fmt.Errorf("failed to plural delete files: %w", err)
		}
	}

	action := entity.ChatRoomDeleteMessageActionWithDeletedBy{
		ChatRoomDeleteMessageActionID: deletedAction.ChatRoomDeleteMessageActionID,
		ChatRoomActionID:              deletedAction.ChatRoomActionID,
		DeletedBy: entity.NullableEntity[entity.SimpleMember]{
			Valid: true,
			Entity: entity.SimpleMember{
				MemberID:       ownerID,
				Name:           owner.Name,
				FirstName:      owner.FirstName,
				LastName:       owner.LastName,
				Email:          owner.Email,
				ProfileImageID: owner.ProfileImageID,
				GradeID:        owner.GradeID,
				GroupID:        owner.GroupID,
			},
		},
	}

	defer func(
		roomID uuid.UUID, belongMemberIDs []uuid.UUID,
		action entity.ChatRoomDeleteMessageActionWithDeletedBy,
		actAttr entity.ChatRoomAction,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomDeletedMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomDeletedMessageEventData{
				ChatRoomID:           roomID,
				Action:               action,
				ChatRoomActionID:     actAttr.ChatRoomActionID,
				ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
				ActedAt:              actAttr.ActedAt,
			})
		}
	}(chatRoomID, belongMemberIDs, action, cra)

	return e, nil
}

// ForceDeleteMessages メッセージを強制削除する。
func (m *ManageMessage) ForceDeleteMessages(
	ctx context.Context,
	messageIDs []uuid.UUID,
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
	msg, err := m.DB.GetPluralMessagesWithSd(
		ctx,
		sd,
		messageIDs,
		parameter.MessageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural messages: %w", err)
	}
	if len(msg.Data) != len(messageIDs) {
		return 0, errhandle.NewModelNotFoundError(MessageTargetMessages)
	}

	ai, err := m.DB.GetPluralAttachedItemsOnMessageWithSd(
		ctx,
		sd,
		messageIDs,
		parameter.AttachedItemOnMessageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural attached items on message: %w", err)
	}

	// readReceiptsはカスケード削除される
	if e, err = m.DB.PluralDeleteMessagesWithSd(ctx, sd, messageIDs); err != nil {
		return 0, fmt.Errorf("failed to plural delete messages: %w", err)
	}

	var imageIDs []uuid.UUID
	var fileIDs []uuid.UUID

	for _, v := range ai.Data {
		if v.AttachableItem.ImageID.Valid {
			imageIDs = append(imageIDs, v.AttachableItem.ImageID.Bytes)
		} else if v.AttachableItem.FileID.Valid {
			fileIDs = append(fileIDs, v.AttachableItem.FileID.Bytes)
		}
	}

	if len(imageIDs) > 0 {
		_, err = pluralDeleteImages(ctx, sd, m.DB, m.Storage, imageIDs, entity.UUID{}, true)
		if err != nil {
			return 0, fmt.Errorf("failed to plural delete images: %w", err)
		}
	}

	if len(fileIDs) > 0 {
		_, err = pluralDeleteFiles(ctx, sd, m.DB, m.Storage, fileIDs, entity.UUID{}, true)
		if err != nil {
			return 0, fmt.Errorf("failed to plural delete files: %w", err)
		}
	}

	return e, nil
}

// DeleteMessagesBefore チャットルームの指定期間以前のメッセージを削除する。
func (m *ManageMessage) DeleteMessagesBefore(
	ctx context.Context,
	chatRoomIDs []uuid.UUID,
	earlierPostedAt time.Time,
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
	cr, err := m.DB.GetPluralChatRoomsWithSd(
		ctx,
		sd,
		chatRoomIDs,
		parameter.ChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get plural chat rooms: %w", err)
	}
	if len(cr.Data) != len(chatRoomIDs) {
		return 0, errhandle.NewModelNotFoundError(MessageTargetChatRoom)
	}
	msg, err := m.DB.GetMessagesWithSd(
		ctx,
		sd,
		parameter.WhereMessageParam{
			WhereInChatRoom:      true,
			InChatRoom:           chatRoomIDs,
			WhereEarlierPostedAt: true,
			EarlierPostedAt:      earlierPostedAt,
		},
		parameter.MessageOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	msgIDs := make([]uuid.UUID, len(msg.Data))
	for i, v := range msg.Data {
		msgIDs[i] = v.MessageID
	}
	if len(msgIDs) == 0 {
		return 0, nil
	}
	return m.ForceDeleteMessages(ctx, msgIDs)
}

// DeleteMessagesBeforeAll は全チャットルームの指定期間以前のメッセージを削除する。
func (m *ManageMessage) DeleteMessagesBeforeAll(
	ctx context.Context,
	earlierPostedAt time.Time,
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
	msg, err := m.DB.GetMessagesWithSd(
		ctx,
		sd,
		parameter.WhereMessageParam{
			WhereEarlierPostedAt: true,
			EarlierPostedAt:      earlierPostedAt,
		},
		parameter.MessageOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	msgIDs := make([]uuid.UUID, len(msg.Data))
	for i, v := range msg.Data {
		msgIDs[i] = v.MessageID
	}
	if len(msgIDs) == 0 {
		return 0, nil
	}
	return m.ForceDeleteMessages(ctx, msgIDs)
}

// EditMessage メッセージを編集する。
func (m *ManageMessage) EditMessage(
	ctx context.Context,
	chatRoomID,
	ownerID, messageID uuid.UUID,
	content string,
) (e entity.Message, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	now := m.Clocker.Now()
	msg, err := m.DB.FindMessageWithChatRoomActionWithSd(ctx, sd, messageID)
	if err != nil {
		var nfe errhandle.ModelNotFoundError
		if errors.As(err, &nfe) {
			return entity.Message{}, errhandle.NewModelNotFoundError(MessageTargetMessages)
		}
		return entity.Message{}, fmt.Errorf("failed to find message: %w", err)
	}
	if msg.SenderID.Bytes != ownerID {
		return entity.Message{}, errhandle.NewCommonError(response.NotMessageOwner, nil)
	}
	if msg.ChatRoomAction.ChatRoomID != chatRoomID {
		return entity.Message{}, errhandle.NewCommonError(response.NotMatchChatRoomMessage, nil)
	}
	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(
		ctx,
		sd,
		chatRoomID,
		parameter.WhereMemberOnChatRoomParam{},
		parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to exists chat room belonging: %w", err)
	}
	var belong bool
	belongMemberIDs := make([]uuid.UUID, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == ownerID {
			belong = true
		}
		belongMemberIDs = append(belongMemberIDs, v.Member.MemberID)
	}
	if !belong {
		return entity.Message{}, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	if e, err = m.DB.UpdateMessageWithSd(
		ctx,
		sd,
		messageID,
		parameter.UpdateMessageParams{
			Body:         content,
			LastEditedAt: now,
		},
	); err != nil {
		return entity.Message{}, fmt.Errorf("failed to update message: %w", err)
	}

	defer func(
		roomID uuid.UUID, msg entity.Message, belongMemberIDs []uuid.UUID,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomEditedMessage, ws.Targets{
				Members: belongMemberIDs,
			}, ws.ChatRoomEditedMessageEventData{
				ChatRoomID: roomID,
				Message:    msg,
			})
		}
	}(chatRoomID, e, belongMemberIDs)

	return e, nil
}

// GetMessagesOnChatRoom チャットルームのメッセージを取得する。
func (m *ManageMessage) GetMessagesOnChatRoom(
	ctx context.Context,
	chatRoomID uuid.UUID,
	whereInSenders []uuid.UUID,
	whereSearchBody string,
	whereEarlierPostedAt time.Time,
	whereLaterPostedAt time.Time,
	whereEarlierLastEditedAt time.Time,
	whereLaterLastEditedAt time.Time,
	order parameter.MessageOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (e store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments], err error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMessageParam{
		WhereInChatRoom:          true,
		InChatRoom:               []uuid.UUID{chatRoomID},
		WhereInSender:            len(whereInSenders) > 0,
		InSender:                 whereInSenders,
		WhereLikeBody:            whereSearchBody != "",
		SearchBody:               whereSearchBody,
		WhereEarlierPostedAt:     !whereEarlierPostedAt.IsZero(),
		EarlierPostedAt:          whereEarlierPostedAt,
		WhereLaterPostedAt:       !whereLaterPostedAt.IsZero(),
		LaterPostedAt:            whereLaterPostedAt,
		WhereEarlierLastEditedAt: !whereEarlierLastEditedAt.IsZero(),
		EarlierLastEditedAt:      whereEarlierLastEditedAt,
		WhereLaterLastEditedAt:   !whereLaterLastEditedAt.IsZero(),
		LaterLastEditedAt:        whereLaterLastEditedAt,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	var message store.ListResult[entity.MessageWithSender]
	if message, err = m.DB.GetMessagesWithSender(
		ctx,
		where,
		order,
		np,
		cp,
		wc,
	); err != nil {
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to get messages: %w", err)
	}
	messageIDs := make([]uuid.UUID, len(message.Data))
	for i, v := range message.Data {
		messageIDs[i] = v.MessageID
	}

	e.CursorPagination = message.CursorPagination
	e.WithCount = message.WithCount

	rs := make(map[uuid.UUID]int64, len(messageIDs))
	reads, err := m.DB.CountReadsOnMessages(
		ctx,
		messageIDs,
		parameter.WhereReadsOnMessageParam{
			WhereIsRead: true,
		},
	)
	if err != nil {
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to count reads: %w", err)
	}
	for _, v := range reads {
		rs[v.MessageID] = v.Count
	}

	ai := make(map[uuid.UUID][]entity.AttachedItemOnMessage, len(messageIDs))
	attachments, err := m.DB.GetPluralAttachedItemsOnMessage(
		ctx,
		messageIDs,
		parameter.AttachedItemOnMessageOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return store.ListResult[entity.MessageWithSenderAndReadReceiptCountAndAttachments]{},
			fmt.Errorf("failed to get plural attached items on message: %w", err)
	}
	for _, v := range attachments.Data {
		ai[v.MessageID] = append(ai[v.MessageID], v)
	}

	e.Data = make([]entity.MessageWithSenderAndReadReceiptCountAndAttachments, len(message.Data))

	for i, v := range message.Data {
		rc, ok := rs[v.MessageID]
		if !ok {
			rc = 0
		}
		att, ok := ai[v.MessageID]
		if !ok {
			att = []entity.AttachedItemOnMessage{}
		}
		e.Data[i] = entity.MessageWithSenderAndReadReceiptCountAndAttachments{
			MessageID:        v.MessageID,
			ChatRoomActionID: v.ChatRoomActionID,
			Sender:           v.Sender,
			Body:             v.Body,
			PostedAt:         v.PostedAt,
			LastEditedAt:     v.LastEditedAt,
			ReadReceiptCount: rc,
			Attachments:      att,
		}
	}

	return e, nil
}

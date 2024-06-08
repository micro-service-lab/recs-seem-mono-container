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

// ManageChatRoom チャットルーム管理サービス。
type ManageChatRoom struct {
	DB      store.Store
	Storage storage.Storage
	Clocker clock.Clock
	WsHub   ws.HubInterface
}

// FindChatRoomByID チャットルームをIDで取得する。
func (m *ManageChatRoom) FindChatRoomByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.ChatRoom, error) {
	e, err := m.DB.FindChatRoomByID(ctx, id)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room by id: %w", err)
	}
	return e, nil
}

// FindChatRoomByIDWithCoverImage チャットルームをIDで取得する。
func (m *ManageChatRoom) FindChatRoomByIDWithCoverImage(
	ctx context.Context,
	id uuid.UUID,
) (entity.ChatRoomWithCoverImage, error) {
	e, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, id)
	if err != nil {
		return entity.ChatRoomWithCoverImage{}, fmt.Errorf("failed to find chat room by id with cover image: %w", err)
	}
	return e, nil
}

// FindPrivateChatRoom プライベートチャットルームを取得する。
func (m *ManageChatRoom) FindPrivateChatRoom(
	ctx context.Context,
	ownerID,
	memberID uuid.UUID,
) (entity.ChatRoom, error) {
	e, err := m.DB.FindChatRoomOnPrivate(ctx, ownerID, memberID)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find private chat room: %w", err)
	}
	return e, nil
}

func createChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	name string,
	coverImage entity.NullableEntity[entity.ImageWithAttachableItem],
	owner entity.Member,
	members []entity.Member,
	fromOrg bool,
) (e entity.ChatRoom, err error) {
	coverImageID := entity.UUID{}
	if coverImage.Valid {
		coverImageID = entity.UUID{Valid: true, Bytes: coverImage.Entity.ImageID}
	}
	if e, err = str.CreateChatRoomWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomParam{
			Name:             entity.String{Valid: true, String: name},
			IsPrivate:        false,
			CoverImageID:     coverImageID,
			OwnerID:          entity.UUID{Valid: true, Bytes: owner.MemberID},
			FromOrganization: fromOrg,
		},
	); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room: %w", err)
	}

	bcrp := make([]parameter.BelongChatRoomParam, 0, len(members)+1)
	bcrp = append(bcrp, parameter.BelongChatRoomParam{
		ChatRoomID: e.ChatRoomID,
		MemberID:   owner.MemberID,
		AddedAt:    now,
	})
	for _, m := range members {
		bcrp = append(bcrp, parameter.BelongChatRoomParam{
			ChatRoomID: e.ChatRoomID,
			MemberID:   m.MemberID,
			AddedAt:    now,
		})
	}

	if _, err = str.BelongChatRoomsWithSd(ctx, sd, bcrp); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	createCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	cra, err := str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: createCraType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	_, err = str.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
		ChatRoomActionID: cra.ChatRoomActionID,
		CreatedBy:        entity.UUID{Valid: true, Bytes: owner.MemberID},
		Name:             entity.String{Valid: true, String: name},
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room create action: %w", err)
	}

	addCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	addOwnerCra, err := str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	ownerCrama, err := str.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
		ChatRoomActionID: addOwnerCra.ChatRoomActionID,
		AddedBy:          entity.UUID{},
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room add member action: %w", err)
	}
	if _, err = str.AddMemberToChatRoomAddMemberActionWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomAddedMemberParam{
			ChatRoomAddMemberActionID: ownerCrama.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: owner.MemberID},
		},
	); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to add member to chat room add member action: %w", err)
	}
	if len(members) > 0 {
		addMembersCra, err := str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           e.ChatRoomID,
			ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.ChatRoom{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		membersCrama, err := str.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
			ChatRoomActionID: addMembersCra.ChatRoomActionID,
			AddedBy:          entity.UUID{Valid: true, Bytes: owner.MemberID},
		})
		if err != nil {
			return entity.ChatRoom{}, fmt.Errorf("failed to create chat room add member action: %w", err)
		}
		membersCramp := make([]parameter.CreateChatRoomAddedMemberParam, 0, len(members))
		for _, m := range members {
			membersCramp = append(membersCramp, parameter.CreateChatRoomAddedMemberParam{
				ChatRoomAddMemberActionID: membersCrama.ChatRoomAddMemberActionID,
				MemberID:                  entity.UUID{Valid: true, Bytes: m.MemberID},
			})
		}
		if _, err = str.AddMembersToChatRoomAddMemberActionWithSd(
			ctx,
			sd,
			membersCramp,
		); err != nil {
			return entity.ChatRoom{}, fmt.Errorf("failed to add members to chat room add member action: %w", err)
		}
	}

	return e, nil
}

func createPrivateChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	owner, member entity.Member,
) (e entity.ChatRoom, err error) {
	if e, err = str.CreateChatRoomWithSd(
		ctx,
		sd,
		parameter.CreateChatRoomParam{
			Name:             entity.String{},
			IsPrivate:        true,
			CoverImageID:     entity.UUID{},
			OwnerID:          entity.UUID{},
			FromOrganization: false,
		},
	); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room: %w", err)
	}

	bcrp := []parameter.BelongChatRoomParam{
		{
			ChatRoomID: e.ChatRoomID,
			MemberID:   owner.MemberID,
			AddedAt:    now,
		},
		{
			ChatRoomID: e.ChatRoomID,
			MemberID:   member.MemberID,
			AddedAt:    now,
		},
	}

	if _, err = str.BelongChatRoomsWithSd(ctx, sd, bcrp); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to belong chat rooms: %w", err)
	}

	createCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	cra, err := str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: createCraType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	_, err = str.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
		ChatRoomActionID: cra.ChatRoomActionID,
		CreatedBy:        entity.UUID{},
		Name:             entity.String{},
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room create action: %w", err)
	}

	addCraType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyAddMember))
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	addOwnerCra, err := str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           e.ChatRoomID,
		ChatRoomActionTypeID: addCraType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	crama, err := str.CreateChatRoomAddMemberActionWithSd(ctx, sd, parameter.CreateChatRoomAddMemberActionParam{
		ChatRoomActionID: addOwnerCra.ChatRoomActionID,
		AddedBy:          entity.UUID{},
	})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to create chat room add member action: %w", err)
	}
	cramp := []parameter.CreateChatRoomAddedMemberParam{
		{
			ChatRoomAddMemberActionID: crama.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: owner.MemberID},
		},
		{
			ChatRoomAddMemberActionID: crama.ChatRoomAddMemberActionID,
			MemberID:                  entity.UUID{Valid: true, Bytes: member.MemberID},
		},
	}
	if _, err = str.AddMembersToChatRoomAddMemberActionWithSd(
		ctx,
		sd,
		cramp,
	); err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to add members to chat room add member action: %w", err)
	}

	return e, nil
}

func updateChatRoom(
	ctx context.Context,
	sd store.Sd,
	now time.Time,
	str store.Store,
	stg storage.Storage,
	chatRoom entity.ChatRoom,
	name string,
	coverImage entity.NullableEntity[entity.ImageWithAttachableItem],
	owner entity.Member,
	force bool,
) (e entity.ChatRoom, action entity.ChatRoomUpdateNameActionWithUpdatedBy,
	actAttr entity.ChatRoomAction, nameUpdated bool, err error,
) {
	e = chatRoom
	if !force && e.FromOrganization {
		name = e.Name.String
	}
	if !force && e.IsPrivate {
		return entity.ChatRoom{}, entity.ChatRoomUpdateNameActionWithUpdatedBy{},
			entity.ChatRoomAction{},
			false, errhandle.NewCommonError(response.CannotUpdatePrivateChatRoom, nil)
	}

	coverImageID := entity.UUID{}
	if coverImage.Valid {
		coverImageID = entity.UUID{Valid: true, Bytes: coverImage.Entity.ImageID}
	}

	if e.Name.String != name {
		craType, err := str.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyUpdateName))
		if err != nil {
			return entity.ChatRoom{}, entity.ChatRoomUpdateNameActionWithUpdatedBy{},
				entity.ChatRoomAction{},
				false, fmt.Errorf("failed to find chat room action type by key: %w", err)
		}
		actAttr, err = str.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           e.ChatRoomID,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.ChatRoom{}, entity.ChatRoomUpdateNameActionWithUpdatedBy{},
				entity.ChatRoomAction{},
				false, fmt.Errorf("failed to create chat room action: %w", err)
		}
		updateAct, err := str.CreateChatRoomUpdateNameActionWithSd(ctx, sd,
			parameter.CreateChatRoomUpdateNameActionParam{
				ChatRoomActionID: actAttr.ChatRoomActionID,
				UpdatedBy:        entity.UUID{Valid: true, Bytes: owner.MemberID},
				Name:             name,
			})
		if err != nil {
			return entity.ChatRoom{}, entity.ChatRoomUpdateNameActionWithUpdatedBy{},
				entity.ChatRoomAction{},
				false, fmt.Errorf("failed to create chat room update name action: %w", err)
		}
		nameUpdated = true
		action = entity.ChatRoomUpdateNameActionWithUpdatedBy{
			ChatRoomUpdateNameActionID: updateAct.ChatRoomUpdateNameActionID,
			ChatRoomActionID:           actAttr.ChatRoomActionID,
			Name:                       name,
			UpdatedBy: entity.NullableEntity[entity.SimpleMember]{
				Valid: true,
				Entity: entity.SimpleMember{
					MemberID:       owner.MemberID,
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
	}
	if e.CoverImageID.Valid && e.CoverImageID.Bytes != coverImageID.Bytes {
		defer func(ownerID, imageID uuid.UUID) {
			if err == nil {
				_, err = pluralDeleteImages(
					ctx,
					sd,
					str,
					stg,
					[]uuid.UUID{imageID},
					entity.UUID{
						Valid: true,
						Bytes: ownerID,
					},
					true,
				)
			}
		}(owner.MemberID, e.CoverImageID.Bytes)
	}
	e, err = str.UpdateChatRoomWithSd(
		ctx,
		sd,
		e.ChatRoomID,
		parameter.UpdateChatRoomParams{
			Name:         entity.String{Valid: true, String: name},
			CoverImageID: coverImageID,
		},
	)
	if err != nil {
		return entity.ChatRoom{}, entity.ChatRoomUpdateNameActionWithUpdatedBy{},
			entity.ChatRoomAction{},
			false, fmt.Errorf("failed to update chat room: %w", err)
	}
	return e, action, actAttr, nameUpdated, nil
}

func deleteChatRoom(
	ctx context.Context,
	sd store.Sd,
	str store.Store,
	stg storage.Storage,
	chatRoom entity.ChatRoom,
	owner entity.Member,
	force bool,
) (c int64, err error) {
	e := chatRoom
	if !force && e.FromOrganization {
		return 0, errhandle.NewCommonError(response.CannotDeleteOrganizationChatRoom, nil)
	}
	if !force && e.IsPrivate {
		return 0, errhandle.NewCommonError(response.CannotDeletePrivateChatRoom, nil)
	}
	attachableItems, err := str.GetAttachedItemsOnChatRoomWithSd(
		ctx, sd, e.ChatRoomID,
		parameter.WhereAttachedItemOnChatRoomParam{},
		parameter.AttachedItemOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get attached items on chat room: %w", err)
	}
	var imageIDs []uuid.UUID
	var fileIDs []uuid.UUID
	for _, v := range attachableItems.Data {
		if v.AttachableItem.ImageID.Valid {
			imageIDs = append(imageIDs, v.AttachableItem.ImageID.Bytes)
		} else if v.AttachableItem.FileID.Valid {
			fileIDs = append(fileIDs, v.AttachableItem.FileID.Bytes)
		}
	}
	if e.CoverImageID.Valid {
		imageIDs = append(imageIDs, e.CoverImageID.Bytes)
	}

	if len(imageIDs) > 0 {
		defer func(imageIDs []uuid.UUID, ownerID uuid.UUID) {
			if err == nil {
				_, err = pluralDeleteImages(ctx, sd, str, stg, imageIDs, entity.UUID{
					Valid: true,
					Bytes: ownerID,
				}, true)
			}
		}(imageIDs, owner.MemberID)
	}

	if len(fileIDs) > 0 {
		defer func(fileIDs []uuid.UUID, ownerID uuid.UUID) {
			if err == nil {
				_, err = pluralDeleteFiles(ctx, sd, str, stg, fileIDs, entity.UUID{
					Valid: true,
					Bytes: ownerID,
				}, true)
			}
		}(fileIDs, owner.MemberID)
	}

	// action, message関連はカスケード削除される
	// chatRoomBelongingはカスケード削除される
	if c, err = str.DeleteChatRoomWithSd(ctx, sd, e.ChatRoomID); err != nil {
		return 0, fmt.Errorf("failed to delete chat room: %w", err)
	}
	return c, nil
}

// CreateChatRoom チャットルームを作成する。
func (m *ManageChatRoom) CreateChatRoom(
	ctx context.Context,
	name string,
	coverImageID entity.UUID,
	ownerID uuid.UUID,
	members []uuid.UUID,
) (e entity.ChatRoom, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	var coverImage entity.NullableEntity[entity.ImageWithAttachableItem]
	if coverImageID.Valid {
		image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetCoverImages)
			}
			return entity.ChatRoom{}, fmt.Errorf("failed to find image with attachable item: %w", err)
		}
		if image.AttachableItem.OwnerID.Valid && image.AttachableItem.OwnerID.Bytes != ownerID {
			return entity.ChatRoom{}, errhandle.NewCommonError(response.NotFileOwner, nil)
		}
		coverImage = entity.NullableEntity[entity.ImageWithAttachableItem]{Valid: true, Entity: image}
	}

	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetOwner)
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	pm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		members,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(pm.Data) != len(members) {
		return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetMembers)
	}
	e, err = createChatRoom(ctx, sd, m.Clocker.Now(), m.DB, name, coverImage, owner, pm.Data, false)

	defer func(room entity.ChatRoom, membersIDs []uuid.UUID) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
				Members: membersIDs,
			}, ws.ChatRoomAddedMeEventData{
				ChatRoom: room,
			})
		}
	}(e, append([]uuid.UUID{ownerID}, members...))

	return e, err
}

// CreatePrivateChatRoom プライベートチャットルームを作成する。
func (m *ManageChatRoom) CreatePrivateChatRoom(
	ctx context.Context,
	ownerID uuid.UUID,
	memberID uuid.UUID,
) (e entity.ChatRoom, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	if ownerID == memberID {
		return entity.ChatRoom{}, errhandle.NewCommonError(response.NotCreateMessageToSelf, nil)
	}
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetOwner)
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	member, err := m.DB.FindMemberByIDWithSd(ctx, sd, memberID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetMembers)
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	_, err = m.DB.FindChatRoomOnPrivateWithSd(ctx, sd, ownerID, memberID)
	exist := true
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			exist = false
		} else {
			return entity.ChatRoom{}, fmt.Errorf("failed to find private chat room: %w", err)
		}
	}
	if exist {
		return entity.ChatRoom{}, errhandle.NewCommonError(response.PrivateChatRoomAlreadyExists, nil)
	}

	e, err = createPrivateChatRoom(ctx, sd, now, m.DB, owner, member)

	defer func(room entity.ChatRoom, membersIDs []uuid.UUID) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
				Members: membersIDs,
			}, ws.ChatRoomAddedMeEventData{
				ChatRoom: room,
			})
		}
	}(e, []uuid.UUID{ownerID, memberID})

	return e, err
}

// UpdateChatRoom チャットルームを更新する。
func (m *ManageChatRoom) UpdateChatRoom(
	ctx context.Context,
	id uuid.UUID,
	name string,
	coverImageID entity.UUID,
	ownerID uuid.UUID,
) (e entity.ChatRoom, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	chatRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, id)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetChatRoom)
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find chat room by id: %w", err)
	}
	var coverImage entity.NullableEntity[entity.ImageWithAttachableItem]
	if coverImageID.Valid {
		image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetCoverImages)
			}
			return entity.ChatRoom{}, fmt.Errorf("failed to find image with attachable item: %w", err)
		}
		if image.AttachableItem.OwnerID.Valid && image.AttachableItem.OwnerID.Bytes != ownerID {
			return entity.ChatRoom{}, errhandle.NewCommonError(response.NotFileOwner, nil)
		}
		coverImage = entity.NullableEntity[entity.ImageWithAttachableItem]{Valid: true, Entity: image}
	}

	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.ChatRoom{}, errhandle.NewModelNotFoundError(ChatRoomTargetOwner)
		}
		return entity.ChatRoom{}, fmt.Errorf("failed to find member by id: %w", err)
	}

	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(ctx, sd, chatRoom.ChatRoomID,
		parameter.WhereMemberOnChatRoomParam{}, parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{}, store.CursorPaginationParam{}, store.WithCountParam{})
	if err != nil {
		return entity.ChatRoom{}, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	var exist bool
	members := make([]uuid.UUID, 0, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == ownerID {
			exist = true
		}
		members = append(members, v.Member.MemberID)
	}
	if !exist {
		return entity.ChatRoom{}, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	var nameUpdated bool
	var action entity.ChatRoomUpdateNameActionWithUpdatedBy
	var actAttr entity.ChatRoomAction
	e, action, actAttr, nameUpdated, err = updateChatRoom(
		ctx, sd, m.Clocker.Now(), m.DB, m.Storage, chatRoom, name, coverImage, owner, false)
	if nameUpdated {
		defer func(
			members []uuid.UUID, chatRoom entity.ChatRoom,
			action entity.ChatRoomUpdateNameActionWithUpdatedBy,
			actAttr entity.ChatRoomAction,
		) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomUpdatedName, ws.Targets{
					Members: members,
				}, ws.ChatRoomUpdatedNameEventData{
					ChatRoomID:           chatRoom.ChatRoomID,
					Action:               action,
					ChatRoomActionID:     actAttr.ChatRoomActionID,
					ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
					ActedAt:              actAttr.ActedAt,
				})
			}
		}(members, e, action, actAttr)
	}
	return e, err
}

// DeleteChatRoom チャットルームを削除する。
func (m *ManageChatRoom) DeleteChatRoom(
	ctx context.Context,
	id uuid.UUID,
	ownerID uuid.UUID,
) (c int64, err error) {
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
	chatRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, id)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomTargetChatRoom)
		}
		return 0, fmt.Errorf("failed to find chat room by id: %w", err)
	}
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return 0, errhandle.NewModelNotFoundError(ChatRoomTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member by id: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnChatRoomWithSd(ctx, sd, chatRoom.ChatRoomID,
		parameter.WhereMemberOnChatRoomParam{}, parameter.MemberOnChatRoomOrderMethodDefault,
		store.NumberedPaginationParam{}, store.CursorPaginationParam{}, store.WithCountParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get members on chat room: %w", err)
	}
	var exist bool
	members := make([]uuid.UUID, 0, len(belongMembers.Data))
	for _, v := range belongMembers.Data {
		if v.Member.MemberID == ownerID {
			exist = true
		}
		members = append(members, v.Member.MemberID)
	}
	if !exist {
		return 0, errhandle.NewCommonError(response.NotChatRoomMember, nil)
	}
	c, err = deleteChatRoom(ctx, sd, m.DB, m.Storage, chatRoom, owner, false)

	defer func(
		members []uuid.UUID, chatRoom entity.ChatRoom, deletedBy entity.Member,
	) {
		if err == nil {
			m.WsHub.Dispatch(ws.EventTypeChatRoomDeleted, ws.Targets{
				Members: members,
			}, ws.ChatRoomDeletedEventData{
				ChatRoom:  chatRoom,
				DeletedBy: deletedBy,
			})
		}
	}(members, chatRoom, owner)

	return c, err
}

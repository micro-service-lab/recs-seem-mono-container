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
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// ManageOrganization オーガナイゼーション管理サービス。
type ManageOrganization struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
	WsHub   ws.HubInterface
}

// Organization オーガナイゼーション。
type Organization struct {
	Name        string
	Description string
	Color       string
}

// WholeOrganization 全体オーガナイゼーション。
var WholeOrganization = Organization{
	Name:        "全体グループ",
	Description: "研究室の全員が所属するグループです",
	Color:       "#9e9e9e",
}

// CreateWholeOrganization 全体グループを作成する。
func (m *ManageOrganization) CreateWholeOrganization(
	ctx context.Context,
	name string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Organization, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	if coverImageID.Valid {
		_, err := m.DB.FindImageByIDWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to find image: %w", err)
		}
	}
	now := m.Clocker.Now()
	cr, err := m.DB.CreateChatRoomWithSd(ctx, sd, parameter.CreateChatRoomParam{
		Name:             entity.String{Valid: true, String: name},
		IsPrivate:        false,
		CoverImageID:     coverImageID,
		OwnerID:          entity.UUID{},
		FromOrganization: true,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           cr.ChatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	_, err = m.DB.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
		ChatRoomActionID: cra.ChatRoomActionID,
		CreatedBy:        entity.UUID{},
		Name:             entity.String{Valid: true, String: name},
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room create action: %w", err)
	}
	p := parameter.CreateOrganizationParam{
		Name:        name,
		Description: description,
		Color:       color,
		IsPersonal:  false,
		IsWhole:     true,
		ChatRoomID: entity.UUID{
			Valid: true,
			Bytes: cr.ChatRoomID,
		},
	}
	e, err = m.DB.CreateOrganizationWithSd(ctx, sd, p)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create organization: %w", err)
	}
	return e, nil
}

// DeleteWholeOrganization 全体グループを削除する。
func (m *ManageOrganization) DeleteWholeOrganization(ctx context.Context) (c int64, err error) {
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
	origin, err := m.DB.FindWholeOrganizationWithSd(ctx, sd)
	if err != nil {
		return 0, fmt.Errorf("failed to find whole organization: %w", err)
	}
	// organizationMemberShipはカスケード削除される
	c, err = m.DB.DeleteOrganizationWithSd(ctx, sd, origin.OrganizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	if origin.ChatRoomID.Valid {
		cr, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, origin.ChatRoomID.Bytes)
		if err != nil {
			return 0, fmt.Errorf("failed to find chat room: %w", err)
		}
		attachableItems, err := m.DB.GetAttachedItemsOnChatRoomWithSd(
			ctx, sd, origin.ChatRoomID.Bytes,
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
			if v.AttachableItem.Image.Valid {
				imageIDs = append(imageIDs, v.AttachableItem.Image.Entity.ImageID)
			} else if v.AttachableItem.File.Valid {
				fileIDs = append(fileIDs, v.AttachableItem.File.Entity.FileID)
			}
		}
		if cr.CoverImage.Valid {
			imageIDs = append(imageIDs, cr.CoverImage.Entity.ImageID)
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
		// action, message関連はカスケード削除される
		// chatRoomBelongingはカスケード削除される
		_, err = m.DB.DeleteChatRoomWithSd(ctx, sd, origin.ChatRoomID.Bytes)
		if err != nil {
			return 0, fmt.Errorf("failed to delete chat room: %w", err)
		}
	}
	return c, nil
}

// UpdateWholeOrganization 全体グループを更新する。
func (m *ManageOrganization) UpdateWholeOrganization(
	ctx context.Context,
	name string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Organization, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	origin, err := m.DB.FindWholeOrganizationWithSd(ctx, sd)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	if coverImageID.Valid {
		_, err := m.DB.FindImageByIDWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to find image: %w", err)
		}
	}
	if origin.ChatRoomID.Valid && origin.Name != name {
		now := m.Clocker.Now()
		craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyUpdateName))
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to find chat room action type by key: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           origin.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomUpdateNameActionWithSd(ctx, sd, parameter.CreateChatRoomUpdateNameActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			UpdatedBy:        entity.UUID{},
			Name:             name,
		})
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to create chat room update name action: %w", err)
		}
	}
	if origin.ChatRoomID.Valid {
		cr, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, origin.ChatRoomID.Bytes)
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to find chat room: %w", err)
		}
		if cr.CoverImage.Valid && cr.CoverImage.Entity.ImageID != coverImageID.Bytes {
			defer func(imageID uuid.UUID) {
				if err == nil {
					_, err = pluralDeleteImages(
						ctx,
						sd,
						m.DB,
						m.Storage,
						[]uuid.UUID{imageID},
						entity.UUID{},
						true,
					)
				}
			}(cr.CoverImage.Entity.ImageID)
		}
		_, err = m.DB.UpdateChatRoomWithSd(ctx, sd, origin.ChatRoomID.Bytes, parameter.UpdateChatRoomParams{
			Name:         entity.String{Valid: true, String: name},
			CoverImageID: coverImageID,
		})
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to update chat room: %w", err)
		}
	}
	p := parameter.UpdateOrganizationParams{
		Name:        name,
		Description: description,
		Color:       color,
	}
	e, err = m.DB.UpdateOrganizationWithSd(ctx, sd, origin.OrganizationID, p)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}
	return e, nil
}

// FindWholeOrganization 全体グループを取得する。
func (m *ManageOrganization) FindWholeOrganization(ctx context.Context) (entity.Organization, error) {
	e, err := m.DB.FindWholeOrganization(ctx)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to find whole organization: %w", err)
	}
	return e, nil
}

// CreateOrganization オーガナイゼーションを作成する。
func (m *ManageOrganization) CreateOrganization(
	ctx context.Context,
	name string,
	description, color entity.String,
	ownerID uuid.UUID,
	members []uuid.UUID,
	withChatRoom bool,
	chatRoomCoverImageID entity.UUID,
) (e entity.Organization, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	var inMembers bool
	for _, v := range members {
		if v == ownerID {
			inMembers = true
			break
		}
	}
	if inMembers {
		return entity.Organization{},
			errhandle.NewCommonError(response.CannotAddSelfToOrganization, nil)
	}

	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Organization{}, errhandle.NewModelNotFoundError(OrganizationTargetOwner)
		}
		return entity.Organization{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	pm, err := m.DB.GetPluralMembersWithSd(
		ctx,
		sd,
		members,
		parameter.MemberOrderMethodDefault,
		store.NumberedPaginationParam{},
	)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to get plural members: %w", err)
	}
	if len(pm.Data) != len(members) {
		return entity.Organization{}, errhandle.NewModelNotFoundError(OrganizationTargetMembers)
	}
	var cr entity.UUID
	if withChatRoom {
		var coverImage entity.NullableEntity[entity.ImageWithAttachableItem]
		if chatRoomCoverImageID.Valid {
			image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, chatRoomCoverImageID.Bytes)
			if err != nil {
				var e errhandle.ModelNotFoundError
				if errors.As(err, &e) {
					return entity.Organization{}, errhandle.NewModelNotFoundError(ChatRoomTargetCoverImages)
				}
				return entity.Organization{}, fmt.Errorf("failed to find image with attachable item by id: %w", err)
			}
			if image.AttachableItem.OwnerID.Valid && image.AttachableItem.OwnerID.Bytes != ownerID {
				return entity.Organization{}, errhandle.NewCommonError(response.NotFileOwner, nil)
			}
			coverImage = entity.NullableEntity[entity.ImageWithAttachableItem]{Valid: true, Entity: image}
		}
		ccr, err := createChatRoom(
			ctx,
			sd,
			now,
			m.DB,
			name,
			coverImage,
			owner,
			pm.Data,
			true,
		)
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to create chat room: %w", err)
		}
		cr = entity.UUID{
			Valid: true,
			Bytes: ccr.ChatRoomID,
		}
		defer func(room entity.ChatRoom, membersIDs []uuid.UUID) {
			if err == nil {
				m.WsHub.Dispatch(ws.EventTypeChatRoomAddedMe, ws.Targets{
					Members: membersIDs,
				}, ws.ChatRoomAddedMeEventData{
					ChatRoom: room,
				})
			}
		}(ccr, append([]uuid.UUID{ownerID}, members...))
	}
	e, err = m.DB.CreateOrganizationWithSd(ctx, sd, parameter.CreateOrganizationParam{
		Name:        name,
		Description: description,
		Color:       color,
		IsPersonal:  false,
		IsWhole:     false,
		ChatRoomID:  cr,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create organization: %w", err)
	}
	bop := make([]parameter.BelongOrganizationParam, 0, len(members)+1)
	bop = append(bop, parameter.BelongOrganizationParam{
		OrganizationID: e.OrganizationID,
		MemberID:       ownerID,
		WorkPositionID: entity.UUID{},
		AddedAt:        now,
	})
	for _, v := range members {
		bop = append(bop, parameter.BelongOrganizationParam{
			OrganizationID: e.OrganizationID,
			MemberID:       v,
			WorkPositionID: entity.UUID{},
			AddedAt:        now,
		})
	}

	if _, err = m.DB.BelongOrganizationsWithSd(
		ctx,
		sd,
		bop,
	); err != nil {
		return entity.Organization{}, fmt.Errorf("failed to belong organizations: %w", err)
	}
	return e, nil
}

// UpdateOrganization オーガナイゼーションを更新する。
func (m *ManageOrganization) UpdateOrganization(
	ctx context.Context,
	id uuid.UUID,
	name string,
	description, color entity.String,
	ownerID uuid.UUID,
) (e entity.Organization, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	origin, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, id)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Organization{}, errhandle.NewModelNotFoundError(OrganizationTargetOrganizations)
		}
		return entity.Organization{}, fmt.Errorf("failed to find organization by id: %w", err)
	}
	if origin.Grade.Valid {
		return entity.Organization{}, errhandle.NewCommonError(response.AttemptOperateGradeOrganization, nil)
	}
	if origin.Group.Valid {
		return entity.Organization{}, errhandle.NewCommonError(response.AttemptOperateGroupOrganization, nil)
	}
	if origin.IsPersonal {
		return entity.Organization{}, errhandle.NewCommonError(response.AttemptOperatePersonalOrganization, nil)
	}
	if origin.IsWhole {
		return entity.Organization{}, errhandle.NewCommonError(response.AttemptOperateWholeOrganization, nil)
	}

	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return entity.Organization{}, errhandle.NewModelNotFoundError(OrganizationTargetOwner)
		}
		return entity.Organization{}, fmt.Errorf("failed to find member by id: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnOrganizationWithSd(
		ctx,
		sd,
		origin.OrganizationID,
		parameter.WhereMemberOnOrganizationParam{},
		parameter.MemberOnOrganizationOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to get members on organization: %w", err)
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
		return entity.Organization{}, errhandle.NewCommonError(response.NotOrganizationMember, nil)
	}
	if origin.ChatRoomID.Valid {
		originRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, origin.ChatRoomID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return entity.Organization{}, errhandle.NewModelNotFoundError(ChatRoomTargetChatRoom)
			}
			return entity.Organization{}, fmt.Errorf("failed to find chat room by id: %w", err)
		}
		var coverImage entity.NullableEntity[entity.ImageWithAttachableItem]
		var cflag bool
		if originRoom.CoverImageID.Valid {
			image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, originRoom.CoverImageID.Bytes)
			if err != nil {
				var e errhandle.ModelNotFoundError
				if errors.As(err, &e) {
					cflag = true
				} else {
					return entity.Organization{}, fmt.Errorf("failed to find image with attachable item by id: %w", err)
				}
			}
			if image.AttachableItem.OwnerID.Valid && image.AttachableItem.OwnerID.Bytes != ownerID {
				return entity.Organization{}, errhandle.NewCommonError(response.NotFileOwner, nil)
			}
			if !cflag {
				coverImage = entity.NullableEntity[entity.ImageWithAttachableItem]{Valid: true, Entity: image}
			}
		}

		var nameUpdated bool
		var action entity.ChatRoomUpdateNameActionWithUpdatedBy
		var actAttr entity.ChatRoomAction
		_, action, actAttr, nameUpdated, err = updateChatRoom(
			ctx,
			sd,
			now,
			m.DB,
			m.Storage,
			originRoom,
			name,
			coverImage,
			owner,
			true,
		)
		if err != nil {
			return entity.Organization{}, fmt.Errorf("failed to update chat room: %w", err)
		}
		if nameUpdated {
			defer func(members []uuid.UUID, chatRoomID uuid.UUID,
				action entity.ChatRoomUpdateNameActionWithUpdatedBy,
				actAttr entity.ChatRoomAction,
			) {
				if err == nil {
					m.WsHub.Dispatch(ws.EventTypeChatRoomUpdatedName, ws.Targets{
						Members: members,
					}, ws.ChatRoomUpdatedNameEventData{
						ChatRoomID:           chatRoomID,
						Action:               action,
						ChatRoomActionID:     actAttr.ChatRoomActionID,
						ChatRoomActionTypeID: actAttr.ChatRoomActionTypeID,
						ActedAt:              actAttr.ActedAt,
					})
				}
			}(members, originRoom.ChatRoomID, action, actAttr)
		}
	}
	p := parameter.UpdateOrganizationParams{
		Name:        name,
		Description: description,
		Color:       color,
	}
	e, err = m.DB.UpdateOrganizationWithSd(ctx, sd, id, p)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}
	return e, nil
}

// DeleteOrganization オーガナイゼーションを削除する。
func (m *ManageOrganization) DeleteOrganization(
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
	origin, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, id)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return 0, errhandle.NewModelNotFoundError(OrganizationTargetOrganizations)
		}
		return 0, fmt.Errorf("failed to find organization by id: %w", err)
	}
	if origin.Grade.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGradeOrganization, nil)
	}
	if origin.Group.Valid {
		return 0, errhandle.NewCommonError(response.AttemptOperateGroupOrganization, nil)
	}
	if origin.IsPersonal {
		return 0, errhandle.NewCommonError(response.AttemptOperatePersonalOrganization, nil)
	}
	if origin.IsWhole {
		return 0, errhandle.NewCommonError(response.AttemptOperateWholeOrganization, nil)
	}
	owner, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			return 0, errhandle.NewModelNotFoundError(OrganizationTargetOwner)
		}
		return 0, fmt.Errorf("failed to find member by id: %w", err)
	}
	belongMembers, err := m.DB.GetMembersOnOrganizationWithSd(
		ctx,
		sd,
		origin.OrganizationID,
		parameter.WhereMemberOnOrganizationParam{},
		parameter.MemberOnOrganizationOrderMethodDefault,
		store.NumberedPaginationParam{},
		store.CursorPaginationParam{},
		store.WithCountParam{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get members on organization: %w", err)
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
		return 0, errhandle.NewCommonError(response.NotOrganizationMember, nil)
	}
	// organizationMemberShipはカスケード削除される
	c, err = m.DB.DeleteOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}

	if origin.ChatRoomID.Valid {
		originRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, origin.ChatRoomID.Bytes)
		if err != nil {
			var e errhandle.ModelNotFoundError
			if errors.As(err, &e) {
				return 0, errhandle.NewModelNotFoundError(ChatRoomTargetChatRoom)
			}
			return 0, fmt.Errorf("failed to find chat room by id: %w", err)
		}
		if _, err = deleteChatRoom(
			ctx,
			sd,
			m.DB,
			m.Storage,
			originRoom,
			owner,
			true,
		); err != nil {
			return 0, fmt.Errorf("failed to delete chat room: %w", err)
		}
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
		}(members, originRoom, owner)
	}
	return c, nil
}

// FindOrganizationByID オーガナイゼーションをIDで取得する。
func (m *ManageOrganization) FindOrganizationByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Organization, error) {
	e, err := m.DB.FindOrganizationByID(ctx, id)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to find organization by id: %w", err)
	}
	return e, nil
}

// FindOrganizationWithChatRoom オーガナイゼーションを取得する。
func (m *ManageOrganization) FindOrganizationWithChatRoom(
	ctx context.Context,
	id uuid.UUID,
) (entity.OrganizationWithChatRoom, error) {
	e, err := m.DB.FindOrganizationWithChatRoom(ctx, id)
	if err != nil {
		return entity.OrganizationWithChatRoom{}, fmt.Errorf("failed to find role with chat room: %w", err)
	}
	return e, nil
}

// FindOrganizationWithDetail オーガナイゼーションを取得する。
func (m *ManageOrganization) FindOrganizationWithDetail(
	ctx context.Context,
	id uuid.UUID,
) (entity.OrganizationWithDetail, error) {
	e, err := m.DB.FindOrganizationWithDetail(ctx, id)
	if err != nil {
		return entity.OrganizationWithDetail{}, fmt.Errorf("failed to find role with detail: %w", err)
	}
	return e, nil
}

// FindOrganizationWithChatRoomAndDetail オーガナイゼーションを取得する。
func (m *ManageOrganization) FindOrganizationWithChatRoomAndDetail(
	ctx context.Context,
	id uuid.UUID,
) (entity.OrganizationWithChatRoomAndDetail, error) {
	e, err := m.DB.FindOrganizationWithChatRoomAndDetail(ctx, id)
	if err != nil {
		return entity.OrganizationWithChatRoomAndDetail{},
			fmt.Errorf("failed to find role with chat room and detail: %w", err)
	}
	return e, nil
}

// GetOrganizations オーガナイゼーションを取得する。
func (m *ManageOrganization) GetOrganizations(
	ctx context.Context,
	whereSearchName string,
	whereOrganizationType parameter.WhereOrganizationType,
	wherePersonalMemberID uuid.UUID,
	order parameter.OrganizationOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Organization], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	if whereOrganizationType == parameter.WhereOrganizationTypePersonal {
		if wherePersonalMemberID == uuid.Nil {
			whereOrganizationType = parameter.WhereOrganizationTypeDefault
		}
	}
	where := parameter.WhereOrganizationParam{
		WhereLikeName:    whereSearchName != "",
		SearchName:       whereSearchName,
		WhereIsWhole:     whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		IsWhole:          whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		WhereIsPersonal:  whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		IsPersonal:       whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		PersonalMemberID: wherePersonalMemberID,
		WhereIsGroup:     whereOrganizationType == parameter.WhereOrganizationTypeGroup,
		WhereIsGrade:     whereOrganizationType == parameter.WhereOrganizationTypeGrade,
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
	r, err := m.DB.GetOrganizations(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Organization]{}, fmt.Errorf("failed to get organizations: %w", err)
	}
	return r, nil
}

// GetOrganizationsWithDetail オーガナイゼーションを取得する。
func (m *ManageOrganization) GetOrganizationsWithDetail(
	ctx context.Context,
	whereSearchName string,
	whereOrganizationType parameter.WhereOrganizationType,
	wherePersonalMemberID uuid.UUID,
	order parameter.OrganizationOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.OrganizationWithDetail], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	if whereOrganizationType == parameter.WhereOrganizationTypePersonal {
		if wherePersonalMemberID == uuid.Nil {
			whereOrganizationType = parameter.WhereOrganizationTypeDefault
		}
	}
	where := parameter.WhereOrganizationParam{
		WhereLikeName:    whereSearchName != "",
		SearchName:       whereSearchName,
		WhereIsWhole:     whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		IsWhole:          whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		WhereIsPersonal:  whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		IsPersonal:       whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		PersonalMemberID: wherePersonalMemberID,
		WhereIsGroup:     whereOrganizationType == parameter.WhereOrganizationTypeGroup,
		WhereIsGrade:     whereOrganizationType == parameter.WhereOrganizationTypeGrade,
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
	r, err := m.DB.GetOrganizationsWithDetail(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.OrganizationWithDetail]{},
			fmt.Errorf("failed to get organizations with detail: %w", err)
	}
	return r, nil
}

// GetOrganizationsWithChatRoom オーガナイゼーションを取得する。
func (m *ManageOrganization) GetOrganizationsWithChatRoom(
	ctx context.Context,
	whereSearchName string,
	whereOrganizationType parameter.WhereOrganizationType,
	wherePersonalMemberID uuid.UUID,
	order parameter.OrganizationOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.OrganizationWithChatRoom], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	if whereOrganizationType == parameter.WhereOrganizationTypePersonal {
		if wherePersonalMemberID == uuid.Nil {
			whereOrganizationType = parameter.WhereOrganizationTypeDefault
		}
	}
	where := parameter.WhereOrganizationParam{
		WhereLikeName:    whereSearchName != "",
		SearchName:       whereSearchName,
		WhereIsWhole:     whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		IsWhole:          whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		WhereIsPersonal:  whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		IsPersonal:       whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		PersonalMemberID: wherePersonalMemberID,
		WhereIsGroup:     whereOrganizationType == parameter.WhereOrganizationTypeGroup,
		WhereIsGrade:     whereOrganizationType == parameter.WhereOrganizationTypeGrade,
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
	r, err := m.DB.GetOrganizationsWithChatRoom(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.OrganizationWithChatRoom]{},
			fmt.Errorf("failed to get organizations with chat room: %w", err)
	}
	return r, nil
}

// GetOrganizationsWithChatRoomAndDetail オーガナイゼーションを取得する。
func (m *ManageOrganization) GetOrganizationsWithChatRoomAndDetail(
	ctx context.Context,
	whereSearchName string,
	whereOrganizationType parameter.WhereOrganizationType,
	wherePersonalMemberID uuid.UUID,
	order parameter.OrganizationOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.OrganizationWithChatRoomAndDetail], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	if whereOrganizationType == parameter.WhereOrganizationTypePersonal {
		if wherePersonalMemberID == uuid.Nil {
			whereOrganizationType = parameter.WhereOrganizationTypeDefault
		}
	}
	where := parameter.WhereOrganizationParam{
		WhereLikeName:    whereSearchName != "",
		SearchName:       whereSearchName,
		WhereIsWhole:     whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		IsWhole:          whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		WhereIsPersonal:  whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		IsPersonal:       whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		PersonalMemberID: wherePersonalMemberID,
		WhereIsGroup:     whereOrganizationType == parameter.WhereOrganizationTypeGroup,
		WhereIsGrade:     whereOrganizationType == parameter.WhereOrganizationTypeGrade,
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
	r, err := m.DB.GetOrganizationsWithChatRoomAndDetail(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.OrganizationWithChatRoomAndDetail]{},
			fmt.Errorf("failed to get organizations with chat room and detail: %w", err)
	}
	return r, nil
}

// GetOrganizationsCount オーガナイゼーションの数を取得する。
func (m *ManageOrganization) GetOrganizationsCount(
	ctx context.Context,
	whereSearchName string,
	whereOrganizationType parameter.WhereOrganizationType,
	wherePersonalMemberID uuid.UUID,
) (int64, error) {
	if whereOrganizationType == parameter.WhereOrganizationTypePersonal {
		if wherePersonalMemberID == uuid.Nil {
			whereOrganizationType = parameter.WhereOrganizationTypeDefault
		}
	}
	p := parameter.WhereOrganizationParam{
		WhereLikeName:    whereSearchName != "",
		SearchName:       whereSearchName,
		WhereIsWhole:     whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		IsWhole:          whereOrganizationType == parameter.WhereOrganizationTypeWhole,
		WhereIsPersonal:  whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		IsPersonal:       whereOrganizationType == parameter.WhereOrganizationTypePersonal,
		PersonalMemberID: wherePersonalMemberID,
		WhereIsGroup:     whereOrganizationType == parameter.WhereOrganizationTypeGroup,
		WhereIsGrade:     whereOrganizationType == parameter.WhereOrganizationTypeGrade,
	}
	c, err := m.DB.CountOrganizations(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get organizations count: %w", err)
	}
	return c, nil
}

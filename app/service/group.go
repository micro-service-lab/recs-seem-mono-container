package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// GroupKey 班キー。
type GroupKey string

const (
	// GroupKeyWeb 班キー: Web。
	GroupKeyWeb GroupKey = "web"
	// GroupKeyGrid 班キー: Grid。
	GroupKeyGrid GroupKey = "grid"
	// GroupKeyNetwork 班キー: Network。
	GroupKeyNetwork GroupKey = "network"
	// GroupKeyProfessor 班キー: Professor。
	GroupKeyProfessor GroupKey = "professor"
)

// Group 班。
type Group struct {
	Key         string
	Name        string
	Description string
	Color       string
}

// Groups 班一覧。
var Groups = []Group{
	{
		Key:         string(GroupKeyWeb),
		Name:        "Web班",
		Description: "Web班",
		Color:       "#77A6F7",
	},
	{
		Key:         string(GroupKeyGrid),
		Name:        "Grid班",
		Description: "Grid班",
		Color:       "#1E90FF",
	},
	{
		Key:         string(GroupKeyNetwork),
		Name:        "Network班",
		Description: "Network班",
		Color:       "#00BFFF",
	},
	{
		Key:         string(GroupKeyProfessor),
		Name:        "教授(班)",
		Description: "教授",
		Color:       "#EEE8AA",
	},
}

// ManageGroup 班管理サービス。
type ManageGroup struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
}

// CreateGroup 班を作成する。
func (m *ManageGroup) CreateGroup(
	ctx context.Context,
	name, key string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Group, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to begin transaction: %w", err)
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
			return entity.Group{}, fmt.Errorf("failed to find image: %w", err)
		}
	}
	now := m.Clocker.Now()
	crp := parameter.CreateChatRoomParam{
		Name:             name,
		IsPrivate:        false,
		CoverImageID:     coverImageID,
		OwnerID:          entity.UUID{},
		FromOrganization: true,
	}
	cr, err := m.DB.CreateChatRoomWithSd(ctx, sd, crp)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           cr.ChatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	_, err = m.DB.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
		ChatRoomActionID: cra.ChatRoomActionID,
		CreatedBy:        entity.UUID{},
		Name:             name,
	})
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to create chat room create action: %w", err)
	}
	op := parameter.CreateOrganizationParam{
		Name:        name,
		Description: description,
		Color:       color,
		IsPersonal:  false,
		IsWhole:     false,
		ChatRoomID:  entity.UUID{Bytes: cr.ChatRoomID, Valid: true},
	}
	o, err := m.DB.CreateOrganizationWithSd(ctx, sd, op)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to create organization: %w", err)
	}
	p := parameter.CreateGroupParam{
		Key:            key,
		OrganizationID: o.OrganizationID,
	}
	e, err = m.DB.CreateGroupWithSd(ctx, sd, p)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to create group: %w", err)
	}
	return e, nil
}

// CreateGroups 班を複数作成する。
func (m *ManageGroup) CreateGroups(
	ctx context.Context, ps []parameter.CreateGroupServiceParam,
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
	var imageIDs []uuid.UUID
	for _, p := range ps {
		if p.CoverImageID.Valid {
			imageIDs = append(imageIDs, p.CoverImageID.Bytes)
		}
	}
	if len(imageIDs) > 0 {
		ci, err := m.DB.GetPluralImagesWithSd(
			ctx, sd, imageIDs,
			parameter.ImageOrderMethodDefault,
			store.NumberedPaginationParam{},
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get plural images: %w", err)
		}
		if len(ci.Data) != len(imageIDs) {
			return 0, fmt.Errorf("failed to get plural images: %w", err)
		}
	}
	now := m.Clocker.Now()
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return 0, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	var p []parameter.CreateGroupParam
	for _, v := range ps {
		crp := parameter.CreateChatRoomParam{
			Name:             v.Name,
			IsPrivate:        false,
			CoverImageID:     v.CoverImageID,
			OwnerID:          entity.UUID{},
			FromOrganization: true,
		}
		cr, err := m.DB.CreateChatRoomWithSd(ctx, sd, crp)
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           cr.ChatRoomID,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			CreatedBy:        entity.UUID{},
			Name:             v.Name,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create chat room create action: %w", err)
		}
		op := parameter.CreateOrganizationParam{
			Name:        v.Name,
			Description: v.Description,
			Color:       v.Color,
			IsPersonal:  false,
			IsWhole:     false,
			ChatRoomID:  entity.UUID{Bytes: cr.ChatRoomID, Valid: true},
		}
		o, err := m.DB.CreateOrganizationWithSd(ctx, sd, op)
		if err != nil {
			return 0, fmt.Errorf("failed to create organization: %w", err)
		}
		p = append(p, parameter.CreateGroupParam{
			Key:            v.Key,
			OrganizationID: o.OrganizationID,
		})
	}
	c, err = m.DB.CreateGroupsWithSd(ctx, sd, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create groups: %w", err)
	}
	return c, nil
}

// DeleteGroup 班を削除する。
func (m *ManageGroup) DeleteGroup(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	e, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find group: %w", err)
	}
	c, err = m.DB.DeleteGroupWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete group: %w", err)
	}
	// organizationMemberShipはカスケード削除される
	_, err = m.DB.DeleteOrganization(ctx, e.Organization.OrganizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	if e.Organization.ChatRoomID.Valid {
		cr, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, e.Organization.ChatRoomID.Bytes)
		if err != nil {
			return 0, fmt.Errorf("failed to find chat room: %w", err)
		}
		attachableItems, err := m.DB.GetAttachedItemsOnChatRoomWithSd(
			ctx, sd, e.Organization.ChatRoomID.Bytes,
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
		if cr.CoverImage.Valid {
			imageIDs = append(imageIDs, cr.CoverImage.Entity.ImageID)
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
		// action, message関連はカスケード削除される
		// chatRoomBelongingはカスケード削除される
		_, err = m.DB.DeleteChatRoomWithSd(ctx, sd, e.Organization.ChatRoomID.Bytes)
		if err != nil {
			return 0, fmt.Errorf("failed to delete chat room: %w", err)
		}
	}
	return c, nil
}

// PluralDeleteGroups 班を複数削除する。
func (m *ManageGroup) PluralDeleteGroups(
	ctx context.Context, ids []uuid.UUID,
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
	es, err := m.DB.GetPluralGroupsWithOrganizationWithSd(
		ctx, sd, ids, parameter.GroupOrderMethodDefault, store.NumberedPaginationParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get plural groups: %w", err)
	}
	if len(es.Data) != len(ids) {
		return 0, fmt.Errorf("failed to get plural groups: %w", err)
	}
	var chatRoomIDs []uuid.UUID
	var organizationIDs []uuid.UUID
	for _, e := range es.Data {
		if e.Organization.ChatRoomID.Valid {
			chatRoomIDs = append(chatRoomIDs, e.Organization.ChatRoomID.Bytes)
			_, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, e.Organization.ChatRoomID.Bytes)
			if err != nil {
				return 0, fmt.Errorf("failed to find chat room: %w", err)
			}
		}
		organizationIDs = append(organizationIDs, e.Organization.OrganizationID)
	}
	c, err = m.DB.PluralDeleteGroupsWithSd(ctx, sd, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete groups: %w", err)
	}
	// organizationMemberShipはカスケード削除される
	_, err = m.DB.PluralDeleteOrganizationsWithSd(ctx, sd, organizationIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete organizations: %w", err)
	}
	if len(chatRoomIDs) > 0 {
		// chatRoomBelongingはカスケード削除される
		var imageIDs []uuid.UUID
		var fileIDs []uuid.UUID
		for _, e := range es.Data {
			attachableItems, err := m.DB.GetAttachedItemsOnChatRoomWithSd(
				ctx, sd, e.Organization.ChatRoomID.Bytes,
				parameter.WhereAttachedItemOnChatRoomParam{},
				parameter.AttachedItemOnChatRoomOrderMethodDefault,
				store.NumberedPaginationParam{},
				store.CursorPaginationParam{},
				store.WithCountParam{},
			)
			if err != nil {
				return 0, fmt.Errorf("failed to get attached items on chat room: %w", err)
			}
			for _, v := range attachableItems.Data {
				if v.AttachableItem.ImageID.Valid {
					imageIDs = append(imageIDs, v.AttachableItem.ImageID.Bytes)
				} else if v.AttachableItem.FileID.Valid {
					fileIDs = append(fileIDs, v.AttachableItem.FileID.Bytes)
				}
			}
			cr, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, e.Organization.ChatRoomID.Bytes)
			if err != nil {
				return 0, fmt.Errorf("failed to find chat room: %w", err)
			}
			if cr.CoverImageID.Valid {
				imageIDs = append(imageIDs, cr.CoverImageID.Bytes)
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
		_, err = m.DB.PluralDeleteChatRoomsWithSd(ctx, sd, chatRoomIDs)
		if err != nil {
			return 0, fmt.Errorf("failed to plural delete chat rooms: %w", err)
		}
	}
	return c, nil
}

// UpdateGroup 班を更新する。
func (m *ManageGroup) UpdateGroup(
	ctx context.Context,
	id uuid.UUID,
	name string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Group, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	fe, err := m.DB.FindGroupWithOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to find group: %w", err)
	}
	if coverImageID.Valid {
		_, err := m.DB.FindImageByIDWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to find image: %w", err)
		}
	}
	if fe.Organization.ChatRoomID.Valid && fe.Organization.Name != name {
		now := m.Clocker.Now()
		craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyUpdateName))
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to find chat room action type: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           fe.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomUpdateNameActionWithSd(ctx, sd, parameter.CreateChatRoomUpdateNameActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			UpdatedBy:        entity.UUID{},
			Name:             name,
		})
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to create chat room update name action: %w", err)
		}
	}

	if fe.Organization.ChatRoomID.Valid {
		cr, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, fe.Organization.ChatRoomID.Bytes)
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to find chat room: %w", err)
		}
		if cr.CoverImage.Valid && cr.CoverImage.Entity.ImageID != coverImageID.Bytes {
			_, err = pluralDeleteImages(
				ctx,
				sd,
				m.DB,
				m.Storage,
				[]uuid.UUID{cr.CoverImage.Entity.ImageID},
				entity.UUID{},
				true,
			)
			if err != nil {
				return entity.Group{}, fmt.Errorf("failed to plural delete images: %w", err)
			}
		}
		crp := parameter.UpdateChatRoomParams{
			Name:         name,
			CoverImageID: coverImageID,
		}
		_, err = m.DB.UpdateChatRoomWithSd(ctx, sd, fe.Organization.ChatRoomID.Bytes, crp)
		if err != nil {
			return entity.Group{}, fmt.Errorf("failed to update chat room: %w", err)
		}
	}
	op := parameter.UpdateOrganizationParams{
		Name:        name,
		Description: description,
		Color:       color,
	}
	_, err = m.DB.UpdateOrganizationWithSd(ctx, sd, fe.Organization.OrganizationID, op)
	if err != nil {
		return entity.Group{}, fmt.Errorf("failed to update organization: %w", err)
	}
	return e, nil
}

// GetGroups 班を取得する。
func (m *ManageGroup) GetGroups(
	ctx context.Context,
	order parameter.GroupOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Group], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereGroupParam{}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetGroups(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Group]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	return r, nil
}

// GetGroupsWithOrganization 班を取得する。
func (m *ManageGroup) GetGroupsWithOrganization(
	ctx context.Context,
	order parameter.GroupOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.GroupWithOrganization], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereGroupParam{}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetGroupsWithOrganization(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.GroupWithOrganization]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	return r, nil
}

// GetGroupsCount 班の数を取得する。
func (m *ManageGroup) GetGroupsCount(
	ctx context.Context,
) (int64, error) {
	p := parameter.WhereGroupParam{}
	c, err := m.DB.CountGroups(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get groups count: %w", err)
	}
	return c, nil
}

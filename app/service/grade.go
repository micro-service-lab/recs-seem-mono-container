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

// GradeKey 年次キー。
type GradeKey string

const (
	// GradeKeyBachelor1 1年生。
	GradeKeyBachelor1 GradeKey = "bachelor1"
	// GradeKeyBachelor2 2年生。
	GradeKeyBachelor2 GradeKey = "bachelor2"
	// GradeKeyBachelor3 3年生。
	GradeKeyBachelor3 GradeKey = "bachelor3"
	// GradeKeyBachelor4 4年生。
	GradeKeyBachelor4 GradeKey = "bachelor4"
	// GradeKeyMaster1 修士1年生。
	GradeKeyMaster1 GradeKey = "master1"
	// GradeKeyMaster2 修士2年生。
	GradeKeyMaster2 GradeKey = "master2"
	// GradeKeyDoctor 博士。
	GradeKeyDoctor GradeKey = "doctor"
	// GradeKeyProfessor 教授。
	GradeKeyProfessor GradeKey = "professor"
)

// Grade 年次。
type Grade struct {
	Key         string
	Name        string
	Description string
	Color       string
}

// Grades 年次一覧。
var Grades = []Grade{
	{
		Key:         string(GradeKeyBachelor1),
		Name:        "B1",
		Description: "学部1回生",
		Color:       "#FF0000",
	},
	{
		Key:         string(GradeKeyBachelor2),
		Name:        "B2",
		Description: "学部2回生",
		Color:       "#00FF00",
	},
	{
		Key:         string(GradeKeyBachelor3),
		Name:        "B3",
		Description: "学部3回生",
		Color:       "#0000FF",
	},
	{
		Key:         string(GradeKeyBachelor4),
		Name:        "B4",
		Description: "学部4回生",
		Color:       "#FFFF00",
	},
	{
		Key:         string(GradeKeyMaster1),
		Name:        "M1",
		Description: "修士1回生",
		Color:       "#FF00FF",
	},
	{
		Key:         string(GradeKeyMaster2),
		Name:        "M2",
		Description: "修士2回生",
		Color:       "#00FFFF",
	},
	{
		Key:         string(GradeKeyDoctor),
		Name:        "ドクター",
		Description: "博士",
		Color:       "#FFA500",
	},
	{
		Key:         string(GradeKeyProfessor),
		Name:        "教授",
		Description: "教授",
		Color:       "#EEE8AA",
	},
}

// ManageGrade 年次管理サービス。
type ManageGrade struct {
	DB      store.Store
	Clocker clock.Clock
	Storage storage.Storage
}

// CreateGrade 年次を作成する。
func (m *ManageGrade) CreateGrade(
	ctx context.Context,
	name, key string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Grade, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to begin transaction: %w", err)
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
			return entity.Grade{}, fmt.Errorf("failed to find image: %w", err)
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
		return entity.Grade{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyCreate))
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to find chat room action type: %w", err)
	}
	cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
		ChatRoomID:           cr.ChatRoomID,
		ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
		ActedAt:              now,
	})
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to create chat room action: %w", err)
	}
	_, err = m.DB.CreateChatRoomCreateActionWithSd(ctx, sd, parameter.CreateChatRoomCreateActionParam{
		ChatRoomActionID: cra.ChatRoomActionID,
		CreatedBy:        entity.UUID{},
		Name:             name,
	})
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to create chat room create action: %w", err)
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
		return entity.Grade{}, fmt.Errorf("failed to create organization: %w", err)
	}
	p := parameter.CreateGradeParam{
		Key:            key,
		OrganizationID: o.OrganizationID,
	}
	e, err = m.DB.CreateGradeWithSd(ctx, sd, p)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to create grade: %w", err)
	}
	return e, nil
}

// CreateGrades 年次を複数作成する。
func (m *ManageGrade) CreateGrades(
	ctx context.Context, ps []parameter.CreateGradeServiceParam,
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
	var p []parameter.CreateGradeParam
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
		p = append(p, parameter.CreateGradeParam{
			Key:            v.Key,
			OrganizationID: o.OrganizationID,
		})
	}
	c, err = m.DB.CreateGradesWithSd(ctx, sd, p)
	if err != nil {
		return 0, fmt.Errorf("failed to create grades: %w", err)
	}
	return c, nil
}

// DeleteGrade 年次を削除する。
func (m *ManageGrade) DeleteGrade(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	e, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find grade: %w", err)
	}
	c, err = m.DB.DeleteGradeWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete grade: %w", err)
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

// PluralDeleteGrades 年次を複数削除する。
func (m *ManageGrade) PluralDeleteGrades(
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
	es, err := m.DB.GetPluralGradesWithOrganizationWithSd(
		ctx, sd, ids, parameter.GradeOrderMethodDefault, store.NumberedPaginationParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get plural grades: %w", err)
	}
	if len(es.Data) != len(ids) {
		return 0, fmt.Errorf("failed to get plural grades: %w", err)
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
	c, err = m.DB.PluralDeleteGradesWithSd(ctx, sd, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete grades: %w", err)
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

// UpdateGrade 年次を更新する。
func (m *ManageGrade) UpdateGrade(
	ctx context.Context,
	id uuid.UUID,
	name string,
	description, color entity.String,
	coverImageID entity.UUID,
) (e entity.Grade, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to begin transaction: %w", err)
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
	fe, err := m.DB.FindGradeWithOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to find grade: %w", err)
	}
	if coverImageID.Valid {
		_, err := m.DB.FindImageByIDWithSd(ctx, sd, coverImageID.Bytes)
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to find image: %w", err)
		}
	}
	if fe.Organization.ChatRoomID.Valid && fe.Organization.Name != name {
		now := m.Clocker.Now()
		craType, err := m.DB.FindChatRoomActionTypeByKeyWithSd(ctx, sd, string(ChatRoomActionTypeKeyUpdateName))
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to find chat room action type: %w", err)
		}
		cra, err := m.DB.CreateChatRoomActionWithSd(ctx, sd, parameter.CreateChatRoomActionParam{
			ChatRoomID:           fe.Organization.ChatRoomID.Bytes,
			ChatRoomActionTypeID: craType.ChatRoomActionTypeID,
			ActedAt:              now,
		})
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to create chat room action: %w", err)
		}
		_, err = m.DB.CreateChatRoomUpdateNameActionWithSd(ctx, sd, parameter.CreateChatRoomUpdateNameActionParam{
			ChatRoomActionID: cra.ChatRoomActionID,
			UpdatedBy:        entity.UUID{},
			Name:             name,
		})
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to create chat room update name action: %w", err)
		}
	}

	if fe.Organization.ChatRoomID.Valid {
		cr, err := m.DB.FindChatRoomByIDWithCoverImage(ctx, fe.Organization.ChatRoomID.Bytes)
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to find chat room: %w", err)
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
				return entity.Grade{}, fmt.Errorf("failed to plural delete images: %w", err)
			}
		}
		crp := parameter.UpdateChatRoomParams{
			Name:         name,
			IsPrivate:    false,
			CoverImageID: coverImageID,
			OwnerID:      entity.UUID{},
		}
		_, err = m.DB.UpdateChatRoomWithSd(ctx, sd, fe.Organization.ChatRoomID.Bytes, crp)
		if err != nil {
			return entity.Grade{}, fmt.Errorf("failed to update chat room: %w", err)
		}
	}
	op := parameter.UpdateOrganizationParams{
		Name:        name,
		Description: description,
		Color:       color,
	}
	_, err = m.DB.UpdateOrganizationWithSd(ctx, sd, fe.Organization.OrganizationID, op)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to update organization: %w", err)
	}
	return e, nil
}

// GetGrades 年次を取得する。
func (m *ManageGrade) GetGrades(
	ctx context.Context,
	order parameter.GradeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Grade], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereGradeParam{}
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
	r, err := m.DB.GetGrades(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Grade]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	return r, nil
}

// GetGradesWithOrganization 年次を取得する。
func (m *ManageGrade) GetGradesWithOrganization(
	ctx context.Context,
	order parameter.GradeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.GradeWithOrganization], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereGradeParam{}
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
	r, err := m.DB.GetGradesWithOrganization(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.GradeWithOrganization]{}, fmt.Errorf("failed to get grades: %w", err)
	}
	return r, nil
}

// GetGradesCount 年次の数を取得する。
func (m *ManageGrade) GetGradesCount(
	ctx context.Context,
) (int64, error) {
	p := parameter.WhereGradeParam{}
	c, err := m.DB.CountGrades(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get grades count: %w", err)
	}
	return c, nil
}

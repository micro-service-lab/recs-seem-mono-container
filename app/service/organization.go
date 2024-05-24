package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ManageOrganization オーガナイゼーション管理サービス。
type ManageOrganization struct {
	DB store.Store
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
	Description: "研究室の全員が所属するグループです。",
	Color:       "#FF0000",
}

// CreateWholeOrganization 全体グループを作成する。
func (m *ManageOrganization) CreateWholeOrganization(
	ctx context.Context,
	name,
	description,
	color string,
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
	cr, err := m.DB.CreateChatRoomWithSd(ctx, sd, parameter.CreateChatRoomParam{
		Name:             WholeOrganization.Name,
		IsPrivate:        false,
		CoverImageID:     entity.UUID{},
		OwnerID:          entity.UUID{},
		FromOrganization: true,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	p := parameter.CreateOrganizationParam{
		Name:        name,
		Description: entity.String{Valid: true, String: description},
		Color:       entity.String{Valid: true, String: color},
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

// CreateOrganization オーガナイゼーションを作成する。
func (m *ManageOrganization) CreateOrganization(
	ctx context.Context, name string, description, color entity.String,
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
	cr, err := m.DB.CreateChatRoomWithSd(ctx, sd, parameter.CreateChatRoomParam{
		Name:             name,
		IsPrivate:        false,
		CoverImageID:     entity.UUID{},
		OwnerID:          entity.UUID{},
		FromOrganization: true,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	p := parameter.CreateOrganizationParam{
		Name:        name,
		Description: description,
		Color:       color,
		IsPersonal:  false,
		IsWhole:     false,
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

// UpdateOrganization オーガナイゼーションを更新する。
func (m *ManageOrganization) UpdateOrganization(
	ctx context.Context, id uuid.UUID, name string, description, color entity.String,
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
	origin, err := m.DB.FindOrganizationWithDetailWithSd(ctx, sd, id)
	if err != nil {
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
	p := parameter.UpdateOrganizationParams{
		Name:        name,
		Description: description,
		Color:       color,
	}
	e, err = m.DB.UpdateOrganizationWithSd(ctx, sd, id, p)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}
	originRoom, err := m.DB.FindChatRoomByIDWithSd(ctx, sd, e.ChatRoomID.Bytes)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to find chat room by id: %w", err)
	}
	_, err = m.DB.UpdateChatRoomWithSd(ctx, sd, originRoom.ChatRoomID, parameter.UpdateChatRoomParams{
		Name:         name,
		IsPrivate:    originRoom.IsPrivate,
		CoverImageID: originRoom.CoverImageID,
		OwnerID:      originRoom.OwnerID,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to update chat room: %w", err)
	}
	return e, nil
}

// DeleteOrganization オーガナイゼーションを削除する。
func (m *ManageOrganization) DeleteOrganization(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	c, err = m.DB.DeleteOrganizationWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete organization: %w", err)
	}
	_, err = m.DB.DeleteChatRoomWithSd(ctx, sd, origin.ChatRoomID.Bytes)
	if err != nil {
		return 0, fmt.Errorf("failed to delete chat room: %w", err)
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

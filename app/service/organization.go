package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageOrganization オーガナイゼーション管理サービス。
type ManageOrganization struct {
	DB store.Store
}

// createOrganization オーガナイゼーションを作成する。
func (m *ManageOrganization) createOrganization(
	ctx context.Context,
	name string, description,
	color entity.String,
	isPersonal, isWhole bool,
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
		FromOrganization: false,
	})
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to create chat room: %w", err)
	}
	p := parameter.CreateOrganizationParam{
		Name:        name,
		Description: description,
		Color:       color,
		IsPersonal:  isPersonal,
		IsWhole:     isWhole,
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

// CreateOrganizations オーガナイゼーションを複数作成する。
func (m *ManageOrganization) CreateOrganizations(
	ctx context.Context, ps []parameter.CreateOrganizationParam,
) (int64, error) {
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
	es, err := m.DB.CreateOrganizations(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create roles: %w", err)
	}
	return es, nil
}

// UpdateOrganization オーガナイゼーションを更新する。
func (m *ManageOrganization) UpdateOrganization(
	ctx context.Context, id uuid.UUID, name, description string,
) (entity.Organization, error) {
	p := parameter.UpdateOrganizationParams{
		Name: name,
		// Description: description,
	}
	e, err := m.DB.UpdateOrganization(ctx, id, p)
	if err != nil {
		return entity.Organization{}, fmt.Errorf("failed to update role: %w", err)
	}
	return e, nil
}

// DeleteOrganization オーガナイゼーションを削除する。
func (m *ManageOrganization) DeleteOrganization(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteOrganization(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete role: %w", err)
	}
	return c, nil
}

// PluralDeleteOrganizations オーガナイゼーションを複数削除する。
func (m *ManageOrganization) PluralDeleteOrganizations(ctx context.Context, ids []uuid.UUID) (int64, error) {
	c, err := m.DB.PluralDeleteOrganizations(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete roles: %w", err)
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
		return entity.Organization{}, fmt.Errorf("failed to find role by id: %w", err)
	}
	return e, nil
}

// GetOrganizations オーガナイゼーションを取得する。
func (m *ManageOrganization) GetOrganizations(
	ctx context.Context,
	whereSearchName string,
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
	where := parameter.WhereOrganizationParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
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
		return store.ListResult[entity.Organization]{}, fmt.Errorf("failed to get roles: %w", err)
	}
	return r, nil
}

// GetOrganizationsCount オーガナイゼーションの数を取得する。
func (m *ManageOrganization) GetOrganizationsCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereOrganizationParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountOrganizations(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get roles count: %w", err)
	}
	return c, nil
}

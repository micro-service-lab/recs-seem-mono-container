package pgadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/query"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

func convGroupWithOrganization(e query.FindGroupByIDWithOrganizationRow) entity.GroupWithOrganization {
	return entity.GroupWithOrganization{
		GroupID: e.GroupID,
		Key:     e.Key,
		Organization: entity.Organization{
			OrganizationID: e.OrganizationID,
			Name:           e.OrganizationName.String,
			Description:    entity.String(e.OrganizationDescription),
			Color:          entity.String(e.OrganizationColor),
			IsPersonal:     e.OrganizationIsPersonal.Bool,
			IsWhole:        e.OrganizationIsWhole.Bool,
			ChatRoomID:     entity.UUID(e.OrganizationChatRoomID),
		},
	}
}

// countGroups はグループ数を取得する内部関数です。
func countGroups(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereGroupParam,
) (int64, error) {
	c, err := qtx.CountGroups(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count groups: %w", err)
	}
	return c, nil
}

// CountGroups はグループ数を取得します。
func (a *PgAdapter) CountGroups(ctx context.Context, where parameter.WhereGroupParam) (int64, error) {
	return countGroups(ctx, a.query, where)
}

// CountGroupsWithSd はSD付きでグループ数を取得します。
func (a *PgAdapter) CountGroupsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGroupParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return countGroups(ctx, qtx, where)
}

// createGroup はグループを作成する内部関数です。
func createGroup(
	ctx context.Context, qtx *query.Queries, param parameter.CreateGroupParam,
) (entity.Group, error) {
	p := query.CreateGroupParams{
		Key:            param.Key,
		OrganizationID: param.OrganizationID,
	}
	e, err := qtx.CreateGroup(ctx, p)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return entity.Group{}, errhandle.NewModelDuplicatedError("group")
		}
		return entity.Group{}, fmt.Errorf("failed to create group: %w", err)
	}
	entity := entity.Group{
		GroupID:        e.GroupID,
		Key:            e.Key,
		OrganizationID: e.OrganizationID,
	}
	return entity, nil
}

// CreateGroup はグループを作成します。
func (a *PgAdapter) CreateGroup(
	ctx context.Context, param parameter.CreateGroupParam,
) (entity.Group, error) {
	return createGroup(ctx, a.query, param)
}

// CreateGroupWithSd はSD付きでグループを作成します。
func (a *PgAdapter) CreateGroupWithSd(
	ctx context.Context, sd store.Sd, param parameter.CreateGroupParam,
) (entity.Group, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Group{}, store.ErrNotFoundDescriptor
	}
	return createGroup(ctx, qtx, param)
}

// createGroups は複数のグループを作成する内部関数です。
func createGroups(
	ctx context.Context, qtx *query.Queries, params []parameter.CreateGroupParam,
) (int64, error) {
	param := make([]query.CreateGroupsParams, len(params))
	for i, p := range params {
		param[i] = query.CreateGroupsParams{
			Key:            p.Key,
			OrganizationID: p.OrganizationID,
		}
	}
	n, err := qtx.CreateGroups(ctx, param)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniquenessViolationCode {
			return 0, errhandle.NewModelDuplicatedError("group")
		}
		return 0, fmt.Errorf("failed to create groups: %w", err)
	}
	return n, nil
}

// CreateGroups は複数のグループを作成します。
func (a *PgAdapter) CreateGroups(
	ctx context.Context, params []parameter.CreateGroupParam,
) (int64, error) {
	return createGroups(ctx, a.query, params)
}

// CreateGroupsWithSd はSD付きで複数のグループを作成します。
func (a *PgAdapter) CreateGroupsWithSd(
	ctx context.Context, sd store.Sd, params []parameter.CreateGroupParam,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return createGroups(ctx, qtx, params)
}

// deleteGroup はグループを削除する内部関数です。
func deleteGroup(ctx context.Context, qtx *query.Queries, groupID uuid.UUID) (int64, error) {
	c, err := qtx.DeleteGroup(ctx, groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete group: %w", err)
	}
	if c != 1 {
		return 0, errhandle.NewModelNotFoundError("group")
	}
	return c, nil
}

// DeleteGroup はグループを削除します。
func (a *PgAdapter) DeleteGroup(ctx context.Context, groupID uuid.UUID) (int64, error) {
	return deleteGroup(ctx, a.query, groupID)
}

// DeleteGroupWithSd はSD付きでグループを削除します。
func (a *PgAdapter) DeleteGroupWithSd(
	ctx context.Context, sd store.Sd, groupID uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return deleteGroup(ctx, qtx, groupID)
}

// pluralDeleteGroups は複数のグループを削除する内部関数です。
func pluralDeleteGroups(ctx context.Context, qtx *query.Queries, groupIDs []uuid.UUID) (int64, error) {
	c, err := qtx.PluralDeleteGroups(ctx, groupIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete groups: %w", err)
	}
	if c != int64(len(groupIDs)) {
		return 0, errhandle.NewModelNotFoundError("group")
	}
	return c, nil
}

// PluralDeleteGroups は複数のグループを削除します。
func (a *PgAdapter) PluralDeleteGroups(ctx context.Context, groupIDs []uuid.UUID) (int64, error) {
	return pluralDeleteGroups(ctx, a.query, groupIDs)
}

// PluralDeleteGroupsWithSd はSD付きで複数のグループを削除します。
func (a *PgAdapter) PluralDeleteGroupsWithSd(
	ctx context.Context, sd store.Sd, groupIDs []uuid.UUID,
) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return 0, store.ErrNotFoundDescriptor
	}
	return pluralDeleteGroups(ctx, qtx, groupIDs)
}

// findGroupByID はグループをIDで取得する内部関数です。
func findGroupByID(
	ctx context.Context, qtx *query.Queries, groupID uuid.UUID,
) (entity.Group, error) {
	e, err := qtx.FindGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Group{}, errhandle.NewModelNotFoundError("group")
		}
		return entity.Group{}, fmt.Errorf("failed to find group: %w", err)
	}
	entity := entity.Group{
		GroupID:        e.GroupID,
		Key:            e.Key,
		OrganizationID: e.OrganizationID,
	}
	return entity, nil
}

// FindGroupByID はグループをIDで取得します。
func (a *PgAdapter) FindGroupByID(ctx context.Context, groupID uuid.UUID) (entity.Group, error) {
	return findGroupByID(ctx, a.query, groupID)
}

// FindGroupByIDWithSd はSD付きでグループをIDで取得します。
func (a *PgAdapter) FindGroupByIDWithSd(
	ctx context.Context, sd store.Sd, groupID uuid.UUID,
) (entity.Group, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Group{}, store.ErrNotFoundDescriptor
	}
	return findGroupByID(ctx, qtx, groupID)
}

func findGroupByKey(
	ctx context.Context, qtx *query.Queries, key string,
) (entity.Group, error) {
	e, err := qtx.FindGroupByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Group{}, errhandle.NewModelNotFoundError("group")
		}
		return entity.Group{}, fmt.Errorf("failed to find group: %w", err)
	}
	entity := entity.Group{
		GroupID:        e.GroupID,
		Key:            e.Key,
		OrganizationID: e.OrganizationID,
	}
	return entity, nil
}

// FindGroupByKey はグループをキーで取得します。
func (a *PgAdapter) FindGroupByKey(ctx context.Context, key string) (entity.Group, error) {
	return findGroupByKey(ctx, a.query, key)
}

// FindGroupByKeyWithSd はSD付きでグループをキーで取得します。
func (a *PgAdapter) FindGroupByKeyWithSd(
	ctx context.Context, sd store.Sd, key string,
) (entity.Group, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.Group{}, store.ErrNotFoundDescriptor
	}
	return findGroupByKey(ctx, qtx, key)
}

// findGroupWithOrganization はグループとオーガナイゼーションを取得する内部関数です。
func findGroupWithOrganization(
	ctx context.Context, qtx *query.Queries, groupID uuid.UUID,
) (entity.GroupWithOrganization, error) {
	e, err := qtx.FindGroupByIDWithOrganization(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.GroupWithOrganization{}, errhandle.NewModelNotFoundError("group")
		}
		return entity.GroupWithOrganization{}, fmt.Errorf("failed to find group with organization: %w", err)
	}
	return convGroupWithOrganization(e), nil
}

// FindGroupWithOrganization はグループとオーガナイゼーションを取得します。
func (a *PgAdapter) FindGroupWithOrganization(
	ctx context.Context, groupID uuid.UUID,
) (entity.GroupWithOrganization, error) {
	return findGroupWithOrganization(ctx, a.query, groupID)
}

// FindGroupWithOrganizationWithSd はSD付きでグループとオーガナイゼーションを取得します。
func (a *PgAdapter) FindGroupWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, groupID uuid.UUID,
) (entity.GroupWithOrganization, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return entity.GroupWithOrganization{}, store.ErrNotFoundDescriptor
	}
	return findGroupWithOrganization(ctx, qtx, groupID)
}

// getGroups はグループを取得する内部関数です。
func getGroups(
	ctx context.Context,
	qtx *query.Queries,
	_ parameter.WhereGroupParam,
	order parameter.GroupOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam,
	wc store.WithCountParam,
) (store.ListResult[entity.Group], error) {
	eConvFunc := func(e query.Group) (entity.Group, error) {
		return entity.Group{
			GroupID:        e.GroupID,
			Key:            e.Key,
			OrganizationID: e.OrganizationID,
		}, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountGroups(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count groups: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]query.Group, error) {
		r, err := qtx.GetGroups(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []query.Group{}, nil
			}
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		return r, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]query.Group, error) {
		p := query.GetGroupsUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetGroupsUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		return r, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]query.Group, error) {
		p := query.GetGroupsUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetGroupsUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		return r, nil
	}
	selector := func(subCursor string, e query.Group) (entity.Int, any) {
		switch subCursor {
		case parameter.GroupDefaultCursorKey:
			return entity.Int(e.MGroupsPkey), nil
		}
		return entity.Int(e.MGroupsPkey), nil
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
		return store.ListResult[entity.Group]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	return res, nil
}

// GetGroups はグループを取得します。
func (a *PgAdapter) GetGroups(
	ctx context.Context, where parameter.WhereGroupParam,
	order parameter.GroupOrderMethod,
	np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Group], error) {
	return getGroups(ctx, a.query, where, order, np, cp, wc)
}

// GetGroupsWithSd はSD付きでグループを取得します。
func (a *PgAdapter) GetGroupsWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGroupParam,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.Group], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Group]{}, store.ErrNotFoundDescriptor
	}
	return getGroups(ctx, qtx, where, order, np, cp, wc)
}

// getPluralGroups は複数のグループを取得する内部関数です。
func getPluralGroups(
	ctx context.Context, qtx *query.Queries, groupIDs []uuid.UUID,
	_ parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Group], error) {
	var e []query.Group
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralGroups(ctx, groupIDs)
	} else {
		e, err = qtx.GetPluralGroupsUseNumberedPaginate(ctx, query.GetPluralGroupsUseNumberedPaginateParams{
			GroupIds: groupIDs,
			Offset:   int32(np.Offset.Int64),
			Limit:    int32(np.Limit.Int64),
		})
	}
	if err != nil {
		return store.ListResult[entity.Group]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	entities := make([]entity.Group, len(e))
	for i, v := range e {
		entities[i] = entity.Group{
			GroupID:        v.GroupID,
			Key:            v.Key,
			OrganizationID: v.OrganizationID,
		}
	}
	return store.ListResult[entity.Group]{Data: entities}, nil
}

// GetPluralGroups は複数のグループを取得します。
func (a *PgAdapter) GetPluralGroups(
	ctx context.Context, groupIDs []uuid.UUID,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Group], error) {
	return getPluralGroups(ctx, a.query, groupIDs, order, np)
}

// GetPluralGroupsWithSd はSD付きで複数のグループを取得します。
func (a *PgAdapter) GetPluralGroupsWithSd(
	ctx context.Context, sd store.Sd, groupIDs []uuid.UUID,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.Group], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.Group]{}, store.ErrNotFoundDescriptor
	}
	return getPluralGroups(ctx, qtx, groupIDs, order, np)
}

func getGroupsWithOrganization(
	ctx context.Context, qtx *query.Queries, _ parameter.WhereGroupParam,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	eConvFunc := func(e entity.GroupWithOrganizationForQuery) (entity.GroupWithOrganization, error) {
		return e.GroupWithOrganization, nil
	}
	runCFunc := func() (int64, error) {
		r, err := qtx.CountGroups(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to count groups: %w", err)
		}
		return r, nil
	}
	runQFunc := func(_ string) ([]entity.GroupWithOrganizationForQuery, error) {
		r, err := qtx.GetGroupsWithOrganization(ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []entity.GroupWithOrganizationForQuery{}, nil
			}
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		e := make([]entity.GroupWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GroupWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGroupsPkey),
				GroupWithOrganization: convGroupWithOrganization(query.FindGroupByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQCPFunc := func(_, _ string,
		limit int32, cursorDir string, cursor int32, _ any,
	) ([]entity.GroupWithOrganizationForQuery, error) {
		p := query.GetGroupsWithOrganizationUseKeysetPaginateParams{
			Limit:           limit,
			CursorDirection: cursorDir,
			Cursor:          cursor,
		}
		r, err := qtx.GetGroupsWithOrganizationUseKeysetPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		e := make([]entity.GroupWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GroupWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGroupsPkey),
				GroupWithOrganization: convGroupWithOrganization(query.FindGroupByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	runQNPFunc := func(_ string, limit, offset int32) ([]entity.GroupWithOrganizationForQuery, error) {
		p := query.GetGroupsWithOrganizationUseNumberedPaginateParams{
			Limit:  limit,
			Offset: offset,
		}
		r, err := qtx.GetGroupsWithOrganizationUseNumberedPaginate(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		e := make([]entity.GroupWithOrganizationForQuery, len(r))
		for i, v := range r {
			e[i] = entity.GroupWithOrganizationForQuery{
				Pkey:                  entity.Int(v.MGroupsPkey),
				GroupWithOrganization: convGroupWithOrganization(query.FindGroupByIDWithOrganizationRow(v)),
			}
		}
		return e, nil
	}
	selector := func(subCursor string, e entity.GroupWithOrganizationForQuery) (entity.Int, any) {
		switch subCursor {
		case parameter.GroupDefaultCursorKey:
			return entity.Int(e.Pkey), nil
		}
		return entity.Int(e.Pkey), nil
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
		return store.ListResult[entity.GroupWithOrganization]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	return res, nil
}

// GetGroupsWithOrganization はグループとオーガナイゼーションを取得します。
func (a *PgAdapter) GetGroupsWithOrganization(
	ctx context.Context, where parameter.WhereGroupParam,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	return getGroupsWithOrganization(ctx, a.query, where, order, np, cp, wc)
}

// GetGroupsWithOrganizationWithSd はSD付きでグループとオーガナイゼーションを取得します。
func (a *PgAdapter) GetGroupsWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, where parameter.WhereGroupParam,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
	cp store.CursorPaginationParam, wc store.WithCountParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.GroupWithOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getGroupsWithOrganization(ctx, qtx, where, order, np, cp, wc)
}

// getPluralGroupsWithOrganization は複数のグループを取得する内部関数です。
func getPluralGroupsWithOrganization(
	ctx context.Context, qtx *query.Queries, groupIDs []uuid.UUID,
	_ parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	var e []query.GetPluralGroupsWithOrganizationRow
	var te []query.GetPluralGroupsWithOrganizationUseNumberedPaginateRow
	var err error
	if !np.Valid {
		e, err = qtx.GetPluralGroupsWithOrganization(ctx, groupIDs)
	} else {
		te, err = qtx.GetPluralGroupsWithOrganizationUseNumberedPaginate(
			ctx, query.GetPluralGroupsWithOrganizationUseNumberedPaginateParams{
				GroupIds: groupIDs,
				Offset:   int32(np.Offset.Int64),
				Limit:    int32(np.Limit.Int64),
			})
		e = make([]query.GetPluralGroupsWithOrganizationRow, len(te))
		for i, v := range te {
			e[i] = query.GetPluralGroupsWithOrganizationRow(v)
		}
	}
	if err != nil {
		return store.ListResult[entity.GroupWithOrganization]{}, fmt.Errorf("failed to get groups: %w", err)
	}
	entities := make([]entity.GroupWithOrganization, len(e))
	for i, v := range e {
		entities[i] = convGroupWithOrganization(query.FindGroupByIDWithOrganizationRow(v))
	}
	return store.ListResult[entity.GroupWithOrganization]{Data: entities}, nil
}

// GetPluralGroupsWithOrganization は複数のグループを取得します。
func (a *PgAdapter) GetPluralGroupsWithOrganization(
	ctx context.Context, groupIDs []uuid.UUID,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	return getPluralGroupsWithOrganization(ctx, a.query, groupIDs, order, np)
}

// GetPluralGroupsWithOrganizationWithSd はSD付きで複数のグループを取得します。
func (a *PgAdapter) GetPluralGroupsWithOrganizationWithSd(
	ctx context.Context, sd store.Sd, groupIDs []uuid.UUID,
	order parameter.GroupOrderMethod, np store.NumberedPaginationParam,
) (store.ListResult[entity.GroupWithOrganization], error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	qtx, ok := a.qtxMap[sd]
	if !ok {
		return store.ListResult[entity.GroupWithOrganization]{}, store.ErrNotFoundDescriptor
	}
	return getPluralGroupsWithOrganization(ctx, qtx, groupIDs, order, np)
}

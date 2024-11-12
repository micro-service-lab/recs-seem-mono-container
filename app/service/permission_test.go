package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils"
)

func TestManagePermission_CreatePermission(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreatePermissionParam
	}
	cases := []struct {
		name        string
		key         string
		description string
		categoryID  uuid.UUID
		want        wants
	}{
		{
			name:        "name",
			key:         "key",
			description: "description",
			categoryID:  testutils.FixedUUID(t, 0),
			want: wants{
				param: parameter.CreatePermissionParam{
					Name:                 "name",
					Key:                  "key",
					Description:          "description",
					PermissionCategoryID: testutils.FixedUUID(t, 0),
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreatePermissionFunc: func(
			_ context.Context, p parameter.CreatePermissionParam,
		) (entity.Permission, error) {
			return entity.Permission{
				PermissionID:         uuid.New(),
				Name:                 p.Name,
				Key:                  p.Key,
				Description:          p.Description,
				PermissionCategoryID: p.PermissionCategoryID,
			}, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePermission(ctx, c.name, c.key, c.description, c.categoryID)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePermissionCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePermission_CreatePermissions(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreatePermissionParam
	}
	cases := []struct {
		param []parameter.CreatePermissionParam
		want  wants
	}{
		{
			param: []parameter.CreatePermissionParam{
				{Name: "name1", Key: "key1", Description: "description1", PermissionCategoryID: testutils.FixedUUID(t, 0)},
				{Name: "name2", Key: "key2", Description: "description2", PermissionCategoryID: testutils.FixedUUID(t, 1)},
			},
			want: wants{
				params: []parameter.CreatePermissionParam{
					{Name: "name1", Key: "key1", Description: "description1", PermissionCategoryID: testutils.FixedUUID(t, 0)},
					{Name: "name2", Key: "key2", Description: "description2", PermissionCategoryID: testutils.FixedUUID(t, 1)},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreatePermissionsFunc: func(
			_ context.Context, _ []parameter.CreatePermissionParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePermissions(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePermissionsCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManagePermission_UpdatePermission(t *testing.T) {
	t.Parallel()
	type wants struct {
		permissionID uuid.UUID
		param        parameter.UpdatePermissionParams
	}
	cases := []struct {
		id          uuid.UUID
		name        string
		key         string
		description string
		categoryID  uuid.UUID
		want        wants
	}{
		{
			id:          testutils.FixedUUID(t, 0),
			name:        "update name",
			key:         "update key",
			description: "update description",
			categoryID:  testutils.FixedUUID(t, 0),
			want: wants{
				permissionID: testutils.FixedUUID(t, 0),
				param: parameter.UpdatePermissionParams{
					Name:                 "update name",
					Key:                  "update key",
					Description:          "update description",
					PermissionCategoryID: testutils.FixedUUID(t, 0),
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdatePermissionFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdatePermissionParams,
		) (entity.Permission, error) {
			return entity.Permission{}, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdatePermission(ctx, c.id, c.name, c.key, c.description, c.categoryID)
		assert.NoError(t, err)
	}

	called := storeMock.UpdatePermissionCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.permissionID, call.PermissionID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePermission_DeletePermission(t *testing.T) {
	t.Parallel()
	type wants struct {
		id uuid.UUID
	}
	cases := []struct {
		id   uuid.UUID
		want wants
	}{
		{
			id: testutils.FixedUUID(t, 0),
			want: wants{
				id: testutils.FixedUUID(t, 0),
			},
		},
	}

	storeMock := &store.StoreMock{
		DeletePermissionFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeletePermission(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeletePermissionCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PermissionID)
	}
}

func TestManagePermission_PluralDeletePermissions(t *testing.T) {
	t.Parallel()
	type wants struct {
		ids []uuid.UUID
	}
	cases := []struct {
		ids  []uuid.UUID
		want wants
	}{
		{
			ids: []uuid.UUID{
				testutils.FixedUUID(t, 0),
				testutils.FixedUUID(t, 1),
			},
			want: wants{
				ids: []uuid.UUID{
					testutils.FixedUUID(t, 0),
					testutils.FixedUUID(t, 1),
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		PluralDeletePermissionsFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeletePermissions(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeletePermissionsCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.PermissionIDs)
	}
}

func TestManagePermission_FindPermissionByID(t *testing.T) {
	t.Parallel()
	type wants struct {
		id uuid.UUID
	}
	cases := []struct {
		id   uuid.UUID
		want wants
	}{
		{
			id: testutils.FixedUUID(t, 0),
			want: wants{
				id: testutils.FixedUUID(t, 0),
			},
		},
	}
	storeMock := &store.StoreMock{
		FindPermissionByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.Permission, error) {
			return entity.Permission{}, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindPermissionByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindPermissionByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PermissionID)
	}
}

func TestManagePermission_GetPermissions(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePermissionParam
		order parameter.PermissionOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName   string
		whereInCategories []uuid.UUID
		order             parameter.PermissionOrderMethod
		pg                parameter.Pagination
		limit             parameter.Limit
		cursor            parameter.Cursor
		offset            parameter.Offset
		withCount         parameter.WithCount
		want              wants
	}{
		{
			whereSearchName:   "",
			whereInCategories: []uuid.UUID{},
			order:             parameter.PermissionOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePermissionParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PermissionOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName:   "",
			whereInCategories: []uuid.UUID{},
			order:             parameter.PermissionOrderMethodDefault,
			pg:                parameter.NumberedPagination,
			limit:             1,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePermissionParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PermissionOrderMethodDefault,
				np: store.NumberedPaginationParam{
					Valid:  true,
					Limit:  entity.Int{Int64: 1},
					Offset: entity.Int{Int64: 0},
				},
				cp: store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName:   "",
			whereInCategories: []uuid.UUID{},
			order:             parameter.PermissionOrderMethodDefault,
			pg:                parameter.NumberedPagination,
			limit:             1,
			cursor:            "",
			offset:            1,
			withCount:         true,
			want: wants{
				where: parameter.WherePermissionParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PermissionOrderMethodDefault,
				np: store.NumberedPaginationParam{
					Valid:  true,
					Limit:  entity.Int{Int64: 1},
					Offset: entity.Int{Int64: 1},
				},
				cp: store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName:   "",
			whereInCategories: []uuid.UUID{},
			order:             parameter.PermissionOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         false,
			want: wants{
				where: parameter.WherePermissionParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PermissionOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName:   "search",
			whereInCategories: []uuid.UUID{},
			order:             parameter.PermissionOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePermissionParam{
					WhereLikeName: true,
					SearchName:    "search",
					InCategories:  []uuid.UUID{},
				},
				order: parameter.PermissionOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			whereInCategories: []uuid.UUID{
				testutils.FixedUUID(t, 0),
				testutils.FixedUUID(t, 1),
			},
			order:     parameter.PermissionOrderMethodDefault,
			pg:        parameter.NonePagination,
			limit:     0,
			cursor:    "",
			offset:    0,
			withCount: true,
			want: wants{
				where: parameter.WherePermissionParam{
					WhereInCategory: true,
					InCategories: []uuid.UUID{
						testutils.FixedUUID(t, 0),
						testutils.FixedUUID(t, 1),
					},
				},
				order: parameter.PermissionOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "search",
			whereInCategories: []uuid.UUID{
				testutils.FixedUUID(t, 0),
				testutils.FixedUUID(t, 1),
			},
			order:     parameter.PermissionOrderMethodDefault,
			pg:        parameter.NonePagination,
			limit:     0,
			cursor:    "",
			offset:    0,
			withCount: true,
			want: wants{
				where: parameter.WherePermissionParam{
					WhereLikeName:   true,
					SearchName:      "search",
					WhereInCategory: true,
					InCategories: []uuid.UUID{
						testutils.FixedUUID(t, 0),
						testutils.FixedUUID(t, 1),
					},
				},
				order: parameter.PermissionOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetPermissionsFunc: func(
			_ context.Context,
			_ parameter.WherePermissionParam,
			_ parameter.PermissionOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.Permission], error) {
			return store.ListResult[entity.Permission]{}, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPermissions(
			ctx, c.whereSearchName, c.whereInCategories, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetPermissionsCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
		assert.Equal(t, c.want.order, call.Order)
		assert.Equal(t, c.want.np, call.Np)
		assert.Equal(t, c.want.cp, call.Cp)
		assert.Equal(t, c.want.wc, call.Wc)
	}
}

func TestManagePermission_GetPermissionsCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePermissionParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WherePermissionParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WherePermissionParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountPermissionsFunc: func(_ context.Context, _ parameter.WherePermissionParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermission{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPermissionsCount(ctx, c.whereSearchName, c.want.where.InCategories)
		assert.NoError(t, err)
	}

	called := storeMock.CountPermissionsCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

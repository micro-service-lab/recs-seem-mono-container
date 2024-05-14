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

func TestManagePermissionCategory_CreatePermissionCategory(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreatePermissionCategoryParam
	}
	cases := []struct {
		name        string
		key         string
		description string
		want        wants
	}{
		{
			name:        "name",
			key:         "key",
			description: "description",
			want: wants{
				param: parameter.CreatePermissionCategoryParam{
					Name:        "name",
					Key:         "key",
					Description: "description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreatePermissionCategoryFunc: func(
			_ context.Context, p parameter.CreatePermissionCategoryParam,
		) (entity.PermissionCategory, error) {
			return entity.PermissionCategory{
				PermissionCategoryID: uuid.New(),
				Name:                 p.Name,
				Key:                  p.Key,
				Description:          p.Description,
			}, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePermissionCategory(ctx, c.name, c.key, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePermissionCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePermissionCategory_CreatePermissionCategories(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreatePermissionCategoryParam
	}
	cases := []struct {
		param []parameter.CreatePermissionCategoryParam
		want  wants
	}{
		{
			param: []parameter.CreatePermissionCategoryParam{
				{Name: "name1", Key: "key1", Description: "description1"},
				{Name: "name2", Key: "key2", Description: "description2"},
			},
			want: wants{
				params: []parameter.CreatePermissionCategoryParam{
					{Name: "name1", Key: "key1", Description: "description1"},
					{Name: "name2", Key: "key2", Description: "description2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreatePermissionCategoriesFunc: func(
			_ context.Context, _ []parameter.CreatePermissionCategoryParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePermissionCategories(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePermissionCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManagePermissionCategory_UpdatePermissionCategory(t *testing.T) {
	t.Parallel()
	type wants struct {
		permissionCategoryID uuid.UUID
		param                parameter.UpdatePermissionCategoryParams
	}
	cases := []struct {
		id          uuid.UUID
		name        string
		key         string
		description string
		want        wants
	}{
		{
			id:          testutils.FixedUUID(t, 0),
			name:        "update name",
			key:         "update key",
			description: "update description",
			want: wants{
				permissionCategoryID: testutils.FixedUUID(t, 0),
				param: parameter.UpdatePermissionCategoryParams{
					Name:        "update name",
					Key:         "update key",
					Description: "update description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdatePermissionCategoryFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdatePermissionCategoryParams,
		) (entity.PermissionCategory, error) {
			return entity.PermissionCategory{}, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdatePermissionCategory(ctx, c.id, c.name, c.key, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.UpdatePermissionCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.permissionCategoryID, call.PermissionCategoryID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePermissionCategory_DeletePermissionCategory(t *testing.T) {
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
		DeletePermissionCategoryFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeletePermissionCategory(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeletePermissionCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PermissionCategoryID)
	}
}

func TestManagePermissionCategory_PluralDeletePermissionCategories(t *testing.T) {
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
		PluralDeletePermissionCategoriesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeletePermissionCategories(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeletePermissionCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.PermissionCategoryIDs)
	}
}

func TestManagePermissionCategory_FindPermissionCategoryByID(t *testing.T) {
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
		FindPermissionCategoryByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.PermissionCategory, error) {
			return entity.PermissionCategory{}, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindPermissionCategoryByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindPermissionCategoryByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PermissionCategoryID)
	}
}

func TestManagePermissionCategory_GetPermissionCategories(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePermissionCategoryParam
		order parameter.PermissionCategoryOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.PermissionCategoryOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.PermissionCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePermissionCategoryParam{},
				order: parameter.PermissionCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.PermissionCategoryOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePermissionCategoryParam{},
				order: parameter.PermissionCategoryOrderMethodDefault,
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
			whereSearchName: "",
			order:           parameter.PermissionCategoryOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WherePermissionCategoryParam{},
				order: parameter.PermissionCategoryOrderMethodDefault,
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
			whereSearchName: "",
			order:           parameter.PermissionCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WherePermissionCategoryParam{},
				order: parameter.PermissionCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.PermissionCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePermissionCategoryParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.PermissionCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetPermissionCategoriesFunc: func(
			_ context.Context,
			_ parameter.WherePermissionCategoryParam,
			_ parameter.PermissionCategoryOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.PermissionCategory], error) {
			return store.ListResult[entity.PermissionCategory]{}, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPermissionCategories(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetPermissionCategoriesCalls()
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

func TestManagePermissionCategory_GetPermissionCategoriesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePermissionCategoryParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WherePermissionCategoryParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WherePermissionCategoryParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountPermissionCategoriesFunc: func(_ context.Context, _ parameter.WherePermissionCategoryParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePermissionCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPermissionCategoriesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountPermissionCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

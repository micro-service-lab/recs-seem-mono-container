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

func TestManageRole_CreateRole(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateRoleParam
	}
	cases := []struct {
		name        string
		description string
		want        wants
	}{
		{
			name:        "name",
			description: "description",
			want: wants{
				param: parameter.CreateRoleParam{
					Name:        "name",
					Description: "description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateRoleFunc: func(
			_ context.Context, p parameter.CreateRoleParam,
		) (entity.Role, error) {
			return entity.Role{
				RoleID:      uuid.New(),
				Name:        p.Name,
				Description: p.Description,
			}, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateRole(ctx, c.name, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.CreateRoleCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageRole_CreateRoles(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateRoleParam
	}
	cases := []struct {
		param []parameter.CreateRoleParam
		want  wants
	}{
		{
			param: []parameter.CreateRoleParam{
				{Name: "name1", Description: "description1"},
				{Name: "name2", Description: "description2"},
			},
			want: wants{
				params: []parameter.CreateRoleParam{
					{Name: "name1", Description: "description1"},
					{Name: "name2", Description: "description2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateRolesFunc: func(
			_ context.Context, _ []parameter.CreateRoleParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateRoles(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateRolesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageRole_UpdateRole(t *testing.T) {
	t.Parallel()
	type wants struct {
		eventTypeID uuid.UUID
		param       parameter.UpdateRoleParams
	}
	cases := []struct {
		id          uuid.UUID
		name        string
		description string
		want        wants
	}{
		{
			id:          testutils.FixedUUID(t, 0),
			name:        "update name",
			description: "update description",
			want: wants{
				eventTypeID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateRoleParams{
					Name:        "update name",
					Description: "update description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateRoleFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateRoleParams,
		) (entity.Role, error) {
			return entity.Role{}, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateRole(ctx, c.id, c.name, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateRoleCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.eventTypeID, call.RoleID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageRole_DeleteRole(t *testing.T) {
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
		DeleteRoleFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeleteRole(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteRoleCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.RoleID)
	}
}

func TestManageRole_PluralDeleteRoles(t *testing.T) {
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
		PluralDeleteRolesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeleteRoles(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteRolesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.RoleIDs)
	}
}

func TestManageRole_FindRoleByID(t *testing.T) {
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
		FindRoleByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.Role, error) {
			return entity.Role{}, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindRoleByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindRoleByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.RoleID)
	}
}

func TestManageRole_GetRoles(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereRoleParam
		order parameter.RoleOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.RoleOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.RoleOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRoleParam{},
				order: parameter.RoleOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.RoleOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRoleParam{},
				order: parameter.RoleOrderMethodDefault,
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
			order:           parameter.RoleOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereRoleParam{},
				order: parameter.RoleOrderMethodDefault,
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
			order:           parameter.RoleOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereRoleParam{},
				order: parameter.RoleOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.RoleOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRoleParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.RoleOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetRolesFunc: func(
			_ context.Context,
			_ parameter.WhereRoleParam,
			_ parameter.RoleOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.Role], error) {
			return store.ListResult[entity.Role]{}, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetRoles(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetRolesCalls()
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

func TestManageRole_GetRolesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereRoleParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereRoleParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereRoleParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountRolesFunc: func(_ context.Context, _ parameter.WhereRoleParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRole{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetRolesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountRolesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

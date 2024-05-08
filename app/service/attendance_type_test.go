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

func TestManageAttendanceType_CreateAttendanceType(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateAttendanceTypeParam
	}
	cases := []struct {
		name  string
		key   string
		color string
		want  wants
	}{
		{
			name:  "name",
			key:   "key",
			color: "color",
			want: wants{
				param: parameter.CreateAttendanceTypeParam{
					Name:  "name",
					Key:   "key",
					Color: "color",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateAttendanceTypeFunc: func(
			_ context.Context, p parameter.CreateAttendanceTypeParam,
		) (entity.AttendanceType, error) {
			return entity.AttendanceType{
				AttendanceTypeID: uuid.New(),
				Name:             p.Name,
				Key:              p.Key,
				Color:            p.Color,
			}, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateAttendanceType(ctx, c.name, c.key, c.color)
		assert.NoError(t, err)
	}

	called := storeMock.CreateAttendanceTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageAttendanceType_CreateAttendanceTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateAttendanceTypeParam
	}
	cases := []struct {
		param []parameter.CreateAttendanceTypeParam
		want  wants
	}{
		{
			param: []parameter.CreateAttendanceTypeParam{
				{Name: "name1", Key: "key1", Color: "color1"},
				{Name: "name2", Key: "key2", Color: "color2"},
			},
			want: wants{
				params: []parameter.CreateAttendanceTypeParam{
					{Name: "name1", Key: "key1", Color: "color1"},
					{Name: "name2", Key: "key2", Color: "color2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateAttendanceTypesFunc: func(
			_ context.Context, _ []parameter.CreateAttendanceTypeParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateAttendanceTypes(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateAttendanceTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageAttendanceType_UpdateAttendanceType(t *testing.T) {
	t.Parallel()
	type wants struct {
		attendTypeID uuid.UUID
		param        parameter.UpdateAttendanceTypeParams
	}
	cases := []struct {
		id    uuid.UUID
		name  string
		key   string
		color string
		want  wants
	}{
		{
			id:    testutils.FixedUUID(t, 0),
			name:  "update name",
			key:   "update key",
			color: "update color",
			want: wants{
				attendTypeID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateAttendanceTypeParams{
					Name:  "update name",
					Key:   "update key",
					Color: "update color",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateAttendanceTypeFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateAttendanceTypeParams,
		) (entity.AttendanceType, error) {
			return entity.AttendanceType{}, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateAttendanceType(ctx, c.id, c.name, c.key, c.color)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateAttendanceTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.attendTypeID, call.AttendanceTypeID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageAttendanceType_DeleteAttendanceType(t *testing.T) {
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
		DeleteAttendanceTypeFunc: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.DeleteAttendanceType(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteAttendanceTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.AttendanceTypeID)
	}
}

func TestManageAttendanceType_PluralDeleteAttendanceTypes(t *testing.T) {
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
		PluralDeleteAttendanceTypesFunc: func(_ context.Context, _ []uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.PluralDeleteAttendanceTypes(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteAttendanceTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.AttendanceTypeIDs)
	}
}

func TestManageAttendanceType_FindAttendanceTypeByID(t *testing.T) {
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
		FindAttendanceTypeByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.AttendanceType, error) {
			return entity.AttendanceType{}, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindAttendanceTypeByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindAttendanceTypeByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.AttendanceTypeID)
	}
}

func TestManageAttendanceType_GetAttendanceTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereAttendanceTypeParam
		order parameter.AttendanceTypeOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.AttendanceTypeOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.AttendanceTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendanceTypeParam{},
				order: parameter.AttendanceTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.AttendanceTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendanceTypeParam{},
				order: parameter.AttendanceTypeOrderMethodDefault,
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
			order:           parameter.AttendanceTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendanceTypeParam{},
				order: parameter.AttendanceTypeOrderMethodDefault,
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
			order:           parameter.AttendanceTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereAttendanceTypeParam{},
				order: parameter.AttendanceTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.AttendanceTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendanceTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.AttendanceTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetAttendanceTypesFunc: func(
			_ context.Context,
			_ parameter.WhereAttendanceTypeParam,
			_ parameter.AttendanceTypeOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.AttendanceType], error) {
			return store.ListResult[entity.AttendanceType]{}, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetAttendanceTypes(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetAttendanceTypesCalls()
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

func TestManageAttendanceType_GetAttendanceTypesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereAttendanceTypeParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereAttendanceTypeParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereAttendanceTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountAttendanceTypesFunc: func(_ context.Context, _ parameter.WhereAttendanceTypeParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendanceType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetAttendanceTypesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountAttendanceTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

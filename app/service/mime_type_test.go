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

func TestManageMimeType_CreateMimeType(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateMimeTypeParam
	}
	cases := []struct {
		name string
		key  string
		kind string
		want wants
	}{
		{
			name: "name",
			key:  "key",
			kind: "kind",
			want: wants{
				param: parameter.CreateMimeTypeParam{
					Name: "name",
					Key:  "key",
					Kind: "kind",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateMimeTypeFunc: func(
			_ context.Context, p parameter.CreateMimeTypeParam,
		) (entity.MimeType, error) {
			return entity.MimeType{
				MimeTypeID: uuid.New(),
				Name:       p.Name,
				Key:        p.Key,
				Kind:       p.Kind,
			}, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateMimeType(ctx, c.name, c.key, c.kind)
		assert.NoError(t, err)
	}

	called := storeMock.CreateMimeTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageMimeType_CreateMimeTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateMimeTypeParam
	}
	cases := []struct {
		param []parameter.CreateMimeTypeParam
		want  wants
	}{
		{
			param: []parameter.CreateMimeTypeParam{
				{Name: "name1", Key: "key1", Kind: "kind1"},
				{Name: "name2", Key: "key2", Kind: "kind2"},
			},
			want: wants{
				params: []parameter.CreateMimeTypeParam{
					{Name: "name1", Key: "key1", Kind: "kind1"},
					{Name: "name2", Key: "key2", Kind: "kind2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateMimeTypesFunc: func(
			_ context.Context, _ []parameter.CreateMimeTypeParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateMimeTypes(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateMimeTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageMimeType_UpdateMimeType(t *testing.T) {
	t.Parallel()
	type wants struct {
		attendTypeID uuid.UUID
		param        parameter.UpdateMimeTypeParams
	}
	cases := []struct {
		id   uuid.UUID
		name string
		key  string
		kind string
		want wants
	}{
		{
			id:   testutils.FixedUUID(t, 0),
			name: "update name",
			key:  "update key",
			kind: "update kind",
			want: wants{
				attendTypeID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateMimeTypeParams{
					Name: "update name",
					Key:  "update key",
					Kind: "update kind",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateMimeTypeFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateMimeTypeParams,
		) (entity.MimeType, error) {
			return entity.MimeType{}, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateMimeType(ctx, c.id, c.name, c.key, c.kind)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateMimeTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.attendTypeID, call.MimeTypeID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageMimeType_DeleteMimeType(t *testing.T) {
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
		DeleteMimeTypeFunc: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.DeleteMimeType(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteMimeTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.MimeTypeID)
	}
}

func TestManageMimeType_PluralDeleteMimeTypes(t *testing.T) {
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
		PluralDeleteMimeTypesFunc: func(_ context.Context, _ []uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.PluralDeleteMimeTypes(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteMimeTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.MimeTypeIDs)
	}
}

func TestManageMimeType_FindMimeTypeByID(t *testing.T) {
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
		FindMimeTypeByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.MimeType, error) {
			return entity.MimeType{}, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindMimeTypeByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindMimeTypeByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.MimeTypeID)
	}
}

func TestManageMimeType_GetMimeTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereMimeTypeParam
		order parameter.MimeTypeOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.MimeTypeOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.MimeTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereMimeTypeParam{},
				order: parameter.MimeTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.MimeTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereMimeTypeParam{},
				order: parameter.MimeTypeOrderMethodDefault,
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
			order:           parameter.MimeTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereMimeTypeParam{},
				order: parameter.MimeTypeOrderMethodDefault,
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
			order:           parameter.MimeTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereMimeTypeParam{},
				order: parameter.MimeTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.MimeTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereMimeTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.MimeTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetMimeTypesFunc: func(
			_ context.Context,
			_ parameter.WhereMimeTypeParam,
			_ parameter.MimeTypeOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.MimeType], error) {
			return store.ListResult[entity.MimeType]{}, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetMimeTypes(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetMimeTypesCalls()
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

func TestManageMimeType_GetMimeTypesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereMimeTypeParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereMimeTypeParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereMimeTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountMimeTypesFunc: func(_ context.Context, _ parameter.WhereMimeTypeParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageMimeType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetMimeTypesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountMimeTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

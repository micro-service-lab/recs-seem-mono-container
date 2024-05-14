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

func TestManageRecordType_CreateRecordType(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateRecordTypeParam
	}
	cases := []struct {
		name string
		key  string
		want wants
	}{
		{
			name: "name",
			key:  "key",
			want: wants{
				param: parameter.CreateRecordTypeParam{
					Name: "name",
					Key:  "key",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateRecordTypeFunc: func(
			_ context.Context, p parameter.CreateRecordTypeParam,
		) (entity.RecordType, error) {
			return entity.RecordType{
				RecordTypeID: uuid.New(),
				Name:         p.Name,
				Key:          p.Key,
			}, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateRecordType(ctx, c.name, c.key)
		assert.NoError(t, err)
	}

	called := storeMock.CreateRecordTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageRecordType_CreateRecordTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateRecordTypeParam
	}
	cases := []struct {
		param []parameter.CreateRecordTypeParam
		want  wants
	}{
		{
			param: []parameter.CreateRecordTypeParam{
				{Name: "name1", Key: "key1"},
				{Name: "name2", Key: "key2"},
			},
			want: wants{
				params: []parameter.CreateRecordTypeParam{
					{Name: "name1", Key: "key1"},
					{Name: "name2", Key: "key2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateRecordTypesFunc: func(
			_ context.Context, _ []parameter.CreateRecordTypeParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateRecordTypes(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateRecordTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageRecordType_UpdateRecordType(t *testing.T) {
	t.Parallel()
	type wants struct {
		attendTypeID uuid.UUID
		param        parameter.UpdateRecordTypeParams
	}
	cases := []struct {
		id   uuid.UUID
		name string
		key  string
		want wants
	}{
		{
			id:   testutils.FixedUUID(t, 0),
			name: "update name",
			key:  "update key",
			want: wants{
				attendTypeID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateRecordTypeParams{
					Name: "update name",
					Key:  "update key",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateRecordTypeFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateRecordTypeParams,
		) (entity.RecordType, error) {
			return entity.RecordType{}, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateRecordType(ctx, c.id, c.name, c.key)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateRecordTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.attendTypeID, call.RecordTypeID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageRecordType_DeleteRecordType(t *testing.T) {
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
		DeleteRecordTypeFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeleteRecordType(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteRecordTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.RecordTypeID)
	}
}

func TestManageRecordType_PluralDeleteRecordTypes(t *testing.T) {
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
		PluralDeleteRecordTypesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeleteRecordTypes(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteRecordTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.RecordTypeIDs)
	}
}

func TestManageRecordType_FindRecordTypeByID(t *testing.T) {
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
		FindRecordTypeByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.RecordType, error) {
			return entity.RecordType{}, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindRecordTypeByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindRecordTypeByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.RecordTypeID)
	}
}

func TestManageRecordType_GetRecordTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereRecordTypeParam
		order parameter.RecordTypeOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.RecordTypeOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.RecordTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRecordTypeParam{},
				order: parameter.RecordTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.RecordTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRecordTypeParam{},
				order: parameter.RecordTypeOrderMethodDefault,
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
			order:           parameter.RecordTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereRecordTypeParam{},
				order: parameter.RecordTypeOrderMethodDefault,
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
			order:           parameter.RecordTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereRecordTypeParam{},
				order: parameter.RecordTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.RecordTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereRecordTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.RecordTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetRecordTypesFunc: func(
			_ context.Context,
			_ parameter.WhereRecordTypeParam,
			_ parameter.RecordTypeOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.RecordType], error) {
			return store.ListResult[entity.RecordType]{}, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetRecordTypes(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetRecordTypesCalls()
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

func TestManageRecordType_GetRecordTypesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereRecordTypeParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereRecordTypeParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereRecordTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountRecordTypesFunc: func(_ context.Context, _ parameter.WhereRecordTypeParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageRecordType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetRecordTypesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountRecordTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

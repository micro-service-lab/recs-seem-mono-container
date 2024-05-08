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

func TestManageEventType_CreateEventType(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateEventTypeParam
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
				param: parameter.CreateEventTypeParam{
					Name:  "name",
					Key:   "key",
					Color: "color",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateEventTypeFunc: func(
			_ context.Context, p parameter.CreateEventTypeParam,
		) (entity.EventType, error) {
			return entity.EventType{
				EventTypeID: uuid.New(),
				Name:        p.Name,
				Key:         p.Key,
				Color:       p.Color,
			}, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateEventType(ctx, c.name, c.key, c.color)
		assert.NoError(t, err)
	}

	called := storeMock.CreateEventTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageEventType_CreateEventTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateEventTypeParam
	}
	cases := []struct {
		param []parameter.CreateEventTypeParam
		want  wants
	}{
		{
			param: []parameter.CreateEventTypeParam{
				{Name: "name1", Key: "key1", Color: "color1"},
				{Name: "name2", Key: "key2", Color: "color2"},
			},
			want: wants{
				params: []parameter.CreateEventTypeParam{
					{Name: "name1", Key: "key1", Color: "color1"},
					{Name: "name2", Key: "key2", Color: "color2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateEventTypesFunc: func(
			_ context.Context, _ []parameter.CreateEventTypeParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateEventTypes(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateEventTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageEventType_UpdateEventType(t *testing.T) {
	t.Parallel()
	type wants struct {
		eventTypeID uuid.UUID
		param       parameter.UpdateEventTypeParams
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
				eventTypeID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateEventTypeParams{
					Name:  "update name",
					Key:   "update key",
					Color: "update color",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateEventTypeFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateEventTypeParams,
		) (entity.EventType, error) {
			return entity.EventType{}, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateEventType(ctx, c.id, c.name, c.key, c.color)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateEventTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.eventTypeID, call.EventTypeID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageEventType_DeleteEventType(t *testing.T) {
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
		DeleteEventTypeFunc: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.DeleteEventType(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteEventTypeCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.EventTypeID)
	}
}

func TestManageEventType_PluralDeleteEventTypes(t *testing.T) {
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
		PluralDeleteEventTypesFunc: func(_ context.Context, _ []uuid.UUID) error {
			return nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		err := s.PluralDeleteEventTypes(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteEventTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.EventTypeIDs)
	}
}

func TestManageEventType_FindEventTypeByID(t *testing.T) {
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
		FindEventTypeByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.EventType, error) {
			return entity.EventType{}, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindEventTypeByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindEventTypeByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.EventTypeID)
	}
}

func TestManageEventType_GetEventTypes(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereEventTypeParam
		order parameter.EventTypeOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.EventTypeOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.EventTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereEventTypeParam{},
				order: parameter.EventTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.EventTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereEventTypeParam{},
				order: parameter.EventTypeOrderMethodDefault,
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
			order:           parameter.EventTypeOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereEventTypeParam{},
				order: parameter.EventTypeOrderMethodDefault,
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
			order:           parameter.EventTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereEventTypeParam{},
				order: parameter.EventTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.EventTypeOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereEventTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.EventTypeOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetEventTypesFunc: func(
			_ context.Context,
			_ parameter.WhereEventTypeParam,
			_ parameter.EventTypeOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.EventType], error) {
			return store.ListResult[entity.EventType]{}, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetEventTypes(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetEventTypesCalls()
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

func TestManageEventType_GetEventTypesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereEventTypeParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereEventTypeParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereEventTypeParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountEventTypesFunc: func(_ context.Context, _ parameter.WhereEventTypeParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageEventType{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetEventTypesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountEventTypesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

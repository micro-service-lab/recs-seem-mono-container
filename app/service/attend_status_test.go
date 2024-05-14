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

func TestManageAttendStatus_CreateAttendStatus(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreateAttendStatusParam
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
				param: parameter.CreateAttendStatusParam{
					Name: "name",
					Key:  "key",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreateAttendStatusFunc: func(
			_ context.Context, p parameter.CreateAttendStatusParam,
		) (entity.AttendStatus, error) {
			return entity.AttendStatus{
				AttendStatusID: uuid.New(),
				Name:           p.Name,
				Key:            p.Key,
			}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateAttendStatus(ctx, c.name, c.key)
		assert.NoError(t, err)
	}

	called := storeMock.CreateAttendStatusCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageAttendStatus_CreateAttendStatuses(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreateAttendStatusParam
	}
	cases := []struct {
		param []parameter.CreateAttendStatusParam
		want  wants
	}{
		{
			param: []parameter.CreateAttendStatusParam{
				{Name: "name1", Key: "key1"},
				{Name: "name2", Key: "key2"},
			},
			want: wants{
				params: []parameter.CreateAttendStatusParam{
					{Name: "name1", Key: "key1"},
					{Name: "name2", Key: "key2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreateAttendStatusesFunc: func(
			_ context.Context, _ []parameter.CreateAttendStatusParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreateAttendStatuses(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreateAttendStatusesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManageAttendStatus_UpdateAttendStatus(t *testing.T) {
	t.Parallel()
	type wants struct {
		attendStatusID uuid.UUID
		param          parameter.UpdateAttendStatusParams
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
				attendStatusID: testutils.FixedUUID(t, 0),
				param: parameter.UpdateAttendStatusParams{
					Name: "update name",
					Key:  "update key",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdateAttendStatusFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdateAttendStatusParams,
		) (entity.AttendStatus, error) {
			return entity.AttendStatus{}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdateAttendStatus(ctx, c.id, c.name, c.key)
		assert.NoError(t, err)
	}

	called := storeMock.UpdateAttendStatusCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.attendStatusID, call.AttendStatusID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManageAttendStatus_DeleteAttendStatus(t *testing.T) {
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
		DeleteAttendStatusFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeleteAttendStatus(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeleteAttendStatusCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.AttendStatusID)
	}
}

func TestManageAttendStatus_PluralDeleteAttendStatuses(t *testing.T) {
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
		PluralDeleteAttendStatusesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeleteAttendStatuses(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeleteAttendStatusesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.AttendStatusIDs)
	}
}

func TestManageAttendStatus_FindAttendStatusByID(t *testing.T) {
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
		FindAttendStatusByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.AttendStatus, error) {
			return entity.AttendStatus{}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindAttendStatusByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindAttendStatusByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.AttendStatusID)
	}
}

func TestManageAttendStatus_GetAttendStatuses(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereAttendStatusParam
		order parameter.AttendStatusOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.AttendStatusOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.AttendStatusOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendStatusParam{},
				order: parameter.AttendStatusOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.AttendStatusOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendStatusParam{},
				order: parameter.AttendStatusOrderMethodDefault,
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
			order:           parameter.AttendStatusOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendStatusParam{},
				order: parameter.AttendStatusOrderMethodDefault,
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
			order:           parameter.AttendStatusOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WhereAttendStatusParam{},
				order: parameter.AttendStatusOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.AttendStatusOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WhereAttendStatusParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.AttendStatusOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetAttendStatusesFunc: func(
			_ context.Context,
			_ parameter.WhereAttendStatusParam,
			_ parameter.AttendStatusOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.AttendStatus], error) {
			return store.ListResult[entity.AttendStatus]{}, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetAttendStatuses(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetAttendStatusesCalls()
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

func TestManageAttendStatus_GetAttendStatusesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WhereAttendStatusParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WhereAttendStatusParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WhereAttendStatusParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountAttendStatusesFunc: func(_ context.Context, _ parameter.WhereAttendStatusParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManageAttendStatus{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetAttendStatusesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountAttendStatusesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

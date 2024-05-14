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

func TestManagePolicyCategory_CreatePolicyCategory(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreatePolicyCategoryParam
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
				param: parameter.CreatePolicyCategoryParam{
					Name:        "name",
					Key:         "key",
					Description: "description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreatePolicyCategoryFunc: func(
			_ context.Context, p parameter.CreatePolicyCategoryParam,
		) (entity.PolicyCategory, error) {
			return entity.PolicyCategory{
				PolicyCategoryID: uuid.New(),
				Name:             p.Name,
				Key:              p.Key,
				Description:      p.Description,
			}, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePolicyCategory(ctx, c.name, c.key, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePolicyCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePolicyCategory_CreatePolicyCategories(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreatePolicyCategoryParam
	}
	cases := []struct {
		param []parameter.CreatePolicyCategoryParam
		want  wants
	}{
		{
			param: []parameter.CreatePolicyCategoryParam{
				{Name: "name1", Key: "key1", Description: "description1"},
				{Name: "name2", Key: "key2", Description: "description2"},
			},
			want: wants{
				params: []parameter.CreatePolicyCategoryParam{
					{Name: "name1", Key: "key1", Description: "description1"},
					{Name: "name2", Key: "key2", Description: "description2"},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreatePolicyCategoriesFunc: func(
			_ context.Context, _ []parameter.CreatePolicyCategoryParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePolicyCategories(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePolicyCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManagePolicyCategory_UpdatePolicyCategory(t *testing.T) {
	t.Parallel()
	type wants struct {
		policyCategoryID uuid.UUID
		param            parameter.UpdatePolicyCategoryParams
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
				policyCategoryID: testutils.FixedUUID(t, 0),
				param: parameter.UpdatePolicyCategoryParams{
					Name:        "update name",
					Key:         "update key",
					Description: "update description",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdatePolicyCategoryFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdatePolicyCategoryParams,
		) (entity.PolicyCategory, error) {
			return entity.PolicyCategory{}, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdatePolicyCategory(ctx, c.id, c.name, c.key, c.description)
		assert.NoError(t, err)
	}

	called := storeMock.UpdatePolicyCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.policyCategoryID, call.PolicyCategoryID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePolicyCategory_DeletePolicyCategory(t *testing.T) {
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
		DeletePolicyCategoryFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeletePolicyCategory(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeletePolicyCategoryCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PolicyCategoryID)
	}
}

func TestManagePolicyCategory_PluralDeletePolicyCategories(t *testing.T) {
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
		PluralDeletePolicyCategoriesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeletePolicyCategories(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeletePolicyCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.PolicyCategoryIDs)
	}
}

func TestManagePolicyCategory_FindPolicyCategoryByID(t *testing.T) {
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
		FindPolicyCategoryByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.PolicyCategory, error) {
			return entity.PolicyCategory{}, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindPolicyCategoryByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindPolicyCategoryByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PolicyCategoryID)
	}
}

func TestManagePolicyCategory_GetPolicyCategories(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePolicyCategoryParam
		order parameter.PolicyCategoryOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName string
		order           parameter.PolicyCategoryOrderMethod
		pg              parameter.Pagination
		limit           parameter.Limit
		cursor          parameter.Cursor
		offset          parameter.Offset
		withCount       parameter.WithCount
		want            wants
	}{
		{
			whereSearchName: "",
			order:           parameter.PolicyCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePolicyCategoryParam{},
				order: parameter.PolicyCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
		{
			whereSearchName: "",
			order:           parameter.PolicyCategoryOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePolicyCategoryParam{},
				order: parameter.PolicyCategoryOrderMethodDefault,
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
			order:           parameter.PolicyCategoryOrderMethodDefault,
			pg:              parameter.NumberedPagination,
			limit:           1,
			cursor:          "",
			offset:          1,
			withCount:       true,
			want: wants{
				where: parameter.WherePolicyCategoryParam{},
				order: parameter.PolicyCategoryOrderMethodDefault,
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
			order:           parameter.PolicyCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       false,
			want: wants{
				where: parameter.WherePolicyCategoryParam{},
				order: parameter.PolicyCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: false,
				},
			},
		},
		{
			whereSearchName: "search",
			order:           parameter.PolicyCategoryOrderMethodDefault,
			pg:              parameter.NonePagination,
			limit:           0,
			cursor:          "",
			offset:          0,
			withCount:       true,
			want: wants{
				where: parameter.WherePolicyCategoryParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
				order: parameter.PolicyCategoryOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetPolicyCategoriesFunc: func(
			_ context.Context,
			_ parameter.WherePolicyCategoryParam,
			_ parameter.PolicyCategoryOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.PolicyCategory], error) {
			return store.ListResult[entity.PolicyCategory]{}, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPolicyCategories(
			ctx, c.whereSearchName, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetPolicyCategoriesCalls()
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

func TestManagePolicyCategory_GetPolicyCategoriesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePolicyCategoryParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WherePolicyCategoryParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WherePolicyCategoryParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountPolicyCategoriesFunc: func(_ context.Context, _ parameter.WherePolicyCategoryParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicyCategory{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPolicyCategoriesCount(ctx, c.whereSearchName)
		assert.NoError(t, err)
	}

	called := storeMock.CountPolicyCategoriesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

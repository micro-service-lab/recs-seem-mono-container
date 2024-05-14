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

func TestManagePolicy_CreatePolicy(t *testing.T) {
	t.Parallel()
	type wants struct {
		param parameter.CreatePolicyParam
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
				param: parameter.CreatePolicyParam{
					Name:             "name",
					Key:              "key",
					Description:      "description",
					PolicyCategoryID: testutils.FixedUUID(t, 0),
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CreatePolicyFunc: func(
			_ context.Context, p parameter.CreatePolicyParam,
		) (entity.Policy, error) {
			return entity.Policy{
				PolicyID:         uuid.New(),
				Name:             p.Name,
				Key:              p.Key,
				Description:      p.Description,
				PolicyCategoryID: p.PolicyCategoryID,
			}, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePolicy(ctx, c.name, c.key, c.description, c.categoryID)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePolicyCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePolicy_CreatePolicies(t *testing.T) {
	t.Parallel()
	type wants struct {
		params []parameter.CreatePolicyParam
	}
	cases := []struct {
		param []parameter.CreatePolicyParam
		want  wants
	}{
		{
			param: []parameter.CreatePolicyParam{
				{Name: "name1", Key: "key1", Description: "description1", PolicyCategoryID: testutils.FixedUUID(t, 0)},
				{Name: "name2", Key: "key2", Description: "description2", PolicyCategoryID: testutils.FixedUUID(t, 1)},
			},
			want: wants{
				params: []parameter.CreatePolicyParam{
					{Name: "name1", Key: "key1", Description: "description1", PolicyCategoryID: testutils.FixedUUID(t, 0)},
					{Name: "name2", Key: "key2", Description: "description2", PolicyCategoryID: testutils.FixedUUID(t, 1)},
				},
			},
		},
	}

	storeMock := &store.StoreMock{
		CreatePoliciesFunc: func(
			_ context.Context, _ []parameter.CreatePolicyParam,
		) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.CreatePolicies(ctx, c.param)
		assert.NoError(t, err)
	}

	called := storeMock.CreatePoliciesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.params, call.Params)
	}
}

func TestManagePolicy_UpdatePolicy(t *testing.T) {
	t.Parallel()
	type wants struct {
		policyID uuid.UUID
		param    parameter.UpdatePolicyParams
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
				policyID: testutils.FixedUUID(t, 0),
				param: parameter.UpdatePolicyParams{
					Name:             "update name",
					Key:              "update key",
					Description:      "update description",
					PolicyCategoryID: testutils.FixedUUID(t, 0),
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		UpdatePolicyFunc: func(
			_ context.Context, _ uuid.UUID, _ parameter.UpdatePolicyParams,
		) (entity.Policy, error) {
			return entity.Policy{}, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.UpdatePolicy(ctx, c.id, c.name, c.key, c.description, c.categoryID)
		assert.NoError(t, err)
	}

	called := storeMock.UpdatePolicyCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.policyID, call.PolicyID)
		assert.Equal(t, c.want.param, call.Param)
	}
}

func TestManagePolicy_DeletePolicy(t *testing.T) {
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
		DeletePolicyFunc: func(_ context.Context, _ uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.DeletePolicy(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.DeletePolicyCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PolicyID)
	}
}

func TestManagePolicy_PluralDeletePolicies(t *testing.T) {
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
		PluralDeletePoliciesFunc: func(_ context.Context, _ []uuid.UUID) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()

	for _, c := range cases {
		_, err := s.PluralDeletePolicies(ctx, c.ids)
		assert.NoError(t, err)
	}

	called := storeMock.PluralDeletePoliciesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.ids, call.PolicyIDs)
	}
}

func TestManagePolicy_FindPolicyByID(t *testing.T) {
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
		FindPolicyByIDFunc: func(_ context.Context, _ uuid.UUID) (entity.Policy, error) {
			return entity.Policy{}, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.FindPolicyByID(ctx, c.id)
		assert.NoError(t, err)
	}

	called := storeMock.FindPolicyByIDCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.id, call.PolicyID)
	}
}

func TestManagePolicy_GetPolicies(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePolicyParam
		order parameter.PolicyOrderMethod
		np    store.NumberedPaginationParam
		cp    store.CursorPaginationParam
		wc    store.WithCountParam
	}
	cases := []struct {
		whereSearchName   string
		whereInCategories []uuid.UUID
		order             parameter.PolicyOrderMethod
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
			order:             parameter.PolicyOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePolicyParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:             parameter.PolicyOrderMethodDefault,
			pg:                parameter.NumberedPagination,
			limit:             1,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePolicyParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:             parameter.PolicyOrderMethodDefault,
			pg:                parameter.NumberedPagination,
			limit:             1,
			cursor:            "",
			offset:            1,
			withCount:         true,
			want: wants{
				where: parameter.WherePolicyParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:             parameter.PolicyOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         false,
			want: wants{
				where: parameter.WherePolicyParam{
					InCategories: []uuid.UUID{},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:             parameter.PolicyOrderMethodDefault,
			pg:                parameter.NonePagination,
			limit:             0,
			cursor:            "",
			offset:            0,
			withCount:         true,
			want: wants{
				where: parameter.WherePolicyParam{
					WhereLikeName: true,
					SearchName:    "search",
					InCategories:  []uuid.UUID{},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:     parameter.PolicyOrderMethodDefault,
			pg:        parameter.NonePagination,
			limit:     0,
			cursor:    "",
			offset:    0,
			withCount: true,
			want: wants{
				where: parameter.WherePolicyParam{
					WhereInCategory: true,
					InCategories: []uuid.UUID{
						testutils.FixedUUID(t, 0),
						testutils.FixedUUID(t, 1),
					},
				},
				order: parameter.PolicyOrderMethodDefault,
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
			order:     parameter.PolicyOrderMethodDefault,
			pg:        parameter.NonePagination,
			limit:     0,
			cursor:    "",
			offset:    0,
			withCount: true,
			want: wants{
				where: parameter.WherePolicyParam{
					WhereLikeName:   true,
					SearchName:      "search",
					WhereInCategory: true,
					InCategories: []uuid.UUID{
						testutils.FixedUUID(t, 0),
						testutils.FixedUUID(t, 1),
					},
				},
				order: parameter.PolicyOrderMethodDefault,
				np:    store.NumberedPaginationParam{},
				cp:    store.CursorPaginationParam{},
				wc: store.WithCountParam{
					Valid: true,
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		GetPoliciesFunc: func(
			_ context.Context,
			_ parameter.WherePolicyParam,
			_ parameter.PolicyOrderMethod,
			_ store.NumberedPaginationParam,
			_ store.CursorPaginationParam,
			_ store.WithCountParam,
		) (store.ListResult[entity.Policy], error) {
			return store.ListResult[entity.Policy]{}, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPolicies(
			ctx, c.whereSearchName, c.whereInCategories, c.order, c.pg, c.limit, c.cursor, c.offset, c.withCount)
		assert.NoError(t, err)
	}

	called := storeMock.GetPoliciesCalls()
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

func TestManagePolicy_GetPoliciesCount(t *testing.T) {
	t.Parallel()
	type wants struct {
		where parameter.WherePolicyParam
	}
	cases := []struct {
		whereSearchName string
		want            wants
	}{
		{
			whereSearchName: "",
			want: wants{
				where: parameter.WherePolicyParam{},
			},
		},
		{
			whereSearchName: "search",
			want: wants{
				where: parameter.WherePolicyParam{
					WhereLikeName: true,
					SearchName:    "search",
				},
			},
		},
	}
	storeMock := &store.StoreMock{
		CountPoliciesFunc: func(_ context.Context, _ parameter.WherePolicyParam) (int64, error) {
			return 0, nil
		},
	}
	s := service.ManagePolicy{
		DB: storeMock,
	}
	ctx := context.Background()
	for _, c := range cases {
		_, err := s.GetPoliciesCount(ctx, c.whereSearchName, c.want.where.InCategories)
		assert.NoError(t, err)
	}

	called := storeMock.CountPoliciesCalls()
	assert.Len(t, called, len(cases))
	for i, call := range called {
		c := cases[i]
		assert.Equal(t, c.want.where, call.Where)
	}
}

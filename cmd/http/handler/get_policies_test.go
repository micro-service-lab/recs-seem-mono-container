package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils/factory"
)

func TestGetPolicies_ServeHTTP(t *testing.T) {
	t.Parallel()
	fdc, err := factory.Generator.NewPolicyCategories(3)
	require.NoError(t, err)
	fd, err := factory.Generator.NewPolicies(10)
	require.NoError(t, err)
	fd, err = fd.WithPolicyCategory(fdc)
	require.NoError(t, err)
	type wants struct {
		resType response.APIResponseType
		data    any
		errAttr response.ApplicationErrorAttributes
	}
	//nolint:lll
	cases := map[string]struct {
		searchName       string
		searchCategories []uuid.UUID
		order            string
		limit            string
		cursor           string
		offset           string
		pagination       string
		withCount        string
		with             []string
		want             wants
	}{
		"simple": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data: fd.ForEntity(),
				},
				errAttr: nil,
			},
		},
		"with count": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data: fd.ForEntity(),
					WithCount: store.WithCountAttribute{
						Count: int64(len(fd)),
					},
				},
				errAttr: nil,
			},
		},
		"with search name": {
			searchName:       fd[0].Name,
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data: fd.FilterByName(fd[0].Name).ForEntity(),
					WithCount: store.WithCountAttribute{
						Count: fd.CountContainsName(fd[0].Name),
					},
				},
				errAttr: nil,
			},
		},
		"maybe no result": {
			searchName:       "no result",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.FilterByName("no result").ForEntity(),
					WithCount:        store.WithCountAttribute{Count: fd.CountContainsName("no result")},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with search categories": {
			searchName:       "",
			searchCategories: []uuid.UUID{fdc[0].PolicyCategoryID},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID}).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID})))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with plural search categories": {
			searchName:       "",
			searchCategories: []uuid.UUID{fdc[0].PolicyCategoryID, fdc[1].PolicyCategoryID},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID, fdc[1].PolicyCategoryID}).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID, fdc[1].PolicyCategoryID})))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with search categories and search name": {
			searchName:       fd[0].Name,
			searchCategories: []uuid.UUID{fdc[0].PolicyCategoryID},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID}).FilterByName(fd[0].Name).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.FilterByPolicyCategories([]uuid.UUID{fdc[0].PolicyCategoryID}).FilterByName(fd[0].Name)))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with order": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "name",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.OrderByNames().ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with reverse order": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "r_name",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.OrderByReverseNames().ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with limit": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "1",
			cursor:           "",
			offset:           "",
			pagination:       "numbered",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.LimitAndOffset(1, 0).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.LimitAndOffset(1, 0)))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with offset": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "1",
			cursor:           "",
			offset:           "1",
			pagination:       "numbered",
			withCount:        "true",
			with:             []string{},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.Policy]{
					Data:             fd.LimitAndOffset(1, 1).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.LimitAndOffset(1, 1)))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with category": {
			searchName:       "",
			searchCategories: []uuid.UUID{},
			order:            "",
			limit:            "",
			cursor:           "",
			offset:           "",
			pagination:       "",
			withCount:        "",
			with:             []string{"category"},
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.PolicyWithCategory]{
					Data: fd.ForEntityWithPolicyCategory(),
				},
				errAttr: nil,
			},
		},
	}
	mockService := &service.ManagerInterfaceMock{
		GetPoliciesFunc: func(
			_ context.Context,
			whereSearchName string,
			whereInCategories []uuid.UUID,
			order parameter.PolicyOrderMethod,
			pg parameter.Pagination,
			limit parameter.Limit,
			_ parameter.Cursor,
			offset parameter.Offset,
			wc parameter.WithCount,
		) (store.ListResult[entity.Policy], error) {
			dd := fd.Copy()
			var wca store.WithCountAttribute
			var cpa store.CursorPaginationAttribute
			if whereSearchName != "" {
				dd = dd.FilterByName(whereSearchName)
			}
			if len(whereInCategories) > 0 {
				dd = dd.FilterByPolicyCategories(whereInCategories)
			}
			switch order {
			case parameter.PolicyOrderMethodName:
				dd = dd.OrderByNames()
			case parameter.PolicyOrderMethodReverseName:
				dd = dd.OrderByReverseNames()
			case parameter.PolicyOrderMethodDefault:
			}
			switch pg {
			case parameter.NumberedPagination:
				dd = dd.LimitAndOffset(int(limit), int(offset))
			case parameter.CursorPagination:
			case parameter.NonePagination:
			}
			if bool(wc) {
				wca = store.WithCountAttribute{
					Count: int64(len(dd)),
				}
			}
			return store.ListResult[entity.Policy]{
				Data:             dd.ForEntity(),
				WithCount:        wca,
				CursorPagination: cpa,
			}, nil
		},
		GetPoliciesWithCategoryFunc: func(
			_ context.Context,
			whereSearchName string,
			whereInCategories []uuid.UUID,
			order parameter.PolicyOrderMethod,
			pg parameter.Pagination,
			limit parameter.Limit,
			_ parameter.Cursor,
			offset parameter.Offset,
			wc parameter.WithCount,
		) (store.ListResult[entity.PolicyWithCategory], error) {
			dd := fd.Copy()
			var wca store.WithCountAttribute
			var cpa store.CursorPaginationAttribute
			if whereSearchName != "" {
				dd = dd.FilterByName(whereSearchName)
			}
			if len(whereInCategories) > 0 {
				dd = dd.FilterByPolicyCategories(whereInCategories)
			}
			switch order {
			case parameter.PolicyOrderMethodName:
				dd = dd.OrderByNames()
			case parameter.PolicyOrderMethodReverseName:
				dd = dd.OrderByReverseNames()
			case parameter.PolicyOrderMethodDefault:
			}
			switch pg {
			case parameter.NumberedPagination:
				dd = dd.LimitAndOffset(int(limit), int(offset))
			case parameter.CursorPagination:
			case parameter.NonePagination:
			}
			if bool(wc) {
				wca = store.WithCountAttribute{
					Count: int64(len(dd)),
				}
			}
			return store.ListResult[entity.PolicyWithCategory]{
				Data:             dd.ForEntityWithPolicyCategory(),
				WithCount:        wca,
				CursorPagination: cpa,
			}, nil
		},
	}
	h := handler.GetPolicies{
		Service: mockService,
	}
	for ni, tc := range cases {
		tcc := tc
		t.Run(ni, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/policies", nil)
			q := r.URL.Query()
			q.Set("search_name", tcc.searchName)
			for _, v := range tcc.searchCategories {
				q.Add("search_categories", v.String())
			}
			q.Set("order", tcc.order)
			q.Set("limit", tcc.limit)
			q.Set("cursor", tcc.cursor)
			q.Set("offset", tcc.offset)
			q.Set("pagination", tcc.pagination)
			q.Set("with_count", tcc.withCount)
			for _, v := range tcc.with {
				q.Add("with", v)
			}
			r.URL.RawQuery = q.Encode()
			h.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			want := response.ApplicationResponse{
				Success:         tcc.want.resType == response.Success,
				Data:            tcc.want.data,
				Code:            tcc.want.resType.Code,
				Message:         tcc.want.resType.Message,
				ErrorAttributes: tcc.want.errAttr,
			}
			wb, err := json.Marshal(want)
			assert.NoError(t, err)
			testutils.AssertResponse(t, resp, tcc.want.resType.StatusCode, wb)
		})
	}
}

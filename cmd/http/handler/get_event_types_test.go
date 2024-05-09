package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetEventTypes_ServeHTTP(t *testing.T) {
	t.Parallel()
	fd, err := factory.Generator.NewEventTypes(10)
	require.NoError(t, err)
	type wants struct {
		resType response.APIResponseType
		data    any
		errAttr response.ApplicationErrorAttributes
	}
	cases := map[string]struct {
		searchName string
		order      string
		limit      string
		cursor     string
		offset     string
		pagination string
		withCount  string
		want       wants
	}{
		"simple": {
			searchName: "",
			order:      "",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data: fd.ForEntity(),
				},
				errAttr: nil,
			},
		},
		"with count": {
			searchName: "",
			order:      "",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data: fd.ForEntity(),
					WithCount: store.WithCountAttribute{
						Count: int64(len(fd)),
					},
				},
				errAttr: nil,
			},
		},
		"with search name": {
			searchName: fd[0].Name,
			order:      "",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data: fd.FilterByName(fd[0].Name).ForEntity(),
					WithCount: store.WithCountAttribute{
						Count: fd.CountContainsName(fd[0].Name),
					},
				},
				errAttr: nil,
			},
		},
		"maybe no result": {
			searchName: "no result",
			order:      "",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data:             fd.FilterByName("no result").ForEntity(),
					WithCount:        store.WithCountAttribute{Count: fd.CountContainsName("no result")},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with order": {
			searchName: "",
			order:      "name",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data:             fd.OrderByNames().ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with reverse order": {
			searchName: "",
			order:      "r_name",
			limit:      "",
			cursor:     "",
			offset:     "",
			pagination: "",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data:             fd.OrderByReverseNames().ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with limit": {
			searchName: "",
			order:      "",
			limit:      "1",
			cursor:     "",
			offset:     "",
			pagination: "numbered",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data:             fd.LimitAndOffset(1, 0).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.LimitAndOffset(1, 0)))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
		"with offset": {
			searchName: "",
			order:      "",
			limit:      "1",
			cursor:     "",
			offset:     "1",
			pagination: "numbered",
			withCount:  "true",
			want: wants{
				resType: response.Success,
				data: store.ListResult[entity.EventType]{
					Data:             fd.LimitAndOffset(1, 1).ForEntity(),
					WithCount:        store.WithCountAttribute{Count: int64(len(fd.LimitAndOffset(1, 1)))},
					CursorPagination: store.CursorPaginationAttribute{},
				},
				errAttr: nil,
			},
		},
	}
	mockService := &service.ManagerInterfaceMock{
		GetEventTypesFunc: func(
			_ context.Context,
			whereSearchName string,
			order parameter.EventTypeOrderMethod,
			pg parameter.Pagination,
			limit parameter.Limit,
			_ parameter.Cursor,
			offset parameter.Offset,
			wc parameter.WithCount,
		) (store.ListResult[entity.EventType], error) {
			dd := fd.Copy()
			var wca store.WithCountAttribute
			var cpa store.CursorPaginationAttribute
			if whereSearchName != "" {
				dd = dd.FilterByName(whereSearchName)
			}
			switch order {
			case parameter.EventTypeOrderMethodName:
				dd = dd.OrderByNames()
			case parameter.EventTypeOrderMethodReverseName:
				dd = dd.OrderByReverseNames()
			case parameter.EventTypeOrderMethodDefault:
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
			return store.ListResult[entity.EventType]{
				Data:             dd.ForEntity(),
				WithCount:        wca,
				CursorPagination: cpa,
			}, nil
		},
	}
	h := handler.GetEventTypes{
		Service: mockService,
	}
	for ni, tc := range cases {
		tcc := tc
		t.Run(ni, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/event_types", nil)
			q := r.URL.Query()
			q.Set("search_name", tcc.searchName)
			q.Set("order", tcc.order)
			q.Set("limit", tcc.limit)
			q.Set("cursor", tcc.cursor)
			q.Set("offset", tcc.offset)
			q.Set("pagination", tcc.pagination)
			q.Set("with_count", tcc.withCount)
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

package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetRoles is a handler for getting roles.
type GetRoles struct {
	Service service.ManagerInterface
}

// GetRolesParam is a parameter for GetRoles.
type GetRolesParam struct {
	SearchName string                    `queryParam:"search_name"`
	Order      parameter.RoleOrderMethod `queryParam:"order"`
	Limit      parameter.Limit           `queryParam:"limit"`
	Offset     parameter.Offset          `queryParam:"offset"`
	Cursor     parameter.Cursor          `queryParam:"cursor"`
	Pagination parameter.Pagination      `queryParam:"pagination"`
	WithCount  parameter.WithCount       `queryParam:"with_count"`
}

var getRolesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.RoleOrderMethodName): parameter.ParseRoleOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):            parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):           parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):          parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):      parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):    parameter.ParseWithCountParam,
}

func (h *GetRoles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetRolesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getRolesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	eventTypes, err := h.Service.GetRoles(
		ctx,
		param.SearchName,
		param.Order,
		param.Pagination,
		param.Limit,
		param.Cursor,
		param.Offset,
		param.WithCount,
	)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, eventTypes, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

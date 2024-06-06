package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetPermissions is a handler for getting permission.
type GetPermissions struct {
	Service service.ManagerInterface
}

// GetPermissionsParam is a parameter for GetPermissions.
type GetPermissionsParam struct {
	SearchName       string                          `queryParam:"search_name"`
	SearchCategories []parameter.EntityID            `queryParam:"search_categories[]"`
	Order            parameter.PermissionOrderMethod `queryParam:"order"`
	Limit            parameter.Limit                 `queryParam:"limit"`
	Offset           parameter.Offset                `queryParam:"offset"`
	Cursor           parameter.Cursor                `queryParam:"cursor"`
	Pagination       parameter.Pagination            `queryParam:"pagination"`
	WithCount        parameter.WithCount             `queryParam:"with_count"`
	With             parameter.PermissionWithParams  `queryParam:"with[]"`
}

var getPermissionsParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PermissionOrderMethodName): parameter.ParsePermissionOrderMethod,
	reflect.TypeOf(parameter.EntityID(uuid.UUID{})):     parameter.ParseEntityIDParam,
	reflect.TypeOf(parameter.Limit(0)):                  parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                 parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):            parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):          parameter.ParseWithCountParam,
	reflect.TypeOf(parameter.PermissionWith{}):          parameter.ParsePermissionWithParam,
}

func (h *GetPermissions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetPermissionsParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getPermissionsParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	inCategories := make([]uuid.UUID, 0, len(param.SearchCategories))
	for _, v := range param.SearchCategories {
		if uuid.UUID(v) != uuid.Nil {
			inCategories = append(inCategories, uuid.UUID(v))
		}
	}
	var permissions any
	switch param.With.Case() {
	case parameter.PermissionWithCaseCategory:
		permissions, err = h.Service.GetPermissionsWithCategory(
			ctx,
			param.SearchName,
			inCategories,
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	case parameter.PermissionWithCaseDefault:
		permissions, err = h.Service.GetPermissions(
			ctx,
			param.SearchName,
			inCategories,
			param.Order,
			param.Pagination,
			param.Limit,
			param.Cursor,
			param.Offset,
			param.WithCount,
		)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, permissions, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

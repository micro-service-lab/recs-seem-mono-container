package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetPoliciesOnRole is a handler for getting policy on role.
type GetPoliciesOnRole struct {
	Service service.ManagerInterface
}

// GetPoliciesOnRoleParam is a parameter for GetPoliciesOnRole.
type GetPoliciesOnRoleParam struct {
	SearchName string                            `queryParam:"search_name"`
	Order      parameter.PolicyOnRoleOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                   `queryParam:"limit"`
	Offset     parameter.Offset                  `queryParam:"offset"`
	Cursor     parameter.Cursor                  `queryParam:"cursor"`
	Pagination parameter.Pagination              `queryParam:"pagination"`
	WithCount  parameter.WithCount               `queryParam:"with_count"`
}

var getPoliciesOnRoleParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PolicyOrderMethodName): parameter.ParsePolicyOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):              parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):             parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):            parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):        parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):      parameter.ParseWithCountParam,
}

func (h *GetPoliciesOnRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "role_id"))
	parse := queryparam.NewParser(r.URL.Query())
	var param GetPoliciesOnRoleParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getPoliciesOnRoleParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	policies, err := h.Service.GetPoliciesOnRole(ctx, id,
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
	err = response.JSONResponseWriter(ctx, w, response.Success, policies, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

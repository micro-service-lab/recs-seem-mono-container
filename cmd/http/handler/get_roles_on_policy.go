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

// GetRoleOnPolicy is a handler for getting role on policy.
type GetRoleOnPolicy struct {
	Service service.ManagerInterface
}

// GetRoleOnPolicyParam is a parameter for GetRoleOnPolicy.
type GetRoleOnPolicyParam struct {
	SearchName string                            `queryParam:"search_name"`
	Order      parameter.RoleOnPolicyOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                   `queryParam:"limit"`
	Offset     parameter.Offset                  `queryParam:"offset"`
	Cursor     parameter.Cursor                  `queryParam:"cursor"`
	Pagination parameter.Pagination              `queryParam:"pagination"`
	WithCount  parameter.WithCount               `queryParam:"with_count"`
}

var getRoleOnPolicyParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PolicyOrderMethodName): parameter.ParsePolicyOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):              parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):             parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):            parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):        parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):      parameter.ParseWithCountParam,
}

func (h *GetRoleOnPolicy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "policy_id"))
	parse := queryparam.NewParser(r.URL.Query())
	var param GetRoleOnPolicyParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getRoleOnPolicyParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	roles, err := h.Service.GetRolesOnPolicy(ctx, id,
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
	err = response.JSONResponseWriter(ctx, w, response.Success, roles, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

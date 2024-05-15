package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetPolicies is a handler for getting policy.
type GetPolicies struct {
	Service service.ManagerInterface
}

// GetPoliciesParam is a parameter for GetPolicies.
type GetPoliciesParam struct {
	SearchName       string                      `queryParam:"search_name"`
	SearchCategories []parameter.EntityID        `queryParam:"search_categories"`
	Order            parameter.PolicyOrderMethod `queryParam:"order"`
	Limit            parameter.Limit             `queryParam:"limit"`
	Offset           parameter.Offset            `queryParam:"offset"`
	Cursor           parameter.Cursor            `queryParam:"cursor"`
	Pagination       parameter.Pagination        `queryParam:"pagination"`
	WithCount        parameter.WithCount         `queryParam:"with_count"`
	With             parameter.PolicyWithParams  `queryParam:"with"`
}

var getPoliciesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PolicyOrderMethodName): parameter.ParsePolicyOrderMethod,
	reflect.TypeOf(parameter.EntityID(uuid.UUID{})): parameter.ParseEntityIDParam,
	reflect.TypeOf(parameter.Limit(0)):              parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):             parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):            parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):        parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):      parameter.ParseWithCountParam,
	reflect.TypeOf(parameter.PolicyWith{}):          parameter.ParsePolicyWithParam,
}

func (h *GetPolicies) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetPoliciesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getPoliciesParseFuncMap,
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
	var policies any
	switch param.With.Case() {
	case parameter.PolicyWithCaseCategory:
		policies, err = h.Service.GetPoliciesWithCategory(
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
	case parameter.PolicyWithCaseDefault:
		policies, err = h.Service.GetPolicies(
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
	err = response.JSONResponseWriter(ctx, w, response.Success, policies, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

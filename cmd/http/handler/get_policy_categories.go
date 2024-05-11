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

// GetPolicyCategories is a handler for getting policy categories.
type GetPolicyCategories struct {
	Service service.ManagerInterface
}

// GetPolicyCategoriesParam is a parameter for GetPolicyCategories.
type GetPolicyCategoriesParam struct {
	SearchName string                              `queryParam:"search_name"`
	Order      parameter.PolicyCategoryOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                     `queryParam:"limit"`
	Offset     parameter.Offset                    `queryParam:"offset"`
	Cursor     parameter.Cursor                    `queryParam:"cursor"`
	Pagination parameter.Pagination                `queryParam:"pagination"`
	WithCount  parameter.WithCount                 `queryParam:"with_count"`
}

var getPolicyCategoriesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PolicyCategoryOrderMethodName): parameter.ParsePolicyCategoryOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                      parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                     parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                    parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):                parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):              parameter.ParseWithCountParam,
}

func (h *GetPolicyCategories) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetPolicyCategoriesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getPolicyCategoriesParseFuncMap,
	})
	if err != nil {
		log.Printf("failed to parse query: %v", err)
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		if !handled {
			if err := response.JSONResponseWriter(ctx, w, response.System, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
		}
		return
	}
	policyCategories, err := h.Service.GetPolicyCategories(
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
		log.Printf("failed to get policy categories: %v", err)
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		if !handled {
			if err := response.JSONResponseWriter(ctx, w, response.System, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, policyCategories, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

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

// GetRecordTypes is a handler for getting record types.
type GetRecordTypes struct {
	Service service.ManagerInterface
}

// GetRecordTypesParam is a parameter for GetRecordTypes.
type GetRecordTypesParam struct {
	SearchName string                          `queryParam:"search_name"`
	Order      parameter.RecordTypeOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                 `queryParam:"limit"`
	Offset     parameter.Offset                `queryParam:"offset"`
	Cursor     parameter.Cursor                `queryParam:"cursor"`
	Pagination parameter.Pagination            `queryParam:"pagination"`
	WithCount  parameter.WithCount             `queryParam:"with_count"`
}

var getRecordTypesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.RecordTypeOrderMethodName): parameter.ParseRecordTypeOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                  parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                 parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):            parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):          parameter.ParseWithCountParam,
}

func (h *GetRecordTypes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetRecordTypesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getRecordTypesParseFuncMap,
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
	attendStatuses, err := h.Service.GetRecordTypes(
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
		log.Printf("failed to get record types: %v", err)
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
	err = response.JSONResponseWriter(ctx, w, response.Success, attendStatuses, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

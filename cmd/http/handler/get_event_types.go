package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// GetEventTypes is a handler for getting event types.
type GetEventTypes struct {
	Service service.ManagerInterface
}

// GetEventTypesParam is a parameter for GetEventTypes.
type GetEventTypesParam struct {
	SearchName string                         `queryParam:"search_name"`
	Order      parameter.EventTypeOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                `queryParam:"limit"`
	Offset     parameter.Offset               `queryParam:"offset"`
	Cursor     parameter.Cursor               `queryParam:"cursor"`
	Pagination parameter.Pagination           `queryParam:"pagination"`
	WithCount  parameter.WithCount            `queryParam:"with_count"`
}

var getEventTypesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.EventTypeOrderMethodName): parameter.ParseEventTypeOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                 parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):               parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):           parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):         parameter.ParseWithCountParam,
}

func (h *GetEventTypes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetEventTypesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getEventTypesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	eventTypes, err := h.Service.GetEventTypes(
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

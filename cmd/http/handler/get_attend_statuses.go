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

// GetAttendStatues is a handler for getting attend statuses.
type GetAttendStatues struct {
	Service service.ManagerInterface
}

// GetAttendStatusesParam is a parameter for GetAttendStatues.
type GetAttendStatusesParam struct {
	SearchName string                            `queryParam:"search_name"`
	Order      parameter.AttendStatusOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                   `queryParam:"limit"`
	Offset     parameter.Offset                  `queryParam:"offset"`
	Cursor     parameter.Cursor                  `queryParam:"cursor"`
	Pagination parameter.Pagination              `queryParam:"pagination"`
	WithCount  parameter.WithCount               `queryParam:"with_count"`
}

var parseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.AttendStatusOrderMethodName): parameter.ParseAttendStatusOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                    parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                   parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                  parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):              parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):            parameter.ParseWithCountParam,
}

func (h *GetAttendStatues) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetAttendStatusesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: parseFuncMap,
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
	attendStatuses, err := h.Service.GetAttendStatuses(
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
		log.Printf("failed to get attend statuses: %v", err)
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

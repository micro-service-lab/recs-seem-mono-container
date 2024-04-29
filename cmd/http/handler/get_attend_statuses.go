package handler

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

type GetAttendStatues struct {
	Service service.ManagerInterface
}

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
	parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: parseFuncMap,
	})
	fmt.Printf("param: %+v\n", param)
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
		errhandle.ErrorHandle(ctx, w, err)
		return
	}
	err = response.JsonResponseWriter(ctx, w, response.Success, attendStatuses, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
	return
}

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

// GetAttendanceTypes is a handler for getting attendance types.
type GetAttendanceTypes struct {
	Service service.ManagerInterface
}

// GetAttendanceTypesParam is a parameter for GetAttendanceTypes.
type GetAttendanceTypesParam struct {
	SearchName string                              `queryParam:"search_name"`
	Order      parameter.AttendanceTypeOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                     `queryParam:"limit"`
	Offset     parameter.Offset                    `queryParam:"offset"`
	Cursor     parameter.Cursor                    `queryParam:"cursor"`
	Pagination parameter.Pagination                `queryParam:"pagination"`
	WithCount  parameter.WithCount                 `queryParam:"with_count"`
}

var getAttendanceTypesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.AttendanceTypeOrderMethodName): parameter.ParseAttendanceTypeOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                      parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                     parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                    parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):                parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):              parameter.ParseWithCountParam,
}

func (h *GetAttendanceTypes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetAttendanceTypesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getAttendanceTypesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	attendanceTypes, err := h.Service.GetAttendanceTypes(
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
	err = response.JSONResponseWriter(ctx, w, response.Success, attendanceTypes, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

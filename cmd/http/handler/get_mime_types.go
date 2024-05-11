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

// GetMimeTypes is a handler for getting mime types.
type GetMimeTypes struct {
	Service service.ManagerInterface
}

// GetMimeTypesParam is a parameter for GetMimeTypes.
type GetMimeTypesParam struct {
	SearchName string                        `queryParam:"search_name"`
	Order      parameter.MimeTypeOrderMethod `queryParam:"order"`
	Limit      parameter.Limit               `queryParam:"limit"`
	Offset     parameter.Offset              `queryParam:"offset"`
	Cursor     parameter.Cursor              `queryParam:"cursor"`
	Pagination parameter.Pagination          `queryParam:"pagination"`
	WithCount  parameter.WithCount           `queryParam:"with_count"`
}

var getMimeTypesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.MimeTypeOrderMethodName): parameter.ParseMimeTypeOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):               parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):              parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):          parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):        parameter.ParseWithCountParam,
}

func (h *GetMimeTypes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetMimeTypesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getMimeTypesParseFuncMap,
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
	mimeTypes, err := h.Service.GetMimeTypes(
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
		log.Printf("failed to get mime types: %v", err)
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
	err = response.JSONResponseWriter(ctx, w, response.Success, mimeTypes, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

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

// GetChatRoomActionTypes is a handler for getting chat room action types.
type GetChatRoomActionTypes struct {
	Service service.ManagerInterface
}

// GetChatRoomActionTypesParam is a parameter for GetChatRoomActionTypes.
type GetChatRoomActionTypesParam struct {
	SearchName string                                  `queryParam:"search_name"`
	Order      parameter.ChatRoomActionTypeOrderMethod `queryParam:"order"`
	Limit      parameter.Limit                         `queryParam:"limit"`
	Offset     parameter.Offset                        `queryParam:"offset"`
	Cursor     parameter.Cursor                        `queryParam:"cursor"`
	Pagination parameter.Pagination                    `queryParam:"pagination"`
	WithCount  parameter.WithCount                     `queryParam:"with_count"`
}

var getChatRoomActionTypesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.ChatRoomActionTypeOrderMethodName): parameter.ParseChatRoomActionTypeOrderMethod,
	reflect.TypeOf(parameter.Limit(0)):                          parameter.ParseLimitParam,
	reflect.TypeOf(parameter.Offset(0)):                         parameter.ParseOffsetParam,
	reflect.TypeOf(parameter.Cursor("")):                        parameter.ParseCursorParam,
	reflect.TypeOf(parameter.NonePagination):                    parameter.ParsePaginationParam,
	reflect.TypeOf(parameter.WithCount(false)):                  parameter.ParseWithCountParam,
}

func (h *GetChatRoomActionTypes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parse := queryparam.NewParser(r.URL.Query())
	var param GetChatRoomActionTypesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getChatRoomActionTypesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	chatRoomActionTypes, err := h.Service.GetChatRoomActionTypes(
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
	err = response.JSONResponseWriter(ctx, w, response.Success, chatRoomActionTypes, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

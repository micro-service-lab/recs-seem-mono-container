package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// GetChatRoomsOnAuth is a handler for getting chat room.
type GetChatRoomsOnAuth struct {
	Service service.ManagerInterface
}

func (h *GetChatRoomsOnAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	parse := queryparam.NewParser(r.URL.Query())
	var param GetChatRoomsOnMemberParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: getChatRoomsOnMemberParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	fmt.Printf("cursor: %v\n", param.Cursor)
	chatRooms, err := h.Service.GetChatRoomsOnMember(
		ctx,
		authUser.MemberID,
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
	err = response.JSONResponseWriter(ctx, w, response.Success, chatRooms, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

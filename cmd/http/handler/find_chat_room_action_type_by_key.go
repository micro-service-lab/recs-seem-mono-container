package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindChatRoomActionTypeByKey is a handler for finding chat room action type.
type FindChatRoomActionTypeByKey struct {
	Service service.ManagerInterface
}

func (h *FindChatRoomActionTypeByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "chat_room_action_type_key")
	chatRoomActionType, err := h.Service.FindChatRoomActionTypeByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, chatRoomActionType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

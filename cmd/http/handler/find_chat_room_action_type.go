package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindChatRoomActionType is a handler for finding chat room action type.
type FindChatRoomActionType struct {
	Service service.ManagerInterface
}

func (h *FindChatRoomActionType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "chat_room_action_type_id"))
	chatRoomActionType, err := h.Service.FindChatRoomActionTypeByID(ctx, id)
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

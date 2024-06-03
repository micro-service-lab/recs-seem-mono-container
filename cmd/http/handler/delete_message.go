package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// DeleteMessage is a handler for deleting message.
type DeleteMessage struct {
	Service service.ManagerInterface
}

func (h *DeleteMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	chatRoomID := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	messageID := uuid.MustParse(chi.URLParam(r, "message_id"))
	_, err := h.Service.DeleteMessage(ctx, chatRoomID, authUser.MemberID, messageID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MessageTargetMessages:
				e.SetTarget("message")
				err = e
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

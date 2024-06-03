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

// ReadMessage is a handler for reading message.
type ReadMessage struct {
	Service service.ManagerInterface
}

// ServeHTTP reads message.
func (h *ReadMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	chatRoomID := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	messageID := uuid.MustParse(chi.URLParam(r, "message_id"))
	read, err := h.Service.ReadMessage(ctx, chatRoomID, authUser.MemberID, messageID)
	if err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MessageTargetMessages:
				e.SetTarget("message")
				err = e
			case service.ReadReceiptTargetReadReceipts:
				e.SetTarget("read receipt")
				err = e
			}
		}
	} else {
		res := struct {
			Read bool `json:"read"`
		}{
			Read: read,
		}
		err = response.JSONResponseWriter(ctx, w, response.Success, res, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

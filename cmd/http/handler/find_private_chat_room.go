package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// FindPrivateChatRoom is a handler for finding event type.
type FindPrivateChatRoom struct {
	Service service.ManagerInterface
}

func (h *FindPrivateChatRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	receiverID := uuid.MustParse(chi.URLParam(r, "member_id"))
	chatRoom, err := h.Service.FindPrivateChatRoom(ctx, authUser.MemberID, receiverID)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, chatRoom, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

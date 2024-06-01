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

// WithdrawOnChatRoom is a handler for remove members on chat room.
type WithdrawOnChatRoom struct {
	Service service.ManagerInterface
}

func (h *WithdrawOnChatRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	id := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	var err error

	if _, err = h.Service.WithdrawMemberFromChatRoom(ctx, id, authUser.MemberID); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomBelongingTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.ChatRoomBelongingTargetChatRoom:
				e.SetTarget("chat_room")
				err = e
			case service.ChatRoomBelongingTargetChatRoomBelongings:
				e.SetTarget("chat_room_belongings")
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

package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// CreatePrivateChatRoom is a handler for creating private chat room.
type CreatePrivateChatRoom struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

func (h *CreatePrivateChatRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	receiverID := uuid.MustParse(chi.URLParam(r, "member_id"))
	var err error

	var message entity.ChatRoom
	if message, err = h.Service.CreatePrivateChatRoom(
		ctx,
		authUser.MemberID,
		receiverID,
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.ChatRoomTargetMembers:
				e.SetTarget("member")
				err = e
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, message, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

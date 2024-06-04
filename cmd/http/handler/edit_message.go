package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// EditMessage is a handler for editing message.
type EditMessage struct {
	Service   service.ManagerInterface
	Validator validation.Validator
}

// EditMessageRequest is a request for EditMessage.
type EditMessageRequest struct {
	Content string `json:"content" validate:"required,max=255" ja:"コンテンツ" en:"Content"`
}

func (h *EditMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	chatRoomID := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	messageID := uuid.MustParse(chi.URLParam(r, "message_id"))
	var err error
	var messageReq EditMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&messageReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			messageReq = EditMessageRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &messageReq, lang.GetLocale(r.Context()))
	} else {
		err = errhandle.NewJSONFormatError()
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var message entity.Message
	if message, err = h.Service.EditMessage(
		ctx,
		chatRoomID,
		authUser.MemberID,
		messageID,
		messageReq.Content,
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MessageTargetMessages:
				e.SetTarget("message")
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

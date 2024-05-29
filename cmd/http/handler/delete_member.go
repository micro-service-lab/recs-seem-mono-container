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

// DeleteMember is a handler for deleting member.
type DeleteMember struct {
	Service service.ManagerInterface
}

func (h *DeleteMember) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "member_id"))
	c, err := h.Service.DeleteMember(ctx, id)
	if err != nil || c == 0 {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

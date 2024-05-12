package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindRole is a handler for finding role.
type FindRole struct {
	Service service.ManagerInterface
}

func (h *FindRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "role_id"))
	eventType, err := h.Service.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrDataNoRecord) {
			if err := response.JSONResponseWriter(ctx, w, response.NotFoundModel, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
			return
		}
		log.Printf("failed to find role: %v", err)
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		if !handled {
			if err := response.JSONResponseWriter(ctx, w, response.System, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, eventType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

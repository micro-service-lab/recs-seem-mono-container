package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindEventTypeByKey is a handler for finding event type.
type FindEventTypeByKey struct {
	Service service.ManagerInterface
}

func (h *FindEventTypeByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "event_type_key")
	eventType, err := h.Service.FindEventTypeByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, eventType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

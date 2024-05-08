package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
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
	attendStatus, err := h.Service.FindEventTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, store.ErrDataNoRecord) {
			if err := response.JSONResponseWriter(ctx, w, response.NotFound, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
			return
		}
		log.Printf("failed to find event type: %v", err)
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
	err = response.JSONResponseWriter(ctx, w, response.Success, attendStatus, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindAttendStatusByKey is a handler for finding attend status.
type FindAttendStatusByKey struct {
	Service service.ManagerInterface
}

func (h *FindAttendStatusByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "attend_status_key")
	attendStatus, err := h.Service.FindAttendStatusByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, attendStatus, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

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

// FindRecordTypeByKey is a handler for finding record type.
type FindRecordTypeByKey struct {
	Service service.ManagerInterface
}

func (h *FindRecordTypeByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "record_type_key")
	recordType, err := h.Service.FindRecordTypeByKey(ctx, key)
	if err != nil {
		if errors.Is(err, store.ErrDataNoRecord) {
			if err := response.JSONResponseWriter(ctx, w, response.NotFoundModel, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
			return
		}
		log.Printf("failed to find record type: %v", err)
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
	err = response.JSONResponseWriter(ctx, w, response.Success, recordType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

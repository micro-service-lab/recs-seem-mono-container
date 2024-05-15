package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindMimeTypeByKey is a handler for finding mime type.
type FindMimeTypeByKey struct {
	Service service.ManagerInterface
}

func (h *FindMimeTypeByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "mime_type_key")
	mimeType, err := h.Service.FindMimeTypeByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, mimeType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

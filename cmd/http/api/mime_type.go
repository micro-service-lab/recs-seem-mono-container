package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// MimeTypeHandler is a handler for mime types.
func MimeTypeHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetMimeTypes{
		Service: svc,
	}
	findHandler := handler.FindMimeType{
		Service: svc,
	}
	findByKeyHandler := handler.FindMimeTypeByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{mime_type_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{mime_type_key}", findByKeyHandler.ServeHTTP)

	return r
}

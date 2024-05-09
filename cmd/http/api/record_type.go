package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// RecordTypeHandler is a handler for record types.
func RecordTypeHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetRecordTypes{
		Service: svc,
	}
	findHandler := handler.FindRecordType{
		Service: svc,
	}
	findByKeyHandler := handler.FindRecordTypeByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{record_type_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{record_type_key}", findByKeyHandler.ServeHTTP)

	return r
}

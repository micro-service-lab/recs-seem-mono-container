package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// EventTypeHandler is a handler for event types.
func EventTypeHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetEventTypes{
		Service: svc,
	}
	findHandler := handler.FindEventType{
		Service: svc,
	}
	findByKeyHandler := handler.FindEventTypeByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{event_type_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{event_type_key}", findByKeyHandler.ServeHTTP)

	return r
}

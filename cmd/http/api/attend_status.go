package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// AttendStatusHandler is a handler for attend statuses.
func AttendStatusHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetAttendStatuses{
		Service: svc,
	}
	findHandler := handler.FindAttendStatus{
		Service: svc,
	}
	findByKeyHandler := handler.FindAttendStatusByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{attend_status_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{attend_status_key}", findByKeyHandler.ServeHTTP)

	return r
}

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// AttendStatusHandler is a handler for attend statuses.
func AttendStatusHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetAttendStatues{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)

	return r
}

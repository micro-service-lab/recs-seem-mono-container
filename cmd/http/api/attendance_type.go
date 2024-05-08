package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// AttendanceTypeHandler is a handler for attend statuses.
func AttendanceTypeHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetAttendanceTypes{
		Service: svc,
	}
	findHandler := handler.FindAttendanceType{
		Service: svc,
	}
	findByKeyHandler := handler.FindAttendanceTypeByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{attendance_type_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{attendance_type_key}", findByKeyHandler.ServeHTTP)

	return r
}

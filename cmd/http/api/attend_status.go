package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

func AttendStatusHandler(svc service.ManagerInterface) func(r chi.Router) {
	getHandler := handler.GetAttendStatues{
		Service: svc,
	}
	return func(r chi.Router) {
		r.Get("/attend_statuses", getHandler.ServeHTTP)
	}
}

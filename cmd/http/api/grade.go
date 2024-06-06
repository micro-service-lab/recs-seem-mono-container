package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// GradeHandler is a handler for grades.
func GradeHandler(
	svc service.ManagerInterface,
) http.Handler {
	getHandler := handler.GetGrades{
		Service: svc,
	}

	r := chi.NewRouter()

	r.Get("/", getHandler.ServeHTTP)

	return r
}

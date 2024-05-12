package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// RoleHandler is a handler for roles.
func RoleHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetRoles{
		Service: svc,
	}
	findHandler := handler.FindRole{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{role_id:uuid}"), findHandler.ServeHTTP)

	return r
}

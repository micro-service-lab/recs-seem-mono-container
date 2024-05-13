package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// RoleHandler is a handler for roles.
func RoleHandler(svc service.ManagerInterface, vd validation.Validator) http.Handler {
	getHandler := handler.GetRoles{
		Service: svc,
	}
	findHandler := handler.FindRole{
		Service: svc,
	}
	createHandler := handler.CreateRole{
		Service:   svc,
		Validator: vd,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Post("/", createHandler.ServeHTTP)
	r.Get(uuidPath("/{role_id:uuid}"), findHandler.ServeHTTP)

	return r
}

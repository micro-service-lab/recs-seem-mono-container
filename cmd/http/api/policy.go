package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// PolicyHandler is a handler for policies.
func PolicyHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetPolicies{
		Service: svc,
	}
	findHandler := handler.FindPolicy{
		Service: svc,
	}
	findByKeyHandler := handler.FindPolicyByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{policy_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{policy_key}", findByKeyHandler.ServeHTTP)

	return r
}

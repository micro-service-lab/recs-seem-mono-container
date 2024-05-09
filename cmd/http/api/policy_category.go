package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// PolicyCategoryHandler is a handler for policy categories.
func PolicyCategoryHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetPolicyCategories{
		Service: svc,
	}
	findHandler := handler.FindPolicyCategory{
		Service: svc,
	}
	findByKeyHandler := handler.FindPolicyCategoryByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{policy_category_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{policy_category_key}", findByKeyHandler.ServeHTTP)

	return r
}

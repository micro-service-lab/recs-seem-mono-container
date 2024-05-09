package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// PermissionCategoryHandler is a handler for permission categories.
func PermissionCategoryHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetPermissionCategories{
		Service: svc,
	}
	findHandler := handler.FindPermissionCategory{
		Service: svc,
	}
	findByKeyHandler := handler.FindPermissionCategoryByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{permission_category_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{permission_category_key}", findByKeyHandler.ServeHTTP)

	return r
}

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// PermissionHandler is a handler for permissions.
func PermissionHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetPermissions{
		Service: svc,
	}
	findHandler := handler.FindPermission{
		Service: svc,
	}
	findByKeyHandler := handler.FindPermissionByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{permission_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{permission_key}", findByKeyHandler.ServeHTTP)

	return r
}

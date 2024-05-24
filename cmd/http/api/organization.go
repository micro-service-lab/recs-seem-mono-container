package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// OrganizationHandler is a handler for roles.
func OrganizationHandler(svc service.ManagerInterface, vd validation.Validator, _ i18n.Translation) http.Handler {
	getHandler := handler.GetOrganizations{
		Service: svc,
	}
	findHandler := handler.FindOrganization{
		Service: svc,
	}
	updateHandler := handler.UpdateOrganization{
		Service:   svc,
		Validator: vd,
	}
	createHandler := handler.CreateOrganization{
		Service:   svc,
		Validator: vd,
	}
	deleteHandler := handler.DeleteOrganization{
		Service: svc,
	}

	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{organization_id:uuid}"), findHandler.ServeHTTP)
	r.Post("/", createHandler.ServeHTTP)
	r.Put(uuidPath("/{organization_id:uuid}"), updateHandler.ServeHTTP)
	r.Delete(uuidPath("/{organization_id:uuid}"), deleteHandler.ServeHTTP)

	return r
}

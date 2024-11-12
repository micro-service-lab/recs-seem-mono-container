package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// OrganizationHandler is a handler for roles.
func OrganizationHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	createHandler := handler.CreateOrganization{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteHandler := handler.DeleteOrganization{
		Service: svc,
	}
	getHandler := handler.GetOrganizations{
		Service: svc,
	}
	findHandler := handler.FindOrganization{
		Service: svc,
	}
	updateHandler := handler.UpdateOrganization{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	addMembers := handler.AddMembersOnOrganization{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	removeMembers := handler.RemoveMembersOnOrganization{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	withdraw := handler.WithdrawOnOrganization{
		Service: svc,
	}
	getMembers := handler.GetMembersOnOrganization{
		Service: svc,
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Post("/", createHandler.ServeHTTP)
		r.Delete(uuidPath("/{organization_id:uuid}"), deleteHandler.ServeHTTP)
		r.Put(uuidPath("/{organization_id:uuid}"), updateHandler.ServeHTTP)
		r.Get("/", getHandler.ServeHTTP)
		r.Get(uuidPath("/{organization_id:uuid}"), findHandler.ServeHTTP)

		r.Route(uuidPath("/{organization_id:uuid}/members"), func(r chi.Router) {
			r.Post("/", addMembers.ServeHTTP)
			r.Post("/withdraw", withdraw.ServeHTTP)
			r.Delete("/", removeMembers.ServeHTTP)
			r.Get("/", getMembers.ServeHTTP)
		})
	})

	return r
}

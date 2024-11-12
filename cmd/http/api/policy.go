package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// PolicyHandler is a handler for policies.
func PolicyHandler(svc service.ManagerInterface, vd validation.Validator, t i18n.Translation) http.Handler {
	getHandler := handler.GetPolicies{
		Service: svc,
	}
	findHandler := handler.FindPolicy{
		Service: svc,
	}
	findByKeyHandler := handler.FindPolicyByKey{
		Service: svc,
	}

	associateRole := handler.AssociateRolesOnPolicy{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	disassociateRole := handler.DisassociateRolesOnPolicy{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	disassociateAllRole := handler.DisassociateAllRolesOnPolicy{
		Service: svc,
	}
	getAssociateRoles := handler.GetRoleOnPolicy{
		Service: svc,
	}

	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{policy_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{policy_key}", findByKeyHandler.ServeHTTP)
	r.Route(uuidPath("/{policy_id:uuid}/roles"), func(r chi.Router) {
		r.Post("/", associateRole.ServeHTTP)
		r.Delete("/", disassociateRole.ServeHTTP)
		r.Delete("/all", disassociateAllRole.ServeHTTP)
		r.Get("/", getAssociateRoles.ServeHTTP)
	})

	return r
}

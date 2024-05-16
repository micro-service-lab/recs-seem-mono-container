package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// RoleHandler is a handler for roles.
func RoleHandler(svc service.ManagerInterface, vd validation.Validator, t i18n.Translation) http.Handler {
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
	updateHandler := handler.UpdateRole{
		Service:   svc,
		Validator: vd,
	}
	deleteHandler := handler.DeleteRole{
		Service: svc,
	}

	associatePolicies := handler.AssociatePoliciesOnRole{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	disassociatePolicies := handler.DisassociatePoliciesOnRole{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	disassociateAllPolicies := handler.DisassociateAllPoliciesOnRole{
		Service: svc,
	}
	getAssociatePolicies := handler.GetPoliciesOnRole{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Put(uuidPath("/{role_id:uuid}"), updateHandler.ServeHTTP)
	r.Delete(uuidPath("/{role_id:uuid}"), deleteHandler.ServeHTTP)
	r.Post("/", createHandler.ServeHTTP)
	r.Get(uuidPath("/{role_id:uuid}"), findHandler.ServeHTTP)
	r.Route(uuidPath("/{role_id:uuid}/policies"), func(r chi.Router) {
		r.Post("/", associatePolicies.ServeHTTP)
		r.Delete("/", disassociatePolicies.ServeHTTP)
		r.Delete("/all", disassociateAllPolicies.ServeHTTP)
		r.Get("/", getAssociatePolicies.ServeHTTP)
	})

	return r
}

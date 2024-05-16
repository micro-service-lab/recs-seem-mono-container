package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// DisassociateAllPoliciesOnRole is a handler for disassociating all policies on role.
type DisassociateAllPoliciesOnRole struct {
	Service service.ManagerInterface
}

func (h *DisassociateAllPoliciesOnRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "role_id"))
	var err error

	if _, err = h.Service.DisassociatePolicyOnRole(ctx, id); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			if e.Target() == service.AssociateRoleTargetRoles {
				e.SetTarget("role")
				err = e
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	}

	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

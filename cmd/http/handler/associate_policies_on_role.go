package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// AssociatePoliciesOnRole is a handler for creating role.
type AssociatePoliciesOnRole struct {
	Service   service.ManagerInterface
	Validator validation.Validator
}

// AssociatePoliciesOnRoleRequest is a request for AssociatePoliciesOnRole.
type AssociatePoliciesOnRoleRequest struct {
	PolicyIDS []uuid.UUID `json:"policy_ids" validate:"required,dive,uuid" ja:"ポリシーID" en:"PolicyIDs"`
}

func (h *AssociatePoliciesOnRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "role_id"))
	var err error
	var req AssociatePoliciesOnRoleRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = AssociatePoliciesOnRoleRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &req, lang.GetLocale(r.Context()))
	} else {
		err = errhandle.NewJSONFormatError()
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	policies := make([]parameter.AssociationRoleParam, len(req.PolicyIDS))
	for i, pID := range req.PolicyIDS {
		policies[i] = parameter.AssociationRoleParam{
			RoleID:   id,
			PolicyID: pID,
		}
	}
	if _, err = h.Service.AssociateRoles(ctx, policies); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			if e.Target() == service.AssociateRoleTargetRoles {
				e.SetTarget("role")
			} else if e.Target() == service.AssociateRoleTargetPolicies {
				e.SetTarget("policy")
			}
			err = e
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

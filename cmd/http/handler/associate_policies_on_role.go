package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// AssociatePoliciesOnRole is a handler for associating policies on role.
type AssociatePoliciesOnRole struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// AssociatePoliciesOnRoleRequest is a request for AssociatePoliciesOnRole.
type AssociatePoliciesOnRoleRequest struct {
	PolicyIDS []uuid.UUID `json:"policy_ids" validate:"required,unique" ja:"ポリシーID" en:"PolicyIDs"`
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
				err = e
			} else if e.Target() == service.AssociateRoleTargetPolicies {
				policyStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "PolicyIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "PolicyIDs",
						Other: "Policy not found",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": policyStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("policy_ids", msgStr)
				err = ve
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

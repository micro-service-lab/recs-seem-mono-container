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

// AssociateRolesOnPolicy is a handler for associating policies on policy.
type AssociateRolesOnPolicy struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// AssociateRolesOnPolicyRequest is a request for AssociateRolesOnPolicy.
type AssociateRolesOnPolicyRequest struct {
	RoleIDS []uuid.UUID `json:"role_ids" validate:"required,unique" ja:"ロールID" en:"RoleIDs"`
}

func (h *AssociateRolesOnPolicy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "policy_id"))
	var err error
	var req AssociateRolesOnPolicyRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = AssociateRolesOnPolicyRequest{}
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
	policies := make([]parameter.AssociationRoleParam, len(req.RoleIDS))
	for i, pID := range req.RoleIDS {
		policies[i] = parameter.AssociationRoleParam{
			PolicyID: id,
			RoleID:   pID,
		}
	}
	if _, err = h.Service.AssociateRoles(ctx, policies); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			if e.Target() == service.AssociateRoleTargetPolicies {
				e.SetTarget("policy")
				err = e
			} else if e.Target() == service.AssociateRoleTargetRoles {
				roleStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "RoleIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "RoleIDs",
						Other: "Role not found",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": roleStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("role_ids", msgStr)
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

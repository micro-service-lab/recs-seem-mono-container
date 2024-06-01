package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// UpdateOrganization is a handler for creating organization.
type UpdateOrganization struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// UpdateOrganizationRequest is a request for UpdateOrganization.
type UpdateOrganizationRequest struct {
	Name        string `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	Description string `json:"description,omitempty" validate:"omitempty" ja:"説明" en:"Description"`
	Color       string `json:"color,omitempty" validate:"omitempty,hexcolor" ja:"色" en:"Color"`
}

func (h *UpdateOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "organization_id"))
	authUser := auth.FromContext(ctx)
	var err error
	var organizationReq UpdateOrganizationRequest
	if err = json.NewDecoder(r.Body).Decode(&organizationReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			organizationReq = UpdateOrganizationRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &organizationReq, lang.GetLocale(r.Context()))
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
	var organization entity.Organization
	if organization, err = h.Service.UpdateOrganization(
		ctx,
		id,
		organizationReq.Name,
		entity.String{
			Valid:  organizationReq.Description != "",
			String: organizationReq.Description,
		},
		entity.String{
			Valid:  organizationReq.Color != "",
			String: organizationReq.Color,
		},
		authUser.MemberID,
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomTargetChatRoom:
				e.SetTarget("chat room")
				err = e
			case service.OrganizationTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.OrganizationTargetOrganizations:
				e.SetTarget("organization")
				err = e
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, organization, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

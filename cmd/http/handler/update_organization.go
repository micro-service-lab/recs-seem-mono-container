package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// UpdateOrganization is a handler for creating organization.
type UpdateOrganization struct {
	Service   service.ManagerInterface
	Validator validation.Validator
}

// UpdateOrganizationRequest is a request for UpdateOrganization.
type UpdateOrganizationRequest struct {
	Name        string `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	Description string `json:"description,omitempty" validate:"omitempty" ja:"説明" en:"Description"`
	Color       string `json:"color,omitempty" validate:"omitempty,hexcolor" ja:"色" en:"Color"`
}

func (h *UpdateOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// id := uuid.MustParse(chi.URLParam(r, "organization_id"))
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
	// dsc := entity.String{
	// 	String: organizationReq.Description,
	// 	Valid:  organizationReq.Description != "",
	// }
	// col := entity.String{
	// 	String: organizationReq.Color,
	// 	Valid:  organizationReq.Color != "",
	// }
	// var organization entity.Organization
	// if organization, err = h.Service.UpdateOrganization(ctx, id, organizationReq.Name, dsc, col); err == nil {
	// 	err = response.JSONResponseWriter(ctx, w, response.Success, organization, nil)
	// }
	// if err != nil {
	// 	handled, err := errhandle.ErrorHandle(ctx, w, err)
	// 	if !handled || err != nil {
	// 		log.Printf("failed to handle error: %v", err)
	// 	}
	// 	return
	// }
}

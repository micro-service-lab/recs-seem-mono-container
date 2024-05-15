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
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// UpdateRole is a handler for creating role.
type UpdateRole struct {
	Service   service.ManagerInterface
	Validator validation.Validator
}

// UpdateRoleRequest is a request for UpdateRole.
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	Description string `json:"description" validate:"required" ja:"説明" en:"Description"`
}

func (h *UpdateRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "role_id"))
	var err error
	var roleReq UpdateRoleRequest
	if err = json.NewDecoder(r.Body).Decode(&roleReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			roleReq = UpdateRoleRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &roleReq, lang.GetLocale(r.Context()))
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
	var role entity.Role
	if role, err = h.Service.UpdateRole(ctx, id, roleReq.Name, roleReq.Description); err == nil {
		err = response.JSONResponseWriter(ctx, w, response.Success, role, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

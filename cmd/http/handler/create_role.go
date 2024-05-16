package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// CreateRole is a handler for creating role.
type CreateRole struct {
	Service   service.ManagerInterface
	Validator validation.Validator
}

// CreateRoleRequest is a request for CreateRole.
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	Description string `json:"description" validate:"required" ja:"説明" en:"Description"`
}

func (h *CreateRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var roleReq CreateRoleRequest
	if err = json.NewDecoder(r.Body).Decode(&roleReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			roleReq = CreateRoleRequest{}
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
	if role, err = h.Service.CreateRole(ctx, roleReq.Name, roleReq.Description); err == nil {
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

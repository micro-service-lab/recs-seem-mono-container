package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// CreateRole is a handler for creating role.
type CreateRole struct {
	Service   service.ManagerInterface
	Validator *validator.Validate
}

// CreateRoleRequest is a request for CreateRole.
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *CreateRole) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role, err := h.Service.CreateRole()
	if err != nil {
		if errors.Is(err, store.ErrDataNoRecord) {
			if err := response.JSONResponseWriter(ctx, w, response.NotFound, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
			return
		}
		log.Printf("failed to find role: %v", err)
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		if !handled {
			if err := response.JSONResponseWriter(ctx, w, response.System, nil, nil); err != nil {
				log.Printf("failed to write response: %v", err)
			}
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, role, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

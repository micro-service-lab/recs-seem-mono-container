package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindPolicyCategory is a handler for finding policy category.
type FindPolicyCategory struct {
	Service service.ManagerInterface
}

func (h *FindPolicyCategory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "policy_category_id"))
	policyCategory, err := h.Service.FindPolicyCategoryByID(ctx, id)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, policyCategory, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

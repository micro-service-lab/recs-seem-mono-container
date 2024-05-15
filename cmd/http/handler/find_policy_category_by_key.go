package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindPolicyCategoryByKey is a handler for finding policy category.
type FindPolicyCategoryByKey struct {
	Service service.ManagerInterface
}

func (h *FindPolicyCategoryByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "policy_category_key")
	policyCategory, err := h.Service.FindPolicyCategoryByKey(ctx, key)
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

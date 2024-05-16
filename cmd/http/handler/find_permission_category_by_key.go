package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindPermissionCategoryByKey is a handler for finding permission category.
type FindPermissionCategoryByKey struct {
	Service service.ManagerInterface
}

func (h *FindPermissionCategoryByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "permission_category_key")
	permissionCategory, err := h.Service.FindPermissionCategoryByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, permissionCategory, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

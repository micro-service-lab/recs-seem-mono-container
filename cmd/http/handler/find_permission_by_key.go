package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindPermissionByKey is a handler for finding permission.
type FindPermissionByKey struct {
	Service service.PermissionManager
}

func (h *FindPermissionByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "permission_key")
	parse := queryparam.NewParser(r.URL.Query())
	var param FindPermissionsParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: findPermissionsParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var permission any
	switch param.With.Case() {
	case parameter.PermissionWithCaseCategory:
		permission, err = h.Service.FindPermissionByKeyWithCategory(
			ctx,
			key,
		)
	case parameter.PermissionWithCaseDefault:
		permission, err = h.Service.FindPermissionByKey(ctx, key)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, permission, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

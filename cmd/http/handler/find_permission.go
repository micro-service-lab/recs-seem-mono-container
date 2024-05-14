package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindPermission is a handler for finding permission.
type FindPermission struct {
	Service service.PermissionManager
}

// FindPermissionsParam is a parameter for FindPermissions.
type FindPermissionsParam struct {
	With parameter.PermissionWithParams `queryParam:"with"`
}

var findPermissionsParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PermissionWith{}): parameter.ParsePermissionWithParam,
}

func (h *FindPermission) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "permission_id"))
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
		permission, err = h.Service.FindPermissionByIDWithCategory(
			ctx,
			id,
		)
	case parameter.PermissionWithCaseDefault:
		permission, err = h.Service.FindPermissionByID(ctx, id)
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

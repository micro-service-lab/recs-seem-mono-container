package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/queryparam"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindOrganization is a handler for finding organization.
type FindOrganization struct {
	Service service.ManagerInterface
}

// FindOrganizationParam is a parameter for FindOrganization.
type FindOrganizationParam struct {
	With parameter.OrganizationWithParams `queryParam:"with"`
}

var findOrganizationParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.OrganizationWith{}): parameter.ParseOrganizationWithParam,
}

func (h *FindOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "organization_id"))
	parse := queryparam.NewParser(r.URL.Query())
	var param FindOrganizationParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: findOrganizationParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var organization any
	switch param.With.Case() {
	case parameter.OrganizationWithCaseDefault:
		organization, err = h.Service.FindOrganizationByID(ctx, id)
	case parameter.OrganizationWithCaseDetail:
		organization, err = h.Service.FindOrganizationWithDetail(
			ctx,
			id,
		)
	case parameter.OrganizationWithCaseChatRoom:
		organization, err = h.Service.FindOrganizationWithChatRoom(
			ctx,
			id,
		)
	case parameter.OrganizationWithCaseChatRoomAndDetail:
		organization, err = h.Service.FindOrganizationWithChatRoomAndDetail(
			ctx,
			id,
		)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, organization, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

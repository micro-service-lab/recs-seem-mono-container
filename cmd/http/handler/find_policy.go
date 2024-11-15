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

// FindPolicy is a handler for finding policy.
type FindPolicy struct {
	Service service.ManagerInterface
}

// FindPoliciesParam is a parameter for FindPolicies.
type FindPoliciesParam struct {
	With parameter.PolicyWithParams `queryParam:"with[]"`
}

var findPoliciesParseFuncMap = map[reflect.Type]queryparam.ParserFunc{
	reflect.TypeOf(parameter.PolicyWith{}): parameter.ParsePolicyWithParam,
}

func (h *FindPolicy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.MustParse(chi.URLParam(r, "policy_id"))
	parse := queryparam.NewParser(r.URL.Query())
	var param FindPoliciesParam
	err := parse.ParseWithOptions(&param, queryparam.Options{
		TagName: "queryParam",
		FuncMap: findPoliciesParseFuncMap,
	})
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var policy any
	switch param.With.Case() {
	case parameter.PolicyWithCaseCategory:
		policy, err = h.Service.FindPolicyByIDWithCategory(
			ctx,
			id,
		)
	case parameter.PolicyWithCaseDefault:
		policy, err = h.Service.FindPolicyByID(ctx, id)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, policy, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

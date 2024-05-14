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

// FindPolicyByKey is a handler for finding policy.
type FindPolicyByKey struct {
	Service service.PolicyManager
}

func (h *FindPolicyByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "policy_key")
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
		policy, err = h.Service.FindPolicyByKeyWithCategory(
			ctx,
			key,
		)
	case parameter.PolicyWithCaseDefault:
		policy, err = h.Service.FindPolicyByKey(ctx, key)
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

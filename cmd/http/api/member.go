package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
)

// MemberHandler is a handler for roles.
func MemberHandler(svc service.ManagerInterface, vd validation.Validator, t i18n.Translation) http.Handler {
	createHandler := handler.CreateMember{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	r := chi.NewRouter()
	r.Post("/", createHandler.ServeHTTP)

	return r
}

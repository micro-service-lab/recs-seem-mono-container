package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// ProfessorHandler is a handler for roles.
func ProfessorHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	createHandler := handler.CreateProfessor{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteHandler := handler.DeleteProfessor{
		Service: svc,
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Post("/", createHandler.ServeHTTP)
		r.Delete("/{professor_id}", deleteHandler.ServeHTTP)
	})

	return r
}

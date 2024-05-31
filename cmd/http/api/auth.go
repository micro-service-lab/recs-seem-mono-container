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

// AuthHandler is a handler for auth.
func AuthHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	loginHandler := handler.Login{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	logoutHandler := handler.Logout{
		Service: svc,
	}
	refreshTokenHandler := handler.RefreshToken{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	retrieveAuth := handler.RetrieveAuth{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Post("/logout", logoutHandler.ServeHTTP)
		r.Get("/me", retrieveAuth.ServeHTTP)
	})
	r.Post("/login", loginHandler.ServeHTTP)
	r.Post("/refresh_token", refreshTokenHandler.ServeHTTP)

	return r
}
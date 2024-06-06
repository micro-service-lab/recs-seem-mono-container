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

// ConfidentialHandler is a handler for confidential actions.
func ConfidentialHandler(
	svc service.ManagerInterface,
	_ validation.Validator,
	_ i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	getChatRooms := handler.GetChatRoomsOnAuth{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Get("/chat-rooms", getChatRooms.ServeHTTP)
	})

	return r
}

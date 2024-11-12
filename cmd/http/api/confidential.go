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
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	deleteAuth := handler.DeleteAuth{
		Service: svc,
	}
	updateAuth := handler.UpdateAuth{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	getChatRooms := handler.GetChatRoomsOnAuth{
		Service: svc,
	}
	getOrganizationHandler := handler.GetOrganizationsOnAuth{
		Service: svc,
	}
	createPrivateChatRoom := handler.CreatePrivateChatRoom{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	findPrivateChatRoom := handler.FindPrivateChatRoom{
		Service: svc,
	}

	createPrivateMessage := handler.CreatePrivateMessage{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))
		r.Delete("/", deleteAuth.ServeHTTP)
		r.Put("/", updateAuth.ServeHTTP)

		r.Route("/chat_rooms", func(r chi.Router) {
			r.Get("/", getChatRooms.ServeHTTP)

			r.Route("/private_members/{member_id}", func(r chi.Router) {
				r.Post("/", createPrivateChatRoom.ServeHTTP)
				r.Get("/", findPrivateChatRoom.ServeHTTP)

				r.Route("/messages", func(r chi.Router) {
					r.Post("/", createPrivateMessage.ServeHTTP)
				})
			})
		})

		r.Route("/organizations", func(r chi.Router) {
			r.Get("/", getOrganizationHandler.ServeHTTP)
		})
	})

	return r
}

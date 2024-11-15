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

// MemberHandler is a handler for roles.
func MemberHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	getHandler := handler.GetMembers{
		Service: svc,
	}
	deleteHandler := handler.DeleteMember{
		Service: svc,
	}
	updateHandler := handler.UpdateMember{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	getChatRooms := handler.GetChatRoomsOnMember{
		Service: svc,
	}

	getOrganizations := handler.GetOrganizationsOnMember{
		Service: svc,
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Get("/", getHandler.ServeHTTP)
		r.Delete(uuidPath("/{member_id}"), deleteHandler.ServeHTTP)
		r.Put(uuidPath("/{member_id}"), updateHandler.ServeHTTP)

		r.Route(uuidPath("/{member_id}/chat_rooms"), func(r chi.Router) {
			r.Get("/", getChatRooms.ServeHTTP)
		})

		r.Route(uuidPath("/{member_id}/organizations"), func(r chi.Router) {
			r.Get("/", getOrganizations.ServeHTTP)
		})
	})

	return r
}

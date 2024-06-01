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

// ChatRoomHandler is a handler for chat rooms.
func ChatRoomHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	createHandler := handler.CreateChatRoom{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteHandler := handler.DeleteChatRoom{
		Service: svc,
	}
	updateHandler := handler.UpdateChatRoom{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	addMembers := handler.AddMembersOnChatRoom{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	removeMembers := handler.RemoveMembersOnChatRoom{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	getMembers := handler.GetMembersOnChatRoom{
		Service: svc,
	}

	getActions := handler.GetChatRoomActionsOnChatRoom{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Post("/", createHandler.ServeHTTP)
		r.Delete(uuidPath("/{chat_room_id:uuid}"), deleteHandler.ServeHTTP)
		r.Put(uuidPath("/{chat_room_id:uuid}"), updateHandler.ServeHTTP)

		r.Route(uuidPath("/{chat_room_id:uuid}/members"), func(r chi.Router) {
			r.Post("/", addMembers.ServeHTTP)
			r.Delete("/", removeMembers.ServeHTTP)
			r.Get("/", getMembers.ServeHTTP)
		})
	})

	r.Route(uuidPath("/{chat_room_id:uuid}/chat_room_actions"), func(r chi.Router) {
		r.Get("/", getActions.ServeHTTP)
	})

	return r
}

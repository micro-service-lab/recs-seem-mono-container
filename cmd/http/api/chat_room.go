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

	readOnChatRoom := handler.ReadMessagesOnChatRoom{
		Service: svc,
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
	withdraw := handler.WithdrawOnChatRoom{
		Service: svc,
	}
	getMembers := handler.GetMembersOnChatRoom{
		Service: svc,
	}

	createMessage := handler.CreateMessage{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteMessage := handler.DeleteMessage{
		Service: svc,
	}
	editMessage := handler.EditMessage{
		Service:   svc,
		Validator: vd,
	}
	readMessage := handler.ReadMessage{
		Service: svc,
	}
	attachItemOnMessage := handler.AttachItemsOnMessage{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	detachItemOnMessage := handler.DetachItemsOnMessage{
		Service:    svc,
		Validator:  vd,
		Translator: t,
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
			r.Post("/withdraw", withdraw.ServeHTTP)
			r.Delete("/", removeMembers.ServeHTTP)
			r.Get("/", getMembers.ServeHTTP)
		})

		r.Route(uuidPath("/{chat_room_id:uuid}/messages"), func(r chi.Router) {
			r.Post("/", createMessage.ServeHTTP)
			r.Delete(uuidPath("/{message_id:uuid}"), deleteMessage.ServeHTTP)
			r.Put(uuidPath("/{message_id:uuid}"), editMessage.ServeHTTP)

			r.Route("/read", func(r chi.Router) {
				r.Post("/", readOnChatRoom.ServeHTTP)
			})

			r.Route(uuidPath("/{message_id:uuid}/read"), func(r chi.Router) {
				r.Post("/", readMessage.ServeHTTP)
			})

			r.Route(uuidPath("/{message_id:uuid}/attachable_items"), func(r chi.Router) {
				r.Post("/", attachItemOnMessage.ServeHTTP)
				r.Delete("/", detachItemOnMessage.ServeHTTP)
			})
		})
	})

	r.Route(uuidPath("/{chat_room_id:uuid}/chat_room_actions"), func(r chi.Router) {
		r.Get("/", getActions.ServeHTTP)
	})

	return r
}

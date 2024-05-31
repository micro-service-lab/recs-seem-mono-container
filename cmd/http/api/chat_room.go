package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// ChatRoomHandler is a handler for chat rooms.
func ChatRoomHandler(svc service.ManagerInterface) http.Handler {
	getActions := handler.GetChatRoomActionsOnChatRoom{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Route(uuidPath("/{chat_room_id:uuid}/chat_room_actions"), func(r chi.Router) {
		r.Get("/", getActions.ServeHTTP)
	})

	return r
}

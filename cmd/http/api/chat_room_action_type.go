package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
)

// ChatRoomActionTypeHandler is a handler for chat room action types.
func ChatRoomActionTypeHandler(svc service.ManagerInterface) http.Handler {
	getHandler := handler.GetChatRoomActionTypes{
		Service: svc,
	}
	findHandler := handler.FindChatRoomActionType{
		Service: svc,
	}
	findByKeyHandler := handler.FindChatRoomActionTypeByKey{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Get("/", getHandler.ServeHTTP)
	r.Get(uuidPath("/{chat_room_action_type_id:uuid}"), findHandler.ServeHTTP)
	r.Get("/{chat_room_action_type_key}", findByKeyHandler.ServeHTTP)

	return r
}

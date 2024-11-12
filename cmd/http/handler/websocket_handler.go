package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// WebsocketHandler Websocket ハンドラを表す構造体。
type WebsocketHandler struct {
	hub ws.HubInterface
}

// NewWebsocketHandler WebsocketHandler を生成して返す。
func NewWebsocketHandler(hub ws.HubInterface) *WebsocketHandler {
	return &WebsocketHandler{
		hub: hub,
	}
}

// Handle WebSocket ハンドラ。
func (h *WebsocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	authUser := auth.FromContext(r.Context())
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
	sock, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := ws.NewClient(sock, *authUser)
	go client.ReadLoop(h.hub.Dispatch, h.hub.UnRegisterClient, h.hub.GetOnlineMembers)
	go client.WriteLoop()
	h.hub.RegisterClient(client)
}

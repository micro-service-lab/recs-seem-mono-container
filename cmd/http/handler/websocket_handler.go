package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
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
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	sock, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := ws.NewClient(sock)
	go client.ReadLoop(h.hub.ReadBroadCast(), h.hub.ReadUnregisterClient())
	go client.WriteLoop()
	h.hub.RegisterClient(client)
}

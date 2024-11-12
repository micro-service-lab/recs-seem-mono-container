package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// Client は WebSocket のクライアントを表す構造体。
type Client struct {
	ws     *websocket.Conn
	sendCh chan []byte
	auth   entity.AuthMember
	uid    uuid.UUID
}

// NewClient クライアントを生成して返す。
func NewClient(ws *websocket.Conn, auth entity.AuthMember) *Client {
	return &Client{
		ws:     ws,
		sendCh: make(chan []byte),
		auth:   auth,
		uid:    uuid.New(),
	}
}

// ReadLoop クライアントの読み込みループ。
func (c *Client) ReadLoop(
	dispatch func(eventType EventType, targets Targets, data any), unregisterClient func(c *Client),
	getOnlineMembers func() []ConnectingMember,
) {
	defer func() {
		c.disconnect(unregisterClient)
	}()

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		var req Request
		if err := json.Unmarshal(msg, &req); err != nil {
			continue
		}
		switch req.RequestType {
		case RequestTypeOnlineMembers:
			dispatch(EventTypeConnectingMembers,
				Targets{Members: []uuid.UUID{c.auth.MemberID}},
				ConnectingMembersEventData{
					ConnectingMembers: getOnlineMembers(),
				})
		}
	}
}

// WriteLoop クライアントの書き込みループ。
func (c *Client) WriteLoop() {
	defer func() {
		if err := c.ws.Close(); err != nil {
			log.Printf("error: %v", err)
		}
	}()

	for {
		message := <-c.sendCh

		w, err := c.ws.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		if _, err := w.Write(message); err != nil {
			return
		}

		for i := 0; i < len(c.sendCh); i++ {
			if _, err := w.Write(<-c.sendCh); err != nil {
				return
			}
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

func (c *Client) disconnect(unregisterClient func(c *Client)) {
	unregisterClient(c)
	if err := c.ws.Close(); err != nil {
		log.Printf("error: %v", err)
	}
}

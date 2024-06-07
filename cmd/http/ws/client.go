package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client は WebSocket のクライアントを表す構造体。
type Client struct {
	ws     *websocket.Conn
	sendCh chan []byte
}

// NewClient クライアントを生成して返す。
func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		ws:     ws,
		sendCh: make(chan []byte),
	}
}

// ReadLoop クライアントの読み込みループ。
func (c *Client) ReadLoop(broadCast chan<- []byte, unregister chan<- *Client) {
	defer func() {
		c.disconnect(unregister)
	}()

	for {
		_, jsonMsg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		broadCast <- jsonMsg
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

func (c *Client) disconnect(unregister chan<- *Client) {
	unregister <- c
	if err := c.ws.Close(); err != nil {
		log.Printf("error: %v", err)
	}
}

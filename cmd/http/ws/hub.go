// Package ws provides a WebSocket application.
package ws

import (
	"context"

	"github.com/micro-service-lab/recs-seem-mono-container/app/pubsub"
)

// HubInterface は WebSocket のハブのインターフェース。
type HubInterface interface {
	RunLoop(ctx context.Context)
	SubscribeMessages(ctx context.Context)
	ReadBroadCast() chan<- []byte
	RegisterClient(c *Client)
	ReadUnregisterClient() chan<- *Client
}

// Hub は WebSocket のハブを表す構造体。
type Hub struct {
	Clients      map[*Client]bool
	RegisterCh   chan *Client
	UnRegisterCh chan *Client
	BroadcastCh  chan []byte
	pubsub       pubsub.Service
}

var _ HubInterface = (*Hub)(nil)

const broadCastChan = "wsBroadcast"

// NewHub Hub を生成して返す。
func NewHub(pubsub pubsub.Service) *Hub {
	return &Hub{
		Clients:      make(map[*Client]bool),
		RegisterCh:   make(chan *Client),
		UnRegisterCh: make(chan *Client),
		BroadcastCh:  make(chan []byte),
		pubsub:       pubsub,
	}
}

// RunLoop ハブのメインループを実行する。
func (h *Hub) RunLoop(ctx context.Context) {
	for {
		select {
		case client := <-h.RegisterCh:
			h.register(client)

		case client := <-h.UnRegisterCh:
			h.unregister(client)

		case msg := <-h.BroadcastCh:
			h.publishMessage(ctx, msg)
		}
	}
}

// SubscribeMessages メッセージを受信する。
func (h *Hub) SubscribeMessages(ctx context.Context) {
	ch := h.pubsub.Subscribe(ctx, broadCastChan)

	for msg := range ch {
		h.broadCastToAllClient([]byte(msg.Payload))
	}
}

// ReadBroadCast ブロードキャスト用のチャネルを返す。
func (h *Hub) ReadBroadCast() chan<- []byte {
	return h.BroadcastCh
}

// RegisterClient クライアントを登録する。
func (h *Hub) RegisterClient(c *Client) {
	h.RegisterCh <- c
}

// ReadUnregisterClient クライアントの登録解除用のチャネルを返す。
func (h *Hub) ReadUnregisterClient() chan<- *Client {
	return h.UnRegisterCh
}

func (h *Hub) publishMessage(ctx context.Context, msg []byte) {
	h.pubsub.Publish(ctx, broadCastChan, msg)
}

func (h *Hub) register(c *Client) {
	h.Clients[c] = true
}

func (h *Hub) unregister(c *Client) {
	delete(h.Clients, c)
}

func (h *Hub) broadCastToAllClient(msg []byte) {
	for c := range h.Clients {
		c.sendCh <- msg
	}
}

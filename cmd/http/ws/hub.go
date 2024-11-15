// Package ws provides a WebSocket application.
package ws

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"sync"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/pubsub"
)

const maxBroadcastSize = 100

// HubInterface は WebSocket のハブのインターフェース。
type HubInterface interface {
	RunLoop(ctx context.Context)
	SubscribeMessages(ctx context.Context)
	GetOnlineMembers() []ConnectingMember
	Dispatch(eventType EventType, targets Targets, data any)
	RegisterClient(c *Client)
	UnRegisterClient(c *Client)
}

// Hub は WebSocket のハブを表す構造体。
type Hub struct {
	Clients       map[uuid.UUID][]*Client
	RegisterCh    chan *Client
	UnRegisterCh  chan *Client
	OnlineMembers map[uuid.UUID][]uuid.UUID
	BroadcastCh   chan *Payload
	pubsub        pubsub.Service
	mu            *sync.RWMutex
	onlineMu      *sync.RWMutex
}

var _ HubInterface = (*Hub)(nil)

const broadCastChan = "wsBroadcast"

// NewHub Hub を生成して返す。
func NewHub(pubsub pubsub.Service) *Hub {
	return &Hub{
		Clients:       make(map[uuid.UUID][]*Client),
		RegisterCh:    make(chan *Client),
		UnRegisterCh:  make(chan *Client),
		OnlineMembers: make(map[uuid.UUID][]uuid.UUID),
		BroadcastCh:   make(chan *Payload, maxBroadcastSize),
		pubsub:        pubsub,
		mu:            &sync.RWMutex{},
		onlineMu:      &sync.RWMutex{},
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
		h.broadCastToClient([]byte(msg.Payload))
	}
}

// Dispatch メッセージを送信する。
func (h *Hub) Dispatch(eventType EventType, targets Targets, data any) {
	payload := Payload{
		EventType: eventType,
		Targets:   targets,
		Data:      data,
	}
	h.BroadcastCh <- &payload
}

// RegisterClient クライアントを登録する。
func (h *Hub) RegisterClient(c *Client) {
	h.RegisterCh <- c
}

// UnRegisterClient クライアントの登録解除用のチャネルを返す。
func (h *Hub) UnRegisterClient(c *Client) {
	h.UnRegisterCh <- c
}

func (h *Hub) publishMessage(ctx context.Context, msg *Payload) {
	bodyBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws: publish message marshal error: %v", err)
		return
	}
	h.pubsub.Publish(ctx, broadCastChan, bodyBytes)
}

// GetOnlineMembers オンラインメンバーを取得する。
func (h *Hub) GetOnlineMembers() []ConnectingMember {
	conMembers := make([]ConnectingMember, 0, len(h.OnlineMembers))
	for memberID, connectIDs := range h.OnlineMembers {
		conMembers = append(conMembers, ConnectingMember{
			MemberID:   memberID,
			ConnectIDs: connectIDs,
		})
	}
	return conMembers
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	h.onlineMu.RLock()
	defer h.mu.Unlock()
	defer h.onlineMu.RUnlock()
	h.Clients[c.auth.MemberID] = append(h.Clients[c.auth.MemberID], c)
	conMembers := h.GetOnlineMembers()
	h.Dispatch(EventTypeConnectingMembers, Targets{
		Members: []uuid.UUID{c.auth.MemberID},
	}, ConnectingMembersEventData{
		ConnectingMembers: conMembers,
	})
	h.Dispatch(EventTypeConnected, Targets{All: true}, ConnectedEventData{
		AuthMemberID: c.auth.MemberID,
		ConnectID:    c.uid,
	})
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	v, ok := h.Clients[c.auth.MemberID]
	if !ok {
		return
	}
	for i, client := range v {
		if client.uid == c.uid {
			h.Clients[c.auth.MemberID] = append(v[:i], v[i+1:]...)
			break
		}
	}
	if len(h.Clients[c.auth.MemberID]) == 0 {
		delete(h.Clients, c.auth.MemberID)
	}
	h.Dispatch(EventTypeDisconnected, Targets{All: true}, DisconnectedEventData{
		AuthMemberID: c.auth.MemberID,
		ConnectID:    c.uid,
	})
}

func (h *Hub) broadCastToClient(msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var payload Payload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("ws: broadcast message unmarshal error: %v", err)
		return
	}
	var data any
	switch payload.EventType {
	case EventTypeConnectingMembers:
		var typedPayload TypedPayload[ConnectingMembersEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeConnected:
		var typedPayload TypedPayload[ConnectedEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		h.addOnlineMemberProcessing(typedPayload.Data)
		data = typedPayload.Data
	case EventTypeDisconnected:
		var typedPayload TypedPayload[DisconnectedEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		h.removeOnlineMemberProcessing(typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomAddedMe:
		var typedPayload TypedPayload[ChatRoomAddedMeEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomAddedMember:
		var typedPayload TypedPayload[ChatRoomAddedMemberEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomRemovedMe:
		var typedPayload TypedPayload[ChatRoomRemovedMeEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomRemovedMember:
		var typedPayload TypedPayload[ChatRoomRemovedMemberEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		data = typedPayload.Data
	case EventTypeChatRoomWithdrawnMe:
		var typedPayload TypedPayload[ChatRoomWithdrawnMeEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomWithdrawnMember:
		var typedPayload TypedPayload[ChatRoomWithdrawnMemberEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomUpdatedName:
		var typedPayload TypedPayload[ChatRoomUpdatedNameEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomDeletedMessage:
		var typedPayload TypedPayload[ChatRoomDeletedMessageEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomEditedMessage:
		var typedPayload TypedPayload[ChatRoomEditedMessageEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomSentMessage:
		var typedPayload TypedPayload[ChatRoomSentMessageEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomReadMessage:
		var typedPayload TypedPayload[ChatRoomReadMessageEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	case EventTypeChatRoomDeleted:
		var typedPayload TypedPayload[ChatRoomDeletedEventData]
		if err := json.Unmarshal(msg, &typedPayload); err != nil {
			log.Printf("ws: broadcast message unmarshal error: %v", err)
			return
		}
		setEmptySlicesToEmptyArrays(&typedPayload.Data)
		data = typedPayload.Data
	}
	res := ResponsePayload{
		EventType: payload.EventType,
		Data:      data,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("ws: broadcast message marshal error: %v", err)
		return
	}
	if payload.Targets.All {
		for _, cs := range h.Clients {
			for _, c := range cs {
				c.sendCh <- resBytes
			}
		}
		return
	}
	for _, memberID := range payload.Targets.Members {
		cs, ok := h.Clients[memberID]
		if !ok {
			continue
		}
		for _, c := range cs {
			c.sendCh <- resBytes
		}
	}
}

func (h *Hub) addOnlineMemberProcessing(data any) {
	h.onlineMu.Lock()
	defer h.onlineMu.Unlock()
	connectedData, ok := data.(ConnectedEventData)
	if !ok {
		return
	}
	memberID := connectedData.AuthMemberID
	connectID := connectedData.ConnectID
	h.OnlineMembers[memberID] = append(h.OnlineMembers[memberID], connectID)
}

func (h *Hub) removeOnlineMemberProcessing(data any) {
	h.onlineMu.Lock()
	defer h.onlineMu.Unlock()
	disconnectedData, ok := data.(DisconnectedEventData)
	if !ok {
		return
	}
	memberID := disconnectedData.AuthMemberID
	connectID := disconnectedData.ConnectID
	v, ok := h.OnlineMembers[memberID]
	if !ok {
		return
	}
	for i, id := range v {
		if id == connectID {
			h.OnlineMembers[memberID] = append(v[:i], v[i+1:]...)
			break
		}
	}
	if len(h.OnlineMembers[memberID]) == 0 {
		delete(h.OnlineMembers, memberID)
	}
}

func setEmptySlicesToEmptyArrays(val any) {
	v := reflect.ValueOf(val).Elem()
	setEmptySlicesToEmptyArraysRecursive(v)
}

func setEmptySlicesToEmptyArraysRecursive(v reflect.Value) {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.Kind() == reflect.Slice && field.IsNil() {
				newSlice := reflect.MakeSlice(field.Type(), 0, 0)
				field.Set(newSlice)
			} else if field.Kind() == reflect.Ptr && !field.IsNil() {
				setEmptySlicesToEmptyArraysRecursive(field.Elem())
			} else if field.Kind() == reflect.Struct {
				setEmptySlicesToEmptyArraysRecursive(field)
			}
		}
	}
}

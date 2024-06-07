// Package pubsub provides a pubsub application.
package pubsub

import (
	"context"
)

// Message メッセージを表す構造体。
type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}

// Service パブサブサービスを表すインターフェース。
type Service interface {
	Publish(ctx context.Context, channel string, payload any)
	Subscribe(ctx context.Context, channel string) <-chan *Message
}

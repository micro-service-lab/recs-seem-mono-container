package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

// RedisService RedisPubSubService を表す構造体。
type RedisService struct {
	client *redis.Client
}

var _ Service = (*RedisService)(nil)

// NewRedisService RedisPubSubService を生成して返す。
func NewRedisService(cfg config.Config) *RedisService {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &RedisService{
		client: cli,
	}
}

// Publish メッセージを送信する。
func (s *RedisService) Publish(ctx context.Context, channel string, payload any) {
	s.client.Publish(ctx, channel, payload)
}

// Subscribe メッセージを受信する。
func (s *RedisService) Subscribe(ctx context.Context, channel string) <-chan *Message {
	chmsg := make(chan *Message)

	go func() {
		defer close(chmsg)

		for msg := range s.client.Subscribe(ctx, channel).Channel() {
			chmsg <- &Message{
				Channel:      msg.Channel,
				Pattern:      msg.Pattern,
				Payload:      msg.Payload,
				PayloadSlice: msg.PayloadSlice,
			}
		}
	}()

	return chmsg
}

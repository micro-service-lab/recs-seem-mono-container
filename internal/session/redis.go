package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

type sessionRedisKey string

const (
	// MemberSession セッション情報を保存する Redis のキー。
	MemberSession sessionRedisKey = "member_session"
	// MemberRefreshSession リフレッシュセッション情報を保存する Redis のキー。
	MemberRefreshSession sessionRedisKey = "member_refresh_session"
)

// RedisManager Redis によるセッション管理を提供する。(memberIDごとのハッシュ管理)
type RedisManager struct {
	redisClient *redis.Client
}

var _ Manager = (*RedisManager)(nil)

// NewRedisManager RedisManager を生成して返す。
func NewRedisManager(ctx context.Context, cfg config.Config) (*RedisManager, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return &RedisManager{
		redisClient: cli,
	}, nil
}

// CheckSession セッションをチェックする。
func (m *RedisManager) CheckSession(
	ctx context.Context, memberID uuid.UUID, sessionID string,
) (bool, error) {
	res, err := m.redisClient.HGet(ctx, string(MemberSession), memberID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("hget: %w", err)
	}
	return res == sessionID, nil
}

// UpdateSession セッションを更新する。
func (m *RedisManager) UpdateSession(
	ctx context.Context, memberID uuid.UUID, sessionID string,
) error {
	// すでにある場合は上書きされる
	if err := m.redisClient.HSet(ctx, string(MemberSession), memberID.String(), sessionID).Err(); err != nil {
		return fmt.Errorf("hset: %w", err)
	}
	return nil
}

// DeleteSession セッションを削除する。
func (m *RedisManager) DeleteSession(ctx context.Context, memberID uuid.UUID) error {
	if err := m.redisClient.HDel(ctx, string(MemberSession), memberID.String()).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("hdel: %w", err)
	}
	return nil
}

// CheckRefreshSession リフレッシュセッションをチェックする。
func (m *RedisManager) CheckRefreshSession(
	ctx context.Context, memberID uuid.UUID, sessionID string,
) (bool, error) {
	res, err := m.redisClient.HGet(ctx, string(MemberRefreshSession), memberID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("hget: %w", err)
	}
	return res == sessionID, nil
}

// UpdateRefreshSession リフレッシュセッションを更新する。
func (m *RedisManager) UpdateRefreshSession(
	ctx context.Context, memberID uuid.UUID, sessionID string,
) error {
	// すでにある場合は上書きされる
	if err := m.redisClient.HSet(ctx, string(MemberRefreshSession), memberID.String(), sessionID).Err(); err != nil {
		return fmt.Errorf("hset: %w", err)
	}
	return nil
}

// DeleteRefreshSession リフレッシュセッションを削除する。
func (m *RedisManager) DeleteRefreshSession(ctx context.Context, memberID uuid.UUID) error {
	if err := m.redisClient.HDel(ctx, string(MemberRefreshSession), memberID.String()).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("hdel: %w", err)
	}
	return nil
}

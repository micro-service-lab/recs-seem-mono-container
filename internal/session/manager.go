package session

import (
	"context"

	"github.com/google/uuid"
)

// Manager セッション管理を提供するインターフェース。
type Manager interface {
	// CheckSession セッションをチェックする。
	CheckSession(ctx context.Context, memberID uuid.UUID, sessionID string) (bool, error)
	// UpdateSession セッションを更新する。
	UpdateSession(ctx context.Context, memberID uuid.UUID, sessionID string) error
	// DeleteSession セッションを削除する。
	DeleteSession(ctx context.Context, memberID uuid.UUID) error
	// CheckRefreshSession リフレッシュセッションをチェックする。
	CheckRefreshSession(ctx context.Context, memberID uuid.UUID, sessionID string) (bool, error)
	// UpdateRefreshSession リフレッシュセッションを更新する。
	UpdateRefreshSession(ctx context.Context, memberID uuid.UUID, sessionID string) error
	// DeleteRefreshSession リフレッシュセッションを削除する。
	DeleteRefreshSession(ctx context.Context, memberID uuid.UUID) error
}

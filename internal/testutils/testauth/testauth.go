// Package testauth provides a mock implementation of auth.Auth for testing.
package testauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

var _ auth.Auth = (*MockAuth)(nil)

// MockAuth auth.Auth の mock 実装。
type MockAuth struct {
	sessionID string
}

// New 新しい MockAuth を生成して返す。
// MockAuth が生成するセッション ID は常に sessionID の値となる。
func New(sessionID string) *MockAuth {
	return &MockAuth{
		sessionID: sessionID,
	}
}

// ParseToken トークン文字列をパースしてセッション情報を返す。
func (s *MockAuth) ParseToken(tokenString string, now time.Time) (*session.Session, error) {
	sess := &session.Session{}
	if err := json.Unmarshal([]byte(tokenString), sess); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	if sess.ID != s.sessionID {
		return nil, errors.New("invalid token")
	}

	if err := sess.Validate(now); err != nil {
		return nil, fmt.Errorf("validate session: %w", err)
	}

	return sess, nil
}

// NewSessionToken セッション ID とトークン文字列を生成する。
// セッション ID は常に固定となる。
func (s *MockAuth) NewSessionToken(
	clientType session.ClientType,
	clientID string,
	now time.Time,
) (string, string, error) {
	sess := &session.Session{
		ID:         s.sessionID,
		ClientType: clientType,
		ClientID:   clientID,
		AuthedAt:   now,
	}

	token, err := json.Marshal(sess)
	if err != nil {
		return "", "", fmt.Errorf("marshal json: %w", err)
	}

	return s.sessionID, string(token), nil
}

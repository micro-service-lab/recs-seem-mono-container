// Package session user session information.
package session

import (
	"errors"
	"fmt"
	"time"
)

// Session セッション情報を表す構造体。
type Session struct {
	ID         string
	ClientType ClientType
	ClientID   string
	AuthedAt   time.Time
}

// Validate セッション情報を検証する。
func (s *Session) Validate(now time.Time) error {
	if s == nil {
		return errors.New("session is nil")
	}

	if s.ID == "" {
		return errors.New("session id is empty")
	}

	if !s.ClientType.IsValid() {
		return errors.New("invalid client type")
	}

	if s.ClientID == "" {
		return errors.New("client id is empty")
	}
	if err := s.validateClientID(); err != nil {
		return fmt.Errorf("validate client id: %w", err)
	}

	if s.AuthedAt.IsZero() {
		return errors.New("auth time is empty")
	}
	if s.AuthedAt.After(now) {
		return errors.New("invalid auth time")
	}

	return nil
}

// validateClientID クライアント ID を検証する。
func (s *Session) validateClientID() error {
	if s.ClientType == ClientTypeUser && FromClientID(s.ClientID) == 0 {
		return errors.New("invalid user id")
	}
	return nil
}

// AdminID 管理者アカウント ID を返す。
// クライアント種別が管理者アカウントでない場合は空文字を返す。
func (s *Session) AdminID() string {
	if s.ClientType != ClientTypeAdmin {
		return ""
	}
	return s.ClientID
}

// UserID ユーザー ID を返す。
// クライアント種別がユーザーでない場合やクライアント ID が不正な場合は 0 を返す。
func (s *Session) UserID() int32 {
	return FromClientID(s.ClientID)
}

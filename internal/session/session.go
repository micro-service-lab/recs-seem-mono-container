// Package session user session information.
package session

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// Session セッション情報を表す構造体。
type Session struct {
	ID         string
	Type       TokenType
	MemberType MemberType
	entity.AuthPayload
	AuthedAt  time.Time
	ExpiresAt time.Time
}

// ValidateAccessToken セッション情報を検証する。
func (s *Session) ValidateAccessToken(now time.Time) error {
	if s == nil {
		return errors.New("session is nil")
	}

	if s.ID == "" {
		return errors.New("session id is empty")
	}

	if !s.MemberType.IsValid() {
		return errors.New("invalid member type")
	}

	if s.MemberID == uuid.Nil {
		return errors.New("member id is empty")
	}

	if s.AuthedAt.IsZero() {
		return errors.New("auth time is empty")
	}
	if s.AuthedAt.After(now) {
		return errors.New("invalid auth time")
	}
	if s.ExpiresAt.IsZero() {
		return errors.New("expires time is empty")
	}
	if s.ExpiresAt.Before(now) {
		return errhandle.NewCommonError(response.ExpireAccessToken, nil)
	}
	if s.Type != TokenTypeAccess {
		return errors.New("invalid token type")
	}

	return nil
}

// ValidateRefreshToken セッション情報を検証する。
func (s *Session) ValidateRefreshToken(now time.Time) error {
	if s == nil {
		return errors.New("session is nil")
	}

	if s.ID == "" {
		return errors.New("session id is empty")
	}

	if !s.MemberType.IsValid() {
		return errors.New("invalid member type")
	}

	if s.MemberID == uuid.Nil {
		return errors.New("member id is empty")
	}

	if s.AuthedAt.IsZero() {
		return errors.New("auth time is empty")
	}
	if s.AuthedAt.After(now) {
		return errors.New("invalid auth time")
	}
	if s.ExpiresAt.IsZero() {
		return errors.New("expires time is empty")
	}
	if s.ExpiresAt.Before(now) {
		return errhandle.NewCommonError(response.ExpireRefreshToken, nil)
	}
	if s.Type != TokenTypeRefresh {
		return errors.New("invalid token type")
	}

	return nil
}

// Package auth provides authentication related functions.
package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

const (
	clientTypeKey = "client_type"
)

// Auth 認証関連の機能を提供するインターフェース。
type Auth interface {
	// ParseToken トークン文字列をパースしてセッション情報を返す。
	ParseToken(tokenString string, now time.Time) (*session.Session, error)
	// NewSessionToken セッション ID とトークン文字列を生成する。
	NewSessionToken(clientType session.ClientType, clientID string, now time.Time) (string, string, error)
}

// jwtAuth JSON Web Token (JWT) による Auth 実装。
type jwtAuth struct {
	secret []byte
	issuer string
}

// New Auth を生成して返す。
func New(secret []byte, issuer string) Auth {
	return &jwtAuth{
		secret: secret,
		issuer: issuer,
	}
}

// ParseToken トークン文字列をパースしてセッション情報を返す。
func (s *jwtAuth) ParseToken(tokenString string, now time.Time) (*session.Session, error) {
	clk := jwt.ClockFunc(func() time.Time {
		return now
	})

	token, err := jwt.ParseString(tokenString,
		jwt.WithKey(jwa.HS512, s.secret), jwt.WithClock(clk), jwt.WithIssuer(s.issuer))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	sess := &session.Session{
		ID:         token.JwtID(),
		ClientType: getClientType(token),
		ClientID:   token.Subject(),
		AuthedAt:   token.IssuedAt(),
	}

	if err := sess.Validate(now); err != nil {
		return nil, fmt.Errorf("validate session: %w", err)
	}

	return sess, nil
}

// getClientType JWT からクライアント種別 (カスタムクレーム) を取得する。
func getClientType(token jwt.Token) session.ClientType {
	value, ok := token.Get(clientTypeKey)
	if !ok {
		return session.ClientTypeInvalid
	}

	str, ok := value.(string)
	if !ok {
		return session.ClientTypeInvalid
	}

	return session.ClientTypeFromString(str)
}

// NewSessionToken セッション ID とトークン文字列を生成する。
func (s *jwtAuth) NewSessionToken(
	clientType session.ClientType,
	clientID string,
	now time.Time,
) (string, string, error) {
	sessID, err := newSessionID()
	if err != nil {
		return "", "", fmt.Errorf("new session id: %w", err)
	}

	sess := &session.Session{
		ID:         sessID,
		ClientType: clientType,
		ClientID:   clientID,
		AuthedAt:   now,
	}

	if err := sess.Validate(now); err != nil {
		return "", "", fmt.Errorf("validate session: %w", err)
	}

	token, err := newToken(sess, s.issuer)
	if err != nil {
		return "", "", fmt.Errorf("new token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, s.secret))
	if err != nil {
		return "", "", fmt.Errorf("sign: %w", err)
	}

	return sessID, string(signed), nil
}

// newSessionID セッション ID を生成する。
func newSessionID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("new uuid v4: %w", err)
	}

	return u.String(), nil
}

// newToken セッション情報から jwt.Token を生成する。
func newToken(sess *session.Session, issuer string) (jwt.Token, error) {
	token := jwt.New()

	if err := token.Set(jwt.IssuedAtKey, sess.AuthedAt); err != nil {
		return nil, fmt.Errorf("set iat: %w", err)
	}

	if err := token.Set(jwt.IssuerKey, issuer); err != nil {
		return nil, fmt.Errorf("set iss: %w", err)
	}

	if err := token.Set(jwt.JwtIDKey, sess.ID); err != nil {
		return nil, fmt.Errorf("set jti: %w", err)
	}

	if err := token.Set(jwt.SubjectKey, sess.ClientID); err != nil {
		return nil, fmt.Errorf("set sub: %w", err)
	}

	if err := token.Set(clientTypeKey, sess.ClientType.String()); err != nil {
		return nil, fmt.Errorf("set client_type: %w", err)
	}

	return token, nil
}

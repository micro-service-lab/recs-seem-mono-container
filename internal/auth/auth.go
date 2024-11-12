// Package auth provides authentication related functions.
package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// AccessTokenCookieKey アクセストークンのクッキー名。
const AccessTokenCookieKey = "access-token"

const (
	memberTypeKey = "member_type"
	tokenTypeKey  = "token_type"
)

// Auth 認証関連の機能を提供するインターフェース。
type Auth interface {
	// ParseAccessToken アクセストークン文字列をパースしてセッション情報を返す。
	ParseAccessToken(tokenString string, now time.Time) (*session.Session, error)
	// ParseRefreshToken リフレッシュトークン文字列をパースしてセッション情報を返す。
	ParseRefreshToken(tokenString string, now time.Time) (*session.Session, error)
	// NewSessionToken セッション ID とトークン文字列を生成する。
	NewSessionToken(
		memberType session.MemberType, memberID uuid.UUID, now time.Time,
	) (string, string, error)
	// NewRefreshToken リフレッシュトークンを生成する。
	NewRefreshToken(
		memberType session.MemberType, memberID uuid.UUID, now time.Time,
	) (string, string, error)
}

// jwtAuth JSON Web Token (JWT) による Auth 実装。
type jwtAuth struct {
	secret                []byte
	refreshSecret         []byte
	issuer                string
	accessTokenExpiresIn  time.Duration
	refreshTokenExpiresIn time.Duration
}

// New Auth を生成して返す。
func New(
	secret []byte,
	refreshSecret []byte,
	issuer string,
	accessTokenExpiresIn,
	refreshTokenExpiresIn time.Duration,
) Auth {
	return &jwtAuth{
		secret:                secret,
		refreshSecret:         refreshSecret,
		issuer:                issuer,
		accessTokenExpiresIn:  accessTokenExpiresIn,
		refreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}

// ParseAccessToken アクセストークン文字列をパースしてセッション情報を返す。
func (s *jwtAuth) ParseAccessToken(tokenString string, now time.Time) (*session.Session, error) {
	clk := jwt.ClockFunc(func() time.Time {
		return now
	})

	token, err := jwt.ParseString(tokenString,
		jwt.WithKey(jwa.HS512, s.secret), jwt.WithClock(clk), jwt.WithIssuer(s.issuer))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	memberID, err := uuid.Parse(token.Subject())
	if err != nil {
		return nil, fmt.Errorf("parse member id: %w", err)
	}

	sess := &session.Session{
		ID:         token.JwtID(),
		MemberType: getMemberType(token),
		Type:       getTokenType(token),
		AuthPayload: entity.AuthPayload{
			MemberID: memberID,
		},
		AuthedAt:  token.IssuedAt(),
		ExpiresAt: token.Expiration(),
	}

	if err := sess.ValidateAccessToken(now); err != nil {
		return nil, fmt.Errorf("validate session: %w", err)
	}

	return sess, nil
}

// ParseRefreshToken リフレッシュトークン文字列をパースしてセッション情報を返す。
func (s *jwtAuth) ParseRefreshToken(tokenString string, now time.Time) (*session.Session, error) {
	clk := jwt.ClockFunc(func() time.Time {
		return now
	})

	token, err := jwt.ParseString(tokenString,
		jwt.WithKey(jwa.HS512, s.refreshSecret), jwt.WithClock(clk), jwt.WithIssuer(s.issuer))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	memberID, err := uuid.Parse(token.Subject())
	if err != nil {
		return nil, fmt.Errorf("parse member id: %w", err)
	}

	sess := &session.Session{
		ID:         token.JwtID(),
		MemberType: getMemberType(token),
		Type:       getTokenType(token),
		AuthPayload: entity.AuthPayload{
			MemberID: memberID,
		},
		AuthedAt:  token.IssuedAt(),
		ExpiresAt: token.Expiration(),
	}

	if err := sess.ValidateRefreshToken(now); err != nil {
		return nil, fmt.Errorf("validate session: %w", err)
	}

	return sess, nil
}

// getMemberType JWT からクライアント種別 (カスタムクレーム) を取得する。
func getMemberType(token jwt.Token) session.MemberType {
	value, ok := token.Get(memberTypeKey)
	if !ok {
		return session.MemberTypeInvalid
	}

	str, ok := value.(string)
	if !ok {
		return session.MemberTypeInvalid
	}

	return session.MemberTypeFromString(str)
}

// getTokenType JWT からトークン種別 (カスタムクレーム) を取得する。
func getTokenType(token jwt.Token) session.TokenType {
	value, ok := token.Get(tokenTypeKey)
	if !ok {
		return session.TokenTypeInvalid
	}

	str, ok := value.(string)
	if !ok {
		return session.TokenTypeInvalid
	}

	return session.TokenTypeFromString(str)
}

// NewSessionToken セッション ID とトークン文字列を生成する。
func (s *jwtAuth) NewSessionToken(
	memberType session.MemberType,
	memberID uuid.UUID,
	now time.Time,
) (string, string, error) {
	sessID, err := newSessionID()
	if err != nil {
		return "", "", fmt.Errorf("new session id: %w", err)
	}

	sess := &session.Session{
		ID:         sessID,
		MemberType: memberType,
		AuthPayload: entity.AuthPayload{
			MemberID: memberID,
		},
		Type:      session.TokenTypeAccess,
		AuthedAt:  now,
		ExpiresAt: now.Add(s.accessTokenExpiresIn),
	}

	if err := sess.ValidateAccessToken(now); err != nil {
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

// NewRefreshToken リフレッシュトークンを生成する。
func (s *jwtAuth) NewRefreshToken(
	memberType session.MemberType,
	memberID uuid.UUID,
	now time.Time,
) (string, string, error) {
	sessID, err := newSessionID()
	if err != nil {
		return "", "", fmt.Errorf("new session id: %w", err)
	}

	sess := &session.Session{
		ID:         sessID,
		MemberType: memberType,
		AuthPayload: entity.AuthPayload{
			MemberID: memberID,
		},
		Type:      session.TokenTypeRefresh,
		AuthedAt:  now,
		ExpiresAt: now.Add(s.refreshTokenExpiresIn),
	}

	if err := sess.ValidateRefreshToken(now); err != nil {
		return "", "", fmt.Errorf("validate session: %w", err)
	}

	token, err := newToken(sess, s.issuer)
	if err != nil {
		return "", "", fmt.Errorf("new token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, s.refreshSecret))
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

	if err := token.Set(jwt.ExpirationKey, sess.ExpiresAt); err != nil {
		return nil, fmt.Errorf("set exp: %w", err)
	}

	if err := token.Set(jwt.IssuerKey, issuer); err != nil {
		return nil, fmt.Errorf("set iss: %w", err)
	}

	if err := token.Set(jwt.JwtIDKey, sess.ID); err != nil {
		return nil, fmt.Errorf("set jti: %w", err)
	}

	if err := token.Set(jwt.SubjectKey, sess.MemberID.String()); err != nil {
		return nil, fmt.Errorf("set sub: %w", err)
	}

	if err := token.Set(tokenTypeKey, sess.Type.String()); err != nil {
		return nil, fmt.Errorf("set token_type: %w", err)
	}

	if err := token.Set(memberTypeKey, sess.MemberType.String()); err != nil {
		return nil, fmt.Errorf("set member_type: %w", err)
	}

	return token, nil
}

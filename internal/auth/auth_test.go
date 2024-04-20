package auth_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

func TestParseToken(t *testing.T) {
	secret := []byte("secret")
	issuer := "test_issuer"
	authService := auth.New(secret, issuer)

	authed := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name            string
		newSessionToken func() (string, string, error) // セッション ID とトークンを返す関数
		now             time.Time
		out             *session.Session
		failed          bool
	}{
		{
			name: "valid admin token",
			newSessionToken: func() (string, string, error) {
				return authService.NewSessionToken(session.ClientTypeAdmin, "admin", authed)
			},
			now: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			out: &session.Session{
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   authed,
			},
		},
		{
			name: "valid user token",
			newSessionToken: func() (string, string, error) {
				return authService.NewSessionToken(session.ClientTypeUser, "12345", authed)
			},
			now: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			out: &session.Session{
				ClientType: session.ClientTypeUser,
				ClientID:   "12345",
				AuthedAt:   authed,
			},
		},
		{
			name: "invalid token (invalid alg)",
			newSessionToken: func() (string, string, error) {
				sessID := "session-id"
				token := jwt.New()
				if err := token.Set(jwt.IssuedAtKey, authed); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.IssuerKey, "test_issuer"); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.JwtIDKey, sessID); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.SubjectKey, "admin"); err != nil {
					return "", "", err
				}
				if err := token.Set("client_type", session.ClientTypeAdmin.String()); err != nil {
					return "", "", err
				}
				signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, secret))
				if err != nil {
					return "", "", err
				}
				return sessID, string(signed), nil
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (signed with invalid secret)",
			newSessionToken: func() (string, string, error) {
				return auth.New([]byte("invalid"), "test_issuer").NewSessionToken(session.ClientTypeAdmin, "admin", authed)
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (issuer not match)",
			newSessionToken: func() (string, string, error) {
				return auth.New(secret, "invalid_issuer").NewSessionToken(session.ClientTypeAdmin, "admin", authed)
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (invalid iss)",
			newSessionToken: func() (string, string, error) {
				sessID := "session-id"
				token := jwt.New()
				if err := token.Set(jwt.IssuedAtKey, authed); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.IssuerKey, "invalid"); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.JwtIDKey, sessID); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.SubjectKey, "admin"); err != nil {
					return "", "", err
				}
				if err := token.Set("client_type", session.ClientTypeAdmin.String()); err != nil {
					return "", "", err
				}
				signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, secret))
				if err != nil {
					return "", "", err
				}
				return sessID, string(signed), nil
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (invalid sub)",
			newSessionToken: func() (string, string, error) {
				sessID := "session-id"
				token := jwt.New()
				if err := token.Set(jwt.IssuedAtKey, authed); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.IssuerKey, "test_issuer"); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.JwtIDKey, sessID); err != nil {
					return "", "", err
				}
				if err := token.Set("client_type", session.ClientTypeAdmin.String()); err != nil {
					return "", "", err
				}
				signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, secret))
				if err != nil {
					return "", "", err
				}
				return sessID, string(signed), nil
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (invalid iat)",
			newSessionToken: func() (string, string, error) {
				return authService.NewSessionToken(session.ClientTypeAdmin, "admin", authed)
			},
			now:    time.Date(2022, 1, 1, 11, 59, 59, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (invalid jti)",
			newSessionToken: func() (string, string, error) {
				sessID := "session-id"
				token := jwt.New()
				if err := token.Set(jwt.IssuedAtKey, authed); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.IssuerKey, "test_issuer"); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.SubjectKey, "admin"); err != nil {
					return "", "", err
				}
				if err := token.Set("client_type", session.ClientTypeAdmin.String()); err != nil {
					return "", "", err
				}
				signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, secret))
				if err != nil {
					return "", "", err
				}
				return sessID, string(signed), nil
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
		{
			name: "invalid token (invalid client_type)",
			newSessionToken: func() (string, string, error) {
				sessID := "session-id"
				token := jwt.New()
				if err := token.Set(jwt.IssuedAtKey, authed); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.IssuerKey, "test_issuer"); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.JwtIDKey, sessID); err != nil {
					return "", "", err
				}
				if err := token.Set(jwt.SubjectKey, "admin"); err != nil {
					return "", "", err
				}
				signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, secret))
				if err != nil {
					return "", "", err
				}
				return sessID, string(signed), nil
			},
			now:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed: true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			sessID, token, err := v.newSessionToken()
			if err != nil {
				tt.Fatalf("failed to generate token: %+v", err)
			}

			out, err := authService.ParseToken(token, v.now)
			switch {
			case err != nil && !v.failed:
				tt.Fatalf("unexpected error: %+v", err)
			case err == nil && v.failed:
				tt.Fatal("unexpected success")
			case err != nil && v.failed:
				// pass
				tt.Logf("expected error: %+v", err)
				return
			}

			v.out.ID = sessID
			if diff := cmp.Diff(v.out, out); diff != "" {
				tt.Errorf("unexpected result:\n%s", diff)
			}
		})
	}
}

func TestNewSessionToken(t *testing.T) {
	secret := []byte("secret")
	issuer := "test_issuer"
	authService := auth.New(secret, issuer)

	cases := []struct {
		name       string
		clientType session.ClientType
		clientID   string
		now        time.Time
		failed     bool
	}{
		{
			name:       "admin",
			clientType: session.ClientTypeAdmin,
			clientID:   "admin",
			now:        time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "user",
			clientType: session.ClientTypeUser,
			clientID:   "12345",
			now:        time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "invalid client",
			clientType: session.ClientTypeInvalid,
			clientID:   "admin",
			now:        time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed:     true,
		},
		{
			name:       "missing client id",
			clientType: session.ClientTypeAdmin,
			clientID:   "",
			now:        time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			failed:     true,
		},
		{
			name:       "now is zero",
			clientType: session.ClientTypeAdmin,
			clientID:   "",
			now:        time.Time{},
			failed:     true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			sessID, token, err := authService.NewSessionToken(v.clientType, v.clientID, v.now)
			switch {
			case err != nil && !v.failed:
				tt.Fatalf("unexpected error: %+v", err)
			case err == nil && v.failed:
				tt.Fatal("unexpected success")
			case err != nil && v.failed:
				// pass
				tt.Logf("expected error: %+v", err)
				return
			}

			sess, err := authService.ParseToken(token, v.now)
			if err != nil {
				tt.Fatalf("failed to parse token: %+v", err)
			}

			if sessID != sess.ID {
				tt.Errorf("expected %q, but got %q", sessID, sess.ID)
			}
		})
	}
}

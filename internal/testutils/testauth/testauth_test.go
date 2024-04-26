package testauth_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils/testauth"
)

func TestMockAuth(t *testing.T) {
	sessionID := "session-id"
	mock := testauth.New(sessionID)

	now := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name       string
		clientType session.ClientType
		clientID   string
		sess       *session.Session
		failed     bool
	}{
		{
			name:       "valid admin token",
			clientType: session.ClientTypeAdmin,
			clientID:   "admin",
			sess: &session.Session{
				ID:         sessionID,
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   now,
			},
		},
		{
			name:       "valid user token",
			clientType: session.ClientTypeUser,
			clientID:   "12345",
			sess: &session.Session{
				ID:         sessionID,
				ClientType: session.ClientTypeUser,
				ClientID:   "12345",
				AuthedAt:   now,
			},
		},
		{
			name:       "invalid token",
			clientType: session.ClientTypeInvalid,
			clientID:   "invalid",
			failed:     true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			id, token, err := mock.NewSessionToken(v.clientType, v.clientID, now)
			if err != nil {
				tt.Fatalf("failed to generate session token: %+v", err)
			}

			if id != sessionID {
				tt.Errorf("expected %q, but got %q", sessionID, id)
			}

			sess, err := mock.ParseToken(token, now)
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

			if diff := cmp.Diff(v.sess, sess); diff != "" {
				tt.Errorf("unexpected result:\n%s", diff)
			}
		})
	}
}

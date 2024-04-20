package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

func TestContext(t *testing.T) {
	cases := []struct {
		name    string
		sess    *session.Session
		adminID string
		userID  int32
	}{
		{
			name:    "no session",
			sess:    nil,
			adminID: "",
			userID:  0,
		},
		{
			name: "admin session",
			sess: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			adminID: "admin",
			userID:  0,
		},
		{
			name: "user session",
			sess: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeUser,
				ClientID:   "12345",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			adminID: "",
			userID:  12345,
		},
		{
			name: "invalid user session",
			sess: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeUser,
				ClientID:   "invalid", // invalid id
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			adminID: "",
			userID:  0,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			ctx := context.Background()
			if v.sess != nil {
				ctx = session.NewContext(ctx, v.sess)
			}

			sess := session.FromContext(ctx)
			if diff := cmp.Diff(v.sess, sess); diff != "" {
				tt.Errorf("unexpected result:\n%s", diff)
			}

			adminID := session.AdminID(ctx)
			if adminID != v.adminID {
				tt.Errorf("expected %q, but got %q", v.adminID, adminID)
			}

			userID := session.UserID(ctx)
			if userID != v.userID {
				tt.Errorf("expected %d, but got %d", v.userID, userID)
			}
		})
	}
}

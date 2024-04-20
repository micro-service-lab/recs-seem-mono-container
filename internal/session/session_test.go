package session_test

import (
	"testing"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

func TestSession(t *testing.T) {
	now := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name        string
		in          *session.Session
		validateErr string
		adminID     string
		userID      int32
	}{
		{
			name: "valid admin session",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			adminID: "admin",
			userID:  0,
		},
		{
			name: "valid user session",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeUser,
				ClientID:   "12345",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			adminID: "",
			userID:  12345,
		},
		{
			name:        "nil session",
			in:          nil,
			validateErr: "session is nil",
		},
		{
			name: "missing session id",
			in: &session.Session{
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			validateErr: "session id is empty",
		},
		{
			name: "invalid client type",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeInvalid,
				ClientID:   "admin",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			validateErr: "invalid client type",
		},
		{
			name: "missing client id",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeAdmin,
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			validateErr: "client id is empty",
		},
		{
			name: "invalid user id",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeUser,
				ClientID:   "invalid",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			validateErr: "validate client id: invalid user id",
		},
		{
			name: "missing auth time",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
			},
			validateErr: "auth time is empty",
		},
		{
			name: "invalid auth time",
			in: &session.Session{
				ID:         "session-id",
				ClientType: session.ClientTypeAdmin,
				ClientID:   "admin",
				AuthedAt:   time.Date(2022, 1, 1, 12, 0, 1, 0, time.UTC),
			},
			validateErr: "invalid auth time",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			err := v.in.Validate(now)
			switch {
			case err != nil && v.validateErr == "":
				tt.Fatalf("unexpected error: %+v", err)
			case err == nil && v.validateErr != "":
				tt.Fatal("unexpected success")
			case err != nil && v.validateErr != "":
				if err.Error() != v.validateErr {
					tt.Fatalf("expected %q, but got %q", v.validateErr, err)
				}
				// pass
				tt.Logf("expected error: %+v", err)
				return
			}

			adminID := v.in.AdminID()
			if adminID != v.adminID {
				tt.Errorf("expected %q, but got %q", v.adminID, adminID)
			}

			userID := v.in.UserID()
			if userID != v.userID {
				tt.Errorf("expected %d, but got %d", v.userID, userID)
			}
		})
	}
}

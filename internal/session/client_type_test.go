package session_test

import (
	"testing"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

func TestClientTypeFromString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  session.ClientType
	}{
		{
			name: "admin",
			in:   "admin",
			out:  session.ClientTypeAdmin,
		},
		{
			name: "user",
			in:   "user",
			out:  session.ClientTypeUser,
		},
		{
			name: "ADMIN",
			in:   "ADMIN",
			out:  session.ClientTypeAdmin,
		},
		{
			name: "USER",
			in:   "USER",
			out:  session.ClientTypeUser,
		},
		{
			name: "empty",
			in:   "",
			out:  session.ClientTypeInvalid,
		},
		{
			name: "invalid string",
			in:   "invalid",
			out:  session.ClientTypeInvalid,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			out := session.ClientTypeFromString(v.in)
			if out != v.out {
				tt.Errorf("expected %d (%s), but got %d (%s)", v.out, v.out, out, out)
			}
		})
	}
}

func TestClientType(t *testing.T) {
	cases := []struct {
		name  string
		in    session.ClientType
		str   string
		valid bool
	}{
		{
			name:  "invalid",
			in:    session.ClientTypeInvalid,
			str:   "invalid",
			valid: false,
		},
		{
			name:  "admin",
			in:    session.ClientTypeAdmin,
			str:   "admin",
			valid: true,
		},
		{
			name:  "user",
			in:    session.ClientTypeUser,
			str:   "user",
			valid: true,
		},
		{
			name:  "unknown value",
			in:    3,
			str:   "invalid",
			valid: false,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			str := v.in.String()
			if str != v.str {
				tt.Errorf("expected %q, but got %q", v.str, str)
			}

			valid := v.in.IsValid()
			if valid != v.valid {
				tt.Errorf("expected %v, but got %v", v.valid, valid)
			}
		})
	}
}

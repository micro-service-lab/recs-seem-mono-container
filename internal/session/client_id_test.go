package session_test

import (
	"testing"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

func TestToClientID(t *testing.T) {
	cases := []struct {
		name string
		in   int32
		out  string
	}{
		{
			name: "zero",
			in:   0,
			out:  "",
		},
		{
			name: "non-zero",
			in:   12345,
			out:  "12345",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			out := session.ToClientID(v.in)
			if out != v.out {
				tt.Errorf("expected %q, but got %q", v.out, out)
			}
		})
	}
}

func TestFromClientID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  int32
	}{
		{
			name: "empty",
			in:   "",
			out:  0,
		},
		{
			name: "string id",
			in:   "admin",
			out:  0,
		},
		{
			name: "numeric id",
			in:   "12345",
			out:  12345,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			out := session.FromClientID(v.in)
			if out != v.out {
				tt.Errorf("expected %d, but got %d", v.out, out)
			}
		})
	}
}

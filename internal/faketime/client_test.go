package faketime_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/faketime"
)

func TestClient(t *testing.T) {
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	clk := fakeclock.New(now)
	api := faketime.NewAPI(clk, "")

	srv := httptest.NewServer(api)
	defer srv.Close()

	client := faketime.NewClient(http.DefaultClient, srv.URL)

	ctx := context.Background()

	t.Run("get", func(tt *testing.T) {
		ts, err := client.Get(ctx)
		if err != nil {
			tt.Fatalf("unexpected error: %+v", err)
		}

		if !ts.Equal(now) {
			tt.Errorf("expected %q, but got %q", now, ts)
		}
	})

	now = now.Add(5 * time.Minute)

	t.Run("set", func(tt *testing.T) {
		err := client.Set(ctx, now)
		if err != nil {
			tt.Fatalf("unexpected error: %+v", err)
		}
	})

	t.Run("re-get", func(tt *testing.T) {
		ts, err := client.Get(ctx)
		if err != nil {
			tt.Fatalf("unexpected error: %+v", err)
		}

		if !ts.Equal(now) {
			tt.Errorf("expected %q, but got %q", now, ts)
		}
	})
}

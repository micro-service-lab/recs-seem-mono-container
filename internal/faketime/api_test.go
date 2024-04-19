package faketime_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/faketime"
)

func TestAPI(t *testing.T) {
	cases := []struct {
		name     string
		bathPath string
	}{
		{
			name: "base path is empty",
		},
		{
			name:     "base path is /",
			bathPath: "/",
		},
		{
			name:     "with base path",
			bathPath: "/prefix",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			clk := fakeclock.New(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC))
			api := faketime.NewAPI(clk, v.bathPath)

			getPath := path.Join(v.bathPath, "/get")
			setPath := path.Join(v.bathPath, "/set")

			tt.Run("get", func(ttt *testing.T) {
				r := httptest.NewRequest(http.MethodGet, getPath, nil)
				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusOK, "2023-01-01T12:00:00Z")
			})

			tt.Run("set", func(ttt *testing.T) {
				r := httptest.NewRequest(http.MethodGet, setPath, nil)
				q := r.URL.Query()
				q.Add("t", time.Date(2023, 1, 31, 12, 0, 0, 0, time.UTC).Format(time.RFC3339))
				r.URL.RawQuery = q.Encode()

				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusOK, "ok")
			})

			tt.Run("set (invalid time)", func(ttt *testing.T) {
				r := httptest.NewRequest(http.MethodGet, setPath, nil)
				q := r.URL.Query()
				q.Add("t", "invalid")
				r.URL.RawQuery = q.Encode()

				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusBadRequest, "Bad Request")
			})

			tt.Run("re-get", func(ttt *testing.T) {
				r := httptest.NewRequest(http.MethodGet, getPath, nil)
				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusOK, "2023-01-31T12:00:00Z")
			})

			tt.Run("invalid method", func(ttt *testing.T) {
				r := httptest.NewRequest(http.MethodPost, getPath, nil)
				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusForbidden, "Forbidden")
			})

			tt.Run("invalid path", func(ttt *testing.T) {
				targetPath := path.Join(v.bathPath, "/invalid")
				r := httptest.NewRequest(http.MethodGet, targetPath, nil)
				w := httptest.NewRecorder()

				api.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				checkResponse(ttt, resp, http.StatusForbidden, "Forbidden")
			})
		})
	}
}

func checkResponse(t *testing.T, resp *http.Response, code int, body string) {
	t.Helper()

	t.Run("status code", func(tt *testing.T) {
		if resp.StatusCode != code {
			tt.Errorf("expected %d, but got %d", code, resp.StatusCode)
		}
	})

	t.Run("content-type", func(tt *testing.T) {
		contentType := resp.Header.Get("Content-Type")

		expected := "text/plain; charset=utf-8"
		if contentType != expected {
			tt.Errorf("expected %q, but got %q", expected, contentType)
		}
	})

	t.Run("body", func(tt *testing.T) {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			tt.Fatalf("failed to read body: %+v", err)
		}

		if string(b) != body {
			tt.Errorf("expected %q, but got %q", body, b)
		}
	})
}

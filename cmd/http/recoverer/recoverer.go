// Package recoverer provides a middleware that recovers from panics,
// logs the panic (and a backtrace), and returns a Custom status if possible.(based on go-chi/chi)
package recoverer

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
)

// Options is a struct for specifying configuration options for the middleware.
type Options struct {
	Handler func(w http.ResponseWriter, r *http.Request, cause any)
}

var defaultOptions = Options{
	Handler: func(w http.ResponseWriter, _ *http.Request, _ any) {
		w.WriteHeader(http.StatusInternalServerError)
	},
}

// New is a middleware that recovers from panics, uses the default options.
func New(next http.Handler) http.Handler {
	return customRecoverer(next, defaultOptions)
}

// NewWithOpts is a middleware that recovers from panics, uses the provided options.
func NewWithOpts(opts Options) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return customRecoverer(next, opts)
	}
}

func customRecoverer(next http.Handler, opts Options) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				//nolint: errorlint
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				if r.Header.Get("Connection") != "Upgrade" {
					opts.Handler(w, r, rvr)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

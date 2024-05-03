// Package main This application is a monolith recs-seem service.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/app"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/api"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/cors"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/faketime"
)

const (
	// corsMaxAge プリフライトリクエストをキャッシュできる時間 (秒)。
	corsMaxAge = 600
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("failed to run: %+v", err)
	}
}

func run(ctx context.Context) error {
	ctr := app.NewContainer()

	if err := ctr.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.Default(),
		NoColor: runtime.GOOS == "windows",
	}))

	if ctr.Config.FakeTime.Enabled {
		log.Println("Fake time mode is enabled")

		fakeClk := fakeclock.New(ctr.Config.FakeTime.Time)

		r.Mount("/faketime", faketime.NewAPI(fakeClk, "/faketime"))
	}

	auth := auth.New([]byte(ctr.Config.AuthSecret), ctr.Config.SecretIssuer)

	apiI := api.NewAPI(ctr.Clocker, auth, ctr.ServiceManager)

	middlewares := make([]func(http.Handler) http.Handler, 0, 2) //nolint:gomnd
	if ctr.Config.ClientOrigin != nil && len(ctr.Config.ClientOrigin) > 0 {
		log.Println("CORS is enabled")
		middlewares = append(middlewares, cors.Handler(cors.Options{
			AllowedOrigins: []string(ctr.Config.ClientOrigin),
			AllowedMethods: []string{http.MethodPost, http.MethodGet},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-Id", "Accept", "X-CSRF-Token"},
			MaxAge:         corsMaxAge,
			ErrorHandler: func(w http.ResponseWriter, _ *http.Request, c cors.Cors, err error) bool {
				_, ok := err.(cors.Error)
				if ok {
					c.Log.Printf("CORS error: %v", err)
					res := struct {
						Message string `json:"message"`
					}{
						Message: "CORS error: " + err.Error(),
					}
					w.Header().Set("Content-Type", "application/json")
					noOrigin := false
					switch {
					case errors.Is(err, &cors.PreflightEmptyOriginError{}):
						fallthrough
					case errors.Is(err, &cors.ActualMissingOriginError{}):
						noOrigin = true
					case errors.Is(err, &cors.PreflightNotOptionMethodError{}):
						fallthrough
					case errors.Is(err, &cors.PreflightNotAllowedMethodError{}):
						fallthrough
					case errors.Is(err, &cors.ActualMethodNotAllowedError{}):
						w.WriteHeader(http.StatusMethodNotAllowed)
					default:
						w.WriteHeader(http.StatusForbidden)
					}
					// For requests that do not conform to the browser's same-origin policy (no Origin header,
					// such as postman, is given), pass through processing.
					if noOrigin {
						return true
					}
					if err := json.NewEncoder(w).Encode(res); err != nil {
						c.Log.Printf("CORS error encoding failed: %v", err)
					}
					return false
				}
				res := struct {
					Message string `json:"message"`
				}{
					Message: "CORS error: An unexpected error has occurred",
				}
				if err := json.NewEncoder(w).Encode(res); err != nil {
					c.Log.Printf("CORS error encoding failed: %v", err)
				}
				return false
			},
			// AllowCredentials: true, // jwtをクッキーに保存する場合はtrueにする
			Debug: ctr.Config.DebugCORS,
		}))
	}
	// middlewares = append(middlewares, app.AuthMiddleware(time.Now, auth, db, app.DefaultAPIBasePath))
	apiI.Use(middlewares...)

	srvApp := api.NewApp(apiI, &api.AppOptions{})
	r.Mount("/", srvApp.Handler())

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctr.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	srv := NewServer(listener, r, ctr.Close)
	return srv.Run(ctx)
}

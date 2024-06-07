// Package main This application is a monolith recs-seem service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	httprateredis "github.com/go-chi/httprate-redis"

	"github.com/micro-service-lab/recs-seem-mono-container/app"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/api"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/cors"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/recoverer"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/faketime"
)

func main() {
	log.Default()
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

	go ctr.WebsocketHub.SubscribeMessages(ctx)
	go ctr.WebsocketHub.RunLoop(ctx)

	r := chi.NewRouter()
	// TODO: slog に変更する
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.Default(),
		NoColor: runtime.GOOS == "windows",
	}))
	r.Use(recoverer.NewWithOpts(recoverer.Options{
		Handler: recoverer.Handler,
	}))
	r.Use(middleware.CleanPath)

	if ctr.Config.FakeTime.Enabled {
		log.Println("Fake time mode is enabled")

		fakeClk := fakeclock.New(ctr.Config.FakeTime.Time)

		r.Mount("/faketime", faketime.NewAPI(fakeClk, "/faketime"))
	}

	apiI := api.NewAPI(
		ctr.Clocker,
		ctr.Auth,
		ctr.Validator,
		ctr.ServiceManager,
		ctr.Translator,
		ctr.SessionManager,
		*ctr.Config,
		ctr.WebsocketHub,
	)

	middlewares := make([]func(http.Handler) http.Handler, 0, 3) //nolint:gomnd
	// CORS ミドルウェアを追加
	if ctr.Config.ClientOrigin != nil && len(ctr.Config.ClientOrigin) > 0 {
		log.Println("CORS is enabled")
		middlewares = append(middlewares, cors.Handler(cors.Options{
			AllowedOrigins: []string(ctr.Config.ClientOrigin),
			AllowedMethods: []string{http.MethodPost, http.MethodGet},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"X-Request-Id",
				"Accept",
				"X-CSRF-Token",
				"X-Requested-With",
			},
			MaxAge:           ctr.Config.CORSMaxAge,
			ErrorHandler:     cors.AppHandler,
			AllowCredentials: true, // jwtをクッキーに保存する場合はtrueにする
			Debug:            ctr.Config.DebugCORS,
		}))
	}
	middlewares = append(middlewares, lang.Handler(string(ctr.Config.DefaultLanguage)))
	middlewares = append(middlewares, httprate.Limit(
		ctr.Config.ThrottleRequestLimit,
		1*time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			limit := w.Header().Get("X-RateLimit-Limit")
			reset := w.Header().Get("X-RateLimit-Reset")
			retryAfter := w.Header().Get("Retry-After")
			errAttr := map[string]any{
				"limit":      limit,
				"resetTime":  reset,
				"retryAfter": retryAfter,
			}
			if err := response.JSONResponseWriter(r.Context(), w, response.ThrottleRequests, nil, errAttr); err != nil {
				log.Printf("[ERROR] response writing failed: %+v", err)
			}
		}),
		httprateredis.WithRedisLimitCounter(&httprateredis.Config{
			Host:     ctr.Config.RedisHost,
			Port:     ctr.Config.RedisPort,
			Password: ctr.Config.RedisPassword,
			DBIndex:  ctr.Config.RedisDB,
		}),
	))
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

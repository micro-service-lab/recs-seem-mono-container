// Package main This application is a monolith recs-seem service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/app"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/api"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/faketime"
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
	// if cfg.ClientOrigin != "" {
	// 	log.Println("CORS is enabled")
	// 	middlewares = append(middlewares, cors.Handler(cors.Options{
	// 		AllowedOrigins: []string{cfg.ClientOrigin},
	// 		AllowedMethods: []string{http.MethodPost},
	// 		AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-Id"},
	// 		MaxAge:         corsMaxAge,
	// 		Debug:          cfg.DebugCORS,
	// 	}))
	// }
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

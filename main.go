// Package main This application is a monolith recs-seem service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/database"
)

const (
	// readTimeout リクエスト読み取りのタイムアウト。
	readTimeout = 5 * time.Second
	// writeTimeout レスポンス書き込みのタイムアウト。
	writeTimeout = 10 * time.Second
)

func main() {
	cfg, err := config.Get()
	ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Open(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.Default(),
		NoColor: runtime.GOOS == "windows",
	}))

	// clk := clock.New()
	// if cfg.FakeTime.Enabled {
	// 	log.Println("Fake time mode is enabled")

	// 	fakeClk := fakeclock.New(cfg.FakeTime.Time)
	// 	clk = fakeClk

	// 	r.Mount("/faketime", faketime.NewAPI(fakeClk, "/faketime"))
	// }

	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	defer srv.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	netAddr := listener.Addr()

	v, ok := netAddr.(*net.TCPAddr)

	if !ok {
		log.Fatal("listener.Addr() is not *net.TCPAddr")
	}

	port := v.Port

	log.Printf("Server listening on port %d", port)
	if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

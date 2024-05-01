package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	// readTimeout リクエスト読み取りのタイムアウト。
	readTimeout = 5 * time.Second
	// writeTimeout レスポンス書き込みのタイムアウト。
	writeTimeout = 10 * time.Second
)

type Server struct {
	srv     *http.Server
	l       net.Listener
	cleanup func() error
}

func NewServer(l net.Listener, mux http.Handler, cleanup func() error) *Server {
	return &Server{
		srv: &http.Server{
			Handler:      mux,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		l:       l,
		cleanup: cleanup,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ePort, ok := s.l.Addr().(*net.TCPAddr)
	if !ok {
		return fmt.Errorf("failed to get port")
	}
	log.Printf("Server listening on port %d", ePort.Port)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed は
		// http.Server.Shutdown() が正常に終了したことを示すので異常ではない。
		if err := s.srv.Serve(s.l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// グレースフルシャットダウンの終了を待つ。
	return eg.Wait()
}

func (s *Server) Close() error {
	if err := s.l.Close(); err != nil {
		return err
	}
	if err := s.srv.Close(); err != nil {
		return err
	}
	if err := s.cleanup(); err != nil {
		return err
	}
	return nil
}

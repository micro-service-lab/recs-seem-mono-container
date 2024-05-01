package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	eg.Go(func() error {
		s := NewServer(l, mux, func() error {
			return nil
		})
		return s.Run(ctx)
	})
	in := "message"
	// どんなポート番号でリッスンしているのか確認
	rawURL := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	rsp, err := http.Get(u.String())
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// サーバの終了動作を検証する
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	// 戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}

// Package faketime は時刻偽装モードの API を提供する。
package faketime

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
)

const (
	getPath = "/get"
	setPath = "/set"
)

var _ http.Handler = (*API)(nil)

// API 時刻偽装モード API を表す。
type API struct {
	clk      *fakeclock.Clock
	basePath string
}

// NewAPI API を生成して返す。
func NewAPI(clk *fakeclock.Clock, basePath string) *API {
	return &API{
		clk:      clk,
		basePath: basePath,
	}
}

// ServeHTTP リクエストに応じてハンドラを実行する。
func (s *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodGet {
		writeHTTPError(w, http.StatusForbidden)
		return
	}

	reqPath := r.URL.Path
	if s.basePath != "" && s.basePath != "/" {
		reqPath = strings.TrimPrefix(reqPath, s.basePath)
	}

	switch reqPath {
	case getPath:
		s.getHandler(w, r)
	case setPath:
		s.setHandler(w, r)
	default:
		writeHTTPError(w, http.StatusForbidden)
	}
}

// getHandler 時刻取得 API の HTTP ハンドラ。
func (s *API) getHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, s.clk.Now().Format(time.RFC3339))
}

// getHandler 時刻設定 API の HTTP ハンドラ。
func (s *API) setHandler(w http.ResponseWriter, r *http.Request) {
	t, err := time.Parse(time.RFC3339, r.URL.Query().Get("t"))
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest)
		return
	}

	s.clk.SetTime(t)

	fmt.Fprint(w, "ok")
}

// writeHTTPError w にエラーレスポンスを書き込む。
// エラーメッセージは HTTP ステータスコードに対応した文字列になる。
func writeHTTPError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	fmt.Fprint(w, http.StatusText(code))
}

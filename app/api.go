package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/app/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// API API を表す構造体。
type API struct {
	// clk 現在時刻を取得するためのインターフェース
	clk clock.Clock
	// auth 認証サービス
	auth auth.Auth
	// db データベースハンドラ
	db store.Store

	// middlewares API で使用する HTTP ミドルウェア
	middlewares []func(http.Handler) http.Handler
}

// NewAPI API を生成して返す。
func NewAPI(clk clock.Clock, auth auth.Auth, db store.Store) *API {
	return &API{
		clk:         clk,
		auth:        auth,
		db:          db,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Use HTTP ミドルウェアを追加する。
func (s *API) Use(middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Handler API 実装を提供する http.Handler を返す。
func (s *API) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(s.middlewares...)
	r.Use(
		allowContentTypeMiddleware("application/json"),
		middleware.RequestID,
		middleware.SetHeader("Content-Type", "application/json"),
	)

	r.Post("/ping", s.pingHandler)

	r.NotFound(s.notFound)
	r.MethodNotAllowed(s.methodNotAllowed)

	return r
}

// PingRequest ping リクエストを表す。
type PingRequest struct {
	// Message 任意の文字列
	Message string `json:"message"`
}

// PingResponse ping レスポンスを表す。
type PingResponse struct {
	// Message 受信した文字列
	Message string `json:"message"`
	// ReceivedTime サーバー受信時刻
	ReceivedTime time.Time `json:"receivedTime"`
}

// pingHandler 疎通確認 API の HTTP ハンドラ。
func (s *API) pingHandler(w http.ResponseWriter, r *http.Request) {
	// 受信時刻を取得する。
	// 現在時刻を参照するときは time.Now() ではなく、s.clk.Now() を使用するようにする。
	// これにより単体テスト等で時刻を偽装することが容易になる。
	receivedTime := s.clk.Now()

	// リクエストボディをデコードする。
	req := &PingRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// デコードに失敗した場合はログ出力して 400 Bad Request を返す。
		log.Printf("[ERROR] request decoding failed: %+v", err)
		errAtr := response.ApplicationErrorAttributes{
			"error": "invalid json",
		}
		err := response.Writer(r.Context(), w, response.Validation, nil, errAtr)
		if err != nil {
			log.Printf("[ERROR] response writing failed: %+v", err)
		}
		return
	}

	// レスポンスボディを表す構造体を生成する。
	resp := &PingResponse{
		Message:      req.Message,
		ReceivedTime: receivedTime,
	}

	err := response.Writer(r.Context(), w, response.Success, resp, nil)
	if err != nil {
		log.Printf("[ERROR] response writing failed: %+v", err)
	}
}

func (s *API) notFound(w http.ResponseWriter, r *http.Request) {
	err := response.Writer(r.Context(), w, response.NotFound, nil, nil)
	if err != nil {
		log.Printf("[ERROR] response writing failed: %+v", err)
	}
}

func (s *API) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	err := response.Writer(r.Context(), w, response.MethodNotAllowed, nil, nil)
	if err != nil {
		log.Printf("[ERROR] response writing failed: %+v", err)
	}
}

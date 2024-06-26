// Package api provides a server application.
package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// API API を表す構造体。
type API struct {
	// clk 現在時刻を取得するためのインターフェース
	clk clock.Clock
	// auth 認証サービス
	auth auth.Auth
	// svc サービスハンドラ
	svc service.ManagerInterface

	// middlewares API で使用する HTTP ミドルウェア
	middlewares []func(http.Handler) http.Handler
}

// NewAPI API を生成して返す。
func NewAPI(clk clock.Clock, auth auth.Auth, svc service.ManagerInterface) *API {
	return &API{
		clk:         clk,
		auth:        auth,
		svc:         svc,
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

	r.Post("/ping", handler.PingHandler(s.clk))

	r.Mount("/attend_statuses", AttendStatusHandler(s.svc))
	r.Mount("/attendance_types", AttendanceTypeHandler(s.svc))
	r.Mount("/event_types", EventTypeHandler(s.svc))
	r.Mount("/permission_categories", PermissionCategoryHandler(s.svc))
	r.Mount("/policy_categories", PolicyCategoryHandler(s.svc))
	r.Mount("/record_types", RecordTypeHandler(s.svc))
	r.Mount("/mime_types", MimeTypeHandler(s.svc))

	r.NotFound(s.notFound)
	r.MethodNotAllowed(s.methodNotAllowed)

	return r
}

func (s *API) notFound(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponseWriter(r.Context(), w, response.NotFound, nil, nil)
	if err != nil {
		log.Printf("[ERROR] response writing failed: %+v", err)
	}
}

func (s *API) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	err := response.JSONResponseWriter(r.Context(), w, response.MethodNotAllowed, nil, nil)
	if err != nil {
		log.Printf("[ERROR] response writing failed: %+v", err)
	}
}

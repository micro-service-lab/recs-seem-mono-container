// Package api provides a server application.
package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// API API を表す構造体。
type API struct {
	// clk 現在時刻を取得するためのインターフェース
	clk clock.Clock
	// auth 認証サービス
	auth auth.Auth
	// validator バリデーションサービス
	validator validation.Validator
	// svc サービスハンドラ
	svc service.ManagerInterface

	// middlewares API で使用する HTTP ミドルウェア
	middlewares []func(http.Handler) http.Handler

	// translator 翻訳サービス
	translator i18n.Translation

	// ssm セッションマネージャ
	ssm session.Manager
}

// NewAPI API を生成して返す。
func NewAPI(
	clk clock.Clock,
	auth auth.Auth,
	validator validation.Validator,
	svc service.ManagerInterface,
	translator i18n.Translation,
	ssm session.Manager,
) *API {
	return &API{
		clk:         clk,
		auth:        auth,
		validator:   validator,
		svc:         svc,
		middlewares: make([]func(http.Handler) http.Handler, 0),
		translator:  translator,
		ssm:         ssm,
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
		// json or multipart/form-data のみ受け付ける
		allowContentTypeMiddleware("application/json", "multipart/form-data"),
		middleware.RequestID,
		middleware.SetHeader("Content-Type", "application/json"),
	)

	r.Post("/ping", handler.PingHandler(s.clk))

	r.Mount("/auth", AuthHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/attend_statuses", AttendStatusHandler(s.svc))
	r.Mount("/attendance_types", AttendanceTypeHandler(s.svc))
	r.Mount("/event_types", EventTypeHandler(s.svc))
	r.Mount("/permission_categories", PermissionCategoryHandler(s.svc))
	r.Mount("/policy_categories", PolicyCategoryHandler(s.svc))
	r.Mount("/record_types", RecordTypeHandler(s.svc))
	r.Mount("/mime_types", MimeTypeHandler(s.svc))
	r.Mount("/permissions", PermissionHandler(s.svc))
	r.Mount("/policies", PolicyHandler(s.svc, s.validator, s.translator))
	r.Mount("/roles", RoleHandler(s.svc, s.validator, s.translator))
	r.Mount("/organizations", OrganizationHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/students", StudentHandler(s.svc, s.validator, s.translator))
	r.Mount("/professors", ProfessorHandler(s.svc, s.validator, s.translator))
	r.Mount("/members", MemberHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/chat_room_action_types", ChatRoomActionTypeHandler(s.svc))
	r.Mount("/chat_rooms", ChatRoomHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/images", ImageHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/files", FileHandler(s.svc, s.validator, s.translator, s.clk, s.auth, s.ssm))
	r.Mount("/read_receipts", ReadReceiptHandler(s.svc, s.clk, s.auth, s.ssm))

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

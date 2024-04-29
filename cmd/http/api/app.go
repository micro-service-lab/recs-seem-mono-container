package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	// DefaultAPIBasePath API のベースパスのデフォルト設定。
	DefaultAPIBasePath = "/"
	// StorageBasePath ストレージのベースパス。
	StorageBasePath = "/_storage"
)

// App サーバーアプリケーションを表す構造体。
type App struct {
	api *API

	// apiBasePath API のベースパス
	apiBasePath string
	// storageBasePath ストレージのベースパス
	storageBasePath string
	// storageLocalPath ストレージのローカルファイルシステム上のパス
	storageLocalPath string

	// middlewares サーバーアプリケーション全体で使用する HTTP ミドルウェア
	middlewares []func(http.Handler) http.Handler
}

// NewApp App を生成して返す。
func NewApp(api *API, options *AppOptions) *App {
	app := &App{
		api:         api,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}

	if options != nil {
		app.apiBasePath = options.APIBasePath
		app.storageLocalPath = options.StorageLocalPath
	}

	app.normalize()

	return app
}

// normalize 未設定のフィールドにデフォルト値を設定する。
func (s *App) normalize() {
	if s.apiBasePath == "" {
		s.apiBasePath = DefaultAPIBasePath
	}
	if s.storageBasePath == "" {
		s.storageBasePath = StorageBasePath
	}
}

// Use HTTP ミドルウェアを追加する。
func (s *App) Use(middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Handler サーバーアプリケーションを提供する http.Handler を返す。
func (s *App) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(s.middlewares...)
	r.Use(middleware.NoCache)

	r.Mount(s.apiBasePath, s.api.Handler())

	r.Get("/", s.root())

	if s.storageLocalPath != "" {
		r.Mount(s.storageBasePath, http.StripPrefix(s.storageBasePath, http.FileServer(http.Dir(s.storageLocalPath))))
	}

	return r
}

// root ルートパスのハンドラ関数を返す。
func (s *App) root() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}
}

// AppOptions サーバーアプリケーションのオプション設定を表す。
type AppOptions struct {
	// APIBasePath API のベースパス
	APIBasePath string
	// StorageBasePath ストレージのベースパス
	StorageBasePath string
	// StorageLocalPath ストレージのローカルファイルシステム上のパス
	StorageLocalPath string
}

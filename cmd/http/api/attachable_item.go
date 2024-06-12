package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// AttachableItemHandler is a handler for attachable item.
func AttachableItemHandler(
	svc service.ManagerInterface,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	downloadHandler := handler.DownloadObject{
		Service: svc,
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Get("/{attachable_item_id}", downloadHandler.ServeHTTP)
	})

	return r
}

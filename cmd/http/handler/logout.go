package handler

import (
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// Logout is a handler for logout.
type Logout struct {
	Service service.ManagerInterface
}

func (h *Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	authUser := auth.FromContext(ctx)
	if authUser == nil {
		err = response.JSONResponseWriter(ctx, w, response.AlreadyLogout, nil, nil)
	} else if err = h.Service.Logout(
		ctx,
		authUser.MemberID,
	); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   auth.AccessTokenCookieKey,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

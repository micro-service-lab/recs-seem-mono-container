package handler

import (
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// GetUnreadReceiptsCount is a handler for getting unread receipts count.
type GetUnreadReceiptsCount struct {
	Service service.ManagerInterface
}

// ServeHTTP reads message.
func (h *GetUnreadReceiptsCount) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	read, err := h.Service.CountUnreadReceiptsOnMember(
		ctx,
		authUser.MemberID,
	)
	if err == nil {
		res := struct {
			Read int64 `json:"read"`
		}{
			Read: read,
		}
		err = response.JSONResponseWriter(ctx, w, response.Success, res, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}
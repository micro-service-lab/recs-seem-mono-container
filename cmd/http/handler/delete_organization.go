package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// DeleteOrganization is a handler for creating organization.
type DeleteOrganization struct {
	Service service.ManagerInterface
}

func (h *DeleteOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	id := uuid.MustParse(chi.URLParam(r, "organization_id"))
	c, err := h.Service.DeleteOrganization(ctx, id, authUser.MemberID)
	if err != nil || c == 0 {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomTargetChatRoom:
				e.SetTarget("chat room")
				err = e
			case service.OrganizationTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.OrganizationTargetOrganizations:
				e.SetTarget("organization")
				err = e
			}
		}
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

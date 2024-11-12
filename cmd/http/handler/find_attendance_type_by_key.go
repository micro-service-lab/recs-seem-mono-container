package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// FindAttendanceTypeByKey is a handler for finding attendance type.
type FindAttendanceTypeByKey struct {
	Service service.ManagerInterface
}

func (h *FindAttendanceTypeByKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "attendance_type_key")
	attendanceType, err := h.Service.FindAttendanceTypeByKey(ctx, key)
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	err = response.JSONResponseWriter(ctx, w, response.Success, attendanceType, nil)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

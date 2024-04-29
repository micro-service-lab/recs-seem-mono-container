package errhandle

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

func ErrorHandle(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
) (bool, error) {
	if err == nil {
		return false, nil
	}
	var e ApplicationError

	if errors.As(err, &e) {
		rType, errAttr := e.ResolveCodeAndAttribute()
		if err := response.JsonResponseWriter(ctx, w, rType, nil, errAttr); err != nil {
			return false, fmt.Errorf("failed to write response: %w", err)
		}
		return true, nil
	}
	response.JsonResponseWriter(ctx, w, response.System, nil, nil)
	return true, nil
}

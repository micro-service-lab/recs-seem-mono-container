// Package errhandle provides error handling for the application.
package errhandle

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ErrorHandle handles errors.
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
		if err := response.JSONResponseWriter(ctx, w, rType, nil, errAttr); err != nil {
			return false, fmt.Errorf("failed to write response: %w", err)
		}
		return true, nil
	}
	er := response.JSONResponseWriter(ctx, w, response.System, nil, nil)
	if er != nil {
		return false, fmt.Errorf("failed to write response: %w", er)
	}
	return true, nil
}

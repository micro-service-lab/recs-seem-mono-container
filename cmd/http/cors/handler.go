package cors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// AppHandler CORS エラーを処理するハンドラ。
func AppHandler(w http.ResponseWriter, _ *http.Request, c Cors, err error) bool {
	_, ok := err.(Error)
	if ok {
		c.Log.Printf("CORS error: %v", err)
		res := struct {
			Message string `json:"message"`
		}{
			Message: "CORS error: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		noOrigin := false
		switch {
		case errors.Is(err, &PreflightEmptyOriginError{}), errors.Is(err, &ActualMissingOriginError{}):
			noOrigin = true
		case errors.Is(err, &PreflightNotOptionMethodError{}),
			errors.Is(err, &PreflightNotAllowedMethodError{}),
			errors.Is(err, &ActualMethodNotAllowedError{}):
			w.WriteHeader(http.StatusMethodNotAllowed)
		default:
			w.WriteHeader(http.StatusForbidden)
		}
		// For requests that do not conform to the browser's same-origin policy (no Origin header,
		// such as postman, is given), pass through processing.
		if noOrigin {
			return true
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			c.Log.Printf("CORS error encoding failed: %v", err)
		}
		return false
	}
	res := struct {
		Message string `json:"message"`
	}{
		Message: "CORS error: An unexpected error has occurred",
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		c.Log.Printf("CORS error encoding failed: %v", err)
	}
	return false
}

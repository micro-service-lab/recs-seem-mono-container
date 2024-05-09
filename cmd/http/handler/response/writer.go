package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONResponseWriter writes a JSON response.
func JSONResponseWriter(
	_ context.Context, w http.ResponseWriter, rType APIResponseType, data any, errAttr ApplicationErrorAttributes,
) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := ApplicationResponse{
		Success:         rType == Success,
		Data:            data,
		Code:            rType.Code,
		Message:         rType.Message,
		ErrorAttributes: errAttr,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		rType = System
		w.WriteHeader(rType.StatusCode)
		rsp := ApplicationResponse{
			Success:         false,
			Data:            nil,
			Code:            rType.Code,
			Message:         rType.Message,
			ErrorAttributes: nil,
		}
		preErr := fmt.Errorf("marshal error response error: %w", err)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			return fmt.Errorf("write error response error: %w", err)
		}
		return preErr
	}
	w.WriteHeader(rType.StatusCode)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		return fmt.Errorf("write response error: %w", err)
	}

	return nil
}

// JSONResponseWriterWithHeader writes a JSON response with a custom header.
func JSONResponseWriterWithHeader(
	_ context.Context,
	w http.ResponseWriter,
	rType APIResponseType,
	data any,
	errAttr ApplicationErrorAttributes,
	header http.Header,
) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range header {
		w.Header()[k] = v
	}
	body := ApplicationResponse{
		Success:         rType == Success,
		Data:            data,
		Code:            rType.Code,
		Message:         rType.Message,
		ErrorAttributes: errAttr,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		rType = System
		w.WriteHeader(rType.StatusCode)
		rsp := ApplicationResponse{
			Success:         false,
			Data:            nil,
			Code:            rType.Code,
			Message:         rType.Message,
			ErrorAttributes: nil,
		}
		preErr := fmt.Errorf("marshal error response error: %w", err)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			return fmt.Errorf("write error response error: %w", err)
		}
		return preErr
	}
	w.WriteHeader(rType.StatusCode)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		return fmt.Errorf("write response error: %w", err)
	}

	return nil
}

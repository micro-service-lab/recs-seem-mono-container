package errhandle

import (
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ModelNotFoundError is an error for model not found.
type ModelNotFoundError struct {
	target string
}

// NewModelNotFoundError creates a new model not found error.
func NewModelNotFoundError(target string) ModelNotFoundError {
	return ModelNotFoundError{target: target}
}

// Error returns the error message.
func (e ModelNotFoundError) Error() string {
	return fmt.Sprintf("model not found: %s", e.target)
}

// Is checks if the target is a model not found error.
func (e ModelNotFoundError) Is(target error) bool {
	_, ok := target.(*ModelNotFoundError)
	return ok
}

// As checks if the target is a model not found error.
func (e ModelNotFoundError) As(target any) bool {
	_, ok := target.(*ModelNotFoundError)
	return ok
}

// Target returns the target.
func (e ModelNotFoundError) Target() string {
	return e.target
}

// SetTarget sets the target.
func (e *ModelNotFoundError) SetTarget(target string) ModelNotFoundError {
	e.target = target
	return *e
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e ModelNotFoundError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.NotFoundModel, map[string]any{
		"target": e.target,
	}
}

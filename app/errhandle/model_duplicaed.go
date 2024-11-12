package errhandle

import (
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ModelDuplicatedError is an error for model duplicated.
type ModelDuplicatedError struct {
	target string
}

// NewModelDuplicatedError creates a new model duplicated error.
func NewModelDuplicatedError(target string) ModelDuplicatedError {
	return ModelDuplicatedError{target: target}
}

// Error returns the error message.
func (e ModelDuplicatedError) Error() string {
	return fmt.Sprintf("model duplicated: %s", e.target)
}

// Is checks if the target is a model duplicated error.
func (e ModelDuplicatedError) Is(target error) bool {
	_, ok := target.(*ModelDuplicatedError)
	return ok
}

// As checks if the target is a model duplicated error.
func (e ModelDuplicatedError) As(target any) bool {
	_, ok := target.(*ModelDuplicatedError)
	return ok
}

// Target returns the target.
func (e ModelDuplicatedError) Target() string {
	return e.target
}

// SetTarget sets the target.
func (e *ModelDuplicatedError) SetTarget(target string) ModelDuplicatedError {
	e.target = target
	return *e
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e ModelDuplicatedError) ResolveCodeAndAttribute() (
	response.APIResponseType, response.ApplicationErrorAttributes,
) {
	return response.ModelConflict, map[string]any{
		"target": e.target,
	}
}

package errhandle

import "github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"

// ValidationError is an error for validation.
type ValidationError struct {
	Err error
}

// NewValidationError creates a new validation error.
func NewValidationError(err error) ValidationError {
	return ValidationError{err}
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the wrapped error.
func (e ValidationError) Unwrap() error {
	return e.Err
}

// Is checks if the target is a validation error.
func (e ValidationError) Is(target error) bool {
	_, ok := target.(ValidationError)
	return ok
}

// As checks if the target is a validation error.
func (e ValidationError) As(target any) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e ValidationError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.Validation, response.ApplicationErrorAttributes{
		"error": e.Error(),
	}
}

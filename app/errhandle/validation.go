package errhandle

import (
	"sync"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ValidationError is an error for validation.
type ValidationError struct {
	mu          *sync.RWMutex
	fieldsError map[string][]string
}

// NewValidationError creates a new validation error.
func NewValidationError(e map[string][]string) ValidationError {
	if e == nil {
		e = make(map[string][]string)
	}
	return ValidationError{
		mu:          &sync.RWMutex{},
		fieldsError: e,
	}
}

// Add adds a field and error message.
func (e *ValidationError) Add(field, message string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.fieldsError[field]; !ok {
		e.fieldsError[field] = make([]string, 0)
	}
	e.fieldsError[field] = append(e.fieldsError[field], message)
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return "validation error"
}

// Is checks if the target is a validation error.
func (e ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// As checks if the target is a validation error.
func (e ValidationError) As(target any) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e ValidationError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.Validation, e.fieldsError
}

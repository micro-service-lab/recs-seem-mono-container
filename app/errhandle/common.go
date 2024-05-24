package errhandle

import (
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// CommonError is an error for model not found.
type CommonError struct {
	code response.APIResponseType
	attr any
}

// NewCommonError creates a new common error.
func NewCommonError(code response.APIResponseType, attr response.ApplicationErrorAttributes) CommonError {
	return CommonError{
		code: code,
		attr: attr,
	}
}

// Error returns the error message.
func (e CommonError) Error() string {
	return fmt.Sprintf("common error: %v", e.code)
}

// Is checks if the target is a common error.
func (e CommonError) Is(target error) bool {
	_, ok := target.(*CommonError)
	return ok
}

// As checks if the target is a common error.
func (e CommonError) As(target any) bool {
	_, ok := target.(*CommonError)
	return ok
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e CommonError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return e.code, e.attr
}

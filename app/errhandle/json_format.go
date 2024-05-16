package errhandle

import (
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// JSONFormatError is an error for json format.
type JSONFormatError struct{}

// NewJSONFormatError creates a new json format error.
func NewJSONFormatError() JSONFormatError {
	return JSONFormatError{}
}

// Error returns the error message.
func (e JSONFormatError) Error() string {
	return "json format error"
}

// Is checks if the target is a json format error.
func (e JSONFormatError) Is(target error) bool {
	_, ok := target.(*JSONFormatError)
	return ok
}

// As checks if the target is a json format error.
func (e JSONFormatError) As(target any) bool {
	_, ok := target.(*JSONFormatError)
	return ok
}

// ResolveCodeAndAttribute resolves code and attribute.
func (e JSONFormatError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.RequestFormatError, nil
}

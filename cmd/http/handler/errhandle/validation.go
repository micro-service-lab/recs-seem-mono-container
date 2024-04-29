package errhandle

import "github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"

type ValidationError struct {
	Err error
}

func NewValidationError(err error) ValidationError {
	return ValidationError{err}
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

func (e ValidationError) Is(target error) bool {
	_, ok := target.(ValidationError)
	return ok
}

func (e ValidationError) As(target interface{}) bool {
	_, ok := target.(*ValidationError)
	return ok
}

func (e ValidationError) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.Validation, response.ApplicationErrorAttributes{
		"error": e.Error(),
	}
}

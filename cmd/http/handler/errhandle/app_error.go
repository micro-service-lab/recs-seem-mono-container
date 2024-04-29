package errhandle

import "github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"

// Application is an interface for application errors.
type ApplicationError interface {
	error
	ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes)
}

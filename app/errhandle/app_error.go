package errhandle

import "github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"

// ApplicationError is an interface for application errors.
type ApplicationError interface {
	error
	// ResolveCodeAndAttribute resolves code and attribute.
	ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes)
}

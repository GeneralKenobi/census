package client

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"net/http"
)

// Configure creates an Api instance that points to host:port and uses the given protocol (http or https). The passed in httpClient is
// used to make all HTTP requests.
func Configure(protocol, host string, port int, httpClient *http.Client) Api {
	return api{
		httpClient:      httpClient,
		requestBasePath: fmt.Sprintf("%s://%s:%d", protocol, host, port),
	}
}

// ApiError is returned if an API operation fails - returns an HTTP error like bad request, produces malformed response body, etc.
type ApiError interface {
	error
	StatusCode() int
	Message() string
}

// Api aggregates all API operations implemented by census.
//
// Correlation ID is included in the requests. It is extracted using mdctx.CorrelationId from the ctx passed to methods.
type Api interface {
	PersonApi
}

type PersonApi interface {
	CreatePerson(ctx context.Context, person apimodel.PersonCreate) (apimodel.PersonCreated, error)
	GetPerson(ctx context.Context, id string) (apimodel.Person, error)
	UpdatePerson(ctx context.Context, id string, person apimodel.PersonUpdate) error
	DeletePerson(ctx context.Context, id string) error
}

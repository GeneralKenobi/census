package e2eutil

import (
	"github.com/GeneralKenobi/census/pkg/api/client"
	"net/http"
)

// IsBadRequest returns true if err is not nil and if err is a client.ApiError with status bad request.
func IsBadRequest(err error) bool {
	return isApiErrorWithStatusCode(err, http.StatusBadRequest)
}

// IsNotFound returns true if err is not nil and if err is a client.ApiError with status not found.
func IsNotFound(err error) bool {
	return isApiErrorWithStatusCode(err, http.StatusNotFound)
}

func isApiErrorWithStatusCode(err error, expectedStatusCode int) bool {
	if err == nil {
		return false
	}

	apiErr, ok := err.(client.ApiError)
	if !ok {
		return false
	}

	return apiErr.StatusCode() == expectedStatusCode
}

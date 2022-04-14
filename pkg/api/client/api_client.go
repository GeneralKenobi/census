package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/util"
	"io"
	"io/ioutil"
	"net/http"
)

type api struct {
	httpClient      *http.Client
	requestBasePath string
}

var _ Api = (*api)(nil) // Interface guard

func (api api) get(ctx context.Context, path string) (*http.Response, error) {
	return api.doRequest(ctx, http.MethodGet, path)
}

func (api api) delete(ctx context.Context, path string) (*http.Response, error) {
	return api.doRequest(ctx, http.MethodDelete, path)
}

func (api api) postPayload(ctx context.Context, path string, body any) (*http.Response, error) {
	return api.doRequestWithPayload(ctx, http.MethodPost, path, body)
}

func (api api) putPayload(ctx context.Context, path string, body any) (*http.Response, error) {
	return api.doRequestWithPayload(ctx, http.MethodPut, path, body)
}

func (api api) doRequest(ctx context.Context, httpMethod, path string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, httpMethod, api.requestBasePath+path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set(correlationIdHeader, mdctx.CorrelationId(ctx))

	response, err := api.httpClient.Do(request)
	if err != nil {
		return response, fmt.Errorf("error executing request: %w", err)
	}
	return response, nil
}

func (api api) doRequestWithPayload(ctx context.Context, httpMethod, path string, payload any) (*http.Response, error) {
	payloadJsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}
	payloadJsonBuffer := bytes.NewBuffer(payloadJsonBytes)

	request, err := http.NewRequestWithContext(ctx, httpMethod, api.requestBasePath+path, payloadJsonBuffer)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set(correlationIdHeader, mdctx.CorrelationId(ctx))
	request.Header.Set(contentTypeHeader, "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return response, fmt.Errorf("error executing request: %w", err)
	}
	return response, nil
}

// processResponse checks if the response is successful (HTTP 2XX), if not it returns an apiError.
//
// Response body is closed within this function.
func processResponse(response *http.Response) error {
	defer response.Body.Close()
	return assertResponseSuccessful(response)
}

// processResponse checks if the response is successful (HTTP 2XX), if not it returns an apiError. It then unmarshalls response body to an
// instance of T, if that fails then apiError is also returned.
//
// Response body is closed within this function.
func processResponseWithPayload[T any](response *http.Response) (T, error) {
	defer response.Body.Close()

	if err := assertResponseSuccessful(response); err != nil {
		return util.ZeroValue[T](), err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		apiErr := apiError{
			cause:      fmt.Errorf("error reading response body: %w", err),
			statusCode: response.StatusCode,
		}
		return util.ZeroValue[T](), apiErr
	}

	var responsePayload T
	err = json.Unmarshal(responseBytes, &responsePayload)
	if err != nil {
		apiErr := apiError{
			cause:      fmt.Errorf("error unmarshaling response body: %w", err),
			statusCode: response.StatusCode,
		}
		return util.ZeroValue[T](), apiErr
	}

	return responsePayload, nil
}

func assertResponseSuccessful(response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	return apiError{
		cause:      fmt.Errorf("expected success status code 2XX"),
		statusCode: response.StatusCode,
		message:    responseBodyAsString(response.Body),
	}
}

func responseBodyAsString(responseBody io.Reader) string {
	responseBytes, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return ""
	}
	return string(responseBytes)
}

type apiError struct {
	cause      error
	statusCode int
	message    string
}

var _ ApiError = (*apiError)(nil) // Interface guard

func (apiErr apiError) Error() string {
	return fmt.Sprintf("%v: status code: %d, message: %s", apiErr.cause, apiErr.statusCode, apiErr.message)
}

func (apiErr apiError) StatusCode() int {
	return apiErr.statusCode
}

func (apiErr apiError) Message() string {
	return apiErr.message
}

const (
	correlationIdHeader = "X-Correlation-ID"
	contentTypeHeader   = "Content-Type"
)

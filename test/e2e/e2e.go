package e2e

import (
	"context"
	"crypto/tls"
	"github.com/GeneralKenobi/census/pkg/api/client"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"net/http"
	"os"
	"strconv"
	"testing"
)

// Environment variables configuring tests.
const (
	censusHostEnv     = "CENSUS_HOST"     // Required env with census host, e.g. localhost
	censusPortEnv     = "CENSUS_PORT"     // Required env with the port census is listening on, e.g. 8443
	censusProtocolEnv = "CENSUS_PROTOCOL" // Required env with the protocol to use when connecting to census, http or https
)

// Context creates a context with a correlation ID and logs that correlation ID to t.
func Context(t *testing.T) context.Context {
	ctx := mdctx.New()
	t.Logf("Correlation ID: %s", mdctx.CorrelationId(ctx))
	return ctx
}

// GetApi configures the API client library using environment variables.
func GetApi() client.Api {
	protocol := getApiProtocol()
	host := getApiHost()
	port := getApiPort()
	httpClient := httpClientWithoutTlsVerification()
	return client.Configure(protocol, host, port, httpClient)
}

func getApiProtocol() string {
	host := os.Getenv(censusProtocolEnv)
	if host == "" {
		panic(censusProtocolEnv + " env is required")
	}
	return host
}

func getApiHost() string {
	host := os.Getenv(censusHostEnv)
	if host == "" {
		panic(censusHostEnv + " env is required")
	}
	return host
}

func getApiPort() int {
	portStr := os.Getenv(censusPortEnv)
	if portStr == "" {
		panic(censusPortEnv + " env is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(censusPortEnv + " has to be an integer")
	}

	return port
}

func httpClientWithoutTlsVerification() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

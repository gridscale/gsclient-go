package gsclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type networkTestCase struct {
	name string
	apiURL string
	httpClient *http.Client
	expectedError string
}

var networkErrorTests []networkTestCase = []networkTestCase{
	{
		name: "Connection timeout",
		apiURL: defaultAPIURL,
		httpClient: &http.Client{Timeout: 1 * time.Nanosecond},
		expectedError: "Maximum number of trials has been exhausted with error: Get %s%s: net/http: request canceled (Client.Timeout exceeded while awaiting headers)",
	},
	{
		name: "Connection refused",
		apiURL: "http://127.0.0.1",
		httpClient: http.DefaultClient,
		expectedError: "Maximum number of trials has been exhausted with error: Get %s%s: dial tcp 127.0.0.1:80: connect: connection refused",
	},
	
}

func TestRequest_RetryNetworkErrors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()
	
	for _, test := range networkErrorTests {
		config := NewConfiguration(test.apiURL, "uuid", "token", true, true, 1, 100, 5)
		config.httpClient = test.httpClient
		client := NewClient(config)
		_, err := client.GetServer(emptyCtx, dummyUUID)
		assert.EqualError(t, err, fmt.Sprintf(test.expectedError, config.apiURL, uri), test.name)
	}
}

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

var getNetworkErrorTests []networkTestCase = []networkTestCase{
	{
		name: "rety the GET requet in case of connection timeout",
		apiURL: "http://127.0.0.1",
		httpClient: &http.Client{Timeout: 1 * time.Nanosecond},
		expectedError: "Maximum number of trials has been exhausted with error: Get %s%s: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)",
	},
	{
		name: "rety the GET requet in case of connection refused",
		apiURL: "http://127.0.0.1",
		httpClient: http.DefaultClient,
		expectedError: "Maximum number of trials has been exhausted with error: Get %s%s: dial tcp 127.0.0.1:80: connect: connection refused",
	},
	{
		name: "rety the GET requet in case of DNS lookup error",
		apiURL: "http://api.unkown.domain",
		httpClient: http.DefaultClient,
		expectedError: "Maximum number of trials has been exhausted with error: Get %s%s: dial tcp: lookup api.unkown.domain: no such host",
	},
}

var postNetworkErrorTests []networkTestCase = []networkTestCase{
	{
		name: "do not rety the POST requet in case of connection timeout",
		apiURL: "http://127.0.0.1",
		httpClient: &http.Client{Timeout: 1 * time.Nanosecond},
		expectedError: "Post %s%s: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)",
	},
	{
		name: "rety the POST requet in case of connection refused",
		apiURL: "http://127.0.0.1",
		httpClient: http.DefaultClient,
		expectedError: "Maximum number of trials has been exhausted with error: Post %s%s: dial tcp 127.0.0.1:80: connect: connection refused",
	},
	{
		name: "rety the POST requet in case of DNS lookup error",
		apiURL: "http://api.unkown.domain",
		httpClient: http.DefaultClient,
		expectedError: "Maximum number of trials has been exhausted with error: Post %s%s: dial tcp: lookup api.unkown.domain: no such host",
	},
}

func TestGetRequest_NetworkErrors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()
	
	for _, test := range getNetworkErrorTests {
		config := NewConfiguration(test.apiURL, "uuid", "token", true, true, 1, 100, 5)
		config.httpClient = test.httpClient
		client := NewClient(config)
		_, err := client.GetServer(emptyCtx, dummyUUID)
		assert.EqualError(t, err, fmt.Sprintf(test.expectedError, config.apiURL, uri), test.name)
	}
}

func TestPostRequest_NetworkErrors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	uri := apiServerBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()
	
	for _, test := range postNetworkErrorTests {
		config := NewConfiguration(test.apiURL, "uuid", "token", true, true, 1, 100, 5)
		config.httpClient = test.httpClient
		client := NewClient(config)
		_, err := client.CreateServer(emptyCtx, ServerCreateRequest{
				Name:            "test",
				Memory:          10,
				Cores:           4,
				LocationUUID:    dummyUUID,
				HardwareProfile: DefaultServerHardware,
				Labels:          []string{"label"},
			})
		assert.EqualError(t, err, fmt.Sprintf(test.expectedError, config.apiURL, uri), test.name)
	}
}

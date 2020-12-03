package gsclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type networkTestCase struct {
	name          string
	apiURL        string
	httpClient    *http.Client
	expectedError string
}

type apiTestCase struct {
	name          string
	statusCode    int
	dummyUUID     string
	expectedError string
}

var getNetworkErrorTests = []networkTestCase{
	{
		name:          "retry the GET request in case of connection timeout",
		apiURL:        "http://127.0.0.1",
		httpClient:    &http.Client{Timeout: 1 * time.Nanosecond},
		expectedError: "Client.Timeout exceeded while awaiting headers",
	},
	{
		name:          "retry the GET request in case of connection refused",
		apiURL:        "http://127.0.0.1",
		httpClient:    http.DefaultClient,
		expectedError: "dial tcp 127.0.0.1:80: connect: connection refused",
	},
	{
		name:          "retry the GET request in case of DNS lookup error",
		apiURL:        "http://api.unknown.domain",
		httpClient:    http.DefaultClient,
		expectedError: "dial tcp: lookup api.unknown.domain",
	},
}

var postNetworkErrorTests = []networkTestCase{
	{
		name:          "do not retry the POST request in case of connection timeout",
		apiURL:        "http://127.0.0.1",
		httpClient:    &http.Client{Timeout: 1 * time.Nanosecond},
		expectedError: "Client.Timeout exceeded while awaiting headers",
	},
	{
		name:          "retry the POST request in case of connection refused",
		apiURL:        "http://127.0.0.1",
		httpClient:    http.DefaultClient,
		expectedError: "dial tcp 127.0.0.1:80: connect: connection refused",
	},
	{
		name:          "retry the POST request in case of DNS lookup error",
		apiURL:        "http://api.unknown.domain",
		httpClient:    http.DefaultClient,
		expectedError: "dial tcp: lookup api.unknown.domain",
	},
}

var apiErrorTests = []apiTestCase{
	{
		name:          "retry the request in case of API error with status code 500",
		statusCode:    500,
		dummyUUID:     "690de890-13c0-4e76-8a01-e10ba8786e53",
		expectedError: "Maximum number of trials has been exhausted with error: Status code: %d. Error: no error message received from server. Request UUID: %s.",
	},
	{
		name:          "retry the request in case of API error with status code 424",
		statusCode:    424,
		dummyUUID:     "690de890-13c0-4e76-8a01-e10ba8786e54",
		expectedError: "Maximum number of trials has been exhausted with error: Status code: %d. Error: no error message received from server. Request UUID: %s. ",
	},
	{
		name:          "retry the request in case of API error with status code 409",
		statusCode:    409,
		dummyUUID:     "690de890-13c0-4e76-8a01-e10ba8786e55",
		expectedError: "Maximum number of trials has been exhausted with error: Status code: %d. Error: no error message received from server. Request UUID: %s.",
	},
}

func TestRequestGet_NetworkErrors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	for _, test := range getNetworkErrorTests {
		config := NewConfiguration(test.apiURL, "uuid", "token", true, true, 100, 5)
		config.httpClient = test.httpClient
		client := NewClient(config)
		_, err := client.GetServer(emptyCtx, dummyUUID)
		assert.Contains(t, fmt.Sprintf("%v", err), test.expectedError, test.name)
	}
}

func TestRequestPost_NetworkErrors(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	uri := apiServerBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	for _, test := range postNetworkErrorTests {
		config := NewConfiguration(test.apiURL, "uuid", "token", true, true, 100, 5)
		config.httpClient = test.httpClient
		client := NewClient(config)
		_, err := client.CreateServer(emptyCtx, ServerCreateRequest{
			Name:            "test",
			Memory:          10,
			Cores:           4,
			HardwareProfile: DefaultServerHardware,
			Labels:          []string{"label"},
		})
		assert.Contains(t, fmt.Sprintf("%v", err), test.expectedError, test.name)
	}
}

func TestRequestGet_APIErrors(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	for _, test := range apiErrorTests {
		uri := path.Join(apiServerBase, test.dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
			w.WriteHeader(test.statusCode)
		})
		_, err := client.GetServer(emptyCtx, test.dummyUUID)
		assert.Contains(t, fmt.Sprintf("%v", err), fmt.Sprintf(test.expectedError, test.statusCode, dummyRequestUUID), test.name)
	}
}

func TestRequestPatch_APIErrors(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	for _, test := range apiErrorTests {
		uri := path.Join(apiServerBase, test.dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
			w.WriteHeader(test.statusCode)
		})
		err := client.UpdateServer(
			emptyCtx,
			test.dummyUUID,
			ServerUpdateRequest{
				Name:   "test",
				Memory: 4,
				Cores:  2,
				Labels: nil,
			})
		assert.Contains(t, fmt.Sprintf("%v", err), fmt.Sprintf(test.expectedError, test.statusCode, dummyRequestUUID), test.name)
	}
}

func TestRequestDelete_APIErrors(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	for _, test := range apiErrorTests {
		uri := path.Join(apiServerBase, test.dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
			w.WriteHeader(test.statusCode)
		})
		err := client.DeleteServer(emptyCtx, test.dummyUUID)
		assert.Contains(t, fmt.Sprintf("%v", err), fmt.Sprintf(test.expectedError, test.statusCode, dummyRequestUUID), test.name)
	}
}

func Test_prepareHTTPRequest(t *testing.T) {
	r := &gsRequest{
		uri:                 path.Join(apiDeletedBase, "networks"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	cfg := DefaultConfiguration("test", "test")
	httpReq, err := r.prepareHTTPRequest(context.Background(), cfg)
	assert.NotNil(t, httpReq)
	assert.Equal(t, r.method, httpReq.Method)
	assert.Equal(t, r.uri, httpReq.URL.RequestURI())
	assert.Nil(t, err)
}

func Test_isErrorHTTPCodeRetryable(t *testing.T) {
	type testCase struct {
		code        int
		isRetryable bool
	}
	testCases := []testCase{
		{
			500,
			true,
		},
		{
			424,
			true,
		},
		{
			429,
			true,
		},
		{
			409,
			true,
		},
		{
			404,
			false,
		},
	}
	for _, test := range testCases {
		isRetryable := isErrorHTTPCodeRetryable(test.code)
		assert.Equal(t, test.isRetryable, isRetryable)
	}
}

func Test_getDelayTimeInMsFromTimestampStr(t *testing.T) {
	type testCase struct {
		successful   bool
		TimestampStr string
	}
	futureTimestamp := time.Now().Add(5*time.Second).UnixNano() / 1000000
	futureTimestampStr := strconv.FormatInt(futureTimestamp, 10)
	testCases := []testCase{
		{
			true,
			futureTimestampStr,
		},
		{
			false,
			"",
		},
	}

	for _, test := range testCases {
		delay, err := getDelayTimeInMsFromTimestampStr(test.TimestampStr)
		if test.successful {
			assert.Nil(t, err)
			assert.GreaterOrEqual(t, delay, int64(0))
		} else {
			assert.NotNil(t, err)
		}
	}
}

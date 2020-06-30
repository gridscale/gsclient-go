package gsclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestClient_waitForRequestCompleted(t *testing.T) {
	requestTestCases := []string{"done", "failed", "pending"}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := NewConfiguration(
		server.URL,
		"uuid",
		"token",
		true,
		true,
		100,
		5,
	)
	client := NewClient(config)
	defer server.Close()
	var isFailed bool
	var reqStatus string
	mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprint(w, fmt.Sprintf(`{"%s": {"status":"%s", "isFailed" : %v}}`, dummyUUID, reqStatus, isFailed))
		}
	})
	for _, reqStatusTest := range requestTestCases {
		reqStatus = reqStatusTest
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testUUID := range uuidCommonTestCases {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				err := client.waitForRequestCompleted(ctx, testUUID.testUUID)
				if isFailed || testUUID.isFailed || reqStatus != requestDoneStatus {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForRequestCompleted returned an error %v", err)
				}
				cancel()
			}
		}
	}
}
func TestClient_WithHeaders(t *testing.T) {
	config := DefaultConfiguration(dummyUUID, "token")
	client := NewClient(config)
	client.WithHTTPHeaders(map[string]string{
		"test_header": "test_header_value",
	})
	assert.Equal(t, client.cfg.httpHeaders["test_header"], "test_header_value")
}

package gsclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_waitForRequestCompleted(t *testing.T) {
	requestTestCases := []string{"done", "failed", "pending"}
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var reqStatus string
	mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprint(w, fmt.Sprintf(`{"%s": {"status":"%s"}}`, dummyUUID, reqStatus))
		}
	})
	for _, reqStatusTest := range requestTestCases {
		reqStatus = reqStatusTest
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testUUID := range uuidCommonTestCases {
				err := client.waitForRequestCompleted(emptyCtx, testUUID.testUUID)
				if isFailed || reqStatus != requestDoneStatus || testUUID.isFailed || reqStatus == requestFailStatus {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForRequestCompleted returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_waitFor404Status(t *testing.T) {
	notfoundTestCases := []bool{true, false}
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isNotFound bool
	var uri = "/test/"
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		if isNotFound {
			w.WriteHeader(404)
		} else {
			fmt.Fprint(w, nil)
		}
	})
	for _, notFound := range notfoundTestCases {
		isNotFound = notFound
		err := client.waitFor404Status(emptyCtx, uri, http.MethodGet)
		if isNotFound {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestClient_waitFor200Status(t *testing.T) {
	foundTestCases := []bool{true, false}
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFound bool
	var uri = "/test/"
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		if isFound {
			fmt.Fprint(w, nil)
		} else {
			w.WriteHeader(404)
		}
	})
	for _, found := range foundTestCases {
		isFound = found
		err := client.waitFor200Status(emptyCtx, uri, http.MethodGet)
		if found {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

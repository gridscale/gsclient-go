package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetLocationList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiLocationBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareLocationListHTTPGet())
	})
	res, err := client.GetLocationList(emptyCtx)
	assert.Nil(t, err, "GetLocationList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLocation()), fmt.Sprintf("%v", res))
}

func TestClient_GetLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareLocationHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLocation returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockLocation()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiLocationBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareLocationCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateLocation(
			emptyCtx,
			LocationCreateRequest{
				Name:               "test",
				ParentLocationUUID: dummyUUID,
				CPUNodeCount:       10,
				ProductNo:          99,
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateLocation returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockLocationCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_UpdateLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiLocationBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodPatch {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				fmt.Fprint(writer, prepareLocationHTTPGet())
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdateLocation(
				emptyCtx,
				test.testUUID,
				LocationUpdateRequest{
					Name: "test",
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateLocation returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiLocationBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodDelete {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				writer.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.DeleteLocation(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteLocation returned an error %v", err)
			}
		}
	}
}

func getMockLocation() Location {
	mock := Location{Properties: LocationProperties{
		Iata:       "fra",
		Status:     "active",
		Labels:     nil,
		Name:       "de/fra",
		ObjectUUID: dummyUUID,
		Country:    "de",
	}}
	return mock
}

func prepareLocationHTTPGet() string {
	location := getMockLocation()
	res, _ := json.Marshal(location)
	return string(res)
}

func prepareLocationListHTTPGet() string {
	location := getMockLocation()
	res, _ := json.Marshal(location.Properties)
	return fmt.Sprintf(`{"locations": {"%s": %s}}`, dummyUUID, string(res))
}

func getMockLocationCreateResponse() CreateResponse {
	mock := CreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareLocationCreateResponse() string {
	res, _ := json.Marshal(getMockLocationCreateResponse())
	return string(res)
}

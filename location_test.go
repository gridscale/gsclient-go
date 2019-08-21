package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetLocationList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiLocationBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLocationListHTTPGet())
	})
	res, err := client.GetLocationList()
	if err != nil {
		t.Errorf("GetLocationList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLocation()), fmt.Sprintf("%v", res))
}

func TestClient_GetLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLocationHTTPGet())
	})
	res, err := client.GetLocation(dummyUuid)
	if err != nil {
		t.Errorf("GetLocation returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockLocation()), fmt.Sprintf("%v", res))
}

func getMockLocation() Location {
	mock := Location{Properties: LocationProperties{
		Iata:       "fra",
		Status:     "active",
		Labels:     nil,
		Name:       "de/fra",
		ObjectUuid: dummyUuid,
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
	return fmt.Sprintf(`{"locations": {"%s": %s}}`, dummyUuid, string(res))
}

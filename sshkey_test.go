package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetSshkeyList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiSshkeyBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareSshkeyListHTTPGet())
	})
	res, err := client.GetSshkeyList()
	assert.Nil(t, err, "GetSshkeyList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkey()), fmt.Sprintf("%v", res))
}

func TestClient_GetSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareSshkeyHTTPGet())
	})
	res, err := client.GetSshkey(dummyUUID)
	assert.Nil(t, err, "GetSshkey returned an error %v", err)
	assert.Equal(t, fmt.Sprintf("%v", getMockSshkey()), fmt.Sprintf("%v", res))
}

func TestClient_CreateSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiSshkeyBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareSshkeyCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateSshkey(SshkeyCreateRequest{
		Name:   "test",
		Sshkey: "example",
		Labels: []string{"label"},
	})
	assert.Nil(t, err, "CreateSshkey returned an error %v", err)
	assert.Equal(t, fmt.Sprintf("%v", getMockSshkeyCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateSshkey(dummyUUID, SshkeyUpdateRequest{
		Name:   "test",
		Sshkey: "example",
	})
	assert.Nil(t, err, "UpdateSshkey returned an error %v", err)
}

func TestClient_DeleteSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteSshkey(dummyUUID)
	assert.Nil(t, err, "DeleteSshkey returned an error %v", err)
}

func TestClient_GetSshkeyEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareEventListHTTPGet())
	})

	res, err := client.GetSshkeyEventList(dummyUUID)
	assert.Nil(t, err, "GetSshkeyEventList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
}

func getMockSshkey() Sshkey {
	mock := Sshkey{Properties: SshkeyProperties{
		Name:       "test",
		ObjectUUID: dummyUUID,
		Status:     "active",
		CreateTime: dummyTime,
		ChangeTime: dummyTime,
		Sshkey:     "example",
		Labels:     []string{"label"},
		UserUUID:   dummyUUID,
	}}
	return mock
}

func getMockSshkeyCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareSshkeyListHTTPGet() string {
	key := getMockSshkey()
	res, _ := json.Marshal(key.Properties)
	return fmt.Sprintf(`{"sshkeys": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareSshkeyHTTPGet() string {
	key := getMockSshkey()
	res, _ := json.Marshal(key)
	return string(res)
}

func prepareSshkeyCreateResponse() string {
	response := getMockSshkeyCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

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
	if err != nil {
		t.Errorf("GetSshkeyList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkey()), fmt.Sprintf("%v", res))
}

func TestClient_GetSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareSshkeyHTTPGet())
	})
	res, err := client.GetSshkey(dummyUuid)
	if err != nil {
		t.Errorf("GetSshkey returned an error %v", err)
	}
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
	if err != nil {
		t.Errorf("CreateSshkey returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockSshkeyCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateSshkey(dummyUuid, SshkeyUpdateRequest{
		Name:   "test",
		Sshkey: "example",
	})
	if err != nil {
		t.Errorf("UpdateSshkey returned an error %v", err)
	}
}

func TestClient_DeleteSshkey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteSshkey(dummyUuid)
	if err != nil {
		t.Errorf("DeleteSshkey returned an error %v", err)
	}
}

func TestClient_GetSshkeyEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareSshkeyEventListHTTPGet())
	})

	res, err := client.GetSshkeyEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetSshkeyEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkeyEvent()), fmt.Sprintf("%v", res))
}

func getMockSshkey() Sshkey {
	mock := Sshkey{Properties:SshkeyProperties{
		Name:       "test",
		ObjectUuid: dummyUuid,
		Status:     "active",
		CreateTime: dummyTime,
		ChangeTime: dummyTime,
		Sshkey:     "example",
		Labels:     []string{"label"},
		UserUuid:   dummyUuid,
	}}
	return mock
}

func getMockSshkeyCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUuid:  dummyUuid,
		RequestUuid: dummyRequestUUID,
	}
	return mock
}

func getMockSshkeyEvent() SshkeyEvent {
	mock := SshkeyEvent{Properties:SshkeyEventProperties{
		ObjectType:    "type",
		RequestUuid:   dummyRequestUUID,
		ObjectUuid:    dummyUuid,
		Activity:      "login",
		RequestType:   "type",
		RequestStatus: "done",
		Change:        "note",
		Timestamp:     dummyTime,
		UserUuid:      dummyUuid,
	}}
	return mock
}

func prepareSshkeyListHTTPGet() string {
	key := getMockSshkey()
	res, _ := json.Marshal(key.Properties)
	return fmt.Sprintf(`{"sshkeys": {"%s": %s}}`, dummyUuid, string(res))
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

func prepareSshkeyEventListHTTPGet() string {
	event := getMockSshkeyEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

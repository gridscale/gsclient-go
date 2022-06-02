package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetSshkeyList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiSSHKeyBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareSshkeyListHTTPGet())
	})
	res, err := client.GetSSHKeyList(emptyCtx)
	assert.Nil(t, err, "GetSshkeyList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkey("active")), fmt.Sprintf("%v", res))
}

func TestClient_GetSshkey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSSHKeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareSshkeyHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSSHKey(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSshkey("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateSshkey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiSSHKeyBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareSshkeyCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateSSHKey(
			emptyCtx,
			SSHKeyCreateRequest{
				Name:   "test",
				SSHKey: "example",
				Labels: []string{"label"},
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateSshkey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSshkeyCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_UpdateSshkey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiSSHKeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodPatch {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				fmt.Fprint(writer, prepareSshkeyHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdateSSHKey(
				emptyCtx,
				test.testUUID,
				SSHKeyUpdateRequest{
					Name:   "test",
					SSHKey: "example",
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateSshkey returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteSshkey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiSSHKeyBase, dummyUUID)
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
			err := client.DeleteSSHKey(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteSshkey returned an error %v", err)
			}
		}
	}
}

func TestClient_GetSshkeyEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSSHKeyBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(writer, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSSHKeyEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkeyEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
}

func getMockSshkey(status string) SSHKey {
	mock := SSHKey{Properties: SSHKeyProperties{
		Name:       "test",
		ObjectUUID: dummyUUID,
		Status:     status,
		CreateTime: dummyTime,
		ChangeTime: dummyTime,
		SSHKey:     "example",
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
	key := getMockSshkey("active")
	res, _ := json.Marshal(key.Properties)
	return fmt.Sprintf(`{"sshkeys": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareSshkeyHTTPGet(status string) string {
	key := getMockSshkey(status)
	res, _ := json.Marshal(key)
	return string(res)
}

func prepareSshkeyCreateResponse() string {
	response := getMockSshkeyCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

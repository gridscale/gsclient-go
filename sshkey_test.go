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
	res, err := client.GetSshkeyList(emptyCtx)
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
<<<<<<< HEAD
	res, err := client.GetSshkey(dummyUUID)
	if err != nil {
		t.Errorf("GetSshkey returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSshkey(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSshkey("active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockSshkey()), fmt.Sprintf("%v", res))
}

func TestClient_CreateSshkey(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiSshkeyBase
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareSshkeyCreateResponse())
			}
		})
		if clientTest {
			httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
			mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			response, err := client.CreateSshkey(
				emptyCtx,
				SshkeyCreateRequest{
					Name:   "test",
					Sshkey: "example",
					Labels: []string{"label"},
				})
			if isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateSshkey returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockSshkeyCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockSshkeyCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateSshkey(t *testing.T) {
<<<<<<< HEAD
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
	if err != nil {
		t.Errorf("UpdateSshkey returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiSshkeyBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
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
				err := client.UpdateSshkey(
					emptyCtx,
					test.testUUID,
					SshkeyUpdateRequest{
						Name:   "test",
						Sshkey: "example",
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateSshkey returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_DeleteSshkey(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteSshkey(dummyUUID)
	if err != nil {
		t.Errorf("DeleteSshkey returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiSshkeyBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
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
				err := client.DeleteSshkey(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteSshkey returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_GetSshkeyEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareSshkeyEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSshkeyEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkeyEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkeyEvent()), fmt.Sprintf("%v", res))
}

<<<<<<< HEAD
func getMockSshkey() Sshkey {
=======
func TestClient_waitForSSHKeyActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareSshkeyHTTPGet("active"))
	})
	err := client.waitForSSHKeyActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForSSHKeyActive returned an error %v", err)
}

func TestClient_waitForSSHKeyDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(404)
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForSSHKeyDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForSSHKeyDeleted returned an error %v", err)
		}
	}
}

func getMockSshkey(status string) Sshkey {
>>>>>>> 8d4aa0e... add `context`
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

func getMockSshkeyEvent() SshkeyEvent {
	mock := SshkeyEvent{Properties: SshkeyEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "login",
		RequestType:   "type",
		RequestStatus: "done",
		Change:        "note",
		Timestamp:     dummyTime,
		UserUUID:      dummyUUID,
	}}
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

func prepareSshkeyEventListHTTPGet() string {
	event := getMockSshkeyEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

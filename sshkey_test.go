package gsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetSshkeyList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiSshkeyBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareSshkeyListHTTPGet())
	})
	res, err := client.GetSshkeyList(context.Background())
	assert.Nil(t, err, "GetSshkeyList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSshkey("active")), fmt.Sprintf("%v", res))
}

func TestClient_GetSshkey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareSshkeyHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSshkey(context.Background(), test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSshkey("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateSshkey(t *testing.T) {
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
			mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			response, err := client.CreateSshkey(
				context.Background(),
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
	}
}

func TestClient_UpdateSshkey(t *testing.T) {
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
					context.Background(),
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
	}
}

func TestClient_DeleteSshkey(t *testing.T) {
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
				err := client.DeleteSshkey(context.Background(), test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteSshkey returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetSshkeyEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSshkeyBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSshkeyEventList(context.Background(), test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSshkeyEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_waitForSSHKeyActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareSshkeyHTTPGet("in-provisioning"))
			} else {
				fmt.Fprint(w, prepareSshkeyHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			err := client.waitForSSHKeyActive(context.Background(), dummyUUID)
			if isFailed || isTimeout {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForSSHKeyActive returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForSSHKeyDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiSshkeyBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareSshkeyHTTPGet("to-be-deleted"))
			} else {
				w.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, test := range uuidCommonTestCases {
				err := client.waitForSSHKeyDeleted(context.Background(), test.testUUID)
				if test.isFailed || isFailed || isTimeout {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForSSHKeyDeleted returned an error %v", err)
				}
			}
		}
	}
}

func getMockSshkey(status string) Sshkey {
	mock := Sshkey{Properties: SshkeyProperties{
		Name:       "test",
		ObjectUUID: dummyUUID,
		Status:     status,
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

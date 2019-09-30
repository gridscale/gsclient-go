package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerIPList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIPListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerIPList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerIPList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIP()), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIPHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			res, err := client.GetServerIP(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetServerIP returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerIP()), fmt.Sprintf("%v", res))
			}
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerIP()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerIP(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiServerBase, dummyUUID, "ips")
		var isFailed bool
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprint(writer, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "ips", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, http.MethodGet, request.Method)
				fmt.Fprintf(writer, prepareServerIPHTTPGet())
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			for _, testServerID := range uuidCommonTestCases {
				for _, testIPID := range uuidCommonTestCases {
					err := client.CreateServerIP(
						emptyCtx,
						testServerID.testUUID,
						ServerIPRelationCreateRequest{
							ObjectUUID: testIPID.testUUID,
						})
					if testServerID.isFailed || testIPID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "CreateServerIP returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteServerIP(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
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
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			for _, testServerID := range uuidCommonTestCases {
				for _, testIPID := range uuidCommonTestCases {
					err := client.DeleteServerIP(emptyCtx, testServerID.testUUID, testIPID.testUUID)
					if testServerID.isFailed || testIPID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "DeleteServerIP returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_LinkIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
>>>>>>> 8d4aa0e... add `context`
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerIP(dummyUUID, ServerIPRelationCreateRequest{
		ObjectUUID: dummyUUID,
	})
<<<<<<< HEAD
	if err != nil {
		t.Errorf("CreateServerIP returned an error %v", err)
	}
=======
	err := client.LinkIP(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "LinkIP returned an error %v", err)
>>>>>>> 8d4aa0e... add `context`
}

func TestClient_DeleteServerIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
<<<<<<< HEAD
	err := client.DeleteServerIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerIP returned an error %v", err)
	}
=======
	err := client.UnlinkIP(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "DeleteServerIP returned an error %v", err)
>>>>>>> 8d4aa0e... add `context`
}

func TestClient_LinkIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
<<<<<<< HEAD
	err := client.LinkIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("LinkIP returned an error %v", err)
=======
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerIPRelCreation(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerIPRelCreation returned an error %v", err)
			}
		}
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_UnlinkIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
<<<<<<< HEAD
	err := client.UnlinkIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("UnlinkIP returned an error %v", err)
=======
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerIPRelDeleted(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerIPRelCreation returned an error %v", err)
			}
		}
>>>>>>> 8d4aa0e... add `context`
	}
}

func getMockServerIP() ServerIPRelationProperties {
	mock := ServerIPRelationProperties{
		ServerUUID: dummyUUID,
		CreateTime: dummyTime,
		Prefix:     "pre",
		Family:     1,
		ObjectUUID: dummyUUID,
		IP:         "192.168.0.1",
	}
	return mock
}

func prepareServerIPListHTTPGet() string {
	ip := getMockServerIP()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relations": [%s]}`, string(res))
}

func prepareServerIPHTTPGet() string {
	ip := getMockServerIP()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relation": %s}`, string(res))
}

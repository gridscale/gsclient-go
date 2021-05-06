package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServerIPList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
}

func TestClient_GetServerIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
}

func TestClient_CreateServerIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	var isFailed bool
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, "")
		}
	})
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

}

func TestClient_DeleteServerIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
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
}

func TestClient_LinkIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(writer, "")
	})
	err := client.LinkIP(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "LinkIP returned an error %v", err)
}

func TestClient_UnlinkIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if request.Method == http.MethodDelete {
			fmt.Fprintf(writer, "")
		} else if request.Method == http.MethodGet {
			writer.WriteHeader(404)
		}
	})
	err := client.UnlinkIP(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "DeleteServerIP returned an error %v", err)
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

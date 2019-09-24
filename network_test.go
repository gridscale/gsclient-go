package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetNetworkList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet(true, "active"))
	})
	res, err := client.GetNetworkList()
	assert.Nil(t, err, "GetNetworkList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
}

func TestClient_GetNetwork(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetwork(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetwork returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateNetwork(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiNetworkBase
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareNetworkCreateResponse())
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
			response, err := client.CreateNetwork(NetworkCreateRequest{
				Name:         "test",
				Labels:       []string{"label"},
				LocationUUID: dummyUUID,
				L2Security:   false,
			})
			if isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateNetwork returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockNetworkCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
		server.Close()
	}
}

func TestClient_UpdateNetwork(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiNetworkBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodPatch {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					fmt.Fprint(writer, prepareNetworkHTTPGet("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateNetwork(test.testUUID, NetworkUpdateRequest{
					Name:       "test",
					L2Security: false,
				})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateNetwork returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteNetwork(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiNetworkBase, dummyUUID)
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
				err := client.DeleteNetwork(test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteNetwork returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetNetworkEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworkEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetworkEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetNetworkPublic(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isPublicNet bool
	pubNetCases := []bool{true, false}
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareNetworkListHTTPGet(isPublicNet, "active"))
		}
	})
	for _, successFailTest := range commonSuccessFailTestCases {
		isFailed = successFailTest.isFailed
		for _, publicNetTest := range pubNetCases {
			isPublicNet = publicNetTest
			res, err := client.GetNetworkPublic()
			if isFailed || !publicNetTest {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetNetworkPublic returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(publicNetTest, "active")), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_GetNetworksByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet(true, "active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworksByLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetworksByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetDeletedNetworks(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedNetworkListHTTPGet("active"))
	})
	res, err := client.GetDeletedNetworks()
	assert.Nil(t, err, "GetDeletedNetworks returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
}

func TestClient_waitForNetworkActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareNetworkHTTPGet("in-provisioning"))
			} else {
				fmt.Fprint(w, prepareNetworkHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			err := client.waitForNetworkActive(dummyUUID)
			if isFailed || isTimeout {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForNetworkActive returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForNetworkDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareNetworkHTTPGet("to-be-deleted"))
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
				err := client.waitForNetworkDeleted(test.testUUID)
				if test.isFailed || isFailed || isTimeout {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForNetworkDeleted returned an error %v", err)
				}
			}
		}
	}
}

func getMockNetwork(isPublic bool, status string) Network {
	mock := Network{Properties: NetworkProperties{
		LocationCountry: "Germany",
		LocationUUID:    "",
		PublicNet:       isPublic,
		ObjectUUID:      dummyUUID,
		NetworkType:     "",
		Name:            "test",
		Status:          status,
		CreateTime:      dummyTime,
		L2Security:      false,
		ChangeTime:      dummyTime,
		LocationName:    "Cologne",
		DeleteBlock:     false,
		Labels:          nil,
		Relations: NetworkRelations{
			Vlans: []NetworkVlan{
				{
					Vlan:       1,
					TenantName: "test",
					TenantUUID: dummyUUID,
				},
			},
		},
	}}
	return mock
}

func prepareNetworkListHTTPGet(isPublic bool, status string) string {
	network := getMockNetwork(isPublic, status)
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"networks": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareNetworkHTTPGet(status string) string {
	network := getMockNetwork(true, status)
	res, _ := json.Marshal(network)
	return string(res)
}

func getMockNetworkCreateResponse() NetworkCreateResponse {
	mock := NetworkCreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareNetworkCreateResponse() string {
	createResponse := getMockNetworkCreateResponse()
	res, _ := json.Marshal(createResponse)
	return string(res)
}

func prepareDeletedNetworkListHTTPGet(status string) string {
	network := getMockNetwork(true, status)
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"deleted_networks": {"%s": %s}}`, dummyUUID, string(res))
}

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
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet())
	})
	res, err := client.GetNetworkList(emptyCtx)
	assert.Nil(t, err, "GetNetworkList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_GetNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkHTTPGet())
	})
<<<<<<< HEAD
	res, err := client.GetNetwork(dummyUUID)
	if err != nil {
		t.Errorf("GetNetwork returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetwork(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetwork returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_CreateNetwork(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareNetworkCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateNetwork(NetworkCreateRequest{
		Name:         "test",
		Labels:       []string{"label"},
		LocationUUID: dummyUUID,
		L2Security:   false,
	})
	if err != nil {
		t.Errorf("CreateNetwork returned an error %v", err)
=======
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
			mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			response, err := client.CreateNetwork(
				emptyCtx,
				NetworkCreateRequest{
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
>>>>>>> 8d4aa0e... add `context`
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockNetworkCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateNetwork(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateNetwork(dummyUUID, NetworkUpdateRequest{
		Name:       "test",
		L2Security: false,
	})
	if err != nil {
		t.Errorf("UpdateNetwork returned an error %v", err)
=======
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
				err := client.UpdateNetwork(
					emptyCtx,
					test.testUUID,
					NetworkUpdateRequest{
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
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_DeleteNetwork(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteNetwork(dummyUUID)
	if err != nil {
		t.Errorf("DeleteNetwork returned an error %v", err)
=======
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
				err := client.DeleteNetwork(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteNetwork returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_GetNetworkEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworkEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetworkEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetworkEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetNetworkPublic(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet())
	})
<<<<<<< HEAD
	res, err := client.GetNetworkPublic()
	if err != nil {
		t.Errorf("GetNetworkPublic returned an error %v", err)
=======
	for _, successFailTest := range commonSuccessFailTestCases {
		isFailed = successFailTest.isFailed
		for _, publicNetTest := range pubNetCases {
			isPublicNet = publicNetTest
			res, err := client.GetNetworkPublic(emptyCtx)
			if isFailed || !publicNetTest {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetNetworkPublic returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(publicNetTest, "active")), fmt.Sprintf("%v", res))
			}
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockNetwork()), fmt.Sprintf("%v", res))
}

<<<<<<< HEAD
func getMockNetwork() Network {
=======
func TestClient_GetNetworksByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet(true, "active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworksByLocation(emptyCtx, test.testUUID)
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
	res, err := client.GetDeletedNetworks(emptyCtx)
	assert.Nil(t, err, "GetDeletedNetworks returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true, "active")), fmt.Sprintf("%v", res))
}

func TestClient_waitForNetworkActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareNetworkHTTPGet("active"))
	})
	err := client.waitForNetworkActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForNetworkActive returned an error %v", err)
}

func TestClient_waitForNetworkDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(404)
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForNetworkDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForNetworkDeleted returned an error %v", err)
		}
	}
}

func getMockNetwork(isPublic bool, status string) Network {
>>>>>>> 8d4aa0e... add `context`
	mock := Network{Properties: NetworkProperties{
		LocationCountry: "Germany",
		LocationUUID:    "",
		PublicNet:       true,
		ObjectUUID:      dummyUUID,
		NetworkType:     "",
		Name:            "test",
		Status:          "active",
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

func prepareNetworkListHTTPGet() string {
	network := getMockNetwork()
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"networks": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareNetworkHTTPGet() string {
	network := getMockNetwork()
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

func getMockNetworkEvent() NetworkEvent {
	mock := NetworkEvent{Properties: NetworkEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "activity",
		RequestType:   "tcp",
		RequestStatus: "done",
		Change:        "change note",
		Timestamp:     dummyTime,
		UserUUID:      dummyUUID,
	}}
	return mock
}

func prepareNetworkEventListHTTPGet() string {
	event := getMockNetworkEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

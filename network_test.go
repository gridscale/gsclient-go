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
	res, err := client.GetNetworkList()
	if err != nil {
		t.Errorf("GetNetworkList returned an error %v", err)
	}
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
	res, err := client.GetNetwork(dummyUUID)
	if err != nil {
		t.Errorf("GetNetwork returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_CreateNetwork(t *testing.T) {
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
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockNetworkCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateNetwork(t *testing.T) {
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
	}
}

func TestClient_DeleteNetwork(t *testing.T) {
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
	res, err := client.GetNetworkEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetNetworkEventList returned an error %v", err)
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
	res, err := client.GetNetworkPublic()
	if err != nil {
		t.Errorf("GetNetworkPublic returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockNetwork()), fmt.Sprintf("%v", res))
}

func TestClient_GetNetworksByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet())
	})
	res, err := client.GetNetworksByLocation(dummyUUID)
	if err != nil {
		t.Errorf("GetNetworksByLocation returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork()), fmt.Sprintf("%v", res))
}

func getMockNetwork() Network {
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

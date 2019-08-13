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
	uri := path.Join(apiNetworkBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkHTTPGet())
	})
	res, err := client.GetNetwork(dummyUuid)
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
		LocationUuid: dummyUuid,
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
	uri := path.Join(apiNetworkBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateNetwork(dummyUuid, NetworkUpdateRequest{
		Name:       "test",
		Labels:     []string{"label"},
		L2Security: false,
	})
	if err != nil {
		t.Errorf("UpdateNetwork returned an error %v", err)
	}
}

func TestClient_DeleteNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteNetwork(dummyUuid)
	if err != nil {
		t.Errorf("DeleteNetwork returned an error %v", err)
	}
}

func TestClient_GetNetworkEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkEventListHTTPGet())
	})
	res, err := client.GetNetworkEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetNetworkEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetworkEvent()), fmt.Sprintf("%v", res))
}

func getMockNetwork() Network {
	mock := Network{Properties:NetworkProperties{
		LocationCountry: "Germany",
		LocationUuid:    "",
		PublicNet:       true,
		ObjectUuid:      dummyUuid,
		NetworkType:     "",
		Name:            "test",
		Status:          "active",
		CreateTime:      dummyTime,
		L2Security:      false,
		ChangeTime:      dummyTime,
		LocationName:    "Cologne",
		DeleteBlock:     false,
		Labels:          nil,
		Relations:       NetworkRelations{
			Vlans: []NetworkVlan{
				{
					Vlan:       1,
					TenantName: "test",
					TenantUuid: dummyUuid,
				},
			},
		},
	}}
	return mock
}

func prepareNetworkListHTTPGet() string {
	network := getMockNetwork()
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"networks": {"%s": %s}}`, dummyUuid, string(res))
}

func prepareNetworkHTTPGet() string {
	network := getMockNetwork()
	res, _ := json.Marshal(network)
	return string(res)
}

func getMockNetworkCreateResponse() NetworkCreateResponse {
	mock := NetworkCreateResponse{
		ObjectUuid:  dummyUuid,
		RequestUuid: dummyRequestUUID,
	}
	return mock
}

func prepareNetworkCreateResponse() string {
	createResponse := getMockNetworkCreateResponse()
	res, _ := json.Marshal(createResponse)
	return string(res)
}

func getMockNetworkEvent() NetworkEvent {
	mock := NetworkEvent{Properties:NetworkEventProperties{
		ObjectType:    "type",
		RequestUuid:   dummyRequestUUID,
		ObjectUuid:    dummyUuid,
		Activity:      "activity",
		RequestType:   "tcp",
		RequestStatus: "done",
		Change:        "change note",
		Timestamp:     dummyTime,
		UserUuid:      dummyUuid,
	}}
	return mock
}

func prepareNetworkEventListHTTPGet() string {
	event := getMockNetworkEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

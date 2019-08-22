package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetIpList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIpBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIpListHTTPGet())
	})
	res, err := client.GetIpList()
	if err != nil {
		t.Errorf("GetIpList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIpBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIpHTTPGet())
	})
	res, err := client.GetIp(dummyUUID)
	if err != nil {
		t.Errorf("GetIp returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_CreateIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIpBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareIpCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateIp(IpCreateRequest{
		Name:         "test",
		Family:       1,
		LocationUUID: dummyUUID,
		Failover:     false,
		ReverseDns:   "8.8.8.8",
	})
	if err != nil {
		t.Errorf("CreateIp returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockIpCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIpBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateIp(dummyUUID, IpUpdateRequest{
		Name:       "test",
		Failover:   false,
		ReverseDns: "8.8.4.4",
	})
	if err != nil {
		t.Errorf("UpdateIp returned an error %v", err)
	}
}

func TestClient_DeleteIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIpBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteIp(dummyUUID)
	if err != nil {
		t.Errorf("DeleteIp returned an error %v", err)
	}
}

func TestClient_GetIpEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIpBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIpEventListHTTPGet())
	})
	res, err := client.GetIpEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetIpEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIpEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetIpVersion(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIpBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIpHTTPGet())
	})
	res := client.GetIpVersion(dummyUUID)
	if res == 0 {
		t.Error("GetIpVersion has an error")
	}
	assert.Equal(t, 1, res)

}

func getMockIP() Ip {
	mock := Ip{Properties: IpProperties{
		Name:            "test",
		LocationCountry: "Germany",
		LocationUUID:    dummyUUID,
		ObjectUUID:      dummyUUID,
		ReverseDns:      "8.8.8.8",
		Family:          1,
		Status:          "active",
		CreateTime:      dummyTime,
		Failover:        false,
		ChangeTime:      dummyTime,
		LocationIata:    "",
		LocationName:    "Cologne",
		Prefix:          "",
		Ip:              "192.168.0.1",
		DeleteBlock:     "",
		UsagesInMinutes: 10,
		CurrentPrice:    0.9,
		Labels:          []string{"label"},
		Relations: IpRelations{
			Loadbalancers: []IpLoadbalancer{
				{
					CreateTime:       dummyTime,
					LoadbalancerName: "test",
					LoadbalancerUUID: dummyUUID,
				},
			},
		},
	}}
	return mock
}

func prepareIpListHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"ips": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareIpHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip)
	return string(res)
}

func getMockIpCreateResponse() IpCreateResponse {
	mock := IpCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
		Prefix:      "ip",
		Ip:          "192.168.0.1",
	}
	return mock
}

func prepareIpCreateResponse() string {
	res, _ := json.Marshal(getMockIpCreateResponse())
	return string(res)
}

func getMockIpEvent() IpEvent {
	mock := IpEvent{Properties: IpEventProperties{
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

func prepareIpEventListHTTPGet() string {
	event := getMockIpEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

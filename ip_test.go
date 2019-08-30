package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetIPList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet())
	})
	res, err := client.GetIPList()
	if err != nil {
		t.Errorf("GetIPList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPHTTPGet())
	})
	res, err := client.GetIP(dummyUUID)
	if err != nil {
		t.Errorf("GetIP returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_CreateIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareIPCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateIP(IPCreateRequest{
		Name:         "test",
		Family:       1,
		LocationUUID: dummyUUID,
		Failover:     false,
		ReverseDNS:   "8.8.8.8",
	})
	if err != nil {
		t.Errorf("CreateIP returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockIPCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateIP(dummyUUID, IPUpdateRequest{
		Name:       "test",
		Failover:   false,
		ReverseDNS: "8.8.4.4",
	})
	if err != nil {
		t.Errorf("UpdateIP returned an error %v", err)
	}
}

func TestClient_DeleteIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteIP(dummyUUID)
	if err != nil {
		t.Errorf("DeleteIP returned an error %v", err)
	}
}

func TestClient_GetIPEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareEventListHTTPGet())
	})
	res, err := client.GetIPEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetIPEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetIPVersion(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPHTTPGet())
	})
	res := client.GetIPVersion(dummyUUID)
	if res == 0 {
		t.Error("GetIPVersion has an error")
	}
	assert.Equal(t, 1, res)

}

func TestClient_GetIPsByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet())
	})
	res, err := client.GetIPsByLocation(dummyUUID)
	if err != nil {
		t.Errorf("GetIPsByLocation returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetDeletedIPs(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedIPListHTTPGet())
	})
	res, err := client.GetDeletedIPs()
	if err != nil {
		t.Errorf("GetDeletedIPs returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func getMockIP() IP {
	mock := IP{Properties: IPProperties{
		Name:            "test",
		LocationCountry: "Germany",
		LocationUUID:    dummyUUID,
		ObjectUUID:      dummyUUID,
		ReverseDNS:      "8.8.8.8",
		Family:          1,
		Status:          "active",
		CreateTime:      dummyTime,
		Failover:        false,
		ChangeTime:      dummyTime,
		LocationIata:    "",
		LocationName:    "Cologne",
		Prefix:          "",
		IP:              "192.168.0.1",
		DeleteBlock:     "",
		UsagesInMinutes: 10,
		CurrentPrice:    0.9,
		Labels:          []string{"label"},
		Relations: IPRelations{
			Loadbalancers: []IPLoadbalancer{
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

func prepareIPListHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"ips": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareIPHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip)
	return string(res)
}

func getMockIPCreateResponse() IPCreateResponse {
	mock := IPCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
		Prefix:      "ip",
		IP:          "192.168.0.1",
	}
	return mock
}

func prepareIPCreateResponse() string {
	res, _ := json.Marshal(getMockIPCreateResponse())
	return string(res)
}

func prepareDeletedIPListHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"deleted_ips": {"%s": %s}}`, dummyUUID, string(res))
}

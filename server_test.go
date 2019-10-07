package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServerList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiServerBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerListHTTPGet())
	})
	res, err := client.GetServerList()
	if err != nil {
		t.Errorf("GetServerList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServer(true)), fmt.Sprintf("%v", res))
}

func TestClient_GetServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(true))
	})
	res, err := client.GetServer(dummyUUID)
	if err != nil {
		t.Errorf("GetServer returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServer(true)), fmt.Sprintf("%v", res))
}

func TestClient_CreateServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiServerBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareServerCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateServer(ServerCreateRequest{
		Name:            "test",
		Memory:          10,
		Cores:           4,
		LocationUUID:    dummyUUID,
		HardwareProfile: "default",
		AvailablityZone: "",
		Labels:          []string{"label"},
		Relations:       &ServerCreateRequestRelations{},
	})
	if err != nil {
		t.Errorf("CreateServer returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockServerCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateServer(dummyUUID, ServerUpdateRequest{
		Name:            "test",
		AvailablityZone: "test zone",
		Memory:          4,
		Cores:           2,
		Labels:          nil,
	})
	if err != nil {
		t.Errorf("UpdateServer returned an error %v", err)
	}
}

func TestClient_DeleteServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServer(dummyUUID)
	if err != nil {
		t.Errorf("DeleteServer returned an error %v", err)
	}
}

func TestClient_GetServerEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerEventListHTTPGet())
	})
	res, err := client.GetServerEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerMetricList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "metrics")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerMetricListHTTPGet())
	})
	res, err := client.GetServerMetricList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerMetricList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerMetric()), fmt.Sprintf("%v", res))
}

func TestClient_IsServerOn(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(true))
	})
	isOn, err := client.IsServerOn(dummyUUID)
	if err != nil {
		t.Errorf("IsServerOn returned an error %v", err)
	}
	assert.Equal(t, true, isOn)
}

func TestClient_setServerPowerState(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	power := true
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
	mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = false
		fmt.Fprint(writer, "")
	})
	err := client.setServerPowerState(dummyUUID, false)
	if err != nil {
		t.Errorf("turnOnOffServer returned an error %v", err)
	}
}

func TestClient_StartServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	power := false
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
	mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = true
		fmt.Fprint(writer, "")
	})
	err := client.StartServer(dummyUUID)
	if err != nil {
		t.Errorf("StartServer returned an error %v", err)
	}
}

func TestClient_StopServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	power := true
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
	mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = false
		fmt.Fprint(writer, "")
	})
	err := client.StopServer(dummyUUID)
	if err != nil {
		t.Errorf("StopServer returned an error %v", err)
	}
}

func TestClient_ShutdownServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	power := true
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
	mux.HandleFunc(uri+"/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = false
		writer.WriteHeader(http.StatusInternalServerError)
 		writer.Write([]byte("☄ HTTP status code returned!"))
		fmt.Fprint(writer, "")
	})

	err := client.ShutdownServer(dummyUUID)
	if err != nil {
		t.Errorf("ShutdownServer returned an error %v", err)
	}
}

func getMockServer(power bool) Server {
	mock := Server{Properties: ServerProperties{
		ObjectUUID:           dummyUUID,
		Name:                 "Test",
		Memory:               2,
		Cores:                4,
		HardwareProfile:      "default",
		Status:               "active",
		LocationUUID:         dummyUUID,
		Power:                power,
		CurrentPrice:         9.5,
		AvailablityZone:      "",
		AutoRecovery:         true,
		Legacy:               false,
		ConsoleToken:         "",
		UsageInMinutesMemory: 47331,
		UsageInMinutesCores:  17476,
		Labels:               []string{"label"},
		Relations: ServerRelations{
			IsoImages: []ServerIsoImageRelationProperties{
				{
					ObjectUUID: dummyUUID,
					ObjectName: "test",
					Private:    false,
					CreateTime: dummyTime,
				},
			},
		},
	}}
	return mock
}

func getMockServerCreateResponse() ServerCreateResponse {
	mock := ServerCreateResponse{
		ObjectUUID:   dummyUUID,
		RequestUUID:  dummyRequestUUID,
		SeverUUID:    dummyUUID,
		NetworkUUIDs: nil,
		StorageUUIDs: nil,
		IPaddrUUIDs:  nil,
	}
	return mock
}

func getMockServerEvent() ServerEvent {
	mock := ServerEvent{Properties: ServerEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "activity",
		RequestType:   "turn on",
		RequestStatus: "done",
		Change:        "change note",
		Timestamp:     dummyTime,
		UserUUID:      dummyUUID,
	}}
	return mock
}

func getMockServerMetric() ServerMetric {
	mock := ServerMetric{Properties: ServerMetricProperties{
		BeginTime:       dummyTime,
		EndTime:         dummyTime,
		PaaSServiceUUID: dummyUUID,
		CoreUsage: struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		}{
			Value: 50.5,
			Unit:  "percentage",
		},
		StorageSize: struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		}{
			Value: 10.5,
			Unit:  "GB",
		},
	}}
	return mock
}

func prepareServerListHTTPGet() string {
	server := getMockServer(true)
	res, _ := json.Marshal(server.Properties)
	return fmt.Sprintf(`{"servers": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareServerHTTPGet(power bool) string {
	server := getMockServer(power)
	res, _ := json.Marshal(server)
	return string(res)
}

func prepareServerCreateResponse() string {
	server := getMockServerCreateResponse()
	res, _ := json.Marshal(server)
	return string(res)
}

func prepareServerEventListHTTPGet() string {
	event := getMockServerEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

func prepareServerMetricListHTTPGet() string {
	metric := getMockServerMetric()
	res, _ := json.Marshal(metric.Properties)
	return fmt.Sprintf(`{"server_metrics": [%s]}`, string(res))
}

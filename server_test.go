package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
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
	uri := path.Join(apiServerBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(true))
	})
	res, err := client.GetServer(dummyUuid)
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
		LocationUuid:    dummyUuid,
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
	uri := path.Join(apiServerBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateServer(dummyUuid, ServerUpdateRequest{
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
	uri := path.Join(apiServerBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServer(dummyUuid)
	if err != nil {
		t.Errorf("DeleteServer returned an error %v", err)
	}
}

func TestClient_GetServerEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerEventListHTTPGet())
	})
	res, err := client.GetServerEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetServerEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerMetricList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "metrics")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerMetricListHTTPGet())
	})
	res, err := client.GetServerMetricList(dummyUuid)
	if err != nil {
		t.Errorf("GetServerMetricList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerMetric()), fmt.Sprintf("%v", res))
}

func TestClient_IsServerOn(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(true))
	})
	isOn, err := client.IsServerOn(dummyUuid)
	if err != nil {
		t.Errorf("IsServerOn returned an error %v", err)
	}
	assert.Equal(t, true, isOn)
}

func TestClient_setServerPowerState(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid)
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
	err := client.setServerPowerState(dummyUuid, false)
	if err != nil {
		t.Errorf("turnOnOffServer returned an error %v", err)
	}
}

func TestClient_StartServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid)
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
	err := client.StartServer(dummyUuid)
	if err != nil {
		t.Errorf("StartServer returned an error %v", err)
	}
}

func TestClient_StopServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid)
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
	err := client.StopServer(dummyUuid)
	if err != nil {
		t.Errorf("StopServer returned an error %v", err)
	}
}

func TestClient_ShutdownServer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid)
	power := true
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
	mux.HandleFunc(uri+"/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = false
		fmt.Fprint(writer, "")
	})
	err := client.ShutdownServer(dummyUuid)
	if err != nil {
		t.Errorf("ShutdownServer returned an error %v", err)
	}
}

func getMockServer(power bool) Server {
	mock := Server{Properties: ServerProperties{
		ObjectUuid:           dummyUuid,
		Name:                 "Test",
		Memory:               2,
		Cores:                4,
		HardwareProfile:      "default",
		Status:               "active",
		LocationUuid:         dummyUuid,
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
					ObjectUuid: dummyUuid,
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
		ObjectUuid:   dummyUuid,
		RequestUuid:  dummyRequestUUID,
		SeverUuid:    dummyUuid,
		NetworkUuids: nil,
		StorageUuids: nil,
		IpaddrUuids:  nil,
	}
	return mock
}

func getMockServerEvent() ServerEvent {
	mock := ServerEvent{Properties: ServerEventProperties{
		ObjectType:    "type",
		RequestUuid:   dummyRequestUUID,
		ObjectUuid:    dummyUuid,
		Activity:      "activity",
		RequestType:   "turn on",
		RequestStatus: "done",
		Change:        "change note",
		Timestamp:     dummyTime,
		UserUuid:      dummyUuid,
	}}
	return mock
}

func getMockServerMetric() ServerMetric {
	mock := ServerMetric{Properties: ServerMetricProperties{
		BeginTime:       dummyTime,
		EndTime:         dummyTime,
		PaaSServiceUuid: dummyUuid,
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
	return fmt.Sprintf(`{"servers": {"%s": %s}}`, dummyUuid, string(res))
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

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
	res, err := client.GetServerList(emptyCtx)
	assert.Nil(t, err, "GetServerList returned an error %v", err)
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
<<<<<<< HEAD
	res, err := client.GetServer(dummyUUID)
	if err != nil {
		t.Errorf("GetServer returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServer(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServer returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockServer(true, "active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServer(true)), fmt.Sprintf("%v", res))
}

func TestClient_CreateServer(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiServerBase
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareServerCreateResponse())
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
			response, err := client.CreateServer(
				emptyCtx,
				ServerCreateRequest{
					Name:            "test",
					Memory:          10,
					Cores:           4,
					LocationUUID:    dummyUUID,
					HardwareProfile: DefaultServerHardware,
					AvailablityZone: "",
					Labels:          []string{"label"},
				})
			if test.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateServer returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockServerCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateServer(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodPatch {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					fmt.Fprint(writer, prepareServerHTTPGet(true, "active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateServer(
					emptyCtx,
					test.testUUID,
					ServerUpdateRequest{
						Name:            "test",
						AvailablityZone: "test zone",
						Memory:          4,
						Cores:           2,
						Labels:          nil,
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateServer returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_DeleteServer(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID)
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
				err := client.DeleteServer(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteServer returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
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
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
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
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerMetricList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerMetricList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerMetric()), fmt.Sprintf("%v", res))
		}
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
	isOn, err := client.IsServerOn(emptyCtx, dummyUUID)
	assert.Nil(t, err, "IsServerOn returned an error %v", err)
	assert.Equal(t, true, isOn)
}

func TestClient_setServerPowerState(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiServerBase, dummyUUID)
		power := true
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodGet, request.Method)
			fmt.Fprintf(writer, prepareServerHTTPGet(power, "active"))
		})
		mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPatch, request.Method)
			power = false
			fmt.Fprint(writer, "")
		})
		err := client.setServerPowerState(emptyCtx, dummyUUID, false)
		assert.Nil(t, err, "turnOnOffServer returned an error %v", err)
		server.Close()
	}
}

func TestClient_StartServer(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiServerBase, dummyUUID)
		power := false
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodGet, request.Method)
			fmt.Fprintf(writer, prepareServerHTTPGet(power, "active"))
		})
		mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPatch, request.Method)
			power = true
			fmt.Fprint(writer, "")
		})
		err := client.StartServer(emptyCtx, dummyUUID)
		assert.Nil(t, err, "StartServer returned an error %v", err)
		server.Close()
	}
}

func TestClient_StopServer(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiServerBase, dummyUUID)
		power := true
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodGet, request.Method)
			fmt.Fprintf(writer, prepareServerHTTPGet(power, "active"))
		})
		mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPatch, request.Method)
			power = false
			fmt.Fprint(writer, "")
		})
		err := client.StopServer(emptyCtx, dummyUUID)
		assert.Nil(t, err, "StopServer returned an error %v", err)
		server.Close()
	}
}

func TestClient_ShutdownServer(t *testing.T) {
	shutdownSuccessTestCases := []bool{true, false}
	for _, clientTest := range syncClientTestCases {
		for _, testCaseShutdownSuccess := range shutdownSuccessTestCases {
			if testCaseShutdownSuccess {
				server, client, mux := setupTestClient(clientTest)
				uri := path.Join(apiServerBase, dummyUUID)
				power := true
				retries := 0
				mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodGet, request.Method)
					fmt.Fprintf(writer, prepareServerHTTPGet(power, "active"))
				})
				mux.HandleFunc(uri+"/shutdown", func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPatch, request.Method)
					if retries < 5 {
						retries++
						writer.WriteHeader(http.StatusInternalServerError)
						writer.Write([]byte("☄ HTTP status code returned!"))
						return
					}
					power = false
					fmt.Fprint(writer, "")
				})
				mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPatch, request.Method)
					power = false
					fmt.Fprint(writer, "")
				})
				err := client.ShutdownServer(emptyCtx, dummyUUID)
				assert.Nil(t, err, "ShutdownServer returned an error %v", err)
				server.Close()
			} else {
				server, client, mux := setupTestClient(true)
				uri := path.Join(apiServerBase, dummyUUID)
				power := true
				mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodGet, request.Method)
					fmt.Fprintf(writer, prepareServerHTTPGet(power, "active"))
				})
				mux.HandleFunc(uri+"/shutdown", func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPatch, request.Method)
					writer.WriteHeader(http.StatusInternalServerError)
					writer.Write([]byte("☄ HTTP status code returned!"))
				})
				mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
					assert.Equal(t, http.MethodPatch, request.Method)
					power = false
					fmt.Fprint(writer, "")
				})
				err := client.ShutdownServer(emptyCtx, dummyUUID)
				assert.Nil(t, err, "ShutdownServer returned an error %v", err)
				server.Close()
			}
		}
	}
}

func TestClient_GetServersByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
>>>>>>> 8d4aa0e... add `context`
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	power := true
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerHTTPGet(power))
	})
<<<<<<< HEAD
	mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = false
		fmt.Fprint(writer, "")
	})
	err := client.setServerPowerState(dummyUUID, false)
	if err != nil {
		t.Errorf("turnOnOffServer returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServersByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServersByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServer(true, "active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	mux.HandleFunc(uri+"/power", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		power = true
		fmt.Fprint(writer, "")
	})
	err := client.StartServer(dummyUUID)
	if err != nil {
		t.Errorf("StartServer returned an error %v", err)
	}
=======
	res, err := client.GetDeletedServers(emptyCtx)
	assert.Nil(t, err, "GetDeletedServers returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServer(true, "active")), fmt.Sprintf("%v", res))
}

func TestClient_waitForServerPowerStatus(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareServerHTTPGet(true, "active"))
	})
	err := client.waitForServerPowerStatus(emptyCtx, dummyUUID, true)
	assert.Nil(t, err, "waitForServerPowerStatus returned an error %v", err)
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	err := client.StopServer(dummyUUID)
	if err != nil {
		t.Errorf("StopServer returned an error %v", err)
	}
=======
	err := client.waitForServerActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForServerActive returned an error %v", err)
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
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
=======
	for _, test := range uuidCommonTestCases {
		err := client.waitForServerDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForServerDeleted returned an error %v", err)
		}
>>>>>>> 8d4aa0e... add `context`
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

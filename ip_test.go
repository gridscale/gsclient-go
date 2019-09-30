package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetIPList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet())
	})
	res, err := client.GetIPList(emptyCtx)
	assert.Nil(t, err, "GetIPList returned an error %v", err)
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
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIP(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIP returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockIP("active")), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_CreateIP(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiIPBase
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareIPCreateResponse())
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
			response, err := client.CreateIP(
				emptyCtx,
				IPCreateRequest{
					Name:         "test",
					Family:       IPv4Type,
					LocationUUID: dummyUUID,
					Failover:     false,
					ReverseDNS:   "8.8.8.8",
				})
			if isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateIP returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockIPCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
		server.Close()
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockIPCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateIP(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiIPBase, dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodPatch {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					fmt.Fprint(writer, prepareIPHTTPGet("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateIP(
					emptyCtx,
					test.testUUID,
					IPUpdateRequest{
						Name:       "test",
						Failover:   false,
						ReverseDNS: "8.8.4.4",
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateIP returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteIP(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiIPBase, dummyUUID)
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
				err := client.DeleteIP(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteIP returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetIPEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIPEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIPEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIPEvent()), fmt.Sprintf("%v", res))
}

func TestClient_GetIPVersion(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareIPHTTPGet("active"))
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res := client.GetIPVersion(emptyCtx, dummyUUID)
		if test.isFailed {
			assert.Equal(t, 0, res)
		} else {
			assert.Equal(t, 1, res)
		}
	}
}

func TestClient_GetIPsByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIPsByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIPsByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockIP("active")), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, res)

func TestClient_GetDeletedIPs(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedIPListHTTPGet("deleted"))
	})
	res, err := client.GetDeletedIPs(emptyCtx)
	assert.Nil(t, err, "GetDeletedIPs returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP("deleted")), fmt.Sprintf("%v", res))
}

func TestClient_waitForIPActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareIPHTTPGet("active"))
	})
	err := client.waitForIPActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForIPActive returned an error %v", err)
}

func TestClient_waitForIPDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(404)
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForIPDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForIPDeleted returned an error %v", err)
		}
	}
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

func getMockIPEvent() IPEvent {
	mock := IPEvent{Properties: IPEventProperties{
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

func prepareIPEventListHTTPGet() string {
	event := getMockIPEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

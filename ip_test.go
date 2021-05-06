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
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareIPListHTTPGet("active"))
	})
	res, err := client.GetIPList(emptyCtx)
	assert.Nil(t, err, "GetIPList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP("active")), fmt.Sprintf("%v", res))
}

func TestClient_GetIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareIPHTTPGet("active"))
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
}

func TestClient_CreateIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareIPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateIP(
			emptyCtx,
			IPCreateRequest{
				Name:       "test",
				Family:     IPv4Type,
				Failover:   false,
				ReverseDNS: "8.8.8.8",
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateIP returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockIPCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_UpdateIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
}

func TestClient_DeleteIP(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
}

func TestClient_GetIPEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareEventListHTTPGet())
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

}

func TestClient_GetIPVersion(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
}

func TestClient_GetDeletedIPs(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareDeletedIPListHTTPGet("deleted"))
	})
	res, err := client.GetDeletedIPs(emptyCtx)
	assert.Nil(t, err, "GetDeletedIPs returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP("deleted")), fmt.Sprintf("%v", res))
}

func getMockIP(status string) IP {
	mock := IP{Properties: IPProperties{
		Name:            "test",
		LocationCountry: "Germany",
		LocationUUID:    dummyUUID,
		ObjectUUID:      dummyUUID,
		ReverseDNS:      "8.8.8.8",
		Family:          1,
		Status:          status,
		CreateTime:      dummyTime,
		Failover:        false,
		ChangeTime:      dummyTime,
		LocationIata:    "",
		LocationName:    "Cologne",
		Prefix:          "",
		IP:              "192.168.0.1",
		DeleteBlock:     false,
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

func prepareIPListHTTPGet(status string) string {
	ip := getMockIP(status)
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"ips": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareIPHTTPGet(status string) string {
	ip := getMockIP(status)
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

func prepareDeletedIPListHTTPGet(status string) string {
	ip := getMockIP(status)
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"deleted_ips": {"%s": %s}}`, dummyUUID, string(res))
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		fmt.Fprint(w, prepareLoadBalancerHTTPCreateResponse())
	})

<<<<<<< HEAD
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	lb := getMockLoadbalancer().Properties
	lbRequest := LoadBalancerCreateRequest{
		Name:                lb.Name,
		Algorithm:           lb.Algorithm,
		LocationUUID:        lb.LocationUUID,
		ListenIPv6UUID:      lb.ListenIPv6UUID,
		ListenIPv4UUID:      lb.ListenIPv4UUID,
		RedirectHTTPToHTTPS: lb.RedirectHTTPToHTTPS,
		ForwardingRules:     lb.ForwardingRules,
		BackendServers:      lb.BackendServers,
		Labels:              lb.Labels,
	}
	response, err := client.CreateLoadBalancer(lbRequest)
	if err != nil {
		t.Errorf("CreateLoadBalancer returned error: %v", err)
=======
		if clientTest {
			httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
			mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, testSuccessFail := range commonSuccessFailTestCases {
			isFailed = testSuccessFail.isFailed
			for _, testLabel := range labelSlices {
				lbRequest := LoadBalancerCreateRequest{
					Name:                lb.Name,
					Algorithm:           LoadbalancerLeastConnAlg,
					LocationUUID:        lb.LocationUUID,
					ListenIPv6UUID:      lb.ListenIPv6UUID,
					ListenIPv4UUID:      lb.ListenIPv4UUID,
					RedirectHTTPToHTTPS: lb.RedirectHTTPToHTTPS,
					ForwardingRules:     lb.ForwardingRules,
					BackendServers:      lb.BackendServers,
					Labels:              testLabel,
				}
				response, err := client.CreateLoadBalancer(emptyCtx, lbRequest)
				if testSuccessFail.isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "CreateLoadBalancer returned error: %v", err)
					assert.Equal(t, fmt.Sprintf("%s", prepareLoadBalancerObjectCreateResponse()), fmt.Sprintf("%s", response))
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%s", prepareLoadBalancerObjectCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_GetLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	expectedObject := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse())
	})
	for _, test := range uuidCommonTestCases {
		loadbalancer, err := client.GetLoadBalancer(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancer returned error: %v", err)
			assert.Equal(t, fmt.Sprintf("%v", expectedObject.Properties), fmt.Sprintf("%v", loadbalancer.Properties))
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", expectedObject.Properties), fmt.Sprintf("%v", loadbalancer.Properties))
}

func TestClient_GetLoadBalancerList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase)
	expectedObjects := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPListResponse())
	})
	loadbalancers, err := client.GetLoadBalancerList(emptyCtx)
	assert.Nil(t, err, "GetLoadBalancerList returned error: %v", err)
	assert.Equal(t, 1, len(loadbalancers))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", loadbalancers))
}

func TestClient_UpdateLoadBalancer(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateLoadBalancer(dummyUUID, LoadBalancerUpdateRequest{
		Name:                "test",
		ListenIPv6UUID:      dummyUUID,
		ListenIPv4UUID:      dummyUUID,
		RedirectHTTPToHTTPS: false,
		Status:              "inactive",
	})
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiLoadBalancerBase, dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			if isFailed {
				w.WriteHeader(400)
			} else {
				if r.Method == http.MethodPatch {
					fmt.Fprintf(w, "")
				} else if r.Method == http.MethodGet {
					fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateLoadBalancer(
					emptyCtx,
					test.testUUID,
					LoadBalancerUpdateRequest{
						Name:                "test",
						ListenIPv6UUID:      dummyUUID,
						ListenIPv4UUID:      dummyUUID,
						RedirectHTTPToHTTPS: false,
						Status:              "inactive",
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateLoadBalancer returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_DeleteLoadBalancer(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteLoadBalancer(dummyUUID)
	if err != nil {
		t.Errorf("DeleteLoadBalancer returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiLoadBalancerBase, dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			if isFailed {
				w.WriteHeader(400)
			} else {
				if r.Method == http.MethodDelete {
					fmt.Fprintf(w, "")
				} else if r.Method == http.MethodGet {
					w.WriteHeader(404)
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.DeleteLoadBalancer(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteLoadBalancer returned an error %v", err)
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_GetLoadBalancerEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID, "events")
	expectedObjects := getMockLoadBalancerEvent()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetLoadBalancerEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancerEventList returned error: %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", response))
		}
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", response))
}

<<<<<<< HEAD
func getMockLoadbalancer() LoadBalancer {
=======
func TestClient_waitForLoadbalancerActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("active"))
	})
	err := client.waitForLoadbalancerActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForLoadbalancerActive returned an error %v", err)
}

func TestClient_waitForLoadbalancerDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(404)
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForLoadbalancerDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForLoadbalancerDeleted returned an error %v", err)
		}
	}
}

func getMockLoadbalancer(status string) LoadBalancer {
>>>>>>> 8d4aa0e... add `context`
	labels := make([]string, 0)
	labels = append(labels, "nice")
	lb := LoadBalancer{
		Properties: LoadBalancerProperties{
			ObjectUUID:          dummyUUID,
			Name:                "go-client-lb",
			Algorithm:           "leastconn",
			LocationUUID:        "45ed677b-3702-4b36-be2a-a2eab9827950",
			ListenIPv6UUID:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			ListenIPv4UUID:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			RedirectHTTPToHTTPS: false,
			ForwardingRules: []ForwardingRule{
				{
					LetsencryptSSL: "",
					ListenPort:     8080,
					Mode:           "http",
					TargetPort:     8000,
				},
			},
			BackendServers: []BackendServer{
				{
					Weight: 100,
					Host:   "185.201.147.176",
				},
			},
			Labels: labels,
		},
	}
	return lb
}

func prepareLoadBalancerHTTPGetResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancer": %s}`, string(res))
}

func prepareLoadBalancerHTTPListResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancers": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareLoadBalancerHTTPCreateResponse() string {
	return fmt.Sprintf(`{"request_uuid": "%s","object_uuid": "%s"}`, dummyRequestUUID, dummyUUID)
}

func prepareLoadBalancerObjectCreateResponse() LoadBalancerCreateResponse {
	return LoadBalancerCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
}

func getMockLoadBalancerEvent() LoadBalancerEvent {
	mock := LoadBalancerEvent{Properties: LoadBalancerEventProperties{
		ObjectType:    "type",
		RequestUUID:   dummyRequestUUID,
		ObjectUUID:    dummyUUID,
		Activity:      "sent",
		RequestType:   "type",
		RequestStatus: "active",
		Change:        "change",
		Timestamp:     dummyTime,
		UserUUID:      dummyUUID,
	}}
	return mock
}

func prepareLoadBalancerEventListHTTPGet() string {
	event := getMockLoadBalancerEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

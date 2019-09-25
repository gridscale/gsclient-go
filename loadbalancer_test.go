package gsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateLoadBalancer(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		lb := getMockLoadbalancer("active").Properties
		labelSlices := [][]string{nil, lb.Labels}
		uri := path.Join(apiLoadBalancerBase)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, http.MethodPost)
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprint(w, prepareLoadBalancerHTTPCreateResponse())
			}
		})

		if clientTest {
			httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
			mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
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
				response, err := client.CreateLoadBalancer(context.Background(), lbRequest)
				if testSuccessFail.isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "CreateLoadBalancer returned error: %v", err)
					assert.Equal(t, fmt.Sprintf("%s", prepareLoadBalancerObjectCreateResponse()), fmt.Sprintf("%s", response))
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient(true)
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	expectedObject := getMockLoadbalancer("active")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("active"))
	})
	for _, test := range uuidCommonTestCases {
		loadbalancer, err := client.GetLoadBalancer(context.Background(), test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancer returned error: %v", err)
			assert.Equal(t, fmt.Sprintf("%v", expectedObject.Properties), fmt.Sprintf("%v", loadbalancer.Properties))
		}
	}
	server.Close()
}

func TestClient_GetLoadBalancerList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase)
	expectedObjects := getMockLoadbalancer("active")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPListResponse("active"))
	})
	loadbalancers, err := client.GetLoadBalancerList(context.Background())
	assert.Nil(t, err, "GetLoadBalancerList returned error: %v", err)
	assert.Equal(t, 1, len(loadbalancers))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", loadbalancers))
}

func TestClient_UpdateLoadBalancer(t *testing.T) {
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
					context.Background(),
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
	}
}

func TestClient_DeleteLoadBalancer(t *testing.T) {
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
				err := client.DeleteLoadBalancer(context.Background(), test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteLoadBalancer returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetLoadBalancerEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID, "events")
	expectedObjects := getMockEvent()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetLoadBalancerEventList(context.Background(), test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancerEventList returned error: %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_waitForLoadbalancerActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("in-provisioning"))
			} else {
				fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			err := client.waitForLoadbalancerActive(context.Background(), dummyUUID)
			if isFailed || isTimeout {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForLoadbalancerActive returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForLoadbalancerDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse("to-be-deleted"))
			} else {
				w.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, test := range uuidCommonTestCases {
				err := client.waitForLoadbalancerDeleted(context.Background(), test.testUUID)
				if test.isFailed || isFailed || isTimeout {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "waitForLoadbalancerDeleted returned an error %v", err)
				}
			}
		}
	}
}

func getMockLoadbalancer(status string) LoadBalancer {
	labels := make([]string, 0)
	labels = append(labels, "nice")
	lb := LoadBalancer{
		Properties: LoadBalancerProperties{
			ObjectUUID:          dummyUUID,
			Name:                "go-client-lb",
			Status:              status,
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

func prepareLoadBalancerHTTPGetResponse(status string) string {
	lb := getMockLoadbalancer(status)
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancer": %s}`, string(res))
}

func prepareLoadBalancerHTTPListResponse(status string) string {
	lb := getMockLoadbalancer(status)
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

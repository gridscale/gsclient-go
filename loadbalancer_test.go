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

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	lb := getMockLoadbalancer().Properties
	lbRequest := LoadBalancerCreateRequest{
		Name:                lb.Name,
		Algorithm:           lb.Algorithm,
		LocationUuid:        lb.LocationUuid,
		ListenIPv6Uuid:      lb.ListenIPv6Uuid,
		ListenIPv4Uuid:      lb.ListenIPv4Uuid,
		RedirectHTTPToHTTPS: lb.RedirectHTTPToHTTPS,
		ForwardingRules:     lb.ForwardingRules,
		BackendServers:      lb.BackendServers,
		Labels:              lb.Labels,
	}
	response, err := client.CreateLoadBalancer(lbRequest)
	if err != nil {
		t.Errorf("CreateLoadBalancer returned error: %v", err)
	}
	assert.Equal(t, fmt.Sprintf("&%s", prepareLoadBalancerObjectCreateResponse()), fmt.Sprintf("%s", response))
}
func TestClient_GetLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUuid)
	expectedObject := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse())
	})
	loadbalancer, err := client.GetLoadBalancer(dummyUuid)
	if err != nil {
		t.Errorf("GetLoadBalancer returned error: %v", err)
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
	loadbalancers, err := client.GetLoadBalancerList()
	if err != nil {
		t.Errorf("GetLoadBalancerList returned error: %v", err)
	}
	assert.Equal(t, 1, len(loadbalancers))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", loadbalancers))
}

func getMockLoadbalancer() LoadBalancer {
	labels := make([]interface{}, 0)
	labels = append(labels, "nice")
	lb := LoadBalancer{
		Properties: LoadBalancerProperties{
			ObjectUuid:          dummyUuid,
			Name:                "go-client-lb",
			Algorithm:           "leastconn",
			LocationUuid:        "45ed677b-3702-4b36-be2a-a2eab9827950",
			ListenIPv6Uuid:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			ListenIPv4Uuid:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			RedirectHTTPToHTTPS: false,
			ForwardingRules: []ForwardingRule{
				{
					LetsencryptSSL: nil,
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
	return fmt.Sprintf(`{"loadbalancers": {"%s": %s}}`, dummyUuid, string(res))
}

func prepareLoadBalancerHTTPCreateResponse() string {
	return fmt.Sprintf(`{"request_uuid": "%s","object_uuid": "%s"}`, dummyRequestUUID, dummyUuid)
}

func prepareLoadBalancerObjectCreateResponse() LoadBalancerCreateResponse {
	return LoadBalancerCreateResponse{
		RequestUuid: dummyRequestUUID,
		ObjectUuid:  dummyUuid,
	}
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetFirewallList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareFirewallListHTTPGet("active"))
	})
	response, err := client.GetFirewallList(emptyCtx)
	assert.Nil(t, err, "GetFirewallList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockFirewall("active")), fmt.Sprintf("%v", response))
}

func TestClient_GetFirewall(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareFirewallHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetFirewall(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetFirewall returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockFirewall("active")), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_CreateFirewall(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase)
	var isFailed bool
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareFirewallCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateFirewall(
			emptyCtx,
			FirewallCreateRequest{
				Name:   "test",
				Labels: []string{"label"},
				Rules: FirewallRules{
					RulesV6In: []FirewallRuleProperties{
						{
							Protocol: TCPTransport,
							DstPort:  "1080",
							SrcPort:  "80",
							Order:    0,
						},
					},
				},
			})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateFirewall returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockFirewallCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_UpdateFirewall(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodPatch {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				fmt.Fprint(w, prepareFirewallHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdateFirewall(
				emptyCtx,
				test.testUUID,
				FirewallUpdateRequest{
					Name:   "test",
					Labels: &[]string{},
					Rules: &FirewallRules{
						RulesV6In: []FirewallRuleProperties{
							{
								Protocol: TCPTransport,
								DstPort:  "1080",
								SrcPort:  "80",
								Order:    0,
							},
						},
					},
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateFirewall returned an error %v", err)
			}
		}
	}

}

func TestClient_DeleteFirewall(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
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
			err := client.DeleteFirewall(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteFirewall returned an error %v", err)
			}
		}
	}
}

func TestClient_GetFirewallEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetFirewallEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetFirewallEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}

	}

}

func getMockFirewall(status string) Firewall {
	mock := Firewall{Properties: FirewallProperties{
		Status:     status,
		Labels:     []string{"label"},
		ObjectUUID: dummyUUID,
		ChangeTime: dummyTime,
		Rules: FirewallRules{
			RulesV6In: []FirewallRuleProperties{
				{
					Protocol: TCPTransport,
					DstPort:  "1080",
					SrcPort:  "80",
					Order:    0,
				},
			},
		},
		CreateTime: dummyTime,
		Private:    true,
		Relations: FirewallRelation{
			Networks: []NetworkInFirewall{
				{
					CreateTime:  dummyTime,
					NetworkUUID: dummyUUID,
					NetworkName: "network",
					ObjectUUID:  dummyUUID,
					ObjectName:  "name",
				},
			},
		},
		Description:  "none",
		LocationName: "Germany",
		Name:         "Test",
	}}
	return mock
}

func prepareFirewallListHTTPGet(status string) string {
	firewall := getMockFirewall(status)
	res, _ := json.Marshal(firewall.Properties)
	return fmt.Sprintf(`{"firewalls": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareFirewallHTTPGet(status string) string {
	firewall := getMockFirewall(status)
	res, _ := json.Marshal(firewall)
	return string(res)
}

func getMockFirewallCreateResponse() FirewallCreateResponse {
	mock := FirewallCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareFirewallCreateResponse() string {
	createRes := getMockFirewallCreateResponse()
	res, _ := json.Marshal(createRes)
	return string(res)
}

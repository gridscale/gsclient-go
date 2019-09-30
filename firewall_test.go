package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetFirewallList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallListHTTPGet())
	})
	response, err := client.GetFirewallList(emptyCtx)
	assert.Nil(t, err, "GetFirewallList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockFirewall()), fmt.Sprintf("%v", response))
}

func TestClient_GetFirewall(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallHTTPGet())
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
	assert.Equal(t, fmt.Sprintf("%v", getMockFirewall()), fmt.Sprintf("%v", response))
}

func TestClient_CreateFirewall(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		uri := path.Join(apiFirewallBase)
		var isFailed bool
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprintf(w, prepareFirewallCreateResponse())
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
			res, err := client.CreateFirewall(
				emptyCtx,
				FirewallCreateRequest{
					Name:   "test",
					Labels: []string{"label"},
					Rules: FirewallRules{
						RulesV6In: []FirewallRuleProperties{
							{
								Protocol: "tcp",
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
		server.Close()
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockFirewallCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateFirewall(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiFirewallBase, dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
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
						Labels: []string{"label"},
						Rules: &FirewallRules{
							RulesV6In: []FirewallRuleProperties{
								{
									Protocol: "tcp",
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
		server.Close()
	}
}

func TestClient_DeleteFirewall(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiFirewallBase, dummyUUID)
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
				err := client.DeleteFirewall(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteFirewall returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetFirewallEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallEventListHTTPGet())
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
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockFirewallEvent()), fmt.Sprintf("%v", response))
}

func TestClient_waitForFirewallActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallHTTPGet("active"))
	})
	err := client.waitForFirewallActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForFirewallActive returned an error %v", err)
}

func TestClient_waitForFirewallDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(404)

	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForFirewallDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForFirewallDeleted returned an error %v", err)
		}
	}
}

func getMockFirewall(status string) Firewall {
	mock := Firewall{Properties: FirewallProperties{
		Status:     "active",
		Labels:     []string{"label"},
		ObjectUUID: dummyUUID,
		ChangeTime: dummyTime,
		Rules: FirewallRules{
			RulesV6In: []FirewallRuleProperties{
				{
					Protocol: "tcp",
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

func prepareFirewallListHTTPGet() string {
	firewall := getMockFirewall()
	res, _ := json.Marshal(firewall.Properties)
	return fmt.Sprintf(`{"firewalls": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareFirewallHTTPGet() string {
	firewall := getMockFirewall()
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

func getMockFirewallEvent() FirewallEvent {
	mock := FirewallEvent{Properties: FirewallEventProperties{
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

func prepareFirewallEventListHTTPGet() string {
	event := getMockFirewallEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

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
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallListHTTPGet())
	})
	response, err := client.GetFirewallList()
	if err != nil {
		t.Errorf("GetFirewallList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockFirewall()), fmt.Sprintf("%v", response))
}

func TestClient_GetFirewall(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallHTTPGet())
	})
	response, err := client.GetFirewall(dummyUuid)
	if err != nil {
		t.Errorf("GetFirewall returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockFirewall()), fmt.Sprintf("%v", response))
}

func TestClient_CreateFirewall(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareFirewallCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateFirewall(FirewallCreateRequest{
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
	if err != nil {
		t.Errorf("CreateFirewall returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockFirewallCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateFirewall(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateFirewall(dummyUuid, FirewallUpdateRequest{
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
	if err != nil {
		t.Errorf("UpdateFirewall returned an error %v", err)
	}
}

func TestClient_DeleteFirewall(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteFirewall(dummyUuid)
	if err != nil {
		t.Errorf("DeleteFirewall returned an error %v", err)
	}
}

func TestClient_GetFirewallEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiFirewallBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareFirewallEventListHTTPGet())
	})
	response, err := client.GetFirewallEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetFirewallEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockFirewallEvent()), fmt.Sprintf("%v", response))
}

func getMockFirewall() Firewall {
	mock := Firewall{Properties: FirewallProperties{
		Status:     "active",
		Labels:     []string{"label"},
		ObjectUuid: dummyUuid,
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
					NetworkUuid: dummyUuid,
					NetworkName: "network",
					ObjectUuid:  dummyUuid,
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
	return fmt.Sprintf(`{"firewalls": {"%s": %s}}`, dummyUuid, string(res))
}

func prepareFirewallHTTPGet() string {
	firewall := getMockFirewall()
	res, _ := json.Marshal(firewall)
	return string(res)
}

func getMockFirewallCreateResponse() FirewallCreateResponse {
	mock := FirewallCreateResponse{
		RequestUuid: dummyRequestUUID,
		ObjectUuid:  dummyUuid,
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
		RequestUuid:   dummyRequestUUID,
		ObjectUuid:    dummyUuid,
		Activity:      "sent",
		RequestType:   "type",
		RequestStatus: "active",
		Change:        "change",
		Timestamp:     dummyTime,
		UserUuid:      dummyUuid,
	}}
	return mock
}

func prepareFirewallEventListHTTPGet() string {
	event := getMockFirewallEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

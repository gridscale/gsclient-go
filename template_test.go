package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetTemplateList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	response, err := client.GetTemplateList()
	if err != nil {
		t.Errorf("GetTemplateList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate()), fmt.Sprintf("%v", response))
}

func TestClient_GetTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateHTTPGet())
	})
	response, err := client.GetTemplate(dummyUUID)
	if err != nil {
		t.Errorf("GetTemplate returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockTemplate()), fmt.Sprintf("%v", response))
}

func TestClient_GetTemplateByName(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	response, err := client.GetTemplateByName("test")
	if err != nil {
		t.Errorf("GetTemplateByName returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockTemplate()), fmt.Sprintf("%v", response))
}

func TestClient_CreateTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareTemplateCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateTemplate(TemplateCreateRequest{
		Name:         "test",
		SnapshotUUID: dummyUUID,
		Labels:       []string{"label"},
	})
	if err != nil {
		t.Errorf("CreateTemplate returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockTemplateCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateTemplate(dummyUUID, TemplateUpdateRequest{
		Name:   "test",
		Labels: []string{"labels"},
	})
	if err != nil {
		t.Errorf("UpdateTemplate returned an error %v", err)
	}
}

func TestClient_DeleteTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteTemplate(dummyUUID)
	if err != nil {
		t.Errorf("DeleteTemplate returned an error %v", err)
	}
}

func TestClient_GetTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateEventListHTTPGet())
	})
	response, err := client.GetTemplateEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetTemplateEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplateEvent()), fmt.Sprintf("%v", response))
}

func getMockTemplate() Template {
	mock := Template{Properties: TemplateProperties{
		Status:           "active",
		Ostype:           "type",
		LocationUUID:     dummyUUID,
		Version:          "1.0",
		LocationIata:     "iata",
		ChangeTime:       dummyTime,
		Private:          true,
		ObjectUUID:       dummyUUID,
		LicenseProductNo: 11111,
		CreateTime:       dummyTime,
		UsageInMinutes:   1000,
		Capacity:         10,
		LocationName:     "Cologne",
		Distro:           "Centos7",
		Description:      "description",
		CurrentPrice:     0,
		LocationCountry:  "Germnany",
		Name:             "test",
		Labels:           []string{"label"},
	}}
	return mock
}

func getMockTemplateCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func getMockTemplateEvent() TemplateEvent {
	mock := TemplateEvent{Properties: TemplateEventProperties{
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

func prepareTemplateListHTTPGet() string {
	template := getMockTemplate()
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"templates": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareTemplateHTTPGet() string {
	template := getMockTemplate()
	res, _ := json.Marshal(template)
	return string(res)
}

func prepareTemplateCreateResponse() string {
	response := getMockTemplateCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareTemplateEventListHTTPGet() string {
	event := getMockTemplateEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

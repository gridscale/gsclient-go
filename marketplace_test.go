package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetMarketplaceTemplateList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiMarketplaceTemplateBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareGetMarketplaceTemplateListHTTPGet())
	})
	res, err := client.GetMarketplaceTemplateList()
	if err != nil {
		t.Errorf("GetMarketplaceTemplateList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockMarketplaceTemplate()), fmt.Sprintf("%v", res))
}

func TestClient_GetMarketplaceTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiMarketplaceTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareGetMarketplaceTemplateHTTPGet())
	})
	res, err := client.GetMarketplaceTemplate(dummyUUID)
	if err != nil {
		t.Errorf("GetMarketplaceTemplate returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceTemplate()), fmt.Sprintf("%v", res))
}

func TestClient_CreateMarketplaceTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiMarketplaceTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareMarketplaceTemplateCreateImportResponseHTTP())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateMarketplaceTemplate(MarketplaceTemplateCreateRequest{
		Name:              "test",
		Labels:            []string{"label"},
		ObjectStoragePath: "s3://example.com",
		Capacity:          1,
	})
	if err != nil {
		t.Errorf("CreateMarketplaceTemplate returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceTemplateCreateImportResponse()), fmt.Sprintf("%v", res))
}

func TestClient_ImportMarketplaceTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiMarketplaceTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareMarketplaceTemplateCreateImportResponseHTTP())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.ImportMarketplaceTemplate(MarketplaceTemplateImportRequest{UniqueHash: "abcd"})
	if err != nil {
		t.Errorf("CreateMarketplaceTemplate returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceTemplateCreateImportResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateMarketplaceTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiMarketplaceTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateMarketplaceTemplate(dummyUUID, MarketplaceTemplateUpdateRequest{
		Name:     "updated name",
		Capacity: 2,
	})
	if err != nil {
		t.Errorf("UpdateMarketplaceTemplate returned an error %v", err)
	}
}

func TestClient_DeleteMarketplaceTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiMarketplaceTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteMarketplaceTemplate(dummyUUID)
	if err != nil {
		t.Errorf("DeleteMarketplaceTemplate returned an error %v", err)
	}
}

func TestClient_GetMarketplaceTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiMarketplaceTemplateBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareMarketplaceTemplateEventListHTTPGet())
	})
	response, err := client.GetMarketplaceTemplateEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetMarketplaceTemplateEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockMarketplaceTemplateEvent()), fmt.Sprintf("%v", response))
}

func getMockMarketplaceTemplate() MarketplaceTemplate {
	mock := MarketplaceTemplate{Properties: MarketplaceTemplateProperties{
		Status:           "active",
		Ostype:           "marketplace",
		LocationUUID:     dummyUUID,
		Version:          "Ubuntu 16.04 LTS",
		LocationIata:     "iata",
		ChangeTime:       dummyTime,
		Private:          false,
		ObjectUUID:       dummyUUID,
		LicenseProductNo: 0,
		CreateTime:       dummyTime,
		UsageInMinutes:   10,
		Capacity:         1,
		LocationName:     "de",
		Distro:           "distro",
		Description:      "des",
		CurrentPrice:     10.99,
		LocationCountry:  "Germany",
		Name:             "Benno",
		Labels:           []string{"label"},
		Metadata: MarketplaceTemplateMetadata{
			OS:   "os",
			Top:  false,
			Icon: "icon",
		},
	}}
	return mock
}

func prepareGetMarketplaceTemplateListHTTPGet() string {
	template := getMockMarketplaceTemplate()
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"templates": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareGetMarketplaceTemplateHTTPGet() string {
	template := getMockMarketplaceTemplate()
	res, _ := json.Marshal(template)
	return string(res)
}

func getMockMarketplaceTemplateCreateImportResponse() MarketplaceTemplateCreateImportResponse {
	mock := MarketplaceTemplateCreateImportResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
		UniqueHash:  "abcd",
	}
	return mock
}

func prepareMarketplaceTemplateCreateImportResponseHTTP() string {
	response := getMockMarketplaceTemplateCreateImportResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func getMockMarketplaceTemplateEvent() MarketplaceTemplateEvent {
	mock := MarketplaceTemplateEvent{Properties: MarketplaceTemplateEventProperties{
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

func prepareMarketplaceTemplateEventListHTTPGet() string {
	event := getMockMarketplaceTemplateEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

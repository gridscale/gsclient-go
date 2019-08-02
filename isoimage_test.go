package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetISOImageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareISOImageHTTPGetList())
	})
	res, err := client.GetISOImageList()
	if err != nil {
		t.Errorf("GetISOImageList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImage()), fmt.Sprintf("%v", res))
}

func TestClient_GetISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareISOImageHTTPGet())
	})
	res, err := client.GetISOImage(dummyUuid)
	if err != nil {
		t.Errorf("GetISOImage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockISOImage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprintf(writer, prepareISOImageHTTPCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateISOImage(ISOImageCreateRequest{
		Name:         "Test",
		SourceUrl:    "http://example.org",
		Labels:       []string{"label"},
		LocationUuid: "aa-bb-cc",
	})
	if err != nil {
		t.Errorf("CreateISOImage returned an error %v", err)
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockISOImageCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateISOImage(dummyUuid, ISOImageUpdateRequest{
		Name:   "test",
		Labels: []string{},
	})
	if err != nil {
		t.Errorf("UpdateISOImage returned an error %v", err)
	}
}

func TestClient_DeleteISOImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteISOImage(dummyUuid)
	if err != nil {
		t.Errorf("DeleteISOImage returned an error %v", err)
	}
}

func TestClient_GetISOImageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiISOBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprint(writer, prepareISOImageHTTPGetEventList())
	})

	res, err := client.GetISOImageEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetISOImageEventList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockISOImageEvent()), fmt.Sprintf("%v", res))
}

func getMockISOImage() ISOImage {
	mock := ISOImage{Properties: ISOImageProperties{
		ObjectUuid: dummyUuid,
		Relations: ISOImageRelation{
			Servers: []ServerinISOImage{
				{
					Bootdevice: true,
					CreateTime: dummyTime,
					ObjectName: "test",
					ObjectUuid: "abc-def",
				},
			},
		},
		Description:     "description",
		LocationName:    "locationName",
		SourceUrl:       "url",
		Labels:          []string{"label"},
		LocationIata:    "iata",
		LocationUuid:    "locUuid",
		Status:          "active",
		CreateTime:      dummyTime,
		Name:            "test",
		Version:         "1.0",
		LocationCountry: "Country",
		UsageInMinutes:  10,
		Private:         false,
		ChangeTime:      dummyTime,
		Capacity:        10,
		CurrentPrice:    9.99,
	}}
	return mock
}

func prepareISOImageHTTPGetList() string {
	iso := getMockISOImage()
	res, _ := json.Marshal(iso.Properties)
	return fmt.Sprintf(`{"isoimages": {"%s": %s}}`, dummyUuid, string(res))
}

func prepareISOImageHTTPGet() string {
	iso := getMockISOImage()
	res, _ := json.Marshal(iso)
	return string(res)
}

func getMockISOImageCreateResponse() ISOImageCreateResponse {
	mock := ISOImageCreateResponse{
		RequestUuid: dummyRequestUUID,
		ObjectUuid:  dummyUuid,
	}
	return mock
}

func prepareISOImageHTTPCreateResponse() string {
	res, _ := json.Marshal(getMockISOImageCreateResponse())
	return string(res)
}

func getMockISOImageEvent() ISOImageEvent {
	mock := ISOImageEvent{Properties: ISOImageEventProperties{
		ObjectType:    "type",
		RequestUuid:   dummyRequestUUID,
		ObjectUuid:    dummyUuid,
		Activity:      "activity",
		RequestType:   "request type",
		RequestStatus: "active",
		Change:        "change description",
		Timestamp:     dummyTime,
		UserUuid:      "user-id",
	}}
	return mock
}

func prepareISOImageHTTPGetEventList() string {
	res, _ := json.Marshal(getMockISOImageEvent().Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetTemplateList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	response, err := client.GetTemplateList(emptyCtx)
	assert.Nil(t, err, "GetTemplateList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate("active")), fmt.Sprintf("%v", response))
}

func TestClient_GetTemplate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareTemplateHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplate(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}

	}
}

func TestClient_GetTemplateByName(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	testCases := []uuidTestCase{
		{
			testUUID: "test",
			isFailed: false,
		},
		{
			testUUID: "",
			isFailed: true,
		},
	}
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	for _, test := range testCases {
		response, err := client.GetTemplateByName(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplateByName returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_CreateTemplate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareTemplateCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateTemplate(
			emptyCtx,
			TemplateCreateRequest{
				Name:         "test",
				SnapshotUUID: dummyUUID,
				Labels:       []string{"label"},
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateTemplate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplateCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_UpdateTemplate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodPatch {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				fmt.Fprint(w, prepareTemplateHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdateTemplate(
				emptyCtx,
				test.testUUID,
				TemplateUpdateRequest{
					Name:   "test",
					Labels: &[]string{"labels"},
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateTemplate returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteTemplate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
			err := client.DeleteTemplate(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteTemplate returned an error %v", err)
			}
		}
	}
}

func TestClient_GetTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplateEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplateEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetTemplatesByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "templates")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplatesByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplatesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetDeletedTemplates(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "templates")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareDeletedTemplateListHTTPGet())
	})
	response, err := client.GetDeletedTemplates(emptyCtx)
	assert.Nil(t, err, "GetDeletedTemplates returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate("deleted")), fmt.Sprintf("%v", response))
}

func getMockTemplate(status string) Template {
	mock := Template{Properties: TemplateProperties{
		Status:           status,
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
		LocationCountry:  "Germany",
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

func prepareTemplateListHTTPGet() string {
	template := getMockTemplate("active")
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"templates": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareTemplateHTTPGet(status string) string {
	template := getMockTemplate(status)
	res, _ := json.Marshal(template)
	return string(res)
}

func prepareTemplateCreateResponse() string {
	response := getMockTemplateCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareDeletedTemplateListHTTPGet() string {
	template := getMockTemplate("deleted")
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"deleted_templates": {"%s": %s}}`, dummyUUID, string(res))
}

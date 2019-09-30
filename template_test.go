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
	response, err := client.GetTemplateList(emptyCtx)
	assert.Nil(t, err, "GetTemplateList returned an error %v", err)
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
<<<<<<< HEAD
	response, err := client.GetTemplate(dummyUUID)
	if err != nil {
		t.Errorf("GetTemplate returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplate(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}

>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	response, err := client.GetTemplateByName("test")
	if err != nil {
		t.Errorf("GetTemplateByName returned an error %v", err)
=======
	for _, test := range testCases {
		response, err := client.GetTemplateByName(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplateByName returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockTemplate()), fmt.Sprintf("%v", response))
}

func TestClient_CreateTemplate(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiTemplateBase
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprintf(w, prepareTemplateCreateResponse())
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
		server.Close()
	}
}

func TestClient_UpdateTemplate(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiTemplateBase, dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
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
						Labels: []string{"labels"},
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateTemplate returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteTemplate(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiTemplateBase, dummyUUID)
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
				err := client.DeleteTemplate(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteTemplate returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
>>>>>>> 8d4aa0e... add `context`
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareTemplateCreateResponse())
	})
<<<<<<< HEAD
=======
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
>>>>>>> 8d4aa0e... add `context`

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
<<<<<<< HEAD

	res, err := client.CreateTemplate(TemplateCreateRequest{
		Name:         "test",
		SnapshotUUID: dummyUUID,
		Labels:       []string{"label"},
	})
	if err != nil {
		t.Errorf("CreateTemplate returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplatesByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplatesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate("active")), fmt.Sprintf("%v", response))
		}
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	err := client.UpdateTemplate(dummyUUID, TemplateUpdateRequest{
		Name:   "test",
		Labels: []string{"labels"},
	})
	if err != nil {
		t.Errorf("UpdateTemplate returned an error %v", err)
	}
=======
	response, err := client.GetDeletedTemplates(emptyCtx)
	assert.Nil(t, err, "GetDeletedTemplates returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate("deleted")), fmt.Sprintf("%v", response))
>>>>>>> 8d4aa0e... add `context`
}

func TestClient_DeleteTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
<<<<<<< HEAD
	err := client.DeleteTemplate(dummyUUID)
	if err != nil {
		t.Errorf("DeleteTemplate returned an error %v", err)
	}
=======
	err := client.waitForTemplateActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForTemplateActive returned an error %v", err)
>>>>>>> 8d4aa0e... add `context`
}

func TestClient_GetTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateEventListHTTPGet())
	})
<<<<<<< HEAD
	response, err := client.GetTemplateEventList(dummyUUID)
	if err != nil {
		t.Errorf("GetTemplateEventList returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		err := client.waitForTemplateDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForTemplateDeleted returned an error %v", err)
		}
>>>>>>> 8d4aa0e... add `context`
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

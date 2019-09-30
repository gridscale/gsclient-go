package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetStorageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	response, err := client.GetStorageList(emptyCtx)
	assert.Nil(t, err, "GetStorageList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_GetStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStorage(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorage returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockStorage("active")), fmt.Sprintf("%v", response))
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_CreateStorage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := apiStorageBase
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprintf(w, prepareStorageCreateResponse())
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
			res, err := client.CreateStorage(
				emptyCtx,
				StorageCreateRequest{
					Capacity:     10,
					LocationUUID: dummyUUID,
					Name:         "test",
					StorageType:  DefaultStorageType,
					Template: &StorageTemplate{
						TemplateUUID: dummyUUID,
						Password:     "pass",
						PasswordType: CryptPasswordType,
						Hostname:     "example.com",
					},
					Labels: []string{"label"},
				})
			if isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateStorage returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
			}
		}
		server.Close()
	}
}

func TestClient_UpdateStorage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			if isFailed {
				w.WriteHeader(400)
			} else {
				if r.Method == http.MethodPatch {
					fmt.Fprintf(w, "")
				} else if r.Method == http.MethodGet {
					fmt.Fprint(w, prepareStorageHTTPGet("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range uuidCommonTestCases {
				err := client.UpdateStorage(
					emptyCtx,
					test.testUUID,
					StorageUpdateRequest{
						Name:     "test",
						Labels:   []string{"label"},
						Capacity: 20,
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateStorage returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteStorage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID)
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
				err := client.DeleteStorage(emptyCtx, test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteStorage returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareFirewallCreateResponse())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStorageEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}
	}
}

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateStorage(StorageCreateRequest{
		Capacity:     10,
		LocationUUID: dummyUUID,
		Name:         "test",
		StorageType:  "storage",
		Template: &StorageTemplate{
			TemplateUUID: dummyUUID,
			Password:     "pass",
			PasswordType: "crypt",
			Hostname:     "example.com",
		},
		Labels: []string{"label"},
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStoragesByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStoragesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage("active")), fmt.Sprintf("%v", response))
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	response, err := client.GetDeletedStorages(emptyCtx)
	assert.Nil(t, err, "GetDeletedStorages returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage("deleted")), fmt.Sprintf("%v", response))
}

func TestClient_DeleteStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.waitForStorageActive(emptyCtx, dummyUUID)
	assert.Nil(t, err, "waitForStorageActive returned an error %v", err)
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		err := client.waitForStorageDeleted(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForStorageDeleted returned an error %v", err)
		}
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageEvent()), fmt.Sprintf("%v", response))
}

func getMockStorage() Storage {
	mock := Storage{Properties: StorageProperties{
		ChangeTime:       dummyTime,
		LocationIata:     "iata",
		Status:           "active",
		LicenseProductNo: 11111,
		LocationCountry:  "Germany",
		UsageInMinutes:   10,
		LastUsedTemplate: dummyUUID,
		CurrentPrice:     9.1,
		Capacity:         10,
		LocationUUID:     dummyUUID,
		StorageType:      "storage",
		ParentUUID:       dummyUUID,
		Name:             "test",
		LocationName:     "Cologne",
		ObjectUUID:       dummyUUID,
		Snapshots: []StorageSnapshotRelation{
			{
				LastUsedTemplate:      dummyUUID,
				ObjectUUID:            dummyUUID,
				StorageUUID:           dummyUUID,
				SchedulesSnapshotName: "test",
				SchedulesSnapshotUUID: dummyUUID,
				ObjectCapacity:        10,
				CreateTime:            dummyTime,
				ObjectName:            "test",
			},
		},
		Relations:  StorageRelations{},
		Labels:     []string{"label"},
		CreateTime: dummyTime,
	}}
	return mock
}

func getMockStorageCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func getMockStorageEvent() StorageEvent {
	mock := StorageEvent{Properties: StorageEventProperties{
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

func prepareStorageListHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"storages": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareStorageHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage)
	return string(res)
}

func prepareStorageCreateResponse() string {
	response := getMockStorageCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareStorageEventListHTTPGet() string {
	event := getMockStorageEvent()
	res, _ := json.Marshal(event.Properties)
	return fmt.Sprintf(`{"events": [%s]}`, string(res))
}

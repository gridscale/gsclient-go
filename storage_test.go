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
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	response, err := client.GetStorageList(emptyCtx)
	assert.Nil(t, err, "GetStorageList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage("active")), fmt.Sprintf("%v", response))
}

func TestClient_GetStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareStorageHTTPGet("active"))
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
}

func TestClient_CreateStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareStorageCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateStorage(
			emptyCtx,
			StorageCreateRequest{
				Capacity:    10,
				Name:        "test",
				StorageType: DefaultStorageType,
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
}

func TestClient_CreateStorageFromBackup(t *testing.T) {
	server, client, muxServer := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, "import")
	muxServer.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareStorageCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateStorageFromBackup(
			emptyCtx,
			dummyUUID,
			"some-storage-name",
		)
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateStorageFromBackup returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_UpdateStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
					Labels:   &[]string{"label"},
					Capacity: 20,
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateStorage returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID)
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
			err := client.DeleteStorage(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteStorage returned an error %v", err)
			}
		}
	}
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareEventListHTTPGet())
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

func TestClient_GetStoragesByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareStorageListHTTPGet())
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
}

func TestClient_GetDeletedStorages(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "storages")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(w, prepareDeletedStorageListHTTPGet())
	})
	response, err := client.GetDeletedStorages(emptyCtx)
	assert.Nil(t, err, "GetDeletedStorages returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage("deleted")), fmt.Sprintf("%v", response))
}

func TestClient_CloneStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "clone")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareStorageCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CloneStorage(emptyCtx, dummyUUID)
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateStorage returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func getMockStorage(status string) Storage {
	mock := Storage{Properties: StorageProperties{
		ChangeTime:       dummyTime,
		LocationIata:     "iata",
		Status:           status,
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

func prepareStorageListHTTPGet() string {
	storage := getMockStorage("active")
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"storages": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareStorageHTTPGet(status string) string {
	storage := getMockStorage(status)
	res, _ := json.Marshal(storage)
	return string(res)
}

func prepareStorageCreateResponse() string {
	response := getMockStorageCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareDeletedStorageListHTTPGet() string {
	storage := getMockStorage("deleted")
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"deleted_storages": {"%s": %s}}`, dummyUUID, string(res))
}

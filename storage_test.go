package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	response, err := client.GetStorageList()
	if err != nil {
		t.Errorf("GetStorageList returned an error %v", err)
	}
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_GetStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageHTTPGet())
	})
	response, err := client.GetStorage(dummyUuid)
	if err != nil {
		t.Errorf("GetStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_CreateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareFirewallCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateStorage(StorageCreateRequest{
		Capacity:     10,
		LocationUuid: dummyUuid,
		Name:         "test",
		StorageType:  "storage",
		Template: &StorageTemplate{
			TemplateUuid: dummyUuid,
			Password:     "pass",
			PasswordType: "crypt",
			Hostname:     "example.com",
		},
		Labels: []string{"label"},
	})
	if err != nil {
		t.Errorf("CreateStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.UpdateStorage(dummyUuid, StorageUpdateRequest{
		Name:     "test",
		Labels:   []string{"label"},
		Capacity: 20,
	})
	if err != nil {
		t.Errorf("UpdateStorage returned an error %v", err)
	}
}

func TestClient_DeleteStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	err := client.DeleteStorage(dummyUuid)
	if err != nil {
		t.Errorf("DeleteStorage returned an error %v", err)
	}
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageEventListHTTPGet())
	})
	response, err := client.GetStorageEventList(dummyUuid)
	if err != nil {
		t.Errorf("GetStorageEventList returned an error %v", err)
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
		LastUsedTemplate: dummyUuid,
		CurrentPrice:     9.1,
		Capacity:         10,
		LocationUuid:     dummyUuid,
		StorageType:      "storage",
		ParentUuid:       dummyUuid,
		Name:             "test",
		LocationName:     "Cologne",
		ObjectUuid:       dummyUuid,
		Snapshots: []StorageSnapshotRelation{
			{
				LastUsedTemplate:      dummyUuid,
				ObjectUuid:            dummyUuid,
				StorageUuid:           dummyUuid,
				SchedulesSnapshotName: "test",
				SchedulesSnapshotUuid: dummyUuid,
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
		ObjectUuid:  dummyUuid,
		RequestUuid: dummyRequestUUID,
	}
	return mock
}

func getMockStorageEvent() StorageEvent {
	mock := StorageEvent{Properties: StorageEventProperties{
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

func prepareStorageListHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"storages": {"%s": %s}}`, dummyUuid, string(res))
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

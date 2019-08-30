package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageSnapshotList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotListHTTPGet())
	})

	res, err := client.GetStorageSnapshotList(dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshotList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot()), fmt.Sprintf("%v", res))
}

func TestClient_GetStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotHTTPGet())
	})

	res, err := client.GetStorageSnapshot(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshot returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshot()), fmt.Sprintf("%v", res))
}

func TestClient_CreateStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprint(w, prepareStorageSnapshotCreateResponseHTTP())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateStorageSnapshot(dummyUUID, StorageSnapshotCreateRequest{
		Name:   "test",
		Labels: []string{"label"},
	})
	if err != nil {
		t.Errorf("CreateStorageSnapshot returned an error: %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotCreateResponse()), fmt.Sprintf("%v", response))
}

func TestClient_UpdateStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.UpdateStorageSnapshot(dummyUUID, dummyUUID, StorageSnapshotUpdateRequest{
		Name:   "test",
		Labels: []string{"label"},
	})
	if err != nil {
		t.Errorf("UpdateStorageSnapshot returned an error %v", err)
	}
}

func TestClient_DeleteStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.DeleteStorageSnapshot(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteStorageSnapshot returned an error %v", err)
	}
}

func TestClient_RollbackStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "rollback")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.RollbackStorage(dummyUUID, dummyUUID, StorageRollbackRequest{Rollback: true})
	if err != nil {
		t.Errorf("RollbackStorage returned an error %v", err)
	}
}

func TestClient_ExportStorageSnapshotToS3(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "export_to_s3")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.ExportStorageSnapshotToS3(dummyUUID, dummyUUID, StorageSnapshotExportToS3Request{
		S3auth: struct {
			Host      string `json:"host"`
			AccessKey string `json:"access_key"`
			SecretKey string `json:"secret_key"`
		}{
			Host:      "example.com",
			AccessKey: "access_key",
			SecretKey: "secret_key",
		},
		S3data: struct {
			Host     string `json:"host"`
			Bucket   string `json:"bucket"`
			Filename string `json:"filename"`
			Private  bool   `json:"private"`
		}{
			Host:     "example.com",
			Bucket:   "bucket",
			Filename: "filename",
			Private:  true,
		},
	})
	if err != nil {
		t.Errorf("ExportStorageSnapshotToS3 returned an error %v", err)
	}
}

func TestClient_GetSnapshotsByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotListHTTPGet())
	})

	res, err := client.GetSnapshotsByLocation(dummyUUID)
	if err != nil {
		t.Errorf("GetSnapshotsByLocation returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot()), fmt.Sprintf("%v", res))
}

func TestClient_GetDeletedSnapshots(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedStorageSnapshotListHTTPGet())
	})

	res, err := client.GetDeletedSnapshots()
	if err != nil {
		t.Errorf("GetSnapshotsByLocation returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot()), fmt.Sprintf("%v", res))
}

func getMockStorageSnapshot() StorageSnapshot {
	mock := StorageSnapshot{Properties: StorageSnapshotProperties{
		Labels:           []string{"label"},
		ObjectUUID:       dummyUUID,
		Name:             "test",
		Status:           "active",
		LocationCountry:  "Germany",
		UsageInMinutes:   60,
		LocationUUID:     dummyUUID,
		ChangeTime:       dummyTime,
		LicenseProductNo: 20,
		CurrentPrice:     0.5,
		CreateTime:       dummyTime,
		Capacity:         10,
		LocationName:     "Cologne",
		LocationIata:     "",
		ParentUUID:       dummyUUID,
	}}
	return mock
}

func prepareStorageSnapshotHTTPGet() string {
	snapshot := getMockStorageSnapshot()
	res, _ := json.Marshal(snapshot)
	return string(res)
}

func prepareStorageSnapshotListHTTPGet() string {
	snapshot := getMockStorageSnapshot()
	res, _ := json.Marshal(snapshot.Properties)
	return fmt.Sprintf(`{"snapshots" : {"%s" : %s}}`, dummyUUID, string(res))
}

func getMockStorageSnapshotCreateResponse() StorageSnapshotCreateResponse {
	mock := StorageSnapshotCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareStorageSnapshotCreateResponseHTTP() string {
	createRes := getMockStorageSnapshotCreateResponse()
	res, _ := json.Marshal(createRes)
	return string(res)
}

func prepareDeletedStorageSnapshotListHTTPGet() string {
	snapshot := getMockStorageSnapshot()
	res, _ := json.Marshal(snapshot.Properties)
	return fmt.Sprintf(`{"deleted_snapshots" : {"%s" : %s}}`, dummyUUID, string(res))
}

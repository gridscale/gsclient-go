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
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotListHTTPGet())
	})

	res, err := client.GetStorageSnapshotList(dummyUuid)
	if err != nil {
		t.Errorf("GetStorageSnapshotList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshotSingle()), fmt.Sprintf("%v", res))
}

func TestClient_GetStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots", dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotHTTPGet())
	})

	res, err := client.GetStorageSnapshot(dummyUuid, dummyUuid)
	if err != nil {
		t.Errorf("GetStorageSnapshot returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotSingle()), fmt.Sprintf("%v", res))
}

func TestClient_CreateStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprint(w, prepareStorageSnapshotCreateResponseHTTP())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	response, err := client.CreateStorageSnapshot(dummyUuid, StorageSnapshotCreateRequest{
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
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots", dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.UpdateStorageSnapshot(dummyUuid, dummyUuid, StorageSnapshotUpdateRequest{
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
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots", dummyUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.DeleteStorageSnapshot(dummyUuid, dummyUuid)
	if err != nil {
		t.Errorf("DeleteStorageSnapshot returned an error %v", err)
	}
}

func TestClient_RollbackStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots", dummyUuid, "rollback")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.RollbackStorage(dummyUuid, dummyUuid, StorageRollbackRequest{Rollback:true})
	if err != nil {
		t.Errorf("RollbackStorage returned an error %v", err)
	}
}

func TestClient_ExportStorageSnapshotToS3(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUuid, "snapshots", dummyUuid, "export_to_s3")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
	err := client.ExportStorageSnapshotToS3(dummyUuid, dummyUuid, StorageSnapshotExportToS3Request{
		S3auth: struct {
			Host       string `json:"host"`
			AccessKeys string `json:"access_keys"`
			SecretKey  string `json:"secret_key"`
		}{
			Host: "example.com",
			AccessKeys: "access_keys",
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

func getMockStorageSnapshotSingle() StorageSnapshotSingle {
	mock := StorageSnapshotSingle{Properties:StorageSnapshotProperties{
		Labels:           []string{"label"},
		ObjectUuid:       dummyUuid,
		Name:             "test",
		Status:           "active",
		LocationCountry:  "Germany",
		UsageInMinutes:   60,
		LocationUuid:     dummyUuid,
		ChangeTime:       dummyTime,
		LicenseProductNo: 20,
		CurrentPrice:     0.5,
		CreateTime:       dummyTime,
		Capacity:         10,
		LocationName:     "Cologne",
		LocationIata:     "",
		ParentUuid:       dummyUuid,
	}}
	return mock
}

func prepareStorageSnapshotHTTPGet() string {
	snapshot := getMockStorageSnapshotSingle()
	res, _ := json.Marshal(snapshot)
	return string(res)
}

func prepareStorageSnapshotListHTTPGet() string {
	snapshot := getMockStorageSnapshotSingle()
	res, _ := json.Marshal(snapshot.Properties)
	return fmt.Sprintf(`{"snapshots" : {"%s" : %s}}`, dummyUuid, string(res))
}

func getMockStorageSnapshotCreateResponse() StorageSnapshotCreateResponse {
	mock := StorageSnapshotCreateResponse{
		RequestUuid: dummyRequestUUID,
		ObjectUuid:  dummyUuid,
	}
	return mock
}

func prepareStorageSnapshotCreateResponseHTTP() string {
	createRes := getMockStorageSnapshotCreateResponse()
	res, _ := json.Marshal(createRes)
	return string(res)
}
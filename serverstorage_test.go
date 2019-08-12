package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerStorageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageListHTTPGet())
	})
	res, err := client.GetServerStorageList(dummyUuid)
	if err != nil {
		t.Errorf("GetServerStorageList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerStorage()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages", dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	res, err := client.GetServerStorage(dummyUuid, dummyUuid)
	if err != nil {
		t.Errorf("GetServerStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerStorage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerStorage(dummyUuid, ServerStorageCreateRequest{
		ObjectUuid: dummyUuid,
		BootDevice: true,
	})
	if err != nil {
		t.Errorf("CreateServerStorage returned an error %v", err)
	}
}

func TestClient_UpdateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages", dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UpdateServerStorage(dummyUuid, dummyUuid, ServerStorageUpdateRequest{
		Ordering:   1,
		BootDevice: true,
	})
	if err != nil {
		t.Errorf("UpdateServerStorage returned an error %v", err)
	}
}

func TestClient_DeleteServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages", dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerStorage(dummyUuid, dummyUuid)
	if err != nil {
		t.Errorf("DeleteServerStorage returned an error %v", err)
	}
}

func TestClient_LinkStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkStorage(dummyUuid, dummyUuid, true)
	if err != nil {
		t.Errorf("LinkStorage returned an error %v", err)
	}
}

func TestClient_UnlinkStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUuid, "storages", dummyUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkStorage(dummyUuid, dummyUuid)
	if err != nil {
		t.Errorf("UnlinkStorage returned an error %v", err)
	}
}

func getMockServerStorage() ServerStorage {
	mock := ServerStorage{
		ObjectUuid:       dummyUuid,
		ObjectName:       "test",
		Capacity:         10,
		StorageType:      "SSD",
		Target:           1,
		Lun:              2,
		Controller:       3,
		CreateTime:       dummyTime,
		BootDevice:       false,
		Bus:              1,
		LastUsedTemplate: dummyUuid,
		LicenseProductNo: 123456789,
		ServerUuid:       dummyUuid,
	}
	return mock
}

func prepareServerStorageListHTTPGet() string {
	serverStorage := getMockServerStorage()
	res, _ := json.Marshal(serverStorage)
	return fmt.Sprintf(`{"storage_relations": [%s]}`, string(res))
}

func prepareServerStorageHTTPGet() string {
	serverStorage := getMockServerStorage()
	res, _ := json.Marshal(serverStorage)
	return fmt.Sprintf(`{"storage_relation": %s}`, string(res))
}
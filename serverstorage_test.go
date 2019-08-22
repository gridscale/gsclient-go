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
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageListHTTPGet())
	})
	res, err := client.GetServerStorageList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerStorageList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerStorage()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	res, err := client.GetServerStorage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetServerStorage returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerStorage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerStorage(dummyUUID, ServerStorageRelationCreateRequest{
		ObjectUUID: dummyUUID,
		BootDevice: true,
	})
	if err != nil {
		t.Errorf("CreateServerStorage returned an error %v", err)
	}
}

func TestClient_UpdateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UpdateServerStorage(dummyUUID, dummyUUID, ServerStorageRelationUpdateRequest{
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
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerStorage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerStorage returned an error %v", err)
	}
}

func TestClient_LinkStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkStorage(dummyUUID, dummyUUID, true)
	if err != nil {
		t.Errorf("LinkStorage returned an error %v", err)
	}
}

func TestClient_UnlinkStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkStorage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("UnlinkStorage returned an error %v", err)
	}
}

func getMockServerStorage() ServerStorageRelationProperties {
	mock := ServerStorageRelationProperties{
		ObjectUUID:       dummyUUID,
		ObjectName:       "test",
		Capacity:         10,
		StorageType:      "SSD",
		Target:           1,
		Lun:              2,
		Controller:       3,
		CreateTime:       dummyTime,
		BootDevice:       false,
		Bus:              1,
		LastUsedTemplate: dummyUUID,
		LicenseProductNo: 123456789,
		ServerUUID:       dummyUUID,
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

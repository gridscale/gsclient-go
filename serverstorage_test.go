package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServerStorageList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareServerStorageListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerStorageList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerStorageList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerStorage()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetServerStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testStorageID := range uuidCommonTestCases {
			res, err := client.GetServerStorage(emptyCtx, testServerID.testUUID, testStorageID.testUUID)
			if testServerID.isFailed || testStorageID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetServerStorage returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerStorage()), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_CreateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, "")
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		for _, testServerID := range uuidCommonTestCases {
			for _, testStorageID := range uuidCommonTestCases {
				err := client.CreateServerStorage(
					emptyCtx,
					testServerID.testUUID,
					ServerStorageRelationCreateRequest{
						ObjectUUID: testStorageID.testUUID,
						BootDevice: true,
					})
				if testServerID.isFailed || testStorageID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "CreateServerStorage returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_UpdateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testStorageID := range uuidCommonTestCases {
			err := client.UpdateServerStorage(
				emptyCtx,
				testServerID.testUUID,
				testStorageID.testUUID,
				ServerStorageRelationUpdateRequest{
					Ordering:   1,
					BootDevice: true,
				})
			if testServerID.isFailed || testStorageID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateServerStorage returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteServerStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodDelete {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				writer.WriteHeader(404)
			}
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		for _, testServerID := range uuidCommonTestCases {
			for _, testStorageID := range uuidCommonTestCases {
				err := client.DeleteServerStorage(emptyCtx, testServerID.testUUID, testStorageID.testUUID)
				if testServerID.isFailed || testStorageID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteServerStorage returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_LinkStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprint(writer, "")
	})
	err := client.LinkStorage(emptyCtx, dummyUUID, dummyUUID, true)
	assert.Nil(t, err, "CreateServerStorage returned an error %v", err)

}

func TestClient_UnlinkStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if request.Method == http.MethodDelete {
			fmt.Fprintf(writer, "")
		} else if request.Method == http.MethodGet {
			writer.WriteHeader(404)
		}
	})
	err := client.UnlinkStorage(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkStorage returned an error %v", err)
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

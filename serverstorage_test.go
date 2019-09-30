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
	assert.Equal(t, fmt.Sprintf("%v", getMockServerStorage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerStorage(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "storages")
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprint(writer, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "storages", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, http.MethodGet, request.Method)
				fmt.Fprintf(writer, prepareServerStorageHTTPGet())
			})
		}
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
		server.Close()
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
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
		server.Close()
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	err := client.LinkStorage(dummyUUID, dummyUUID, true)
	if err != nil {
		t.Errorf("LinkStorage returned an error %v", err)
	}
=======
	mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "storages", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	err := client.LinkStorage(emptyCtx, dummyUUID, dummyUUID, true)
	assert.Nil(t, err, "CreateServerStorage returned an error %v", err)

>>>>>>> 8d4aa0e... add `context`
}

func TestClient_UnlinkStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkStorage(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkStorage returned an error %v", err)
}

<<<<<<< HEAD
=======
func TestClient_waitForServerStorageRelCreation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerStorageRelCreation(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerStorageRelCreation returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForServerStorageRelDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(404)
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerStorageRelDeleted(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerStorageRelDeleted returned an error %v", err)
			}
		}
	}
}

>>>>>>> 8d4aa0e... add `context`
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

package gsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerStorageList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerStorageList(context.Background(), test.testUUID)
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
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testStorageID := range uuidCommonTestCases {
			res, err := client.GetServerStorage(context.Background(), testServerID.testUUID, testStorageID.testUUID)
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
						context.Background(),
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
	}
}

func TestClient_UpdateServerStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testStorageID := range uuidCommonTestCases {
			err := client.UpdateServerStorage(
				context.Background(),
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
					err := client.DeleteServerStorage(context.Background(), testServerID.testUUID, testStorageID.testUUID)
					if testServerID.isFailed || testStorageID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "DeleteServerStorage returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_LinkStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "storages", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerStorageHTTPGet())
	})
	err := client.LinkStorage(context.Background(), dummyUUID, dummyUUID, true)
	assert.Nil(t, err, "CreateServerStorage returned an error %v", err)

}

func TestClient_UnlinkStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {

		if request.Method == http.MethodDelete {
			fmt.Fprintf(writer, "")
		} else if request.Method == http.MethodGet {
			writer.WriteHeader(404)
		}
	})
	err := client.UnlinkStorage(context.Background(), dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkStorage returned an error %v", err)
}

func TestClient_waitForServerStorageRelCreation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				writer.WriteHeader(404)
			} else {
				fmt.Fprintf(writer, prepareServerIsoImageHTTPGet())
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, testServerID := range uuidCommonTestCases {
				for _, testIPID := range uuidCommonTestCases {
					err := client.waitForServerStorageRelCreation(context.Background(), testServerID.testUUID, testIPID.testUUID)
					if testServerID.isFailed || testIPID.isFailed || isFailed || isTimeout {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "waitForServerStorageRelCreation returned an error %v", err)
					}
				}
			}
		}
	}
}

func TestClient_waitForServerStorageRelDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiServerBase, dummyUUID, "storages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprintf(writer, prepareServerStorageHTTPGet())
			} else {
				writer.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, testServerID := range uuidCommonTestCases {
				for _, testIPID := range uuidCommonTestCases {
					err := client.waitForServerStorageRelDeleted(context.Background(), testServerID.testUUID, testIPID.testUUID)
					if testServerID.isFailed || testIPID.isFailed || isFailed || isTimeout {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "waitForServerStorageRelDeleted returned an error %v", err)
					}
				}
			}
		}
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

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
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageSnapshotListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageSnapshotList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageSnapshotList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageSnapshotHTTPGet("active"))
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testSnapshotID := range uuidCommonTestCases {
			res, err := client.GetStorageSnapshot(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID)
			if testStorageID.isFailed || testSnapshotID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetStorageSnapshot returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshot("active")), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_CreateStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprint(w, prepareStorageSnapshotCreateResponseHTTP())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		for _, test := range uuidCommonTestCases {
			response, err := client.CreateStorageSnapshot(
				emptyCtx,
				test.testUUID,
				StorageSnapshotCreateRequest{
					Name:   "test",
					Labels: []string{"label"},
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateStorageSnapshot returned an error: %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotCreateResponse()), fmt.Sprintf("%v", response))
			}
		}
	}

}

func TestClient_UpdateStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodPatch {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				fmt.Fprint(w, prepareStorageSnapshotHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, testStorageID := range uuidCommonTestCases {
			for _, testSnapshotID := range uuidCommonTestCases {
				err := client.UpdateStorageSnapshot(
					emptyCtx,
					testStorageID.testUUID,
					testSnapshotID.testUUID,
					StorageSnapshotUpdateRequest{
						Name:   "test",
						Labels: []string{"label"},
					})
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateStorageSnapshot returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_DeleteStorageSnapshot(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
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
		for _, testStorageID := range uuidCommonTestCases {
			for _, testSnapshotID := range uuidCommonTestCases {
				err := client.DeleteStorageSnapshot(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID)
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteStorageSnapshot returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_RollbackStorage(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "rollback")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, "")
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, testStorageID := range uuidCommonTestCases {
			for _, testSnapshotID := range uuidCommonTestCases {
				err := client.RollbackStorage(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID, StorageRollbackRequest{Rollback: true})
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "RollbackStorage returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_ExportStorageSnapshotToS3(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "export_to_s3")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, "")
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, testStorageID := range uuidCommonTestCases {
			for _, testSnapshotID := range uuidCommonTestCases {
				err := client.ExportStorageSnapshotToS3(
					emptyCtx,
					testStorageID.testUUID,
					testSnapshotID.testUUID,
					StorageSnapshotExportToS3Request{
						S3auth: S3auth{
							Host:      "example.com",
							AccessKey: "access_key",
							SecretKey: "secret_key",
						},
						S3data: S3data{
							Host:     "example.com",
							Bucket:   "bucket",
							Filename: "filename",
							Private:  true,
						},
					})
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "ExportStorageSnapshotToS3 returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_GetSnapshotsByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageSnapshotListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSnapshotsByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSnapshotsByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetDeletedSnapshots(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "snapshots")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareDeletedStorageSnapshotListHTTPGet())
	})

	res, err := client.GetDeletedSnapshots(emptyCtx)
	assert.Nil(t, err, "GetSnapshotsByLocation returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("deleted")), fmt.Sprintf("%v", res))
}

func getMockStorageSnapshot(status string) StorageSnapshot {
	mock := StorageSnapshot{Properties: StorageSnapshotProperties{
		Labels:           []string{"label"},
		ObjectUUID:       dummyUUID,
		Name:             "test",
		Status:           status,
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

func prepareStorageSnapshotHTTPGet(status string) string {
	snapshot := getMockStorageSnapshot(status)
	res, _ := json.Marshal(snapshot)
	return string(res)
}

func prepareStorageSnapshotListHTTPGet() string {
	snapshot := getMockStorageSnapshot("active")
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
	snapshot := getMockStorageSnapshot("deleted")
	res, _ := json.Marshal(snapshot.Properties)
	return fmt.Sprintf(`{"deleted_snapshots" : {"%s" : %s}}`, dummyUUID, string(res))
}

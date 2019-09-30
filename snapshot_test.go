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
<<<<<<< HEAD

	res, err := client.GetStorageSnapshotList(dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshotList returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageSnapshotList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageSnapshotList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD

	res, err := client.GetStorageSnapshot(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetStorageSnapshot returned an error %v", err)
=======
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
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshot()), fmt.Sprintf("%v", res))
}

func TestClient_CreateStorageSnapshot(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshots")
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprint(w, prepareStorageSnapshotCreateResponseHTTP())
			}
		})
		if clientTest {
			httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
			mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
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
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotCreateResponse()), fmt.Sprintf("%v", response))
}

func TestClient_UpdateStorageSnapshot(t *testing.T) {
<<<<<<< HEAD
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
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
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
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_DeleteStorageSnapshot(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
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
		server.Close()
	}
}

func TestClient_RollbackStorage(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "rollback")
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprintf(w, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID), func(w http.ResponseWriter, r *http.Request) {
				if isFailed {
					w.WriteHeader(400)
				} else {
					fmt.Fprint(w, prepareStorageSnapshotHTTPGet("active"))
				}
			})
		}
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
		server.Close()
	}
}

func TestClient_ExportStorageSnapshotToS3(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "export_to_s3")
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			if isFailed {
				w.WriteHeader(400)
			} else {
				fmt.Fprintf(w, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID), func(w http.ResponseWriter, r *http.Request) {
				if isFailed {
					w.WriteHeader(400)
				} else {
					fmt.Fprint(w, prepareStorageSnapshotHTTPGet("active"))
				}
			})
		}
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testStorageID := range uuidCommonTestCases {
				for _, testSnapshotID := range uuidCommonTestCases {
					err := client.ExportStorageSnapshotToS3(
						emptyCtx,
						testStorageID.testUUID,
						testSnapshotID.testUUID,
						StorageSnapshotExportToS3Request{
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
					if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "ExportStorageSnapshotToS3 returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_GetSnapshotsByLocation(t *testing.T) {
	server, client, mux := setupTestClient(true)
>>>>>>> 8d4aa0e... add `context`
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprint(w, "")
	})
<<<<<<< HEAD
	err := client.DeleteStorageSnapshot(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteStorageSnapshot returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSnapshotsByLocation(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSnapshotsByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("active")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
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
<<<<<<< HEAD
	err := client.RollbackStorage(dummyUUID, dummyUUID, StorageRollbackRequest{Rollback: true})
	if err != nil {
		t.Errorf("RollbackStorage returned an error %v", err)
	}
=======

	res, err := client.GetDeletedSnapshots(emptyCtx)
	assert.Nil(t, err, "GetSnapshotsByLocation returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshot("deleted")), fmt.Sprintf("%v", res))
>>>>>>> 8d4aa0e... add `context`
}

func TestClient_ExportStorageSnapshotToS3(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID, "export_to_s3")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprint(w, "")
	})
<<<<<<< HEAD
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
=======
	err := client.waitForSnapshotActive(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "waitForSnapshotActive returned an error %v", err)
}

func TestClient_waitForSnapshotDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshots", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(404)
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testSnapshotID := range uuidCommonTestCases {
			err := client.waitForSnapshotDeleted(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID)
			if testStorageID.isFailed || testSnapshotID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForSnapshotDeleted returned an error %v", err)
			}
		}
>>>>>>> 8d4aa0e... add `context`
	}
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

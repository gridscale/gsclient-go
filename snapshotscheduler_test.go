package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetStorageSnapshotScheduleList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageSnapshotScheduleList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageSnapshotScheduleList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshotSchedule("active")), fmt.Sprintf("%v", res))
		}
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshotSchedule()), fmt.Sprintf("%v", res))
}

func TestClient_GetStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleHTTPGet())
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testScheduleID := range uuidCommonTestCases {
			res, err := client.GetStorageSnapshotSchedule(emptyCtx, testStorageID.testUUID, testScheduleID.testUUID)
			if testStorageID.isFailed || testScheduleID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetStorageSnapshotSchedule returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotSchedule("active")), fmt.Sprintf("%v", res))
			}
		}
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotSchedule()), fmt.Sprintf("%v", res))
}

func TestClient_CreateStorageSnapshotSchedule(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules")
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprintf(writer, prepareStorageSnapshotScheduleHTTPCreateResponse())
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
				response, err := client.CreateStorageSnapshotSchedule(
					emptyCtx,
					test.testUUID,
					StorageSnapshotScheduleCreateRequest{
						Name:          "test",
						Labels:        []string{"test"},
						RunInterval:   60,
						KeepSnapshots: 1,
						NextRuntime:   &dummyTime,
					})
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "CreateStorageSnapshotSchedule returned an error %v", err)
					assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotScheduleHTTPCreateResponse()), fmt.Sprintf("%s", response))
				}
			}
		}
		server.Close()
	}

	assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotScheduleHTTPCreateResponse()), fmt.Sprintf("%s", response))
}

func TestClient_UpdateStorageSnapshotSchedule(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			if isFailed {
				writer.WriteHeader(400)
			} else {
				if request.Method == http.MethodPatch {
					fmt.Fprintf(writer, "")
				} else if request.Method == http.MethodGet {
					fmt.Fprint(writer, prepareStorageSnapshotScheduleHTTPGet("active"))
				}
			}
		})
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testStorageID := range uuidCommonTestCases {
				for _, testScheduleID := range uuidCommonTestCases {
					err := client.UpdateStorageSnapshotSchedule(
						emptyCtx,
						testStorageID.testUUID,
						testScheduleID.testUUID,
						StorageSnapshotScheduleUpdateRequest{
							Name:          "test",
							Labels:        []string{"label"},
							RunInterval:   60,
							KeepSnapshots: 1,
							NextRuntime:   &dummyTime,
						})
					if testStorageID.isFailed || testScheduleID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "UpdateStorageSnapshotSchedule returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_DeleteStorageSnapshotSchedule(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
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
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, testStorageID := range uuidCommonTestCases {
				for _, testScheduleID := range uuidCommonTestCases {
					err := client.DeleteStorageSnapshotSchedule(emptyCtx, testStorageID.testUUID, testScheduleID.testUUID)
					if testStorageID.isFailed || testScheduleID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "DeleteStorageSnapshotSchedule returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_waitForSnapshotScheduleActive(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})

	err := client.UpdateStorageSnapshotSchedule(dummyUUID, dummyUUID, StorageSnapshotScheduleUpdateRequest{
		Name:          "test",
		Labels:        []string{"label"},
		RunInterval:   60,
		KeepSnapshots: 1,
		NextRuntime:   dummyTime,
	})
	err := client.waitForSnapshotScheduleActive(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "waitForSnapshotScheduleActive returned an error %v", err)
}

func TestClient_DeleteStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testScheduleID := range uuidCommonTestCases {
			err := client.waitForSnapshotScheduleDeleted(emptyCtx, testStorageID.testUUID, testScheduleID.testUUID)
			if testStorageID.isFailed || testScheduleID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForSnapshotScheduleDeleted returned an error %v", err)
			}
		}
	}
}

func getMockStorageSnapshotSchedule() StorageSnapshotSchedule {
	mock := StorageSnapshotSchedule{Properties: StorageSnapshotScheduleProperties{
		ChangeTime:    dummyTime,
		CreateTime:    dummyTime,
		KeepSnapshots: 1,
		Labels:        []string{"label"},
		Name:          "test",
		NextRuntime:   dummyTime,
		ObjectUUID:    dummyUUID,
		Relations: StorageSnapshotScheduleRelations{Snapshots: []StorageSnapshotScheduleRelation{
			{
				CreateTime: dummyTime,
				Name:       "test",
				ObjectUUID: dummyUUID,
			},
		}},
		RunInterval: 60,
		Status:      "active",
		StorageUUID: dummyUUID,
	}}
	return mock
}

func prepareStorageSnapshotScheduleListHTTPGet() string {
	scheduler := getMockStorageSnapshotSchedule()
	res, _ := json.Marshal(scheduler.Properties)
	return fmt.Sprintf(`{"snapshot_schedules" : {"%s" : %s}}`, dummyUUID, string(res))
}

func prepareStorageSnapshotScheduleHTTPGet() string {
	scheduler := getMockStorageSnapshotSchedule()
	res, _ := json.Marshal(scheduler)
	return string(res)
}

func getMockStorageSnapshotScheduleHTTPCreateResponse() StorageSnapshotScheduleCreateResponse {
	mock := StorageSnapshotScheduleCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareStorageSnapshotScheduleHTTPCreateResponse() string {
	res, _ := json.Marshal(getMockStorageSnapshotScheduleHTTPCreateResponse())
	return string(res)
}

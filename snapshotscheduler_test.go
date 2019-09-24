package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageSnapshotScheduleList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageSnapshotScheduleList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageSnapshotScheduleList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageSnapshotSchedule("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetStorageSnapshotSchedule(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareStorageSnapshotScheduleHTTPGet("active"))
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testScheduleID := range uuidCommonTestCases {
			res, err := client.GetStorageSnapshotSchedule(testStorageID.testUUID, testScheduleID.testUUID)
			if testStorageID.isFailed || testScheduleID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetStorageSnapshotSchedule returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageSnapshotSchedule("active")), fmt.Sprintf("%v", res))
			}
		}
	}
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
			mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, httpResponse)
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			for _, test := range uuidCommonTestCases {
				response, err := client.CreateStorageSnapshotSchedule(test.testUUID, StorageSnapshotScheduleCreateRequest{
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
					err := client.UpdateStorageSnapshotSchedule(testStorageID.testUUID, testScheduleID.testUUID, StorageSnapshotScheduleUpdateRequest{
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
					err := client.DeleteStorageSnapshotSchedule(testStorageID.testUUID, testScheduleID.testUUID)
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
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(writer, prepareStorageSnapshotScheduleHTTPGet("in-provisioning"))
			} else {
				fmt.Fprint(writer, prepareStorageSnapshotScheduleHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			err := client.waitForSnapshotScheduleActive(dummyUUID, dummyUUID)
			if isFailed || isTimeout {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForSnapshotScheduleActive returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForSnapshotScheduleDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiStorageBase, dummyUUID, "snapshot_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprint(writer, prepareStorageSnapshotScheduleHTTPGet("to-be-deleted"))
			} else {
				writer.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, testStorageID := range uuidCommonTestCases {
				for _, testScheduleID := range uuidCommonTestCases {
					err := client.waitForSnapshotScheduleDeleted(testStorageID.testUUID, testScheduleID.testUUID)
					if testStorageID.isFailed || testScheduleID.isFailed || isFailed || isTimeout {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "waitForSnapshotScheduleDeleted returned an error %v", err)
					}
				}
			}
		}
	}
}

func getMockStorageSnapshotSchedule(status string) StorageSnapshotSchedule {
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
		Status:      status,
		StorageUUID: dummyUUID,
	}}
	return mock
}

func prepareStorageSnapshotScheduleListHTTPGet() string {
	scheduler := getMockStorageSnapshotSchedule("active")
	res, _ := json.Marshal(scheduler.Properties)
	return fmt.Sprintf(`{"snapshot_schedules" : {"%s" : %s}}`, dummyUUID, string(res))
}

func prepareStorageSnapshotScheduleHTTPGet(status string) string {
	scheduler := getMockStorageSnapshotSchedule(status)
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

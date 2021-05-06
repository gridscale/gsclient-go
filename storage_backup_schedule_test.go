package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetStorageBackupScheduleList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "backup_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageBackupScheduleListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageBackupScheduleList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageBackupScheduleList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageBackupSchedule("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetStorageBackupSchedule(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "backup_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageBackupScheduleHTTPGet("active"))
	})
	for _, testStorageID := range uuidCommonTestCases {
		for _, testScheduleID := range uuidCommonTestCases {
			res, err := client.GetStorageBackupSchedule(emptyCtx, testStorageID.testUUID, testScheduleID.testUUID)
			if testStorageID.isFailed || testScheduleID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetStorageBackupSchedule returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageBackupSchedule("active")), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_CreateStorageBackupSchedule(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "backup_schedules")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareStorageBackupScheduleHTTPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		for _, test := range uuidCommonTestCases {
			response, err := client.CreateStorageBackupSchedule(
				emptyCtx,
				test.testUUID,
				StorageBackupScheduleCreateRequest{
					Name:        "test",
					RunInterval: 60,
					KeepBackups: 1,
					NextRuntime: dummyTime,
					Active:      true,
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateStorageBackupSchedule returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockStorageBackupScheduleHTTPCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
	}
}

func TestClient_UpdateStorageBackupSchedule(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "backup_schedules", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodPatch {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				fmt.Fprint(writer, prepareStorageBackupScheduleHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, testStorageID := range uuidCommonTestCases {
			for _, testScheduleID := range uuidCommonTestCases {
				err := client.UpdateStorageBackupSchedule(
					emptyCtx,
					testStorageID.testUUID,
					testScheduleID.testUUID,
					StorageBackupScheduleUpdateRequest{
						Name:        "test",
						RunInterval: 60,
						KeepBackups: 1,
						NextRuntime: &dummyTime,
					})
				if testStorageID.isFailed || testScheduleID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "UpdateStorageBackupSchedule returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_DeleteStorageBackupSchedule(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "backup_schedules", dummyUUID)
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
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, testStorageID := range uuidCommonTestCases {
			for _, testScheduleID := range uuidCommonTestCases {
				err := client.DeleteStorageBackupSchedule(emptyCtx, testStorageID.testUUID, testScheduleID.testUUID)
				if testStorageID.isFailed || testScheduleID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteStorageBackupSchedule returned an error %v", err)
				}
			}
		}
	}
}

func getMockStorageBackupSchedule(status string) StorageBackupSchedule {
	mock := StorageBackupSchedule{Properties: StorageBackupScheduleProperties{
		ChangeTime:  dummyTime,
		CreateTime:  dummyTime,
		KeepBackups: 1,
		Name:        "test",
		NextRuntime: dummyTime,
		ObjectUUID:  dummyUUID,
		Relations: StorageBackupScheduleRelations{StorageBackups: []StorageBackupScheduleRelation{
			{
				CreateTime: dummyTime,
				Name:       "test",
				ObjectUUID: dummyUUID,
			},
		}},
		RunInterval: 60,
		Status:      status,
		StorageUUID: dummyUUID,
		Active:      status == "active",
	}}
	return mock
}

func prepareStorageBackupScheduleListHTTPGet() string {
	scheduler := getMockStorageBackupSchedule("active")
	res, _ := json.Marshal(scheduler.Properties)
	return fmt.Sprintf(`{"schedule_storage_backups" : {"%s" : %s}}`, dummyUUID, string(res))
}

func prepareStorageBackupScheduleHTTPGet(status string) string {
	scheduler := getMockStorageBackupSchedule(status)
	res, _ := json.Marshal(scheduler)
	return string(res)
}

func getMockStorageBackupScheduleHTTPCreateResponse() StorageBackupScheduleCreateResponse {
	mock := StorageBackupScheduleCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
	return mock
}

func prepareStorageBackupScheduleHTTPCreateResponse() string {
	res, _ := json.Marshal(getMockStorageBackupScheduleHTTPCreateResponse())
	return string(res)
}

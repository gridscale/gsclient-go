package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetStorageBackupList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "backups")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareStorageBackupListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetStorageBackupList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageBackupList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorageBackup()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_DeleteStorageBackup(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "backups", dummyUUID)
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
				err := client.DeleteStorageBackup(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID)
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteStorageBackup returned an error %v", err)
				}
			}
		}
	}
}

func TestClient_RollbackStorageBackup(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiStorageBase, dummyUUID, "backups", dummyUUID, "rollback")
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
				err := client.RollbackStorageBackup(emptyCtx, testStorageID.testUUID, testSnapshotID.testUUID, StorageRollbackRequest{Rollback: true})
				if testStorageID.isFailed || testSnapshotID.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "RollbackStorageBackup returned an error %v", err)
				}
			}
		}
	}
}

func getMockStorageBackup() StorageBackup {
	mock := StorageBackup{Properties: StorageBackupProperties{
		ObjectUUID: dummyUUID,
		Name:       "test",
		CreateTime: dummyTime,
		Capacity:   10,
	}}
	return mock
}

func prepareStorageBackupListHTTPGet() string {
	snapshot := getMockStorageBackup()
	res, _ := json.Marshal(snapshot.Properties)
	return fmt.Sprintf(`{"backups" : {"%s" : %s}}`, dummyUUID, string(res))
}

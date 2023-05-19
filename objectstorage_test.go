package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetObjectStorageAccessKeyList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareObjectStorageAccessKeyListHTTPGet())
	})

	res, err := client.GetObjectStorageAccessKeyList(emptyCtx)
	assert.Nil(t, err, "GetObjectStorageAccessKeyList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockObjectStorageAccessKey()), fmt.Sprintf("%v", res))
}

func TestClient_GetObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "access_keys", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareObjectStorageAccessKeyHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetObjectStorageAccessKey(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetObjectStorageAccessKey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockObjectStorageAccessKey()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, prepareObjectStorageAccessKeyHTTPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateObjectStorageAccessKey(emptyCtx)
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateObjectStorageAccessKey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockObjectStorageAccessKeyCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_AdvancedCreateObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, prepareObjectStorageAccessKeyHTTPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.AdvancedCreateObjectStorageAccessKey(emptyCtx, ObjectStorageAccessKeyCreateRequest{})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "AdvancedCreateObjectStorageAccessKey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockObjectStorageAccessKeyCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_UpdateObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, prepareObjectStorageAccessKeyHTTPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		err := client.UpdateObjectStorageAccessKey(emptyCtx, dummyUUID, ObjectStorageAccessKeyUpdateRequest{})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "AdvancedCreateObjectStorageAccessKey returned an error %v", err)
		}
	}
}

func TestClient_DeleteObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiObjectStorageBase, "access_keys", dummyUUID)
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
		for _, test := range uuidCommonTestCases {
			err := client.DeleteObjectStorageAccessKey(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteObjectStorageAccessKey returned an error %v", err)
			}
		}
	}
}

func TestClient_GetObjectStorageBucketList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "buckets")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareObjectStorageBucketListHTTPGet())
	})

	res, err := client.GetObjectStorageBucketList(emptyCtx)
	assert.Nil(t, err, "GetObjectStorageBucketList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockObjectStorageBucket()), fmt.Sprintf("%v", res))
}

func getMockObjectStorageAccessKey() ObjectStorageAccessKey {
	mock := ObjectStorageAccessKey{Properties: ObjectStorageAccessKeyProperties{
		SecretKey: "dummy-secret-key",
		AccessKey: "dummy-access-key",
		User:      "any",
	}}
	return mock
}

func prepareObjectStorageAccessKeyListHTTPGet() string {
	accessKey := getMockObjectStorageAccessKey()
	res, _ := json.Marshal(accessKey.Properties)
	return fmt.Sprintf(`{"access_keys": [%s]}`, string(res))
}

func prepareObjectStorageAccessKeyHTTPGet() string {
	accessKey := getMockObjectStorageAccessKey()
	res, _ := json.Marshal(accessKey)
	return string(res)
}

func getMockObjectStorageAccessKeyCreateResponse() ObjectStorageAccessKeyCreateResponse {
	mock := ObjectStorageAccessKeyCreateResponse{
		AccessKey: struct {
			SecretKey string `json:"secret_key"`
			AccessKey string `json:"access_key"`
		}{
			SecretKey: "dummy-secret-key",
			AccessKey: "dummy-access-key",
		},
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareObjectStorageAccessKeyHTTPCreateResponse() string {
	createRes := getMockObjectStorageAccessKeyCreateResponse()
	res, _ := json.Marshal(createRes)
	return string(res)
}

func getMockObjectStorageBucket() ObjectStorageBucket {
	mock := ObjectStorageBucket{Properties: ObjectStorageBucketProperties{
		Name: "test",
		Usage: struct {
			SizeKb     int `json:"size_kb"`
			NumObjects int `json:"num_objects"`
		}{
			SizeKb:     1000000,
			NumObjects: 10,
		},
	}}
	return mock
}

func prepareObjectStorageBucketListHTTPGet() string {
	bucket := getMockObjectStorageBucket()
	res, _ := json.Marshal(bucket.Properties)
	return fmt.Sprintf(`{"buckets": [%s]}`, string(res))
}

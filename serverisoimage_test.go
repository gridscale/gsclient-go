package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerIsoImageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageListHTTPGet())
	})
<<<<<<< HEAD
	res, err := client.GetServerIsoImageList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerIsoImageList returned an error %v", err)
=======
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerIsoImageList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerIsoImageList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIsoImage("test")), fmt.Sprintf("%v", res))
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIsoImage()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageHTTPget())
	})
<<<<<<< HEAD
	res, err := client.GetServerIsoImage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetServerIsoImage returned an error %v", err)
=======
	for _, testServerID := range uuidCommonTestCases {
		for _, testISOImageID := range uuidCommonTestCases {
			res, err := client.GetServerIsoImage(emptyCtx, testServerID.testUUID, testISOImageID.testUUID)
			if testServerID.isFailed || testISOImageID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetServerIsoImage returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerIsoImage("test")), fmt.Sprintf("%v", res))
			}
		}
>>>>>>> 8d4aa0e... add `context`
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerIsoImage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerIsoImage(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerIsoImage(dummyUUID, ServerIsoImageRelationCreateRequest{
		ObjectUUID: dummyUUID,
	})
	if err != nil {
		t.Errorf("CreateServerIsoImage returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "isoimages")
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprint(writer, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, http.MethodGet, request.Method)
				fmt.Fprintf(writer, prepareServerIPHTTPGet())
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			for _, testServerID := range uuidCommonTestCases {
				for _, testISOImageID := range uuidCommonTestCases {
					err := client.CreateServerIsoImage(
						emptyCtx,
						testServerID.testUUID,
						ServerIsoImageRelationCreateRequest{
							ObjectUUID: testISOImageID.testUUID,
						})
					if testServerID.isFailed || testISOImageID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "CreateServerIsoImage returned an error %v", err)
					}
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_UpdateServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testISOImageID := range uuidCommonTestCases {
			err := client.UpdateServerIsoImage(
				emptyCtx,
				testServerID.testUUID,
				testISOImageID.testUUID,
				ServerIsoImageRelationUpdateRequest{
					BootDevice: true,
					Name:       "test",
				})
			if testServerID.isFailed || testISOImageID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateServerIsoImage returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteServerIsoImage(t *testing.T) {
<<<<<<< HEAD
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerIsoImage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerIsoImage returned an error %v", err)
=======
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
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
				for _, testISOImageID := range uuidCommonTestCases {
					err := client.DeleteServerIsoImage(emptyCtx, testServerID.testUUID, testISOImageID.testUUID)
					if testServerID.isFailed || testISOImageID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "DeleteServerIsoImage returned an error %v", err)
					}
				}
			}
		}
		server.Close()
>>>>>>> 8d4aa0e... add `context`
	}
}

func TestClient_LinkIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
<<<<<<< HEAD
	err := client.LinkIsoImage(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("LinkIsoImage returned an error %v", err)
	}
=======
	mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIPHTTPGet())
	})
	err := client.LinkIsoImage(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "LinkIsoImage returned an error %v", err)

>>>>>>> 8d4aa0e... add `context`
}

func TestClient_UnlinkIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkIsoImage(emptyCtx, dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkIsoImage returned an error %v", err)
}

<<<<<<< HEAD
func getMockServerIsoImage() ServerIsoImageRelationProperties {
=======
func TestClient_waitForServerISOImageRelCreation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerISOImageRelCreation(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerISOImageRelCreation returned an error %v", err)
			}
		}
	}
}

func TestClient_waitForServerISOImageRelDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(404)
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testIPID := range uuidCommonTestCases {
			err := client.waitForServerISOImageRelDeleted(emptyCtx, testServerID.testUUID, testIPID.testUUID)
			if testServerID.isFailed || testIPID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "waitForServerISOImageRelDeleted returned an error %v", err)
			}
		}
	}
}

func getMockServerIsoImage(name string) ServerIsoImageRelationProperties {
>>>>>>> 8d4aa0e... add `context`
	mock := ServerIsoImageRelationProperties{
		ObjectUUID: dummyUUID,
		ObjectName: "test",
		Private:    false,
		CreateTime: dummyTime,
		Bootdevice: true,
	}
	return mock
}

func prepareServerIsoImageListHTTPGet() string {
	iso := getMockServerIsoImage()
	res, _ := json.Marshal(iso)
	return fmt.Sprintf(`{"isoimage_relations": [%s]}`, string(res))
}

func prepareServerIsoImageHTTPget() string {
	iso := getMockServerIsoImage()
	res, _ := json.Marshal(iso)
	return fmt.Sprintf(`{"isoimage_relation": %s}`, string(res))
}

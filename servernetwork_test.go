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

func TestClient_GetServerNetworkList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerNetworkList(context.Background(), test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerNetworkList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerNetwork("test")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			res, err := client.GetServerNetwork(context.Background(), testServerID.testUUID, testNetworkID.testUUID)
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetServerNetwork returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerNetwork("test")), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_CreateServerNetwork(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "networks")
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprint(writer, "")
			}
		})
		if clientTest {
			mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "networks", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, http.MethodGet, request.Method)
				fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
			})
		}
		for _, test := range commonSuccessFailTestCases {
			isFailed = test.isFailed
			for _, testServerID := range uuidCommonTestCases {
				for _, testNetworkID := range uuidCommonTestCases {
					err := client.CreateServerNetwork(
						context.Background(),
						testServerID.testUUID,
						ServerNetworkRelationCreateRequest{
							ObjectUUID:           testNetworkID.testUUID,
							Ordering:             1,
							BootDevice:           false,
							L3security:           nil,
							FirewallTemplateUUID: dummyUUID,
						})
					if testServerID.isFailed || testNetworkID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "CreateServerNetwork returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_UpdateServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient(false)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			err := client.UpdateServerNetwork(
				context.Background(),
				testServerID.testUUID,
				testNetworkID.testUUID,
				ServerNetworkRelationUpdateRequest{
					Ordering:             0,
					BootDevice:           true,
					FirewallTemplateUUID: dummyUUID,
				})
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateServerNetwork returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteServerNetwork(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
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
				for _, testNetworkID := range uuidCommonTestCases {
					err := client.DeleteServerNetwork(context.Background(), testServerID.testUUID, testNetworkID.testUUID)
					if testServerID.isFailed || testNetworkID.isFailed || isFailed {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "DeleteServerNetwork returned an error %v", err)
					}
				}
			}
		}
		server.Close()
	}
}

func TestClient_LinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, "")
		}
	})
	mux.HandleFunc(path.Join(apiServerBase, dummyUUID, "networks", dummyUUID), func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
	})
	err := client.LinkNetwork(context.Background(), dummyUUID, dummyUUID, dummyUUID, false, 1, nil, nil)
	assert.Nil(t, err, "LinkNetwork returned an error %v", err)
}

func TestClient_UnlinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			fmt.Fprintf(writer, "")
		} else if request.Method == http.MethodGet {
			writer.WriteHeader(404)
		}
	})
	err := client.UnlinkNetwork(context.Background(), dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkNetwork returned an error %v", err)
}

func TestClient_waitForServerNetworkRelCreation(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				writer.WriteHeader(404)
			} else {
				fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, isTimeoutTest := range timeoutTestCases {
			isTimeout = isTimeoutTest
			for _, testServerID := range uuidCommonTestCases {
				for _, testIPID := range uuidCommonTestCases {
					err := client.waitForServerNetworkRelCreation(context.Background(), testServerID.testUUID, testIPID.testUUID)
					if testServerID.isFailed || testIPID.isFailed || isFailed || isTimeout {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "waitForServerNetworkRelCreation returned an error %v", err)
					}
				}
			}
		}
	}
}

func TestClient_waitForServerNetworkRelDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	var isTimeout bool
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if isTimeout {
				fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
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
					err := client.waitForServerNetworkRelDeleted(context.Background(), testServerID.testUUID, testIPID.testUUID)
					if testServerID.isFailed || testIPID.isFailed || isFailed || isTimeout {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err, "waitForServerNetworkRelDeleted returned an error %v", err)
					}
				}
			}
		}
	}
}

func getMockServerNetwork(name string) ServerNetworkRelationProperties {
	mock := ServerNetworkRelationProperties{
		L2security:           true,
		ServerUUID:           dummyUUID,
		CreateTime:           dummyTime,
		PublicNet:            false,
		FirewallTemplateUUID: dummyUUID,
		ObjectName:           name,
		Mac:                  "",
		BootDevice:           true,
		PartnerUUID:          dummyUUID,
		Ordering:             0,
		Firewall:             "",
		NetworkType:          "",
		NetworkUUID:          dummyUUID,
		ObjectUUID:           dummyUUID,
		L3security:           nil,
	}
	return mock
}

func prepareServerNetworkListHTTPGet() string {
	net := getMockServerNetwork("test")
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relations": [%s]}`, string(res))
}

func prepareServerNetworkHTTPGet() string {
	net := getMockServerNetwork("test")
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relation": %s}`, string(res))
}

func prepareServerNetworkHTTPGetForUpdate(value ServerNetworkRelationProperties) string {
	res, _ := json.Marshal(value)
	return fmt.Sprintf(`{"network_relation": %s}`, string(res))
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerIpList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIpListHTTPGet())
	})
	res, err := client.GetServerIpList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerIpList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIp()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIpHTTPGet())
	})
	res, err := client.GetServerIp(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetServerIp returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerIp()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerIp(dummyUUID, ServerIpRelationCreateRequest{
		ObjectUUID: dummyUUID,
	})
	if err != nil {
		t.Errorf("CreateServerIp returned an error %v", err)
	}
}

func TestClient_DeleteServerIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerIp(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerIp returned an error %v", err)
	}
}

func TestClient_LinkIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkIp(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("LinkIp returned an error %v", err)
	}
}

func TestClient_UnlinkIp(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkIp(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("UnlinkIp returned an error %v", err)
	}
}

func getMockServerIp() ServerIpRelationProperties {
	mock := ServerIpRelationProperties{
		ServerUUID: dummyUUID,
		CreateTime: dummyTime,
		Prefix:     "pre",
		Family:     1,
		ObjectUUID: dummyUUID,
		Ip:         "192.168.0.1",
	}
	return mock
}

func prepareServerIpListHTTPGet() string {
	ip := getMockServerIp()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relations": [%s]}`, string(res))
}

func prepareServerIpHTTPGet() string {
	ip := getMockServerIp()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relation": %s}`, string(res))
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerIPList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIPListHTTPGet())
	})
	res, err := client.GetServerIPList(dummyUUID)
	if err != nil {
		t.Errorf("GetServerIPList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIPHTTPGet())
	})
	res, err := client.GetServerIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("GetServerIP returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockServerIP()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerIP(dummyUUID, ServerIPRelationCreateRequest{
		ObjectUUID: dummyUUID,
	})
	if err != nil {
		t.Errorf("CreateServerIP returned an error %v", err)
	}
}

func TestClient_DeleteServerIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("DeleteServerIP returned an error %v", err)
	}
}

func TestClient_LinkIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("LinkIP returned an error %v", err)
	}
}

func TestClient_UnlinkIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "ips", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkIP(dummyUUID, dummyUUID)
	if err != nil {
		t.Errorf("UnlinkIP returned an error %v", err)
	}
}

func getMockServerIP() ServerIPRelationProperties {
	mock := ServerIPRelationProperties{
		ServerUUID: dummyUUID,
		CreateTime: dummyTime,
		Prefix:     "pre",
		Family:     1,
		ObjectUUID: dummyUUID,
		IP:         "192.168.0.1",
	}
	return mock
}

func prepareServerIPListHTTPGet() string {
	ip := getMockServerIP()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relations": [%s]}`, string(res))
}

func prepareServerIPHTTPGet() string {
	ip := getMockServerIP()
	res, _ := json.Marshal(ip)
	return fmt.Sprintf(`{"ip_relation": %s}`, string(res))
}

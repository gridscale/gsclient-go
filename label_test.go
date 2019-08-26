package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetLabelList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiLabelBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLabelListHTTPGet())
	})
	res, err := client.GetLabelList()
	if err != nil {
		t.Errorf("GetLabelList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLabel()), fmt.Sprintf("%v", res))
}

func TestClient_CreateLabel(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiLabelBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, prepareLabelCreateResponse())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	res, err := client.CreateLabel(LabelCreateRequest{Label:"test"})
	if err != nil {
		t.Errorf("CreateLabel returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockLabelCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_DeleteLabel(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLabelBase, "test")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteLabel("test")
	if err != nil {
		t.Errorf("DeleteLabel returned an error %v", err)
	}
}

func getMockLabel() Label {
	mock := Label{Properties: LabelProperties{
		Label:      "test",
		CreateTime: dummyTime,
		ChangeTime: dummyTime,
		Relations:  nil,
		Status:     "active",
	}}
	return mock
}

func getMockLabelCreateResponse() CreateResponse {
	mock := CreateResponse{
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareLabelListHTTPGet() string {
	label := getMockLabel()
	res, _ := json.Marshal(label.Properties)
	return fmt.Sprintf(`{"labels": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareLabelCreateResponse() string {
	response := getMockLabelCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

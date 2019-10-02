package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

var labelTestCases = []uuidTestCase{
	{
		testUUID: "test",
		isFailed: false,
	},
	{
		testUUID: "",
		isFailed: true,
	},
}

func TestClient_GetLabelList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiLabelBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLabelListHTTPGet("test"))
	})
	res, err := client.GetLabelList()
	assert.Nil(t, err, "GetLabelList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLabel("test")), fmt.Sprintf("%v", res))
}

func TestClient_CreateLabel(t *testing.T) {
	server, client, mux := setupTestClient(true)
	var isFailed bool
	uri := apiLabelBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, prepareLabelCreateResponse())
		}
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc(requestBase, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateLabel(LabelCreateRequest{Label: "test"})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateLabel returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockLabelCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
	server.Close()
}

func TestClient_waitForLabelDeleted(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiLabelBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareLabelListHTTPGet("not-test"))
	})
	for _, test := range labelTestCases {
		err := client.waitForLabelDeleted(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "waitForFirewallDeleted returned an error %v", err)
		}
	}
}

func TestClient_DeleteLabel(t *testing.T) {
	for _, clientTest := range syncClientTestCases {
		server, client, mux := setupTestClient(clientTest)
		var isFailed bool
		uri := path.Join(apiLabelBase, "test")
		mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodDelete, request.Method)
			if isFailed {
				writer.WriteHeader(400)
			} else {
				fmt.Fprint(writer, "")
			}
		})
		if clientTest {
			mux.HandleFunc(apiLabelBase, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				fmt.Fprint(w, prepareLabelListHTTPGet("not-test"))
			})
		}
		for _, serverTest := range commonSuccessFailTestCases {
			isFailed = serverTest.isFailed
			for _, test := range labelTestCases {
				err := client.DeleteLabel(test.testUUID)
				if test.isFailed || isFailed {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err, "DeleteLabel returned an error %v", err)
				}
			}
		}
		server.Close()
	}
}

func getMockLabel(label string) Label {
	mock := Label{Properties: LabelProperties{
		Label:      label,
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

func prepareLabelListHTTPGet(labelName string) string {
	label := getMockLabel(labelName)
	res, _ := json.Marshal(label.Properties)
	return fmt.Sprintf(`{"labels": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareLabelCreateResponse() string {
	response := getMockLabelCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

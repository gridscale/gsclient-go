package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, prepareLabelListHTTPGet("test"))
	})
	res, err := client.GetLabelList(emptyCtx)
	assert.Nil(t, err, "GetLabelList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLabel("test")), fmt.Sprintf("%v", res))
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

func prepareLabelListHTTPGet(labelName string) string {
	label := getMockLabel(labelName)
	res, _ := json.Marshal(label.Properties)
	return fmt.Sprintf(`{"labels": {"%s": %s}}`, dummyUUID, string(res))
}

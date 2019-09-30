package gsclient

import (
	"context"
	"net/http"
	"net/http/httptest"
)

const (
	dummyUUID        = "690de890-13c0-4e76-8a01-e10ba8786e53"
	dummyTime        = "2018-04-28T09:47:41Z"
	dummyRequestUUID = "x123xx1x-123x-1x12-123x-123xxx123x1x"
)
<<<<<<< HEAD

func setupTestClient() (*httptest.Server, *Client, *http.ServeMux) {
=======
var emptyCtx = context.Background()
var dummyTimeOriginal, _ = time.Parse(gsTimeLayout, "2018-04-28T09:47:41Z")
var dummyTime = GSTime{dummyTimeOriginal}

type uuidTestCase struct {
	isFailed bool
	testUUID string
}

type successFailTestCase struct {
	isFailed bool
}

var commonSuccessFailTestCases []successFailTestCase = []successFailTestCase{
	{
		isFailed: true,
	},
	{
		isFailed: false,
	},
}

var uuidCommonTestCases []uuidTestCase = []uuidTestCase{
	{
		testUUID: dummyUUID,
		isFailed: false,
	},
	{
		testUUID: "",
		isFailed: true,
	},
}

var syncClientTestCases []bool = []bool{true, false}
var timeoutTestCases []bool = []bool{true, false}

func setupTestClient(sync bool) (*httptest.Server, *Client, *http.ServeMux) {
>>>>>>> 8d4aa0e... add `context`
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := NewConfiguration(server.URL, "uuid", "token", true)
	return server, NewClient(config), mux
}

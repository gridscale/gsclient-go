package gsclient

import (
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	dummyUUID        = "690de890-13c0-4e76-8a01-e10ba8786e53"
	dummyRequestUUID = "x123xx1x-123x-1x12-123x-123xxx123x1x"
)

var dummyTimeOriginal, _ = time.Parse(gsTimeLayout, "2018-04-28T09:47:41Z")
var dummyTime = JSONTime{dummyTimeOriginal}

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

func setupTestClient() (*httptest.Server, *Client, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := NewConfiguration(server.URL, "uuid", "token", true, 5, 500, 50)
	return server, NewClient(config), mux
}

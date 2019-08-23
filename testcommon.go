package gsclient

import (
	"net/http"
	"net/http/httptest"
)

const (
	dummyUUID        = "690de890-13c0-4e76-8a01-e10ba8786e53"
	dummyTime        = "2018-04-28T09:47:41Z"
	dummyRequestUUID = "x123xx1x-123x-1x12-123x-123xxx123x1x"
)

func setupTestClient() (*httptest.Server, *Client, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := NewConfiguration(server.URL, "uuid", "token", true)
	return server, NewClient(config), mux
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetMarketplaceApplicationList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiMarketplaceApplicationBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareMarketplaceApplicationListHTTPGet())
	})
	response, err := client.GetMarketplaceApplicationList(emptyCtx)
	assert.Nil(t, err, "GetMarketplaceApplicationList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockMarketplaceApplication("active")), fmt.Sprintf("%v", response))
}

func TestClient_GetMarketplaceApplication(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiMarketplaceApplicationBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareMarketplaceApplicationHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetMarketplaceApplication(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetMarketplaceApplication returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceApplication("active")), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_CreateMarketplaceApplication(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiMarketplaceApplicationBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareMarketplaceApplicationCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateMarketplaceApplication(
			emptyCtx,
			MarketplaceApplicationCreateRequest{
				Name:              "test",
				ObjectStoragePath: "s3://test/export/test.gz",
				Category:          MarketplaceApplicationCloudStorageCategory,
				Setup: MarketplaceApplicationSetup{
					Cores:    1,
					Memory:   2,
					Capacity: 10,
				},
				Metadata: nil,
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateMarketplaceApplication returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceApplicationCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_ImportMarketplaceApplication(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiMarketplaceApplicationBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, prepareMarketplaceApplicationCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.ImportMarketplaceApplication(
			emptyCtx,
			MarketplaceApplicationImportRequest{
				UniqueHash: "hash",
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "ImportMarketplaceApplication returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockMarketplaceApplicationCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_UpdateMarketplaceApplication(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiMarketplaceApplicationBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodPatch {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				fmt.Fprint(w, prepareMarketplaceApplicationHTTPGet("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdateMarketplaceApplication(
				emptyCtx,
				test.testUUID,
				MarketplaceApplicationUpdateRequest{
					Name:              "test new",
					ObjectStoragePath: "s3://test/export/test_new.gz",
					Category:          MarketplaceApplicationAdminpanelCategory,
					Setup: &MarketplaceApplicationSetup{
						Cores:    2,
						Memory:   4,
						Capacity: 20,
					},
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateMarketplaceApplication returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteMarketplaceApplication(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiMarketplaceApplicationBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodDelete {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				w.WriteHeader(404)
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.DeleteMarketplaceApplication(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteMarketplaceApplication returned an error %v", err)
			}
		}
	}
}

func TestClient_GetMarketplaceApplicationEventList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiMarketplaceApplicationBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetMarketplaceApplicationEventList(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetMarketplaceApplicationEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}
	}
}

func getMockMarketplaceApplication(status string) MarketplaceApplication {
	mock := MarketplaceApplication{Properties: MarketplaceApplicationProperties{
		Name:               "test",
		UniqueHash:         "hash",
		ObjectStoragePath:  "s3://test/export/test.gz",
		IsApplicationOwner: true,
		Setup: MarketplaceApplicationSetup{
			Cores:    1,
			Memory:   2,
			Capacity: 10,
		},
		Category: "Cloud Storage",
		Metadata: MarketplaceApplicationMetadata{
			License:    "test",
			OS:         "test",
			Overview:   "test",
			Hints:      "test",
			Icon:       "test",
			Features:   "test",
			TermsOfUse: "test",
			Author:     "test",
			Advices:    "test",
		},
		ChangeTime:      dummyTime,
		CreateTime:      dummyTime,
		ObjectUUID:      dummyUUID,
		Status:          status,
		ApplicationType: "test",
	}}
	return mock
}

func getMockMarketplaceApplicationCreateResponse() MarketplaceApplicationCreateResponse {
	mock := MarketplaceApplicationCreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareMarketplaceApplicationListHTTPGet() string {
	marketApp := getMockMarketplaceApplication("active")
	res, _ := json.Marshal(marketApp.Properties)
	return fmt.Sprintf(`{"applications": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareMarketplaceApplicationHTTPGet(status string) string {
	marketApp := getMockMarketplaceApplication(status)
	res, _ := json.Marshal(marketApp)
	return string(res)
}

func prepareMarketplaceApplicationCreateResponse() string {
	response := getMockMarketplaceApplicationCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetSSLCertificateList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := apiSSLCertificateBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareSSLCertificateListHTTPGet())
	})
	res, err := client.GetSSLCertificateList(emptyCtx)
	assert.Nil(t, err, "GetSSLCertificateList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSSLCertificate("active")), fmt.Sprintf("%v", res))
}

func TestClient_GetSSLCertificate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiSSLCertificateBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		fmt.Fprintf(writer, prepareSSLCertificateHTTPGet("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetSSLCertificate(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetSSLCertificate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSSLCertificate("active")), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_SSLCertificate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := apiSSLCertificateBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareSSLCertificateCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateSSLCertificate(
			emptyCtx,
			SSLCertificateCreateRequest{
				Name:            "test",
				PrivateKey:      "-----BEGIN RSA PRIVATE KEY-----abc-----END RSA PRIVATE KEY-----",
				LeafCertificate: "-----BEGIN CERTIFICATE-----abc-----END CERTIFICATE-----",
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateSSLCertificate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSSLCertificateCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_DeleteSSLCertificate(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiSSLCertificateBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeader, dummyRequestUUID)
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
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.DeleteSSLCertificate(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteSSLCertificate returned an error %v", err)
			}
		}
	}
}

func getMockSSLCertificate(status string) SSLCertificate {
	mock := SSLCertificate{Properties: SSLCertificateProperties{
		Name:          "test",
		ObjectUUID:    dummyUUID,
		Status:        status,
		CreateTime:    dummyTime,
		ChangeTime:    dummyTime,
		Labels:        []string{"label"},
		CommonName:    "www.example.com",
		NotValidAfter: dummyTime,
		Fingerprints: FingerprintProperties{
			MD5:    "6F:92:07:C2:01:2A:AF:A2:0C:0F:2E:21:4A:10:DB:64",
			SHA1:   "98:7D:24:03:E6:E3:B8:10:68:CD:89:C6:3A:78:66:3F:E2:61:A5:03",
			SHA256: "13:4B:BF:D0:16:9E:7F:15:56:48:E6:E4:8E:38:6A:B6:E6:BE:84:4F:C4:6B:C6:4B:57:C0:41:22:2F:0A:C0:3A",
		},
	}}
	return mock
}

func getMockSSLCertificateCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareSSLCertificateListHTTPGet() string {
	key := getMockSSLCertificate("active")
	res, _ := json.Marshal(key.Properties)
	return fmt.Sprintf(`{"certificates": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareSSLCertificateHTTPGet(status string) string {
	key := getMockSSLCertificate(status)
	res, _ := json.Marshal(key)
	return string(res)
}

func prepareSSLCertificateCreateResponse() string {
	response := getMockSSLCertificateCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

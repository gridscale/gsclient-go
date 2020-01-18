package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetPaaSServiceList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "services")
	expectedObj := getMockPaaSService("active")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, preparePaaSHTTPGetListResponse("active"))
	})
	paasList, err := client.GetPaaSServiceList(emptyCtx)
	assert.Nil(t, err, "GetPaaSServiceList returned an error %v", err)
	assert.Equal(t, 1, len(paasList))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObj), fmt.Sprintf("%v", paasList))
}

func TestClient_GetPaaSService(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "services", dummyUUID)
	expectedObj := getMockPaaSService("active")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, preparePaaSHTTPGetResponse("active"))
	})
	for _, test := range uuidCommonTestCases {
		paas, err := client.GetPaaSService(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetPaaSService returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", expectedObj.Properties), fmt.Sprintf("%v", paas.Properties))
		}
	}
}

func TestClient_CreatePaaSService(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "services")
	expectedRespObj := getMockPaaSServiceCreateResponse()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprintf(w, preparePaaSHTTPCreateResponse())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreatePaaSService(
			emptyCtx,
			PaaSServiceCreateRequest{
				Name:                    "test",
				PaaSServiceTemplateUUID: "test-template",
				Labels:                  []string{"label"},
				PaaSSecurityZoneUUID:    "test-security-zone-id",
				ResourceLimits: []ResourceLimit{
					{
						Resource: "cpu",
						Limit:    2,
					},
				},
				Parameters: nil,
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreatePaaSService returned error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", expectedRespObj), fmt.Sprintf("%v", response))
		}
	}

}

func TestClient_UpdatePaaSService(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "services", dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			w.WriteHeader(400)
		} else {
			if r.Method == http.MethodPatch {
				fmt.Fprintf(w, "")
			} else if r.Method == http.MethodGet {
				fmt.Fprint(w, preparePaaSHTTPGetResponse("active"))
			}
		}
	})
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdatePaaSService(
				emptyCtx,
				test.testUUID,
				PaaSServiceUpdateRequest{
					Name:       "test",
					Labels:     &[]string{"label"},
					Parameters: parameters,
					ResourceLimits: []ResourceLimit{
						{
							Resource: "cpu",
							Limit:    2,
						},
					},
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdatePaaSService returned an error %v", err)
			}
		}
	}
}

func TestClient_DeletePaaSService(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "services", dummyUUID)
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
			err := client.DeletePaaSService(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeletePaaSService returned an error %v", err)
			}
		}
	}
}

func TestClient_GetPaaSServiceMetrics(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "services", dummyUUID, "metrics")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, preparePaaSHTTPGetMetricsResponse())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetPaaSServiceMetrics(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetPaaSServiceMetrics returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockPaaSServiceMetric()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetPaaSTemplateList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "service_templates")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, preparePaaSHTTPGetTemplatesResponse())
	})
	res, err := client.GetPaaSTemplateList(emptyCtx)
	assert.Nil(t, err, "GetPaaSTemplateList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockPaasTemplate()), fmt.Sprintf("%v", res))
}

func TestClient_GetPaaSSecurityZoneList(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "security_zones")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, preparePaaSHTTPGetSecurityZoneList("active"))
	})
	res, err := client.GetPaaSSecurityZoneList(emptyCtx)
	assert.Nil(t, err, "GetPaaSSecurityZone returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSecurityZone("active")), fmt.Sprintf("%v", res))
}

func TestClient_CreatePaaSSecurityZone(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "security_zones")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodPost)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, preparePaaSHTTPCreateSecurityZone())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreatePaaSSecurityZone(
			emptyCtx,
			PaaSSecurityZoneCreateRequest{
				Name: "test",
			})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreatePaaSSecurityZone returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockPaaSSecurityZoneCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetPaaSSecurityZone(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiPaaSBase, "security_zones", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprintf(writer, preparePaaSHTTPGetSecurityZone("active"))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetPaaSSecurityZone(emptyCtx, test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetPaaSSecurityZone returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockSecurityZone("active")), fmt.Sprintf("%s", res))
		}
	}
}

func TestClient_UpdatePaaSSecurityZone(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "security_zones", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			if request.Method == http.MethodPatch {
				fmt.Fprintf(writer, "")
			} else if request.Method == http.MethodGet {
				fmt.Fprint(writer, preparePaaSHTTPGetSecurityZone("active"))
			}
		}
	})
	for _, serverTest := range commonSuccessFailTestCases {
		isFailed = serverTest.isFailed
		for _, test := range uuidCommonTestCases {
			err := client.UpdatePaaSSecurityZone(
				emptyCtx,
				test.testUUID,
				PaaSSecurityZoneUpdateRequest{
					Name:                 "test",
					PaaSSecurityZoneUUID: dummyUUID,
				})
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdatePaaSSecurityZone returned an error %v", err)
			}
		}
	}
}

func TestClient_DeletePaaSSecurityZone(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiPaaSBase, "security_zones", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
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
			err := client.DeletePaaSSecurityZone(emptyCtx, test.testUUID)
			if test.isFailed || isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeletePaaSSecurityZone returned an error %v", err)
			}
		}
	}
}

func TestClient_GetDeletedPaaSServices(t *testing.T) {
	server, client, mux := setupTestClient(true)
	defer server.Close()
	uri := path.Join(apiDeletedBase, "paas_services")
	expectedObj := getMockPaaSService("active")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		w.Header().Set(requestUUIDHeaderParam, dummyRequestUUID)
		fmt.Fprint(w, prepareDeletedPaaSHTTPGetListResponse("active"))
	})
	paasList, err := client.GetDeletedPaaSServices(emptyCtx)
	assert.Nil(t, err, "GetDeletedPaaSServices returned an error %v", err)
	assert.Equal(t, 1, len(paasList))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObj), fmt.Sprintf("%v", paasList))
}

func getMockPaaSService(status string) PaaSService {
	listenPort := make(map[string]map[string]int)
	portmap := make(map[string]int)
	portmap["mysql"] = 3306
	listenPort["fcfc::1:305e:6eff:fe62:4503"] = portmap
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	mock := PaaSService{
		Properties: PaaSServiceProperties{
			ObjectUUID: dummyUUID,
			Labels:     []string{"label"},
			Credentials: []Credential{
				{
					Username: "username",
					Password: "password",
					Type:     "type",
				},
			},
			CreateTime:          dummyTime,
			ListenPorts:         listenPort,
			SecurityZoneUUID:    "d711fc50-ad96-4070-b769-6fe2bf93792c",
			ServiceTemplateUUID: "504e2d11-7255-4712-b744-fcb093a4e613",
			UsageInMinutes:      999,
			CurrentPrice:        5.789,
			ChangeTime:          dummyTime,
			Status:              status,
			Name:                "test",
			ResourceLimits: []ResourceLimit{
				{
					Resource: "cpu",
					Limit:    2,
				},
			},
			Parameters: parameters,
		},
	}
	return mock
}

func getMockPaaSServiceMetric() PaaSServiceMetric {
	mock := PaaSServiceMetric{Properties: PaaSMetricProperties{
		BeginTime:       dummyTime,
		EndTime:         dummyTime,
		PaaSServiceUUID: dummyUUID,
		CoreUsage: PaaSMetricValue{
			Value: 50,
			Unit:  "percentage",
		},
		StorageSize: PaaSMetricValue{
			Value: 128,
			Unit:  "GB",
		},
	}}
	return mock
}

func preparePaaSHTTPGetMetricsResponse() string {
	metrics := getMockPaaSServiceMetric()
	res, _ := json.Marshal(metrics.Properties)
	return fmt.Sprintf(`{"paas_service_metrics": [%s]}`, string(res))
}

func preparePaaSHTTPGetListResponse(status string) string {
	paas := getMockPaaSService(status)
	res, _ := json.Marshal(paas.Properties)
	return fmt.Sprintf(`{"paas_services": {"%s" : %s}}`, dummyUUID, string(res))
}

func preparePaaSHTTPGetResponse(status string) string {
	paas := getMockPaaSService(status)
	res, _ := json.Marshal(paas)
	return string(res)
}

func preparePaaSHTTPCreateResponse() string {
	paasCreateResponse := getMockPaaSServiceCreateResponse()
	res, _ := json.Marshal(paasCreateResponse)
	return string(res)
}

func getMockPaaSServiceCreateResponse() PaaSServiceCreateResponse {
	listenPort := make(map[string]map[string]int)
	portmap := make(map[string]int)
	portmap["mysql"] = 3306
	portmap["http"] = 80
	listenPort["fcfc::1:aaaa:bbbb:cccc:dddd"] = portmap
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	return PaaSServiceCreateResponse{
		RequestUUID:     dummyRequestUUID,
		ListenPorts:     listenPort,
		PaaSServiceUUID: dummyUUID,
		Credentials: []Credential{
			{
				Username: "username",
				Password: "password",
				Type:     "type",
			},
		},
		ObjectUUID: dummyUUID,
		ResourceLimits: []ResourceLimit{
			{
				Resource: "cpu",
				Limit:    2,
			},
		},
		Parameters: parameters,
	}
}

func getMockPaasTemplate() PaaSTemplate {
	mock := PaaSTemplate{Properties: PaaSTemplateProperties{
		Name:       "test",
		ObjectUUID: "d711fc50-ad96-4070-b769-6fe2bf93792c",
		Category:   "database",
		ProductNo:  0,
		Labels:     []string{"label"},
		Resources: Resource{
			Memory:      10,
			Connections: 10,
		},
		Status: "active",
	}}
	return mock
}

func preparePaaSHTTPGetTemplatesResponse() string {
	template := getMockPaasTemplate()
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"paas_service_templates": {"%s" : %s}}`, "d711fc50-ad96-4070-b769-6fe2bf93792c", string(res))
}

func getMockSecurityZone(status string) PaaSSecurityZone {
	mock := PaaSSecurityZone{Properties: PaaSSecurityZoneProperties{
		LocationCountry: "Germany",
		CreateTime:      dummyTime,
		LocationIata:    "none",
		ObjectUUID:      "aa-bb-cc-dd",
		Labels:          []string{"label"},
		LocationName:    "Bonn",
		Status:          status,
		LocationUUID:    "cc-dd-ee",
		ChangeTime:      dummyTime,
		Name:            "test",
		Relation:        PaaSRelationService{Services: []ServiceObject{{ObjectUUID: "ff-gg-hh"}}},
	}}
	return mock
}

func preparePaaSHTTPGetSecurityZoneList(status string) string {
	zone := getMockSecurityZone(status)
	res, _ := json.Marshal(zone.Properties)
	return fmt.Sprintf(`{"paas_security_zones": {"%s": %s}}`, "test", string(res))
}

func preparePaaSHTTPGetSecurityZone(status string) string {
	zone := getMockSecurityZone(status)
	res, _ := json.Marshal(zone)
	return string(res)
}

func getMockPaaSSecurityZoneCreateResponse() PaaSSecurityZoneCreateResponse {
	return PaaSSecurityZoneCreateResponse{
		RequestUUID:          dummyRequestUUID,
		PaaSSecurityZoneUUID: dummyUUID,
		ObjectUUID:           dummyUUID,
	}
}

func preparePaaSHTTPCreateSecurityZone() string {
	res, _ := json.Marshal(getMockPaaSSecurityZoneCreateResponse())
	return string(res)
}

func prepareDeletedPaaSHTTPGetListResponse(status string) string {
	paas := getMockPaaSService(status)
	res, _ := json.Marshal(paas.Properties)
	return fmt.Sprintf(`{"deleted_paas_services": {"%s" : %s}}`, dummyUUID, string(res))
}

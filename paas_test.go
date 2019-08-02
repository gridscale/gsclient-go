package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

const (
	paasServiceUuid  = "690de890-13c0-4e76-8a01-e10ba8786e53"
	securityZoneUuid = "abc-123"
)

func TestClient_GetPaaSServiceList(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services")
	expectedObj := getMockPaaSService()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, preparePaaSHTTPGetListResponse())
	})
	paasList, err := client.GetPaaSServiceList()
	if err != nil {
		t.Errorf("GetPaaSServiceList returned an error %v", err)
	}
	assert.Equal(t, 1, len(paasList))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObj), fmt.Sprintf("%v", paasList))
}

func TestClient_GetPaaSService(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services", paasServiceUuid)
	expectedObj := getMockPaaSService()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, preparePaaSHTTPGetResponse())
	})
	paas, err := client.GetPaaSService(paasServiceUuid)
	if err != nil {
		t.Errorf("GetPaaSService returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", expectedObj.Properties), fmt.Sprintf("%v", paas.Properties))
}

func TestClient_CreatePaaSService(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services")
	expectedRespObj := getMockPaaSServiceCreateResponse()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		fmt.Fprintf(w, preparePaaSHTTPCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, requestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	response, err := client.CreatePaaSService(PaaSServiceCreateRequest{
		Name:                    "test",
		PaaSServiceTemplateUuid: "test-template",
		Labels:                  []string{"label"},
		PaaSSecurityZoneUuid:    "test-security-zone-id",
		ResourceLimits: []ResourceLimit{
			{
				Resource: "cpu",
				Limit:    2,
			},
		},
		Parameters: nil,
	})
	if err != nil {
		t.Errorf("CreatePaaSService returned error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", expectedRespObj), fmt.Sprintf("%v", *response))
}

func TestClient_UpdatePaaSService(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services", paasServiceUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPatch)
		fmt.Fprintf(w, "")
	})
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	err := client.UpdatePaaSService(paasServiceUuid, PaaSServiceUpdateRequest{
		Name:       "test",
		Labels:     []string{"label"},
		Parameters: parameters,
		ResourceLimits: []ResourceLimit{
			{
				Resource: "cpu",
				Limit:    2,
			},
		},
	})
	if err != nil {
		t.Errorf("UpdatePaaSService returned an error %v", err)
	}
}

func TestClient_DeletePaaSService(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services", paasServiceUuid)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodDelete)
		fmt.Fprintf(w, "")
	})
	err := client.DeletePaaSService(paasServiceUuid)
	if err != nil {
		t.Errorf("DeletePaaSService returned an error %v", err)
	}
}

func TestClient_GetPaaSServiceMetrics(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services", paasServiceUuid, "metrics")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		fmt.Fprintf(writer, preparePaaSHTTPGetMetricsResponse())
	})
	res, err := client.GetPaaSServiceMetrics(paasServiceUuid)
	if err != nil {
		t.Errorf("GetPaaSServiceMetrics returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockPaaSServiceMetric()), fmt.Sprintf("%v", res))
}

func TestClient_GetPaaSTemplateList(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "service_templates")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		fmt.Fprintf(writer, preparePaaSHTTPGetTemplatesResponse())
	})
	res, err := client.GetPaaSTemplateList()
	if err != nil {
		t.Errorf("GetPaaSTemplateList returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockPaasTemplate()), fmt.Sprintf("%v", res))
}

func TestClient_GetSecurityZoneList(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "security_zones")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		fmt.Fprintf(writer, preparePaaSHTTPGetSecurityZoneList())
	})
	res, err := client.GetPaaSSecurityZoneList()

	if err != nil {
		t.Errorf("GetPaaSSecurityZone returned an error %v", err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockSecurityZone()), fmt.Sprintf("%v", res))
}

func TestClient_CreatePaaSSecurityZone(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "security_zone")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodPost)
		fmt.Fprintf(writer, preparePaaSHTTPCreateSecurityZone())
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, requestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	res, err := client.CreatePaaSSecurityZone(PaaSSecurityZoneCreateRequest{
		Name:         "test",
		LocationUuid: "aa-bb-cc",
	})
	if err != nil {
		t.Errorf("CreatePaaSSecurityZone returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockPaaSSecurityZoneCreateResponse()), fmt.Sprintf("%v", *res))
}

func TestClient_GetPaaSSecurityZone(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "security_zones", securityZoneUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		fmt.Fprintf(writer, preparePaaSHTTPGetSecurityZone())
	})
	res, err := client.GetPaaSSecurityZone(securityZoneUuid)
	if err != nil {
		t.Errorf("GetPaaSSecurityZone returned an error %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", getMockSecurityZone()), fmt.Sprintf("%s", res))
}

func TestClient_UpdatePaaSSecurityZone(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "security_zones", securityZoneUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodPatch)
		fmt.Fprint(writer, "")
	})
	err := client.UpdatePaaSSecurityZone(securityZoneUuid, PaaSSecurityZoneUpdateRequest{
		Name:                 "test",
		LocationUuid:         "a-b-c",
		PaaSSecurityZoneUuid: securityZoneUuid,
	})
	if err != nil {
		t.Errorf("UpdatePaaSSecurityZone returned an error %v", err)
	}
}

func TestClient_DeletePaaSSecurityZone(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "security_zones", securityZoneUuid)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodDelete)
		fmt.Fprint(writer, "")
	})
	err := client.DeletePaaSSecurityZone(securityZoneUuid)
	if err != nil {
		t.Errorf("DeletePaaSSecurityZone returned an error %v", err)
	}
}

func getMockPaaSService() PaaSService {
	listenPort := make(map[string]map[string]int)
	portmap := make(map[string]int)
	portmap["mysql"] = 3306
	listenPort["fcfc::1:305e:6eff:fe62:4503"] = portmap
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	mock := PaaSService{
		Properties: PaaSServiceProperties{
			ObjectUuid: paasServiceUuid,
			Labels:     []string{"label"},
			Credentials: []Credential{
				{
					Username: "username",
					Password: "password",
					Type:     "type",
				},
			},
			CreateTime:          "2018-04-28T09:47:41Z",
			ListenPorts:         listenPort,
			SecurityZoneUuid:    "d711fc50-ad96-4070-b769-6fe2bf93792c",
			ServiceTemplateUuid: "504e2d11-7255-4712-b744-fcb093a4e613",
			UsageInMinutes:      999,
			CurrentPrice:        5.789,
			ChangeTime:          "2018-04-29T09:47:41Z",
			Status:              "active",
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
		BeginTime:       "2018-04-28T09:47:41Z",
		EndTime:         "2018-04-28T09:47:41Z",
		PaaSServiceUuid: paasServiceUuid,
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

func preparePaaSHTTPGetListResponse() string {
	paas := getMockPaaSService()
	res, _ := json.Marshal(paas.Properties)
	return fmt.Sprintf(`{"paas_services": {"%s" : %s}}`, paasServiceUuid, string(res))
}

func preparePaaSHTTPGetResponse() string {
	paas := getMockPaaSService()
	res, _ := json.Marshal(paas)
	return string(res)
}

func preparePaaSHTTPCreateResponse() string {
	paasCreateResponse := getMockPaaSServiceCreateResponse()
	res, _ := json.Marshal(paasCreateResponse)
	return string(res)
}

func getMockPaaSServiceCreateResponse() PaaSServiceCreateResponse {
	listenPort := make(map[string]map[string]string)
	portmap := make(map[string]string)
	portmap["mysql"] = "3306"
	portmap["http"] = "80"
	listenPort["fcfc::1:aaaa:bbbb:cccc:dddd"] = portmap
	parameters := make(map[string]interface{})
	parameters["TEST_PARAM"] = "param value"
	return PaaSServiceCreateResponse{
		RequestUuid:     requestUUID,
		ListenPorts:     listenPort,
		PaaSServiceUuid: paasServiceUuid,
		Credentials: []Credential{
			{
				Username: "username",
				Password: "password",
				Type:     "type",
			},
		},
		ObjectUuid: paasServiceUuid,
		ResourceLimits: []ResourceLimit{
			{
				Resource: "cpu",
				Limit:    2,
			},
		},
		Parameters: parameters,
	}
}

func getMockPaasTemplate() PaasTemplate {
	mock := PaasTemplate{Properties: PaaSTemplateProperties{
		Name:       "test",
		ObjectUuid: "d711fc50-ad96-4070-b769-6fe2bf93792c",
		Category:   "database",
		ProductNo:  0,
		Labels:     []string{"label"},
		Resources: []Resource{
			{
				Memory:      10,
				Connections: 10,
			},
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

func getMockSecurityZone() PaaSSecurityZone {
	mock := PaaSSecurityZone{Properties: PaaSSecurityZoneProperties{
		LocationCountry: "Germany",
		CreateTime:      "2018-04-28T09:47:41Z",
		LocationIata:    "none",
		ObjectUuid:      "aa-bb-cc-dd",
		Labels:          []string{"label"},
		LocationName:    "Bonn",
		Status:          "active",
		LocationUuid:    "cc-dd-ee",
		ChangeTime:      "2018-04-28T09:47:41Z",
		Name:            "test",
		Relation:        PaaSRelationService{Services: []ServiceObject{{ObjectUuid: "ff-gg-hh"}}},
	}}
	return mock
}

func preparePaaSHTTPGetSecurityZoneList() string {
	zone := getMockSecurityZone()
	res, _ := json.Marshal(zone.Properties)
	return fmt.Sprintf(`{"paas_security_zones": {"%s": %s}}`, "test", string(res))
}

func preparePaaSHTTPGetSecurityZone() string {
	zone := getMockSecurityZone()
	res, _ := json.Marshal(zone)
	return string(res)
}

func getMockPaaSSecurityZoneCreateResponse() PaaSSecurityZoneCreateResponse {
	return PaaSSecurityZoneCreateResponse{
		RequestUuid:          requestUUID,
		PaaSSecurityZoneUuid: securityZoneUuid,
		ObjectUuid:           securityZoneUuid,
	}
}

func preparePaaSHTTPCreateSecurityZone() string {
	res, _ := json.Marshal(getMockPaaSSecurityZoneCreateResponse())
	return string(res)
}

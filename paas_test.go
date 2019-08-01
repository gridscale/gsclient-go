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
	securityZoneUuid = ""
)

func TestGetPaaSServiceList(t *testing.T) {
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

func TestGetPaaSService(t *testing.T) {
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

func TestCreatePaaSService(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiPaaSBase, "services")
	expectedRespObj := getMockPaaSObjectRespone()
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

func preparePaaSHTTPGetListResponse() string {
	paas := getMockPaaSService()
	res, _ := json.Marshal(paas.Properties)
	return fmt.Sprintf(`{"paas_services": {"%s" : %s}}`, paasServiceUuid, string(res))
}

func preparePaaSHTTPGetResponse() string {
	paas := getMockPaaSService()
	res, _ := json.Marshal(paas.Properties)
	return fmt.Sprintf(`{"paas_service": %s}`, string(res))
}

func preparePaaSHTTPCreateResponse() string {
	paasCreateResponse := getMockPaaSObjectRespone()
	res, _ := json.Marshal(paasCreateResponse)
	return string(res)
}

func getMockPaaSObjectRespone() PaaSServiceCreateResponse {
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

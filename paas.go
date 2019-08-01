package gsclient

//PaaSServices is the JSON struct of a list of PaaS services
type PaaSServices struct {
	List map[string]PaaSServiceProperties `json:"paas_services"`
}

//PaaSService is the JSON struct of a single PaaS service
type PaaSService struct {
	Properties PaaSServiceProperties `json:"paas_service"`
}

//PaaSServiceProperties is the properties of a single PaaS service
type PaaSServiceProperties struct {
	ObjectUuid          string                    `json:"object_uuid"`
	Labels              []string                  `json:"labels"`
	Credentials         []Credential              `json:"credentials"`
	CreateTime          string                    `json:"create_time"`
	ListenPorts         map[string]map[string]int `json:"listen_ports"`
	SecurityZoneUuid    string                    `json:"security_zone_uuid"`
	ServiceTemplateUuid string                    `json:"service_template_uuid"`
	UsageInMinutes      int                       `json:"usage_in_minutes"`
	//UsageInMinutesStorage int                       `json:"usage_in_minutes_storage"`
	//UsageInMinutesCores   int                       `json:"usage_in_minutes_cores"`
	CurrentPrice   float64                `json:"current_price"`
	ChangeTime     string                 `json:"change_time"`
	Status         string                 `json:"status"`
	Name           string                 `json:"name"`
	ResourceLimits []ResourceLimit        `json:"resource_limits"`
	Parameters     map[string]interface{} `json:"parameters"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type PaaSServiceCreateRequest struct {
	Name                    string                 `json:"name"`
	PaaSServiceTemplateUuid string                 `json:"paas_service_template_uuid"`
	Labels                  []string               `json:"labels"`
	PaaSSecurityZoneUuid    string                 `json:"paas_security_zone_uuid"`
	ResourceLimits          []ResourceLimit        `json:"resource_limits"`
	Parameters              map[string]interface{} `json:"parameters"`
}

type ResourceLimit struct {
	Resource string `json:"resource"`
	Limit    int    `json:"limit"`
}

type PaaSServiceCreateResponse struct {
	RequestUuid     string                 `json:"request_uuid"`
	ListenPorts     map[string]string      `json:"listen_ports"`
	PaaSServiceUuid string                 `json:"paas_service_uuid"`
	Credentials     []Credential           `json:"credentials"`
	ObjectUuid      string                 `json:"object_uuid"`
	ResourceLimits  []ResourceLimit        `json:"resource_limits"`
	Parameters      map[string]interface{} `json:"parameters"`
}

type PaaSTemplates struct {
	List map[string]PaaSTemplateProperties `json:"paas_service_templates"`
}

type PaasTemplate struct {
	Properties PaaSTemplateProperties `json:"paas_service_template"`
}

type PaaSTemplateProperties struct {
	Name             string               `json:"name"`
	ObjectUuid       string               `json:"object_uuid"`
	Category         string               `json:"category"`
	ProductNo        int                  `json:"product_no"`
	Labels           []string             `json:"labels"`
	Resources        []Resource           `json:"resources"`
	Status           string               `json:"status"`
	ParametersSchema map[string]Parameter `json:"parameters_schema"`
}

type Parameter struct {
	Required    bool        `json:"required"`
	Empty       bool        `json:"empty"`
	Description string      `json:"description"`
	Max         int         `json:"max"`
	Min         int         `json:"min"`
	Default     interface{} `json:"default"`
	Type        string      `json:"type"`
	Allowed     []string    `json:"allowed"`
	Regex       string      `json:"regex"`
}

type Resource struct {
	Memory      int `json:"memory"`
	Connections int `json:"connections"`
}

type PaaSServiceUpdateRequest struct {
	Name           string                 `json:"name"`
	Labels         []string               `json:"labels"`
	Parameters     map[string]interface{} `json:"parameters"`
	ResourceLimits []ResourceLimit        `json:"resource_limits"`
}

type PaaSServiceMetrics struct {
	List []PaaSMetricProperties `json:"paas_service_metrics"`
}

type PaaSServiceMetric struct {
	Properties PaaSMetricProperties `json:"paas_service_metrics"`
}

type PaaSMetricProperties struct {
	BeginTime       string          `json:"begin_time"`
	EndTime         string          `json:"end_time"`
	PaaSServiceUuid string          `json:"paas_service_uuid"`
	CoreUsage       PaaSMetricValue `json:"core_usage"`
	StorageSize     PaaSMetricValue `json:"storage_size"`
}

type PaaSMetricValue struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type PaaSSecurityZones struct {
	List map[string]PaaSSecurityZoneProperties `json:"paas_security_zones"`
}

type PaaSSecurityZone struct {
	Properties PaaSSecurityZoneProperties `json:"paas_security_zone"`
}

type PaaSSecurityZoneProperties struct {
	LocationCountry string              `json:"location_country"`
	CreateTime      string              `json:"create_time"`
	LocationIata    string              `json:"location_iata"`
	ObjectUuid      string              `json:"object_uuid"`
	Labels          []string            `json:"labels"`
	LocationName    string              `json:"location_name"`
	Status          string              `json:"status"`
	LocationUuid    string              `json:"location_uuid"`
	ChangeTime      string              `json:"change_time"`
	Name            string              `json:"name"`
	Relation        PaaSRelationService `json:"relation"`
}

type PaaSRelationService struct {
	Services []ServiceObject `json:"services"`
}

type ServiceObject struct {
	ObjectUuid string `json:"object_uuid"`
}

type PaaSSecurityZoneCreateRequest struct {
	Name         string `json:"name"`
	LocationUuid string `json:"location_uuid"`
}

type PaaSSecurityZoneCreateResponse struct {
	RequestUuid          string `json:"request_uuid"`
	PaaSSecurityZoneUuid string `json:"paas_security_zone_uuid"`
	ObjectUuid           string `json:"object_uuid"`
}

type PaaSSecurityZoneUpdateRequest struct {
	Name                 string `json:"name"`
	LocationUuid         string `json:"location_uuid"`
	PaaSSecurityZoneUuid string `json:"paas_security_zone_uuid"`
}

//GetPaaSServiceList returns a list of PaaS Services
func (c *Client) GetPaaSServiceList() ([]PaaSService, error) {
	r := Request{
		uri:    apiPaaSBase + "/services",
		method: "GET",
	}

	response := new(PaaSServices)
	err := r.execute(*c, &response)

	list := []PaaSService{}
	for _, properties := range response.List {
		list = append(list, PaaSService{
			Properties: properties,
		})
	}

	return list, err
}

//CreatePaaSService creates a new PaaS service
func (c *Client) CreatePaaSService(body PaaSServiceCreateRequest) (*PaaSServiceCreateResponse, error) {
	r := Request{
		uri:    apiPaaSBase + "/services",
		method: "POST",
		body:   body,
	}

	response := new(PaaSServiceCreateResponse)
	err := r.execute(*c, response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

//GetPaaSService returns a specific PaaS Service based on given id
func (c *Client) GetPaaSService(id string) (PaaSService, error) {
	r := Request{
		uri:    apiPaaSBase + "/services/" + id,
		method: "GET",
	}

	response := PaaSService{}
	err := r.execute(*c, &response)
	return response, err
}

//UpdatePaaSService updates a specific PaaS Service based on a given id
func (c *Client) UpdatePaaSService(id string, body PaaSServiceUpdateRequest) error {
	r := Request{
		uri:    apiPaaSBase + "/services/" + id,
		method: "PATCH",
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeletePaaSService deletes a PaaS service
func (c *Client) DeletePaaSService(id string) error {
	r := Request{
		uri:    apiPaaSBase + "/services/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

//GetPaaSServiceMetrics get a specific PaaS Service's metrics based on a given id
func (c *Client) GetPaaSServiceMetrics(id string) ([]PaaSServiceMetric, error) {
	r := Request{
		uri:    apiPaaSBase + "/services/" + id + "/metrics",
		method: "GET",
	}
	response := new(PaaSServiceMetrics)
	err := r.execute(*c, response)
	list := []PaaSServiceMetric{}
	for _, properties := range response.List {
		list = append(list, PaaSServiceMetric{
			Properties: properties,
		})
	}
	return list, err
}

//GetPaaSTemplateList returns a list of PaaS service templates
func (c *Client) GetPaaSTemplateList() ([]PaasTemplate, error) {
	r := Request{
		uri:    apiPaaSBase + "/service_templates",
		method: "GET",
	}

	response := new(PaaSTemplates)
	err := r.execute(*c, response)

	list := []PaasTemplate{}
	for _, properties := range response.List {
		paasTemplate := PaasTemplate{
			Properties: properties,
		}
		list = append(list, paasTemplate)
	}

	return list, err
}

//GetSecurityZones get available security zones
func (c *Client) GetSecurityZoneList() ([]PaaSSecurityZone, error) {
	r := Request{
		uri:    apiPaaSBase + "/security_zones",
		method: "GET",
	}
	response := new(PaaSSecurityZones)
	err := r.execute(*c, response)
	list := []PaaSSecurityZone{}
	for _, properties := range response.List {
		list = append(list, PaaSSecurityZone{
			Properties: properties,
		})
	}
	return list, err
}

//CreateSecurityZone creates a new PaaS security zone
func (c *Client) CreatePaaSSecurityZone(body PaaSSecurityZoneCreateRequest) (*PaaSSecurityZoneCreateResponse, error) {
	r := Request{
		uri:    apiPaaSBase + "/security_zone",
		method: "POST",
		body:   body,
	}
	response := new(PaaSSecurityZoneCreateResponse)
	err := r.execute(*c, response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//GetSecurityZone get a specific PaaS Security Zone based on given id
func (c *Client) GetPaaSSecurityZone(id string) (PaaSSecurityZone, error) {
	r := Request{
		uri:    apiPaaSBase + "/security_zones/" + id,
		method: "GET",
	}
	response := PaaSSecurityZone{}
	err := r.execute(*c, &response)
	return response, err
}

//UpdatePaaSSecurityZone update a specific PaaS security zone based on given id
func (c *Client) UpdatePaaSSecurityZone(id string, body PaaSSecurityZoneUpdateRequest) error {
	r := Request{
		uri:    apiPaaSBase + "/security_zones/" + id,
		method: "PATCH",
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeletePaaSSecurityZone delete a specific PaaS Security Zone based on given id
func (c *Client) DeletePaaSSecurityZone(id string) error {
	r := Request{
		uri:    apiPaaSBase + "/security_zones/" + id,
		method: "DELETE",
	}
	return r.execute(*c, nil)
}

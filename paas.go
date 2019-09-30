package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

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
	ObjectUUID          string                    `json:"object_uuid"`
	Labels              []string                  `json:"labels"`
	Credentials         []Credential              `json:"credentials"`
	CreateTime          string                    `json:"create_time"`
	ListenPorts         map[string]map[string]int `json:"listen_ports"`
	SecurityZoneUUID    string                    `json:"security_zone_uuid"`
	ServiceTemplateUUID string                    `json:"service_template_uuid"`
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

//Credential is JSON struct of credential
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

//PaaSServiceCreateRequest is JSON struct of a request for creating a PaaS service
type PaaSServiceCreateRequest struct {
	Name                    string                 `json:"name"`
	PaaSServiceTemplateUUID string                 `json:"paas_service_template_uuid"`
	Labels                  []string               `json:"labels,omitempty"`
	PaaSSecurityZoneUUID    string                 `json:"paas_security_zone_uuid,omitempty"`
	ResourceLimits          []ResourceLimit        `json:"resource_limits,omitempty"`
	Parameters              map[string]interface{} `json:"parameters,omitempty"`
}

//ResourceLimit is JSON struct of resource limit
type ResourceLimit struct {
	Resource string `json:"resource"`
	Limit    int    `json:"limit"`
}

//PaaSServiceCreateResponse is JSON struct of a response for creating a PaaS service
type PaaSServiceCreateResponse struct {
	RequestUUID     string                       `json:"request_uuid"`
	ListenPorts     map[string]map[string]string `json:"listen_ports"`
	PaaSServiceUUID string                       `json:"paas_service_uuid"`
	Credentials     []Credential                 `json:"credentials"`
	ObjectUUID      string                       `json:"object_uuid"`
	ResourceLimits  []ResourceLimit              `json:"resource_limits"`
	Parameters      map[string]interface{}       `json:"parameters"`
}

//PaaSTemplates is JSON struct of a list of PaaS Templates
type PaaSTemplates struct {
	List map[string]PaaSTemplateProperties `json:"paas_service_templates"`
}

//PaaSTemplate is JSON struct for a single PaaS Template
type PaaSTemplate struct {
	Properties PaaSTemplateProperties `json:"paas_service_template"`
}

//PaaSTemplateProperties is JSOn struct of properties of a PaaS template
type PaaSTemplateProperties struct {
	Name             string               `json:"name"`
	ObjectUUID       string               `json:"object_uuid"`
	Category         string               `json:"category"`
	ProductNo        int                  `json:"product_no"`
	Labels           []string             `json:"labels"`
	Resources        []Resource           `json:"resources"`
	Status           string               `json:"status"`
	ParametersSchema map[string]Parameter `json:"parameters_schema"`
}

//Parameter JSON of a parameter
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

//Resource JSON of a resource
type Resource struct {
	Memory      int `json:"memory"`
	Connections int `json:"connections"`
}

//PaaSServiceUpdateRequest JSON of a request for updating a PaaS service
type PaaSServiceUpdateRequest struct {
	Name           string                 `json:"name,omitempty"`
	Labels         []string               `json:"labels,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	ResourceLimits []ResourceLimit        `json:"resource_limits,omitempty"`
}

//PaaSServiceMetrics JSON of a list of PaaS metrics
type PaaSServiceMetrics struct {
	List []PaaSMetricProperties `json:"paas_service_metrics"`
}

//PaaSServiceMetric JSON of a single PaaS metric
type PaaSServiceMetric struct {
	Properties PaaSMetricProperties `json:"paas_service_metric"`
}

//PaaSMetricProperties JSON of properties of a PaaS metric
type PaaSMetricProperties struct {
	BeginTime       string          `json:"begin_time"`
	EndTime         string          `json:"end_time"`
	PaaSServiceUUID string          `json:"paas_service_uuid"`
	CoreUsage       PaaSMetricValue `json:"core_usage"`
	StorageSize     PaaSMetricValue `json:"storage_size"`
}

//PaaSMetricValue JSON of a metric value
type PaaSMetricValue struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

//PaaSSecurityZones JSON struct of a list of PaaS security zones
type PaaSSecurityZones struct {
	List map[string]PaaSSecurityZoneProperties `json:"paas_security_zones"`
}

//PaaSSecurityZone JSON struct of a single PaaS security zone
type PaaSSecurityZone struct {
	Properties PaaSSecurityZoneProperties `json:"paas_security_zone"`
}

//PaaSSecurityZoneProperties JSOn struct of properties of a PaaS security zone
type PaaSSecurityZoneProperties struct {
	LocationCountry string              `json:"location_country"`
	CreateTime      string              `json:"create_time"`
	LocationIata    string              `json:"location_iata"`
	ObjectUUID      string              `json:"object_uuid"`
	Labels          []string            `json:"labels"`
	LocationName    string              `json:"location_name"`
	Status          string              `json:"status"`
	LocationUUID    string              `json:"location_uuid"`
	ChangeTime      string              `json:"change_time"`
	Name            string              `json:"name"`
	Relation        PaaSRelationService `json:"relation"`
}

//PaaSRelationService JSON struct of a relation between a PaaS service and a service
type PaaSRelationService struct {
	Services []ServiceObject `json:"services"`
}

//ServiceObject JSON struct of a service object
type ServiceObject struct {
	ObjectUUID string `json:"object_uuid"`
}

//PaaSSecurityZoneCreateRequest JSON struct of a request for creating a PaaS security zone
type PaaSSecurityZoneCreateRequest struct {
	Name         string `json:"name,omitempty"`
	LocationUUID string `json:"location_uuid,omitempty"`
}

//PaaSSecurityZoneCreateResponse JSON struct of a response for creating a PaaS security zone
type PaaSSecurityZoneCreateResponse struct {
	RequestUUID          string `json:"request_uuid"`
	PaaSSecurityZoneUUID string `json:"paas_security_zone_uuid"`
	ObjectUUID           string `json:"object_uuid"`
}

//PaaSSecurityZoneUpdateRequest JSON struct of a request for updating a PaaS security zone
type PaaSSecurityZoneUpdateRequest struct {
	Name                 string `json:"name,omitempty"`
	LocationUUID         string `json:"location_uuid,omitempty"`
	PaaSSecurityZoneUUID string `json:"paas_security_zone_uuid,omitempty"`
}

//GetPaaSServiceList returns a list of PaaS Services
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServices
func (c *Client) GetPaaSServiceList(ctx context.Context) ([]PaaSService, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "services"),
		method: http.MethodGet,
	}
	var response PaaSServices
	var paasServices []PaaSService
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		paasServices = append(paasServices, PaaSService{
			Properties: properties,
		})
	}
	return paasServices, err
}

//CreatePaaSService creates a new PaaS service
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createPaasService
func (c *Client) CreatePaaSService(ctx context.Context, body PaaSServiceCreateRequest) (PaaSServiceCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "services"),
		method: http.MethodPost,
		body:   body,
	}
	var response PaaSServiceCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return PaaSServiceCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
	return response, err
}

//GetPaaSService returns a specific PaaS Service based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasService
func (c *Client) GetPaaSService(ctx context.Context, id string) (PaaSService, error) {
	if !isValidUUID(id) {
		return PaaSService{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodGet,
	}
	var response PaaSService
	err := r.execute(ctx, *c, &response)
	return response, err
}

//UpdatePaaSService updates a specific PaaS Service based on a given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updatePaasService
func (c *Client) UpdatePaaSService(ctx context.Context, id string, body PaaSServiceUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForPaaSServiceActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//DeletePaaSService deletes a PaaS service
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deletePaasService
func (c *Client) DeletePaaSService(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForPaaSServiceDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//GetPaaSServiceMetrics get a specific PaaS Service's metrics based on a given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServiceMetrics
func (c *Client) GetPaaSServiceMetrics(ctx context.Context, id string) ([]PaaSServiceMetric, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id, "metrics"),
		method: http.MethodGet,
	}
	var response PaaSServiceMetrics
	var metrics []PaaSServiceMetric
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		metrics = append(metrics, PaaSServiceMetric{
			Properties: properties,
		})
	}
	return metrics, err
}

//GetPaaSTemplateList returns a list of PaaS service templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServiceTemplates
func (c *Client) GetPaaSTemplateList(ctx context.Context) ([]PaaSTemplate, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "service_templates"),
		method: http.MethodGet,
	}
	var response PaaSTemplates
	var paasTemplates []PaaSTemplate
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		paasTemplate := PaaSTemplate{
			Properties: properties,
		}
		paasTemplates = append(paasTemplates, paasTemplate)
	}
	return paasTemplates, err
}

//GetPaaSSecurityZoneList get available security zones
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasSecurityZones
func (c *Client) GetPaaSSecurityZoneList(ctx context.Context) ([]PaaSSecurityZone, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones"),
		method: http.MethodGet,
	}
	var response PaaSSecurityZones
	var securityZones []PaaSSecurityZone
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		securityZones = append(securityZones, PaaSSecurityZone{
			Properties: properties,
		})
	}
	return securityZones, err
}

//CreatePaaSSecurityZone creates a new PaaS security zone
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createPaasSecurityZone
func (c *Client) CreatePaaSSecurityZone(ctx context.Context, body PaaSSecurityZoneCreateRequest) (PaaSSecurityZoneCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones"),
		method: http.MethodPost,
		body:   body,
	}
	var response PaaSSecurityZoneCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return PaaSSecurityZoneCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
	return response, err
}

//GetPaaSSecurityZone get a specific PaaS Security Zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasSecurityZone
func (c *Client) GetPaaSSecurityZone(ctx context.Context, id string) (PaaSSecurityZone, error) {
	if !isValidUUID(id) {
		return PaaSSecurityZone{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodGet,
	}
	var response PaaSSecurityZone
	err := r.execute(ctx, *c, &response)
	return response, err
}

//UpdatePaaSSecurityZone update a specific PaaS security zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updatePaasSecurityZone
func (c *Client) UpdatePaaSSecurityZone(ctx context.Context, id string, body PaaSSecurityZoneUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSecurityZoneActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//DeletePaaSSecurityZone delete a specific PaaS Security Zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deletePaasSecurityZone
func (c *Client) DeletePaaSSecurityZone(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSecurityZoneDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//GetDeletedPaaSServices returns a list of deleted PaaS Services
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedPaasServices
func (c *Client) GetDeletedPaaSServices(ctx context.Context) ([]PaaSService, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "paas_services"),
		method: http.MethodGet,
	}
	var response DeletedPaaSServices
	var paasServices []PaaSService
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		paasServices = append(paasServices, PaaSService{
			Properties: properties,
		})
	}
	return paasServices, err
}

//waitForPaaSServiceActive allows to wait until the PaaS service's status is active
func (c *Client) waitForPaaSServiceActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		paas, err := c.GetPaaSService(ctx, id)
		return paas.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForPaaSServiceDeleted allows to wait until the PaaS service is deleted
func (c *Client) waitForPaaSServiceDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiPaaSBase, "services", id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}

//waitForSecurityZoneActive allows to wait until the security zone's status is active
func (c *Client) waitForSecurityZoneActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		secZone, err := c.GetPaaSSecurityZone(ctx, id)
		return secZone.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForSecurityZoneDeleted allows to wait until the security zone is deleted
func (c *Client) waitForSecurityZoneDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiPaaSBase, "security_zones", id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}

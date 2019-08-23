package gsclient

import (
	"net/http"
	"path"
)

//MarketplaceTemplateList JSON struct of a list of marketplace templates
type MarketplaceTemplateList struct {
	List map[string]MarketplaceTemplateProperties `json:"templates"`
}

//MarketplaceTemplate JSON struct of a single marketplace template
type MarketplaceTemplate struct {
	Properties MarketplaceTemplateProperties `json:"template"`
}

//MarketplaceTemplateProperties JSON struct of properties of a marketplace template
type MarketplaceTemplateProperties struct {
	Status           string                      `json:"status"`
	Ostype           string                      `json:"ostype"`
	LocationUUID     string                      `json:"location_uuid"`
	Version          string                      `json:"version"`
	LocationIata     string                      `json:"location_iata"`
	ChangeTime       string                      `json:"change_time"`
	Private          bool                        `json:"private"`
	ObjectUUID       string                      `json:"object_uuid"`
	LicenseProductNo int                         `json:"license_product_no"`
	CreateTime       string                      `json:"create_time"`
	UsageInMinutes   int                         `json:"usage_in_minutes"`
	Capacity         int                         `json:"capacity"`
	LocationName     string                      `json:"location_name"`
	Distro           string                      `json:"distro"`
	Description      string                      `json:"description"`
	CurrentPrice     float64                     `json:"current_price"`
	LocationCountry  string                      `json:"location_country"`
	Name             string                      `json:"name"`
	Labels           []string                    `json:"labels"`
	Metadata         MarketplaceTemplateMetadata `json:"metadata"`
}

//MarketplaceTemplateMetadata JSON struct of metadata of a marketplace template
type MarketplaceTemplateMetadata struct {
	OS               string                 `json:"os"`
	Top              bool                   `json:"top"`
	Icon             string                 `json:"icon"`
	Setup            map[string]interface{} `json:"setup"`
	License          string                 `json:"license"`
	Version          string                 `json:"version"`
	Category         string                 `json:"category"`
	Publisher        string                 `json:"publisher"`
	Description      map[string]string      `json:"description"`
	OtherSoftware    map[string]string      `json:"other_software"`
	ShortDescription map[string]string      `json:"short_description"`
}

//MarketplaceTemplateCreateRequest JSON struct of a request for creating a new marketplace template
type MarketplaceTemplateCreateRequest struct {
	Name              string   `json:"name"`
	Labels            []string `json:"labels,omitempty"`
	ObjectStoragePath string   `json:"object_storage_path"`
	Capacity          int      `json:"capacity"`
}

//MarketplaceTemplateCreateImportResponse JSON struct of a response of a marketplace template creation/importing
type MarketplaceTemplateCreateImportResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
	UniqueHash  string `json:"unique_hash"`
}

//MarketplaceTemplateImportRequest JSON struct of a request for importing  a marketplace template
type MarketplaceTemplateImportRequest struct {
	UniqueHash string `json:"unique_hash"`
}

//MarketplaceTemplateUpdateRequest JSON struct of a request for updating a marketplace template
type MarketplaceTemplateUpdateRequest struct {
	Name              string   `json:"name"`
	Labels            []string `json:"labels,omitempty"`
	ObjectStoragePath string   `json:"object_storage_path"`
	Capacity          int      `json:"capacity"`
	TemplateUUID      string   `json:"template_uuid"`
}

//MarketplaceTemplateEventList JSON struct of a list of events of a marketplace template
type MarketplaceTemplateEventList struct {
	List []MarketplaceTemplateEventProperties `json:"events"`
}

//MarketplaceTemplateEvent JSON struct of an event of a marketplace template
type MarketplaceTemplateEvent struct {
	Properties MarketplaceTemplateEventProperties `json:"event"`
}

//MarketplaceTemplateEventProperties JSON struct of properties of an event
type MarketplaceTemplateEventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUUID   string `json:"request_uuid"`
	ObjectUUID    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUUID      string `json:"user_uuid"`
}

//GetMarketplaceTemplateList gets a list of marketplace templates
func (c *Client) GetMarketplaceTemplateList() ([]MarketplaceTemplate, error) {
	r := Request{
		uri:    apiMarketplaceTemplateBase,
		method: http.MethodGet,
	}
	var response MarketplaceTemplateList
	var templates []MarketplaceTemplate
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, MarketplaceTemplate{Properties: properties})
	}
	return templates, err
}

//GetMarketplaceTemplate gets a specific marketplace template
func (c *Client) GetMarketplaceTemplate(id string) (MarketplaceTemplate, error) {
	r := Request{
		uri:    path.Join(apiMarketplaceTemplateBase, id),
		method: http.MethodGet,
	}
	var template MarketplaceTemplate
	err := r.execute(*c, &template)
	return template, err
}

//CreateMarketplaceTemplate creates a new marketplace template
func (c *Client) CreateMarketplaceTemplate(body MarketplaceTemplateCreateRequest) (
	MarketplaceTemplateCreateImportResponse, error) {
	r := Request{
		uri:    apiMarketplaceTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var response MarketplaceTemplateCreateImportResponse
	err := r.execute(*c, &response)
	if err != nil {
		return response, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//ImportMarketplaceTemplate imports a new marketplace template
func (c *Client) ImportMarketplaceTemplate(body MarketplaceTemplateImportRequest) (
	MarketplaceTemplateCreateImportResponse, error) {
	r := Request{
		uri:    apiMarketplaceTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var response MarketplaceTemplateCreateImportResponse
	err := r.execute(*c, &response)
	if err != nil {
		return response, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//UpdateMarketplaceTemplate updates a marketplace template
func (c *Client) UpdateMarketplaceTemplate(id string, body MarketplaceTemplateUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiMarketplaceTemplateBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteMarketplaceTemplate deletes a marketplace template
func (c *Client) DeleteMarketplaceTemplate(id string) error {
	r := Request{
		uri:    path.Join(apiMarketplaceTemplateBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetMarketplaceTemplateEventList gets a list of events of a marketplace template
func (c *Client) GetMarketplaceTemplateEventList(id string) ([]MarketplaceTemplateEvent, error) {
	r := Request{
		uri:    path.Join(apiMarketplaceTemplateBase, id, "events"),
		method: http.MethodGet,
	}
	var response MarketplaceTemplateEventList
	var events []MarketplaceTemplateEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		events = append(events, MarketplaceTemplateEvent{Properties: properties})
	}
	return events, err
}

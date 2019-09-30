package gsclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

//TemplateList JSON struct of a list of templates
type TemplateList struct {
	List map[string]TemplateProperties `json:"templates"`
}

//Template JSON struct of a single template
type Template struct {
	Properties TemplateProperties `json:"template"`
}

//TemplateProperties JSOn struct of properties of a template
type TemplateProperties struct {
	Status           string   `json:"status"`
	Ostype           string   `json:"ostype"`
	LocationUUID     string   `json:"location_uuid"`
	Version          string   `json:"version"`
	LocationIata     string   `json:"location_iata"`
	ChangeTime       string   `json:"change_time"`
	Private          bool     `json:"private"`
	ObjectUUID       string   `json:"object_uuid"`
	LicenseProductNo int      `json:"license_product_no"`
	CreateTime       string   `json:"create_time"`
	UsageInMinutes   int      `json:"usage_in_minutes"`
	Capacity         int      `json:"capacity"`
	LocationName     string   `json:"location_name"`
	Distro           string   `json:"distro"`
	Description      string   `json:"description"`
	CurrentPrice     float64  `json:"current_price"`
	LocationCountry  string   `json:"location_country"`
	Name             string   `json:"name"`
	Labels           []string `json:"labels"`
}

//TemplateEventList JSON struct of a list of a template's events
type TemplateEventList struct {
	List []TemplateEventProperties `json:"events"`
}

//TemplateEvent JSON struct of an event of a template
type TemplateEvent struct {
	Properties TemplateEventProperties `json:"event"`
}

//TemplateEventProperties JSON struct of properties of an event of a template
type TemplateEventProperties struct {
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

//TemplateCreateRequest JSON struct of a request for creating a template
type TemplateCreateRequest struct {
	Name         string   `json:"name"`
	SnapshotUUID string   `json:"snapshot_uuid"`
	Labels       []string `json:"labels,omitempty"`
}

//TemplateUpdateRequest JSON struct of a request for updating a template
type TemplateUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//GetTemplate gets a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplate
func (c *Client) GetTemplate(ctx context.Context, id string) (Template, error) {
	if !isValidUUID(id) {
		return Template{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodGet,
	}
	var response Template
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetTemplateList gets a list of templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplates
func (c *Client) GetTemplateList(ctx context.Context) ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{
			Properties: properties,
		})
	}
	return templates, err
}

//GetTemplateByName gets a template by its name
func (c *Client) GetTemplateByName(ctx context.Context, name string) (Template, error) {
	if name == "" {
		return Template{}, errors.New("'name' is required")
	}
	templates, err := c.GetTemplateList(ctx)
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return Template{Properties: template.Properties}, nil
		}
	}
	return Template{}, fmt.Errorf("Template %v not found", name)
}

//CreateTemplate creates a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createTemplate
func (c *Client) CreateTemplate(ctx context.Context, body TemplateCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return CreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
	return response, err
}

//UpdateTemplate updates a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateTemplate
func (c *Client) UpdateTemplate(ctx context.Context, id string, body TemplateUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForTemplateActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//DeleteTemplate deletes a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteTemplate
func (c *Client) DeleteTemplate(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForTemplateDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//GetTemplateEventList gets a list of a template's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplateEvents
func (c *Client) GetTemplateEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var templateEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		templateEvents = append(templateEvents, TemplateEvent{Properties: properties})
	}
	return templateEvents, err
}

//GetTemplatesByLocation gets a list of templates by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationTemplates
func (c *Client) GetTemplatesByLocation(ctx context.Context, id string) ([]Template, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "templates"),
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}

//GetDeletedTemplates gets a list of deleted templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedTemplates
func (c *Client) GetDeletedTemplates(ctx context.Context) ([]Template, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "templates"),
		method: http.MethodGet,
	}
	var response DeletedTemplateList
	var templates []Template
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}

//waitForTemplateActive allows to wait until the template's status is active
func (c *Client) waitForTemplateActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		template, err := c.GetTemplate(ctx, id)
		return template.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForTemplateDeleted allows to wait until the template is deleted
func (c *Client) waitForTemplateDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiTemplateBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}

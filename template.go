package gsclient

import (
	"fmt"
	"net/http"
	"path"
)

type TemplateList struct {
	List map[string]TemplateProperties `json:"templates"`
}

type Template struct {
	Properties TemplateProperties `json:"template"`
}

type TemplateProperties struct {
	Status           string   `json:"status"`
	Ostype           string   `json:"ostype"`
	LocationUuid     string   `json:"location_uuid"`
	Version          string   `json:"version"`
	LocationIata     string   `json:"location_iata"`
	ChangeTime       string   `json:"change_time"`
	Private          bool     `json:"private"`
	ObjectUuid       string   `json:"object_uuid"`
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

//GetTemplate gets a template
func (c *Client) GetTemplate(id string) (Template, error) {
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodGet,
	}
	var response Template
	err := r.execute(*c, &response)

	return response, err
}

//GetTemplateList gets a list of templates
func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{
			Properties: properties,
		})
	}
	return templates, err
}

//GetTemplateByName gets a template by its name
func (c *Client) GetTemplateByName(name string) (Template, error) {
	templates, err := c.GetTemplateList()
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return c.GetTemplate(template.Properties.ObjectUuid)
		}
	}

	return Template{}, fmt.Errorf("Template %v not found", name)
}

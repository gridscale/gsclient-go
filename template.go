package gsclient

import (
	"fmt"
)

type Templates struct {
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

func (c *Client) GetTemplate(id string) (*Template, error) {
	r := Request{
		uri:    apiTemplateBase + "/" + id,
		method: "GET",
	}

	response := new(Template)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: "GET",
	}

	response := new(Templates)
	err := r.execute(*c, &response)

	list := []Template{}
	for _, properties := range response.List {
		template := Template{
			Properties: properties,
		}
		list = append(list, template)
	}

	return list, err
}

func (c *Client) GetTemplateByName(name string) (*Template, error) {
	templates, err := c.GetTemplateList()
	if err != nil {
		return nil, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return c.GetTemplate(template.Properties.ObjectUuid)
		}
	}

	return nil, fmt.Errorf("Template %v not found", name)
}

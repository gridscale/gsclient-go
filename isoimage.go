package gsclient

import (
	"net/http"
	"path"
)

//ISOImageList is JSON struct of a list of ISO images
type ISOImageList struct {
	List map[string]ISOImageProperties `json:"isoimages"`
}

//ISOImage is JSON struct of a list an ISO image
type ISOImage struct {
	Properties ISOImageProperties `json:"isoimage"`
}

//ISOImageProperties is JSON struct of properties of an ISO image
type ISOImageProperties struct {
	ObjectUuid      string           `json:"object_uuid"`
	Relations       ISOImageRelation `json:"relations"`
	Description     string           `json:"description"`
	LocationName    string           `json:"location_name"`
	SourceUrl       string           `json:"source_url"`
	Labels          []string         `json:"labels"`
	LocationIata    string           `json:"location_iata"`
	LocationUuid    string           `json:"location_uuid"`
	Status          string           `json:"status"`
	CreateTime      string           `json:"create_time"`
	Name            string           `json:"name"`
	Version         string           `json:"version"`
	LocationCountry string           `json:"location_country"`
	UsageInMinutes  int              `json:"usage_in_minutes"`
	Private         bool             `json:"private"`
	ChangeTime      string           `json:"change_time"`
	Capacity        int              `json:"capacity"`
	CurrentPrice    float64          `json:"current_price"`
}

type ISOImageRelation struct {
	Servers []ServerinISOImage `json:"servers"`
}

type ServerinISOImage struct {
	Bootdevice bool   `json:"bootdevice"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
	ObjectUuid string `json:"object_uuid"`
}

type ISOImageCreateRequest struct {
	Name         string   `json:"name"`
	SourceUrl    string   `json:"source_url"`
	Labels       []string `json:"labels"`
	LocationUuid string   `json:"location_uuid"`
}

type ISOImageCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

type ISOImageUpdateRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

type ISOImageEventList struct {
	List []ISOImageEventProperties `json:"events"`
}

type ISOImageEvent struct {
	Properties ISOImageEventProperties `json:"event"`
}

type ISOImageEventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUuid   string `json:"request_uuid"`
	ObjectUuid    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUuid      string `json:"user_uuid"`
}

//GetISOImageList returns a list of available ISO images
func (c *Client) GetISOImageList() ([]ISOImage, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodGet,
	}
	response := ISOImageList{}
	list := []ISOImage{}
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ISOImage{Properties: properties})
	}
	return list, err
}

//GetISOImage returns a specific ISO image based on given id
func (c *Client) GetISOImage(id string) (ISOImage, error) {
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodGet,
	}
	response := ISOImage{}
	err := r.execute(*c, &response)
	return response, err
}

//CreateISOImage creates an ISO image
func (c *Client) CreateISOImage(body ISOImageCreateRequest) (ISOImageCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodPost,
		body:   body,
	}
	response := ISOImageCreateResponse{}
	err := r.execute(*c, &response)
	if err != nil {
		return ISOImageCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//UpdateISOImage updates a specific ISO Image
func (c *Client) UpdateISOImage(id string, body ISOImageUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteISOImage deletes a specific ISO image
func (c *Client) DeleteISOImage(id string) error {
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetISOImageEvents returns a list of events of an ISO image
func (c *Client) GetISOImageEventList(id string) ([]ISOImageEvent, error) {
	r := Request{
		uri:    path.Join(apiISOBase, id, "events"),
		method: http.MethodGet,
	}
	response := ISOImageEventList{}
	list := []ISOImageEvent{}
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ISOImageEvent{Properties: properties})
	}
	return list, err
}

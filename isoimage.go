package gsclient

import (
<<<<<<< HEAD
=======
	"context"
	"errors"
>>>>>>> 8d4aa0e... add `context`
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
	ObjectUUID      string           `json:"object_uuid"`
	Relations       ISOImageRelation `json:"relations"`
	Description     string           `json:"description"`
	LocationName    string           `json:"location_name"`
	SourceURL       string           `json:"source_url"`
	Labels          []string         `json:"labels"`
	LocationIata    string           `json:"location_iata"`
	LocationUUID    string           `json:"location_uuid"`
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

//ISOImageRelation is JSON struct of a list of an ISO-Image's relations
type ISOImageRelation struct {
	Servers []ServerinISOImage `json:"servers"`
}

//ServerinISOImage is JSON struct of a relation between an ISO-Image and a Server
type ServerinISOImage struct {
	Bootdevice bool   `json:"bootdevice"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
	ObjectUUID string `json:"object_uuid"`
}

//ISOImageCreateRequest is JSON struct of a request for creating an ISO-Image
type ISOImageCreateRequest struct {
	Name         string   `json:"name"`
	SourceURL    string   `json:"source_url"`
	Labels       []string `json:"labels,omitempty"`
	LocationUUID string   `json:"location_uuid"`
}

//ISOImageCreateResponse is JSON struct of a response for creating an ISO-Image
type ISOImageCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//ISOImageUpdateRequest is JSON struct of a request for updating an ISO-Image
type ISOImageUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//ISOImageEventList is JSON struct of a list of an ISO-Image's events
type ISOImageEventList struct {
	List []ISOImageEventProperties `json:"events"`
}

//ISOImageEvent is JSON struct of a single event of an ISO-Image
type ISOImageEvent struct {
	Properties ISOImageEventProperties `json:"event"`
}

//ISOImageEventProperties is JSON struct of an ISO-Image event
type ISOImageEventProperties struct {
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

//GetISOImageList returns a list of available ISO images
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIsoimages
func (c *Client) GetISOImageList(ctx context.Context) ([]ISOImage, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodGet,
	}
	var response ISOImageList
	var isoImages []ISOImage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}

//GetISOImage returns a specific ISO image based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIsoimage
func (c *Client) GetISOImage(ctx context.Context, id string) (ISOImage, error) {
	if !isValidUUID(id) {
		return ISOImage{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodGet,
	}
	var response ISOImage
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateISOImage creates an ISO image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createIsoimage
func (c *Client) CreateISOImage(ctx context.Context, body ISOImageCreateRequest) (ISOImageCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodPost,
		body:   body,
	}
	var response ISOImageCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return ISOImageCreateResponse{}, err
	}
<<<<<<< HEAD
	err = c.WaitForRequestCompletion(response.RequestUUID)
=======
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
>>>>>>> 8d4aa0e... add `context`
	return response, err
}

//UpdateISOImage updates a specific ISO Image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateIsoimage
func (c *Client) UpdateISOImage(ctx context.Context, id string, body ISOImageUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodPatch,
		body:   body,
	}
<<<<<<< HEAD
	return r.execute(*c, nil)
=======
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForISOImageActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//DeleteISOImage deletes a specific ISO image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteIsoimage
func (c *Client) DeleteISOImage(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodDelete,
	}
<<<<<<< HEAD
	return r.execute(*c, nil)
=======
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForISOImageDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetISOImageEventList returns a list of events of an ISO image
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIsoimageEvents
func (c *Client) GetISOImageEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var isoImageEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		isoImageEvents = append(isoImageEvents, ISOImageEvent{Properties: properties})
	}
	return isoImageEvents, err
}
<<<<<<< HEAD
=======

//GetISOImagesByLocation gets a list of ISO images by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationIsoimages
func (c *Client) GetISOImagesByLocation(ctx context.Context, id string) ([]ISOImage, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "isoimages"),
		method: http.MethodGet,
	}
	var response ISOImageList
	var isoImages []ISOImage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}

//GetDeletedISOImages gets a list of deleted ISO images
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedIsoimages
func (c *Client) GetDeletedISOImages(ctx context.Context) ([]ISOImage, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "isoimages"),
		method: http.MethodGet,
	}
	var response DeletedISOImageList
	var isoImages []ISOImage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}

//waitForISOImageActive allows to wait until the ISO-Image's status is active
func (c *Client) waitForISOImageActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		img, err := c.GetISOImage(ctx, id)
		return img.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForISOImageDeleted allows to wait until the ISO-Image id deleted
func (c *Client) waitForISOImageDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiISOBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

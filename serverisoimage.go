package gsclient

import (
	"net/http"
	"path"
)

type ServerIsoImageList struct {
	List []ServerIsoImage `json:"isoimage_relations"`
}

type ServerIsoImageSingle struct {
	Properties ServerIsoImage `json:"isoimage_relation"`
}

type ServerIsoImage struct {
	ObjectUuid string `json:"object_uuid"`
	ObjectName string `json:"object_name"`
	Private    bool   `json:"private"`
	CreateTime string `json:"create_time"`
	Bootdevice bool   `json:"bootdevice"`
}

type ServerIsoImageCreateRequest struct {
	ObjectUuid string `json:"object_uuid"`
}

type ServerIsoImageUpdateRequest struct {
	BootDevice bool   `json:"bootdevice"`
	Name       string `json:"name"`
}

//GetServerIsoImageList gets a list of a specific server's ISO images
func (c *Client) GetServerIsoImageList(id string) ([]ServerIsoImage, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodGet,
	}
	var response ServerIsoImageList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIsoImage gets an ISO image of a specific server
func (c *Client) GetServerIsoImage(serverId, isoImageId string) (ServerIsoImage, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "isoimages", isoImageId),
		method: http.MethodGet,
	}
	var response ServerIsoImageSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerIsoImage updates a link between a storage and an ISO image
func (c *Client) UpdateServerIsoImage(serverId, isoImageId string, body ServerIsoImageUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "isoimages", isoImageId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerIsoImage creates a link between a server and an ISO image
func (c *Client) CreateServerIsoImage(id string, body ServerIsoImageCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerIsoImage deletes a link between an ISO image and a server
func (c *Client) DeleteServerIsoImage(serverId, isoImageId string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "isoimages", isoImageId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkIsoimage attaches an ISO image to a server
func (c *Client) LinkIsoimage(serverid string, isoimageid string) error {
	body := ServerIsoImageCreateRequest{
		ObjectUuid:isoimageid,
	}
	return c.CreateServerIsoImage(serverid, body)
}

//UnlinkIsoimage removes the link between an ISO image and a server
func (c *Client) UnlinkIsoimage(serverid string, isoimageid string) error {
	return c.DeleteServerIsoImage(serverid, isoimageid)
}

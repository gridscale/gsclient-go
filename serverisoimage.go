package gsclient

import (
	"net/http"
	"path"
)

//ServerIsoImageRelationList JSON struct of a list of relations between a server and ISO-Images
type ServerIsoImageRelationList struct {
	List []ServerIsoImageRelationProperties `json:"isoimage_relations"`
}

//ServerIsoImageRelation JSON struct of a single relation between a server and an ISO-Image
type ServerIsoImageRelation struct {
	Properties ServerIsoImageRelationProperties `json:"isoimage_relation"`
}

//ServerIsoImageRelationProperties JSON struct of properties of a relation between a server and an ISO-Image
type ServerIsoImageRelationProperties struct {
	ObjectUuid string `json:"object_uuid"`
	ObjectName string `json:"object_name"`
	Private    bool   `json:"private"`
	CreateTime string `json:"create_time"`
	Bootdevice bool   `json:"bootdevice"`
}

//ServerIsoImageRelationCreateRequest JSON struct of a request for creating a relation between a server and an ISO-Image
type ServerIsoImageRelationCreateRequest struct {
	ObjectUuid string `json:"object_uuid"`
}

//ServerIsoImageRelationUpdateRequest JSON struct of a request for updating a relation between a server and an ISO-Image
type ServerIsoImageRelationUpdateRequest struct {
	BootDevice bool   `json:"bootdevice"`
	Name       string `json:"name"`
}

//GetServerIsoImageList gets a list of a specific server's ISO images
func (c *Client) GetServerIsoImageList(id string) ([]ServerIsoImageRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodGet,
	}
	var response ServerIsoImageRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIsoImage gets an ISO image of a specific server
func (c *Client) GetServerIsoImage(serverId, isoImageId string) (ServerIsoImageRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "isoimages", isoImageId),
		method: http.MethodGet,
	}
	var response ServerIsoImageRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerIsoImage updates a link between a storage and an ISO image
func (c *Client) UpdateServerIsoImage(serverId, isoImageId string, body ServerIsoImageRelationUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "isoimages", isoImageId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerIsoImage creates a link between a server and an ISO image
func (c *Client) CreateServerIsoImage(id string, body ServerIsoImageRelationCreateRequest) error {
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

//LinkIsoImage attaches an ISO image to a server
func (c *Client) LinkIsoImage(serverId string, isoimageId string) error {
	body := ServerIsoImageRelationCreateRequest{
		ObjectUuid: isoimageId,
	}
	return c.CreateServerIsoImage(serverId, body)
}

//UnlinkIsoImage removes the link between an ISO image and a server
func (c *Client) UnlinkIsoImage(serverId string, isoimageId string) error {
	return c.DeleteServerIsoImage(serverId, isoimageId)
}

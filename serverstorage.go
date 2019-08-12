package gsclient

import (
	"net/http"
	"path"
)

type ServerStorageList struct {
	List []ServerStorage `json:"storage_relations"`
}

type ServerStorageSingle struct {
	Properties ServerStorage `json:"storage_relation"`
}

type ServerStorage struct {
	ObjectUuid       string `json:"object_uuid"`
	ObjectName       string `json:"object_name"`
	Capacity         int    `json:"capacity"`
	StorageType      string `json:"storage_type"`
	Target           int    `json:"target"`
	Lun              int    `json:"lun"`
	Controller       int    `json:"controller"`
	CreateTime       string `json:"create_time"`
	BootDevice       bool   `json:"bootdevice"`
	Bus              int    `json:"bus"`
	LastUsedTemplate string `json:"last_used_template"`
	LicenseProductNo int    `json:"license_product_no"`
	ServerUuid       string `json:"server_uuid"`
}

type ServerStorageCreateRequest struct {
	ObjectUuid string `json:"object_uuid"`
	BootDevice bool   `json:"bootdevice"`
}

type ServerStorageUpdateRequest struct {
	Ordering   int      `json:"ordering"`
	BootDevice bool     `json:"bootdevice"`
	L3security []string `json:"l3security"`
}

//GetServerStorageList gets a list of a specific server's storages
func (c *Client) GetServerStorageList(id string) ([]ServerStorage, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodGet,
	}
	var response ServerStorageList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerStorage gets a storage of a specific server
func (c *Client) GetServerStorage(serverId, storageId string) (ServerStorage, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodGet,
	}
	var response ServerStorageSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerStorage updates a link between a storage and a server
func (c *Client) UpdateServerStorage(serverId, storageId string, body ServerStorageUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerStorage create a link between a server and a storage
func (c *Client) CreateServerStorage(id string, body ServerStorageCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerStorage delete a link between a storage and a server
func (c *Client) DeleteServerStorage(serverId, storageId string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkStorage attaches a storage to a server
func (c *Client) LinkStorage(serverid string, storageid string, bootdevice bool) error {
	body := ServerStorageCreateRequest{
		ObjectUuid: storageid,
		BootDevice: bootdevice,
	}
	return c.CreateServerStorage(serverid, body)
}

//UnlinkStorage remove a storage from a server
func (c *Client) UnlinkStorage(serverid string, storageid string) error {
	return c.DeleteServerStorage(serverid, storageid)
}

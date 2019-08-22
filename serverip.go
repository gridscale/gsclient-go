package gsclient

import (
	"net/http"
	"path"
)

//ServerIpRelationList JSON struct of a list of relations between a server and IP addresses
type ServerIpRelationList struct {
	List []ServerIpRelationProperties `json:"ip_relations"`
}

//ServerIpRelation JSON struct of a single relation between a server and a IP address
type ServerIpRelation struct {
	Properties ServerIpRelationProperties `json:"ip_relation"`
}

//ServerIpRelationProperties JSON struct of properties of a relation between a server and a IP address
type ServerIpRelationProperties struct {
	ServerUUID string `json:"server_uuid"`
	CreateTime string `json:"create_time"`
	Prefix     string `json:"prefix"`
	Family     int    `json:"family"`
	ObjectUUID string `json:"object_uuid"`
	Ip         string `json:"ip"`
}

//ServerIpRelationCreateRequest JSON struct of request for creating a relation between a server and a IP address
type ServerIpRelationCreateRequest struct {
	ObjectUUID string `json:"object_uuid"`
}

//GetServerIpList gets a list of a specific server's IPs
func (c *Client) GetServerIpList(id string) ([]ServerIpRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodGet,
	}
	var response ServerIpRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIp gets an IP of a specific server
func (c *Client) GetServerIp(serverId, ipId string) (ServerIpRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "ips", ipId),
		method: http.MethodGet,
	}
	var response ServerIpRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//CreateServerIp create a link between a server and an IP
func (c *Client) CreateServerIp(id string, body ServerIpRelationCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerIp delete a link between a server and an IP
func (c *Client) DeleteServerIp(serverId, ipID string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "ips", ipID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkIp attaches an IP to a server
func (c *Client) LinkIp(serverId string, ipID string) error {
	body := ServerIpRelationCreateRequest{
		ObjectUUID: ipID,
	}
	return c.CreateServerIp(serverId, body)
}

//UnlinkIp removes a link between an IP and a server
func (c *Client) UnlinkIp(serverId string, ipID string) error {
	return c.DeleteServerIp(serverId, ipID)
}

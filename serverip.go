package gsclient

import (
	"net/http"
	"path"
)

type ServerIpList struct {
	List []ServerIp `json:"ip_relations"`
}

type ServerIpSingle struct {
	Properties ServerIp `json:"ip_relation"`
}

type ServerIp struct {
	ServerUuid string `json:"server_uuid"`
	CreateTime string `json:"create_time"`
	Prefix     string `json:"prefix"`
	Family     int    `json:"family"`
	ObjectUuid string `json:"object_uuid"`
	Ip         string `json:"ip"`
}

type ServerIpCreateRequest struct {
	ObjectUuid string `json:"object_uuid"`
}


//GetServerIpList gets a list of a specific server's IPs
func (c *Client) GetServerIpList(id string) ([]ServerIp, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodGet,
	}
	var response ServerIpList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIp gets an IP of a specific server
func (c *Client) GetServerIp(serverId, ipId string) (ServerIp, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "ips", ipId),
		method: http.MethodGet,
	}
	var response ServerIpSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//CreateServerIp create a link between a server and an IP
func (c *Client) CreateServerIp(id string, body ServerIpCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerId delete a link between a storage and an IP
func (c *Client) DeleteServerId(serverId, ipId string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "ips", ipId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkIp attaches an IP to a server
func (c *Client) LinkIp(serverid string, ipid string) error {
	body := ServerIpCreateRequest{
		ObjectUuid:ipid,
	}
	return c.CreateServerIp(serverid, body)
}

//UnlinkIp removes a link between an IP and a server
func (c *Client) UnlinkIp(serverid string, ipid string) error {
	return c.DeleteServerId(serverid, ipid)
}

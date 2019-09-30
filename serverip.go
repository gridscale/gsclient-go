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

//ServerIPRelationList JSON struct of a list of relations between a server and IP addresses
type ServerIPRelationList struct {
	List []ServerIPRelationProperties `json:"ip_relations"`
}

//ServerIPRelation JSON struct of a single relation between a server and a IP address
type ServerIPRelation struct {
	Properties ServerIPRelationProperties `json:"ip_relation"`
}

//ServerIPRelationProperties JSON struct of properties of a relation between a server and a IP address
type ServerIPRelationProperties struct {
	ServerUUID string `json:"server_uuid"`
	CreateTime string `json:"create_time"`
	Prefix     string `json:"prefix"`
	Family     int    `json:"family"`
	ObjectUUID string `json:"object_uuid"`
	IP         string `json:"ip"`
}

//ServerIPRelationCreateRequest JSON struct of request for creating a relation between a server and a IP address
type ServerIPRelationCreateRequest struct {
	ObjectUUID string `json:"object_uuid"`
}

//GetServerIPList gets a list of a specific server's IPs
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIps
func (c *Client) GetServerIPList(ctx context.Context, id string) ([]ServerIPRelationProperties, error) {
	if id == "" {
		return nil, errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodGet,
	}
	var response ServerIPRelationList
	err := r.execute(ctx, *c, &response)
	return response.List, err
}

//GetServerIP gets an IP of a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedIp
func (c *Client) GetServerIP(ctx context.Context, serverID, ipID string) (ServerIPRelationProperties, error) {
	if serverID == "" || ipID == "" {
		return ServerIPRelationProperties{}, errors.New("'serverID' and 'ipID' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "ips", ipID),
		method: http.MethodGet,
	}
	var response ServerIPRelation
	err := r.execute(ctx, *c, &response)
	return response.Properties, err
}

//CreateServerIP create a link between a server and an IP
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/linkIpToServer
func (c *Client) CreateServerIP(ctx context.Context, id string, body ServerIPRelationCreateRequest) error {
	if id == "" || body.ObjectUUID == "" {
		return errors.New("'server_id' and 'ip_id' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodPost,
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
		return c.waitForServerIPRelCreation(ctx, id, body.ObjectUUID)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//DeleteServerIP delete a link between a server and an IP
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/unlinkIpFromServer
func (c *Client) DeleteServerIP(ctx context.Context, serverID, ipID string) error {
	if serverID == "" || ipID == "" {
		return errors.New("'serverID' and 'ipID' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "ips", ipID),
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
		return c.waitForServerIPRelDeleted(ctx, serverID, ipID)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//LinkIP attaches an IP to a server
func (c *Client) LinkIP(ctx context.Context, serverID string, ipID string) error {
	body := ServerIPRelationCreateRequest{
		ObjectUUID: ipID,
	}
	return c.CreateServerIP(ctx, serverID, body)
}

//UnlinkIP removes a link between an IP and a server
func (c *Client) UnlinkIP(ctx context.Context, serverID string, ipID string) error {
	return c.DeleteServerIP(ctx, serverID, ipID)
}
<<<<<<< HEAD
=======

//waitForServerIPRelCreation allows to wait until the relation between a server and an IP address is created
func (c *Client) waitForServerIPRelCreation(ctx context.Context, serverID, ipID string) error {
	if !isValidUUID(serverID) || !isValidUUID(ipID) {
		return errors.New("'serverID' and 'ipID' are required")
	}
	uri := path.Join(apiServerBase, serverID, "ips", ipID)
	method := http.MethodGet
	return c.waitFor200Status(ctx, uri, method)
}

//waitForServerIPRelDeleted allows to wait until the relation between a server and an IP address is deleted
func (c *Client) waitForServerIPRelDeleted(ctx context.Context, serverID, ipID string) error {
	if !isValidUUID(serverID) || !isValidUUID(ipID) {
		return errors.New("'serverID' and 'ipID' are required")
	}
	uri := path.Join(apiServerBase, serverID, "ips", ipID)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

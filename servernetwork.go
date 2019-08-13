package gsclient

import (
	"net/http"
	"path"
)

type ServerNetworkList struct {
	List []ServerNetwork `json:"network_relations"`
}

type ServerNetworkSingle struct {
	Properties ServerNetwork `json:"network_relation"`
}

type ServerNetwork struct {
	L2security           bool     `json:"l2security"`
	ServerUuid           string   `json:"server_uuid"`
	CreateTime           string   `json:"create_time"`
	PublicNet            bool     `json:"public_net"`
	FirewallTemplateUuid string   `json:"firewall_template_uuid,omitempty"`
	ObjectName           string   `json:"object_name"`
	Mac                  string   `json:"mac"`
	BootDevice           bool     `json:"bootdevice"`
	PartnerUuid          string   `json:"partner_uuid"`
	Ordering             int      `json:"ordering"`
	Firewall             string   `json:"firewall,omitempty"`
	NetworkType          string   `json:"network_type"`
	NetworkUuid          string   `json:"network_uuid"`
	ObjectUuid           string   `json:"object_uuid"`
	L3security           []string `json:"l3security"`
	//Vlan                 int          `json:"vlan,omitempty"`
	//Vxlan                int          `json:"vxlan,omitempty"`
	//Mcast                string       `json:"mcast, omitempty"`
}

type ServerNetworkCreateRequest struct {
	ObjectUuid           string        `json:"object_uuid"`
	Ordering             int           `json:"ordering"`
	BootDevice           bool          `json:"bootdevice"`
	L3security           []string      `json:"l3security"`
	Firewall             FirewallRules `json:"firewall"`
	FirewallTemplateUuid string        `json:"firewall_template_uuid"`
}

type ServerNetworkUpdateRequest struct {
	Ordering             int           `json:"ordering"`
	BootDevice           bool          `json:"bootdevice"`
	L3security           []string      `json:"l3security"`
	Firewall             FirewallRules `json:"firewall"`
	FirewallTemplateUuid string        `json:"firewall_template_uuid"`
}

//GetServerNetworkList gets a list of a specific server's networks
func (c *Client) GetServerNetworkList(id string) ([]ServerNetwork, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodGet,
	}
	var response ServerNetworkList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerNetwork gets a network of a specific server
func (c *Client) GetServerNetwork(serverId, networkId string) (ServerNetwork, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "networks", networkId),
		method: http.MethodGet,
	}
	var response ServerNetworkSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerNetwork updates a link between a network and a server
func (c *Client) UpdateServerNetwork(serverId, networkId string, body ServerNetworkUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "networks", networkId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerNetwork creates a link between a network and a storage
func (c *Client) CreateServerNetwork(id string, body ServerNetworkCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerNetwork deletes a link between a network and a server
func (c *Client) DeleteServerNetwork(serverId, networkId string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "networks", networkId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkNetwork attaches a network to a server
func (c *Client) LinkNetwork(serverid, networkid, firewallTemplate string, bootdevice bool, order int,
	l3security []string) error {
	body := ServerNetworkCreateRequest{
		ObjectUuid:           networkid,
		Ordering:             order,
		BootDevice:           bootdevice,
		L3security:           l3security,
		FirewallTemplateUuid: firewallTemplate,
		//Firewall:
	}
	return c.CreateServerNetwork(serverid, body)
}

//UnlinkNetwork removes the link between a network and a server
func (c *Client) UnlinkNetwork(serverid string, networkid string) error {
	return c.DeleteServerNetwork(serverid, networkid)
}

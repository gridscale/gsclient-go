package gsclient

import (
	"net/http"
	"path"
)

//ServerNetworkRelationList JSON struct of a list of relations between a server and networks
type ServerNetworkRelationList struct {
	List []ServerNetworkRelationProperties `json:"network_relations"`
}

//ServerNetworkRelation JSON struct of a single relation between a server and a network
type ServerNetworkRelation struct {
	Properties ServerNetworkRelationProperties `json:"network_relation"`
}

//ServerNetworkRelationProperties JSON struct of properties of a relation between a server and a network
type ServerNetworkRelationProperties struct {
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

//ServerNetworkRelationCreateRequest JSON struct of a request for creating a relation between a server and a network
type ServerNetworkRelationCreateRequest struct {
	ObjectUuid           string        `json:"object_uuid"`
	Ordering             int           `json:"ordering,omitempty"`
	BootDevice           bool          `json:"bootdevice,omitempty"`
	L3security           []string      `json:"l3security,omitempty"`
	Firewall             FirewallRules `json:"firewall,omitempty"`
	FirewallTemplateUuid string        `json:"firewall_template_uuid,omitempty"`
}

//ServerNetworkRelationUpdateRequest JSON struct of a request for updating a relation between a server and a network
type ServerNetworkRelationUpdateRequest struct {
	Ordering             int           `json:"ordering"`
	BootDevice           bool          `json:"bootdevice"`
	L3security           []string      `json:"l3security"`
	Firewall             FirewallRules `json:"firewall"`
	FirewallTemplateUuid string        `json:"firewall_template_uuid"`
}

//GetServerNetworkList gets a list of a specific server's networks
func (c *Client) GetServerNetworkList(id string) ([]ServerNetworkRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodGet,
	}
	var response ServerNetworkRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerNetwork gets a network of a specific server
func (c *Client) GetServerNetwork(serverId, networkId string) (ServerNetworkRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "networks", networkId),
		method: http.MethodGet,
	}
	var response ServerNetworkRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerNetwork updates a link between a network and a server
func (c *Client) UpdateServerNetwork(serverId, networkId string, body ServerNetworkRelationUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "networks", networkId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerNetwork creates a link between a network and a storage
func (c *Client) CreateServerNetwork(id string, body ServerNetworkRelationCreateRequest) error {
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
func (c *Client) LinkNetwork(serverId, networkId, firewallTemplate string, bootdevice bool, order int,
	l3security []string, firewall FirewallRules) error {
	body := ServerNetworkRelationCreateRequest{
		ObjectUuid:           networkId,
		Ordering:             order,
		BootDevice:           bootdevice,
		L3security:           l3security,
		FirewallTemplateUuid: firewallTemplate,
		Firewall:             firewall,
	}
	return c.CreateServerNetwork(serverId, body)
}

//UnlinkNetwork removes the link between a network and a server
func (c *Client) UnlinkNetwork(serverId string, networkId string) error {
	return c.DeleteServerNetwork(serverId, networkId)
}

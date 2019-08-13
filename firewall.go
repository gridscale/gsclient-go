package gsclient

import (
	"net/http"
	"path"
)

//FirewallList is JSON structure of a list of firewalls
type FirewallList struct {
	List map[string]FirewallProperties `json:"firewalls"`
}

//Firewall is JSON structure of a single firewall
type Firewall struct {
	Properties FirewallProperties `json:"firewall"`
}

type FirewallProperties struct {
	Status       string           `json:"status"`
	Labels       []string         `json:"labels"`
	ObjectUuid   string           `json:"object_uuid"`
	ChangeTime   string           `json:"change_time"`
	Rules        FirewallRules    `json:"rules"`
	CreateTime   string           `json:"create_time"`
	Private      bool             `json:"private"`
	Relations    FirewallRelation `json:"relations"`
	Description  string           `json:"description"`
	LocationName string           `json:"location_name"`
	Name         string           `json:"name"`
}

type FirewallRules struct {
	RulesV6In  []FirewallRuleProperties `json:"rules-v6-in"`
	RulesV6Out []FirewallRuleProperties `json:"rules-v6-out"`
	RulesV4In  []FirewallRuleProperties `json:"rules-v4-in"`
	RulesV4Out []FirewallRuleProperties `json:"rules-v4-out"`
}

type FirewallRuleProperties struct {
	Protocol string `json:"protocol"`
	DstPort  string `json:"dst_port"`
	SrcPort  string `json:"src_port"`
	SrcCidr  string `json:"src_cidr"`
	Action   string `json:"action"`
	Comment  string `json:"comment"`
	DstCidr  string `json:"dst_cidr"`
	Order    int    `json:"order"`
}

type FirewallRelation struct {
	Networks []NetworkInFirewall `json:"networks"`
}

type NetworkInFirewall struct {
	CreateTime  string `json:"create_time"`
	NetworkUuid string `json:"network_uuid"`
	NetworkName string `json:"network_name"`
	ObjectUuid  string `json:"object_uuid"`
	ObjectName  string `json:"object_name"`
}

type FirewallCreateRequest struct {
	Name   string        `json:"name"`
	Labels []string      `json:"labels"`
	Rules  FirewallRules `json:"rules"`
}

type FirewallCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

type FirewallUpdateRequest struct {
	Name   string        `json:"name"`
	Labels []string      `json:"labels"`
	Rules  FirewallRules `json:"rules"`
}

type FirewallEventList struct {
	List []FirewallEventProperties `json:"events"`
}

type FirewallEvent struct {
	Properties FirewallEventProperties `json:"event"`
}

type FirewallEventProperties struct {
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

//GetFirewallList gets a list of available firewalls
func (c *Client) GetFirewallList() ([]Firewall, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodGet,
	}
	var response FirewallList
	var firewalls []Firewall
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		firewalls = append(firewalls, Firewall{Properties: properties})
	}
	return firewalls, err
}

//GetFirewall gets a specific firewall based on given id
func (c *Client) GetFirewall(id string) (Firewall, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodGet,
	}
	var response Firewall
	err := r.execute(*c, &response)
	return response, err
}

//CreateFirewall creates a new firewall
func (c *Client) CreateFirewall(body FirewallCreateRequest) (FirewallCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodPost,
		body:   body,
	}
	var response FirewallCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return FirewallCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//UpdateFirewall update a specific firewall
func (c *Client) UpdateFirewall(id string, body FirewallUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteFirewall delete a specific firewall
func (c *Client) DeleteFirewall(id string) error {
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetFirewallEventList get list of a firewall's events
func (c *Client) GetFirewallEventList(id string) ([]FirewallEvent, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase, id, "events"),
		method: http.MethodGet,
	}
	var response FirewallEventList
	var firewallEvents []FirewallEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		firewallEvents = append(firewallEvents, FirewallEvent{Properties: properties})
	}
	return firewallEvents, err
}

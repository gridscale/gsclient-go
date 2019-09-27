package gsclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"
)

//FirewallList is JSON structure of a list of firewalls
type FirewallList struct {
	//Array of firewalls
	List map[string]FirewallProperties `json:"firewalls"`
}

//Firewall is JSON structure of a single firewall
type Firewall struct {
	//Properties of a firewall
	Properties FirewallProperties `json:"firewall"`
}

//FirewallProperties is JSON struct of a firewall's properties
type FirewallProperties struct {
	//Status indicates the status of the object
	Status string `json:"status"`

	//List of labels
	Labels []string `json:"labels"`

	//The UUID of an object is always unique, and refers to a specific object
	ObjectUUID string `json:"object_uuid"`

	//Defines the date and time of the last object change
	ChangeTime GSTime `json:"change_time"`

	//FirewallRules
	Rules FirewallRules `json:"rules"`

	//Defines the date and time the object was initially created
	CreateTime GSTime `json:"create_time"`

	//If this is a private or public Firewall-Template
	Private bool `json:"private"`

	//The information about other object which are related to this Firewall. The object could be Network.
	Relations FirewallRelation `json:"relations"`

	//Description of the ISO-Image release
	Description string `json:"description"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`
}

//FirewallRules is JSON struct of a list of firewall's rules
type FirewallRules struct {
	//Firewall template rules for inbound traffic - covers ipv6 addresses.
	RulesV6In []FirewallRuleProperties `json:"rules-v6-in,omitempty"`

	//Firewall template tules for outbound traffic - covers ipv6 addresses.
	RulesV6Out []FirewallRuleProperties `json:"rules-v6-out,omitempty"`

	//Firewall template rules for inbound traffic - covers ipv4 addresses.
	RulesV4In []FirewallRuleProperties `json:"rules-v4-in,omitempty"`

	//Firewall template tules for outbound traffic - covers ipv4 addresses.
	RulesV4Out []FirewallRuleProperties `json:"rules-v4-out,omitempty"`
}

//FirewallRuleProperties is JSON struct of a firewall's rule properties
type FirewallRuleProperties struct {
	//Enum:"udp" "tcp". Either udp or tcp
	Protocol string `json:"protocol,omitempty"`

	//A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
	DstPort string `json:"dst_port,omitempty"`

	//A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
	SrcPort string `json:"src_port,omitempty"`

	//A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
	SrcCidr string `json:"src_cidr,omitempty"`

	//Enum:"accept" "drop". This defines what the firewall will do. Either accept or drop.
	Action string `json:"action"`

	//Description
	Comment string `json:"comment,omitempty"`

	//Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then all IPs have access to this service.
	DstCidr string `json:"dst_cidr,omitempty"`

	//The order at which the firewall will compare packets against its rules,
	// a packet will be compared against the first rule, it will either allow it to pass
	// or block it and it won t be matched against any other rules.
	// However, if it does no match the rule, then it will proceed onto rule 2.
	// Packets that do not match any rules are blocked by default.
	Order int `json:"order"`
}

//FirewallRelation is a JSON struct of a list of firewall's relations
type FirewallRelation struct {
	//Array of object (NetworkinFirewall)
	Networks []NetworkInFirewall `json:"networks"`
}

//NetworkInFirewall is a JSON struct of a firewall's relation
type NetworkInFirewall struct {
	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The UUID of the network you're requesting.
	NetworkUUID string `json:"network_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	NetworkName string `json:"network_name"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`
}

//FirewallCreateRequest is JSON struct of a request for creating a firewall
type FirewallCreateRequest struct {
	//Name of firewall being created
	Name string `json:"name"`

	//Labels. Can be nil
	Labels []string `json:"labels,omitempty"`

	//FirewallRules
	Rules FirewallRules `json:"rules"`
}

//FirewallCreateResponse is JSON struct of a response for creating a firewall
type FirewallCreateResponse struct {
	//Request UUID
	RequestUUID string `json:"request_uuid"`

	//The UUID of the firewall being created
	ObjectUUID string `json:"object_uuid"`
}

//FirewallUpdateRequest is JSON struct of a request for updating a firewall
type FirewallUpdateRequest struct {
	//New name. Leave it if you do not want to update the name
	Name string `json:"name,omitempty"`

	//New list of labels. Leave it if you do not want to update the Labels
	Labels []string `json:"labels,omitempty"`

	//FirewallRules. Leave it if you do not want to update the firewall rules
	Rules *FirewallRules `json:"rules,omitempty"`
}

//GetFirewallList gets a list of available firewalls
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getFirewalls
func (c *Client) GetFirewallList(ctx context.Context) ([]Firewall, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodGet,
	}
	var response FirewallList
	var firewalls []Firewall
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		firewalls = append(firewalls, Firewall{Properties: properties})
	}
	return firewalls, err
}

//GetFirewall gets a specific firewall based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getFirewall
func (c *Client) GetFirewall(ctx context.Context, id string) (Firewall, error) {
	if !isValidUUID(id) {
		return Firewall{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodGet,
	}
	var response Firewall
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateFirewall creates a new firewall
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createFirewall
func (c *Client) CreateFirewall(ctx context.Context, body FirewallCreateRequest) (FirewallCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodPost,
		body:   body,
	}
	var response FirewallCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return FirewallCreateResponse{}, err
	}
	//Block until the request is finished
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
	return response, err
}

//UpdateFirewall update a specific firewall
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateFirewall
func (c *Client) UpdateFirewall(ctx context.Context, id string, body FirewallUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForFirewallActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//DeleteFirewall delete a specific firewall
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteFirewall
func (c *Client) DeleteFirewall(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForFirewallDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//GetFirewallEventList get list of a firewall's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getFirewallEvents
func (c *Client) GetFirewallEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var firewallEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		firewallEvents = append(firewallEvents, Event{Properties: properties})
	}
	return firewallEvents, err
}

//waitForFirewallActive allows to wait until the firewall's status is active
func (c *Client) waitForFirewallActive(ctx context.Context, id string) error {
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			errorMessage := fmt.Sprintf("Timeout reached when waiting for firewall %v to be active", id)
			c.cfg.logger.Error(errorMessage)
			return errors.New(errorMessage)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			fw, err := c.GetFirewall(ctx, id)
			if err != nil {
				return err
			}
			if fw.Properties.Status == resourceActiveStatus {
				return nil
			}
		}
	}
}

//waitForFirewallDeleted allows to wait until the firewall is deleted
func (c *Client) waitForFirewallDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			errorMessage := fmt.Sprintf("Timeout reached when waiting for firewall %v to be deleted", id)
			c.cfg.logger.Error(errorMessage)
			return errors.New(errorMessage)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			r := Request{
				uri:          path.Join(apiFirewallBase, id),
				method:       http.MethodGet,
				skipPrint404: true,
			}
			err := r.execute(ctx, *c, nil)
			if err != nil {
				if requestError, ok := err.(RequestError); ok {
					if requestError.StatusCode == 404 {
						return nil
					}
				}
				return err
			}
		}
	}
}

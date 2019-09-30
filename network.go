package gsclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

//NetworkList is JSON struct of a list of networks
type NetworkList struct {
	List map[string]NetworkProperties `json:"networks"`
}

//Network is JSON struct of a single network
type Network struct {
	Properties NetworkProperties `json:"network"`
}

//NetworkProperties is JSON struct of a network's properties
type NetworkProperties struct {
	LocationCountry string           `json:"location_country"`
	LocationUUID    string           `json:"location_uuid"`
	PublicNet       bool             `json:"public_net"`
	ObjectUUID      string           `json:"object_uuid"`
	NetworkType     string           `json:"network_type"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	CreateTime      string           `json:"create_time"`
	L2Security      bool             `json:"l2security"`
	ChangeTime      string           `json:"change_time"`
	LocationIata    string           `json:"location_iata"`
	LocationName    string           `json:"location_name"`
	DeleteBlock     bool             `json:"delete_block"`
	Labels          []string         `json:"labels"`
	Relations       NetworkRelations `json:"relations"`
}

//NetworkRelations is JSON struct of a list of a network's relations
type NetworkRelations struct {
	Vlans   []NetworkVlan   `json:"vlans"`
	Servers []NetworkServer `json:"servers"`
}

//NetworkVlan is JSON struct of a relation between a network and a VLAN
type NetworkVlan struct {
	Vlan       int    `json:"vlan"`
	TenantName string `json:"tenant_name"`
	TenantUUID string `json:"tenant_uuid"`
}

//NetworkServer is JSON struct of a relation between a network and a server
type NetworkServer struct {
	ObjectUUID  string   `json:"object_uuid"`
	Mac         string   `json:"mac"`
	Bootdevice  bool     `json:"bootdevice"`
	CreateTime  string   `json:"create_time"`
	L3security  []string `json:"l3security"`
	ObjectName  string   `json:"object_name"`
	NetworkUUID string   `json:"network_uuid"`
	Ordering    int      `json:"ordering"`
}

//NetworkCreateRequest is JSON of a request for creating a network
type NetworkCreateRequest struct {
	Name         string   `json:"name"`
	Labels       []string `json:"labels,omitempty"`
	LocationUUID string   `json:"location_uuid"`
	L2Security   bool     `json:"l2security,omitempty"`
}

//NetworkCreateResponse is JSON of a response for creating a network
type NetworkCreateResponse struct {
	ObjectUUID  string `json:"object_uuid"`
	RequestUUID string `json:"request_uuid"`
}

//NetworkUpdateRequest is JSON of a request for updating a network
type NetworkUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	L2Security bool   `json:"l2security"`
}

//NetworkEventList is JSON struct of a list of a network's events
type NetworkEventList struct {
	List []NetworkEventProperties `json:"events"`
}

//NetworkEvent is JSON struct of a single event of a network
type NetworkEvent struct {
	Properties NetworkEventProperties `json:"event"`
}

//NetworkEventProperties is JSON struct of properties of an event
type NetworkEventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUUID   string `json:"request_uuid"`
	ObjectUUID    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUUID      string `json:"user_uuid"`
}

//GetNetwork get a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getNetwork
func (c *Client) GetNetwork(ctx context.Context, id string) (Network, error) {
	if !isValidUUID(id) {
		return Network{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodGet,
	}
	var response Network
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateNetwork creates a network
//
//See: https://gridscale.io/en//api-documentation/index.html#tag/network
func (c *Client) CreateNetwork(ctx context.Context, body NetworkCreateRequest) (NetworkCreateResponse, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodPost,
		body:   body,
	}
	var response NetworkCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return NetworkCreateResponse{}, err
	}
<<<<<<< HEAD
	err = c.WaitForRequestCompletion(response.RequestUUID)
=======
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
>>>>>>> 8d4aa0e... add `context`
	return response, err
}

//DeleteNetwork deletes a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteNetwork
func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
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
		//Block until the request is finished
		return c.waitForNetworkDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//UpdateNetwork updates a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateNetwork
func (c *Client) UpdateNetwork(ctx context.Context, id string, body NetworkUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodPatch,
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
		//Block until the request is finished
		return c.waitForNetworkActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetNetworkList gets a list of available networks
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getNetworks
func (c *Client) GetNetworkList(ctx context.Context) ([]Network, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodGet,
	}
	var response NetworkList
	var networks []Network
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{
			Properties: properties,
		})
	}
	return networks, err
}

//GetNetworkEventList gets a list of a network's events
//
//See: https://gridscale.io/en//api-documentation/index.html#tag/network
func (c *Client) GetNetworkEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var networkEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		networkEvents = append(networkEvents, NetworkEvent{Properties: properties})
	}
	return networkEvents, err
}

//GetNetworkPublic gets public network
func (c *Client) GetNetworkPublic(ctx context.Context) (Network, error) {
	networks, err := c.GetNetworkList(ctx)
	if err != nil {
		return Network{}, err
	}
	for _, network := range networks {
		if network.Properties.PublicNet {
			return Network{Properties: network.Properties}, nil
		}
	}
	return Network{}, fmt.Errorf("Public Network not found")
}
<<<<<<< HEAD
=======

//GetNetworksByLocation gets a list of networks by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedNetworks
func (c *Client) GetNetworksByLocation(ctx context.Context, id string) ([]Network, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "networks"),
		method: http.MethodGet,
	}
	var response NetworkList
	var networks []Network
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{Properties: properties})
	}
	return networks, err
}

//GetDeletedNetworks gets a list of deleted networks
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedNetworks
func (c *Client) GetDeletedNetworks(ctx context.Context) ([]Network, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "networks"),
		method: http.MethodGet,
	}
	var response DeletedNetworkList
	var networks []Network
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{Properties: properties})
	}
	return networks, err
}

//waitForNetworkActive allows to wait until the network's status is active
func (c *Client) waitForNetworkActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		net, err := c.GetNetwork(ctx, id)
		return net.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForNetworkDeleted allows to wait until the network is deleted
func (c *Client) waitForNetworkDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiNetworkBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

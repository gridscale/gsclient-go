package gsclient

import (
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
func (c *Client) GetNetwork(id string) (Network, error) {
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodGet,
	}
	var response Network
	err := r.execute(*c, &response)
	return response, err
}

//CreateNetwork creates a network
func (c *Client) CreateNetwork(body NetworkCreateRequest) (NetworkCreateResponse, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodPost,
		body:   body,
	}
	var response NetworkCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return NetworkCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//DeleteNetwork deletes a specific network based on given id
func (c *Client) DeleteNetwork(id string) error {
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateNetwork updates a specific network based on given id
func (c *Client) UpdateNetwork(id string, body NetworkUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodPatch,
		body:   body,
	}

	return r.execute(*c, nil)
}

//GetNetworkList gets a list of available networks
func (c *Client) GetNetworkList() ([]Network, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodGet,
	}
	var response NetworkList
	var networks []Network
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{
			Properties: properties,
		})
	}
	return networks, err
}

//GetNetworkEventList gets a list of a network's events
func (c *Client) GetNetworkEventList(id string) ([]NetworkEvent, error) {
	r := Request{
		uri:    path.Join(apiNetworkBase, id, "events"),
		method: http.MethodGet,
	}
	var response NetworkEventList
	var networkEvents []NetworkEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networkEvents = append(networkEvents, NetworkEvent{Properties: properties})
	}
	return networkEvents, err
}

//GetNetworkPublic gets public network
func (c *Client) GetNetworkPublic() (Network, error) {
	networks, err := c.GetNetworkList()
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

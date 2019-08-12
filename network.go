package gsclient

import (
	"net/http"
	"path"
)

type Networks struct {
	List map[string]NetworkProperties `json:"networks"`
}

type Network struct {
	Properties NetworkProperties `json:"network"`
}

type NetworkProperties struct {
	LocationCountry string           `json:"location_country"`
	LocationUuid    string           `json:"location_uuid"`
	PublicNet       bool             `json:"public_net"`
	ObjectUuid      string           `json:"object_uuid"`
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

type NetworkRelations struct {
	Vlans   []NetworkVlan   `json:"vlans"`
	Servers []NetworkServer `json:"servers"`
}

type NetworkVlan struct {
	Vlan       int    `json:"vlan"`
	TenantName string `json:"tenant_name"`
	TenantUuid string `json:"tenant_uuid"`
}

type NetworkServer struct {
	ObjectUuid  string   `json:"object_uuid"`
	Mac         string   `json:"mac"`
	Bootdevice  bool     `json:"bootdevice"`
	CreateTime  string   `json:"create_time"`
	L3security  []string `json:"l3security"`
	ObjectName  string   `json:"object_name"`
	NetworkUuid string   `json:"network_uuid"`
	Ordering    int      `json:"ordering"`
}

type NetworkCreateRequest struct {
	Name         string   `json:"name"`
	Labels       []string `json:"labels,omitempty"`
	LocationUuid string   `json:"location_uuid"`
	L2Security   bool     `json:"l2security,omitempty"`
}

type NetworkCreateResponse struct {
	ObjectUuid  string `json:"object_uuid"`
	RequestUuid string `json:"request_uuid"`
}

type NetworkUpdateRequest struct {
	Name       string   `json:"name,omitempty"`
	Labels     []string `json:"labels"`
	L2Security bool     `json:"l2security"`
}

type NetworkEventList struct {
	List []NetworkEventProperties `json:"events"`
}

type NetworkEvent struct {
	Properties NetworkEventProperties `json:"event"`
}

type NetworkEventProperties struct {
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

	err = c.WaitForRequestCompletion(response.RequestUuid)

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

	var response Networks
	err := r.execute(*c, &response)

	list := []Network{}
	for _, properties := range response.List {
		list = append(list, Network{
			Properties: properties,
		})
	}

	return list, err
}

//GetNetworkEventList gets a list of a network's events
func (c *Client) GetNetworkEventList(id string) ([]NetworkEvent, error) {
	r := Request{
		uri:    path.Join(apiNetworkBase, id, "events"),
		method: http.MethodGet,
	}
	var response NetworkEventList
	var list []NetworkEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, NetworkEvent{Properties: properties})
	}
	return list, err
}

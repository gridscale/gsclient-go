package gsclient

import (
	"net/http"
	"path"
)

//IpList is JSON struct of a list of IPs
type IpList struct {
	List map[string]IpProperties `json:"ips"`
}

//Ip is JSON struct if a single IP
type Ip struct {
	Properties IpProperties `json:"ip"`
}

//IpProperties is JSON struct of an IP's properties
type IpProperties struct {
	Name            string      `json:"name"`
	LocationCountry string      `json:"location_country"`
	LocationUuid    string      `json:"location_uuid"`
	ObjectUuid      string      `json:"object_uuid"`
	ReverseDns      string      `json:"reverse_dns"`
	Family          int         `json:"family"`
	Status          string      `json:"status"`
	CreateTime      string      `json:"create_time"`
	Failover        bool        `json:"failover"`
	ChangeTime      string      `json:"change_time"`
	LocationIata    string      `json:"location_iata"`
	LocationName    string      `json:"location_name"`
	Prefix          string      `json:"prefix"`
	Ip              string      `json:"ip"`
	DeleteBlock     string      `json:"delete_block"`
	UsagesInMinutes float64     `json:"usage_in_minutes"`
	CurrentPrice    float64     `json:"current_price"`
	Labels          []string    `json:"labels"`
	Relations       IpRelations `json:"relations"`
}

//IpRelations is JSON struct of a list of an IP's relations
type IpRelations struct {
	Loadbalancers []IpLoadbalancer                  `json:"loadbalancers"`
	Servers       []IpServer                        `json:"servers"`
	PublicIps     []ServerIpRelationProperties      `json:"public_ips"`
	Storages      []ServerStorageRelationProperties `json:"storages"`
}

//IpLoadbalancer is JSON struct of the relation between an IP and a Load Balancer
type IpLoadbalancer struct {
	CreateTime       string `json:"create_time"`
	LoadbalancerName string `json:"loadbalancer_name"`
	LoadbalancerUuid string `json:"loadbalancer_uuid"`
}

//IpServer is JSON struct of the relation between an IP and a Server
type IpServer struct {
	CreateTime string `json:"create_time"`
	ServerName string `json:"server_name"`
	ServerUuid string `json:"server_uuid"`
}

//IpCreateResponse is JSON struct of a response for creating an IP
type IpCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
	Prefix      string `json:"prefix"`
	Ip          string `json:"ip"`
}

//IpCreateRequest is JSON struct of a request for creating an IP
type IpCreateRequest struct {
	Name         string   `json:"name,omitempty"`
	Family       int      `json:"family"`
	LocationUuid string   `json:"location_uuid"`
	Failover     bool     `json:"failover,omitempty"`
	ReverseDns   string   `json:"reverse_dns,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

//IpUpdateRequest is JSON struct of a request for updating an IP
type IpUpdateRequest struct {
	Name       string   `json:"name,omitempty"`
	Failover   bool     `json:"failover"`
	ReverseDns string   `json:"reverse_dns,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

//IpEventList is JSON struct of a list of an IP's events
type IpEventList struct {
	List []IpEventProperties `json:"events"`
}

//IpEvent is JSON struct of a single IP
type IpEvent struct {
	Properties IpEventProperties `json:"event"`
}

//IpEventProperties is JSON struct of an IP's properties
type IpEventProperties struct {
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

//GetIp get a specific IP based on given id
func (c *Client) GetIp(id string) (Ip, error) {
	r := Request{
		uri:    path.Join(apiIpBase, id),
		method: http.MethodGet,
	}

	var response Ip
	err := r.execute(*c, &response)

	return response, err
}

//GetIpList gets a list of available IPs
func (c *Client) GetIpList() ([]Ip, error) {
	r := Request{
		uri:    apiIpBase,
		method: http.MethodGet,
	}

	var response IpList
	var IPs []Ip
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, Ip{Properties: properties})
	}

	return IPs, err
}

//CreateIp creates an IP
func (c *Client) CreateIp(body IpCreateRequest) (IpCreateResponse, error) {
	r := Request{
		uri:    apiIpBase,
		method: http.MethodPost,
		body:   body,
	}

	var response IpCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return IpCreateResponse{}, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

//DeleteIp deletes a specific IP based on given id
func (c *Client) DeleteIp(id string) error {
	r := Request{
		uri:    path.Join(apiIpBase, id),
		method: http.MethodDelete,
	}

	return r.execute(*c, nil)
}

//UpdateIp updates a specific IP based on given id
func (c *Client) UpdateIp(id string, body IpUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiIpBase, id),
		method: http.MethodPatch,
		body:   body,
	}

	return r.execute(*c, nil)
}

//GetIpEventList gets a list of an IP's events
func (c *Client) GetIpEventList(id string) ([]IpEvent, error) {
	r := Request{
		uri:    path.Join(apiNetworkBase, id, "events"),
		method: http.MethodGet,
	}
	var response IpEventList
	var IPEvents []IpEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		IPEvents = append(IPEvents, IpEvent{Properties: properties})
	}
	return IPEvents, err
}

//GetIpVersion gets IP's version, returns 0 if an error was encountered
func (c *Client) GetIpVersion(id string) int {
	ip, err := c.GetIp(id)
	if err != nil {
		return 0
	}
	return ip.Properties.Family
}

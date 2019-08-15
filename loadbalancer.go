package gsclient

import (
	"net/http"
	"path"
)

//LoadBalancers is the JSON struct of a list of load balancers
type LoadBalancers struct {
	List map[string]LoadBalancerProperties `json:"loadbalancers"`
}

//LoadBalancer is the JSON struct of a load balancer
type LoadBalancer struct {
	Properties LoadBalancerProperties `json:"loadbalancer"`
}

//LoadBalancerProperties is the properties of a load balancer
type LoadBalancerProperties struct {
	ObjectUuid          string           `json:"object_uuid"`
	LocationSite        int              `json:"location_site"`
	Name                string           `json:"name"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	LocationIata        string           `json:"location_iata"`
	LocationUuid        string           `json:"location_uuid"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	ChangeTime          string           `json:"change_time"`
	Status              string           `json:"status"`
	CurrentPrice        float64          `json:"current_price"`
	LocationCountry     string           `json:"location_country"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Labels              []string         `json:"labels"`
	LocationName        string           `json:"location_name"`
	UsageInMinutes      int              `json:"usage_in_minutes"`
	Algorithm           string           `json:"algorithm"`
	CreateTime          string           `json:"create_time"`
	ListenIPv6Uuid      string           `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string           `json:"listen_ipv4_uuid"`
}

//BackendServer is the JSON struct of backend server
type BackendServer struct {
	Weight int    `json:"weight"`
	Host   string `json:"host"`
}

//ForwardingRule is the JSON struct of forwarding rule
type ForwardingRule struct {
	LetsencryptSSL interface{} `json:"letsencrypt_ssl"`
	ListenPort     int         `json:"listen_port"`
	Mode           string      `json:"mode"`
	TargetPort     int         `json:"target_port"`
}

//LoadBalancerCreateRequest is the JSON struct for creating a load balancer request
type LoadBalancerCreateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6Uuid      string           `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []string         `json:"labels"`
	LocationUuid        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerUpdateRequest is the JSON struct for updating a load balancer request
type LoadBalancerUpdateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6Uuid      string           `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []string         `json:"labels"`
	LocationUuid        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerCreateResponse is the JSON struct for a load balancer response
type LoadBalancerCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

//LoadBalancerEventList is the JSON struct for a load alancer's events
type LoadBalancerEventList struct {
	List []LoadBalancerEventProperties `json:"events"`
}

//LoadBalancerEvent is JSON struct for a load balancer
type LoadBalancerEvent struct {
	Properties LoadBalancerEventProperties `json:"event"`
}

//LoadBalancerEventProperties is the properties of a load balancer's event
type LoadBalancerEventProperties struct {
	ObjectUuid    string `json:"object_uuid"`
	ObjectType    string `json:"object_type"`
	RequestUuid   string `json:"request_uuid"`
	RequestType   string `json:"request_type"`
	Activity      string `json:"activity"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUuid      string `json:"user_uuid"`
}

//GetLoadBalancerList returns a list of load balancers
func (c *Client) GetLoadBalancerList() ([]LoadBalancer, error) {
	r := Request{
		uri:    apiLoadBalancerBase,
		method: http.MethodGet,
	}
	var response LoadBalancers
	var loadBalancers []LoadBalancer
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		loadBalancers = append(loadBalancers, LoadBalancer{Properties: properties})
	}
	return loadBalancers, err
}

//GetLoadBalancer returns a load balancer of a given uuid
func (c *Client) GetLoadBalancer(id string) (LoadBalancer, error) {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodGet,
	}
	var response LoadBalancer
	err := r.execute(*c, &response)
	return response, err
}

//CreateLoadBalancer creates a new load balancer
func (c *Client) CreateLoadBalancer(body LoadBalancerCreateRequest) (LoadBalancerCreateResponse, error) {
	r := Request{
		uri:    apiLoadBalancerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response LoadBalancerCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return LoadBalancerCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//UpdateLoadBalancer update configuration of a load balancer
func (c *Client) UpdateLoadBalancer(id string, body LoadBalancerUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetLoadBalancerEventList retrieves events of a given uuid
func (c *Client) GetLoadBalancerEventList(id string) ([]LoadBalancerEvent, error) {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id, "events"),
		method: http.MethodGet,
	}
	var response LoadBalancerEventList
	var loadBalancerEvents []LoadBalancerEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		loadBalancerEvents = append(loadBalancerEvents, LoadBalancerEvent{Properties: properties})
	}
	return loadBalancerEvents, err
}

//DeleteLoadBalancer deletes a load balancer
func (c *Client) DeleteLoadBalancer(id string) error {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

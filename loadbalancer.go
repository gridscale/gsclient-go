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

//LoadBalancers is the JSON struct of a list of loadbalancers
type LoadBalancers struct {
	List map[string]LoadBalancerProperties `json:"loadbalancers"`
}

//LoadBalancer is the JSON struct of a loadbalancer
type LoadBalancer struct {
	Properties LoadBalancerProperties `json:"loadbalancer"`
}

//LoadBalancerProperties is the properties of a loadbalancer
type LoadBalancerProperties struct {
	ObjectUUID          string           `json:"object_uuid"`
	LocationSite        int              `json:"location_site"`
	Name                string           `json:"name"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	LocationIata        string           `json:"location_iata"`
	LocationUUID        string           `json:"location_uuid"`
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
	ListenIPv6UUID      string           `json:"listen_ipv6_uuid"`
	ListenIPv4UUID      string           `json:"listen_ipv4_uuid"`
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

//LoadBalancerCreateRequest is the JSON struct for creating a loadbalancer request
type LoadBalancerCreateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6UUID      string           `json:"listen_ipv6_uuid"`
	ListenIPv4UUID      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []string         `json:"labels"`
	LocationUUID        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerUpdateRequest is the JSON struct for updating a loadbalancer request
type LoadBalancerUpdateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6UUID      string           `json:"listen_ipv6_uuid"`
	ListenIPv4UUID      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []string         `json:"labels"`
	LocationUUID        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerCreateResponse is the JSON struct for a loadbalancer response
type LoadBalancerCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//LoadBalancerEventList is the JSON struct for a loadbalancer's events
type LoadBalancerEventList struct {
	List []LoadBalancerEventProperties `json:"events"`
}

//LoadBalancerEvent is JSON struct for a load balancer
type LoadBalancerEvent struct {
	Properties LoadBalancerEventProperties `json:"event"`
}

//LoadBalancerEventProperties is the properties of a loadbalancer's event
type LoadBalancerEventProperties struct {
	ObjectUUID    string `json:"object_uuid"`
	ObjectType    string `json:"object_type"`
	RequestUUID   string `json:"request_uuid"`
	RequestType   string `json:"request_type"`
	Activity      string `json:"activity"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUUID      string `json:"user_uuid"`
}

//GetLoadBalancerList returns a list of loadbalancers
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancers
func (c *Client) GetLoadBalancerList(ctx context.Context) ([]LoadBalancer, error) {
	r := Request{
		uri:    apiLoadBalancerBase,
		method: http.MethodGet,
	}
	var response LoadBalancers
	var loadBalancers []LoadBalancer
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		loadBalancers = append(loadBalancers, LoadBalancer{Properties: properties})
	}
	return loadBalancers, err
}

//GetLoadBalancer returns a loadbalancer of a given uuid
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancer
func (c *Client) GetLoadBalancer(ctx context.Context, id string) (LoadBalancer, error) {
	if !isValidUUID(id) {
		return LoadBalancer{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodGet,
	}
	var response LoadBalancer
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateLoadBalancer creates a new loadbalancer
//
//Note: loadbalancer's algorithm can only be either `LoadbalancerRoundrobinAlg` or `LoadbalancerLeastConnAlg`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createLoadbalancer
func (c *Client) CreateLoadBalancer(ctx context.Context, body LoadBalancerCreateRequest) (LoadBalancerCreateResponse, error) {
	if body.Labels == nil {
		body.Labels = make([]string, 0)
	}
	r := Request{
		uri:    apiLoadBalancerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response LoadBalancerCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return LoadBalancerCreateResponse{}, err
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

//UpdateLoadBalancer update configuration of a loadbalancer
//
//Note: loadbalancer's algorithm can only be either `LoadbalancerRoundrobinAlg` or `LoadbalancerLeastConnAlg`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateLoadbalancer
func (c *Client) UpdateLoadBalancer(ctx context.Context, id string, body LoadBalancerUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	if body.Labels == nil {
		body.Labels = make([]string, 0)
	}
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
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
		return c.waitForLoadbalancerActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetLoadBalancerEventList retrieves events of a given uuid
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLoadbalancerEvents
func (c *Client) GetLoadBalancerEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var loadBalancerEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		loadBalancerEvents = append(loadBalancerEvents, LoadBalancerEvent{Properties: properties})
	}
	return loadBalancerEvents, err
}

//DeleteLoadBalancer deletes a loadbalancer
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteLoadbalancer
func (c *Client) DeleteLoadBalancer(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: http.MethodDelete,
	}
<<<<<<< HEAD
	return r.execute(*c, nil)
}
=======
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForLoadbalancerDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//waitForLoadbalancerActive allows to wait until the loadbalancer's status is active
func (c *Client) waitForLoadbalancerActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		lb, err := c.GetLoadBalancer(ctx, id)
		return lb.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForLoadbalancerDeleted allows to wait until the loadbalancer is deleted
func (c *Client) waitForLoadbalancerDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiLoadBalancerBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

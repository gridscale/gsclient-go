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

//IPList is JSON struct of a list of IPs
type IPList struct {
	List map[string]IPProperties `json:"ips"`
}

//IP is JSON struct if a single IP
type IP struct {
	Properties IPProperties `json:"ip"`
}

//IPProperties is JSON struct of an IP's properties
type IPProperties struct {
	Name            string      `json:"name"`
	LocationCountry string      `json:"location_country"`
	LocationUUID    string      `json:"location_uuid"`
	ObjectUUID      string      `json:"object_uuid"`
	ReverseDNS      string      `json:"reverse_dns"`
	Family          int         `json:"family"`
	Status          string      `json:"status"`
	CreateTime      string      `json:"create_time"`
	Failover        bool        `json:"failover"`
	ChangeTime      string      `json:"change_time"`
	LocationIata    string      `json:"location_iata"`
	LocationName    string      `json:"location_name"`
	Prefix          string      `json:"prefix"`
	IP              string      `json:"ip"`
	DeleteBlock     string      `json:"delete_block"`
	UsagesInMinutes float64     `json:"usage_in_minutes"`
	CurrentPrice    float64     `json:"current_price"`
	Labels          []string    `json:"labels"`
	Relations       IPRelations `json:"relations"`
}

//IPRelations is JSON struct of a list of an IP's relations
type IPRelations struct {
	Loadbalancers []IPLoadbalancer                  `json:"loadbalancers"`
	Servers       []IPServer                        `json:"servers"`
	PublicIPs     []ServerIPRelationProperties      `json:"public_ips"`
	Storages      []ServerStorageRelationProperties `json:"storages"`
}

//IPLoadbalancer is JSON struct of the relation between an IP and a Load Balancer
type IPLoadbalancer struct {
	CreateTime       string `json:"create_time"`
	LoadbalancerName string `json:"loadbalancer_name"`
	LoadbalancerUUID string `json:"loadbalancer_uuid"`
}

//IPServer is JSON struct of the relation between an IP and a Server
type IPServer struct {
	CreateTime string `json:"create_time"`
	ServerName string `json:"server_name"`
	ServerUUID string `json:"server_uuid"`
}

//IPCreateResponse is JSON struct of a response for creating an IP
type IPCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
	Prefix      string `json:"prefix"`
	IP          string `json:"ip"`
}

//IPCreateRequest is JSON struct of a request for creating an IP
type IPCreateRequest struct {
	Name         string   `json:"name,omitempty"`
	Family       int      `json:"family"`
	LocationUUID string   `json:"location_uuid"`
	Failover     bool     `json:"failover,omitempty"`
	ReverseDNS   string   `json:"reverse_dns,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

//IPUpdateRequest is JSON struct of a request for updating an IP
type IPUpdateRequest struct {
	Name       string   `json:"name,omitempty"`
	Failover   bool     `json:"failover"`
	ReverseDNS string   `json:"reverse_dns,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

//IPEventList is JSON struct of a list of an IP's events
type IPEventList struct {
	List []IPEventProperties `json:"events"`
}

//IPEvent is JSON struct of a single IP
type IPEvent struct {
	Properties IPEventProperties `json:"event"`
}

//IPEventProperties is JSON struct of an IP's properties
type IPEventProperties struct {
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

//GetIP get a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIp
func (c *Client) GetIP(ctx context.Context, id string) (IP, error) {
	if !isValidUUID(id) {
		return IP{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiIPBase, id),
		method: http.MethodGet,
	}

	var response IP
	err := r.execute(ctx, *c, &response)

	return response, err
}

//GetIPList gets a list of available IPs
//
//https://gridscale.io/en//api-documentation/index.html#operation/getIps
func (c *Client) GetIPList(ctx context.Context) ([]IP, error) {
	r := Request{
		uri:    apiIPBase,
		method: http.MethodGet,
	}

	var response IPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}

	return IPs, err
}

//CreateIP creates an IP
//
//Note: IP address family can only be either `IPv4Type` or `IPv6Type`
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createIp
func (c *Client) CreateIP(ctx context.Context, body IPCreateRequest) (IPCreateResponse, error) {
	r := Request{
		uri:    apiIPBase,
		method: http.MethodPost,
		body:   body,
	}

	var response IPCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return IPCreateResponse{}, err
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

//DeleteIP deletes a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteIp
func (c *Client) DeleteIP(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiIPBase, id),
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
		return c.waitForIPDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//UpdateIP updates a specific IP based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateIp
func (c *Client) UpdateIP(ctx context.Context, id string, body IPUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiIPBase, id),
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
		return c.waitForIPActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetIPEventList gets a list of an IP's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getIpEvents
func (c *Client) GetIPEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiIPBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var IPEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPEvents = append(IPEvents, IPEvent{Properties: properties})
	}
	return IPEvents, err
}

//GetIPVersion gets IP's version, returns 0 if an error was encountered
func (c *Client) GetIPVersion(ctx context.Context, id string) int {
	ip, err := c.GetIP(ctx, id)
	if err != nil {
		return 0
	}
	return ip.Properties.Family
}
<<<<<<< HEAD
=======

//GetIPsByLocation gets a list of IPs by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationIps
func (c *Client) GetIPsByLocation(ctx context.Context, id string) ([]IP, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "ips"),
		method: http.MethodGet,
	}
	var response IPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}
	return IPs, err
}

//GetDeletedIPs gets a list of deleted IPs
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedIps
func (c *Client) GetDeletedIPs(ctx context.Context) ([]IP, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "ips"),
		method: http.MethodGet,
	}
	var response DeletedIPList
	var IPs []IP
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		IPs = append(IPs, IP{Properties: properties})
	}
	return IPs, err
}

//waitForIPActive allows to wait until the IP address's status is active
func (c *Client) waitForIPActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		ip, err := c.GetIP(ctx, id)
		return ip.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForIPDeleted allows to wait until the IP address is deleted
func (c *Client) waitForIPDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiIPBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

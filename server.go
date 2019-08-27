package gsclient

import (
	"net/http"
	"path"
)

//ServerList JSON struct of a list of servers
type ServerList struct {
	List map[string]ServerProperties `json:"servers"`
}

//Server JSON struct of a single server
type Server struct {
	Properties ServerProperties `json:"server"`
}

//ServerProperties JSON struct of properties of a server
type ServerProperties struct {
	ObjectUUID           string          `json:"object_uuid"`
	Name                 string          `json:"name"`
	Memory               int             `json:"memory"`
	Cores                int             `json:"cores"`
	HardwareProfile      string          `json:"hardware_profile"`
	Status               string          `json:"status"`
	LocationUUID         string          `json:"location_uuid"`
	Power                bool            `json:"power"`
	CurrentPrice         float64         `json:"current_price"`
	AvailablityZone      string          `json:"availability_zone"`
	AutoRecovery         bool            `json:"auto_recovery"`
	Legacy               bool            `json:"legacy"`
	ConsoleToken         string          `json:"console_token"`
	UsageInMinutesMemory int             `json:"usage_in_minutes_memory"`
	UsageInMinutesCores  int             `json:"usage_in_minutes_cores"`
	Labels               []string        `json:"labels"`
	Relations            ServerRelations `json:"relations"`
}

//ServerRelations JSON struct of a list of server relations
type ServerRelations struct {
	IsoImages []ServerIsoImageRelationProperties `json:"isoimages"`
	Networks  []ServerNetworkRelationProperties  `json:"networks"`
	PublicIPs []ServerIPRelationProperties       `json:"public_ips"`
	Storages  []ServerStorageRelationProperties  `json:"storages"`
}

//ServerCreateRequest JSON struct of a request for creating a server
type ServerCreateRequest struct {
	Name            string   `json:"name"`
	Memory          int      `json:"memory"`
	Cores           int      `json:"cores"`
	LocationUUID    string   `json:"location_uuid"`
	HardwareProfile string   `json:"hardware_profile,omitempty"`
	AvailablityZone string   `json:"availability_zone,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	Status          string   `json:"status,omitempty"`
	AutoRecovery    *bool    `json:"auto_recovery,omitempty"`
}

//ServerCreateResponse JSON struct of a response for creating a server
type ServerCreateResponse struct {
	ObjectUUID   string   `json:"object_uuid"`
	RequestUUID  string   `json:"request_uuid"`
	SeverUUID    string   `json:"sever_uuid"`
	NetworkUUIDs []string `json:"network_uuids"`
	StorageUUIDs []string `json:"storage_uuids"`
	IPaddrUUIDs  []string `json:"ipaddr_uuids"`
}

//ServerPowerUpdateRequest JSON struct of a request for updating server's power state
type ServerPowerUpdateRequest struct {
	Power bool `json:"power"`
}

//ServerCreateRequestStorage JSON struct of a relation between a server and a storage
type ServerCreateRequestStorage struct {
	StorageUUID string `json:"storage_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

//ServerCreateRequestNetwork JSON struct of a relation between a server and a network
type ServerCreateRequestNetwork struct {
	NetworkUUID string `json:"network_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

//ServerCreateRequestIP JSON struct of a relation between a server and an IP address
type ServerCreateRequestIP struct {
	IPaddrUUID string `json:"ipaddr_uuid,omitempty"`
}

//ServerCreateRequestIsoimage JSON struct of a relation between a server and an ISO-Image
type ServerCreateRequestIsoimage struct {
	IsoimageUUID string `json:"isoimage_uuid,omitempty"`
}

//ServerUpdateRequest JSON of a request for updating a server
type ServerUpdateRequest struct {
	Name            string   `json:"name,omitempty"`
	AvailablityZone string   `json:"availability_zone,omitempty"`
	Memory          int      `json:"memory,omitempty"`
	Cores           int      `json:"cores,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	AutoRecovery    *bool    `json:"auto_recovery,omitempty"`
}

//ServerEventList JSON struct of a list of a server's events
type ServerEventList struct {
	List []ServerEventProperties `json:"events"`
}

//ServerEvent JSON struct of a single event of a server
type ServerEvent struct {
	Properties ServerEventProperties `json:"event"`
}

//ServerEventProperties JSON struct of properties of a server
type ServerEventProperties struct {
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

//ServerMetricList JSON struct of a list of a server's metrics
type ServerMetricList struct {
	List []ServerMetricProperties `json:"server_metrics"`
}

//ServerMetric JSON struct of a single metric of a server
type ServerMetric struct {
	Properties ServerMetricProperties `json:"server_metric"`
}

//ServerMetricProperties JSON stru
type ServerMetricProperties struct {
	BeginTime       string `json:"begin_time"`
	EndTime         string `json:"end_time"`
	PaaSServiceUUID string `json:"paas_service_uuid"`
	CoreUsage       struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"core_usage"`
	StorageSize struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"storage_size"`
}

//GetServer gets a specific server based on given list
func (c *Client) GetServer(id string) (Server, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodGet,
	}
	var response Server
	err := r.execute(*c, &response)
	return response, err
}

//GetServerList gets a list of available servers
func (c *Client) GetServerList() ([]Server, error) {
	r := Request{
		uri:    apiServerBase,
		method: http.MethodGet,
	}
	var response ServerList
	var servers []Server
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{
			Properties: properties,
		})
	}
	return servers, err
}

//CreateServer create a server
func (c *Client) CreateServer(body ServerCreateRequest) (ServerCreateResponse, error) {
	r := Request{
		uri:    apiServerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response ServerCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return ServerCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//DeleteServer deletes a specific server
func (c *Client) DeleteServer(id string) error {
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateServer updates a specific server
func (c *Client) UpdateServer(id string, body ServerUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetServerEventList gets a list of a specific server's events
func (c *Client) GetServerEventList(id string) ([]ServerEvent, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "events"),
		method: http.MethodGet,
	}
	var response ServerEventList
	var serverEvents []ServerEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		serverEvents = append(serverEvents, ServerEvent{Properties: properties})
	}
	return serverEvents, err
}

//GetServerMetricList gets a list of a specific server's metrics
func (c *Client) GetServerMetricList(id string) ([]ServerMetric, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "metrics"),
		method: http.MethodGet,
	}
	var response ServerMetricList
	var serverMetrics []ServerMetric
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		serverMetrics = append(serverMetrics, ServerMetric{Properties: properties})
	}
	return serverMetrics, err
}

//IsServerOn returns true if the server's power is on, otherwise returns false
func (c *Client) IsServerOn(id string) (bool, error) {
	server, err := c.GetServer(id)
	if err != nil {
		return false, err
	}
	return server.Properties.Power, nil
}

//setServerPowerState turn on/off a specific server.
//turnOn=true to turn on, turnOn=false to turn off
func (c *Client) setServerPowerState(id string, powerState bool) error {
	isOn, err := c.IsServerOn(id)
	if err != nil {
		return err
	}
	if isOn == powerState {
		return nil
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "power"),
		method: http.MethodPatch,
		body: ServerPowerUpdateRequest{
			Power: powerState,
		},
	}
	err = r.execute(*c, nil)
	if err != nil {
		return err
	}
	return c.WaitForServerPowerStatus(id, powerState)
}

//StartServer starts a server
func (c *Client) StartServer(id string) error {
	return c.setServerPowerState(id, true)
}

//StopServer stops a server
func (c *Client) StopServer(id string) error {
	return c.setServerPowerState(id, false)
}

//ShutdownServer shutdowns a specific server
func (c *Client) ShutdownServer(id string) error {
	//Make sure the server exists and that it isn't already in the state we need it to be
	server, err := c.GetServer(id)
	if err != nil {
		return err
	}
	if !server.Properties.Power {
		return nil
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "shutdown"),
		method: http.MethodPatch,
		body:   new(map[string]string),
	}
	err = r.execute(*c, nil)
	if err != nil {
		return err
	}
	//If we get an error, which includes a timeout, power off the server instead
	err = c.WaitForServerPowerStatus(id, false)
	if err != nil {
		return c.setServerPowerState(id, false)
	}
	return nil
}

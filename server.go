package gsclient

import (
	"net/http"
	"path"
)

type Servers struct {
	List map[string]ServerProperties `json:"servers"`
}

type Server struct {
	Properties ServerProperties `json:"server"`
}

type ServerProperties struct {
	ObjectUuid           string          `json:"object_uuid"`
	Name                 string          `json:"name"`
	Memory               int             `json:"memory"`
	Cores                int             `json:"cores"`
	HardwareProfile      string          `json:"hardware_profile"`
	Status               string          `json:"status"`
	LocationUuid         string          `json:"location_uuid"`
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

type ServerRelations struct {
	IsoImages []ServerIsoImage `json:"isoimages"`
	Networks  []ServerNetwork  `json:"networks"`
	PublicIps []ServerIp       `json:"public_ips"`
	Storages  []ServerStorage  `json:"storages"`
}

type ServerCreateRequest struct {
	Name            string                       `json:"name"`
	Memory          int                          `json:"memory"`
	Cores           int                          `json:"cores"`
	LocationUuid    string                       `json:"location_uuid"`
	HardwareProfile string                       `json:"hardware_profile,omitempty"`
	AvailablityZone string                       `json:"availability_zone,omitempty"`
	Labels          []string                     `json:"labels,omitempty"`
	Relations       ServerCreateRequestRelations `json:"relations,omitempty"`
}

type ServerCreateResponse struct {
	ObjectUuid   string   `json:"object_uuid"`
	RequestUuid  string   `json:"request_uuid"`
	SeverUuid    string   `json:"sever_uuid"`
	NetworkUuids []string `json:"network_uuids"`
	StorageUuids []string `json:"storage_uuids"`
	IpaddrUuids  []string `json:"ipaddr_uuids"`
}

type ServerPowerUpdateRequest struct {
	Power bool `json:"power"`
}

type ServerCreateRequestRelations struct {
	IsoImages []ServerCreateRequestIsoimage `json:"isoimages"`
	Networks  []ServerCreateRequestNetwork  `json:"networks"`
	PublicIps []ServerCreateRequestIp       `json:"public_ips"`
	Storages  []ServerCreateRequestStorage  `json:"storages"`
}

type ServerCreateRequestStorage struct {
	StorageUuid string `json:"storage_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerCreateRequestNetwork struct {
	NetworkUuid string `json:"network_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerCreateRequestIp struct {
	IpaddrUuid string `json:"ipaddr_uuid,omitempty"`
}

type ServerCreateRequestIsoimage struct {
	IsoimageUuid string `json:"isoimage_uuid,omitempty"`
}

type ServerUpdateRequest struct {
	Name            string   `json:"name,omitempty"`
	AvailablityZone string   `json:"availability_zone,omitempty"`
	Memory          int      `json:"memory,omitempty"`
	Cores           int      `json:"cores,omitempty"`
	Labels          []string `json:"labels"`
}

type ServerEventList struct {
	List []ServerEventProperties `json:"events"`
}

type ServerEvent struct {
	Properties ServerEventProperties `json:"event"`
}

type ServerEventProperties struct {
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

type ServerMetricList struct {
	List []ServerMetricProperties `json:"server_metrics"`
}

type ServerMetric struct {
	Properties ServerMetricProperties `json:"server_metric"`
}

type ServerMetricProperties struct {
	BeginTime       string `json:"begin_time"`
	EndTime         string `json:"end_time"`
	PaaSServiceUuid string `json:"paas_service_uuid"`
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

	var response Servers
	err := r.execute(*c, &response)

	list := []Server{}
	for _, properties := range response.List {
		list = append(list, Server{
			Properties: properties,
		})
	}

	return list, err
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

	err = c.WaitForRequestCompletion(response.RequestUuid)

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
	var list []ServerEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ServerEvent{Properties: properties})
	}
	return list, err
}

//GetServerMetricList gets a list of a specific server's metrics
func (c *Client) GetServerMetricList(id string) ([]ServerMetric, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "metrics"),
		method: http.MethodGet,
	}
	var response ServerMetricList
	var list []ServerMetric
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ServerMetric{Properties: properties})
	}
	return list, err
}

//IsServerOn returns true if the server's power is on, otherwise returns false
func (c *Client) IsServerOn(id string) (bool, error) {
	server, err := c.GetServer(id)
	if err != nil {
		return false, err
	}
	return server.Properties.Power, nil
}

//turnOnOffServer turn on/off a specific server.
//turnOn=true to turn on, turnOn=false to turn off
func (c *Client) turnOnOffServer(id string, turnOn bool) error {
	isOn, err := c.IsServerOn(id)
	if err != nil {
		return err
	}
	if isOn == turnOn {
		return nil
	}

	r := Request{
		uri:    path.Join(apiServerBase, id, "power"),
		method: http.MethodPatch,
		body: ServerPowerUpdateRequest{
			Power: turnOn,
		},
	}

	err = r.execute(*c, nil)
	if err != nil {
		return err
	}

	return c.WaitForServerPowerStatus(id, turnOn)
}

//StartServer starts a server
func (c *Client) StartServer(id string) error {
	return c.turnOnOffServer(id, true)
}

//StopServer stops a server
func (c *Client) StopServer(id string) error {
	return c.turnOnOffServer(id, false)
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
		return c.turnOnOffServer(id, false)
	}

	return nil
}

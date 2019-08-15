package gsclient

import (
	"net/http"
	"path"
)

type StorageList struct {
	List map[string]StorageProperties `json:"storages"`
}

type Storage struct {
	Properties StorageProperties `json:"storage"`
}

type StorageProperties struct {
	ChangeTime       string                    `json:"change_time"`
	LocationIata     string                    `json:"location_iata"`
	Status           string                    `json:"status"`
	LicenseProductNo int                       `json:"license_product_no"`
	LocationCountry  string                    `json:"location_country"`
	UsageInMinutes   int                       `json:"usage_in_minutes"`
	LastUsedTemplate string                    `json:"last_used_template"`
	CurrentPrice     float64                   `json:"current_price"`
	Capacity         int                       `json:"capacity"`
	LocationUuid     string                    `json:"location_uuid"`
	StorageType      string                    `json:"storage_type"`
	ParentUuid       string                    `json:"parent_uuid"`
	Name             string                    `json:"name"`
	LocationName     string                    `json:"location_name"`
	ObjectUuid       string                    `json:"object_uuid"`
	Snapshots        []StorageSnapshotRelation `json:"snapshots"`
	Relations        StorageRelations          `json:"relations"`
	Labels           []string                  `json:"labels"`
	CreateTime       string                    `json:"create_time"`
}

type StorageRelations struct {
	Servers           []StorageServer            `json:"servers"`
	SnapshotSchedules []StorageSnapshotSchedules `json:"snapshot_schedules"`
}

type StorageServer struct {
	Bootdevice bool   `json:"bootdevice"`
	Target     int    `json:"target"`
	Controller int    `json:"controller"`
	Bus        int    `json:"bus"`
	ObjectUuid string `json:"object_uuid"`
	Lun        int    `json:"lun"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
}

type StorageSnapshotRelation struct {
	LastUsedTemplate      string `json:"last_used_template"`
	ObjectUuid            string `json:"object_uuid"`
	StorageUuid           string `json:"storage_uuid"`
	SchedulesSnapshotName string `json:"schedules_snapshot_name"`
	SchedulesSnapshotUuid string `json:"schedules_snapshot_uuid"`
	ObjectCapacity        int    `json:"object_capacity"`
	CreateTime            string `json:"create_time"`
	ObjectName            string `json:"object_name"`
}

type StorageSnapshotSchedules struct {
	RunInterval   int    `json:"run_interval"`
	KeepSnapshots int    `json:"keep_snapshots"`
	ObjectName    string `json:"object_name"`
	NextRuntime   string `json:"next_runtime"`
	ObjectUuid    int    `json:"object_uuid"`
	Name          string `json:"name"`
	CreateTime    string `json:"create_time"`
}
type StorageTemplate struct {
	Sshkeys      []string `json:"sshkeys,omitempty"`
	TemplateUuid string   `json:"template_uuid,omitempty"`
	Password     string   `json:"password,omitempty"`
	PasswordType string   `json:"password_type,omitempty"`
	Hostname     string   `json:"hostname,omitempty"`
}

type StorageCreateRequest struct {
	Capacity     int             `json:"capacity"`
	LocationUuid string          `json:"location_uuid"`
	Name         string          `json:"name"`
	StorageType  string          `json:"storage_type,omitempty"`
	Template     StorageTemplate `json:"template,omitempty"`
	Labels       []string        `json:"labels,omitempty"`
}

type StorageUpdateRequest struct {
	Name     string   `json:"name,omitempty"`
	Labels   []string `json:"labels"`
	Capacity int      `json:"capacity,omitempty"`
}

type StorageEventList struct {
	List []StorageEventProperties `json:"events"`
}

type StorageEvent struct {
	Properties StorageEventProperties `json:"event"`
}

type StorageEventProperties struct {
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

//GetStorage get a storage
func (c *Client) GetStorage(id string) (Storage, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodGet,
	}
	var response Storage
	err := r.execute(*c, &response)
	return response, err
}

//GetStorageList gets a list of available storages
func (c *Client) GetStorageList() ([]Storage, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodGet,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{
			Properties: properties,
		})
	}
	return storages, err
}

//CreateStorage create a storage
func (c *Client) CreateStorage(body StorageCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return CreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//DeleteStorage delete a storage
func (c *Client) DeleteStorage(id string) error {
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateStorage update a storage
func (c *Client) UpdateStorage(id string, body StorageUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetStorageEventList get list of a storage's events
func (c *Client) GetStorageEventList(id string) ([]StorageEvent, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id, "events"),
		method: http.MethodGet,
	}
	var response StorageEventList
	var storageEvents []StorageEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storageEvents = append(storageEvents, StorageEvent{Properties: properties})
	}
	return storageEvents, err
}

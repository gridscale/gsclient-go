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

//StorageList JSON struct of a list of storages
type StorageList struct {
	List map[string]StorageProperties `json:"storages"`
}

//Storage JSON struct of a single storage
type Storage struct {
	Properties StorageProperties `json:"storage"`
}

//StorageProperties JSON struct of properties of a storage
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
	LocationUUID     string                    `json:"location_uuid"`
	StorageType      string                    `json:"storage_type"`
	ParentUUID       string                    `json:"parent_uuid"`
	Name             string                    `json:"name"`
	LocationName     string                    `json:"location_name"`
	ObjectUUID       string                    `json:"object_uuid"`
	Snapshots        []StorageSnapshotRelation `json:"snapshots"`
	Relations        StorageRelations          `json:"relations"`
	Labels           []string                  `json:"labels"`
	CreateTime       string                    `json:"create_time"`
}

//StorageRelations JSON struct of a list of a storage's relations
type StorageRelations struct {
	Servers           []StorageServerRelation              `json:"servers"`
	SnapshotSchedules []StorageAndSnapshotScheduleRelation `json:"snapshot_schedules"`
}

//StorageServerRelation JSON struct of a relation between a storage and a server
type StorageServerRelation struct {
	Bootdevice bool   `json:"bootdevice"`
	Target     int    `json:"target"`
	Controller int    `json:"controller"`
	Bus        int    `json:"bus"`
	ObjectUUID string `json:"object_uuid"`
	Lun        int    `json:"lun"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
}

//StorageSnapshotRelation JSON struct of a relation between a storage and a snapshot
type StorageSnapshotRelation struct {
	LastUsedTemplate      string `json:"last_used_template"`
	ObjectUUID            string `json:"object_uuid"`
	StorageUUID           string `json:"storage_uuid"`
	SchedulesSnapshotName string `json:"schedules_snapshot_name"`
	SchedulesSnapshotUUID string `json:"schedules_snapshot_uuid"`
	ObjectCapacity        int    `json:"object_capacity"`
	CreateTime            string `json:"create_time"`
	ObjectName            string `json:"object_name"`
}

//StorageAndSnapshotScheduleRelation JSON struct of a relation between a storage and a snapshot schedule
type StorageAndSnapshotScheduleRelation struct {
	RunInterval   int    `json:"run_interval"`
	KeepSnapshots int    `json:"keep_snapshots"`
	ObjectName    string `json:"object_name"`
	NextRuntime   string `json:"next_runtime"`
	ObjectUUID    int    `json:"object_uuid"`
	Name          string `json:"name"`
	CreateTime    string `json:"create_time"`
}

//StorageTemplate JSON struct of a storage template
type StorageTemplate struct {
	Sshkeys      []string `json:"sshkeys,omitempty"`
	TemplateUUID string   `json:"template_uuid"`
	Password     string   `json:"password,omitempty"`
	PasswordType string   `json:"password_type,omitempty"`
	Hostname     string   `json:"hostname,omitempty"`
}

//StorageCreateRequest JSON struct of a request for creating a storage
type StorageCreateRequest struct {
	Capacity     int              `json:"capacity"`
	LocationUUID string           `json:"location_uuid"`
	Name         string           `json:"name"`
	StorageType  string           `json:"storage_type,omitempty"`
	Template     *StorageTemplate `json:"template,omitempty"`
	Labels       []string         `json:"labels,omitempty"`
}

//StorageUpdateRequest JSON struct of a request for updating a storage
type StorageUpdateRequest struct {
	Name     string   `json:"name,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Capacity int      `json:"capacity,omitempty"`
}

//StorageEventList JSON struct of a list of a storage's events
type StorageEventList struct {
	List []StorageEventProperties `json:"events"`
}

//StorageEvent JSON struct of an event of a storage
type StorageEvent struct {
	Properties StorageEventProperties `json:"event"`
}

//StorageEventProperties JSON struct of properties of an event of a storage
type StorageEventProperties struct {
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

//GetStorage get a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorage
func (c *Client) GetStorage(ctx context.Context, id string) (Storage, error) {
	if !isValidUUID(id) {
		return Storage{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodGet,
	}
	var response Storage
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetStorageList gets a list of available storages
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorages
func (c *Client) GetStorageList(ctx context.Context) ([]Storage, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodGet,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{
			Properties: properties,
		})
	}
	return storages, err
}

//CreateStorage create a storage
//
//NOTE:
//
// - Allowed value for `StorageType`: nil, DefaultStorageType, HighStorageType, InsaneStorageType.
//
// - Allowed value for `PasswordType`: nil, PlainPasswordType, CryptPasswordType.
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createStorage
func (c *Client) CreateStorage(ctx context.Context, body StorageCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return CreateResponse{}, err
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

//DeleteStorage delete a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteStorage
func (c *Client) DeleteStorage(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
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
		return c.waitForStorageDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//UpdateStorage update a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateStorage
func (c *Client) UpdateStorage(ctx context.Context, id string, body StorageUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
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
		return c.waitForStorageActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetStorageEventList get list of a storage's event
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getStorageEvents
func (c *Client) GetStorageEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var storageEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storageEvents = append(storageEvents, StorageEvent{Properties: properties})
	}
	return storageEvents, err
}
<<<<<<< HEAD
=======

//GetStoragesByLocation gets a list of storages by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationStorages
func (c *Client) GetStoragesByLocation(ctx context.Context, id string) ([]Storage, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "storages"),
		method: http.MethodGet,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}

//GetDeletedStorages gets a list of deleted storages
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedStorages
func (c *Client) GetDeletedStorages(ctx context.Context) ([]Storage, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "storages"),
		method: http.MethodGet,
	}
	var response DeletedStorageList
	var storages []Storage
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}

//waitForStorageActive allows to wait until the storage's status is active
func (c *Client) waitForStorageActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		storage, err := c.GetStorage(ctx, id)
		return storage.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForStorageDeleted allows to wait until the storage is deleted
func (c *Client) waitForStorageDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiStorageBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

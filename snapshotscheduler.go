package gsclient

import (
	"net/http"
	"path"
)

type StorageSnapshotScheduleList struct {
	List map[string]StorageSnapshotSchedulerProperties `json:"snapshot_schedules"`
}

type StorageSnapshotSchedule struct {
	Properties StorageSnapshotSchedulerProperties `json:"snapshot_schedule"`
}

type StorageSnapshotSchedulerProperties struct {
	ChangeTime    string                            `json:"change_time"`
	CreateTime    string                            `json:"create_time"`
	KeepSnapshots int                               `json:"keep_snapshots"`
	Labels        []string                          `json:"labels"`
	Name          string                            `json:"name"`
	NextRuntime   string                            `json:"next_runtime"`
	ObjectUuid    string                            `json:"object_uuid"`
	Relations     StorageSnapshotSchedulerRelations `json:"relations"`
	RunInterval   int                               `json:"run_interval"`
	Status        string                            `json:"status"`
	StorageUuid   string                            `json:"storage_uuid"`
}

type StorageSnapshotSchedulerRelations struct {
	Snapshots []StorageSnapshotSchedulerRelation `json:"snapshots"`
}

type StorageSnapshotSchedulerRelation struct {
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
	ObjectUuid string `json:"object_uuid"`
}

type StorageSnapshotScheduleCreateRequest struct {
	Name          string   `json:"name"`
	Labels        []string `json:"labels"`
	RunInterval   int      `json:"run_interval"`
	KeepSnapshots int      `json:"keep_snapshots"`
	NextRuntime   string   `json:"next_runtime"`
}

type StorageSnapshotScheduleCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

type StorageSnapshotScheduleUpdateRequest struct {
	Name          string   `json:"name"`
	Labels        []string `json:"labels"`
	RunInterval   int      `json:"run_interval"`
	KeepSnapshots int      `json:"keep_snapshots"`
	NextRuntime   string   `json:"next_runtime"`
}

//GetStorageSnapshotScheduleList gets a list of available storage snapshot schedules based on a given storage's id
func (c *Client) GetStorageSnapshotScheduleList(id string) ([]StorageSnapshotSchedule, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshot_schedules"),
		method: http.MethodGet,
	}
	var response StorageSnapshotScheduleList
	var list []StorageSnapshotSchedule
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, StorageSnapshotSchedule{Properties: properties})
	}
	return list, err
}

//GetStorageSnapshotSchedule gets a specific storage snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) GetStorageSnapshotSchedule(storageId, scheduleId string) (StorageSnapshotSchedule, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshot_schedules", scheduleId),
		method: http.MethodGet,
	}
	var response StorageSnapshotSchedule
	err := r.execute(*c, &response)
	return response, err
}

//CreateStorageSnapshotSchedule create a storage's snapshot scheduler
func (c *Client) CreateStorageSnapshotSchedule(id string, body StorageSnapshotScheduleCreateRequest) (
	StorageSnapshotScheduleCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshot_schedules"),
		method: http.MethodPost,
		body:   body,
	}
	var response StorageSnapshotScheduleCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return StorageSnapshotScheduleCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//UpdateStorageSnapshotSchedule updates specific Storage's snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) UpdateStorageSnapshotSchedule(storageId, scheduleId string,
	body StorageSnapshotScheduleUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshot_schedules", scheduleId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteStorageSnapshotSchedule deletes specific Storage's snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) DeleteStorageSnapshotSchedule(storageId, scheduleId string) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshot_schedules", scheduleId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

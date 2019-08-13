package gsclient

import (
	"net/http"
	"path"
)

//StorageSnapshotList is JSON structure of a list of storage snapshots
type StorageSnapshotList struct {
	List map[string]StorageSnapshotProperties `json:"snapshots"`
}

//StorageSnapshotSingle is JSON structure of a single storage snapshot
type StorageSnapshotSingle struct {
	Properties StorageSnapshotProperties `json:"snapshot"`
}

type StorageSnapshotProperties struct {
	Labels           []string `json:"labels"`
	ObjectUuid       string   `json:"object_uuid"`
	Name             string   `json:"name"`
	Status           string   `json:"status"`
	LocationCountry  string   `json:"location_country"`
	UsageInMinutes   int      `json:"usage_in_minutes"`
	LocationUuid     string   `json:"location_uuid"`
	ChangeTime       string   `json:"change_time"`
	LicenseProductNo int      `json:"license_product_no"`
	CurrentPrice     float64  `json:"current_price"`
	CreateTime       string   `json:"create_time"`
	Capacity         int      `json:"capacity"`
	LocationName     string   `json:"location_name"`
	LocationIata     string   `json:"location_iata"`
	ParentUuid       string   `json:"parent_uuid"`
}

type StorageSnapshotCreateRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

type StorageSnapshotCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

type StorageSnapshotUpdateRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

type StorageRollbackRequest struct {
	Rollback bool `json:"rollback"`
}

type StorageSnapshotExportToS3Request struct {
	S3auth struct {
		Host       string `json:"host"`
		AccessKeys string `json:"access_keys"`
		SecretKey  string `json:"secret_key"`
	} `json:"s3auth"`
	S3data struct {
		Host     string `json:"host"`
		Bucket   string `json:"bucket"`
		Filename string `json:"filename"`
		Private  bool   `json:"private"`
	} `json:"s3data"`
}

//GetStorageSnapshotList gets a list of storage snapshots
func (c *Client) GetStorageSnapshotList(id string) ([]StorageSnapshotSingle, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodGet,
	}
	var response StorageSnapshotList
	var list []StorageSnapshotSingle
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, StorageSnapshotSingle{Properties: properties})
	}
	return list, err
}

//GetStorageSnapshot gets a specific storage's snapshot based on given storage id and snapshot id.
func (c *Client) GetStorageSnapshot(storageId, snapshotId string) (StorageSnapshotSingle, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshots", snapshotId),
		method: http.MethodGet,
	}
	var response StorageSnapshotSingle
	err := r.execute(*c, &response)
	return response, err
}

//CreateStorageSnapshot creates a new storage's snapshot
func (c *Client) CreateStorageSnapshot(id string, body StorageSnapshotCreateRequest) (StorageSnapshotCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodPost,
		body:   body,
	}
	var response StorageSnapshotCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return StorageSnapshotCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//UpdateStorageSnapshot updates a specific storage's snapshot
func (c *Client) UpdateStorageSnapshot(storageId, snapshotId string, body StorageSnapshotUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshots", snapshotId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteStorageSnapshot deletes a specific storage's snapshot
func (c *Client) DeleteStorageSnapshot(storageId, snapshotId string) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshots", snapshotId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//RollbackStorage rollbacks a storage
func (c *Client) RollbackStorage(storageId, snapshotId string, body StorageRollbackRequest) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshots", snapshotId, "rollback"),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//ExportStorageSnapshotToS3 export a storage's snapshot to S3
func (c *Client) ExportStorageSnapshotToS3(storageId, snapshotId string, body StorageSnapshotExportToS3Request) error {
	r := Request{
		uri:    path.Join(apiStorageBase, storageId, "snapshots", snapshotId, "export_to_s3"),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

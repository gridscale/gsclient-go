package gsclient

import (
	"net/http"
	"path"
)

//ObjectStorageAccessKeyList is JSON structure of a list of Object Storage Access Keys
type ObjectStorageAccessKeyList struct {
	List []ObjectStorageAccessKeyProperties `json:"access_keys"`
}

//ObjectStorageAccessKey is JSON structure of a single Object Storage Access Key
type ObjectStorageAccessKey struct {
	Properties ObjectStorageAccessKeyProperties `json:"access_key"`
}

type ObjectStorageAccessKeyProperties struct {
	SecretKey string `json:"secret_key"`
	AccessKey string `json:"access_key"`
	User      string `json:"user"`
}

type ObjectStorageAccessKeyResponse struct {
	AccessKey struct {
		SecretKey string `json:"secret_key"`
		AccessKey string `json:"access_key"`
	} `json:"access_key"`
	RequestUuid string `json:"request_uuid"`
}

type ObjectStorageBucketList struct {
	List []ObjectStorageBucketProperties `json:"buckets"`
}

type ObjectStorageBucket struct {
	Properties ObjectStorageBucketProperties `json:"bucket"`
}

type ObjectStorageBucketProperties struct {
	Name  string `json:"name"`
	Usage struct {
		SizeKb     int `json:"size_kb"`
		NumObjects int `json:"num_objects"`
	} `json:"usage"`
}

//GetObjectStorageAccessKeyList gets a list of available object storage access keys
func (c *Client) GetObjectStorageAccessKeyList() ([]ObjectStorageAccessKey, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKeyList
	var list []ObjectStorageAccessKey
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ObjectStorageAccessKey{Properties: properties})
	}
	return list, err
}

//GetObjectStorageAccessKey gets a specific object storage access key based on given id
func (c *Client) GetObjectStorageAccessKey(id string) (ObjectStorageAccessKey, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKey
	err := r.execute(*c, &response)
	return response, err
}

//CreateObjectStorageAccessKey creates an object storage access key
func (c *Client) CreateObjectStorageAccessKey() (ObjectStorageAccessKeyResponse, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodPost,
	}
	var response ObjectStorageAccessKeyResponse
	err := r.execute(*c, &response)
	if err != nil {
		return ObjectStorageAccessKeyResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUuid)
	return response, err
}

//DeleteObjectStorageAccessKey deletes a specific object storage access key based on given id
func (c *Client) DeleteObjectStorageAccessKey(id string) error {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetObjectStorageBucketList gets a list of object storage buckets
func (c *Client) GetObjectStorageBucketList() ([]ObjectStorageBucket, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "buckets"),
		method: http.MethodGet,
	}
	var response ObjectStorageBucketList
	var list []ObjectStorageBucket
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		list = append(list, ObjectStorageBucket{Properties:properties})
	}
	return list, err
}

package gsclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

//ObjectStorageAccessKeyList is JSON structure of a list of Object Storage Access Keys
type ObjectStorageAccessKeyList struct {
	//Array of Object Storages' access keys
	List []ObjectStorageAccessKeyProperties `json:"access_keys"`
}

//ObjectStorageAccessKey is JSON structure of a single Object Storage Access Key
type ObjectStorageAccessKey struct {
	//Properties of an object storage access key
	Properties ObjectStorageAccessKeyProperties `json:"access_key"`
}

//ObjectStorageAccessKeyProperties is JSON struct of properties of an object storage access key
type ObjectStorageAccessKeyProperties struct {
	//The object storage secret_key.
	SecretKey string `json:"secret_key"`

	//The object storage access_key.
	AccessKey string `json:"access_key"`

	//Account this credentials belong to.
	User string `json:"user"`
}

//ObjectStorageAccessKeyCreateResponse is JSON struct of a response for creating an object storage access key
type ObjectStorageAccessKeyCreateResponse struct {
	AccessKey struct {
		////The object storage secret_key.
		SecretKey string `json:"secret_key"`

		//The object storage secret_key.
		AccessKey string `json:"access_key"`
	} `json:"access_key"`

	//UUID of the request
	RequestUUID string `json:"request_uuid"`
}

//ObjectStorageBucketList is JSON struct of a list of buckets
type ObjectStorageBucketList struct {
	//Array of Buckets
	List []ObjectStorageBucketProperties `json:"buckets"`
}

//ObjectStorageBucket is JSON struct of a single bucket
type ObjectStorageBucket struct {
	//Properties of a bucket
	Properties ObjectStorageBucketProperties `json:"bucket"`
}

//ObjectStorageBucketProperties is JSON struct of properties of a bucket
type ObjectStorageBucketProperties struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The current usage of the bucket.
	Usage struct {
		//The size of the the bucket (in kb).
		SizeKb int `json:"size_kb"`

		//The number of files in the bucket.
		NumObjects int `json:"num_objects"`
	} `json:"usage"`
}

//GetObjectStorageAccessKeyList gets a list of available object storage access keys
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getAccessKeys
func (c *Client) GetObjectStorageAccessKeyList(ctx context.Context) ([]ObjectStorageAccessKey, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKeyList
	var accessKeys []ObjectStorageAccessKey
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		accessKeys = append(accessKeys, ObjectStorageAccessKey{Properties: properties})
	}
	return accessKeys, err
}

//GetObjectStorageAccessKey gets a specific object storage access key based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getAccessKey
func (c *Client) GetObjectStorageAccessKey(ctx context.Context, id string) (ObjectStorageAccessKey, error) {
	if strings.TrimSpace(id) == "" {
		return ObjectStorageAccessKey{}, errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKey
	err := r.execute(ctx, *c, &response)
	return response, err
}

//CreateObjectStorageAccessKey creates an object storage access key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createAccessKey
func (c *Client) CreateObjectStorageAccessKey(ctx context.Context) (ObjectStorageAccessKeyCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodPost,
	}
	var response ObjectStorageAccessKeyCreateResponse
	err := r.execute(ctx, *c, &response)
	if err != nil {
		return ObjectStorageAccessKeyCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(ctx, response.RequestUUID)
	}
	return response, err
}

//DeleteObjectStorageAccessKey deletes a specific object storage access key based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteAccessKey
func (c *Client) DeleteObjectStorageAccessKey(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(ctx, *c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForObjectStorageAccessKeyDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
}

//GetObjectStorageBucketList gets a list of object storage buckets
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getBuckets
func (c *Client) GetObjectStorageBucketList(ctx context.Context) ([]ObjectStorageBucket, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "buckets"),
		method: http.MethodGet,
	}
	var response ObjectStorageBucketList
	var buckets []ObjectStorageBucket
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		buckets = append(buckets, ObjectStorageBucket{Properties: properties})
	}
	return buckets, err
}

//waitForObjectStorageAccessKeyDeleted allows to wait until the object storage's access key is deleted
func (c *Client) waitForObjectStorageAccessKeyDeleted(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("'id' is required")
	}
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			errorMessage := fmt.Sprintf("Timeout reached when waiting for access key %v to be deleted", id)
			c.cfg.logger.Error(errorMessage)
			return errors.New(errorMessage)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			r := Request{
				uri:          path.Join(apiObjectStorageBase, "access_keys", id),
				method:       http.MethodGet,
				skipPrint404: true,
			}
			err := r.execute(ctx, *c, nil)
			if err != nil {
				if requestError, ok := err.(RequestError); ok {
					if requestError.StatusCode == 404 {
						return nil
					}
				}
				return err
			}
		}
	}
}

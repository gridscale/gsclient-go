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

//SshkeyList JSON struct of a list of SSH-keys
type SshkeyList struct {
	List map[string]SshkeyProperties `json:"sshkeys"`
}

//Sshkey JSON struct of a single SSH-key
type Sshkey struct {
	Properties SshkeyProperties `json:"sshkey"`
}

//SshkeyProperties JSON struct of properties of a single SSH-key
type SshkeyProperties struct {
	Name       string   `json:"name"`
	ObjectUUID string   `json:"object_uuid"`
	Status     string   `json:"status"`
	CreateTime string   `json:"create_time"`
	ChangeTime string   `json:"change_time"`
	Sshkey     string   `json:"sshkey"`
	Labels     []string `json:"labels"`
	UserUUID   string   `json:"user_uuid"`
}

//SshkeyCreateRequest JSON struct of a request for creating a SSH-key
type SshkeyCreateRequest struct {
	Name   string   `json:"name"`
	Sshkey string   `json:"sshkey"`
	Labels []string `json:"labels,omitempty"`
}

//SshkeyUpdateRequest JSON struct of a request for updating a SSH-key
type SshkeyUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Sshkey string   `json:"sshkey,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//SshkeyEventList JSON struct of a list of a SSH-key's events
type SshkeyEventList struct {
	List []SshkeyEventProperties `json:"events"`
}

//SshkeyEvent JSON struct of an event of a SSH-key
type SshkeyEvent struct {
	Properties SshkeyEventProperties `json:"event"`
}

//SshkeyEventProperties JSON struct of properties of an event of a SSH-key
type SshkeyEventProperties struct {
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

//GetSshkey gets a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKey
func (c *Client) GetSshkey(ctx context.Context, id string) (Sshkey, error) {
	if !isValidUUID(id) {
		return Sshkey{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodGet,
	}
	var response Sshkey
	err := r.execute(ctx, *c, &response)
	return response, err
}

//GetSshkeyList gets a list of ssh keys
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeys
func (c *Client) GetSshkeyList(ctx context.Context) ([]Sshkey, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: http.MethodGet,
	}

	var response SshkeyList
	var sshKeys []Sshkey
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshKeys = append(sshKeys, Sshkey{Properties: properties})
	}
	return sshKeys, err
}

//CreateSshkey creates a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createSshKey
func (c *Client) CreateSshkey(ctx context.Context, body SshkeyCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: "POST",
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

//DeleteSshkey deletes a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteSshKey
func (c *Client) DeleteSshkey(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
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
		return c.waitForSSHKeyDeleted(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//UpdateSshkey updates a ssh key
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateSshKey
func (c *Client) UpdateSshkey(ctx context.Context, id string, body SshkeyUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
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
		return c.waitForSSHKeyActive(ctx, id)
	}
	return r.execute(ctx, *c, nil)
>>>>>>> 8d4aa0e... add `context`
}

//GetSshkeyEventList gets a ssh key's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeyEvents
func (c *Client) GetSshkeyEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiSshkeyBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var sshEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshEvents = append(sshEvents, SshkeyEvent{Properties: properties})
	}
	return sshEvents, err
}
<<<<<<< HEAD
=======

//waitForSSHKeyActive allows to wait until the SSH-Key's status is active
func (c *Client) waitForSSHKeyActive(ctx context.Context, id string) error {
	return retryWithTimeout(func() (bool, error) {
		key, err := c.GetSshkey(ctx, id)
		return key.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForSSHKeyDeleted allows to wait until the SSH-Key is deleted
func (c *Client) waitForSSHKeyDeleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiSshkeyBase, id)
	method := http.MethodGet
	return c.waitFor404Status(ctx, uri, method)
}
>>>>>>> 8d4aa0e... add `context`

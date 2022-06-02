package gsclient

import (
	"context"
	"errors"
	"net/http"
	"path"
)

// SSHKeyOperator provides an interface for operations on SSH keys.
type SSHKeyOperator interface {
	GetSSHKey(ctx context.Context, id string) (SSHKey, error)
	GetSSHKeyList(ctx context.Context) ([]SSHKey, error)
	CreateSSHKey(ctx context.Context, body SSHKeyCreateRequest) (CreateResponse, error)
	DeleteSSHKey(ctx context.Context, id string) error
	UpdateSSHKey(ctx context.Context, id string, body SSHKeyUpdateRequest) error
	GetSSHKeyEventList(ctx context.Context, id string) ([]Event, error)
}

// SSHKeyList holds a list of SSH keys.
type SSHKeyList struct {
	// Array of SSH keys.
	List map[string]SSHKeyProperties `json:"sshkeys"`
}

// SSHKey represents a single SSH key.
type SSHKey struct {
	// Properties of a SSH key.
	Properties SSHKeyProperties `json:"sshkey"`
}

// SSHKeyProperties holds properties of a single SSH key.
// A SSH key can be retrieved when creating new storages and attaching them to
// servers.
type SSHKeyProperties struct {
	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1).
	SSHKey string `json:"sshkey"`

	// List of labels.
	Labels []string `json:"labels"`

	// The User-UUID of the account which created this SSH Key.
	UserUUID string `json:"user_uuid"`
}

// SSHKeyCreateRequest represents a request for creating a SSH key.
type SSHKeyCreateRequest struct {
	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1).
	SSHKey string `json:"sshkey"`

	// List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

// SSHKeyUpdateRequest represents a request for updating a SSH key.
type SSHKeyUpdateRequest struct {
	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	// Optional.
	Name string `json:"name,omitempty"`

	// The OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1). Optional.
	SSHKey string `json:"sshkey,omitempty"`

	// List of labels. Optional.
	Labels *[]string `json:"labels,omitempty"`
}

// GetSSHKey gets a single SSH key object.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKey
func (c *Client) GetSSHKey(ctx context.Context, id string) (SSHKey, error) {
	if !isValidUUID(id) {
		return SSHKey{}, errors.New("'id' is invalid")
	}
	r := gsRequest{
		uri:                 path.Join(apiSSHKeyBase, id),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response SSHKey
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetSSHKeyList gets the list of SSH keys in the project.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeys
func (c *Client) GetSSHKeyList(ctx context.Context) ([]SSHKey, error) {
	r := gsRequest{
		uri:                 apiSSHKeyBase,
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}

	var response SSHKeyList
	var sshKeys []SSHKey
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshKeys = append(sshKeys, SSHKey{Properties: properties})
	}
	return sshKeys, err
}

// CreateSSHKey creates a new SSH key.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/createSshKey
func (c *Client) CreateSSHKey(ctx context.Context, body SSHKeyCreateRequest) (CreateResponse, error) {
	r := gsRequest{
		uri:    apiSSHKeyBase,
		method: "POST",
		body:   body,
	}
	var response CreateResponse
	err := r.execute(ctx, *c, &response)
	return response, err
}

// DeleteSSHKey removes a single SSH key.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/deleteSshKey
func (c *Client) DeleteSSHKey(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := gsRequest{
		uri:    path.Join(apiSSHKeyBase, id),
		method: http.MethodDelete,
	}
	return r.execute(ctx, *c, nil)
}

// UpdateSSHKey updates a SSH key.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/updateSshKey
func (c *Client) UpdateSSHKey(ctx context.Context, id string, body SSHKeyUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := gsRequest{
		uri:    path.Join(apiSSHKeyBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(ctx, *c, nil)
}

// GetSSHKeyEventList gets a SSH key's events.
//
// See: https://gridscale.io/en//api-documentation/index.html#operation/getSshKeyEvents
func (c *Client) GetSSHKeyEventList(ctx context.Context, id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := gsRequest{
		uri:                 path.Join(apiSSHKeyBase, id, "events"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
	}
	var response EventList
	var sshEvents []Event
	err := r.execute(ctx, *c, &response)
	for _, properties := range response.List {
		sshEvents = append(sshEvents, Event{Properties: properties})
	}
	return sshEvents, err
}

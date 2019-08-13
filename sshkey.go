package gsclient

import (
	"net/http"
	"path"
)

type SshkeyList struct {
	List map[string]SshkeyProperties `json:"sshkeys"`
}

type Sshkey struct {
	Properties SshkeyProperties `json:"sshkey"`
}

type SshkeyProperties struct {
	Name       string   `json:"name"`
	ObjectUuid string   `json:"object_uuid"`
	Status     string   `json:"status"`
	CreateTime string   `json:"create_time"`
	ChangeTime string   `json:"change_time"`
	Sshkey     string   `json:"sshkey"`
	Labels     []string `json:"labels"`
	UserUuid   string   `json:"user_uuid"`
}

type SshkeyCreateRequest struct {
	Name   string   `json:"name"`
	Sshkey string   `json:"sshkey"`
	Labels []string `json:"labels,omitempty"`
}

type SshkeyUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Sshkey string   `json:"sshkey,omitempty"`
	Labels []string `json:"labels"`
}

type SshkeyEventList struct {
	List []SshkeyEventProperties `json:"events"`
}

type SshkeyEvent struct {
	Properties SshkeyEventProperties `json:"event"`
}

type SshkeyEventProperties struct {
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

//GetSshkey gets a ssh key
func (c *Client) GetSshkey(id string) (Sshkey, error) {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodGet,
	}
	var response Sshkey
	err := r.execute(*c, &response)
	return response, err
}

//GetSshkeyList gets a list of ssh keys
func (c *Client) GetSshkeyList() ([]Sshkey, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: http.MethodGet,
	}

	var response SshkeyList
	var sshKeys []Sshkey
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		sshKeys = append(sshKeys, Sshkey{Properties: properties})
	}
	return sshKeys, err
}

//CreateSshkey creates a ssh key
func (c *Client) CreateSshkey(body SshkeyCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiSshkeyBase,
		method: "POST",
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

//DeleteSshkey deletes a ssh key
func (c *Client) DeleteSshkey(id string) error {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateSshkey updates a ssh key
func (c *Client) UpdateSshkey(id string, body SshkeyUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetSshkeyEventList gets a ssh key's events
func (c *Client) GetSshkeyEventList(id string) ([]SshkeyEvent, error) {
	r := Request{
		uri:    path.Join(apiSshkeyBase, id, "events"),
		method: http.MethodGet,
	}
	var response SshkeyEventList
	var sshEvents []SshkeyEvent
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		sshEvents = append(sshEvents, SshkeyEvent{Properties:properties})
	}
	return sshEvents, err
}
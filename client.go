package gsclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"time"
)

const (
	requestBase          = "/requests/"
	apiServerBase        = "/objects/servers"
	apiStorageBase       = "/objects/storages"
	apiNetworkBase       = "/objects/networks"
	apiIPBase            = "/objects/ips"
	apiSshkeyBase        = "/objects/sshkeys"
	apiTemplateBase      = "/objects/templates"
	apiLoadBalancerBase  = "/objects/loadbalancers"
	apiPaaSBase          = "/objects/paas"
	apiISOBase           = "/objects/isoimages"
	apiObjectStorageBase = "/objects/objectstorages"
	apiFirewallBase      = "/objects/firewalls"
	apiLocationBase      = "/objects/locations"
	apiEventBase         = "/objects/events"
	apiLabelBase         = "/objects/labels"
	apiDeletedBase       = "/objects/deleted"
)

//Client struct of a gridscale golang client
type Client struct {
	cfg *Config
}

//NewClient creates new gridscale golang client
func NewClient(c *Config) *Client {
	client := &Client{
		cfg: c,
	}
	return client
}

//getLogger returns logger
func (c *Client) getLogger() logrus.Logger {
	return c.cfg.logger
}

//getHttpClient returns http client
func (c *Client) getHttpClient() *http.Client {
	return c.cfg.httpClient
}

//isSynchronous returns if the client is sync or not
func (c *Client) isSynchronous() bool {
	return c.cfg.sync
}

//getRequestCheckTimeout returns request check timeout
func (c *Client) getRequestCheckTimeout() time.Duration {
	return c.cfg.requestCheckTimeoutSecs
}

//getDelayInterval returns request delay interval
func (c *Client) getDelayInterval() time.Duration {
	return c.cfg.delayInterval
}

//getMaxNumberOfRetries returns max number of retries
func (c *Client) getMaxNumberOfRetries() int {
	return c.cfg.maxNumberOfRetries
}

//getAPIURL returns api URL
func (c *Client) getAPIURL() string {
	return c.cfg.apiURL
}

//getUserAgent returns user agent
func (c *Client) getUserAgent() string {
	return c.cfg.userAgent
}

//getUserUUID returns user UUID
func (c *Client) getUserUUID() string {
	return c.cfg.userUUID
}

//getAPIToken returns api token
func (c *Client) getAPIToken() string {
	return c.cfg.apiToken
}

//waitForRequestCompleted allows to wait for a request to complete
func (c *Client) waitForRequestCompleted(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	logger := c.getLogger()
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:    path.Join(requestBase, id),
			method: "GET",
		}
		var response RequestStatus
		err := r.execute(ctx, *c, &response)
		if err != nil {
			return false, err
		}
		if response[id].Status == requestDoneStatus {
			logger.Info("Done with creating")
			return false, nil
		} else if response[id].Status == requestFailStatus {
			errMessage := fmt.Sprintf("request %s failed with error %s", id, response[id].Message)
			return false, errors.New(errMessage)
		}
		return true, nil
	}, c.getRequestCheckTimeout(), c.getDelayInterval())
}

//waitFor404Status waits until server returns 404 status code
func (c *Client) waitFor404Status(ctx context.Context, uri, method string) error {
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:          uri,
			method:       method,
			skipPrint404: true,
		}
		err := r.execute(ctx, *c, nil)
		if err != nil {
			if requestError, ok := err.(RequestError); ok {
				if requestError.StatusCode == 404 {
					return false, nil
				}
			}
			return false, err
		}
		return true, nil
	}, c.getRequestCheckTimeout(), c.getDelayInterval())
}

//waitFor200Status waits until server returns 200 (OK) status code
func (c *Client) waitFor200Status(ctx context.Context, uri, method string) error {
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:          uri,
			method:       method,
			skipPrint404: true,
		}
		err := r.execute(ctx, *c, nil)
		if err != nil {
			if requestError, ok := err.(RequestError); ok {
				if requestError.StatusCode == 404 {
					return true, nil
				}
			}
			return false, err
		}
		return false, nil
	}, c.getRequestCheckTimeout(), c.getDelayInterval())
}

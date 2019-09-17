package gsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

//Request gridscale's custom request struct
type Request struct {
	uri    string
	method string
	body   interface{}
}

//CreateResponse common struct of a response for creation
type CreateResponse struct {
	ObjectUUID  string `json:"object_uuid"`
	RequestUUID string `json:"request_uuid"`
}

//RequestStatus status of a request
type RequestStatus map[string]RequestStatusProperties

//RequestStatusProperties JSON struct of properties of a request's status
type RequestStatusProperties struct {
	Status     string   `json:"status"`
	Message    string   `json:"message"`
	CreateTime JSONTime `json:"create_time"`
}

//RequestError error of a request
type RequestError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusCode  int
}

//Error just returns error as string
func (r RequestError) Error() string {
	message := r.Description
	if message == "" {
		message = "no error message received from server"
	}
	return fmt.Sprintf("statuscode %v returned: %s", r.StatusCode, message)
}

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(c Client, output interface{}) error {
	url := c.cfg.apiURL + r.uri
	c.cfg.logger.Debugf("%v request sent to URL: %v", r.method, url)

	//Convert the body of the request to json
	jsonBody := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(jsonBody).Encode(r.body)
		if err != nil {
			return err
		}
	}

	//Add authentication headers and content type
	request, err := http.NewRequest(r.method, url, jsonBody)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", c.cfg.userAgent)
	request.Header.Add("X-Auth-UserID", c.cfg.userUUID)
	request.Header.Add("X-Auth-Token", c.cfg.apiToken)
	request.Header.Add("Content-Type", "application/json")
	c.cfg.logger.Debugf("Request body: %v", request.Body)

	retryNo := 0
	maxNumOfRetries := c.cfg.maxNumberOfRetries
	delayInterval := c.cfg.delayInterval
	var latestRetryErr error
RETRY:
	for retryNo <= maxNumOfRetries {
		//execute the request
		result, err := c.cfg.httpClient.Do(request)
		if err != nil {
			c.cfg.logger.Errorf("Error while executing the request: %v", err)
			return err
		}

		iostream, err := ioutil.ReadAll(result.Body)
		if err != nil {
			c.cfg.logger.Errorf("Error while reading the response's body: %v", err)
			return err
		}

		c.cfg.logger.Debugf("Status code returned: %v", result.StatusCode)

		if result.StatusCode >= 300 {
			var errorMessage RequestError //error messages have a different structure, so they are read with a different struct
			errorMessage.StatusCode = result.StatusCode
			json.Unmarshal(iostream, &errorMessage)
			if result.StatusCode >= 500 {
				latestRetryErr = errorMessage
				time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
				retryNo++
				c.cfg.logger.Errorf("RETRY no %d ! Error message: %v. Title: %v. Code: %v.", retryNo, errorMessage.Description, errorMessage.Title, errorMessage.StatusCode)
				continue RETRY //continue the RETRY loop
			}
			c.cfg.logger.Errorf("Error message: %v. Title: %v. Code: %v.", errorMessage.Description, errorMessage.Title, errorMessage.StatusCode)
			return errorMessage
		}
		c.cfg.logger.Debugf("Response body: %v", string(iostream))
		//if output is set
		if output != nil {
			err = json.Unmarshal(iostream, output) //Edit the given struct
			if err != nil {
				c.cfg.logger.Errorf("Error while marshaling JSON: %v", err)
				return err
			}
		}
		return nil
	}
	return latestRetryErr
}

//WaitForRequestCompletion allows to wait for a request to complete. Timeouts are currently hardcoded
func (c *Client) WaitForRequestCompletion(id string) error {
	r := Request{
		uri:    path.Join("/requests/", id),
		method: "GET",
	}
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			c.cfg.logger.Errorf("Timeout reached when waiting for request %v to complete", id)
			return fmt.Errorf("Timeout reached when waiting for request %v to complete", id)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			var response RequestStatus
			r.execute(*c, &response)
			if response[id].Status == "done" {
				c.cfg.logger.Info("Done with creating")
				return nil
			}
		}
	}
}

//WaitForServerPowerStatus  allows to wait for a server changing its power status. Timeouts are currently hardcoded
func (c *Client) WaitForServerPowerStatus(id string, status bool) error {
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			c.cfg.logger.Errorf("Timeout reached when trying to shut down system with id %v", id)
			return fmt.Errorf("Timeout reached when trying to shut down system with id %v", id)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			server, err := c.GetServer(id)
			if err != nil {
				return err
			}
			if server.Properties.Power == status {
				c.cfg.logger.Infof("The power status of the server with id %v has changed to %t", id, status)
				return nil
			}
		}
	}
}

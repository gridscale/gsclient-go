package gsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Request gridscale's custom request struct
type Request struct {
	uri          string
	method       string
	skipPrint404 bool
	body         interface{}
}

//CreateResponse common struct of a response for creation
type CreateResponse struct {
	//UUID of the object being created
	ObjectUUID string `json:"object_uuid"`

	//UUID of the request
	RequestUUID string `json:"request_uuid"`
}

//RequestStatus status of a request
type RequestStatus map[string]RequestStatusProperties

//RequestStatusProperties JSON struct of properties of a request's status
type RequestStatusProperties struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime GSTime `json:"create_time"`
}

//RequestError error of a request
type RequestError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusCode  int
	RequestUUID string
}

//Error just returns error as string
func (r RequestError) Error() string {
	message := r.Description
	if message == "" {
		message = "no error message received from server"
	}
	errorMessageFormat := "Status code: %v. Error: %s. Request UUID: %s. "
	if r.StatusCode >= 500 {
		errorMessageFormat += "Please report this error along with the request UUID."
	}
	return fmt.Sprintf(errorMessageFormat, r.StatusCode, message, r.RequestUUID)
}

const requestUUIDHeaderParam = "X-Request-Id"

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(ctx context.Context, c Client, output interface{}) error {
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
	request = request.WithContext(ctx)
	request.Header.Set("User-Agent", c.cfg.userAgent)
	request.Header.Add("X-Auth-UserID", c.cfg.userUUID)
	request.Header.Add("X-Auth-Token", c.cfg.apiToken)
	request.Header.Add("Content-Type", "application/json")
	c.cfg.logger.Debugf("Request body: %v", request.Body)
	return retryWithLimitedNumOfRetries(func() (bool, error) {
		//execute the request
		result, err := c.cfg.httpClient.Do(request)
		if err != nil {
			c.cfg.logger.Errorf("Error while executing the request: %v", err)
			return false, err
		}
		statusCode := result.StatusCode
		requestUUID := result.Header.Get(requestUUIDHeaderParam)
		iostream, err := ioutil.ReadAll(result.Body)
		if err != nil {
			c.cfg.logger.Errorf("Error while reading the response's body: %v", err)
			return false, err
		}

		c.cfg.logger.Debugf("Status code: %v. Request UUID: %v.", statusCode, requestUUID)

		if result.StatusCode >= 300 {
			var errorMessage RequestError //error messages have a different structure, so they are read with a different struct
			errorMessage.StatusCode = statusCode
			errorMessage.RequestUUID = requestUUID
			json.Unmarshal(iostream, &errorMessage)
			//if internal server error OR object is in status that does not allow the request, retry
			if result.StatusCode >= 500 || result.StatusCode == 424 {
				return true, errorMessage
			}
			if r.skipPrint404 && result.StatusCode == 404 {
				c.cfg.logger.Debug("Skip 404 error code.")
				return false, errorMessage
			}
			c.cfg.logger.Errorf(
				"Error message: %v. Title: %v. Code: %v. Request UUID: %v.",
				errorMessage.Description,
				errorMessage.Title,
				errorMessage.StatusCode,
				errorMessage.RequestUUID,
			)
			return false, errorMessage
		}
		c.cfg.logger.Debugf("Response body: %v", string(iostream))
		//if output is set
		if output != nil {
			err = json.Unmarshal(iostream, output) //Edit the given struct
			if err != nil {
				c.cfg.logger.Errorf("Error while marshaling JSON: %v", err)
				return false, err
			}
		}
		return false, nil
	}, c.cfg.maxNumberOfRetries, c.cfg.delayInterval)
}

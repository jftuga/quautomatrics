package rest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// used when attempting a network operation
const maxAttempts int = 240
const attemptDelay int = 60

type Rest struct {
	token   string
	BaseURL string
}

// New - return a Rest struct containing the API token and API URL
func New(token, datacenter string) *Rest {
	url := strings.Replace(DataCenterTemplateURL, "{DATACENTER}", datacenter, -1)
	return &Rest{
		token:   token,
		BaseURL: url,
	}
}

// action - a backend, generic function called by Get and Delete
func (rest Rest) action(verb, path string) string {
	var body []byte
	success := false
	for attempts := 0; attempts <= maxAttempts; attempts++ {
		url := rest.BaseURL + path
		req, err := http.NewRequest(verb, url, nil)
		if err != nil {
			log.Fatalf("Error #32043: %s:%s", path, err)
			return ""
		}
		req.Header.Add("X-API-TOKEN", rest.token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Warning #59632: error on response.\n[ERR] - %s\nWill try again in %d seconds...\n", err, attemptDelay)
			time.Sleep(time.Duration(attemptDelay) * time.Second)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Warning #59635: error on response.\n[ERR] - %s\nWill try again in %d seconds...\n", err, attemptDelay)
			time.Sleep(time.Duration(attemptDelay) * time.Second)
			continue
		}
		// action is successful
		success = true
		break
	}
	if !success {
		log.Fatalf("Error #59638: Unable to complete HTTP '%s'.\nThis was attempted %d times with a delay of %d seconds in between each attempt.\n", verb, maxAttempts, attemptDelay)
		return ""
	}
	return string(body)
}

// Get - execute a Get request to an API end point
func (rest Rest) Get(path string) string {
	return rest.action("GET", path)
}

// Delete - execute a Delete request to an API end point
func (rest Rest) Delete(path string) string {
	return rest.action("DELETE", path)
}

// Post - execute a Post request to an API end point
func (rest Rest) Post(path string, jsonData []byte) string {
	var body []byte
	success := false
	for attempts := 0; attempts <= maxAttempts; attempts++ {
		url := rest.BaseURL + path
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error #32347: %s:%s", path, err)
			return ""
		}
		req.Header.Add("X-API-TOKEN", rest.token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Warning #59411: error on response.\n[ERR] - %s\nWill try again in %d seconds...\n", err, attemptDelay)
			time.Sleep(time.Duration(attemptDelay) * time.Second)
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Warning #59415: error on response.\n[ERR] - %s\nWill try again in %d seconds...\n", err, attemptDelay)
			time.Sleep(time.Duration(attemptDelay) * time.Second)
			continue
		}
		// post is successful
		success = true
		break
	}
	if !success {
		log.Fatalf("Error #59419: Unable to complete HTTP POST.\nThis was attempted %d times with a delay of %d seconds in between each attempt.\n", maxAttempts, attemptDelay)
		return ""
	}
	return string(body)
}

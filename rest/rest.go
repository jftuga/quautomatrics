package rest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

// action - generic function called by Get and Delete
func (rest Rest) action(verb, path string) string {
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
		log.Println("Error #59632: error on response.\n[ERR] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
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
		log.Println("Error #59411: error on response.\n[ERR] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

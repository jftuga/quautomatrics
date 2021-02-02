package rest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// https://stackoverflow.com/a/51453196

type Rest struct {
	token   string
	BaseURL string
}

func New(token, datacenter string) *Rest {
	url := strings.Replace(DataCenterTemplateURL, "{DATACENTER}", datacenter, -1)
	return &Rest{
		token:   token,
		BaseURL: url,
	}
}

func (rest Rest) Action(verb, path string) string {
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

func (rest Rest) Get(path string) string {
	return rest.Action("GET", path)
}

func (rest Rest) Delete(path string) string {
	return rest.Action("DELETE", path)
}

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

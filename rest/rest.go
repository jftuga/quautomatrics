package rest

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// https://stackoverflow.com/a/51453196

type Rest struct {
	token string
	BaseURL string
}

func New(token, datacenter string) *Rest {
	url := strings.Replace(DataCenterTemplateURL,"{DATACENTER}", datacenter, -1)
	return &Rest {
		token: token,
		BaseURL: url,
	}
}

func (rest Rest) Get(path string) string {
	url := rest.BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("%s:%s", path, err)
		return ""
	}
	req.Header.Add("X-API-TOKEN", rest.token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERR] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
package library

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/jftuga/quautomatrics/rest"
	"github.com/spf13/viper"
	"log"
)

type Connection struct {
	Token      string
	Datacenter string
	Rest       rest.Rest
}

type Library struct {
	Id   string
	Name string
	Conn Connection
}

// New - return a Library struct containing the name and connection information
func New(name string) *Library {
	id := generic(name, "/libraries", map[string]string {"name": "libraryName","id": "libraryId"})
	if len(id) == 0 {
		log.Fatalf("Error #47750: Library not found: %s\n", name)
	}
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	return &Library{
		Id:   id,
		Name: name,
		Conn: Connection{token, dc, *r},
	}
}

func (lib Library) GetLibraryMessage(name string) string {
	value := generic(name, fmt.Sprintf("/libraries/%s/messages?offset=0", lib.Id), map[string]string {"name": "description","id": "id"})
	if len(value) == 0 {
		log.Fatalf("Error #47752: Library Message not found: %s\n", name)
	}
	return value
}

// generic - make API call to to get the value of objectName
// keys must contain two keys: name, id
func generic(objectName, path string, keys map[string]string) string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	allObjects := r.Get(path)

	// validate JSON returned from API
	if json.Valid([]byte(allObjects)) != true {
		log.Fatalf("Error #20828: Invalid JSON returned from API:\n%s\n", allObjects)
	}

	result, _, _, err := jsonparser.Get([]byte(allObjects), "result")
	if err != nil {
		log.Fatalf("Error #20837: parsing JSON for key='result'\n%s\n", result)
	}

	// iterate through all mailing list entries and search for key["name"]
	hasList := false
	var name, objectValue string
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err = jsonparser.GetString(value, keys["name"])
		if err != nil {
			log.Fatalf("Error #20865: parsing JSON for objectName='%s', key='%s'\n%s\n", objectName, keys["name"], value)
		}
		if name == objectName { // mailing list is found, get its 'id'
			hasList = true
			objectValue, err = jsonparser.GetString(value, keys["id"])
			if err != nil {
				log.Fatalf("Error #20845: parsing JSON for key='%s'\n%s\n", keys["id"], value)
			}
			return
		}
	}, "elements")
	if err != nil {
		log.Fatalf("Error #20809: %s", err)
	}

	if hasList == false {
		log.Printf("Warning #20812: object does not exist: %s\n", objectName)
		return ""
	}
	return objectValue
}

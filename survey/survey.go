package survey

import (
	"encoding/json"
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
	id := generic(name, "/surveys", map[string]string {"name": "name","id": "id"})
	if len(id) == 0 {
		log.Fatalf("Error #49950: Survey not found: %s\n", name)
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

// generic - make API call to to get the value of objectName
// keys must contain two keys: name, id
func generic(objectName, path string, keys map[string]string) string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	allObjects := r.Get(path)

	// validate JSON returned from API
	if json.Valid([]byte(allObjects)) != true {
		log.Fatalf("Error #49928: Invalid JSON returned from API:\n%s\n", allObjects)
	}

	result, _, _, err := jsonparser.Get([]byte(allObjects), "result")
	if err != nil {
		log.Fatalf("Error #49937: parsing JSON for key='result'\n%s\n", result)
	}

	// iterate through all mailing list entries and search for key["name"]
	hasList := false
	var name, objectValue string
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err = jsonparser.GetString(value, keys["name"])
		if err != nil {
			log.Fatalf("Error #49965: parsing JSON for key='%s'\n%s\n", keys["name"], value)
		}
		if name == objectName { // mailing list is found, get its 'id'
			hasList = true
			objectValue, err = jsonparser.GetString(value, keys["id"])
			if err != nil {
				log.Fatalf("Error #49945: parsing JSON for key='%s'\n%s\n", keys["id"], value)
			}
			return
		}
	}, "elements")
	if err != nil {
		log.Fatalf("Error #49909: %s", err)
	}

	if hasList == false {
		log.Printf("Warning #49912: object does not exist: %s\n", objectName)
		return ""
	}
	return objectValue
}


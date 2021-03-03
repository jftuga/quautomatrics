package rest

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/spf13/viper"
	"log"
)

// Generic - make API call to to get the value of objectName
// keys must contain two keys: name, id
func Generic(objectName, path string, keys map[string]string) string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := New(token, dc)
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

// GenericMap - return all key/values from the API listed by path, keys
// Example: list all Surveys...
//          path="/surveys"
//          keys=map[string]string {"name": "name","id": "id"}
func GenericMap(path string, keys map[string]string) map[string]string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := New(token, dc)
	allObjects := r.Get(path)

	// validate JSON returned from API
	if json.Valid([]byte(allObjects)) != true {
		log.Fatalf("Error #32828: Invalid JSON returned from API:\n%s\n", allObjects)
	}

	result, _, _, err := jsonparser.Get([]byte(allObjects), "result")
	if err != nil {
		log.Fatalf("Error #32837: parsing JSON for key='result'\n%s\n", result)
	}

	var name, objectValue string
	returnedData := make(map[string]string)
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err = jsonparser.GetString(value, keys["name"])
		if err != nil {
			log.Fatalf("Error #32865: parsing JSON for name='%s', key='%s'\n%s\n%s\n", name, keys["name"], value, err)
		}
		objectValue, err = jsonparser.GetString(value, keys["id"])
		if err != nil {
			log.Fatalf("Error #32845: parsing JSON for key='%s'\n%s\n%s\n", keys["id"], value, err)
		}
		//fmt.Println(name, objectValue)
		returnedData[name] = objectValue
	}, "elements")
	if err != nil {
		log.Fatalf("Error #32809: %s", err)
	}
	return returnedData
}

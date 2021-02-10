package rest

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/jftuga/OLD--qualtrics_survey/rest"
	"github.com/spf13/viper"
	"log"
)

// GenericMap - return all key/values from the API listed by path, keys
// Example: list all Surveys...
//          path="/surveys"
//          keys=map[string]string {"name": "name","id": "id"}
func GenericMap(path string, keys map[string]string) map[string]string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
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
			log.Fatalf("Error #32865: parsing JSON for name='%s', key='%s'\n%s\n", name, keys["name"], value)
		}
		objectValue, err = jsonparser.GetString(value, keys["id"])
		if err != nil {
			log.Fatalf("Error #32845: parsing JSON for key='%s'\n%s\n", keys["id"], value)
		}
		//fmt.Println(name, objectValue)
		returnedData[name] = objectValue
	}, "elements")
	if err != nil {
		log.Fatalf("Error #32809: %s", err)
	}
	return returnedData
}

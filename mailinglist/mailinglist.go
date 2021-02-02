package mailingList

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

type MailingList struct {
	Id   string
	Name string
	Conn Connection
}

type Contact struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
}

func extractItem(value []byte, key string) string {
	item, err := jsonparser.GetString(value, key)
	if err != nil {
		log.Fatalf("Error #80211: unable to get '%s' from:\n%s\n", key, value)
	}
	return item
}

func New(name string) *MailingList {
	id := getMailingListID(name)
	if len(id) == 0 {
		log.Fatalf("Error #58025: Mailing list not found: %s\n", name)
	}
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	return &MailingList{
		Id:   id,
		Name: name,
		Conn: Connection{token, dc, *r},
	}
}

func getMailingListID(mailingListName string) string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	path := "mailinglists?offset=0"
	allMailingLists := r.Get(path)

	// validate JSON returned from API
	if json.Valid([]byte(allMailingLists)) != true {
		log.Fatalf("Error #45873: Invalid JSON returned from API:\n%s\n", allMailingLists)
	}

	result, _, _, err := jsonparser.Get([]byte(allMailingLists), "result")
	if err != nil {
		log.Fatalf("Error #46237: parsing JSON for key='result'\n%s\n", result)
	}

	// iterate through all mailing list entries and search for 'mailingListName'
	hasList := false
	var name, id string
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err = jsonparser.GetString(value, "name")
		if err != nil {
			log.Fatalf("Error #46885: parsing JSON for key='name'\n%s\n", value)
		}
		if name == mailingListName { // mailing list is found, get its 'id'
			hasList = true
			id, err = jsonparser.GetString(value, "id")
			if err != nil {
				log.Fatalf("Error #46005: parsing JSON for key='id'\n%s\n", value)
			}
			return
		}
	}, "elements")
	if err != nil {
		log.Fatalf("Error #38932: %s", err)
	}

	if hasList == false {
		log.Printf("Warning #46376: mailing list does not exist: %s\n", mailingListName)
		return ""
	}
	return id
}

func (mList MailingList) GetAllContacts() []Contact {
	path := fmt.Sprintf("/mailinglists/%s/contacts", mList.Id)
	request := mList.Conn.Rest.Get(path)
	result, _, _, err := jsonparser.Get([]byte(request), "result")
	if err != nil {
		log.Fatalf("Error #73021: parsing JSON for key='result'\n%s\n", result)
	}
	var allContacts []Contact
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id := extractItem(value, "id")
		email := extractItem(value, "email")
		first := extractItem(value, "firstName")
		last := extractItem(value, "lastName")
		con := new(Contact)
		con.Id = id
		con.Email = email
		con.LastName = last
		con.FirstName = first
		allContacts = append(allContacts, *con)
	}, "elements")
	if err != nil {
		log.Fatalf("Error #77502: %s", err)
	}

	//fmt.Println("result:", allContacts)
	return allContacts
}

func (mList MailingList) DeleteContact(contactId string) bool {
	path := fmt.Sprintf("/mailinglists/%s/contacts/%s", mList.Id, contactId)
	request := mList.Conn.Rest.Delete(path)
	meta, _, _, err := jsonparser.Get([]byte(request), "meta")
	if err != nil {
		log.Fatalf("Error #73639: parsing JSON for key='meta'\n%s\n", meta)
	}
	successfulDelete := false
	err = jsonparser.ObjectEach(meta, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		//fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
		if string(key) == "httpStatus" {
			if string(value) == "200 - OK" {
				successfulDelete = true
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error #77531: %s", err)
	}
	return successfulDelete
}

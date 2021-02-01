package mailingList

import (
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

func New(name, id string) *MailingList {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)
	return &MailingList{
		Id:   id,
		Name: name,
		Conn: Connection{token, dc, *r},
	}
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

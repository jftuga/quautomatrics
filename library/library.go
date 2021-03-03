package library

import (
	"fmt"
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
	id := rest.Generic(name, "/libraries", map[string]string {"name": "libraryName","id": "libraryId"})
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

// GetLibraryMessage - used by the listLibraries and createDistribution cli commands
// given the library id from lib.Id, return the library message
func (lib Library) GetLibraryMessage(name string) string {
	value := rest.Generic(name, fmt.Sprintf("/libraries/%s/messages?offset=0", lib.Id), map[string]string {"name": "description","id": "id"})
	if len(value) == 0 {
		log.Fatalf("Error #47752: Library Message not found: %s\n", name)
	}
	return value
}

// GetAllLibraryMessage - used by the listLibraries -M option
// given the library id from lib.Id, return all of library messages
func (lib Library) GetAllLibraryMessage() map[string]string {
	value := rest.GenericMap(fmt.Sprintf("/libraries/%s/messages?offset=0", lib.Id), map[string]string {"name": "description", "id": "id"})
	if len(value) == 0 {
		log.Fatalln("Error #47852: No Library Messages were found.")
	}
	return value
}

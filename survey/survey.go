package survey

import (
	"github.com/jftuga/quautomatrics/library"
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
	id := library.Generic(name, "/surveys", map[string]string {"name": "name","id": "id"})
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

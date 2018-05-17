package sync

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

var (
	err     error
	c       *mgo.Collection
	service Service
)

type Service struct {
	Name string
	// Port        string
	Address     string
	MongoServer string
	Redirect    string
	Partners    []string
	TLS         bool
}

func Serve(serviceVars Service) error {
	service = serviceVars
	if service.Name == "" {
		return errors.New("A service name must be provided")
	}

	session, err := mgo.Dial(service.MongoServer)
	if err != nil {
		return err
	}
	defer session.Close()
	c = session.DB(service.Name).C("master")

	http.HandleFunc("/in", in)
	http.HandleFunc("/out", out)
	http.HandleFunc("/forward", forward)
	http.HandleFunc("/back", back)
	if service.TLS {
		fmt.Println("Serving:", service.Name, "on port: 443")
		return http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	}
	fmt.Println("Serving:", service.Name, "on port: 80")
	return http.ListenAndServe(":80", nil)
	// fmt.Println("Serving:", service.Name, "on port:", service.Port)
	// return http.ListenAndServe(":"+service.Port, nil)

}
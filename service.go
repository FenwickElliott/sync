package sync

import (
	"gopkg.in/mgo.v2"
)

// Service struct containing all service specific data
type Service struct {
	Name        string `bson:"_id,omitempty"`
	Address     string `bson:"host"`
	Port        string `bson:"port"`
	Redirect    string `bson:"redirect"`
	MongoServer string `bson:"mongoserver"`
	TLS         bool   `bson:"tls"`
	collection  *mgo.Collection
}

func GetService(name string) (Service, error) {
	session, err := mgo.Dial("cookies.fenwickelliott.io")
	check(err)
	servicesCollection := session.DB("services").C("services")
	res := Service{}
	err = servicesCollection.FindId(name).One(&res)
	if err != nil {
		return Service{}, err
	}
	return res, nil
}

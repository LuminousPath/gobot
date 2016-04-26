package ohayou

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var dbAddress string = "107170.com"

type Test struct {
	Name    string
	Ohayous string
}

func getOhayous(nick string) string {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	result := Test{}

	err = q.Find(bson.M{"username": nick}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result.Ohayous
}

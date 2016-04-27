package ohayou

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbAddress string = "localhost"

// returns a user's document as a User{} type
func getUser(nick string) User {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	result := User{}

	err = q.Find(bson.M{"username": nick}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// saves the new amount of ohayous after a user has ohayou'd
func saveOhayous(nick string, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	save := bson.M{"$set": bson.M{"ohayous": amt, "last": time.Now()}}

	err = q.Update(bson.M{"username": nick}, save)
	if err != nil {
		log.Fatal(err)
	}
}

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
		log.Println("getUser:" + err.Error())
	}

	return result
}

func newUser(nick string, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	save := bson.M{"username": nick,
		"last":          t.In(est),
		"ohayous":       amt,
		"cumOhayous":    amt,
		"add":           0,
		"timesOhayoued": 1,
		"items":         bson.M{},
		"itemsMultiply": bson.M{}}

	err = q.Insert(save)
	if err != nil {
		log.Println("newUser:" + err.Error())
	}
}

// saves the new amount of ohayous after a user has ohayou'd
func saveOhayous(user User, amt int) int {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	totalAmt := user.Ohayous + amt

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	save := bson.M{"$set": bson.M{"ohayous": totalAmt,
		"last":          t.In(est),
		"cumOhayous":    user.CumOhayous + amt,
		"timesOhayoued": user.TimesOhayoued + 1}}

	err = q.Update(bson.M{"username": user.Username}, save)
	if err != nil {
		log.Println("saveOhayous:" + err.Error())
	}

	return totalAmt

}

// saves an item when it's purchased
func saveItem(user User, item string, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB("ircbot").C("ohayou")

	save := bson.M{"$inc": bson.M{"items." + item: amt}}

	err = q.Update(bson.M{"username": nick}, save)
	if err != nil {
		log.Println("saveItem:" + err.Error())
	}

}

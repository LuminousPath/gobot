package ohayou

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// db consts
var (
	dbAddress string = "localhost"
	dbName    string = "ircbot"
	ohyCol    string = "ohayou"
	itemCol   string = "items"
)

// returns a user's document as a User{} type
func getUser(nick string) *User {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	result := User{}

	err = q.Find(bson.M{"username": nick}).One(&result)
	if err != nil {
		log.Println("getUser: " + err.Error())
	}

	return &result
}

func newUser(nick string, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	save := bson.M{
		"username":      nick,
		"last":          t.In(est),
		"ohayous":       amt,
		"cumOhayous":    amt,
		"add":           0,
		"timesOhayoued": 1,
		"items":         bson.M{},
		"itemMultiply":  bson.M{}}

	err = q.Insert(save)
	if err != nil {
		log.Println("newUser: " + err.Error())
	}
}

// saves the new amount of ohayous after a user has ohayou'd
func saveOhayous(user *User, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	save := bson.M{"$set": bson.M{
		"ohayous":       amt,
		"last":          t.In(est),
		"cumOhayous":    user.CumOhayous + amt,
		"timesOhayoued": user.TimesOhayoued + 1}}

	err = q.Update(bson.M{"username": user.Username}, save)
	if err != nil {
		log.Println("saveOhayous: " + err.Error())
	}
}

// get all data for item
func getItem(item string) *Item {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(itemCol)

	result := Item{}

	err = q.Find(bson.M{"name": item}).One(&result)
	if err != nil {
		log.Println("getItem: " + err.Error())
	}

	return &result
}

// saves an item when it's purchased
func saveItem(user *User, item string, amt int) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{"items." + item: amt}}

	err = q.Update(bson.M{"username": item}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}

}

func getCategories() string {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(itemCol)

	var result []string

	err = q.Find(nil).Distinct("category", &result)
	if err != nil {
		log.Println("getCategories: " + err.Error())
	}

	return strings.Join(append(result), ", ")
}

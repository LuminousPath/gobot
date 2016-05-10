package ohayou

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// returns a user's document as a User{} type
func getUser(nick string) bool {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	err = q.Find(bson.M{"username": nick}).One(&USER)
	if err != nil {
		log.Println("getUser: " + err.Error())
		return false
	}
	return true
}

func newUser(nick string, amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)
	t = time.Now()

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
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	t = time.Now()

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
func getItem(item string) bool {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	err = q.Find(bson.M{"name": item}).One(&ITEM)
	if err != nil {
		log.Println("getItem: " + err.Error())
		return false
	}
	return true
}

// saves an item when it's purchased
func saveItem(user *User, item string, amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{"items." + item: amt}}

	err = q.Update(bson.M{"username": user.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

func setLastUsed(user *User, item string) {
	s := session.Copy()
	defer s.Close()
	t = time.Now()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"lastUsed." + item: t.In(est)}}

	err = q.Update(bson.M{"username": user.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

// returns a concatenated string of all categories
func listCategories() string {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	var result []string

	err = q.Find(nil).Distinct("category", &result)
	if err != nil {
		log.Println("getCategories: " + err.Error())
	}

	// cate, gor, ies
	return strings.Join(append(result), ", ")
}

// returns basic information about all items in a category
func getCategory(name string) []string {
	s := session.Copy()
	defer session.Close()
	q := s.DB(dbName).C(itemCol)

	var result []Item

	err = q.Find(bson.M{"category": name}).Sort("price").All(&result)
	if err != nil {
		log.Println("getCategory: " + err.Error())
	}

	items := make([]string, len(result))

	// just get the necessary info
	for j, item := range result {
		items[j] = fmt.Sprintf("%s - %d ohayous - %s",
			item.Name, item.Price, item.Desc)
	}

	return items
}

func resetLast(user *User) {
	s := session.Copy()
	defer session.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"last": 0}}

	err = q.Update(bson.M{"username": user.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

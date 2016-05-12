package ohayou

import (
	"fmt"
	"log"
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

	save = bson.M{
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
func (u *User) saveOhayous(amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)
	t = time.Now()

	save = bson.M{"$set": bson.M{
		"ohayous":       amt,
		"last":          t.In(est),
		"cumOhayous":    u.CumOhayous + amt,
		"timesOhayoued": u.TimesOhayoued + 1}}

	err = q.Update(bson.M{"username": u.Username}, save)
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

// mostly used for simple incrementation for an item, ie bottle
func (u *User) saveItem(itm string, amt int) {
	if getItem(itm) {
		s := session.Copy()
		defer s.Close()
		q := s.DB(dbName).C(ohyCol)

		if ITEM.Multiplies != "" {
			save = bson.M{"$inc": bson.M{
				"ohayous":                         -ITEM.Price * amt,
				"add":                             ITEM.Add * amt,
				"items." + itm:                    amt,
				"itemMultiply." + ITEM.Multiplies: ITEM.Multiply}}
		} else {
			save = bson.M{"$inc": bson.M{
				"ohayous":      -ITEM.Price * amt,
				"add":          ITEM.Add * amt,
				"items." + itm: amt}}
		}

		err = q.Update(bson.M{"username": u.Username}, save)
		if err != nil {
			log.Println("saveItem: " + err.Error())
		}
	}
}

func (u *User) setLastUsed(item string) {
	s := session.Copy()
	defer s.Close()
	t = time.Now()
	q := s.DB(dbName).C(ohyCol)

	save = bson.M{"$set": bson.M{"lastUsed." + item: t.In(est)}}

	err = q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

// returns a concatenated string of all categories
func listCategories() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	err = q.Find(nil).Distinct("category", &itemCats)
	if err != nil {
		log.Println("getCategories: " + err.Error())
	}
}

// returns basic information about all items in a category
func getCategory(name string) []string {
	s := session.Copy()
	defer s.Close()
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

func (u *User) resetLast() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save = bson.M{"$set": bson.M{"last": 0}}

	err = q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

func (u *User) savePin(pn int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save = bson.M{"$set": bson.M{"pin": pn}}

	err = q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

func getTop() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save = bson.M{"username": 1, "ohayous": 1}

	err = q.Find(nil).Select(save).Sort("-ohayous").Limit(5).Iter().All(&top)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

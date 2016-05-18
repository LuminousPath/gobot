package ohayou

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// db consts
const (
	dbAddress string = "localhost"
	dbName    string = "ircbot"
	ohyCol    string = "ohayou"
	itemCol   string = "items"
)

// returns a user's document as a User{} type
func GetUser(nick string) (User, bool) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	user := User{}

	err := q.Find(bson.M{"username": nick}).One(&user)
	if err != nil {
		log.Println("getUser: " + err.Error())
		return user, false
	}
	return user, true
}

func NewUser(nick string, amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{
		"username":      nick,
		"last":          time.Now().In(est),
		"ohayous":       amt,
		"cumOhayous":    amt,
		"add":           0,
		"timesOhayoued": 1,
		"items":         bson.M{},
		"itemMultiply":  bson.M{}}

	err := q.Insert(save)
	if err != nil {
		log.Println("newUser: " + err.Error())
	}
}

// saves the new amount of ohayous after a user has ohayou'd
func (u *User) SaveOhayous(amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{
		"ohayous":       amt,
		"last":          time.Now().In(est),
		"cumOhayous":    u.CumOhayous + amt,
		"timesOhayoued": u.TimesOhayoued + 1}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveOhayous: " + err.Error())
	}
}

// get all data for item
func GetItem(itm string) (Item, bool) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	item := Item{}

	err := q.Find(bson.M{"name": itm}).One(&item)
	if err != nil {
		log.Println("getItem: " + err.Error())
		return item, false
	}
	return item, true
}

// mostly used for simple incrementation for an item, ie bottle
func (u *User) SaveItem(item Item, amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	var save bson.M

	if item.Multiplies != "" {
		save = bson.M{"$inc": bson.M{
			"ohayous":                         -item.Price * amt,
			"add":                             item.Add * amt,
			"items." + item.Name:              amt,
			"itemMultiply." + item.Multiplies: item.Multiply}}
	} else {
		save = bson.M{"$inc": bson.M{
			"ohayous":            -item.Price * amt,
			"add":                item.Add * amt,
			"items." + item.Name: amt}}
	}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}
}

func (u *User) SetLastUsed(item string) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"lastUsed." + item: time.Now().In(est)}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

// returns a concatenated string of all categories
func setCategories() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	err := q.Find(nil).Distinct("category", &itemCtgs)
	if err != nil {
		log.Println("setCategories: " + err.Error())
	}
}

// returns basic information about all items in a category
func ItemCategory(name string) []string {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(itemCol)

	var result []Item

	err := q.Find(bson.M{"category": name}).Sort("price").All(&result)
	if err != nil {
		log.Println("getCategory: " + err.Error())
	}

	items := make([]string, len(result))

	// just get the necessary info
	for i, item := range result {
		items[i] = fmt.Sprintf("%s - %d ohayous - %s",
			item.Name, item.Price, item.Desc)
	}

	return items
}

func (u *User) ResetLast() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"last": 0}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

func (u *User) SavePin(pin int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"pin": pin}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}

func Top() string {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	var top []UserOhayous
	var result string

	query := bson.M{"username": 1, "ohayous": 1}

	err := q.Find(nil).Select(query).Sort("-ohayous").Limit(5).Iter().All(&top)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}

	for i := range top {
		result += fmt.Sprintf("%s: %d, ", top[i].Username, top[i].Ohayous)
	}
	// trim trailing ", "
	return result[:len(top)-2]
}

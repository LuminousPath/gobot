package ohayou

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func useItem(nick, nickRaw, itemName, useOn string) string {
	user := getUser(nick)
	item := getItem(itemName)

	if user.TimesOhayoued == 0 {
		return "You don't have any items because you've never ohayoued!" +
			" Get your first ration by typing " + p + "ohayou"
	}

	if item.Name == "" {
		return itemName + " isn't an item. Type " + p + "items to look and what's" +
			" available, and " + p + "inventory to see what items you have."
	}

	if user.Items[item.Name] == 0 {
		return "You don't any of that."
	}

	if !item.Useable {
		return item.Name + " is passive. It can't be used"
	}

	if item.Consume {
		consumeItem(nick, item.Name)
	}

	if item.HasFunction != "" {
		// uh oh
	}

	return nickRaw + " " + strings.Replace(item.Effect, "%s", useOn, -1)
}

func consumeItem(nick, itemName string) {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{
		"items." + itemName: -1}}

	err = q.Update(bson.M{"username": nick}, save)
	if err != nil {
		log.Println("consumeItem: " + err.Error())
	}
}

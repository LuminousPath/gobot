package ohayou

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func useItem(nick, nickRaw, itemName, useOn string) string {
	if !getUser(nick) {
		return "You don't have any items because you've never ohayoued!" +
			" Get your first ration by typing " + p + "ohayou"
	}

	if !getItem(itemName) {
		return itemName + " isn't an item. Type " + p + "items to look and what's" +
			" available, and " + p + "inventory to see what items you have."
	}

	if USER.Items[ITEM.Name] == 0 {
		return "You don't any of that."
	}

	if !ITEM.Useable {
		return ITEM.Name + " is passive. It can't be used"
	}

	if ITEM.Consume {
		go consumeItem(USER, ITEM.Name)
	}

	if ITEM.HasFunction != "" {
		doIt := itemFuncs[ITEM.HasFunction]
		extra = doIt(USER, ITEM.Name)
	}

	return nickRaw + " " + strings.Replace(ITEM.Effect, "%s", useOn, -1) + extra
}

func consumeItem(user *User, itemName string) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{
		"items." + itemName: -1}}

	err = q.Update(bson.M{"username": USER.Username}, save)
	if err != nil {
		log.Println("consumeItem: " + err.Error())
	}
}

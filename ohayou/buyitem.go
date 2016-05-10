package ohayou

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func buyItem(nick, item string, amt int) string {
	if !getUser(nick) {
		return "You haven't ohayoued yet! Type " + p + "ohayou to get your first ration."
	}

	// item not found
	if !getItem(item) {
		return "I don't have that in stock."
	}

	// item cannot be purchased
	if !ITEM.Purchase {
		return "That's not for sale."
	}

	if USER.Ohayous < ITEM.Price*amt {
		return "You can't afford that."
	}

	// user is already at the limit for that item
	if ITEM.Limit > 0 && USER.Items[item] >= ITEM.Limit {
		return fmt.Sprintf("You can't purchase any more of that. You can only have"+
			" %d %s", ITEM.Limit, item)
	}

	// this purchase (presumeably batch purchase) would push them over the limit
	if ITEM.Limit > 0 && USER.Items[item]+amt > ITEM.Limit {
		return fmt.Sprintf("You can't purchase that much. You can only have"+
			" %d %s", ITEM.Limit, item)
	}

	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{}

	// check if this item multiplies another item
	if ITEM.Multiplies != "" {
		save = bson.M{"$inc": bson.M{
			"ohayous":                         -ITEM.Price * amt,
			"add":                             ITEM.Add * amt,
			"items." + item:                   amt,
			"itemMultiply." + ITEM.Multiplies: ITEM.Multiply}}
	} else {
		save = bson.M{"$inc": bson.M{
			"ohayous":       -ITEM.Price * amt,
			"add":           ITEM.Add * amt,
			"items." + item: amt}}
	}

	q.Update(bson.M{"username": nick}, save)
	if err != nil {
		log.Println("getCategory: " + err.Error())
	}

	if amt > 1 {
		return fmt.Sprintf("You purchased %d %ss for %d ohayous. "+
			"You have %d ohayous left.",
			amt, item, ITEM.Price*amt, USER.Ohayous-(ITEM.Price*amt))
	} else {
		return fmt.Sprintf("You purchased %d %s for %d ohayous. "+
			"You have %d ohayous left.",
			amt, item, ITEM.Price*amt, USER.Ohayous-(ITEM.Price*amt))
	}
}

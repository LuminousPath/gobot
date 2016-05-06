package ohayou

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func buyItem(nick, item string, amt int) string {
	user := getUser(nick)
	itemData := getItem(item)

	if user.Username == "" {
		return "You haven't ohayoued yet! Type " + p + "ohayou to get your first ration."
	}

	// item not found
	if itemData.Name == "" {
		return "I don't have that in stock."
	}

	// item cannot be purchased
	if !itemData.Purchase {
		return "That's not for sale."
	}

	if user.Ohayous < itemData.Price*amt {
		return "You can't afford that."
	}

	// user is already at the limit for that item
	if itemData.Limit > 0 && user.Items[item] >= itemData.Limit {
		return fmt.Sprintf("You can't purchase any more of that. You can only have"+
			" %d %s", itemData.Limit, item)
	}

	// this purchase (presumeably batch purchase) would push them over the limit
	if itemData.Limit > 0 && user.Items[item]+amt > itemData.Limit {
		return fmt.Sprintf("You can't purchase that much. You can only have"+
			" %d %s", itemData.Limit, item)
	}

	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	q := session.DB(dbName).C(ohyCol)

	save := bson.M{}

	// check if this item multiplies another item
	if itemData.Multiplies != "" {
		save = bson.M{"$inc": bson.M{
			"ohayous":                             -itemData.Price * amt,
			"add":                                 itemData.Add * amt,
			"items." + item:                       amt,
			"itemMultiply." + itemData.Multiplies: itemData.Multiply}}
	} else {
		save = bson.M{"$inc": bson.M{
			"ohayous":       -itemData.Price * amt,
			"add":           itemData.Add * amt,
			"items." + item: amt}}
	}

	q.Update(bson.M{"username": nick}, save)
	if err != nil {
		log.Println("getCategory: " + err.Error())
	}

	if amt > 1 {
		return fmt.Sprintf("You purchased %d %ss for %d ohayous. "+
			"You have %d ohayous left.",
			amt, item, itemData.Price*amt, user.Ohayous-(itemData.Price*amt))
	} else {
		return fmt.Sprintf("You purchased %d %s for %d ohayous. "+
			"You have %d ohayous left.",
			amt, item, itemData.Price*amt, user.Ohayous-(itemData.Price*amt))
	}
}

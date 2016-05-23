package ohayou

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func (u *User) Use(nick, itm, on string) string {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	item, ok := GetItem(itm)
	var extra string
	// item not found
	if !ok {
		return itm + " isn't an item. Type " + p + "items to look at what's" +
			" available, and " + p + "inventory to see what items you have."
	}

	if u.Items[item.Name] == 0 {
		return "You don't have any of that."
	}

	if !item.Useable {
		return item.Name + " is a passive item. It can't be used."
	}

	if item.Consume {
		u.ConsumeItem(item.Name)
	}

	if item.Name == "vault" && u.Vault.Installed {
		return ""
	}

	if item.HasFunction != "" {
		do := itemFuncs[item.HasFunction]
		extra = do(u, item.Name)
	}

	if canAdoptCat && item.HasFunction == "adoptCat" {
		return nick + " offers the cat a " + item.Name + "..."
	}

	return nick + " " + strings.Replace(item.Effect, "%s", on, -1) + extra
}

func (u *User) ConsumeItem(itm string) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{
		"items." + itm: -1}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("consumeItem: " + err.Error())
	}
}

package ohayou

import (
	"fmt"
	"time"
)

var (
	typeResponse   string
	ohayous        int // new ohayous
	itemOhayous    int // extra ohayous given to user from items
	itemMultiplier int // any item multipliers
	totalOhayous   int // added all up

	adj = [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous",
		"Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}
)

// main function that distributes ohayous
func newOhayou(nick string) string {
	ohayous = randNum(0, 6)
	switch ohayous {
	case 0:
		typeResponse = "But not good enough. You get 0 ohayous."
	case 1:
		typeResponse = "You get 1 ohayou."
	case 6:
		typeResponse = "Wow! You get 6 ohayous!"
	default:
		typeResponse = fmt.Sprintf("You get %d ohayous!", ohayous)
	}
	// get their data
	last = time.Now()
	// dont allow ohayou if they have ohayou'd today
	if !getUser(lowNick) {
		newUser(lowNick, ohayous)
		return "Congratulations on your first ohayou " + nick + "!!! " +
			typeResponse + " Type " + p + "help ohayou if you don't know what " +
			"this is."
	} else if USER.Last.In(est).Format("20060102") >= last.In(est).Format("20060102") {
		return "You already got your ohayou ration today, " + nick + "."
	} else {
		itemMultiplier = 1
		itemOhayous = 0
		for itm, amt := range USER.Items {
			if USER.Items[itm] == 0 {
				continue
			}
			// check if user has item(s) that multiply another item
			if USER.ItemMultiply[itm] != 0 {
				itemMultiplier = USER.ItemMultiply[itm]
			}
			getItem(itm)
			itemOhayous += (ITEM.Add * amt) * itemMultiplier
		}
		if USER.Ohayous <= 0 {
			totalOhayous = ohayous + itemOhayous
		} else {
			totalOhayous = USER.Ohayous + ohayous + itemOhayous
		}
		// store it
		USER.saveOhayous(totalOhayous)
		if itemOhayous == 0 {
			return fmt.Sprintf("%s ohayou %s!!! %s You have %d ohayous.",
				adj[randNum(0, 10)], nick, typeResponse, totalOhayous)
		} else {
			return fmt.Sprintf("%s ohayou %s!!! %s Your items increased "+
				"that to %d. You have %d ohayous.",
				adj[randNum(0, 10)], nick, typeResponse,
				ohayous+itemOhayous, totalOhayous)
		}
	}
}

package ohayou

import (
	"fmt"
	"strings"
	"time"
)

var adj = [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous",
	"Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}

// main function that distributes ohayous
func Ohayou(nick string) string {
	ohayous := randNum(0, 6)
	var typeResponse string
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
	user, ok := GetUser(nick)
	// dont allow ohayou if they have ohayou'd today
	if !ok {
		NewUser(nick, ohayous)
		return "Congratulations on your first ohayou " + nick + "!!! " +
			typeResponse + " Type " + p + "help ohayou if you don't know what " +
			"this is."
	} else if user.Last.In(est).Format("20060102") >= time.Now().In(est).Format("20060102") {
		return "You already got your ohayou ration today, " + nick + "."
	} else {
		itemOhayous, totalOhayous := 0, 0
		for itm, amt := range user.Items {
			itemMultiplier := 1
			// check if user has item(s) that multiply another item
			if user.ItemMultiply[itm] != 0 {
				itemMultiplier = user.ItemMultiply[itm]
			}
			item, _ := GetItem(itm)
			itemOhayous += (item.Add * amt) * itemMultiplier
		}
		if user.Ohayous <= 0 {
			totalOhayous = ohayous + itemOhayous
		} else {
			totalOhayous = user.Ohayous + ohayous + itemOhayous
		}

		if doubleOhayou {
			var ohayouError int = randNum(13, 21)
			totalOhayous = int(totalOhayous * (ohayouError / 10))
		}
		// store it
		user.SaveOhayous(totalOhayous)
		if doubleOhayou {
			return fmt.Sprintf("ERROR C0045: <%s> <%d> OHAYOUS DISPENSED",
				strings.ToUpper(nick), totalOhayous)
		}
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

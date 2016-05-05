package ohayou

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mferrera/go-ircevent"
)

func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}

	return false
}

func randNum(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// main function that distributes ohayous
func newOhayou(nick string) string {
	adj := [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous",
		"Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}

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
	user := getUser(strings.ToLower(nick))

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	// dont allow ohayou if they have ohayou'd today
	if user.Last.Format("20060102") >= t.In(est).Format("20060102") {
		return "You already got your ohayou ration today, " + nick
	} else if user.TimesOhayoued == 0 {
		newUser(strings.ToLower(nick), ohayous)

		return "Congratulations on your first ohayou " + nick + "!!! " +
			typeResponse + " Type .ohayouhelp if you don't know what this is."
	} else {
		var itemOhayous int

		for item, amt := range user.Items {
			itemMultiplier := 1

			// check if user has item(s) that multiply another item
			if user.ItemMultiply[item] != 0 {
				itemMultiplier = user.ItemMultiply[item]
			}

			itemData := getItem(item)
			itemOhayous += (itemData.Add * amt) * itemMultiplier
		}

		totalOhayous := user.Ohayous + ohayous + itemOhayous

		// store it
		saveOhayous(user, totalOhayous)

		if itemOhayous == 0 {
			return fmt.Sprintf("%s ohayou %s!!! %s You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse, totalOhayous)
		} else {
			return fmt.Sprintf("%s ohayou %s!!! %s Your items increased "+
				"that to %d. You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse,
				ohayous+itemOhayous, totalOhayous)
		}
	}
}

func Run(bot *irc.Connection, p, cmd, channel, nick string, word []string, admin bool) {
	say := bot.Privmsg

	// main command to acquire new ohayous
	if cmd == p+"ohayou" && !hasArgs(word) {
		say(channel, newOhayou(nick))
	}

	// respond to channel with how many ohayous X has
	if cmd == p+"ohayou" && hasArgs(word) {
		user := getUser(strings.ToLower(word[1]))

		if user.Username != "" {
			say(channel, fmt.Sprintf("%s has %d ohayous.", word[1], user.Ohayous))
		} else {
			say(channel, word[1]+" hasn't ohayoued yet!")
		}
	}

	if cmd == p+"buy" && !hasArgs(word) {
		say(channel, "Usage: "+p+"buy <item> will buy you one <item>."+
			p+"buy <item> 3 will buy you 3 of <item>, if you can afford it.")
	} else if cmd == p+"buy" && hasArgs(word) {
		// if a purchase quantity is given
		if len(word) > 2 {
			// try to convert it to an integer
			amt, err := strconv.Atoi(word[2])
			if err != nil {
				say(channel, "You didn't give a valid quantity. Usage: "+p+
					"buy <item> will buy you one <item>. "+p+"buy <item>"+
					" 3 will buy you 3 of <item>, if you can afford it.")
			} else {
				say(channel, buyItem(strings.ToLower(nick),
					strings.ToLower(word[1]),
					amt))
			}
		} else {
			say(channel, buyItem(strings.ToLower(nick),
				strings.ToLower(word[1]),
				1))
		}
	}

	// just shows how to use .items and lists item categories
	if cmd == p+"items" && !hasArgs(word) {
		say(channel, "Type "+p+"item <category> to get a list of items by "+
			"category. Categories: "+listCategories())
	}

	// PMs all items in a category
	if cmd == p+"items" && hasArgs(word) {
		data := getCategory(strings.ToLower(word[1]))

		for _, item := range data {
			say(nick, item)
		}
	}

	// returns information about an item
	if cmd == p+"item" && hasArgs(word) {
		data := getItem(strings.ToLower(word[1]))

		if data.Name != "" {
			say(channel, fmt.Sprintf("%s: %s - Price: %d ohayous",
				data.Name, data.Desc, data.Price))
		} else {
			say(channel, "I don't carry that.")
		}
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		user := getUser(strings.ToLower(nick))

		if user.TimesOhayoued == 0 {
			say(channel, "You haven't ohayoued yet! Type .ohayou to "+
				"get your first ration.")
		} else if len(user.Items) > 0 {
			inv := "You have "

			for item, amt := range user.Items {
				if amt > 1 {
					inv += fmt.Sprintf("%d %ss ", amt, item)
				} else {
					inv += fmt.Sprintf("%d %s ", amt, item)
				}
			}

			say(nick, inv)
		} else {
			say(nick, "You don't have any items yet. Keep saving!")
		}
	}
}

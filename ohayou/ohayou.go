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
	return min + rand.Intn(max-min)
}

// main function that distributes ohayous
func newOhayou(nick string) string {
	adj := [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous", "Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}

	ohayous := randNum(0, 7)
	var typeResponse string

	switch ohayous {
	case 0:
		typeResponse = "But not good enough. You get 0 ohayous."
	case 1:
		typeResponse = "You get 1 ohayou."
	case 6:
		typeResponse = "Wow! You get 6 ohayous!"
	default:
		typeResponse = "You get " + strconv.Itoa(ohayous) + " ohayous!"
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
		return fmt.Sprintf("You already got your ohayou ration today, "+
			"%s.", nick)
	}

	if user.TimesOhayoued == 0 {
		newUser(strings.ToLower(nick), ohayous)
		return fmt.Sprintf("Congratulations on your first ohayou %s!!! "+
			"%s Type .ohayouhelp if you don't know what this is.",
			nick, typeResponse)
	} else {
		var itemOhayous int

		for item, amt := range user.Items {
			itemMult := 1

			// if an item is multiplied by another item
			if user.ItemMultiply[item] != 0 {
				itemMult = user.ItemMultiply[item]
			}

			itemData := getItem(item)
			itemOhayous += (itemData.Add * amt) * itemMult
		}

		totalOhayous := user.Ohayous + ohayous + itemOhayous

		saveOhayous(user, totalOhayous)

		if itemOhayous == 0 {
			return fmt.Sprintf("%s ohayou %s!!! %s You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse, totalOhayous)
		} else {
			return fmt.Sprintf("%s ohayou %s!!! %s Your items increased "+
				"that to %d. You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse, itemOhayous, totalOhayous)
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
		data := getUser(strings.ToLower(word[1]))
		if data.Username != "" {
			say(channel, fmt.Sprintf("%s has %d ohayous.", word[1], data.Ohayous))
		} else {
			say(channel, word[1]+" hasn't ohayoued yet!")
		}
	}

	if cmd == p+"items" {
		say(channel, "Type "+p+"item <category> to get a list of items by "+
			"category. Categories: "+getCategories())
	}

	// returns information about an item
	if cmd == p+"item" && hasArgs(word) {
		data := getItem(strings.ToLower(word[1]))
		if data.Name != "" {
			say(channel, fmt.Sprintf("%s: %s - Price: %d ohayous", data.Name, data.Desc, data.Price))
		} else {
			say(channel, "I don't carry that.")
		}
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		data := getUser(strings.ToLower(nick))
		inv := "You have "

		for item, qty := range data.Items {
			if qty > 1 {
				inv += strconv.Itoa(qty) + " " + item + "s "
			} else {
				inv += strconv.Itoa(qty) + " " + item + " "
			}
		}

		say(nick, inv)
	}
}
